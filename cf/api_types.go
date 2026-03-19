// cf/api_types.go v1
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

// cf/api_types.go v1
