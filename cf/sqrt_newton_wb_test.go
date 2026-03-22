// cf/sqrt_newton_wb_test.go v1
package cf

import (
	"math/big"
	"testing"
)

func TestWB_NewtonUpperSqrtBound_TwoFromTwo_IsThreeOverTwo(t *testing.T) {
	next, err := newtonUpperSqrtBound(
		Rational{Num: big.NewInt(2), Den: big.NewInt(1)},
		Rational{Num: big.NewInt(2), Den: big.NewInt(1)},
	)
	if err != nil {
		t.Fatalf("newtonUpperSqrtBound() error: %v", err)
	}

	assertBigIntString(t, "next.Num", next.Num, "3")
	assertBigIntString(t, "next.Den", next.Den, "2")
}

func TestWB_SqrtEnclosureFromUpperBound_TwoFromThreeOverTwo_IsOpenFourThirdsToThreeHalves(t *testing.T) {
	zr, err := sqrtEnclosureFromUpperBound(
		Rational{Num: big.NewInt(2), Den: big.NewInt(1)},
		Rational{Num: big.NewInt(3), Den: big.NewInt(2)},
	)
	if err != nil {
		t.Fatalf("sqrtEnclosureFromUpperBound() error: %v", err)
	}
	if zr.IsExact() {
		t.Fatalf("sqrtEnclosureFromUpperBound() should not be exact")
	}
	if zr.Outside {
		t.Fatalf("sqrtEnclosureFromUpperBound() should return an inside range")
	}

	assertBigIntString(t, "zr.Lo.Value.Num", zr.Lo.Value.Num, "4")
	assertBigIntString(t, "zr.Lo.Value.Den", zr.Lo.Value.Den, "3")
	assertBigIntString(t, "zr.Hi.Value.Num", zr.Hi.Value.Num, "3")
	assertBigIntString(t, "zr.Hi.Value.Den", zr.Hi.Value.Den, "2")
	if zr.Lo.Closed || zr.Hi.Closed {
		t.Fatalf("irrational sqrt enclosure should stay open")
	}
}

func TestWB_SqrtEnclosureFromUpperBound_TwoFromSeventeenTwelfths_IsOpenTwentyFourSeventeenthsToSeventeenTwelfths(t *testing.T) {
	zr, err := sqrtEnclosureFromUpperBound(
		Rational{Num: big.NewInt(2), Den: big.NewInt(1)},
		Rational{Num: big.NewInt(17), Den: big.NewInt(12)},
	)
	if err != nil {
		t.Fatalf("sqrtEnclosureFromUpperBound() error: %v", err)
	}
	if zr.IsExact() {
		t.Fatalf("sqrtEnclosureFromUpperBound() should not be exact")
	}
	if zr.Outside {
		t.Fatalf("sqrtEnclosureFromUpperBound() should return an inside range")
	}

	assertBigIntString(t, "zr.Lo.Value.Num", zr.Lo.Value.Num, "24")
	assertBigIntString(t, "zr.Lo.Value.Den", zr.Lo.Value.Den, "17")
	assertBigIntString(t, "zr.Hi.Value.Num", zr.Hi.Value.Num, "17")
	assertBigIntString(t, "zr.Hi.Value.Den", zr.Hi.Value.Den, "12")
	if zr.Lo.Closed || zr.Hi.Closed {
		t.Fatalf("irrational sqrt enclosure should stay open")
	}
}

func TestWB_ExactPositiveSqrtNewtonFeedback_ExactTwoWithIdentity_CertifiesOneUsingUpperTwo(t *testing.T) {
	got, ok, err := exactPositiveSqrtNewtonFeedback(
		Rational{Num: big.NewInt(2), Den: big.NewInt(1)},
		identityUnaryLFT(),
		Rational{},
	)
	if err != nil {
		t.Fatalf("exactPositiveSqrtNewtonFeedback() error: %v", err)
	}
	if !ok {
		t.Fatalf("exactPositiveSqrtNewtonFeedback() should succeed")
	}
	if got.Candidate == nil {
		t.Fatalf("candidate should be non-nil")
	}

	assertBigIntString(t, "got.Candidate", got.Candidate, "1")
	assertBigIntString(t, "got.UpperBound.Num", got.UpperBound.Num, "2")
	assertBigIntString(t, "got.UpperBound.Den", got.UpperBound.Den, "1")

	if got.ImageRange.IsExact() {
		t.Fatalf("image range should not be exact")
	}
	assertBigIntString(t, "got.ImageRange.Lo.Value.Num", got.ImageRange.Lo.Value.Num, "1")
	assertBigIntString(t, "got.ImageRange.Lo.Value.Den", got.ImageRange.Lo.Value.Den, "1")
	assertBigIntString(t, "got.ImageRange.Hi.Value.Num", got.ImageRange.Hi.Value.Num, "2")
	assertBigIntString(t, "got.ImageRange.Hi.Value.Den", got.ImageRange.Hi.Value.Den, "1")
	if got.ImageRange.Lo.Closed || got.ImageRange.Hi.Closed {
		t.Fatalf("sqrt(2) enclosure should stay open")
	}
}

func TestWB_ExactPositiveSqrtNewtonFeedback_ExactTwoAfterEmitOne_CertifiesTwoUsingUpperThreeOverTwo(t *testing.T) {
	got, ok, err := exactPositiveSqrtNewtonFeedback(
		Rational{Num: big.NewInt(2), Den: big.NewInt(1)},
		identityUnaryLFT().emit(big.NewInt(1)),
		Rational{},
	)
	if err != nil {
		t.Fatalf("exactPositiveSqrtNewtonFeedback() error: %v", err)
	}
	if !ok {
		t.Fatalf("exactPositiveSqrtNewtonFeedback() should succeed")
	}
	if got.Candidate == nil {
		t.Fatalf("candidate should be non-nil")
	}

	assertBigIntString(t, "got.Candidate", got.Candidate, "2")
	assertBigIntString(t, "got.UpperBound.Num", got.UpperBound.Num, "3")
	assertBigIntString(t, "got.UpperBound.Den", got.UpperBound.Den, "2")

	if got.ImageRange.IsExact() {
		t.Fatalf("image range should not be exact")
	}
	assertBigIntString(t, "got.ImageRange.Lo.Value.Num", got.ImageRange.Lo.Value.Num, "2")
	assertBigIntString(t, "got.ImageRange.Lo.Value.Den", got.ImageRange.Lo.Value.Den, "1")
	assertBigIntString(t, "got.ImageRange.Hi.Value.Num", got.ImageRange.Hi.Value.Num, "3")
	assertBigIntString(t, "got.ImageRange.Hi.Value.Den", got.ImageRange.Hi.Value.Den, "1")
	if got.ImageRange.Lo.Closed || got.ImageRange.Hi.Closed {
		t.Fatalf("post-emit image range should stay open")
	}
}

func TestWB_ExactPositiveSqrtNewtonFeedback_ExactTwoAfterEmitOneTwo_CertifiesTwoUsingUpperSeventeenTwelfths(t *testing.T) {
	got, ok, err := exactPositiveSqrtNewtonFeedback(
		Rational{Num: big.NewInt(2), Den: big.NewInt(1)},
		identityUnaryLFT().emit(big.NewInt(1)).emit(big.NewInt(2)),
		Rational{},
	)
	if err != nil {
		t.Fatalf("exactPositiveSqrtNewtonFeedback() error: %v", err)
	}
	if !ok {
		t.Fatalf("exactPositiveSqrtNewtonFeedback() should succeed")
	}
	if got.Candidate == nil {
		t.Fatalf("candidate should be non-nil")
	}

	assertBigIntString(t, "got.Candidate", got.Candidate, "2")
	assertBigIntString(t, "got.UpperBound.Num", got.UpperBound.Num, "17")
	assertBigIntString(t, "got.UpperBound.Den", got.UpperBound.Den, "12")

	if got.ImageRange.IsExact() {
		t.Fatalf("image range should not be exact")
	}
	assertBigIntString(t, "got.ImageRange.Lo.Value.Num", got.ImageRange.Lo.Value.Num, "7")
	assertBigIntString(t, "got.ImageRange.Lo.Value.Den", got.ImageRange.Lo.Value.Den, "3")
	assertBigIntString(t, "got.ImageRange.Hi.Value.Num", got.ImageRange.Hi.Value.Num, "5")
	assertBigIntString(t, "got.ImageRange.Hi.Value.Den", got.ImageRange.Hi.Value.Den, "2")
	if got.ImageRange.Lo.Closed || got.ImageRange.Hi.Closed {
		t.Fatalf("post-emit image range should stay open")
	}
}

func TestWB_SqrtStrategy_Emit_PreservesCachedUpperBound(t *testing.T) {
	s := sqrtStrategy{
		post:  identityUnaryLFT(),
		upper: Rational{Num: big.NewInt(17), Den: big.NewInt(12)},
	}

	next, ok := s.Emit(big.NewInt(2)).(sqrtStrategy)
	if !ok {
		t.Fatalf("Emit() should return sqrtStrategy")
	}

	assertBigIntString(t, "next.upper.Num", next.upper.Num, "17")
	assertBigIntString(t, "next.upper.Den", next.upper.Den, "12")
}

func TestWB_ExactPositiveSqrtNewtonFeedback_ExactTwoAfterEmitOneTwoTwo_CertifiesTwoUsingUpperFiveHundredSeventySevenOverFourHundredEight(t *testing.T) {
	got, ok, err := exactPositiveSqrtNewtonFeedback(
		Rational{Num: big.NewInt(2), Den: big.NewInt(1)},
		identityUnaryLFT().emit(big.NewInt(1)).emit(big.NewInt(2)).emit(big.NewInt(2)),
		Rational{Num: big.NewInt(17), Den: big.NewInt(12)},
	)
	if err != nil {
		t.Fatalf("exactPositiveSqrtNewtonFeedback() error: %v", err)
	}
	if !ok {
		t.Fatalf("exactPositiveSqrtNewtonFeedback() should succeed")
	}
	if got.Candidate == nil {
		t.Fatalf("candidate should be non-nil")
	}

	assertBigIntString(t, "got.Candidate", got.Candidate, "2")
	assertBigIntString(t, "got.UpperBound.Num", got.UpperBound.Num, "577")
	assertBigIntString(t, "got.UpperBound.Den", got.UpperBound.Den, "408")
}

func TestWB_ExactPositiveSqrtNewtonFeedback_ExactTwoAfterEmitOneTwoTwoTwo_CertifiesTwoUsingUpperSixSixFiveEightFiveSevenOverFourSevenZeroEightThreeTwo(t *testing.T) {
	got, ok, err := exactPositiveSqrtNewtonFeedback(
		Rational{Num: big.NewInt(2), Den: big.NewInt(1)},
		identityUnaryLFT().emit(big.NewInt(1)).emit(big.NewInt(2)).emit(big.NewInt(2)).emit(big.NewInt(2)),
		Rational{Num: big.NewInt(577), Den: big.NewInt(408)},
	)
	if err != nil {
		t.Fatalf("exactPositiveSqrtNewtonFeedback() error: %v", err)
	}
	if !ok {
		t.Fatalf("exactPositiveSqrtNewtonFeedback() should succeed")
	}
	if got.Candidate == nil {
		t.Fatalf("candidate should be non-nil")
	}

	assertBigIntString(t, "got.Candidate", got.Candidate, "2")
	assertBigIntString(t, "got.UpperBound.Num", got.UpperBound.Num, "665857")
	assertBigIntString(t, "got.UpperBound.Den", got.UpperBound.Den, "470832")
}

// cf/sqrt_newton_wb_test.go v1
