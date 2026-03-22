// cfsource/sqrt2.go v1
package cfsource

import (
	"math/big"

	"github.com/egp/GoGCF/cf"
)

type sqrt2Source struct {
	index int
}

func Sqrt2() cf.GCFSource {
	return sqrt2Source{index: 0}
}

func (s sqrt2Source) NextPQ() (cf.PQTerm, cf.GCFSource, error) {
	q := int64(2)
	if s.index == 0 {
		q = 1
	}

	term := cf.PQTerm{
		Kind: cf.PQValue,
		P:    big.NewInt(1),
		Q:    big.NewInt(q),
	}

	return term, sqrt2Source{index: s.index + 1}, nil
}

func (s sqrt2Source) CurrentRange() cf.Range {
	switch s.index {
	case 0:
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
					Num: big.NewInt(3),
					Den: big.NewInt(2),
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
					Num: big.NewInt(2),
					Den: big.NewInt(1),
				},
				Closed: false,
			},
			Hi: cf.Bound{
				Kind: cf.BoundFinite,
				Value: cf.Rational{
					Num: big.NewInt(5),
					Den: big.NewInt(2),
				},
				Closed: false,
			},
			Outside: false,
		}
	}
}
