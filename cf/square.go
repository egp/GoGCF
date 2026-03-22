// cf/square.go v1
package cf

import "math/big"

type squareStrategy struct {
	op diagLFT
}

func Square(x GCF) GCF {
	return strategyUnaryGCF{
		strategy: squareStrategy{op: squareDiagLFT()},
		child:    x,
	}
}

func squareDiagLFT() diagLFT {
	return diagLFT{
		a: big.NewInt(1),
		b: big.NewInt(0),
		c: big.NewInt(0),
		d: big.NewInt(0),
		e: big.NewInt(0),
		f: big.NewInt(1),
	}
}

func (s squareStrategy) effectiveOp() diagLFT {
	if s.op.a == nil || s.op.b == nil || s.op.c == nil ||
		s.op.d == nil || s.op.e == nil || s.op.f == nil {
		return squareDiagLFT()
	}
	return s.op
}

func (s squareStrategy) RangeFromOperand(xr Range) (Range, error) {
	op := s.effectiveOp()

	x, exact := exactFiniteRangeValue(xr)
	if exact {
		r, err := op.evalAt(x)
		if err != nil {
			return Range{}, err
		}
		return exactRangeFromRational(r), nil
	}

	xLo, xHi, ok := finiteInsideRangeEndpoints(xr)
	if !ok {
		return Range{}, nil
	}

	loSq, err := op.evalAt(xLo)
	if err != nil {
		return Range{}, err
	}
	hiSq, err := op.evalAt(xHi)
	if err != nil {
		return Range{}, err
	}

	type endpointImage struct {
		value  Rational
		closed bool
	}

	candidates := []endpointImage{
		{value: loSq, closed: xr.Lo.Closed},
		{value: hiSq, closed: xr.Hi.Closed},
	}
	if insideRangeIncludesZero(xr) {
		candidates = append(candidates, endpointImage{
			value:  RationalZero(),
			closed: true,
		})
	}

	lo := candidates[0].value
	loClosed := candidates[0].closed
	hi := candidates[0].value
	hiClosed := candidates[0].closed

	for _, candidate := range candidates[1:] {
		switch cmp := rationalCmp(candidate.value, lo); {
		case cmp < 0:
			lo = candidate.value
			loClosed = candidate.closed
		case cmp == 0:
			loClosed = loClosed || candidate.closed
		}

		switch cmp := rationalCmp(candidate.value, hi); {
		case cmp > 0:
			hi = candidate.value
			hiClosed = candidate.closed
		case cmp == 0:
			hiClosed = hiClosed || candidate.closed
		}
	}

	loN, err := normalizedRational(lo)
	if err != nil {
		return Range{}, err
	}
	hiN, err := normalizedRational(hi)
	if err != nil {
		return Range{}, err
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
	}, nil
}

func (s squareStrategy) EmitCandidateFromOperand(xr Range) (*big.Int, bool, error) {
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

func (s squareStrategy) Emit(term *big.Int) unaryStrategy {
	return diagLFTStrategy{op: s.effectiveOp().emit(term)}
}

func (s squareStrategy) ExactEval(x Rational) (Rational, bool, error) {
	r, err := s.effectiveOp().evalAt(x)
	if err != nil {
		return Rational{}, false, err
	}
	return r, true, nil
}

// cf/square.go v1
