package cf

import (
	"math/big"
	"testing"
)

type finiteGeneralizedSource struct {
	terms []PQTerm
	index int
}

func (s finiteGeneralizedSource) NextPQ() (PQTerm, GCFSource, error) {
	if s.index >= len(s.terms) {
		return PQTerm{Kind: PQEOF}, s, nil
	}

	term := s.terms[s.index]
	return PQTerm{
			Kind: PQValue,
			P:    new(big.Int).Set(term.P),
			Q:    new(big.Int).Set(term.Q),
		},
		finiteGeneralizedSource{
			terms: s.terms,
			index: s.index + 1,
		},
		nil
}

func TestBB_FromSource_FiniteGeneralizedSource_AdaptsToRegularCF(t *testing.T) {
	src := finiteGeneralizedSource{
		terms: []PQTerm{
			{Kind: PQValue, P: big.NewInt(1), Q: big.NewInt(1)},
			{Kind: PQValue, P: big.NewInt(2), Q: big.NewInt(3)},
		},
	}

	g := FromSource(src)

	prefix, err := g.Take(2)
	if err != nil {
		t.Fatalf("Take(2) error: %v", err)
	}

	cur := prefix
	for i, want := range []string{"1", "3"} {
		term, rest, err := cur.Next()
		if err != nil {
			t.Fatalf("term %d Next() error: %v", i, err)
		}
		if !term.IsValue() {
			t.Fatalf("term %d kind = %v, want value", i, term.Kind)
		}

		value, ok := term.BigInt()
		if !ok {
			t.Fatalf("term %d missing BigInt value", i)
		}
		if got := value.String(); got != want {
			t.Fatalf("term %d = %s, want %s", i, got, want)
		}

		cur = rest
	}

	eof, _, err := cur.Next()
	if err != nil {
		t.Fatalf("EOF Next() error: %v", err)
	}
	if !eof.IsEOF() {
		t.Fatalf("final term kind = %v, want EOF", eof.Kind)
	}

	r, err := prefix.Convergent()
	if err != nil {
		t.Fatalf("prefix.Convergent() error: %v", err)
	}
	if got, want := r.Num.String(), "4"; got != want {
		t.Fatalf("numerator = %s, want %s", got, want)
	}
	if got, want := r.Den.String(), "3"; got != want {
		t.Fatalf("denominator = %s, want %s", got, want)
	}
}
