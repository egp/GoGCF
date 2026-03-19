// cf/source_wb_test.go v1
package cf

import (
	"math/big"
	"testing"
)

func TestWB_Int64Source_ConcreteStateTransitionsToEmitted(t *testing.T) {
	src := Int64(7)

	s0, ok := src.(int64Source)
	if !ok {
		t.Fatalf("Int64() concrete type = %T, want int64Source", src)
	}
	if s0.value != 7 {
		t.Fatalf("initial value = %d, want 7", s0.value)
	}
	if s0.emitted {
		t.Fatalf("initial emitted = true, want false")
	}

	term, rest, err := s0.NextPQ()
	if err != nil {
		t.Fatalf("NextPQ() error: %v", err)
	}
	if !term.IsValue() {
		t.Fatalf("first term should be a value")
	}

	s1, ok := rest.(int64Source)
	if !ok {
		t.Fatalf("remainder concrete type = %T, want int64Source", rest)
	}
	if s1.value != 7 {
		t.Fatalf("remainder value = %d, want 7", s1.value)
	}
	if !s1.emitted {
		t.Fatalf("remainder emitted = false, want true")
	}

	if s0.emitted {
		t.Fatalf("original source mutated: emitted = true, want false")
	}
}

func TestWB_FromSource_Sqrt2_TakeDoesNotMutateOriginalEvaluator(t *testing.T) {
	g := FromSource(Sqrt2())

	sb, ok := g.(sourceBackedGCF)
	if !ok {
		t.Fatalf("FromSource(Sqrt2()) concrete type = %T, want sourceBackedGCF", g)
	}

	src0, ok := sb.src.(sqrt2Source)
	if !ok {
		t.Fatalf("sourceBackedGCF src concrete type = %T, want sqrt2Source", sb.src)
	}
	if src0.index != 0 {
		t.Fatalf("initial sqrt2 index = %d, want 0", src0.index)
	}

	_, _ = g.Take(3)

	src1, ok := sb.src.(sqrt2Source)
	if !ok {
		t.Fatalf("sourceBackedGCF src concrete type after Take = %T, want sqrt2Source", sb.src)
	}
	if src1.index != 0 {
		t.Fatalf("Take mutated original evaluator source index = %d, want 0", src1.index)
	}

	term, rest, err := g.Next()
	if err != nil {
		t.Fatalf("original evaluator Next() error: %v", err)
	}
	if !term.IsValue() {
		t.Fatalf("original evaluator first term should still be a value")
	}
	value, ok := term.BigInt()
	if !ok {
		t.Fatalf("original evaluator first term should expose a BigInt")
	}
	if got, want := value.String(), "1"; got != want {
		t.Fatalf("original evaluator first term = %s, want %s", got, want)
	}

	sb1, ok := rest.(sourceBackedGCF)
	if !ok {
		t.Fatalf("remainder concrete type = %T, want sourceBackedGCF", rest)
	}
	srcAfterNext, ok := sb1.src.(sqrt2Source)
	if !ok {
		t.Fatalf("remainder src concrete type = %T, want sqrt2Source", sb1.src)
	}
	if srcAfterNext.index != 1 {
		t.Fatalf("remainder sqrt2 index = %d, want 1", srcAfterNext.index)
	}
}
func TestWB_RCFPrefixGCF_ConvergentOnConcretePrefix_SevenOverFive(t *testing.T) {
	prefix := rcfPrefixGCF{
		terms: []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(2)},
		index: 0,
	}

	r, err := prefix.Convergent()
	if err != nil {
		t.Fatalf("Convergent() error: %v", err)
	}

	if got, want := r.Num.String(), "7"; got != want {
		t.Fatalf("numerator = %s, want %s", got, want)
	}
	if got, want := r.Den.String(), "5"; got != want {
		t.Fatalf("denominator = %s, want %s", got, want)
	}
}

func TestWB_RCFPrefixGCF_ConvergentUsesRemainingTermsFromIndex(t *testing.T) {
	prefix := rcfPrefixGCF{
		terms: []*big.Int{big.NewInt(9), big.NewInt(1), big.NewInt(2), big.NewInt(2)},
		index: 1,
	}

	r, err := prefix.Convergent()
	if err != nil {
		t.Fatalf("Convergent() error: %v", err)
	}

	if got, want := r.Num.String(), "7"; got != want {
		t.Fatalf("numerator = %s, want %s", got, want)
	}
	if got, want := r.Den.String(), "5"; got != want {
		t.Fatalf("denominator = %s, want %s", got, want)
	}
}
func TestWB_Rat64Source_ConcreteStateStartsAtIndexZero(t *testing.T) {
	src := Rat64(7, 5)

	s0, ok := src.(rat64Source)
	if !ok {
		t.Fatalf("Rat64() concrete type = %T, want rat64Source", src)
	}
	if s0.index != 0 {
		t.Fatalf("initial index = %d, want 0", s0.index)
	}
	if got, want := s0.num, int64(7); got != want {
		t.Fatalf("stored num = %d, want %d", got, want)
	}
	if got, want := s0.den, int64(5); got != want {
		t.Fatalf("stored den = %d, want %d", got, want)
	}
	if len(s0.terms) != 3 {
		t.Fatalf("initial terms len = %d, want 3", len(s0.terms))
	}
	if got, want := s0.terms[0].String(), "1"; got != want {
		t.Fatalf("terms[0] = %s, want %s", got, want)
	}
	if got, want := s0.terms[1].String(), "2"; got != want {
		t.Fatalf("terms[1] = %s, want %s", got, want)
	}
	if got, want := s0.terms[2].String(), "2"; got != want {
		t.Fatalf("terms[2] = %s, want %s", got, want)
	}
}

func TestWB_Rat64Source_ConcreteRemainderAdvancesWithoutMutatingOriginal(t *testing.T) {
	src := Rat64(7, 5)

	s0, ok := src.(rat64Source)
	if !ok {
		t.Fatalf("Rat64() concrete type = %T, want rat64Source", src)
	}

	term1, rest1, err := s0.NextPQ()
	if err != nil {
		t.Fatalf("first NextPQ() error: %v", err)
	}
	if !term1.IsValue() {
		t.Fatalf("first term should be a value")
	}
	if got, want := term1.Q.String(), "1"; got != want {
		t.Fatalf("first term Q = %s, want %s", got, want)
	}

	s1, ok := rest1.(rat64Source)
	if !ok {
		t.Fatalf("remainder concrete type = %T, want rat64Source", rest1)
	}
	if s1.index != 1 {
		t.Fatalf("remainder index = %d, want 1", s1.index)
	}
	if s0.index != 0 {
		t.Fatalf("original source mutated: index = %d, want 0", s0.index)
	}
}

// cf/source_wb_test.go v1
