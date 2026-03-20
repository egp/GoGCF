// cf/api_types.go v3

package cf

import "math/big"

type PQTermKind int

const (
	PQValue PQTermKind = iota
	PQEOF
)

type PQTerm struct {
	Kind PQTermKind
	P    *big.Int
	Q    *big.Int
}

func (t PQTerm) IsValue() bool {
	return t.Kind == PQValue
}

func (t PQTerm) IsEOF() bool {
	return t.Kind == PQEOF
}

type RCFTermKind int

const (
	RCFValue RCFTermKind = iota
	RCFEOF
)

type RCFTerm struct {
	Kind  RCFTermKind
	Value *big.Int
}

func (t RCFTerm) IsValue() bool {
	return t.Kind == RCFValue
}

func (t RCFTerm) IsEOF() bool {
	return t.Kind == RCFEOF
}

func (t RCFTerm) BigInt() (*big.Int, bool) {
	if t.Kind != RCFValue || t.Value == nil {
		return nil, false
	}
	return t.Value, true
}

type Rational struct {
	Num *big.Int
	Den *big.Int
}

type BoundKind int

const (
	BoundFinite BoundKind = iota
	BoundNegInf
	BoundPosInf
)

type Bound struct {
	Kind   BoundKind
	Value  Rational
	Closed bool
}

type Range struct {
	Lo      Bound
	Hi      Bound
	Outside bool
}

func (r Range) IsInside() bool {
	return !r.Outside
}

func (r Range) IsOutside() bool {
	return r.Outside
}

func (r Range) IsExact() bool {
	if r.Outside {
		return false
	}
	if r.Lo.Kind != BoundFinite || r.Hi.Kind != BoundFinite {
		return false
	}
	if !r.Lo.Closed || !r.Hi.Closed {
		return false
	}
	if r.Lo.Value.Num == nil || r.Lo.Value.Den == nil || r.Hi.Value.Num == nil || r.Hi.Value.Den == nil {
		return false
	}

	left := new(big.Int).Mul(r.Lo.Value.Num, r.Hi.Value.Den)
	right := new(big.Int).Mul(r.Hi.Value.Num, r.Lo.Value.Den)
	return left.Cmp(right) == 0
}

type GCFSource interface {
	NextPQ() (PQTerm, GCFSource, error)
}

type GCF interface {
	Next() (RCFTerm, GCF, error)
	Range() Range
	Take(n int) (GCF, error)
	Convergent() (Rational, error)
}

func (r Range) Cmp(other Range) int {
	rValid := rangeWellFormed(r)
	otherValid := rangeWellFormed(other)

	switch {
	case rValid && !otherValid:
		return -1
	case !rValid && otherValid:
		return 1
	case !rValid && !otherValid:
		return 0
	}

	if rangeEquivalent(r, other) {
		return 0
	}

	rSubOther := rangeAllowedSubsetOf(r, other)
	otherSubR := rangeAllowedSubsetOf(other, r)

	switch {
	case rSubOther && !otherSubR:
		return -1
	case otherSubR && !rSubOther:
		return 1
	}

	rClass := rangeQualityClass(r)
	otherClass := rangeQualityClass(other)
	if rClass != otherClass {
		if rClass < otherClass {
			return -1
		}
		return 1
	}

	switch {
	case !r.Outside && !other.Outside:
		return compareInsideWidth(r, other)
	case r.Outside && other.Outside:
		return compareOutsideWidth(r, other)
	default:
		return 0
	}
}

func rangeEquivalent(a, b Range) bool {
	return a.Outside == b.Outside &&
		boundEqual(a.Lo, b.Lo) &&
		boundEqual(a.Hi, b.Hi)
}

func rangeAllowedSubsetOf(a, b Range) bool {
	switch {
	case !a.Outside && !b.Outside:
		return insideSubsetOf(a, b)
	case a.Outside && b.Outside:
		return outsideSubsetOf(a, b)
	case !a.Outside && b.Outside:
		return insideSubsetOfOutside(a, b)
	case a.Outside && !b.Outside:
		return outsideSubsetOfInside(a, b)
	default:
		return false
	}
}

func insideSubsetOf(a, b Range) bool {
	return boundLE(b.Lo, a.Lo) && boundLE(a.Hi, b.Hi)
}

func outsideSubsetOf(a, b Range) bool {
	return boundLE(a.Lo, b.Lo) && boundLE(b.Hi, a.Hi)
}

func insideSubsetOfOutside(inside, outside Range) bool {
	return intervalDisjoint(inside.Lo, inside.Hi, outside.Lo, outside.Hi)
}

func outsideSubsetOfInside(outside, inside Range) bool {
	return isInsideFullLine(inside)
}

func isInsideFullLine(r Range) bool {
	if r.Outside {
		return false
	}
	return r.Lo.Kind == BoundNegInf && r.Hi.Kind == BoundPosInf
}

func intervalDisjoint(aLo, aHi, bLo, bHi Bound) bool {
	return boundLT(aHi, bLo) || boundLT(bHi, aLo)
}

func rangeQualityClass(r Range) int {
	if r.IsExact() {
		return 0
	}
	if !r.Outside {
		return 1
	}
	return 2
}

func compareInsideWidth(a, b Range) int {
	spanCmp, ok := compareFiniteSpan(a, b)
	if ok {
		return spanCmp
	}
	return 0
}

func compareOutsideWidth(a, b Range) int {
	spanCmp, ok := compareFiniteSpan(a, b)
	if ok {
		return -spanCmp
	}
	return 0
}

func compareFiniteSpan(a, b Range) (int, bool) {
	aSpan, aOK := finiteSpan(a)
	bSpan, bOK := finiteSpan(b)

	switch {
	case aOK && bOK:
		return aSpan.Cmp(bSpan), true
	case aOK && !bOK:
		return -1, true
	case !aOK && bOK:
		return 1, true
	default:
		return 0, false
	}
}

func finiteSpan(r Range) (*big.Rat, bool) {
	if r.Lo.Kind != BoundFinite || r.Hi.Kind != BoundFinite {
		return nil, false
	}
	if r.Lo.Value.Num == nil || r.Lo.Value.Den == nil || r.Hi.Value.Num == nil || r.Hi.Value.Den == nil {
		return nil, false
	}

	lo := rationalToBigRat(r.Lo.Value)
	hi := rationalToBigRat(r.Hi.Value)
	return new(big.Rat).Sub(hi, lo), true
}

func rationalToBigRat(r Rational) *big.Rat {
	return new(big.Rat).SetFrac(
		new(big.Int).Set(r.Num),
		new(big.Int).Set(r.Den),
	)
}
func boundEqual(a, b Bound) bool {
	if a.Kind != b.Kind || a.Closed != b.Closed {
		return false
	}
	if a.Kind != BoundFinite {
		return true
	}
	if !rationalWellFormed(a.Value) || !rationalWellFormed(b.Value) {
		return false
	}
	return rationalCmp(a.Value, b.Value) == 0
}

func boundLT(a, b Bound) bool {
	return boundCmp(a, b) < 0
}

func boundLE(a, b Bound) bool {
	return boundCmp(a, b) <= 0
}

func boundCmp(a, b Bound) int {
	if a.Kind != b.Kind {
		return boundKindOrder(a.Kind) - boundKindOrder(b.Kind)
	}
	if a.Kind != BoundFinite {
		return 0
	}
	return rationalCmp(a.Value, b.Value)
}

func boundKindOrder(k BoundKind) int {
	switch k {
	case BoundNegInf:
		return -1
	case BoundFinite:
		return 0
	case BoundPosInf:
		return 1
	default:
		return 0
	}
}
func rationalCmp(a, b Rational) int {
	switch {
	case rationalWellFormed(a) && !rationalWellFormed(b):
		return 1
	case !rationalWellFormed(a) && rationalWellFormed(b):
		return -1
	case !rationalWellFormed(a) && !rationalWellFormed(b):
		return 0
	}

	left := new(big.Int).Mul(
		new(big.Int).Set(a.Num),
		new(big.Int).Set(b.Den),
	)
	right := new(big.Int).Mul(
		new(big.Int).Set(b.Num),
		new(big.Int).Set(a.Den),
	)
	return left.Cmp(right)
}
func rangeWellFormed(r Range) bool {
	return boundWellFormed(r.Lo) && boundWellFormed(r.Hi)
}

func boundWellFormed(b Bound) bool {
	switch b.Kind {
	case BoundNegInf, BoundPosInf:
		return true
	case BoundFinite:
		return rationalWellFormed(b.Value)
	default:
		return false
	}
}

func rationalWellFormed(r Rational) bool {
	return r.Num != nil && r.Den != nil && r.Den.Sign() != 0
}

// cf/api_types.go v3
