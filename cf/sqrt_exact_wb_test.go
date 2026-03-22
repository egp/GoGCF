// cf/sqrt_exact_wb_test.go v1
package cf

import (
	"math/big"
	"testing"
)

func TestWB_SqrtStrategy_ExactEval_Zero_IsZero(t *testing.T) {
	s := sqrtStrategy{}

	r, ok, err := s.ExactEval(Rational{
		Num: big.NewInt(0),
		Den: big.NewInt(1),
	})
	if err != nil {
		t.Fatalf("ExactEval() error: %v", err)
	}
	if !ok {
		t.Fatalf("ExactEval() should succeed for zero")
	}

	assertBigIntString(t, "r.Num", r.Num, "0")
	assertBigIntString(t, "r.Den", r.Den, "1")
}

func TestWB_SqrtStrategy_ExactEval_Four_IsTwo(t *testing.T) {
	s := sqrtStrategy{}

	r, ok, err := s.ExactEval(Rational{
		Num: big.NewInt(4),
		Den: big.NewInt(1),
	})
	if err != nil {
		t.Fatalf("ExactEval() error: %v", err)
	}
	if !ok {
		t.Fatalf("ExactEval() should succeed for four")
	}

	assertBigIntString(t, "r.Num", r.Num, "2")
	assertBigIntString(t, "r.Den", r.Den, "1")
}

func TestWB_SqrtStrategy_ExactEval_NineOverFour_IsThreeOverTwo(t *testing.T) {
	s := sqrtStrategy{}

	r, ok, err := s.ExactEval(Rational{
		Num: big.NewInt(9),
		Den: big.NewInt(4),
	})
	if err != nil {
		t.Fatalf("ExactEval() error: %v", err)
	}
	if !ok {
		t.Fatalf("ExactEval() should succeed for nine over four")
	}

	assertBigIntString(t, "r.Num", r.Num, "3")
	assertBigIntString(t, "r.Den", r.Den, "2")
}

func TestWB_SqrtStrategy_ExactEval_NegativeOne_IsUndefined(t *testing.T) {
	s := sqrtStrategy{}

	_, _, err := s.ExactEval(Rational{
		Num: big.NewInt(-1),
		Den: big.NewInt(1),
	})
	if err != ErrUndefined {
		t.Fatalf("ExactEval() error = %v, want %v", err, ErrUndefined)
	}
}

func TestWB_SqrtStrategy_RangeFromOperand_ExactFour_IsExactTwo(t *testing.T) {
	s := sqrtStrategy{}

	xr := exactRangeFromRational(Rational{
		Num: big.NewInt(4),
		Den: big.NewInt(1),
	})

	zr, err := s.RangeFromOperand(xr)
	if err != nil {
		t.Fatalf("RangeFromOperand() error: %v", err)
	}
	if !zr.IsExact() {
		t.Fatalf("RangeFromOperand() should return an exact point")
	}

	assertBigIntString(t, "zr.Lo.Value.Num", zr.Lo.Value.Num, "2")
	assertBigIntString(t, "zr.Lo.Value.Den", zr.Lo.Value.Den, "1")
	assertBigIntString(t, "zr.Hi.Value.Num", zr.Hi.Value.Num, "2")
	assertBigIntString(t, "zr.Hi.Value.Den", zr.Hi.Value.Den, "1")
}

func TestWB_SqrtStrategy_RangeFromOperand_ExactTwo_IsOpenOneToTwo(t *testing.T) {
	s := sqrtStrategy{}

	xr := exactRangeFromRational(Rational{
		Num: big.NewInt(2),
		Den: big.NewInt(1),
	})

	zr, err := s.RangeFromOperand(xr)
	if err != nil {
		t.Fatalf("RangeFromOperand() error: %v", err)
	}
	if zr.IsExact() {
		t.Fatalf("RangeFromOperand() should not be exact")
	}
	if zr.Outside {
		t.Fatalf("RangeFromOperand() should return an inside range")
	}

	assertBigIntString(t, "zr.Lo.Value.Num", zr.Lo.Value.Num, "1")
	assertBigIntString(t, "zr.Lo.Value.Den", zr.Lo.Value.Den, "1")
	assertBigIntString(t, "zr.Hi.Value.Num", zr.Hi.Value.Num, "2")
	assertBigIntString(t, "zr.Hi.Value.Den", zr.Hi.Value.Den, "1")
	if zr.Lo.Closed || zr.Hi.Closed {
		t.Fatalf("sqrt(2) should be strictly between 1 and 2")
	}
}

func TestWB_SqrtStrategy_EmitCandidateFromOperand_ExactTwo_CertifiesOne(t *testing.T) {
	s := sqrtStrategy{}

	xr := exactRangeFromRational(Rational{
		Num: big.NewInt(2),
		Den: big.NewInt(1),
	})

	q, ok, err := s.EmitCandidateFromOperand(xr)
	if err != nil {
		t.Fatalf("EmitCandidateFromOperand() error: %v", err)
	}
	if !ok {
		t.Fatalf("EmitCandidateFromOperand() should certify an output term")
	}

	assertBigIntString(t, "q", q, "1")
}

func TestWB_SqrtStrategy_EmitCandidateFromOperand_ExactTwoAfterEmitOne_CertifiesTwo(t *testing.T) {
	s := sqrtStrategy{
		post: identityUnaryLFT().emit(big.NewInt(1)),
	}

	xr := exactRangeFromRational(Rational{
		Num: big.NewInt(2),
		Den: big.NewInt(1),
	})

	q, ok, err := s.EmitCandidateFromOperand(xr)
	if err != nil {
		t.Fatalf("EmitCandidateFromOperand() error: %v", err)
	}
	if !ok {
		t.Fatalf("EmitCandidateFromOperand() should certify an output term")
	}

	assertBigIntString(t, "q", q, "2")
}

func TestWB_SqrtStrategy_RangeFromOperand_ExactTwoAfterEmitOne_IsInsideTwoToThree(t *testing.T) {
	s := sqrtStrategy{
		post: identityUnaryLFT().emit(big.NewInt(1)),
	}

	xr := exactRangeFromRational(Rational{
		Num: big.NewInt(2),
		Den: big.NewInt(1),
	})

	zr, err := s.RangeFromOperand(xr)
	if err != nil {
		t.Fatalf("RangeFromOperand() error: %v", err)
	}
	if zr.IsExact() {
		t.Fatalf("RangeFromOperand() should not be exact")
	}
	if zr.Outside {
		t.Fatalf("RangeFromOperand() should return an inside range")
	}

	assertBigIntString(t, "zr.Lo.Value.Num", zr.Lo.Value.Num, "2")
	assertBigIntString(t, "zr.Lo.Value.Den", zr.Lo.Value.Den, "1")
	assertBigIntString(t, "zr.Hi.Value.Num", zr.Hi.Value.Num, "3")
	assertBigIntString(t, "zr.Hi.Value.Den", zr.Hi.Value.Den, "1")
	if zr.Lo.Closed || zr.Hi.Closed {
		t.Fatalf("post-emit sqrt(2) remainder should stay open")
	}
}

func TestWB_SqrtStrategy_EmitCandidateFromOperand_ExactTwoAfterEmitOneTwo_CertifiesTwo(t *testing.T) {
	s := sqrtStrategy{
		post: identityUnaryLFT().emit(big.NewInt(1)).emit(big.NewInt(2)),
	}

	xr := exactRangeFromRational(Rational{
		Num: big.NewInt(2),
		Den: big.NewInt(1),
	})

	q, ok, err := s.EmitCandidateFromOperand(xr)
	if err != nil {
		t.Fatalf("EmitCandidateFromOperand() error: %v", err)
	}
	if !ok {
		t.Fatalf("EmitCandidateFromOperand() should certify an output term")
	}

	assertBigIntString(t, "q", q, "2")
}

func TestWB_SqrtStrategy_RangeFromOperand_ExactTwoAfterEmitOneTwo_IsInsideSevenThirdsToFiveHalves(t *testing.T) {
	s := sqrtStrategy{
		post: identityUnaryLFT().emit(big.NewInt(1)).emit(big.NewInt(2)),
	}

	xr := exactRangeFromRational(Rational{
		Num: big.NewInt(2),
		Den: big.NewInt(1),
	})

	zr, err := s.RangeFromOperand(xr)
	if err != nil {
		t.Fatalf("RangeFromOperand() error: %v", err)
	}
	if zr.IsExact() {
		t.Fatalf("RangeFromOperand() should not be exact")
	}
	if zr.Outside {
		t.Fatalf("RangeFromOperand() should return an inside range")
	}

	assertBigIntString(t, "zr.Lo.Value.Num", zr.Lo.Value.Num, "7")
	assertBigIntString(t, "zr.Lo.Value.Den", zr.Lo.Value.Den, "3")
	assertBigIntString(t, "zr.Hi.Value.Num", zr.Hi.Value.Num, "5")
	assertBigIntString(t, "zr.Hi.Value.Den", zr.Hi.Value.Den, "2")
	if zr.Lo.Closed || zr.Hi.Closed {
		t.Fatalf("post-emit sqrt(2) remainder should stay open")
	}
}

func TestWB_SqrtStrategy_EmitCandidateFromOperand_ExactTwoAfterEmitOneTwoTwo_CertifiesTwo(t *testing.T) {
	s := sqrtStrategy{
		post:  identityUnaryLFT().emit(big.NewInt(1)).emit(big.NewInt(2)).emit(big.NewInt(2)),
		upper: Rational{Num: big.NewInt(17), Den: big.NewInt(12)},
	}

	xr := exactRangeFromRational(Rational{
		Num: big.NewInt(2),
		Den: big.NewInt(1),
	})

	q, ok, err := s.EmitCandidateFromOperand(xr)
	if err != nil {
		t.Fatalf("EmitCandidateFromOperand() error: %v", err)
	}
	if !ok {
		t.Fatalf("EmitCandidateFromOperand() should certify an output term")
	}

	assertBigIntString(t, "q", q, "2")
	assertBigIntString(t, "s.upper.Num", s.upper.Num, "17")
	assertBigIntString(t, "s.upper.Den", s.upper.Den, "12")
}

// cf/sqrt_exact_wb_test.go v1
