// cf/lft_wb_test.go v1
package cf

import "testing"

func TestWB_IdentityUnaryLFT_HasCoefficients_1_0_0_1(t *testing.T) {
	u := identityUnaryLFT()

	assertBigIntString(t, "u.a", u.a, "1")
	assertBigIntString(t, "u.b", u.b, "0")
	assertBigIntString(t, "u.c", u.c, "0")
	assertBigIntString(t, "u.d", u.d, "1")
}

func TestWB_IdentityBinaryLFT_HasCoefficients_1_0_0_0_0_0_0_1(t *testing.T) {
	b := identityBinaryLFT()

	assertBigIntString(t, "b.a", b.a, "1")
	assertBigIntString(t, "b.b", b.b, "0")
	assertBigIntString(t, "b.c", b.c, "0")
	assertBigIntString(t, "b.d", b.d, "0")
	assertBigIntString(t, "b.e", b.e, "0")
	assertBigIntString(t, "b.f", b.f, "0")
	assertBigIntString(t, "b.g", b.g, "0")
	assertBigIntString(t, "b.h", b.h, "1")
}

func TestWB_BinaryDecision_HasExactlyThreeValidPrimaryActions(t *testing.T) {
	actions := []binaryDecision{
		decisionIngestLeft,
		decisionIngestRight,
		decisionEmitOutput,
	}

	seen := make(map[binaryDecision]bool, len(actions))
	for _, action := range actions {
		if seen[action] {
			t.Fatalf("duplicate action value detected: %v", action)
		}
		seen[action] = true

		if !action.isValid() {
			t.Fatalf("expected action %v to be valid", action)
		}
	}

	if got, want := len(seen), 3; got != want {
		t.Fatalf("number of distinct primary actions = %d, want %d", got, want)
	}
}

func TestWB_BinaryDecision_InvalidValueIsNotValid(t *testing.T) {
	invalid := binaryDecision(99)
	if invalid.isValid() {
		t.Fatalf("invalid decision value should not be valid")
	}
}

func assertBigIntString(t *testing.T, label string, got interface{ String() string }, want string) {
	t.Helper()

	if got == nil {
		t.Fatalf("%s is nil, want %s", label, want)
	}
	if got.String() != want {
		t.Fatalf("%s = %s, want %s", label, got.String(), want)
	}
}

func TestWB_BinaryStepState_Choose_EmitOutputWhenAvailable(t *testing.T) {
	state := binaryStepState{
		canEmitOutput:  true,
		canIngestLeft:  true,
		canIngestRight: true,
	}

	got, err := state.choose()
	if err != nil {
		t.Fatalf("choose() error: %v", err)
	}
	if got != decisionEmitOutput {
		t.Fatalf("choose() = %v, want %v", got, decisionEmitOutput)
	}
}

func TestWB_BinaryStepState_Choose_IngestLeftWhenOnlyLeftAvailable(t *testing.T) {
	state := binaryStepState{
		canEmitOutput:  false,
		canIngestLeft:  true,
		canIngestRight: false,
	}

	got, err := state.choose()
	if err != nil {
		t.Fatalf("choose() error: %v", err)
	}
	if got != decisionIngestLeft {
		t.Fatalf("choose() = %v, want %v", got, decisionIngestLeft)
	}
}

func TestWB_BinaryStepState_Choose_IngestRightWhenOnlyRightAvailable(t *testing.T) {
	state := binaryStepState{
		canEmitOutput:  false,
		canIngestLeft:  false,
		canIngestRight: true,
	}

	got, err := state.choose()
	if err != nil {
		t.Fatalf("choose() error: %v", err)
	}
	if got != decisionIngestRight {
		t.Fatalf("choose() = %v, want %v", got, decisionIngestRight)
	}
}

func TestWB_BinaryStepState_Choose_PrefersLeftOnSymmetricIngestTie(t *testing.T) {
	state := binaryStepState{
		canEmitOutput:  false,
		canIngestLeft:  true,
		canIngestRight: true,
	}

	got, err := state.choose()
	if err != nil {
		t.Fatalf("choose() error: %v", err)
	}
	if got != decisionIngestLeft {
		t.Fatalf("choose() = %v, want %v", got, decisionIngestLeft)
	}
}

func TestWB_BinaryStepState_Choose_ReturnsErrStuckWhenNoPrimaryActionAvailable(t *testing.T) {
	state := binaryStepState{
		canEmitOutput:  false,
		canIngestLeft:  false,
		canIngestRight: false,
	}

	_, err := state.choose()
	if err != ErrStuck {
		t.Fatalf("choose() error = %v, want %v", err, ErrStuck)
	}
}

// cf/lft_wb_test.go v1
