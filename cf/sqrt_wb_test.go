// cf/sqrt_wb_test.go v1
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

func TestWB_StrategyUnaryGCF_WithSqrtStrategy_OpenFourToTwentyFiveOverFour_EmitsTwo(t *testing.T) {
	g := strategyUnaryGCF{
		strategy: sqrtStrategy{},
		child: scriptedRangeGCF{
			term: RCFTerm{Kind: RCFValue, Value: big.NewInt(99)},
			r: Range{
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
			},
		},
	}

	term, _, err := g.Next()
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

// cf/sqrt_wb_test.go v1
