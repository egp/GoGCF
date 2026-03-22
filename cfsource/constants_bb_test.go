// cfsource/constants_bb_test.go v2
package cfsource_test

import (
	"testing"

	"github.com/egp/GoGCF/cf"
	"github.com/egp/GoGCF/cfsource"
)

func TestBB_FromSource_E_FirstSixTermsAre_2_1_2_1_1_4(t *testing.T) {
	g := cf.FromSource(cfsource.E())

	want := []string{"2", "1", "2", "1", "1", "4"}
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

func TestBB_FromSource_E_RangeIsInsideOpen_TwoToThree(t *testing.T) {
	g := cf.FromSource(cfsource.E())

	r := g.Range()

	if !r.IsInside() {
		t.Fatalf("e range should be inside")
	}
	if r.IsOutside() {
		t.Fatalf("e range should not be outside")
	}
	if r.IsExact() {
		t.Fatalf("e range should not be exact")
	}
	if r.Lo.Kind != cf.BoundFinite || r.Hi.Kind != cf.BoundFinite {
		t.Fatalf("e range should have finite bounds")
	}
	if r.Lo.Closed || r.Hi.Closed {
		t.Fatalf("e initial bounds should be open")
	}
	if got, want := r.Lo.Value.Num.String(), "2"; got != want {
		t.Fatalf("lo numerator = %s, want %s", got, want)
	}
	if got, want := r.Lo.Value.Den.String(), "1"; got != want {
		t.Fatalf("lo denominator = %s, want %s", got, want)
	}
	if got, want := r.Hi.Value.Num.String(), "3"; got != want {
		t.Fatalf("hi numerator = %s, want %s", got, want)
	}
	if got, want := r.Hi.Value.Den.String(), "1"; got != want {
		t.Fatalf("hi denominator = %s, want %s", got, want)
	}
}

func TestBB_FromSource_E_RangeAfterOneNext_IsInsideOpen_OneToTwo(t *testing.T) {
	g := cf.FromSource(cfsource.E())

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
	if got, want := value.String(), "2"; got != want {
		t.Fatalf("first term = %s, want %s", got, want)
	}

	r := rest.Range()

	if !r.IsInside() {
		t.Fatalf("e remainder range should be inside")
	}
	if r.IsOutside() {
		t.Fatalf("e remainder range should not be outside")
	}
	if r.IsExact() {
		t.Fatalf("e remainder range should not be exact")
	}
	if r.Lo.Kind != cf.BoundFinite || r.Hi.Kind != cf.BoundFinite {
		t.Fatalf("e remainder range should have finite bounds")
	}
	if r.Lo.Closed || r.Hi.Closed {
		t.Fatalf("e remainder bounds should be open")
	}
	if got, want := r.Lo.Value.Num.String(), "1"; got != want {
		t.Fatalf("lo numerator = %s, want %s", got, want)
	}
	if got, want := r.Lo.Value.Den.String(), "1"; got != want {
		t.Fatalf("lo denominator = %s, want %s", got, want)
	}
	if got, want := r.Hi.Value.Num.String(), "2"; got != want {
		t.Fatalf("hi numerator = %s, want %s", got, want)
	}
	if got, want := r.Hi.Value.Den.String(), "1"; got != want {
		t.Fatalf("hi denominator = %s, want %s", got, want)
	}
}

// EOF cfsource/constants_bb_test.go v2
