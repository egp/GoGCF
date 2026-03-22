// cf/sqrt_interval_wb_test.go v1
package cf

import (
	"math/big"
	"testing"
)

func TestWB_SqrtStrategy_RangeFromOperand_ClosedFourToNine_IsClosedTwoToThree(t *testing.T) {
	s := sqrtStrategy{}

	xr := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(4),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(9),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Outside: false,
	}

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
	if !zr.Lo.Closed || !zr.Hi.Closed {
		t.Fatalf("closed interval should stay closed under exact endpoint square roots")
	}
}

func TestWB_SqrtStrategy_RangeFromOperand_OpenFourToNine_IsOpenTwoToThree(t *testing.T) {
	s := sqrtStrategy{}

	xr := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(4),
				Den: big.NewInt(1),
			},
			Closed: false,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(9),
				Den: big.NewInt(1),
			},
			Closed: false,
		},
		Outside: false,
	}

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
		t.Fatalf("open interval should stay open under exact endpoint square roots")
	}
}

func TestWB_SqrtStrategy_RangeFromOperand_CertainlyNegativeInterval_IsUndefined(t *testing.T) {
	s := sqrtStrategy{}

	xr := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(-9),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(-4),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Outside: false,
	}

	_, err := s.RangeFromOperand(xr)
	if err != ErrUndefined {
		t.Fatalf("RangeFromOperand() error = %v, want %v", err, ErrUndefined)
	}
}

func TestWB_SqrtStrategy_RangeFromOperand_StraddlingZero_ReturnsUnknownWithoutError(t *testing.T) {
	s := sqrtStrategy{}

	xr := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(-1),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(4),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Outside: false,
	}

	zr, err := s.RangeFromOperand(xr)
	if err != nil {
		t.Fatalf("RangeFromOperand() error: %v", err)
	}
	if rangeWellFormed(zr) {
		t.Fatalf("RangeFromOperand() should stay uncertified when the operand may be negative")
	}
}

func TestWB_SqrtStrategy_EmitCandidateFromOperand_OpenFourToTwentyFiveOverFour_CertifiesTwo(t *testing.T) {
	s := sqrtStrategy{}

	xr := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(4),
				Den: big.NewInt(1),
			},
			Closed: false,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(25),
				Den: big.NewInt(4),
			},
			Closed: false,
		},
		Outside: false,
	}

	q, ok, err := s.EmitCandidateFromOperand(xr)
	if err != nil {
		t.Fatalf("EmitCandidateFromOperand() error: %v", err)
	}
	if !ok {
		t.Fatalf("EmitCandidateFromOperand() should certify an output term")
	}

	assertBigIntString(t, "q", q, "2")
}

func TestWB_SqrtStrategy_RangeFromOperand_ClosedTwoToThree_IsOpenOneToTwo(t *testing.T) {
	s := sqrtStrategy{}

	xr := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(2),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(3),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Outside: false,
	}

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
		t.Fatalf("sqrt([2,3]) should stay strictly between 1 and 2")
	}
}

func TestWB_SqrtStrategy_EmitCandidateFromOperand_ClosedTwoToThree_CertifiesOne(t *testing.T) {
	s := sqrtStrategy{}

	xr := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(2),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(3),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Outside: false,
	}

	q, ok, err := s.EmitCandidateFromOperand(xr)
	if err != nil {
		t.Fatalf("EmitCandidateFromOperand() error: %v", err)
	}
	if !ok {
		t.Fatalf("EmitCandidateFromOperand() should certify an output term")
	}

	assertBigIntString(t, "q", q, "1")
}

func TestWB_SqrtStrategy_RangeFromOperand_ClosedTwoToNineOverFourAfterEmitOne_IsClosedTwoToOpenThree(t *testing.T) {
	s := sqrtStrategy{
		post: identityUnaryLFT().emit(big.NewInt(1)),
	}

	xr := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(2),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(9),
				Den: big.NewInt(4),
			},
			Closed: true,
		},
		Outside: false,
	}

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
	if !zr.Lo.Closed {
		t.Fatalf("lower endpoint 2 should be closed because sqrt(9/4)=3/2 is included")
	}
	if zr.Hi.Closed {
		t.Fatalf("upper endpoint 3 should stay open because it comes from the irrational side")
	}
}

func TestWB_SqrtStrategy_EmitCandidateFromOperand_ClosedTwoToNineOverFourAfterEmitOne_CertifiesTwo(t *testing.T) {
	s := sqrtStrategy{
		post: identityUnaryLFT().emit(big.NewInt(1)),
	}

	xr := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(2),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(9),
				Den: big.NewInt(4),
			},
			Closed: true,
		},
		Outside: false,
	}

	q, ok, err := s.EmitCandidateFromOperand(xr)
	if err != nil {
		t.Fatalf("EmitCandidateFromOperand() error: %v", err)
	}
	if !ok {
		t.Fatalf("EmitCandidateFromOperand() should certify an output term")
	}

	assertBigIntString(t, "q", q, "2")
}

func TestWB_SqrtStrategy_RangeFromOperand_ClosedOneToTwo_PreservesClosedLowerEndpoint(t *testing.T) {
	s := sqrtStrategy{}

	xr := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(1),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(2),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Outside: false,
	}

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

	if !zr.Lo.Closed {
		t.Fatalf("sqrt([1,2]) should preserve the exact lower endpoint 1 as closed")
	}
	if zr.Hi.Closed {
		t.Fatalf("sqrt([1,2]) upper endpoint should stay open because sqrt(2) is irrational")
	}
}

func TestWB_SqrtStrategy_RangeFromOperand_ClosedZeroToTwo_PreservesClosedZeroLowerEndpoint(t *testing.T) {
	s := sqrtStrategy{}

	xr := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(0),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(2),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Outside: false,
	}

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

	assertBigIntString(t, "zr.Lo.Value.Num", zr.Lo.Value.Num, "0")
	assertBigIntString(t, "zr.Lo.Value.Den", zr.Lo.Value.Den, "1")
	assertBigIntString(t, "zr.Hi.Value.Num", zr.Hi.Value.Num, "2")
	assertBigIntString(t, "zr.Hi.Value.Den", zr.Hi.Value.Den, "1")

	if !zr.Lo.Closed {
		t.Fatalf("sqrt([0,2]) should preserve the exact lower endpoint 0 as closed")
	}
	if zr.Hi.Closed {
		t.Fatalf("sqrt([0,2]) upper endpoint should stay open because sqrt(2) is irrational")
	}
}

func TestWB_SqrtStrategy_EmitCandidateFromOperand_ClosedOneToFourAfterEmitOne_StaysUncertifiedWithoutError(t *testing.T) {
	s := sqrtStrategy{
		post: identityUnaryLFT().emit(big.NewInt(1)),
	}

	xr := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(1),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(4),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Outside: false,
	}

	q, ok, err := s.EmitCandidateFromOperand(xr)
	if err != nil {
		t.Fatalf("EmitCandidateFromOperand() error: %v", err)
	}
	if ok {
		t.Fatalf("EmitCandidateFromOperand() should stay uncertified")
	}
	if q != nil {
		t.Fatalf("candidate should be nil when uncertified")
	}
}

func TestWB_SqrtStrategy_EmitCandidateFromOperand_ClosedTwoToNineOverFourAfterEmitOneTwo_StaysUncertifiedWithoutError(t *testing.T) {
	s := sqrtStrategy{
		post: identityUnaryLFT().emit(big.NewInt(1)).emit(big.NewInt(2)),
	}

	xr := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(2),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(9),
				Den: big.NewInt(4),
			},
			Closed: true,
		},
		Outside: false,
	}

	q, ok, err := s.EmitCandidateFromOperand(xr)
	if err != nil {
		t.Fatalf("EmitCandidateFromOperand() error: %v", err)
	}
	if ok {
		t.Fatalf("EmitCandidateFromOperand() should stay uncertified once an included endpoint reaches the post-emit pole")
	}
	if q != nil {
		t.Fatalf("candidate should be nil when uncertified")
	}
}

func TestWB_SqrtStrategy_RangeFromOperand_ClosedFourToNineAfterEmitTwo_ReturnsUnknownWithoutError(t *testing.T) {
	s := sqrtStrategy{
		post: identityUnaryLFT().emit(big.NewInt(2)),
	}

	xr := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(4),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(9),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Outside: false,
	}

	zr, err := s.RangeFromOperand(xr)
	if err != nil {
		t.Fatalf("RangeFromOperand() error: %v", err)
	}
	if rangeWellFormed(zr) {
		t.Fatalf("RangeFromOperand() should stay uncertified when an included exact endpoint reaches the post-emit pole")
	}
}

func TestWB_SqrtStrategy_EmitCandidateFromOperand_ClosedFourToNineAfterEmitTwo_StaysUncertifiedWithoutError(t *testing.T) {
	s := sqrtStrategy{
		post: identityUnaryLFT().emit(big.NewInt(2)),
	}

	xr := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(4),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(9),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Outside: false,
	}

	q, ok, err := s.EmitCandidateFromOperand(xr)
	if err != nil {
		t.Fatalf("EmitCandidateFromOperand() error: %v", err)
	}
	if ok {
		t.Fatalf("EmitCandidateFromOperand() should stay uncertified when an included exact endpoint reaches the post-emit pole")
	}
	if q != nil {
		t.Fatalf("candidate should be nil when uncertified")
	}
}

func TestWB_SqrtStrategy_RangeFromOperand_OpenFortyNineOverTwentyFiveToOneHundredOverFortyNineAfterEmitOneTwo_IsOpenTwoToThree(t *testing.T) {
	s := sqrtStrategy{
		post: identityUnaryLFT().emit(big.NewInt(1)).emit(big.NewInt(2)),
	}

	xr := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(49),
				Den: big.NewInt(25),
			},
			Closed: false,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(100),
				Den: big.NewInt(49),
			},
			Closed: false,
		},
		Outside: false,
	}

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
		t.Fatalf("post-emit interval should stay open")
	}
}

func TestWB_SqrtStrategy_EmitCandidateFromOperand_OpenFortyNineOverTwentyFiveToOneHundredOverFortyNineAfterEmitOneTwo_CertifiesTwo(t *testing.T) {
	s := sqrtStrategy{
		post: identityUnaryLFT().emit(big.NewInt(1)).emit(big.NewInt(2)),
	}

	xr := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(49),
				Den: big.NewInt(25),
			},
			Closed: false,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(100),
				Den: big.NewInt(49),
			},
			Closed: false,
		},
		Outside: false,
	}

	q, ok, err := s.EmitCandidateFromOperand(xr)
	if err != nil {
		t.Fatalf("EmitCandidateFromOperand() error: %v", err)
	}
	if !ok {
		t.Fatalf("EmitCandidateFromOperand() should certify an output term")
	}

	assertBigIntString(t, "q", q, "2")
}

// cf/sqrt_interval_wb_test.go v1
