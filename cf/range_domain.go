// cf/range_domain.go v2
package cf

import "math/big"

func rangeMaybeIncludesZero(r Range) bool {
	if !rangeWellFormed(r) {
		return true
	}

	insideHasZero := insideRangeIncludesZero(r)
	if r.Outside {
		return !insideHasZero
	}
	return insideHasZero
}

func rangeCertainlyNegative(r Range) bool {
	if !rangeWellFormed(r) {
		return false
	}

	if !r.Outside {
		return upperBoundIsStrictlyNegative(r.Hi)
	}

	if r.Hi.Kind != BoundPosInf {
		return false
	}
	loCmp, ok := boundNumberCmpZero(r.Lo)
	if !ok {
		return false
	}
	return loCmp < 0 || (loCmp == 0 && r.Lo.Closed)
}

func rangeCertainlyNonNegative(r Range) bool {
	if !rangeWellFormed(r) {
		return false
	}

	if !r.Outside {
		return lowerBoundIsNonNegative(r.Lo)
	}

	if r.Lo.Kind != BoundNegInf {
		return false
	}
	hiCmp, ok := boundNumberCmpZero(r.Hi)
	if !ok {
		return false
	}
	return hiCmp >= 0
}

func insideRangeIncludesZero(r Range) bool {
	return lowerBoundAllowsZero(r.Lo) && upperBoundAllowsZero(r.Hi)
}

func lowerBoundAllowsZero(b Bound) bool {
	switch b.Kind {
	case BoundNegInf:
		return true
	case BoundPosInf:
		return false
	case BoundFinite:
		cmp, ok := boundNumberCmpZero(b)
		if !ok {
			return false
		}
		return cmp < 0 || (cmp == 0 && b.Closed)
	default:
		return false
	}
}

func upperBoundAllowsZero(b Bound) bool {
	switch b.Kind {
	case BoundPosInf:
		return true
	case BoundNegInf:
		return false
	case BoundFinite:
		cmp, ok := boundNumberCmpZero(b)
		if !ok {
			return false
		}
		return cmp > 0 || (cmp == 0 && b.Closed)
	default:
		return false
	}
}

func upperBoundIsStrictlyNegative(b Bound) bool {
	switch b.Kind {
	case BoundNegInf:
		return true
	case BoundPosInf:
		return false
	case BoundFinite:
		cmp, ok := boundNumberCmpZero(b)
		if !ok {
			return false
		}
		return cmp < 0 || (cmp == 0 && !b.Closed)
	default:
		return false
	}
}

func lowerBoundIsNonNegative(b Bound) bool {
	switch b.Kind {
	case BoundPosInf:
		return true
	case BoundNegInf:
		return false
	case BoundFinite:
		cmp, ok := boundNumberCmpZero(b)
		if !ok {
			return false
		}
		return cmp > 0 || (cmp == 0 && b.Closed)
	default:
		return false
	}
}

func boundNumberCmpZero(b Bound) (int, bool) {
	switch b.Kind {
	case BoundNegInf:
		return -1, true
	case BoundPosInf:
		return 1, true
	case BoundFinite:
		if !rationalWellFormed(b.Value) {
			return 0, false
		}
		return rationalCmp(b.Value, RationalZero()), true
	default:
		return 0, false
	}
}

func RationalZero() Rational {
	return Rational{Num: bigZero(), Den: bigOne()}
}

func bigZero() *big.Int {
	return new(big.Int)
}

func bigOne() *big.Int {
	return big.NewInt(1)
}

// cf/range_domain.go v2
