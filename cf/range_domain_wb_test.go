// cf/range_domain_wb_test.go v1
package cf

import (
	"math/big"
	"testing"
)

func TestWB_RangeMaybeIncludesZero_ClosedZeroToThreeHalves_IsTrue(t *testing.T) {
	r := Range{
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
				Num: big.NewInt(3),
				Den: big.NewInt(2),
			},
			Closed: true,
		},
		Outside: false,
	}

	if !rangeMaybeIncludesZero(r) {
		t.Fatalf("closed interval touching zero should maybe include zero")
	}
}

func TestWB_RangeMaybeIncludesZero_OpenPositiveInterval_IsFalse(t *testing.T) {
	r := Range{
		Lo: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(1),
				Den: big.NewInt(10),
			},
			Closed: false,
		},
		Hi: Bound{
			Kind: BoundFinite,
			Value: Rational{
				Num: big.NewInt(3),
				Den: big.NewInt(2),
			},
			Closed: false,
		},
		Outside: false,
	}

	if rangeMaybeIncludesZero(r) {
		t.Fatalf("open strictly positive interval should not maybe include zero")
	}
}

func TestWB_RangeMaybeIncludesZero_OutsideMinusOneToOne_IsFalse(t *testing.T) {
	r := Range{
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
				Num: big.NewInt(1),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Outside: true,
	}

	if rangeMaybeIncludesZero(r) {
		t.Fatalf("outside [-1,1] should exclude zero")
	}
}

func TestWB_RangeCertainlyNegative_ExactMinusOne_IsTrue(t *testing.T) {
	r := Range{
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
				Num: big.NewInt(-1),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Outside: false,
	}

	if !rangeCertainlyNegative(r) {
		t.Fatalf("exact -1 should be certainly negative")
	}
}

func TestWB_RangeCertainlyNonNegative_ClosedZeroToThreeHalves_IsTrue(t *testing.T) {
	r := Range{
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
				Num: big.NewInt(3),
				Den: big.NewInt(2),
			},
			Closed: true,
		},
		Outside: false,
	}

	if !rangeCertainlyNonNegative(r) {
		t.Fatalf("[0, 3/2] should be certainly nonnegative")
	}
}

func TestWB_RangeCertainlyNonNegative_StraddlingZero_IsFalse(t *testing.T) {
	r := Range{
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
				Num: big.NewInt(1),
				Den: big.NewInt(1),
			},
			Closed: true,
		},
		Outside: false,
	}

	if rangeCertainlyNonNegative(r) {
		t.Fatalf("[-1, 1] should not be certainly nonnegative")
	}
}

// cf/range_domain_wb_test.go v1
