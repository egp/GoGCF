// cf/binary.go v2
package cf

import "math/big"

type binaryGCF struct {
	op    binaryLFT
	left  GCF
	right GCF
}

func Add(x, y GCF) GCF {
	return binaryGCF{
		op:    additionBinaryLFT(),
		left:  x,
		right: y,
	}
}

func Sub(x, y GCF) GCF {
	return binaryGCF{
		op:    subtractionBinaryLFT(),
		left:  x,
		right: y,
	}
}

func Mul(x, y GCF) GCF {
	return binaryGCF{
		op:    multiplicationBinaryLFT(),
		left:  x,
		right: y,
	}
}

func Div(x, y GCF) GCF {
	return binaryGCF{
		op:    divisionBinaryLFT(),
		left:  x,
		right: y,
	}
}

func (g binaryGCF) Next() (RCFTerm, GCF, error) {
	return RCFTerm{Kind: RCFEOF}, g, nil
}

func (g binaryGCF) Range() Range {
	return Range{}
}

func (g binaryGCF) Take(n int) (GCF, error) {
	return nil, ErrUndefined
}

func (g binaryGCF) Convergent() (Rational, error) {
	return Rational{}, ErrUndefined
}

func pqFromRCFTerm(term RCFTerm) (PQTerm, error) {
	if term.IsEOF() {
		return PQTerm{Kind: PQEOF}, nil
	}
	if !term.IsValue() || term.Value == nil {
		return PQTerm{}, ErrUndefined
	}

	return PQTerm{
		Kind: PQValue,
		P:    big.NewInt(1),
		Q:    new(big.Int).Set(term.Value),
	}, nil
}

func (g binaryGCF) ingestLeftStep() (binaryGCF, error) {
	term, rest, err := g.left.Next()
	if err != nil {
		return g, err
	}

	pq, err := pqFromRCFTerm(term)
	if err != nil {
		return g, err
	}
	if pq.IsEOF() {
		return g, ErrStuck
	}

	return binaryGCF{
		op:    g.op.ingestLeft(pq.P, pq.Q),
		left:  rest,
		right: g.right,
	}, nil
}

func (g binaryGCF) ingestRightStep() (binaryGCF, error) {
	term, rest, err := g.right.Next()
	if err != nil {
		return g, err
	}

	pq, err := pqFromRCFTerm(term)
	if err != nil {
		return g, err
	}
	if pq.IsEOF() {
		return g, ErrStuck
	}

	return binaryGCF{
		op:    g.op.ingestRight(pq.P, pq.Q),
		left:  g.left,
		right: rest,
	}, nil
}

func additionBinaryLFT() binaryLFT {
	return binaryLFT{
		a: big.NewInt(0),
		b: big.NewInt(1),
		c: big.NewInt(1),
		d: big.NewInt(0),
		e: big.NewInt(0),
		f: big.NewInt(0),
		g: big.NewInt(0),
		h: big.NewInt(1),
	}
}

func subtractionBinaryLFT() binaryLFT {
	return binaryLFT{
		a: big.NewInt(0),
		b: big.NewInt(1),
		c: big.NewInt(-1),
		d: big.NewInt(0),
		e: big.NewInt(0),
		f: big.NewInt(0),
		g: big.NewInt(0),
		h: big.NewInt(1),
	}
}

func multiplicationBinaryLFT() binaryLFT {
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

func divisionBinaryLFT() binaryLFT {
	return binaryLFT{
		a: big.NewInt(0),
		b: big.NewInt(1),
		c: big.NewInt(0),
		d: big.NewInt(0),
		e: big.NewInt(0),
		f: big.NewInt(0),
		g: big.NewInt(1),
		h: big.NewInt(0),
	}
}

// cf/binary.go v2
