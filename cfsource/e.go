// cfsource/e.go v1
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

func (s eSource) NextPQ() (cf.PQTerm, cf.GCFSource, error) {
	var q int64
	switch {
	case s.index == 0:
		q = 2
	case s.index%3 == 2:
		k := int64((s.index + 1) / 3)
		q = 2 * k
	default:
		q = 1
	}

	term := cf.PQTerm{
		Kind: cf.PQValue,
		P:    big.NewInt(1),
		Q:    big.NewInt(q),
	}

	return term, eSource{index: s.index + 1}, nil
}

func (s eSource) CurrentRange() cf.Range {
	switch s.index {
	case 0:
		return cf.Range{
			Lo: cf.Bound{
				Kind: cf.BoundFinite,
				Value: cf.Rational{
					Num: big.NewInt(2),
					Den: big.NewInt(1),
				},
				Closed: false,
			},
			Hi: cf.Bound{
				Kind: cf.BoundFinite,
				Value: cf.Rational{
					Num: big.NewInt(3),
					Den: big.NewInt(1),
				},
				Closed: false,
			},
			Outside: false,
		}
	default:
		return cf.Range{
			Lo: cf.Bound{
				Kind: cf.BoundFinite,
				Value: cf.Rational{
					Num: big.NewInt(1),
					Den: big.NewInt(1),
				},
				Closed: false,
			},
			Hi: cf.Bound{
				Kind: cf.BoundFinite,
				Value: cf.Rational{
					Num: big.NewInt(2),
					Den: big.NewInt(1),
				},
				Closed: false,
			},
			Outside: false,
		}
	}
}

// EOF cfsource/e.go v1
