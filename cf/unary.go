// cf/unary.go v3
package cf

type unaryGCF struct {
	op       unaryLFT
	child    GCF
	resolved GCF
}

func (g unaryGCF) asStrategyUnary() strategyUnaryGCF {
	return strategyUnaryGCF{
		strategy: unaryLFTStrategy{op: g.op},
		child:    g.child,
		resolved: g.resolved,
	}
}

func unaryGCFfromStrategy(sg strategyUnaryGCF) unaryGCF {
	strategy, ok := sg.strategy.(unaryLFTStrategy)
	if !ok {
		return unaryGCF{}
	}
	return unaryGCF{
		op:       strategy.op,
		child:    sg.child,
		resolved: sg.resolved,
	}
}

func (g unaryGCF) Next() (RCFTerm, GCF, error) {
	term, rest, err := g.asStrategyUnary().Next()
	if err != nil {
		return RCFTerm{}, g, err
	}

	nextSG, ok := rest.(strategyUnaryGCF)
	if !ok {
		return RCFTerm{}, g, ErrUndefined
	}

	return term, unaryGCFfromStrategy(nextSG), nil
}

func (g unaryGCF) Range() Range {
	return g.asStrategyUnary().Range()
}

func (g unaryGCF) Take(n int) (GCF, error) {
	return g.asStrategyUnary().Take(n)
}

func (g unaryGCF) Convergent() (Rational, error) {
	return g.asStrategyUnary().Convergent()
}

// cf/unary.go v3
