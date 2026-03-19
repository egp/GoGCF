// cf/sqrt2_source_wb_test.go v1
package cf

import "testing"

func TestWB_Sqrt2_ConcreteStateProgressesByIndex(t *testing.T) {
	src := Sqrt2()

	s0, ok := src.(sqrt2Source)
	if !ok {
		t.Fatalf("Sqrt2() concrete type = %T, want sqrt2Source", src)
	}
	if s0.index != 0 {
		t.Fatalf("initial sqrt2Source index = %d, want 0", s0.index)
	}

	term1, rest1, err := s0.NextPQ()
	if err != nil {
		t.Fatalf("first NextPQ() error: %v", err)
	}
	if !term1.IsValue() {
		t.Fatalf("first term should be a value")
	}
	if got, want := term1.P.String(), "1"; got != want {
		t.Fatalf("first term P = %s, want %s", got, want)
	}
	if got, want := term1.Q.String(), "1"; got != want {
		t.Fatalf("first term Q = %s, want %s", got, want)
	}

	s1, ok := rest1.(sqrt2Source)
	if !ok {
		t.Fatalf("remainder concrete type = %T, want sqrt2Source", rest1)
	}
	if s1.index != 1 {
		t.Fatalf("remainder index = %d, want 1", s1.index)
	}

	term2, rest2, err := s1.NextPQ()
	if err != nil {
		t.Fatalf("second NextPQ() error: %v", err)
	}
	if !term2.IsValue() {
		t.Fatalf("second term should be a value")
	}
	if got, want := term2.P.String(), "1"; got != want {
		t.Fatalf("second term P = %s, want %s", got, want)
	}
	if got, want := term2.Q.String(), "2"; got != want {
		t.Fatalf("second term Q = %s, want %s", got, want)
	}

	s2, ok := rest2.(sqrt2Source)
	if !ok {
		t.Fatalf("second remainder concrete type = %T, want sqrt2Source", rest2)
	}
	if s2.index != 2 {
		t.Fatalf("second remainder index = %d, want 2", s2.index)
	}

	if s0.index != 0 {
		t.Fatalf("original source mutated: index = %d, want 0", s0.index)
	}
}

// cf/sqrt2_source_wb_test.go v1
