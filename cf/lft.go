// cf/lft.go v7
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

func mul(x, y *big.Int) *big.Int {
	return new(big.Int).Mul(x, y)
}

func addMul(x, y, z *big.Int) *big.Int {
	return new(big.Int).Add(new(big.Int).Mul(x, y), z)
}

// cf/lft.go v7
