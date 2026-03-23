// cf/binary_bb_test.go v1
package cf_test

import (
	"testing"

	"github.com/egp/GoGCF/cf"
)

func TestBB_Add_Int64_ThreePlusFive_EmitsEightThenEOF(t *testing.T) {
	g := cf.Add(cf.FromSource(cf.Int64(3)), cf.FromSource(cf.Int64(5)))

	term1, rest, err := g.Next()
	if err != nil {
		t.Fatalf("first Next() error: %v", err)
	}
	if !term1.IsValue() {
		t.Fatalf("first term should be a value")
	}
	value1, ok := term1.BigInt()
	if !ok {
		t.Fatalf("first term should expose a BigInt")
	}
	if got, want := value1.String(), "8"; got != want {
		t.Fatalf("first term = %s, want %s", got, want)
	}

	term2, _, err := rest.Next()
	if err != nil {
		t.Fatalf("second Next() error: %v", err)
	}
	if !term2.IsEOF() {
		t.Fatalf("second term should be EOF")
	}
}

func TestBB_Div_Int64_TwelveOverFive_Emits_2_2_2_ThenEOF(t *testing.T) {
	g := cf.Div(cf.FromSource(cf.Int64(12)), cf.FromSource(cf.Int64(5)))

	want := []string{"2", "2", "2"}
	cur := g
	for i, w := range want {
		term, rest, err := cur.Next()
		if err != nil {
			t.Fatalf("Next() #%d error: %v", i+1, err)
		}
		if !term.IsValue() {
			t.Fatalf("term #%d should be a value", i+1)
		}
		value, ok := term.BigInt()
		if !ok {
			t.Fatalf("term #%d should expose a BigInt", i+1)
		}
		if got := value.String(); got != w {
			t.Fatalf("term #%d = %s, want %s", i+1, got, w)
		}
		cur = rest
	}

	term4, _, err := cur.Next()
	if err != nil {
		t.Fatalf("final Next() error: %v", err)
	}
	if !term4.IsEOF() {
		t.Fatalf("final term should be EOF")
	}
}

func TestBB_Sub_Int64_EightMinusThree_EmitsFiveThenEOF(t *testing.T) {
	g := cf.Sub(cf.FromSource(cf.Int64(8)), cf.FromSource(cf.Int64(3)))

	term1, rest, err := g.Next()
	if err != nil {
		t.Fatalf("first Next() error: %v", err)
	}
	if !term1.IsValue() {
		t.Fatalf("first term should be a value")
	}
	value1, ok := term1.BigInt()
	if !ok {
		t.Fatalf("first term should expose a BigInt")
	}
	if got, want := value1.String(), "5"; got != want {
		t.Fatalf("first term = %s, want %s", got, want)
	}

	term2, _, err := rest.Next()
	if err != nil {
		t.Fatalf("second Next() error: %v", err)
	}
	if !term2.IsEOF() {
		t.Fatalf("second term should be EOF")
	}
}

func TestBB_Sub_Int64_ThreeMinusEight_EmitsMinusFiveThenEOF(t *testing.T) {
	g := cf.Sub(cf.FromSource(cf.Int64(3)), cf.FromSource(cf.Int64(8)))

	term1, rest, err := g.Next()
	if err != nil {
		t.Fatalf("first Next() error: %v", err)
	}
	if !term1.IsValue() {
		t.Fatalf("first term should be a value")
	}
	value1, ok := term1.BigInt()
	if !ok {
		t.Fatalf("first term should expose a BigInt")
	}
	if got, want := value1.String(), "-5"; got != want {
		t.Fatalf("first term = %s, want %s", got, want)
	}

	term2, _, err := rest.Next()
	if err != nil {
		t.Fatalf("second Next() error: %v", err)
	}
	if !term2.IsEOF() {
		t.Fatalf("second term should be EOF")
	}
}

func TestBB_Div_Int64_FiveOverTwelve_Emits_0_2_2_2_ThenEOF(t *testing.T) {
	g := cf.Div(cf.FromSource(cf.Int64(5)), cf.FromSource(cf.Int64(12)))

	want := []string{"0", "2", "2", "2"}
	cur := g
	for i, w := range want {
		term, rest, err := cur.Next()
		if err != nil {
			t.Fatalf("Next() #%d error: %v", i+1, err)
		}
		if !term.IsValue() {
			t.Fatalf("term #%d should be a value", i+1)
		}
		value, ok := term.BigInt()
		if !ok {
			t.Fatalf("term #%d should expose a BigInt", i+1)
		}
		if got := value.String(); got != w {
			t.Fatalf("term #%d = %s, want %s", i+1, got, w)
		}
		cur = rest
	}

	term5, _, err := cur.Next()
	if err != nil {
		t.Fatalf("final Next() error: %v", err)
	}
	if !term5.IsEOF() {
		t.Fatalf("final term should be EOF")
	}
}

func TestBB_Add_Int64_TakeOneReturnsEightThenEOF(t *testing.T) {
	g := cf.Add(cf.FromSource(cf.Int64(3)), cf.FromSource(cf.Int64(5)))

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
	if got, want := value1.String(), "8"; got != want {
		t.Fatalf("prefix first term = %s, want %s", got, want)
	}

	term2, _, err := rest.Next()
	if err != nil {
		t.Fatalf("prefix second Next() error: %v", err)
	}
	if !term2.IsEOF() {
		t.Fatalf("prefix second term should be EOF")
	}
}

func TestBB_Div_Int64_TwelveOverFive_ConvergentIsTwelveOverFive(t *testing.T) {
	g := cf.Div(cf.FromSource(cf.Int64(12)), cf.FromSource(cf.Int64(5)))

	r, err := g.Convergent()
	if err != nil {
		t.Fatalf("Convergent() error: %v", err)
	}

	if got, want := r.Num.String(), "12"; got != want {
		t.Fatalf("numerator = %s, want %s", got, want)
	}
	if got, want := r.Den.String(), "5"; got != want {
		t.Fatalf("denominator = %s, want %s", got, want)
	}
}

// func TestBB_Add_Int64_ThreePlusSqrt2_FirstTwoTermsAre_4_2(t *testing.T) {
// 	g := cf.Add(cf.FromSource(cf.Int64(3)), cf.FromSource(cf.Sqrt2()))

// 	term1, rest1, err := g.Next()
// 	if err != nil {
// 		t.Fatalf("first Next() error: %v", err)
// 	}
// 	if !term1.IsValue() {
// 		t.Fatalf("first term should be a value")
// 	}
// 	value1, ok := term1.BigInt()
// 	if !ok {
// 		t.Fatalf("first term should expose a BigInt")
// 	}
// 	if got, want := value1.String(), "4"; got != want {
// 		t.Fatalf("first term = %s, want %s", got, want)
// 	}

// 	term2, _, err := rest1.Next()
// 	if err != nil {
// 		t.Fatalf("second Next() error: %v", err)
// 	}
// 	if !term2.IsValue() {
// 		t.Fatalf("second term should be a value")
// 	}
// 	value2, ok := term2.BigInt()
// 	if !ok {
// 		t.Fatalf("second term should expose a BigInt")
// 	}
// 	if got, want := value2.String(), "2"; got != want {
// 		t.Fatalf("second term = %s, want %s", got, want)
// 	}
// }

// cf/binary_bb_test.go v1
