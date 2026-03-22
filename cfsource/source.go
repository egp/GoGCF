// cfsource/source.go v2
package cfsource

import (
	"math/big"

	"github.com/egp/GoGCF/cf"
)

type int64Source struct {
	value   int64
	emitted bool
}

type rat64Source struct {
	terms []*big.Int
	index int
	num   int64
	den   int64
}

func Int64(v int64) cf.GCFSource {
	return int64Source{value: v}
}

func Rat64(num, den int64) cf.GCFSource {
	return rat64Source{
		terms: rcfTermsFromRat64(num, den),
		index: 0,
		num:   num,
		den:   den,
	}
}

func (s int64Source) NextPQ() (cf.PQTerm, cf.GCFSource, error) {
	if s.emitted {
		return cf.PQTerm{Kind: cf.PQEOF}, s, nil
	}

	term := cf.PQTerm{
		Kind: cf.PQValue,
		P:    big.NewInt(1),
		Q:    big.NewInt(s.value),
	}

	return term, int64Source{value: s.value, emitted: true}, nil
}

func (s int64Source) CurrentRange() cf.Range {
	if !s.emitted {
		return exactRationalRange(s.value, 1)
	}
	return cf.Range{}
}

func (s rat64Source) NextPQ() (cf.PQTerm, cf.GCFSource, error) {
	if s.index >= len(s.terms) {
		return cf.PQTerm{Kind: cf.PQEOF}, s, nil
	}

	term := cf.PQTerm{
		Kind: cf.PQValue,
		P:    big.NewInt(1),
		Q:    new(big.Int).Set(s.terms[s.index]),
	}

	return term, rat64Source{
		terms: s.terms,
		index: s.index + 1,
		num:   s.num,
		den:   s.den,
	}, nil
}

func (s rat64Source) CurrentRange() cf.Range {
	r, ok := convergentFromTerms(s.terms, s.index)
	if !ok {
		return cf.Range{}
	}
	return exactRangeFromRational(r)
}

func convergentFromTerms(terms []*big.Int, index int) (cf.Rational, bool) {
	if index < 0 || index >= len(terms) {
		return cf.Rational{}, false
	}

	num := new(big.Int).Set(terms[len(terms)-1])
	den := big.NewInt(1)

	for i := len(terms) - 2; i >= index; i-- {
		nextNum := new(big.Int).Mul(terms[i], num)
		nextNum.Add(nextNum, den)
		num, den = nextNum, num
	}

	return cf.Rational{
		Num: num,
		Den: den,
	}, true
}

func exactRangeFromRational(r cf.Rational) cf.Range {
	b := cf.Bound{
		Kind:   cf.BoundFinite,
		Value:  r,
		Closed: true,
	}
	return cf.Range{
		Lo:      b,
		Hi:      b,
		Outside: false,
	}
}

func exactRationalRange(num, den int64) cf.Range {
	n := big.NewInt(num)
	d := big.NewInt(den)
	if d.Sign() < 0 {
		n.Neg(n)
		d.Neg(d)
	}

	return exactRangeFromRational(cf.Rational{
		Num: n,
		Den: d,
	})
}

func rcfTermsFromRat64(num, den int64) []*big.Int {
	if den == 0 {
		return nil
	}

	n := num
	d := den
	if d < 0 {
		n = -n
		d = -d
	}

	terms := make([]*big.Int, 0, 8)
	for d != 0 {
		q := floorDivInt64(n, d)
		terms = append(terms, big.NewInt(q))

		r := n - q*d
		if r == 0 {
			break
		}

		n, d = d, r
	}

	return terms
}

func floorDivInt64(n, d int64) int64 {
	q := n / d
	r := n % d
	if r != 0 && ((r > 0) != (d > 0)) {
		q--
	}
	return q
}

// EOF cfsource/source.go v2
