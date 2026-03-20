// cf/binary.go v9
package cf

import "math/big"

type binaryGCF struct {
	op       binaryLFT
	left     GCF
	right    GCF
	resolved GCF
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
	resolved := g.resolved
	if resolved == nil {
		prefix, err := g.completeToPrefix()
		if err != nil {
			return RCFTerm{}, g, err
		}
		resolved = prefix
	}

	term, rest, err := resolved.Next()
	if err != nil {
		return RCFTerm{}, g, err
	}

	return term, binaryGCF{
		op:       g.op,
		left:     g.left,
		right:    g.right,
		resolved: rest,
	}, nil
}

func (g binaryGCF) Range() Range {
	return Range{}
}

func (g binaryGCF) Take(n int) (GCF, error) {
	prefix, err := g.completeToPrefix()
	if err != nil {
		return nil, err
	}
	return prefix.Take(n)
}

func (g binaryGCF) Convergent() (Rational, error) {
	return g.completeToRational()
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

func (g binaryGCF) step() (binaryDecision, binaryGCF, error) {
	canLeft, err := canIngestFromChild(g.left)
	if err != nil {
		return 0, g, err
	}

	canRight, err := canIngestFromChild(g.right)
	if err != nil {
		return 0, g, err
	}

	state := binaryStepState{
		canEmitOutput:  false,
		canIngestLeft:  canLeft,
		canIngestRight: canRight,
	}

	action, err := state.choose()
	if err != nil {
		return 0, g, err
	}

	switch action {
	case decisionIngestLeft:
		next, err := g.ingestLeftStep()
		if err != nil {
			return 0, g, err
		}
		return decisionIngestLeft, next, nil
	case decisionIngestRight:
		next, err := g.ingestRightStep()
		if err != nil {
			return 0, g, err
		}
		return decisionIngestRight, next, nil
	default:
		return 0, g, ErrStuck
	}
}

func (g binaryGCF) completeToRational() (Rational, error) {
	cur := g
	for {
		_, next, err := cur.step()
		if err == nil {
			cur = next
			continue
		}
		if err != ErrStuck {
			return Rational{}, err
		}

		u := cur.op.collapseRightEOF()
		return u.collapseEOFToRational()
	}
}

func (g binaryGCF) completeToPrefix() (GCF, error) {
	r, err := g.completeToRational()
	if err != nil {
		return nil, err
	}
	return prefixGCFfromRational(r)
}

func canIngestFromChild(child GCF) (bool, error) {
	term, _, err := child.Next()
	if err != nil {
		return false, err
	}
	return !term.IsEOF(), nil
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
		b: big.NewInt(-1),
		c: big.NewInt(1),
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
		b: big.NewInt(0),
		c: big.NewInt(1),
		d: big.NewInt(0),
		e: big.NewInt(0),
		f: big.NewInt(1),
		g: big.NewInt(0),
		h: big.NewInt(0),
	}
}

// cf/binary.go v9
