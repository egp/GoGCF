// cf/sqrt_bb_test.go v1
package cf_test

import (
	"testing"

	"github.com/egp/GoGCF/cf"
)

func TestBB_Sqrt_Int64_Zero_EmitsZeroThenEOF(t *testing.T) {
	g := cf.Sqrt(cf.FromSource(cf.Int64(0)))

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
	if got, want := value1.String(), "0"; got != want {
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

func TestBB_Sqrt_Int64_Four_EmitsTwoThenEOF(t *testing.T) {
	g := cf.Sqrt(cf.FromSource(cf.Int64(4)))

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
	if got, want := value1.String(), "2"; got != want {
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

func TestBB_Sqrt_Rat64_NineOverFour_Emits_1_2_ThenEOF(t *testing.T) {
	g := cf.Sqrt(cf.FromSource(cf.Rat64(9, 4)))

	want := []string{"1", "2"}
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

	term3, _, err := cur.Next()
	if err != nil {
		t.Fatalf("final Next() error: %v", err)
	}
	if !term3.IsEOF() {
		t.Fatalf("final term should be EOF")
	}
}

func TestBB_Sqrt_Int64_NegativeOne_NextReturnsErrUndefined(t *testing.T) {
	g := cf.Sqrt(cf.FromSource(cf.Int64(-1)))

	_, _, err := g.Next()
	if err != cf.ErrUndefined {
		t.Fatalf("Next() error = %v, want %v", err, cf.ErrUndefined)
	}
}
func TestBB_Sqrt_Int64_Two_FirstTermIsOne(t *testing.T) {
	g := cf.Sqrt(cf.FromSource(cf.Int64(2)))

	term1, _, err := g.Next()
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
	if got, want := value1.String(), "1"; got != want {
		t.Fatalf("first term = %s, want %s", got, want)
	}
}
func TestBB_Sqrt_Int64_Two_FirstTwoTermsAre_1_2(t *testing.T) {
	g := cf.Sqrt(cf.FromSource(cf.Int64(2)))

	want := []string{"1", "2"}
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
}
func TestBB_Sqrt_Int64_Two_FirstThreeTermsAre_1_2_2(t *testing.T) {
	g := cf.Sqrt(cf.FromSource(cf.Int64(2)))

	want := []string{"1", "2", "2"}
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
}
func TestBB_Sqrt_Int64_Two_FirstFourTermsAre_1_2_2_2(t *testing.T) {
	g := cf.Sqrt(cf.FromSource(cf.Int64(2)))

	want := []string{"1", "2", "2", "2"}
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
}
func TestBB_Sqrt_Int64_Two_FirstFiveTermsAre_1_2_2_2_2(t *testing.T) {
	g := cf.Sqrt(cf.FromSource(cf.Int64(2)))

	want := []string{"1", "2", "2", "2", "2"}
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
}

func TestBB_Sqrt_Int64_Two_FirstTenTermsAre_1_FollowedByNineTwos(t *testing.T) {
	g := cf.Sqrt(cf.FromSource(cf.Int64(2)))

	want := []string{"1", "2", "2", "2", "2", "2", "2", "2", "2", "2"}
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
}

// cf/sqrt_bb_test.go v1
