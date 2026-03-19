// cf/source.go v4
package cf

import (
	"errors"
	"io"
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

type rcfPrefixGCF struct {
	terms []*big.Int
	index int
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

func (g sourceBackedGCF) Convergent() (Rational, error) {
	return Rational{}, ErrUndefined
}

func (g rcfPrefixGCF) Next() (RCFTerm, GCF, error) {
	if g.index >= len(g.terms) {
		return RCFTerm{Kind: RCFEOF}, g, nil
	}

	term := RCFTerm{
		Kind:  RCFValue,
		Value: new(big.Int).Set(g.terms[g.index]),
	}

	return term, rcfPrefixGCF{
		terms: g.terms,
		index: g.index + 1,
	}, nil
}

func (g rcfPrefixGCF) Range() Range {
	return Range{}
}

func (g rcfPrefixGCF) Take(n int) (GCF, error) {
	if n < 0 {
		return nil, ErrUndefined
	}

	remaining := len(g.terms) - g.index
	if n > remaining {
		n = remaining
	}

	prefixTerms := make([]*big.Int, 0, n)
	for i := 0; i < n; i++ {
		prefixTerms = append(prefixTerms, new(big.Int).Set(g.terms[g.index+i]))
	}

	prefix := rcfPrefixGCF{
		terms: prefixTerms,
		index: 0,
	}

	if n < remaining {
		return prefix, nil
	}
	return prefix, io.EOF
}

func (g rcfPrefixGCF) Convergent() (Rational, error) {
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

// cf/source.go v4
