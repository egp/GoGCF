// cf/lft.go v10
package cf

import "math/big"

type binaryDecision int

const (
	decisionIngestLeft binaryDecision = iota
	decisionIngestRight
	decisionEmitOutput
)

type unaryLFT struct {
	a *big.Int
	b *big.Int
	c *big.Int
	d *big.Int
}

type binaryLFT struct {
	a *big.Int
	b *big.Int
	c *big.Int
	d *big.Int
	e *big.Int
	f *big.Int
	g *big.Int
	h *big.Int
}

type binaryStepState struct {
	canEmitOutput  bool
	canIngestLeft  bool
	canIngestRight bool
}

func identityUnaryLFT() unaryLFT {
	return unaryLFT{
		a: big.NewInt(1),
		b: big.NewInt(0),
		c: big.NewInt(0),
		d: big.NewInt(1),
	}
}

func identityBinaryLFT() binaryLFT {
	return binaryLFT{
		a: big.NewInt(1),
		b: big.NewInt(0),
		c: big.NewInt(0),
		d: big.NewInt(0),
		e: big.NewInt(0),
		f: big.NewInt(0),
		g: big.NewInt(0),
		h: big.NewInt(1),
	}
}

func (d binaryDecision) isValid() bool {
	switch d {
	case decisionIngestLeft, decisionIngestRight, decisionEmitOutput:
		return true
	default:
		return false
	}
}

func (s binaryStepState) choose() (binaryDecision, error) {
	switch {
	case s.canEmitOutput:
		return decisionEmitOutput, nil
	case s.canIngestLeft:
		return decisionIngestLeft, nil
	case s.canIngestRight:
		return decisionIngestRight, nil
	default:
		return 0, ErrStuck
	}
}

func (b binaryLFT) ingestLeft(p, q *big.Int) binaryLFT {
	return binaryLFT{
		a: addMul(b.a, q, b.c),
		b: addMul(b.b, q, b.d),
		c: mul(b.a, p),
		d: mul(b.b, p),
		e: addMul(b.e, q, b.g),
		f: addMul(b.f, q, b.h),
		g: mul(b.e, p),
		h: mul(b.f, p),
	}
}

func (b binaryLFT) ingestRight(p, q *big.Int) binaryLFT {
	return binaryLFT{
		a: addMul(b.a, q, b.b),
		b: mul(b.a, p),
		c: addMul(b.c, q, b.d),
		d: mul(b.c, p),
		e: addMul(b.e, q, b.f),
		f: mul(b.e, p),
		g: addMul(b.g, q, b.h),
		h: mul(b.g, p),
	}
}

func (b binaryLFT) exactQuotient() (Rational, bool) {
	u := b.collapseRightEOF()
	r, err := u.collapseEOFToRational()
	if err != nil {
		return Rational{}, false
	}
	return r, true
}

func (b binaryLFT) emit(term *big.Int) binaryLFT {
	if term == nil {
		return b
	}

	return binaryLFT{
		a: new(big.Int).Set(b.e),
		b: new(big.Int).Set(b.f),
		c: new(big.Int).Set(b.g),
		d: new(big.Int).Set(b.h),
		e: new(big.Int).Sub(new(big.Int).Set(b.a), new(big.Int).Mul(new(big.Int).Set(term), b.e)),
		f: new(big.Int).Sub(new(big.Int).Set(b.b), new(big.Int).Mul(new(big.Int).Set(term), b.f)),
		g: new(big.Int).Sub(new(big.Int).Set(b.c), new(big.Int).Mul(new(big.Int).Set(term), b.g)),
		h: new(big.Int).Sub(new(big.Int).Set(b.d), new(big.Int).Mul(new(big.Int).Set(term), b.h)),
	}
}

func (b binaryLFT) collapseLeftEOF() unaryLFT {
	return unaryLFT{
		a: new(big.Int).Set(b.a),
		b: new(big.Int).Set(b.b),
		c: new(big.Int).Set(b.e),
		d: new(big.Int).Set(b.f),
	}
}

func (b binaryLFT) collapseRightEOF() unaryLFT {
	return unaryLFT{
		a: new(big.Int).Set(b.a),
		b: new(big.Int).Set(b.c),
		c: new(big.Int).Set(b.e),
		d: new(big.Int).Set(b.g),
	}
}

func (u unaryLFT) collapseEOFToRational() (Rational, error) {
	if u.c == nil || u.c.Sign() == 0 {
		return Rational{}, ErrUndefined
	}

	num := new(big.Int).Set(u.a)
	den := new(big.Int).Set(u.c)

	if den.Sign() < 0 {
		num.Neg(num)
		den.Neg(den)
	}

	return Rational{
		Num: num,
		Den: den,
	}, nil
}

func mul(x, y *big.Int) *big.Int {
	return new(big.Int).Mul(x, y)
}

func addMul(x, y, z *big.Int) *big.Int {
	return new(big.Int).Add(new(big.Int).Mul(x, y), z)
}

func (b binaryLFT) evalAt(x, y Rational) (Rational, error) {
	if b.a == nil || b.b == nil || b.c == nil || b.d == nil ||
		b.e == nil || b.f == nil || b.g == nil || b.h == nil {
		return Rational{}, ErrUndefined
	}

	xr, err := normalizedRational(x)
	if err != nil {
		return Rational{}, err
	}
	yr, err := normalizedRational(y)
	if err != nil {
		return Rational{}, err
	}

	xn := xr.Num
	xd := xr.Den
	yn := yr.Num
	yd := yr.Den

	xyNum := new(big.Int).Mul(new(big.Int).Set(xn), new(big.Int).Set(yn))
	commonDen := new(big.Int).Mul(new(big.Int).Set(xd), new(big.Int).Set(yd))

	xScaled := new(big.Int).Mul(new(big.Int).Set(xn), new(big.Int).Set(yd))
	yScaled := new(big.Int).Mul(new(big.Int).Set(yn), new(big.Int).Set(xd))

	num := new(big.Int).Mul(new(big.Int).Set(b.a), xyNum)
	num.Add(num, new(big.Int).Mul(new(big.Int).Set(b.b), xScaled))
	num.Add(num, new(big.Int).Mul(new(big.Int).Set(b.c), yScaled))
	num.Add(num, new(big.Int).Mul(new(big.Int).Set(b.d), commonDen))

	den := new(big.Int).Mul(new(big.Int).Set(b.e), xyNum)
	den.Add(den, new(big.Int).Mul(new(big.Int).Set(b.f), xScaled))
	den.Add(den, new(big.Int).Mul(new(big.Int).Set(b.g), yScaled))
	den.Add(den, new(big.Int).Mul(new(big.Int).Set(b.h), commonDen))

	if den.Sign() == 0 {
		return Rational{}, ErrUndefined
	}

	return normalizedRational(Rational{
		Num: num,
		Den: den,
	})
}

func (b binaryLFT) rangeFromXYRanges(xr, yr Range) (Range, error) {
	x, xExact := exactFiniteRangeValue(xr)
	y, yExact := exactFiniteRangeValue(yr)
	if xExact && yExact {
		z, err := b.evalAt(x, y)
		if err != nil {
			return Range{}, err
		}
		return exactRangeFromRational(z), nil
	}

	xLo, xHi, xOK := finiteInsideRangeEndpoints(xr)
	yLo, yHi, yOK := finiteInsideRangeEndpoints(yr)
	if !xOK || !yOK {
		return Range{}, nil
	}

	corners := [][2]Rational{
		{xLo, yLo},
		{xLo, yHi},
		{xHi, yLo},
		{xHi, yHi},
	}

	values := make([]Rational, 0, 4)
	for _, corner := range corners {
		v, err := b.evalAt(corner[0], corner[1])
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

	allClosed := xr.Lo.Closed && xr.Hi.Closed && yr.Lo.Closed && yr.Hi.Closed
	return finiteInsideRange(lo, hi, allClosed), nil
}

func normalizedRational(r Rational) (Rational, error) {
	if r.Num == nil || r.Den == nil || r.Den.Sign() == 0 {
		return Rational{}, ErrUndefined
	}

	num := new(big.Int).Set(r.Num)
	den := new(big.Int).Set(r.Den)

	if den.Sign() < 0 {
		num.Neg(num)
		den.Neg(den)
	}

	g := new(big.Int).GCD(nil, nil, num, den)
	if g.Sign() != 0 && g.Cmp(big.NewInt(1)) != 0 {
		num.Quo(num, g)
		den.Quo(den, g)
	}

	return Rational{
		Num: num,
		Den: den,
	}, nil
}

func exactFiniteRangeValue(r Range) (Rational, bool) {
	if !r.IsExact() {
		return Rational{}, false
	}
	if r.Lo.Kind != BoundFinite || r.Hi.Kind != BoundFinite {
		return Rational{}, false
	}
	v, err := normalizedRational(r.Lo.Value)
	if err != nil {
		return Rational{}, false
	}
	return v, true
}

func exactRangeFromRational(r Rational) Range {
	v, err := normalizedRational(r)
	if err != nil {
		return Range{}
	}

	b := Bound{
		Kind:   BoundFinite,
		Value:  v,
		Closed: true,
	}

	return Range{
		Lo:      b,
		Hi:      b,
		Outside: false,
	}
}

func (b binaryLFT) emitCandidateFromRanges(xr, yr Range) (*big.Int, bool, error) {
	zr, err := b.rangeFromXYRanges(xr, yr)
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

func finiteInsideRangeEndpoints(r Range) (Rational, Rational, bool) {
	if r.Outside {
		return Rational{}, Rational{}, false
	}
	if r.Lo.Kind != BoundFinite || r.Hi.Kind != BoundFinite {
		return Rational{}, Rational{}, false
	}

	lo, err := normalizedRational(r.Lo.Value)
	if err != nil {
		return Rational{}, Rational{}, false
	}
	hi, err := normalizedRational(r.Hi.Value)
	if err != nil {
		return Rational{}, Rational{}, false
	}
	return lo, hi, true
}

func finiteInsideRange(lo, hi Rational, closed bool) Range {
	loN, err := normalizedRational(lo)
	if err != nil {
		return Range{}
	}
	hiN, err := normalizedRational(hi)
	if err != nil {
		return Range{}
	}

	return Range{
		Lo: Bound{
			Kind:   BoundFinite,
			Value:  loN,
			Closed: closed,
		},
		Hi: Bound{
			Kind:   BoundFinite,
			Value:  hiN,
			Closed: closed,
		},
		Outside: false,
	}
}

type diagLFT struct {
	a *big.Int
	b *big.Int
	c *big.Int
	d *big.Int
	e *big.Int
	f *big.Int
}

func identityDiagLFT() diagLFT {
	return diagLFT{
		a: big.NewInt(0),
		b: big.NewInt(1),
		c: big.NewInt(0),
		d: big.NewInt(0),
		e: big.NewInt(0),
		f: big.NewInt(1),
	}
}

func (d diagLFT) evalAt(x Rational) (Rational, error) {
	if d.a == nil || d.b == nil || d.c == nil ||
		d.d == nil || d.e == nil || d.f == nil {
		return Rational{}, ErrUndefined
	}

	xr, err := normalizedRational(x)
	if err != nil {
		return Rational{}, err
	}

	xn := xr.Num
	xd := xr.Den

	x2Num := new(big.Int).Mul(new(big.Int).Set(xn), new(big.Int).Set(xn))
	x2Den := new(big.Int).Mul(new(big.Int).Set(xd), new(big.Int).Set(xd))
	xScaled := new(big.Int).Mul(new(big.Int).Set(xn), new(big.Int).Set(xd))

	num := new(big.Int).Mul(new(big.Int).Set(d.a), x2Num)
	num.Add(num, new(big.Int).Mul(new(big.Int).Set(d.b), xScaled))
	num.Add(num, new(big.Int).Mul(new(big.Int).Set(d.c), x2Den))

	den := new(big.Int).Mul(new(big.Int).Set(d.d), x2Num)
	den.Add(den, new(big.Int).Mul(new(big.Int).Set(d.e), xScaled))
	den.Add(den, new(big.Int).Mul(new(big.Int).Set(d.f), x2Den))

	if den.Sign() == 0 {
		return Rational{}, ErrUndefined
	}

	return normalizedRational(Rational{
		Num: num,
		Den: den,
	})
}

func (b binaryLFT) diagonal() diagLFT {
	return diagLFT{
		a: new(big.Int).Set(b.a),
		b: new(big.Int).Add(new(big.Int).Set(b.b), b.c),
		c: new(big.Int).Set(b.d),
		d: new(big.Int).Set(b.e),
		e: new(big.Int).Add(new(big.Int).Set(b.f), b.g),
		f: new(big.Int).Set(b.h),
	}
}

// cf/lft.go v10
