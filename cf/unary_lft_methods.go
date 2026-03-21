// cf/unary_lft_methods.go v2
package cf

import "math/big"

func (u unaryLFT) evalAt(x Rational) (Rational, error) {
	if u.a == nil || u.b == nil || u.c == nil || u.d == nil {
		return Rational{}, ErrUndefined
	}

	xr, err := normalizedRational(x)
	if err != nil {
		return Rational{}, err
	}

	xn := xr.Num
	xd := xr.Den

	num := new(big.Int).Mul(new(big.Int).Set(u.a), new(big.Int).Set(xn))
	num.Add(num, new(big.Int).Mul(new(big.Int).Set(u.b), new(big.Int).Set(xd)))

	den := new(big.Int).Mul(new(big.Int).Set(u.c), new(big.Int).Set(xn))
	den.Add(den, new(big.Int).Mul(new(big.Int).Set(u.d), new(big.Int).Set(xd)))

	if den.Sign() == 0 {
		return Rational{}, ErrUndefined
	}

	return normalizedRational(Rational{
		Num: num,
		Den: den,
	})
}

func (u unaryLFT) emit(term *big.Int) unaryLFT {
	if term == nil {
		return u
	}

	return unaryLFT{
		a: new(big.Int).Set(u.c),
		b: new(big.Int).Set(u.d),
		c: new(big.Int).Sub(new(big.Int).Set(u.a), new(big.Int).Mul(new(big.Int).Set(term), u.c)),
		d: new(big.Int).Sub(new(big.Int).Set(u.b), new(big.Int).Mul(new(big.Int).Set(term), u.d)),
	}
}

func (u unaryLFT) rangeFromXRange(xr Range) (Range, error) {
	x, exact := exactFiniteRangeValue(xr)
	if exact {
		z, err := u.evalAt(x)
		if err != nil {
			return Range{}, err
		}
		return exactRangeFromRational(z), nil
	}

	xLo, xHi, ok := finiteInsideRangeEndpoints(xr)
	if !ok {
		return Range{}, nil
	}
	if unaryLFTPoleTouchesIntervalClosure(u, xLo, xHi) {
		return Range{}, nil
	}

	values := make([]Rational, 0, 2)
	for _, x := range []Rational{xLo, xHi} {
		v, err := u.evalAt(x)
		if err != nil {
			continue
		}
		values = append(values, v)
	}

	if len(values) == 0 {
		return Range{}, nil
	}

	lo := values[0]
	hi := values[0]
	for _, v := range values[1:] {
		if rationalCmp(v, lo) < 0 {
			lo = v
		}
		if rationalCmp(v, hi) > 0 {
			hi = v
		}
	}

	allClosed := xr.Lo.Closed && xr.Hi.Closed
	return finiteInsideRange(lo, hi, allClosed), nil
}

func (u unaryLFT) emitCandidateFromRange(xr Range) (*big.Int, bool, error) {
	zr, err := u.rangeFromXRange(xr)
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

func unaryLFTPoleTouchesIntervalClosure(u unaryLFT, lo, hi Rational) bool {
	if u.c == nil || u.d == nil || u.c.Sign() == 0 {
		return false
	}

	root, err := normalizedRational(Rational{
		Num: new(big.Int).Neg(new(big.Int).Set(u.d)),
		Den: new(big.Int).Set(u.c),
	})
	if err != nil {
		return false
	}

	return rationalCmp(lo, root) <= 0 && rationalCmp(root, hi) <= 0
}

// cf/unary_lft_methods.go v2
