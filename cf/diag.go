// cf/diag.go v2
package cf

import (
	"io"
	"math/big"
)

type diagGCF struct {
	op       diagLFT
	child    GCF
	resolved GCF
}

func (g diagGCF) Next() (RCFTerm, GCF, error) {
	if g.resolved != nil {
		term, rest, err := g.resolved.Next()
		if err != nil {
			return RCFTerm{}, g, err
		}
		return term, diagGCF{
			op:       g.op,
			child:    g.child,
			resolved: rest,
		}, nil
	}

	xr := g.child.Range()
	q, ok, err := g.op.emitCandidateFromRange(xr)
	if err != nil {
		return RCFTerm{}, g, err
	}
	if !ok {
		return RCFTerm{}, g, ErrUndefined
	}

	term := RCFTerm{
		Kind:  RCFValue,
		Value: new(big.Int).Set(q),
	}

	if x, exact := exactFiniteRangeValue(xr); exact {
		z, err := g.op.evalAt(x)
		if err != nil {
			return RCFTerm{}, g, err
		}
		_, rem := floorDivModBig(z.Num, z.Den)
		if rem.Sign() == 0 {
			return term, diagGCF{
				op:       g.op,
				child:    g.child,
				resolved: rcfPrefixGCF{terms: nil, index: 0},
			}, nil
		}
	}

	return term, diagGCF{
		op:    g.op.emit(q),
		child: g.child,
	}, nil
}

func (g diagGCF) Range() Range {
	if g.resolved != nil {
		return g.resolved.Range()
	}

	zr, err := g.op.rangeFromXRange(g.child.Range())
	if err != nil {
		return Range{}
	}
	return zr
}

func (g diagGCF) Take(n int) (GCF, error) {
	if n < 0 {
		return nil, ErrUndefined
	}

	terms := make([]*big.Int, 0, n)
	cur := GCF(g)

	for len(terms) < n {
		term, rest, err := cur.Next()
		if err != nil {
			return nil, err
		}

		if term.IsEOF() {
			return rcfPrefixGCF{terms: terms, index: 0}, io.EOF
		}

		value, ok := term.BigInt()
		if !ok {
			return nil, ErrUndefined
		}

		terms = append(terms, new(big.Int).Set(value))
		cur = rest
	}

	return rcfPrefixGCF{terms: terms, index: 0}, nil
}

func (g diagGCF) Convergent() (Rational, error) {
	return Rational{}, ErrUndefined
}

// cf/diag.go v2
