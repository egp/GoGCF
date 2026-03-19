// cf/source_bb_test.go v1
package cf_test

import (
	"io"
	"testing"

	"github.com/egp/GoGCF/cf"
)

func TestBB_Int64_SourceEmitsSinglePQValueThenEOF(t *testing.T) {
	src := cf.Int64(7)

	first, rest, err := src.NextPQ()
	if err != nil {
		t.Fatalf("first NextPQ() error: %v", err)
	}
	if !first.IsValue() {
		t.Fatalf("first Int64 term should be a value")
	}
	if first.P == nil || first.Q == nil {
		t.Fatalf("first Int64 term should have non-nil P and Q")
	}
	if got, want := first.P.String(), "1"; got != want {
		t.Fatalf("first Int64 P = %s, want %s", got, want)
	}
	if got, want := first.Q.String(), "7"; got != want {
		t.Fatalf("first Int64 Q = %s, want %s", got, want)
	}

	second, _, err := rest.NextPQ()
	if err != nil {
		t.Fatalf("second NextPQ() error: %v", err)
	}
	if !second.IsEOF() {
		t.Fatalf("second Int64 term should be EOF")
	}
}

func TestBB_FromSource_Int64_FirstNextEmitsRCFValue(t *testing.T) {
	g := cf.FromSource(cf.Int64(7))

	term, _, err := g.Next()
	if err != nil {
		t.Fatalf("Next() error: %v", err)
	}
	if !term.IsValue() {
		t.Fatalf("first emitted term should be a value")
	}
	value, ok := term.BigInt()
	if !ok {
		t.Fatalf("first emitted term should expose a BigInt")
	}
	if got, want := value.String(), "7"; got != want {
		t.Fatalf("first emitted value = %s, want %s", got, want)
	}
}

func TestBB_FromSource_Int64_SecondNextEmitsEOF(t *testing.T) {
	g := cf.FromSource(cf.Int64(7))

	_, rest, err := g.Next()
	if err != nil {
		t.Fatalf("first Next() error: %v", err)
	}

	term, _, err := rest.Next()
	if err != nil {
		t.Fatalf("second Next() error: %v", err)
	}
	if !term.IsEOF() {
		t.Fatalf("second emitted term should be EOF")
	}
}

func TestBB_FromSource_Int64_TakeOneReturnsFinitePrefix(t *testing.T) {
	g := cf.FromSource(cf.Int64(7))

	prefix, err := g.Take(1)
	if err != nil {
		t.Fatalf("Take(1) error: %v", err)
	}

	term1, rest, err := prefix.Next()
	if err != nil {
		t.Fatalf("prefix first Next() error: %v", err)
	}
	if !term1.IsValue() {
		t.Fatalf("prefix first term should be a value")
	}
	value1, ok := term1.BigInt()
	if !ok {
		t.Fatalf("prefix first term should expose a BigInt")
	}
	if got, want := value1.String(), "7"; got != want {
		t.Fatalf("prefix first value = %s, want %s", got, want)
	}

	term2, _, err := rest.Next()
	if err != nil {
		t.Fatalf("prefix second Next() error: %v", err)
	}
	if !term2.IsEOF() {
		t.Fatalf("prefix second term should be EOF")
	}
}

func TestBB_FromSource_Int64_TakeThreeReturnsEOFAndShortFinitePrefix(t *testing.T) {
	g := cf.FromSource(cf.Int64(7))

	prefix, err := g.Take(3)
	if err != io.EOF {
		t.Fatalf("Take(3) error = %v, want io.EOF", err)
	}

	term1, rest, err := prefix.Next()
	if err != nil {
		t.Fatalf("prefix first Next() error: %v", err)
	}
	if !term1.IsValue() {
		t.Fatalf("prefix first term should be a value")
	}
	value1, ok := term1.BigInt()
	if !ok {
		t.Fatalf("prefix first term should expose a BigInt")
	}
	if got, want := value1.String(), "7"; got != want {
		t.Fatalf("prefix first value = %s, want %s", got, want)
	}

	term2, _, err := rest.Next()
	if err != nil {
		t.Fatalf("prefix second Next() error: %v", err)
	}
	if !term2.IsEOF() {
		t.Fatalf("prefix second term should be EOF")
	}
}

func TestBB_FromSource_Sqrt2_TakeThreeReturnsFinitePrefix_1_2_2(t *testing.T) {
	g := cf.FromSource(cf.Sqrt2())

	prefix, err := g.Take(3)
	if err != nil {
		t.Fatalf("Take(3) error: %v", err)
	}

	want := []string{"1", "2", "2"}
	cur := prefix
	for i, w := range want {
		term, rest, err := cur.Next()
		if err != nil {
			t.Fatalf("prefix Next() #%d error: %v", i+1, err)
		}
		if !term.IsValue() {
			t.Fatalf("prefix term #%d should be a value", i+1)
		}
		value, ok := term.BigInt()
		if !ok {
			t.Fatalf("prefix term #%d should expose a BigInt", i+1)
		}
		if got := value.String(); got != w {
			t.Fatalf("prefix term #%d = %s, want %s", i+1, got, w)
		}
		cur = rest
	}

	term4, _, err := cur.Next()
	if err != nil {
		t.Fatalf("prefix final Next() error: %v", err)
	}
	if !term4.IsEOF() {
		t.Fatalf("prefix final term should be EOF")
	}
}

// cf/source_bb_test.go v1
