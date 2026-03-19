// cf/sqrt2_source_bb_test.go v1
package cf_test

import (
	"testing"

	"github.com/egp/GoGCF/cf"
)

func TestBB_Sqrt2_FirstStepIsValue(t *testing.T) {
	src := cf.Sqrt2()

	term, _, err := src.NextPQ()
	if err != nil {
		t.Fatalf("NextPQ() error: %v", err)
	}
	if !term.IsValue() {
		t.Fatalf("first sqrt2 term should be a value, got EOF")
	}
}

func TestBB_Sqrt2_IsImmutableAtFirstStep(t *testing.T) {
	src := cf.Sqrt2()

	t1, _, err := src.NextPQ()
	if err != nil {
		t.Fatalf("first NextPQ() error: %v", err)
	}

	t2, _, err := src.NextPQ()
	if err != nil {
		t.Fatalf("second NextPQ() on original source error: %v", err)
	}

	if t1.IsEOF() || t2.IsEOF() {
		t.Fatalf("expected value terms on repeated first-step reads, got EOF")
	}

	if t1.P.String() != t2.P.String() || t1.Q.String() != t2.Q.String() {
		t.Fatalf("immutable source violated: first-step terms differ")
	}
}

func TestBB_Sqrt2_RemainderAlsoProducesValue(t *testing.T) {
	src := cf.Sqrt2()

	_, rest, err := src.NextPQ()
	if err != nil {
		t.Fatalf("first NextPQ() error: %v", err)
	}

	term2, _, err := rest.NextPQ()
	if err != nil {
		t.Fatalf("remainder NextPQ() error: %v", err)
	}
	if !term2.IsValue() {
		t.Fatalf("second sqrt2 term should be a value, got EOF")
	}
}

// cf/sqrt2_source_bb_test.go v1
