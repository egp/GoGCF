// cf/diag_wb_test.go v1
package cf

import (
	"testing"
)

func TestWB_DiagGCF_Range_OnIdentityOverFiveOverTwo_IsExactFiveOverTwo(t *testing.T) {
	g := diagGCF{
		op:    identityDiagLFT(),
		child: FromSource(Rat64(5, 2)),
	}

	r := g.Range()
	if !r.IsExact() {
		t.Fatalf("Range() should be exact")
	}

	assertBigIntString(t, "r.Lo.Value.Num", r.Lo.Value.Num, "5")
	assertBigIntString(t, "r.Lo.Value.Den", r.Lo.Value.Den, "2")
	assertBigIntString(t, "r.Hi.Value.Num", r.Hi.Value.Num, "5")
	assertBigIntString(t, "r.Hi.Value.Den", r.Hi.Value.Den, "2")
}

func TestWB_DiagGCF_Next_OnIdentityOverTwelveOverFive_EmitsTwoThenLeavesExactFiveOverTwo(t *testing.T) {
	g := diagGCF{
		op:    identityDiagLFT(),
		child: FromSource(Rat64(12, 5)),
	}

	term, rest, err := g.Next()
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

	nextDG, ok := rest.(diagGCF)
	if !ok {
		t.Fatalf("remainder concrete type = %T, want diagGCF", rest)
	}
	if nextDG.resolved != nil {
		t.Fatalf("remainder should stay stepwise, not resolved")
	}

	r := nextDG.Range()
	if !r.IsExact() {
		t.Fatalf("remainder Range() should be exact")
	}

	assertBigIntString(t, "r.Lo.Value.Num", r.Lo.Value.Num, "5")
	assertBigIntString(t, "r.Lo.Value.Den", r.Lo.Value.Den, "2")
	assertBigIntString(t, "r.Hi.Value.Num", r.Hi.Value.Num, "5")
	assertBigIntString(t, "r.Hi.Value.Den", r.Hi.Value.Den, "2")
}

func TestWB_DiagGCF_Next_DoesNotMutateOriginalNode(t *testing.T) {
	g := diagGCF{
		op:    identityDiagLFT(),
		child: FromSource(Rat64(12, 5)),
	}

	_, _, err := g.Next()
	if err != nil {
		t.Fatalf("Next() error: %v", err)
	}

	if g.resolved != nil {
		t.Fatalf("original diag node should remain unresolved")
	}

	r := g.child.Range()
	if !r.IsExact() {
		t.Fatalf("original child should remain unchanged")
	}
	assertBigIntString(t, "r.Lo.Value.Num", r.Lo.Value.Num, "12")
	assertBigIntString(t, "r.Lo.Value.Den", r.Lo.Value.Den, "5")
}

// cf/diag_wb_test.go v1
