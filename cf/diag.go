// cf/diag.go v3
package cf

type diagGCF struct {
	op       diagLFT
	child    GCF
	resolved GCF
}

func (g diagGCF) asStrategyUnary() strategyUnaryGCF {
	return strategyUnaryGCF{
		strategy: diagLFTStrategy{op: g.op},
		child:    g.child,
		resolved: g.resolved,
	}
}

func diagGCFfromStrategy(sg strategyUnaryGCF) diagGCF {
	strategy, ok := sg.strategy.(diagLFTStrategy)
	if !ok {
		return diagGCF{}
	}
	return diagGCF{
		op:       strategy.op,
		child:    sg.child,
		resolved: sg.resolved,
	}
}

func (g diagGCF) Next() (RCFTerm, GCF, error) {
	term, rest, err := g.asStrategyUnary().Next()
	if err != nil {
		return RCFTerm{}, g, err
	}

	nextSG, ok := rest.(strategyUnaryGCF)
	if !ok {
		return RCFTerm{}, g, ErrUndefined
	}

	return term, diagGCFfromStrategy(nextSG), nil
}

func (g diagGCF) Range() Range {
	return g.asStrategyUnary().Range()
}

func (g diagGCF) Take(n int) (GCF, error) {
	return g.asStrategyUnary().Take(n)
}

func (g diagGCF) Convergent() (Rational, error) {
	return g.asStrategyUnary().Convergent()
}

// cf/diag.go v3
