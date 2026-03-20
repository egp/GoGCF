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

func TestBB_Sqrt2_FirstTermIsPQ_1_1(t *testing.T) {
	src := cf.Sqrt2()

	term, _, err := src.NextPQ()
	if err != nil {
		t.Fatalf("NextPQ() error: %v", err)
	}
	if !term.IsValue() {
		t.Fatalf("first sqrt2 term should be a value")
	}
	if term.P == nil || term.Q == nil {
		t.Fatalf("first sqrt2 term should have non-nil P and Q")
	}
	if got, want := term.P.String(), "1"; got != want {
		t.Fatalf("first sqrt2 P = %s, want %s", got, want)
	}
	if got, want := term.Q.String(), "1"; got != want {
		t.Fatalf("first sqrt2 Q = %s, want %s", got, want)
	}
}

func TestBB_Sqrt2_SecondAndThirdTermsArePQ_1_2(t *testing.T) {
	src := cf.Sqrt2()

	_, rest1, err := src.NextPQ()
	if err != nil {
		t.Fatalf("first NextPQ() error: %v", err)
	}

	second, rest2, err := rest1.NextPQ()
	if err != nil {
		t.Fatalf("second NextPQ() error: %v", err)
	}
	if !second.IsValue() {
		t.Fatalf("second sqrt2 term should be a value")
	}
	if got, want := second.P.String(), "1"; got != want {
		t.Fatalf("second sqrt2 P = %s, want %s", got, want)
	}
	if got, want := second.Q.String(), "2"; got != want {
		t.Fatalf("second sqrt2 Q = %s, want %s", got, want)
	}

	third, _, err := rest2.NextPQ()
	if err != nil {
		t.Fatalf("third NextPQ() error: %v", err)
	}
	if !third.IsValue() {
		t.Fatalf("third sqrt2 term should be a value")
	}
	if got, want := third.P.String(), "1"; got != want {
		t.Fatalf("third sqrt2 P = %s, want %s", got, want)
	}
	if got, want := third.Q.String(), "2"; got != want {
		t.Fatalf("third sqrt2 Q = %s, want %s", got, want)
	}
}
func TestBB_FromSource_Sqrt2_RangeIsInsideOpen_OneToThreeHalves(t *testing.T) {
	g := cf.FromSource(cf.Sqrt2())

	r := g.Range()

	if !r.IsInside() {
		t.Fatalf("sqrt2 range should be inside")
	}
	if r.IsOutside() {
		t.Fatalf("sqrt2 range should not be outside")
	}
	if r.IsExact() {
		t.Fatalf("sqrt2 range should not be exact")
	}
	if r.Lo.Kind != cf.BoundFinite || r.Hi.Kind != cf.BoundFinite {
		t.Fatalf("sqrt2 range should have finite bounds")
	}
	if r.Lo.Closed || r.Hi.Closed {
		t.Fatalf("sqrt2 initial bounds should be open")
	}

	if got, want := r.Lo.Value.Num.String(), "1"; got != want {
		t.Fatalf("lo numerator = %s, want %s", got, want)
	}
	if got, want := r.Lo.Value.Den.String(), "1"; got != want {
		t.Fatalf("lo denominator = %s, want %s", got, want)
	}
	if got, want := r.Hi.Value.Num.String(), "3"; got != want {
		t.Fatalf("hi numerator = %s, want %s", got, want)
	}
	if got, want := r.Hi.Value.Den.String(), "2"; got != want {
		t.Fatalf("hi denominator = %s, want %s", got, want)
	}
}

func TestBB_FromSource_Sqrt2_RangeAfterOneNext_IsInsideOpen_TwoToFiveHalves(t *testing.T) {
	g := cf.FromSource(cf.Sqrt2())

	term, rest, err := g.Next()
	if err != nil {
		t.Fatalf("first Next() error: %v", err)
	}
	if !term.IsValue() {
		t.Fatalf("first term should be a value")
	}
	value, ok := term.BigInt()
	if !ok {
		t.Fatalf("first term should expose a BigInt")
	}
	if got, want := value.String(), "1"; got != want {
		t.Fatalf("first term = %s, want %s", got, want)
	}

	r := rest.Range()

	if !r.IsInside() {
		t.Fatalf("remainder sqrt2-tail range should be inside")
	}
	if r.IsOutside() {
		t.Fatalf("remainder sqrt2-tail range should not be outside")
	}
	if r.IsExact() {
		t.Fatalf("remainder sqrt2-tail range should not be exact")
	}
	if r.Lo.Kind != cf.BoundFinite || r.Hi.Kind != cf.BoundFinite {
		t.Fatalf("remainder sqrt2-tail range should have finite bounds")
	}
	if r.Lo.Closed || r.Hi.Closed {
		t.Fatalf("remainder sqrt2-tail bounds should be open")
	}

	if got, want := r.Lo.Value.Num.String(), "2"; got != want {
		t.Fatalf("lo numerator = %s, want %s", got, want)
	}
	if got, want := r.Lo.Value.Den.String(), "1"; got != want {
		t.Fatalf("lo denominator = %s, want %s", got, want)
	}
	if got, want := r.Hi.Value.Num.String(), "5"; got != want {
		t.Fatalf("hi numerator = %s, want %s", got, want)
	}
	if got, want := r.Hi.Value.Den.String(), "2"; got != want {
		t.Fatalf("hi denominator = %s, want %s", got, want)
	}
}

// cf/sqrt2_source_bb_test.go v1
