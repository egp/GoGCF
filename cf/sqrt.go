// cf/sqrt.go v4
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

func (s sqrtStrategy) effectivePost() unaryLFT {
	if s.post.a == nil || s.post.b == nil || s.post.c == nil || s.post.d == nil {
		return identityUnaryLFT()
	}
	return s.post
}

func (s sqrtStrategy) RangeFromOperand(xr Range) (Range, error) {
	post := s.effectivePost()

	x, exact := exactFiniteRangeValue(xr)
	if exact {
		r, ok, err := s.ExactEval(x)
		if err != nil {
			return Range{}, err
		}
		if ok {
			return exactRangeFromRational(r), nil
		}

		zr, ok, err := exactPositiveSqrtImageEnclosure(x, post)
		if err != nil {
			return Range{}, err
		}
		if ok {
			return zr, nil
		}
		return Range{}, nil
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

	lo, okLo, err := exactSqrtRational(xLo)
	if err != nil {
		return Range{}, err
	}
	if !okLo {
		return Range{}, nil
	}

	hi, okHi, err := exactSqrtRational(xHi)
	if err != nil {
		return Range{}, err
	}
	if !okHi {
		return Range{}, nil
	}

	base := Range{
		Lo: Bound{
			Kind:   BoundFinite,
			Value:  lo,
			Closed: xr.Lo.Closed,
		},
		Hi: Bound{
			Kind:   BoundFinite,
			Value:  hi,
			Closed: xr.Hi.Closed,
		},
		Outside: false,
	}

	return post.rangeFromXRange(base)
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

func exactPositiveSqrtImageEnclosure(x Rational, post unaryLFT) (Range, bool, error) {
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

	base, err := exactPositiveSqrtEnclosure(xr)
	if err != nil {
		return Range{}, false, err
	}

	if q, ok, err := post.emitCandidateFromRange(base); err != nil {
		return Range{}, false, err
	} else if ok {
		return openUnitIntervalAt(q), true, nil
	}

	lo, hi, ok := finiteInsideRangeEndpoints(base)
	if !ok {
		return Range{}, false, nil
	}

	for i := 0; i < 64; i++ {
		mid, err := rationalMidpoint(lo, hi)
		if err != nil {
			return Range{}, false, err
		}

		cmp, err := rationalSquareCmp(mid, xr)
		if err != nil {
			return Range{}, false, err
		}
		if cmp < 0 {
			lo = mid
		} else if cmp > 0 {
			hi = mid
		} else {
			r, err := post.evalAt(mid)
			if err != nil {
				return Range{}, false, err
			}
			return exactRangeFromRational(r), true, nil
		}

		refined := Range{
			Lo: Bound{
				Kind:   BoundFinite,
				Value:  lo,
				Closed: false,
			},
			Hi: Bound{
				Kind:   BoundFinite,
				Value:  hi,
				Closed: false,
			},
			Outside: false,
		}

		if q, ok, err := post.emitCandidateFromRange(refined); err != nil {
			return Range{}, false, err
		} else if ok {
			return openUnitIntervalAt(q), true, nil
		}
	}

	return Range{}, false, nil
}

func exactPositiveSqrtEnclosure(x Rational) (Range, error) {
	xr, err := normalizedRational(x)
	if err != nil {
		return Range{}, err
	}
	if xr.Num.Sign() < 0 {
		return Range{}, ErrUndefined
	}

	if root, ok, err := exactSqrtRational(xr); err != nil {
		return Range{}, err
	} else if ok {
		return exactRangeFromRational(root), nil
	}

	lo, err := normalizedRational(Rational{
		Num: floorBigIntSqrt(xr.Num),
		Den: ceilBigIntSqrt(xr.Den),
	})
	if err != nil {
		return Range{}, err
	}
	hi, err := normalizedRational(Rational{
		Num: ceilBigIntSqrt(xr.Num),
		Den: floorBigIntSqrt(xr.Den),
	})
	if err != nil {
		return Range{}, err
	}

	return Range{
		Lo: Bound{
			Kind:   BoundFinite,
			Value:  lo,
			Closed: false,
		},
		Hi: Bound{
			Kind:   BoundFinite,
			Value:  hi,
			Closed: false,
		},
		Outside: false,
	}, nil
}

func openUnitIntervalAt(q *big.Int) Range {
	lo := new(big.Int).Set(q)
	hi := new(big.Int).Add(new(big.Int).Set(q), big.NewInt(1))

	return Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: lo,
				Den: big.NewInt(1),
			},
			Closed: false,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: hi,
				Den: big.NewInt(1),
			},
			Closed: false,
		},
		Outside: false,
	}
}

func rationalMidpoint(a, b Rational) (Rational, error) {
	ar, err := normalizedRational(a)
	if err != nil {
		return Rational{}, err
	}
	br, err := normalizedRational(b)
	if err != nil {
		return Rational{}, err
	}

	num := new(big.Int).Mul(new(big.Int).Set(ar.Num), new(big.Int).Set(br.Den))
	num.Add(num, new(big.Int).Mul(new(big.Int).Set(br.Num), new(big.Int).Set(ar.Den)))

	den := new(big.Int).Mul(new(big.Int).Set(ar.Den), new(big.Int).Set(br.Den))
	den.Mul(den, big.NewInt(2))

	return normalizedRational(Rational{
		Num: num,
		Den: den,
	})
}

func rationalSquareCmp(a, x Rational) (int, error) {
	ar, err := normalizedRational(a)
	if err != nil {
		return 0, err
	}
	xr, err := normalizedRational(x)
	if err != nil {
		return 0, err
	}

	left := new(big.Int).Mul(new(big.Int).Set(ar.Num), new(big.Int).Set(ar.Num))
	left.Mul(left, xr.Den)

	right := new(big.Int).Mul(new(big.Int).Set(ar.Den), new(big.Int).Set(ar.Den))
	right.Mul(right, xr.Num)

	return left.Cmp(right), nil
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

// cf/sqrt.go v4
