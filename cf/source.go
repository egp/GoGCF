// cf/source.go v3
package cf

import (
	"errors"
	"math/big"
)

var ErrUndefined = errors.New("undefined")
var ErrStuck = errors.New("stuck")

type int64Source struct {
	value   int64
	emitted bool
}

func Int64(v int64) GCFSource {
	return int64Source{value: v}
}

func (s int64Source) NextPQ() (PQTerm, GCFSource, error) {
	if s.emitted {
		return PQTerm{Kind: PQEOF}, s, nil
	}

	term := PQTerm{
		Kind: PQValue,
		P:    big.NewInt(1),
		Q:    big.NewInt(s.value),
	}

	return term, int64Source{value: s.value, emitted: true}, nil
}

type sourceBackedGCF struct {
	src GCFSource
}

func FromSource(src GCFSource) GCF {
	return sourceBackedGCF{src: src}
}

func (g sourceBackedGCF) Next() (RCFTerm, GCF, error) {
	pq, rest, err := g.src.NextPQ()
	if err != nil {
		return RCFTerm{}, g, err
	}

	if pq.IsEOF() {
		return RCFTerm{Kind: RCFEOF}, sourceBackedGCF{src: rest}, nil
	}

	if pq.P == nil || pq.Q == nil {
		return RCFTerm{}, g, ErrUndefined
	}
	if pq.P.Cmp(big.NewInt(1)) != 0 {
		return RCFTerm{}, g, ErrUndefined
	}

	return RCFTerm{
		Kind:  RCFValue,
		Value: new(big.Int).Set(pq.Q),
	}, sourceBackedGCF{src: rest}, nil
}

func (g sourceBackedGCF) Range() Range {
	switch src := g.src.(type) {
	case int64Source:
		if !src.emitted {
			return exactInt64Range(src.value)
		}
	}

	return Range{}
}

func (g sourceBackedGCF) Take(n int) (GCF, error) {
	return nil, ErrUndefined
}

func (g sourceBackedGCF) Convergent() (Rational, error) {
	return Rational{}, ErrUndefined
}

func exactInt64Range(v int64) Range {
	r := Rational{
		Num: big.NewInt(v),
		Den: big.NewInt(1),
	}

	b := Bound{
		Kind:   BoundFinite,
		Value:  r,
		Closed: true,
	}

	return Range{
		Lo:      b,
		Hi:      b,
		Outside: false,
	}
}

// cf/source.go v3
