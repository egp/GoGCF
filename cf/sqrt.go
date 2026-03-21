// cf/sqrt.go v2
package cf

import "math/big"

type sqrtStrategy struct {
	post unaryLFT
}

func Sqrt(x GCF) GCF {
	return strategyUnaryGCF{
		strategy: sqrtStrategy{post: identityUnaryLFT()},
		child:    x,
	}
}

func (s sqrtStrategy) RangeFromOperand(xr Range) (Range, error) {
	x, exact := exactFiniteRangeValue(xr)
	if !exact {
		return Range{}, nil
	}

	r, ok, err := s.ExactEval(x)
	if err != nil {
		return Range{}, err
	}
	if !ok {
		return Range{}, nil
	}

	return exactRangeFromRational(r), nil
}

func (s sqrtStrategy) EmitCandidateFromOperand(xr Range) (*big.Int, bool, error) {
	zr, err := s.RangeFromOperand(xr)
	if err != nil {
		return nil, false, err
	}
	if !rangeWellFormed(zr) {
		return nil, false, nil
	}
	if zr.Outside {
		return nil, false, nil
	}
	if zr.Lo.Kind != BoundFinite || zr.Hi.Kind != BoundFinite {
		return nil, false, nil
	}
	if !rationalWellFormed(zr.Lo.Value) || !rationalWellFormed(zr.Hi.Value) {
		return nil, false, nil
	}

	lo, err := normalizedRational(zr.Lo.Value)
	if err != nil {
		return nil, false, err
	}
	hi, err := normalizedRational(zr.Hi.Value)
	if err != nil {
		return nil, false, err
	}

	qLo, _ := floorDivModBig(lo.Num, lo.Den)
	qHi, hiRem := floorDivModBig(hi.Num, hi.Den)
	if !zr.Hi.Closed && hiRem.Sign() == 0 {
		qHi.Sub(qHi, big.NewInt(1))
	}

	if qLo.Cmp(qHi) != 0 {
		return nil, false, nil
	}

	return qLo, true, nil
}

func (s sqrtStrategy) Emit(term *big.Int) unaryStrategy {
	return sqrtStrategy{
		post: s.effectivePost().emit(term),
	}
}

func (s sqrtStrategy) ExactEval(x Rational) (Rational, bool, error) {
	xr, err := normalizedRational(x)
	if err != nil {
		return Rational{}, false, err
	}

	if xr.Num.Sign() < 0 {
		return Rational{}, false, ErrUndefined
	}

	rootNum, okNum := exactBigIntSqrt(xr.Num)
	if !okNum {
		return Rational{}, false, nil
	}
	rootDen, okDen := exactBigIntSqrt(xr.Den)
	if !okDen {
		return Rational{}, false, nil
	}

	y, err := normalizedRational(Rational{
		Num: rootNum,
		Den: rootDen,
	})
	if err != nil {
		return Rational{}, false, err
	}

	r, err := s.effectivePost().evalAt(y)
	if err != nil {
		return Rational{}, false, err
	}
	return r, true, nil
}

func exactBigIntSqrt(n *big.Int) (*big.Int, bool) {
	if n == nil || n.Sign() < 0 {
		return nil, false
	}

	root := new(big.Int).Sqrt(new(big.Int).Set(n))
	sq := new(big.Int).Mul(new(big.Int).Set(root), new(big.Int).Set(root))
	if sq.Cmp(n) != 0 {
		return nil, false
	}
	return root, true
}

func (s sqrtStrategy) effectivePost() unaryLFT {
	if s.post.a == nil || s.post.b == nil || s.post.c == nil || s.post.d == nil {
		return identityUnaryLFT()
	}
	return s.post
}

// cf/sqrt.go v2
