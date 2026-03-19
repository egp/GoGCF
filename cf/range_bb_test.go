// cf/range_bb_test.go v1
package cf_test

import (
	"math/big"
	"testing"

	"github.com/egp/GoGCF/cf"
)

func TestBB_Range_ExactClosedFinitePoint(t *testing.T) {
	r := cf.Range{
		Lo:      finiteClosedBound(3, 2),
		Hi:      finiteClosedBound(3, 2),
		Outside: false,
	}

	if !r.IsInside() {
		t.Fatalf("exact closed finite point should be inside")
	}
	if r.IsOutside() {
		t.Fatalf("exact closed finite point should not be outside")
	}
	if !r.IsExact() {
		t.Fatalf("exact closed finite point should be exact")
	}
}

func TestBB_Range_OutsideFlagControlsInsideOutside(t *testing.T) {
	r := cf.Range{
		Lo:      finiteClosedBound(1, 1),
		Hi:      finiteClosedBound(2, 1),
		Outside: true,
	}

	if r.IsInside() {
		t.Fatalf("outside range should not be inside")
	}
	if !r.IsOutside() {
		t.Fatalf("outside range should be outside")
	}
	if r.IsExact() {
		t.Fatalf("outside range should not be exact")
	}
}

func TestBB_Range_HalfOpenSamePointIsNotExact(t *testing.T) {
	r := cf.Range{
		Lo:      finiteOpenBound(5, 1),
		Hi:      finiteClosedBound(5, 1),
		Outside: false,
	}

	if !r.IsInside() {
		t.Fatalf("half-open same-point range should still be inside")
	}
	if r.IsOutside() {
		t.Fatalf("half-open same-point range should not be outside")
	}
	if r.IsExact() {
		t.Fatalf("half-open same-point range should not be exact")
	}
}

func TestBB_Range_InfiniteEndpointIsNotExact(t *testing.T) {
	r := cf.Range{
		Lo: cf.Bound{
			Kind:   cf.BoundNegInf,
			Closed: false,
		},
		Hi:      finiteClosedBound(9, 1),
		Outside: false,
	}

	if !r.IsInside() {
		t.Fatalf("range with infinite endpoint should still be inside when Outside is false")
	}
	if r.IsExact() {
		t.Fatalf("range with infinite endpoint should not be exact")
	}
}

func TestBB_FromSource_Int64_RangeIsExactPoint(t *testing.T) {
	g := cf.FromSource(cf.Int64(7))

	r := g.Range()

	if !r.IsInside() {
		t.Fatalf("Int64 evaluator range should be inside")
	}
	if r.IsOutside() {
		t.Fatalf("Int64 evaluator range should not be outside")
	}
	if !r.IsExact() {
		t.Fatalf("Int64 evaluator range should be exact")
	}

	if r.Lo.Kind != cf.BoundFinite || r.Hi.Kind != cf.BoundFinite {
		t.Fatalf("Int64 evaluator range should have finite bounds")
	}
	if !r.Lo.Closed || !r.Hi.Closed {
		t.Fatalf("Int64 evaluator range should have closed bounds")
	}

	if got, want := r.Lo.Value.Num.String(), "7"; got != want {
		t.Fatalf("lo numerator = %s, want %s", got, want)
	}
	if got, want := r.Lo.Value.Den.String(), "1"; got != want {
		t.Fatalf("lo denominator = %s, want %s", got, want)
	}
	if got, want := r.Hi.Value.Num.String(), "7"; got != want {
		t.Fatalf("hi numerator = %s, want %s", got, want)
	}
	if got, want := r.Hi.Value.Den.String(), "1"; got != want {
		t.Fatalf("hi denominator = %s, want %s", got, want)
	}
}

func finiteClosedBound(num, den int64) cf.Bound {
	return cf.Bound{
		Kind: cf.BoundFinite,
		Value: cf.Rational{
			Num: int64Big(num),
			Den: int64Big(den),
		},
		Closed: true,
	}
}

func finiteOpenBound(num, den int64) cf.Bound {
	return cf.Bound{
		Kind: cf.BoundFinite,
		Value: cf.Rational{
			Num: int64Big(num),
			Den: int64Big(den),
		},
		Closed: false,
	}
}

func int64Big(v int64) *big.Int {
	return big.NewInt(v)
}

// cf/range_bb_test.go v1
