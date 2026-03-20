// cf/binary_next_wb_test.go v1
package cf

import (
	"math/big"
	"testing"
)

func TestWB_BinaryGCF_Next_FirstCallOnDivTwelveOverFive_LeavesStepwiseBinaryRemainder(t *testing.T) {
	g := Div(FromSource(Int64(12)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Div() concrete type = %T, want binaryGCF", g)
	}

	term, rest, err := bg.Next()
	if err != nil {
		t.Fatalf("Next() error: %v", err)
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

	nextBG, ok := rest.(binaryGCF)
	if !ok {
		t.Fatalf("remainder concrete type = %T, want binaryGCF", rest)
	}
	if nextBG.resolved != nil {
		t.Fatalf("first Next() should leave a stepwise binary remainder, not a resolved prefix")
	}

	r := nextBG.Range()
	if !r.IsExact() {
		t.Fatalf("remainder Range() should be exact")
	}
	assertBigIntString(t, "r.Lo.Value.Num", r.Lo.Value.Num, "5")
	assertBigIntString(t, "r.Lo.Value.Den", r.Lo.Value.Den, "2")
	assertBigIntString(t, "r.Hi.Value.Num", r.Hi.Value.Num, "5")
	assertBigIntString(t, "r.Hi.Value.Den", r.Hi.Value.Den, "2")
}
func TestWB_BinaryGCF_Next_SecondCallOnDivTwelveOverFive_LeavesStepwiseBinaryRemainder(t *testing.T) {
	g := Div(FromSource(Int64(12)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Div() concrete type = %T, want binaryGCF", g)
	}

	_, rest1, err := bg.Next()
	if err != nil {
		t.Fatalf("first Next() error: %v", err)
	}

	term2, rest2, err := rest1.Next()
	if err != nil {
		t.Fatalf("second Next() error: %v", err)
	}
	if !term2.IsValue() {
		t.Fatalf("second term should be a value")
	}
	value2, ok := term2.BigInt()
	if !ok {
		t.Fatalf("second term should expose a BigInt")
	}
	if got, want := value2.String(), "2"; got != want {
		t.Fatalf("second term = %s, want %s", got, want)
	}

	nextBG, ok := rest2.(binaryGCF)
	if !ok {
		t.Fatalf("remainder concrete type = %T, want binaryGCF", rest2)
	}
	if nextBG.resolved != nil {
		t.Fatalf("second Next() should still leave a stepwise binary remainder")
	}

	r := nextBG.Range()
	if !r.IsExact() {
		t.Fatalf("second remainder Range() should be exact")
	}
	assertBigIntString(t, "r.Lo.Value.Num", r.Lo.Value.Num, "2")
	assertBigIntString(t, "r.Lo.Value.Den", r.Lo.Value.Den, "1")
	assertBigIntString(t, "r.Hi.Value.Num", r.Hi.Value.Num, "2")
	assertBigIntString(t, "r.Hi.Value.Den", r.Hi.Value.Den, "1")
}

func TestWB_BinaryGCF_Next_OnResolvedNodeWrapsResolvedRemainderBackIntoBinaryGCF(t *testing.T) {
	prefix, err := prefixGCFfromRational(Rational{
		Num: big.NewInt(12),
		Den: big.NewInt(5),
	})
	if err != nil {
		t.Fatalf("prefixGCFfromRational() error: %v", err)
	}

	bg := binaryGCF{resolved: prefix}

	term, rest, err := bg.Next()
	if err != nil {
		t.Fatalf("Next() error: %v", err)
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

	nextBG, ok := rest.(binaryGCF)
	if !ok {
		t.Fatalf("remainder concrete type = %T, want binaryGCF", rest)
	}
	if nextBG.resolved == nil {
		t.Fatalf("resolved remainder should stay wrapped in binaryGCF")
	}
}

func TestWB_BinaryGCF_Next_DoesNotMutateOriginalNodeWhenStepping(t *testing.T) {
	g := Div(FromSource(Int64(12)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Div() concrete type = %T, want binaryGCF", g)
	}

	_, _, err := bg.Next()
	if err != nil {
		t.Fatalf("Next() error: %v", err)
	}

	if bg.resolved != nil {
		t.Fatalf("original node should remain immutable and unresolved")
	}

	leftTerm, _, err := bg.left.Next()
	if err != nil {
		t.Fatalf("original left.Next() error: %v", err)
	}
	if !leftTerm.IsValue() {
		t.Fatalf("original left child should remain unconsumed")
	}

	rightTerm, _, err := bg.right.Next()
	if err != nil {
		t.Fatalf("original right.Next() error: %v", err)
	}
	if !rightTerm.IsValue() {
		t.Fatalf("original right child should remain unconsumed")
	}
}

// cf/binary_next_wb_test.go v1
