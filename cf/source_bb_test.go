// cf/source_bb_test.go v1
package cf_test

import (
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

// cf/source_bb_test.go v1
