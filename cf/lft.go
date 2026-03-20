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
		a: addMul(b.a, q, b.b),
		b: addMul(b.c, q, b.d),
		c: mul(b.a, p),
		d: mul(b.b, p),
		e: addMul(b.e, q, b.f),
		f: addMul(b.g, q, b.h),
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

// cf/lft.go v10
