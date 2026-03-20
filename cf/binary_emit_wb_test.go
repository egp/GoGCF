// cf/binary_emit_wb_test.go v1
package cf

import "testing"

func TestWB_BinaryGCF_EmitStep_AfterBothIngestsOnDivTwelveByFive_EmitsTwoAndLeavesFiveOverTwo(t *testing.T) {
	g := Div(FromSource(Int64(12)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Div() concrete type = %T, want binaryGCF", g)
	}

	_, bg1, err := bg.step()
	if err != nil {
		t.Fatalf("first step() error: %v", err)
	}

	_, bg2, err := bg1.step()
	if err != nil {
		t.Fatalf("second step() error: %v", err)
	}

	term, next, err := bg2.emitStep()
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

	r, ok := next.op.exactQuotient()
	if !ok {
		t.Fatalf("remainder exactQuotient() should succeed after emit")
	}
	assertBigIntString(t, "r.Num", r.Num, "5")
	assertBigIntString(t, "r.Den", r.Den, "2")
}

func TestWB_BinaryGCF_Step_ChoosesEmitWhenExactEmitIsAvailable(t *testing.T) {
	g := Div(FromSource(Int64(12)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Div() concrete type = %T, want binaryGCF", g)
	}

	_, bg1, err := bg.step()
	if err != nil {
		t.Fatalf("first step() error: %v", err)
	}

	_, bg2, err := bg1.step()
	if err != nil {
		t.Fatalf("second step() error: %v", err)
	}

	action, next, err := bg2.step()
	if err != nil {
		t.Fatalf("third step() error: %v", err)
	}
	if action != decisionEmitOutput {
		t.Fatalf("third step() action = %v, want %v", action, decisionEmitOutput)
	}

	r, ok := next.op.exactQuotient()
	if !ok {
		t.Fatalf("remainder exactQuotient() should succeed after emit")
	}
	assertBigIntString(t, "r.Num", r.Num, "5")
	assertBigIntString(t, "r.Den", r.Den, "2")
}

// cf/binary_emit_wb_test.go v1
