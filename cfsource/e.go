// cfsource/e.go v2
package cfsource

import (
	"math/big"

	"github.com/egp/GoGCF/cf"
)

type eSource struct {
	index int
}

func E() cf.GCFSource {
	return eSource{index: 0}
}

func (s eSource) currentQ() int64 {
	switch {
	case s.index == 0:
		return 2
	case s.index%3 == 2:
		k := int64((s.index + 1) / 3)
		return 2 * k
	default:
		return 1
	}
}

func (s eSource) NextPQ() (cf.PQTerm, cf.GCFSource, error) {
	q := s.currentQ()

	term := cf.PQTerm{
		Kind: cf.PQValue,
		P:    big.NewInt(1),
		Q:    big.NewInt(q),
	}

	return term, eSource{index: s.index + 1}, nil
}

func (s eSource) CurrentRange() cf.Range {
	q := s.currentQ()

	return cf.Range{
		Lo: cf.Bound{
			Kind: cf.BoundFinite,
			Value: cf.Rational{
				Num: big.NewInt(q),
				Den: big.NewInt(1),
			},
			Closed: false,
		},
		Hi: cf.Bound{
			Kind: cf.BoundFinite,
			Value: cf.Rational{
				Num: big.NewInt(q + 1),
				Den: big.NewInt(1),
			},
			Closed: false,
		},
		Outside: false,
	}
}

// cfsource/e.go v2
