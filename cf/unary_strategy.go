// cf/unary_strategy.go v2
package cf

import (
	"io"
	"math/big"
)

type unaryStrategy interface {
	RangeFromOperand(xr Range) (Range, error)
	EmitCandidateFromOperand(xr Range) (*big.Int, bool, error)
	Emit(term *big.Int) unaryStrategy
	ExactEval(x Rational) (Rational, bool, error)
}

type strategyUnaryGCF struct {
	strategy unaryStrategy
	child    GCF
	resolved GCF
}

func (g strategyUnaryGCF) Next() (RCFTerm, GCF, error) {
	if g.resolved != nil {
		term, rest, err := g.resolved.Next()
		if err != nil {
			return RCFTerm{}, g, err
		}
		return term, strategyUnaryGCF{
			strategy: g.strategy,
			child:    g.child,
			resolved: rest,
		}, nil
	}

	xr := g.child.Range()
	q, ok, err := g.strategy.EmitCandidateFromOperand(xr)
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
		z, exactOK, err := g.strategy.ExactEval(x)
		if err != nil {
			return RCFTerm{}, g, err
		}
		if exactOK {
			_, rem := floorDivModBig(z.Num, z.Den)
			if rem.Sign() == 0 {
				return term, strategyUnaryGCF{
					strategy: g.strategy,
					child:    g.child,
					resolved: rcfPrefixGCF{terms: nil, index: 0},
				}, nil
			}
		}
	}

	return term, strategyUnaryGCF{
		strategy: g.strategy.Emit(q),
		child:    g.child,
	}, nil
}

func (g strategyUnaryGCF) Range() Range {
	if g.resolved != nil {
		return g.resolved.Range()
	}

	zr, err := g.strategy.RangeFromOperand(g.child.Range())
	if err != nil {
		return Range{}
	}
	return zr
}

func (g strategyUnaryGCF) Take(n int) (GCF, error) {
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

func (g strategyUnaryGCF) Convergent() (Rational, error) {
	return Rational{}, ErrUndefined
}

// cf/unary_strategy.go v2
