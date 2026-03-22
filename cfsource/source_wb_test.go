// cfsource/source_wb_test.go v1
package cfsource

import "testing"

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
