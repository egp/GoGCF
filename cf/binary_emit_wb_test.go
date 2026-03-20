// cf/binary_emit_wb_test.go v1
package cf

import "testing"

func TestWB_BinaryGCF_EmitStep_OnDivTwelveByFive_EmitsTwoAndLeavesExactFiveOverTwoRange(t *testing.T) {
	g := Div(FromSource(Int64(12)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Div() concrete type = %T, want binaryGCF", g)
	}

	term, next, err := bg.emitStep()
	if err != nil {
		t.Fatalf("emitStep() error: %v", err)
	}
	if !term.IsValue() {
		t.Fatalf("emitStep() should emit a value term")
	}

	value, ok := term.BigInt()
	if !ok {
		t.Fatalf("emitStep() value term should expose a BigInt")
	}
	if got, want := value.String(), "2"; got != want {
		t.Fatalf("emitStep() term = %s, want %s", got, want)
	}

	r := next.Range()
	if !r.IsExact() {
		t.Fatalf("remainder Range() should be exact after emit")
	}
	assertBigIntString(t, "r.Lo.Value.Num", r.Lo.Value.Num, "5")
	assertBigIntString(t, "r.Lo.Value.Den", r.Lo.Value.Den, "2")
	assertBigIntString(t, "r.Hi.Value.Num", r.Hi.Value.Num, "5")
	assertBigIntString(t, "r.Hi.Value.Den", r.Hi.Value.Den, "2")
}
func TestWB_BinaryGCF_Step_ChoosesEmitWhenExactEmitIsAvailable(t *testing.T) {
	g := Div(FromSource(Int64(12)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Div() concrete type = %T, want binaryGCF", g)
	}

	action, next, err := bg.step()
	if err != nil {
		t.Fatalf("step() error: %v", err)
	}
	if action != decisionEmitOutput {
		t.Fatalf("step() action = %v, want %v", action, decisionEmitOutput)
	}

	r := next.Range()
	if !r.IsExact() {
		t.Fatalf("remainder Range() should be exact after emit")
	}
	assertBigIntString(t, "r.Lo.Value.Num", r.Lo.Value.Num, "5")
	assertBigIntString(t, "r.Lo.Value.Den", r.Lo.Value.Den, "2")
	assertBigIntString(t, "r.Hi.Value.Num", r.Hi.Value.Num, "5")
	assertBigIntString(t, "r.Hi.Value.Den", r.Hi.Value.Den, "2")
}

func TestWB_BinaryGCF_Step_ChoosesEmitFromPropagatedRangeWhenIntegerPartIsCertified(t *testing.T) {
	bg := binaryGCF{
		op:    divisionBinaryLFT(),
		left:  FromSource(Rat64(5, 1)),
		right: FromSource(Rat64(2, 1)),
	}

	action, _, err := bg.step()
	if err != nil {
		t.Fatalf("step() error: %v", err)
	}
	if action != decisionEmitOutput {
		t.Fatalf("step() action = %v, want %v", action, decisionEmitOutput)
	}
}
func TestWB_BinaryGCF_Step_ChoosesEmitFromPropagatedRange_ForThreePlusSqrt2(t *testing.T) {
	g := Add(FromSource(Int64(3)), FromSource(Sqrt2()))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Add() concrete type = %T, want binaryGCF", g)
	}

	action, _, err := bg.step()
	if err != nil {
		t.Fatalf("step() error: %v", err)
	}
	if action != decisionEmitOutput {
		t.Fatalf("step() action = %v, want %v", action, decisionEmitOutput)
	}
}

// cf/binary_emit_wb_test.go v1
