// cf/lft_wb_test.go v1
package cf

import (
	"math/big"
	"testing"
)

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

func TestWB_BinaryLFT_CollapseLeftEOF_Selects_AXplusB_Over_EXplusF(t *testing.T) {
	b := binaryLFT{
		a: big.NewInt(2),
		b: big.NewInt(3),
		c: big.NewInt(5),
		d: big.NewInt(7),
		e: big.NewInt(11),
		f: big.NewInt(13),
		g: big.NewInt(17),
		h: big.NewInt(19),
	}

	u := b.collapseLeftEOF()

	assertBigIntString(t, "u.a", u.a, "2")
	assertBigIntString(t, "u.b", u.b, "3")
	assertBigIntString(t, "u.c", u.c, "11")
	assertBigIntString(t, "u.d", u.d, "13")
}

func TestWB_BinaryLFT_CollapseRightEOF_Selects_AXplusC_Over_EXplusG(t *testing.T) {
	b := binaryLFT{
		a: big.NewInt(2),
		b: big.NewInt(3),
		c: big.NewInt(5),
		d: big.NewInt(7),
		e: big.NewInt(11),
		f: big.NewInt(13),
		g: big.NewInt(17),
		h: big.NewInt(19),
	}

	u := b.collapseRightEOF()

	assertBigIntString(t, "u.a", u.a, "2")
	assertBigIntString(t, "u.b", u.b, "5")
	assertBigIntString(t, "u.c", u.c, "11")
	assertBigIntString(t, "u.d", u.d, "17")
}

func TestWB_BinaryLFT_CollapseEOF_DoesNotMutateOriginal(t *testing.T) {
	b := binaryLFT{
		a: big.NewInt(2),
		b: big.NewInt(3),
		c: big.NewInt(5),
		d: big.NewInt(7),
		e: big.NewInt(11),
		f: big.NewInt(13),
		g: big.NewInt(17),
		h: big.NewInt(19),
	}

	_ = b.collapseLeftEOF()
	_ = b.collapseRightEOF()

	assertBigIntString(t, "b.a", b.a, "2")
	assertBigIntString(t, "b.b", b.b, "3")
	assertBigIntString(t, "b.c", b.c, "5")
	assertBigIntString(t, "b.d", b.d, "7")
	assertBigIntString(t, "b.e", b.e, "11")
	assertBigIntString(t, "b.f", b.f, "13")
	assertBigIntString(t, "b.g", b.g, "17")
	assertBigIntString(t, "b.h", b.h, "19")
}

func TestWB_UnaryLFT_CollapseEOFToRational_SevenOverFive(t *testing.T) {
	u := unaryLFT{
		a: big.NewInt(7),
		b: big.NewInt(1),
		c: big.NewInt(5),
		d: big.NewInt(0),
	}

	r, err := u.collapseEOFToRational()
	if err != nil {
		t.Fatalf("collapseEOFToRational() error: %v", err)
	}

	assertBigIntString(t, "r.Num", r.Num, "7")
	assertBigIntString(t, "r.Den", r.Den, "5")
}

func TestWB_UnaryLFT_CollapseEOFToRational_EightOverOne(t *testing.T) {
	u := unaryLFT{
		a: big.NewInt(8),
		b: big.NewInt(5),
		c: big.NewInt(1),
		d: big.NewInt(0),
	}

	r, err := u.collapseEOFToRational()
	if err != nil {
		t.Fatalf("collapseEOFToRational() error: %v", err)
	}

	assertBigIntString(t, "r.Num", r.Num, "8")
	assertBigIntString(t, "r.Den", r.Den, "1")
}

func TestWB_AddPath_AfterBothIngests_CollapseRightEOFThenUnaryEOF_GivesEight(t *testing.T) {
	b := additionBinaryLFT()
	b = b.ingestLeft(big.NewInt(1), big.NewInt(3))
	b = b.ingestRight(big.NewInt(1), big.NewInt(5))

	u := b.collapseRightEOF()
	r, err := u.collapseEOFToRational()
	if err != nil {
		t.Fatalf("collapse path error: %v", err)
	}

	assertBigIntString(t, "r.Num", r.Num, "8")
	assertBigIntString(t, "r.Den", r.Den, "1")
}

// append to cf/lft_wb_test.go v2
func TestWB_BinaryLFT_ExactQuotient_AfterAddThreeAndFiveIngests_IsEightOverOne(t *testing.T) {
	b := additionBinaryLFT()
	b = b.ingestLeft(big.NewInt(1), big.NewInt(3))
	b = b.ingestRight(big.NewInt(1), big.NewInt(5))

	r, ok := b.exactQuotient()
	if !ok {
		t.Fatalf("exactQuotient() should succeed after both finite ingests")
	}

	assertBigIntString(t, "r.Num", r.Num, "8")
	assertBigIntString(t, "r.Den", r.Den, "1")
}

func TestWB_BinaryLFT_ExactQuotient_AfterDivTwelveByFiveIngests_IsTwelveOverFive(t *testing.T) {
	b := divisionBinaryLFT()
	b = b.ingestLeft(big.NewInt(1), big.NewInt(12))
	b = b.ingestRight(big.NewInt(1), big.NewInt(5))

	r, ok := b.exactQuotient()
	if !ok {
		t.Fatalf("exactQuotient() should succeed after both finite ingests")
	}

	assertBigIntString(t, "r.Num", r.Num, "12")
	assertBigIntString(t, "r.Den", r.Den, "5")
}

func TestWB_BinaryLFT_Emit_TwelveOverFiveFirstDigitTwo_LeavesFiveOverTwo(t *testing.T) {
	b := divisionBinaryLFT()
	b = b.ingestLeft(big.NewInt(1), big.NewInt(12))
	b = b.ingestRight(big.NewInt(1), big.NewInt(5))

	next := b.emit(big.NewInt(2))

	r, ok := next.exactQuotient()
	if !ok {
		t.Fatalf("exactQuotient() should succeed after emitting from exact finite state")
	}

	assertBigIntString(t, "r.Num", r.Num, "5")
	assertBigIntString(t, "r.Den", r.Den, "2")
}

// append to cf/lft_wb_test.go v2

// cf/lft_wb_test.go v1
