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

// cf/lft_wb_test.go v1
