// cf/cfsource_smoke_bb_test.go v1
package cf_test

import (
	"testing"

	"github.com/egp/GoGCF/cf"
	"github.com/egp/GoGCF/cfsource"
)

func TestBB_CFSource_Int64_Three_EmitsThreeThenEOF(t *testing.T) {
	g := cf.FromSource(cfsource.Int64(3))

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
	if got, want := value1.String(), "3"; got != want {
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

func TestBB_CFSource_Rat64_ThreeHalves_Emits_1_2_ThenEOF(t *testing.T) {
	g := cf.FromSource(cfsource.Rat64(3, 2))

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

func TestBB_CFSource_Sqrt2_FirstTwoTermsAre_1_2(t *testing.T) {
	g := cf.FromSource(cfsource.Sqrt2())

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

// cf/cfsource_smoke_bb_test.go v1
