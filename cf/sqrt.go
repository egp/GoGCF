// cf/sqrt.go v6
package cf

import "math/big"

type sqrtStrategy struct {
	post  unaryLFT
	upper Rational
}

func Sqrt(x GCF) GCF {
	return strategyUnaryGCF{
		strategy: sqrtStrategy{post: identityUnaryLFT()},
		child:    x,
	}
}

func (s sqrtStrategy) effectivePost() unaryLFT {
	if s.post.a == nil || s.post.b == nil || s.post.c == nil || s.post.d == nil {
		return identityUnaryLFT()
	}
	return s.post
}

func (s sqrtStrategy) OuroborosFeedback(xr Range) (Range, *big.Int, unaryStrategy, bool, error) {
	x, exact := exactFiniteRangeValue(xr)
	if !exact {
		return Range{}, nil, s, false, nil
	}

	xrNorm, err := normalizedRational(x)
	if err != nil {
		return Range{}, nil, s, false, err
	}
	if xrNorm.Num.Sign() < 0 {
		return Range{}, nil, s, false, ErrUndefined
	}

	if _, ok, err := exactSqrtRational(xrNorm); err != nil {
		return Range{}, nil, s, false, err
	} else if ok {
		return Range{}, nil, s, false, nil
	}

	got, ok, err := exactPositiveSqrtNewtonFeedback(xrNorm, s.effectivePost(), s.upper)
	if err != nil {
		return Range{}, nil, s, false, err
	}
	if !ok {
		return Range{}, nil, s, false, nil
	}

	next := sqrtStrategy{
		post:  s.effectivePost(),
		upper: got.UpperBound,
	}
	return got.ImageRange, got.Candidate, next, true, nil
}

func (s sqrtStrategy) RangeFromOperand(xr Range) (Range, error) {
	if x, exact := exactFiniteRangeValue(xr); exact {
		return s.rangeFromExactOperand(x)
	}

	if rangeCertainlyNegative(xr) {
		return Range{}, ErrUndefined
	}
	if !rangeCertainlyNonNegative(xr) {
		return Range{}, nil
	}

	xLo, xHi, ok := finiteInsideRangeEndpoints(xr)
	if !ok {
		return Range{}, nil
	}

	return s.rangeFromFinitePositiveInterval(xLo, xHi, xr.Lo.Closed, xr.Hi.Closed)
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
		post:  s.effectivePost().emit(term),
		upper: s.upper,
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

	root, ok, err := exactSqrtRational(xr)
	if err != nil {
		return Rational{}, false, err
	}
	if !ok {
		return Rational{}, false, nil
	}

	r, err := s.effectivePost().evalAt(root)
	if err != nil {
		return Rational{}, false, err
	}
	return r, true, nil
}

func exactSqrtRational(x Rational) (Rational, bool, error) {
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

	r, err := normalizedRational(Rational{
		Num: rootNum,
		Den: rootDen,
	})
	if err != nil {
		return Rational{}, false, err
	}
	return r, true, nil
}

func exactPositiveSqrtImageEnclosure(x Rational, post unaryLFT, seedUpper Rational) (Range, bool, error) {
	xr, err := normalizedRational(x)
	if err != nil {
		return Range{}, false, err
	}
	if xr.Num.Sign() < 0 {
		return Range{}, false, ErrUndefined
	}

	if root, ok, err := exactSqrtRational(xr); err != nil {
		return Range{}, false, err
	} else if ok {
		r, err := post.evalAt(root)
		if err != nil {
			return Range{}, false, err
		}
		return exactRangeFromRational(r), true, nil
	}

	got, ok, err := exactPositiveSqrtNewtonFeedback(xr, post, seedUpper)
	if err != nil {
		return Range{}, false, err
	}
	if !ok {
		return Range{}, false, nil
	}
	return got.ImageRange, true, nil
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

func floorBigIntSqrt(n *big.Int) *big.Int {
	if n == nil || n.Sign() < 0 {
		return nil
	}
	return new(big.Int).Sqrt(new(big.Int).Set(n))
}

func ceilBigIntSqrt(n *big.Int) *big.Int {
	if n == nil || n.Sign() < 0 {
		return nil
	}
	floor := floorBigIntSqrt(n)
	sq := new(big.Int).Mul(new(big.Int).Set(floor), new(big.Int).Set(floor))
	if sq.Cmp(n) == 0 {
		return floor
	}
	return new(big.Int).Add(floor, big.NewInt(1))
}

func combineConservativeEndpointImageRanges(a, b Range) (Range, bool) {
	if !rangeWellFormed(a) || !rangeWellFormed(b) {
		return Range{}, false
	}
	if a.Outside || b.Outside {
		return Range{}, false
	}
	if a.Lo.Kind != BoundFinite || a.Hi.Kind != BoundFinite || b.Lo.Kind != BoundFinite || b.Hi.Kind != BoundFinite {
		return Range{}, false
	}
	if !rationalWellFormed(a.Lo.Value) || !rationalWellFormed(a.Hi.Value) ||
		!rationalWellFormed(b.Lo.Value) || !rationalWellFormed(b.Hi.Value) {
		return Range{}, false
	}

	lo := a.Lo.Value
	loClosed := a.Lo.Closed
	switch cmp := rationalCmp(b.Lo.Value, lo); {
	case cmp < 0:
		lo = b.Lo.Value
		loClosed = b.Lo.Closed
	case cmp == 0:
		loClosed = loClosed || b.Lo.Closed
	}

	hi := a.Hi.Value
	hiClosed := a.Hi.Closed
	switch cmp := rationalCmp(b.Hi.Value, hi); {
	case cmp > 0:
		hi = b.Hi.Value
		hiClosed = b.Hi.Closed
	case cmp == 0:
		hiClosed = hiClosed || b.Hi.Closed
	}

	loN, err := normalizedRational(lo)
	if err != nil {
		return Range{}, false
	}
	hiN, err := normalizedRational(hi)
	if err != nil {
		return Range{}, false
	}

	return Range{
		Lo: Bound{
			Kind:   BoundFinite,
			Value:  loN,
			Closed: loClosed,
		},
		Hi: Bound{
			Kind:   BoundFinite,
			Value:  hiN,
			Closed: hiClosed,
		},
		Outside: false,
	}, true
}

func sqrtEndpointImageRangeForInterval(x Rational, post unaryLFT) (Range, bool, error) {
	xr, err := normalizedRational(x)
	if err != nil {
		return Range{}, false, err
	}
	if xr.Num.Sign() < 0 {
		return Range{}, false, ErrUndefined
	}

	if root, ok, err := exactSqrtRational(xr); err != nil {
		return Range{}, false, err
	} else if ok {
		r, err := post.evalAt(root)
		if err != nil {
			if err == ErrUndefined {
				return Range{}, false, nil
			}
			return Range{}, false, err
		}
		return exactRangeFromRational(r), true, nil
	}

	got, ok, err := exactPositiveSqrtNewtonFeedback(xr, post, Rational{})
	if err != nil {
		if err == ErrUndefined {
			return Range{}, false, nil
		}
		return Range{}, false, err
	}
	if !ok {
		return Range{}, false, nil
	}

	return got.ImageRange, true, nil
}
func (s sqrtStrategy) rangeFromFinitePositiveInterval(xLo, xHi Rational, loClosed, hiClosed bool) (Range, error) {
	post := s.effectivePost()

	loRoot, loExact, err := exactSqrtRational(xLo)
	if err != nil {
		return Range{}, err
	}
	hiRoot, hiExact, err := exactSqrtRational(xHi)
	if err != nil {
		return Range{}, err
	}

	if loExact && hiExact {
		base := Range{
			Lo: Bound{
				Kind:   BoundFinite,
				Value:  loRoot,
				Closed: loClosed,
			},
			Hi: Bound{
				Kind:   BoundFinite,
				Value:  hiRoot,
				Closed: hiClosed,
			},
			Outside: false,
		}
		return post.rangeFromXRange(base)
	}

	loImg, loOK, err := sqrtEndpointImageRangeForInterval(xLo, post)
	if err != nil {
		return Range{}, err
	}

	hiImg, hiOK, err := sqrtEndpointImageRangeForInterval(xHi, post)
	if err != nil {
		return Range{}, err
	}

	if !loOK || !hiOK {
		return Range{}, nil
	}

	zr, ok := combineConservativeEndpointImageRanges(loImg, hiImg)
	if !ok {
		return Range{}, nil
	}
	return zr, nil
}

func (s sqrtStrategy) rangeFromExactOperand(x Rational) (Range, error) {
	post := s.effectivePost()

	r, ok, err := s.ExactEval(x)
	if err != nil {
		return Range{}, err
	}
	if ok {
		return exactRangeFromRational(r), nil
	}

	zr, ok, err := exactPositiveSqrtImageEnclosure(x, post, s.upper)
	if err != nil {
		return Range{}, err
	}
	if ok {
		return zr, nil
	}

	return Range{}, nil
}

// cf/sqrt.go v6
