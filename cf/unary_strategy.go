// cf/unary_strategy.go v5
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

type ouroborosUnaryStrategy interface {
	OuroborosFeedback(xr Range) (Range, *big.Int, unaryStrategy, bool, error)
}

type childIngestingUnaryStrategy interface {
	IngestChild(term *big.Int) (unaryStrategy, error)
}

type childExactUnaryStrategy interface {
	ExactFromChild(child GCF) (Rational, bool, error)
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

	if resolver, ok := g.strategy.(childExactUnaryStrategy); ok {
		r, exact, err := resolver.ExactFromChild(g.child)
		if err != nil {
			return RCFTerm{}, g, err
		}
		if exact {
			prefix, err := prefixGCFfromRational(r)
			if err != nil {
				return RCFTerm{}, g, err
			}
			term, rest, err := prefix.Next()
			if err != nil {
				return RCFTerm{}, g, err
			}
			return term, strategyUnaryGCF{
				strategy: g.strategy,
				child:    g.child,
				resolved: rest,
			}, nil
		}
	}

	cur := g
	for {
		xr := cur.child.Range()

		strategyForEmit := cur.strategy
		var q *big.Int
		var ok bool
		var err error

		if ouro, has := cur.strategy.(ouroborosUnaryStrategy); has {
			_, q, strategyForEmit, ok, err = ouro.OuroborosFeedback(xr)
			if err != nil {
				return RCFTerm{}, g, err
			}
		}

		if !ok {
			q, ok, err = cur.strategy.EmitCandidateFromOperand(xr)
			if err != nil {
				return RCFTerm{}, g, err
			}
			if ok {
				strategyForEmit = cur.strategy
			}
		}

		if ok {
			term := RCFTerm{
				Kind:  RCFValue,
				Value: new(big.Int).Set(q),
			}

			if x, exact := exactFiniteRangeValue(xr); exact {
				z, exactOK, err := strategyForEmit.ExactEval(x)
				if err != nil {
					return RCFTerm{}, g, err
				}
				if exactOK {
					_, rem := floorDivModBig(z.Num, z.Den)
					if rem.Sign() == 0 {
						return term, strategyUnaryGCF{
							strategy: strategyForEmit,
							child:    cur.child,
							resolved: rcfPrefixGCF{terms: nil, index: 0},
						}, nil
					}
				}
			}

			return term, strategyUnaryGCF{
				strategy: strategyForEmit.Emit(q),
				child:    cur.child,
			}, nil
		}

		ingestor, has := cur.strategy.(childIngestingUnaryStrategy)
		if !has {
			return RCFTerm{}, g, ErrUndefined
		}

		childTerm, childRest, err := cur.child.Next()
		if err != nil {
			return RCFTerm{}, g, err
		}
		if childTerm.IsEOF() {
			return RCFTerm{}, g, ErrUndefined
		}

		value, ok := childTerm.BigInt()
		if !ok {
			return RCFTerm{}, g, ErrUndefined
		}

		nextStrategy, err := ingestor.IngestChild(value)
		if err != nil {
			return RCFTerm{}, g, err
		}

		cur = strategyUnaryGCF{
			strategy: nextStrategy,
			child:    childRest,
		}
	}
}

func (g strategyUnaryGCF) Range() Range {
	if g.resolved != nil {
		return g.resolved.Range()
	}

	if resolver, ok := g.strategy.(childExactUnaryStrategy); ok {
		r, exact, err := resolver.ExactFromChild(g.child)
		if err == nil && exact {
			return exactRangeFromRational(r)
		}
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

// cf/unary_strategy.go v5
