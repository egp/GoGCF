// cf/binary.go v12
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
	if g.resolved != nil {
		term, rest, err := g.resolved.Next()
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

	cur := g
	for {
		action, next, err := cur.step()
		if err != nil {
			if err != ErrStuck {
				return RCFTerm{}, g, err
			}

			prefix, collapseErr := cur.completeToPrefix()
			if collapseErr != nil {
				return RCFTerm{}, g, collapseErr
			}

			term, rest, nextErr := prefix.Next()
			if nextErr != nil {
				return RCFTerm{}, g, nextErr
			}

			return term, binaryGCF{
				op:       cur.op,
				left:     cur.left,
				right:    cur.right,
				resolved: rest,
			}, nil
		}

		switch action {
		case decisionEmitOutput:
			term, emittedNext, emitErr := cur.emitStep()
			if emitErr != nil {
				return RCFTerm{}, g, emitErr
			}
			return term, emittedNext, nil

		case decisionIngestLeft, decisionIngestRight:
			cur = next

		default:
			return RCFTerm{}, g, ErrStuck
		}
	}
}

func (g binaryGCF) Range() Range {
	if g.resolved != nil {
		return g.resolved.Range()
	}

	canLeft, errLeft := canIngestFromChild(g.left)
	canRight, errRight := canIngestFromChild(g.right)

	if errLeft == nil && errRight == nil && !canLeft && !canRight {
		if r, ok := g.op.exactQuotient(); ok {
			return exactRangeFromRational(r)
		}
	}

	zr, err := g.op.rangeFromXYRanges(g.left.Range(), g.right.Range())
	if err != nil {
		return Range{}
	}
	return zr
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
func (g binaryGCF) emitStep() (RCFTerm, binaryGCF, error) {
	xr := g.left.Range()
	yr := g.right.Range()

	if q, ok, err := g.op.emitCandidateFromRanges(xr, yr); err != nil {
		return RCFTerm{}, g, err
	} else if ok {
		zr, err := g.op.rangeFromXYRanges(xr, yr)
		if err != nil {
			return RCFTerm{}, g, err
		}

		term := RCFTerm{
			Kind:  RCFValue,
			Value: new(big.Int).Set(q),
		}

		if z, exact := exactFiniteRangeValue(zr); exact {
			_, rem := floorDivModBig(z.Num, z.Den)
			if rem.Sign() == 0 {
				return term, binaryGCF{
					op:       g.op,
					left:     g.left,
					right:    g.right,
					resolved: rcfPrefixGCF{terms: nil, index: 0},
				}, nil
			}
		}

		return term, binaryGCF{
			op:    g.op.emit(q),
			left:  g.left,
			right: g.right,
		}, nil
	}

	r, ok := g.op.exactQuotient()
	if !ok {
		return RCFTerm{}, g, ErrStuck
	}

	q, rem := floorDivModBig(r.Num, r.Den)

	term := RCFTerm{
		Kind:  RCFValue,
		Value: new(big.Int).Set(q),
	}

	if rem.Sign() == 0 {
		return term, binaryGCF{
			op:       g.op,
			left:     g.left,
			right:    g.right,
			resolved: rcfPrefixGCF{terms: nil, index: 0},
		}, nil
	}

	return term, binaryGCF{
		op:    g.op.emit(q),
		left:  g.left,
		right: g.right,
	}, nil
}

func (g binaryGCF) step() (binaryDecision, binaryGCF, error) {
	xr := g.left.Range()
	yr := g.right.Range()

	canEmitOutput := false
	if _, ok, err := g.op.emitCandidateFromRanges(xr, yr); err != nil {
		return 0, g, err
	} else {
		canEmitOutput = ok
	}

	canLeft, err := canIngestFromChild(g.left)
	if err != nil {
		return 0, g, err
	}

	canRight, err := canIngestFromChild(g.right)
	if err != nil {
		return 0, g, err
	}

	if !canEmitOutput && canLeft && canRight {
		leftNext, err := g.ingestLeftStep()
		if err != nil {
			return 0, g, err
		}
		rightNext, err := g.ingestRightStep()
		if err != nil {
			return 0, g, err
		}

		switch leftNext.Range().Cmp(rightNext.Range()) {
		case -1:
			return decisionIngestLeft, leftNext, nil
		case 1:
			return decisionIngestRight, rightNext, nil
		default:
			return decisionIngestLeft, leftNext, nil
		}
	}

	if !canEmitOutput && !canLeft && !canRight {
		_, canEmitOutput = g.op.exactQuotient()
	}

	state := binaryStepState{
		canEmitOutput:  canEmitOutput,
		canIngestLeft:  canLeft,
		canIngestRight: canRight,
	}

	action, err := state.choose()
	if err != nil {
		return 0, g, err
	}

	switch action {
	case decisionEmitOutput:
		_, next, err := g.emitStep()
		if err != nil {
			return 0, g, err
		}
		return decisionEmitOutput, next, nil

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

func (g binaryGCF) ingestOnlyStep() (binaryGCF, error) {
	canLeft, err := canIngestFromChild(g.left)
	if err != nil {
		return g, err
	}
	if canLeft {
		return g.ingestLeftStep()
	}

	canRight, err := canIngestFromChild(g.right)
	if err != nil {
		return g, err
	}
	if canRight {
		return g.ingestRightStep()
	}

	return g, ErrStuck
}

func (g binaryGCF) completeToRational() (Rational, error) {
	cur := g
	for {
		next, err := cur.ingestOnlyStep()
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

// cf/binary.go v12
