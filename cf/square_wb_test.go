// cf/square_wb_test.go v1
package cf

import (
	"math/big"
	"testing"
)

func TestWB_SquareDiagLFT_HasCoefficients_1_0_0_0_0_1(t *testing.T) {
	d := squareDiagLFT()

	assertBigIntString(t, "d.a", d.a, "1")
	assertBigIntString(t, "d.b", d.b, "0")
	assertBigIntString(t, "d.c", d.c, "0")
	assertBigIntString(t, "d.d", d.d, "0")
	assertBigIntString(t, "d.e", d.e, "0")
	assertBigIntString(t, "d.f", d.f, "1")
}

func TestWB_SquareStrategy_ExactEval_ThreeHalves_IsNineOverFour(t *testing.T) {
	s := squareStrategy{}

	r, ok, err := s.ExactEval(Rational{
		Num: big.NewInt(3),
		Den: big.NewInt(2),
	})
	if err != nil {
		t.Fatalf("ExactEval() error: %v", err)
	}
	if !ok {
		t.Fatalf("ExactEval() should succeed")
	}

	assertBigIntString(t, "r.Num", r.Num, "9")
	assertBigIntString(t, "r.Den", r.Den, "4")
}

func TestWB_SquareStrategy_RangeFromOperand_ClosedOneToThreeHalves_IsClosedOneToNineOverFour(t *testing.T) {
	s := squareStrategy{}

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
				Num: big.NewInt(3),
				Den: big.NewInt(2),
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
	assertBigIntString(t, "zr.Hi.Value.Num", zr.Hi.Value.Num, "9")
	assertBigIntString(t, "zr.Hi.Value.Den", zr.Hi.Value.Den, "4")
	if !zr.Lo.Closed || !zr.Hi.Closed {
		t.Fatalf("closed positive interval should stay closed under square")
	}
}

func TestWB_SquareStrategy_RangeFromOperand_ClosedMinusOneToTwo_IsClosedZeroToFour(t *testing.T) {
	s := squareStrategy{}

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
	assertBigIntString(t, "zr.Hi.Value.Num", zr.Hi.Value.Num, "4")
	assertBigIntString(t, "zr.Hi.Value.Den", zr.Hi.Value.Den, "1")
	if !zr.Lo.Closed || !zr.Hi.Closed {
		t.Fatalf("closed interval crossing zero should square to a closed interval")
	}
}

// cf/square_wb_test.go v1
