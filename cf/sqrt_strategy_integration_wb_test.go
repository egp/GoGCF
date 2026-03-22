// cf/sqrt_strategy_integration_wb_test.go v1
package cf

import (
	"math/big"
	"testing"
)

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

func TestWB_StrategyUnaryGCF_WithSqrtStrategy_ClosedTwoToThree_EmitsOne(t *testing.T) {
	g := strategyUnaryGCF{
		strategy: sqrtStrategy{},
		child: scriptedRangeGCF{
			term: RCFTerm{Kind: RCFValue, Value: big.NewInt(99)},
			r: Range{
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
	if got, want := value.String(), "1"; got != want {
		t.Fatalf("first term = %s, want %s", got, want)
	}
}

func TestWB_StrategyUnaryGCF_WithSqrtStrategy_ClosedTwoToNineOverFour_FirstTwoTermsAre_1_2(t *testing.T) {
	g := strategyUnaryGCF{
		strategy: sqrtStrategy{},
		child: scriptedRangeGCF{
			term: RCFTerm{Kind: RCFValue, Value: big.NewInt(99)},
			r: Range{
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
			},
		},
	}

	want := []string{"1", "2"}
	cur := GCF(g)
	for i, w := range want {
		term, rest, err := cur.Next()
		if err != nil {
			t.Fatalf("Next() #%d error: %v", i+1, err)
		}
		if !term.IsValue() {
			t.Fatalf("term #%d should be a value", i+1)
		}
		value, ok := term.BigInt()
		if !ok {
			t.Fatalf("term #%d should expose a BigInt", i+1)
		}
		if got := value.String(); got != w {
			t.Fatalf("term #%d = %s, want %s", i+1, got, w)
		}
		cur = rest
	}
}

func TestWB_StrategyUnaryGCF_WithSqrtStrategy_OpenFortyNineOverTwentyFiveToOneHundredOverFortyNine_FirstThreeTermsAre_1_2_2(t *testing.T) {
	g := strategyUnaryGCF{
		strategy: sqrtStrategy{},
		child: scriptedRangeGCF{
			term: RCFTerm{Kind: RCFValue, Value: big.NewInt(99)},
			r: Range{
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
			},
		},
	}

	want := []string{"1", "2", "2"}
	cur := GCF(g)
	for i, w := range want {
		term, rest, err := cur.Next()
		if err != nil {
			t.Fatalf("Next() #%d error: %v", i+1, err)
		}
		if !term.IsValue() {
			t.Fatalf("term #%d should be a value", i+1)
		}
		value, ok := term.BigInt()
		if !ok {
			t.Fatalf("term #%d should expose a BigInt", i+1)
		}
		if got := value.String(); got != w {
			t.Fatalf("term #%d = %s, want %s", i+1, got, w)
		}
		cur = rest
	}
}

func TestWB_SqrtOfSquare_Range_OnClosedOneToTwo_PreservesOriginalInterval(t *testing.T) {
	x := scriptedRangeGCF{
		term: RCFTerm{Kind: RCFValue, Value: big.NewInt(99)},
		r: Range{
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
		},
	}

	zr := Sqrt(Square(x)).Range()
	if !rangeWellFormed(zr) {
		t.Fatalf("Range() should be well formed")
	}
	if zr.IsExact() {
		t.Fatalf("Range() should not be exact")
	}
	if zr.Outside {
		t.Fatalf("Range() should return an inside range")
	}

	assertBigIntString(t, "zr.Lo.Value.Num", zr.Lo.Value.Num, "1")
	assertBigIntString(t, "zr.Lo.Value.Den", zr.Lo.Value.Den, "1")
	assertBigIntString(t, "zr.Hi.Value.Num", zr.Hi.Value.Num, "2")
	assertBigIntString(t, "zr.Hi.Value.Den", zr.Hi.Value.Den, "1")
	if !zr.Lo.Closed || !zr.Hi.Closed {
		t.Fatalf("sqrt(square([1,2])) should preserve the original closed interval [1,2]")
	}
}

func TestWB_SquareOfSqrt_Range_OnClosedOneToTwo_PreservesOriginalInterval(t *testing.T) {
	x := scriptedRangeGCF{
		term: RCFTerm{Kind: RCFValue, Value: big.NewInt(99)},
		r: Range{
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
		},
	}

	zr := Square(Sqrt(x)).Range()
	if !rangeWellFormed(zr) {
		t.Fatalf("Range() should be well formed")
	}
	if zr.IsExact() {
		t.Fatalf("Range() should not be exact")
	}
	if zr.Outside {
		t.Fatalf("Range() should return an inside range")
	}

	assertBigIntString(t, "zr.Lo.Value.Num", zr.Lo.Value.Num, "1")
	assertBigIntString(t, "zr.Lo.Value.Den", zr.Lo.Value.Den, "1")
	assertBigIntString(t, "zr.Hi.Value.Num", zr.Hi.Value.Num, "2")
	assertBigIntString(t, "zr.Hi.Value.Den", zr.Hi.Value.Den, "1")
	if !zr.Lo.Closed || !zr.Hi.Closed {
		t.Fatalf("square(sqrt([1,2])) should preserve the original closed interval [1,2]")
	}
}

func TestWB_SquareOfSqrt_Next_OnClosedThreeHalvesToSevenQuarters_EmitsOne(t *testing.T) {
	x := scriptedRangeGCF{
		term: RCFTerm{Kind: RCFValue, Value: big.NewInt(99)},
		r: Range{
			Lo: Bound{
				Kind: BoundFinite,
				Value: Rational{
					Num: big.NewInt(3),
					Den: big.NewInt(2),
				},
				Closed: true,
			},
			Hi: Bound{
				Kind: BoundFinite,
				Value: Rational{
					Num: big.NewInt(7),
					Den: big.NewInt(4),
				},
				Closed: true,
			},
			Outside: false,
		},
	}

	g := Square(Sqrt(x))

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
	if got, want := value.String(), "1"; got != want {
		t.Fatalf("first term = %s, want %s", got, want)
	}
}

func TestWB_SquareOfSqrt_Next_OnClosedThreeHalvesToSevenQuarters_LeavesClosedFourThirdsToTwoRemainder(t *testing.T) {
	x := scriptedRangeGCF{
		term: RCFTerm{Kind: RCFValue, Value: big.NewInt(99)},
		r: Range{
			Lo: Bound{
				Kind: BoundFinite,
				Value: Rational{
					Num: big.NewInt(3),
					Den: big.NewInt(2),
				},
				Closed: true,
			},
			Hi: Bound{
				Kind: BoundFinite,
				Value: Rational{
					Num: big.NewInt(7),
					Den: big.NewInt(4),
				},
				Closed: true,
			},
			Outside: false,
		},
	}

	g := Square(Sqrt(x))

	term, rest, err := g.Next()
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
	if got, want := value.String(), "1"; got != want {
		t.Fatalf("first term = %s, want %s", got, want)
	}

	zr := rest.Range()
	if !rangeWellFormed(zr) {
		t.Fatalf("remainder Range() should be well formed")
	}
	if zr.IsExact() {
		t.Fatalf("remainder Range() should not be exact")
	}
	if zr.Outside {
		t.Fatalf("remainder Range() should return an inside range")
	}

	assertBigIntString(t, "zr.Lo.Value.Num", zr.Lo.Value.Num, "4")
	assertBigIntString(t, "zr.Lo.Value.Den", zr.Lo.Value.Den, "3")
	assertBigIntString(t, "zr.Hi.Value.Num", zr.Hi.Value.Num, "2")
	assertBigIntString(t, "zr.Hi.Value.Den", zr.Hi.Value.Den, "1")
	if !zr.Lo.Closed || !zr.Hi.Closed {
		t.Fatalf("remainder of [3/2,7/4] after emitting 1 should be the closed interval [4/3,2]")
	}
}

func TestWB_SquareOfSqrt_Next_OnOpenFourThirdsToClosedThreeHalves_FirstTwoTermsAre_1_2(t *testing.T) {
	x := scriptedRangeGCF{
		term: RCFTerm{Kind: RCFValue, Value: big.NewInt(99)},
		r: Range{
			Lo: Bound{
				Kind: BoundFinite,
				Value: Rational{
					Num: big.NewInt(4),
					Den: big.NewInt(3),
				},
				Closed: false,
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
		},
	}

	g := Square(Sqrt(x))

	want := []string{"1", "2"}
	cur := g
	for i, w := range want {
		term, rest, err := cur.Next()
		if err != nil {
			t.Fatalf("Next() #%d error: %v", i+1, err)
		}
		if !term.IsValue() {
			t.Fatalf("term #%d should be a value", i+1)
		}
		value, ok := term.BigInt()
		if !ok {
			t.Fatalf("term #%d should expose a BigInt", i+1)
		}
		if got := value.String(); got != w {
			t.Fatalf("term #%d = %s, want %s", i+1, got, w)
		}
		cur = rest
	}
}

func TestWB_SquareOfSqrt_Next_OnOpenFourThirdsToClosedThreeHalves_LeavesClosedTwoToOpenThreeRemainder(t *testing.T) {
	x := scriptedRangeGCF{
		term: RCFTerm{Kind: RCFValue, Value: big.NewInt(99)},
		r: Range{
			Lo: Bound{
				Kind: BoundFinite,
				Value: Rational{
					Num: big.NewInt(4),
					Den: big.NewInt(3),
				},
				Closed: false,
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
		},
	}

	g := Square(Sqrt(x))

	term, rest, err := g.Next()
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
	if got, want := value.String(), "1"; got != want {
		t.Fatalf("first term = %s, want %s", got, want)
	}

	zr := rest.Range()
	if !rangeWellFormed(zr) {
		t.Fatalf("remainder Range() should be well formed")
	}
	if zr.IsExact() {
		t.Fatalf("remainder Range() should not be exact")
	}
	if zr.Outside {
		t.Fatalf("remainder Range() should return an inside range")
	}

	assertBigIntString(t, "zr.Lo.Value.Num", zr.Lo.Value.Num, "2")
	assertBigIntString(t, "zr.Lo.Value.Den", zr.Lo.Value.Den, "1")
	assertBigIntString(t, "zr.Hi.Value.Num", zr.Hi.Value.Num, "3")
	assertBigIntString(t, "zr.Hi.Value.Den", zr.Hi.Value.Den, "1")
	if !zr.Lo.Closed {
		t.Fatalf("lower endpoint 2 should be closed because x=3/2 is included")
	}
	if zr.Hi.Closed {
		t.Fatalf("upper endpoint 3 should be open because x=4/3 is excluded")
	}
}

func TestWB_SquareOfSqrt_Next_OnClosedSevenFifthsToOpenTenSevenths_FirstThreeTermsAre_1_2_2(t *testing.T) {
	x := scriptedRangeGCF{
		term: RCFTerm{Kind: RCFValue, Value: big.NewInt(99)},
		r: Range{
			Lo: Bound{
				Kind: BoundFinite,
				Value: Rational{
					Num: big.NewInt(7),
					Den: big.NewInt(5),
				},
				Closed: true,
			},
			Hi: Bound{
				Kind: BoundFinite,
				Value: Rational{
					Num: big.NewInt(10),
					Den: big.NewInt(7),
				},
				Closed: false,
			},
			Outside: false,
		},
	}

	g := Square(Sqrt(x))

	want := []string{"1", "2", "2"}
	cur := g
	for i, w := range want {
		term, rest, err := cur.Next()
		if err != nil {
			t.Fatalf("Next() #%d error: %v", i+1, err)
		}
		if !term.IsValue() {
			t.Fatalf("term #%d should be a value", i+1)
		}
		value, ok := term.BigInt()
		if !ok {
			t.Fatalf("term #%d should expose a BigInt", i+1)
		}
		if got := value.String(); got != w {
			t.Fatalf("term #%d = %s, want %s", i+1, got, w)
		}
		cur = rest
	}
}

func TestWB_SquareOfSqrt_Next_OnClosedSevenFifthsToOpenTenSevenths_AfterTwoTermsLeavesClosedTwoToOpenThreeRemainder(t *testing.T) {
	x := scriptedRangeGCF{
		term: RCFTerm{Kind: RCFValue, Value: big.NewInt(99)},
		r: Range{
			Lo: Bound{
				Kind: BoundFinite,
				Value: Rational{
					Num: big.NewInt(7),
					Den: big.NewInt(5),
				},
				Closed: true,
			},
			Hi: Bound{
				Kind: BoundFinite,
				Value: Rational{
					Num: big.NewInt(10),
					Den: big.NewInt(7),
				},
				Closed: false,
			},
			Outside: false,
		},
	}

	g := Square(Sqrt(x))

	term1, rest, err := g.Next()
	if err != nil {
		t.Fatalf("first Next() error: %v", err)
	}
	if !term1.IsValue() {
		t.Fatalf("first term should be a value")
	}
	value1, ok := term1.BigInt()
	if !ok {
		t.Fatalf("first term should expose a BigInt")
	}
	if got, want := value1.String(), "1"; got != want {
		t.Fatalf("first term = %s, want %s", got, want)
	}

	term2, rest, err := rest.Next()
	if err != nil {
		t.Fatalf("second Next() error: %v", err)
	}
	if !term2.IsValue() {
		t.Fatalf("second term should be a value")
	}
	value2, ok := term2.BigInt()
	if !ok {
		t.Fatalf("second term should expose a BigInt")
	}
	if got, want := value2.String(), "2"; got != want {
		t.Fatalf("second term = %s, want %s", got, want)
	}

	zr := rest.Range()
	if !rangeWellFormed(zr) {
		t.Fatalf("remainder Range() should be well formed")
	}
	if zr.IsExact() {
		t.Fatalf("remainder Range() should not be exact")
	}
	if zr.Outside {
		t.Fatalf("remainder Range() should return an inside range")
	}

	assertBigIntString(t, "zr.Lo.Value.Num", zr.Lo.Value.Num, "2")
	assertBigIntString(t, "zr.Lo.Value.Den", zr.Lo.Value.Den, "1")
	assertBigIntString(t, "zr.Hi.Value.Num", zr.Hi.Value.Num, "3")
	assertBigIntString(t, "zr.Hi.Value.Den", zr.Hi.Value.Den, "1")
	if !zr.Lo.Closed {
		t.Fatalf("lower endpoint 2 should be closed because x=7/5 is included")
	}
	if zr.Hi.Closed {
		t.Fatalf("upper endpoint 3 should be open because x=10/7 is excluded")
	}
}

// cf/sqrt_strategy_integration_wb_test.go v1
