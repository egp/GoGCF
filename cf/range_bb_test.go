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

func TestBB_Range_Cmp_EqualExactPointsCompareEqual(t *testing.T) {
	a := cf.Range{
		Lo:      finiteClosedBound(3, 2),
		Hi:      finiteClosedBound(3, 2),
		Outside: false,
	}
	b := cf.Range{
		Lo:      finiteClosedBound(3, 2),
		Hi:      finiteClosedBound(3, 2),
		Outside: false,
	}

	if got := a.Cmp(b); got != 0 {
		t.Fatalf("a.Cmp(b) = %d, want 0", got)
	}
	if got := b.Cmp(a); got != 0 {
		t.Fatalf("b.Cmp(a) = %d, want 0", got)
	}
}

func TestBB_Range_Cmp_InsideNarrowBeatsInsideWide(t *testing.T) {
	narrow := cf.Range{
		Lo:      finiteClosedBound(2, 1),
		Hi:      finiteClosedBound(3, 1),
		Outside: false,
	}
	wide := cf.Range{
		Lo:      finiteClosedBound(1, 1),
		Hi:      finiteClosedBound(4, 1),
		Outside: false,
	}

	if got := narrow.Cmp(wide); got >= 0 {
		t.Fatalf("narrow.Cmp(wide) = %d, want < 0", got)
	}
	if got := wide.Cmp(narrow); got <= 0 {
		t.Fatalf("wide.Cmp(narrow) = %d, want > 0", got)
	}
}

func TestBB_Range_Cmp_OutsideWideBeatsOutsideNarrow(t *testing.T) {
	wideExcluded := cf.Range{
		Lo:      finiteClosedBound(1, 1),
		Hi:      finiteClosedBound(4, 1),
		Outside: true,
	}
	narrowExcluded := cf.Range{
		Lo:      finiteClosedBound(2, 1),
		Hi:      finiteClosedBound(3, 1),
		Outside: true,
	}

	if got := wideExcluded.Cmp(narrowExcluded); got >= 0 {
		t.Fatalf("wideExcluded.Cmp(narrowExcluded) = %d, want < 0", got)
	}
	if got := narrowExcluded.Cmp(wideExcluded); got <= 0 {
		t.Fatalf("narrowExcluded.Cmp(wideExcluded) = %d, want > 0", got)
	}
}

func TestBB_Range_Cmp_CrossClassSubset_InsideRefinementBeatsOutside(t *testing.T) {
	outside := cf.Range{
		Lo:      finiteClosedBound(-1, 1),
		Hi:      finiteClosedBound(1, 1),
		Outside: true,
	}
	insideRefinement := cf.Range{
		Lo:      finiteClosedBound(2, 1),
		Hi:      finiteClosedBound(3, 1),
		Outside: false,
	}

	if got := insideRefinement.Cmp(outside); got >= 0 {
		t.Fatalf("insideRefinement.Cmp(outside) = %d, want < 0", got)
	}
	if got := outside.Cmp(insideRefinement); got <= 0 {
		t.Fatalf("outside.Cmp(insideRefinement) = %d, want > 0", got)
	}
}

func TestBB_Range_Cmp_CrossClassSubset_OutsideRefinementBeatsInsideFullLine(t *testing.T) {
	insideFullLine := cf.Range{
		Lo:      negInfOpenBound(),
		Hi:      posInfOpenBound(),
		Outside: false,
	}
	outsideRefinement := cf.Range{
		Lo:      finiteClosedBound(-1, 1),
		Hi:      finiteClosedBound(1, 1),
		Outside: true,
	}

	if got := outsideRefinement.Cmp(insideFullLine); got >= 0 {
		t.Fatalf("outsideRefinement.Cmp(insideFullLine) = %d, want < 0", got)
	}
	if got := insideFullLine.Cmp(outsideRefinement); got <= 0 {
		t.Fatalf("insideFullLine.Cmp(outsideRefinement) = %d, want > 0", got)
	}
}

func TestBB_Range_Cmp_Fallback_InsideNarrowBeatsOutsideWideWhenIncomparable(t *testing.T) {
	insideNarrow := cf.Range{
		Lo:      finiteClosedBound(2, 1),
		Hi:      finiteClosedBound(3, 1),
		Outside: false,
	}
	outsideWide := cf.Range{
		Lo:      finiteClosedBound(10, 1),
		Hi:      finiteClosedBound(20, 1),
		Outside: true,
	}

	if got := insideNarrow.Cmp(outsideWide); got >= 0 {
		t.Fatalf("insideNarrow.Cmp(outsideWide) = %d, want < 0", got)
	}
	if got := outsideWide.Cmp(insideNarrow); got <= 0 {
		t.Fatalf("outsideWide.Cmp(insideNarrow) = %d, want > 0", got)
	}
}

func negInfOpenBound() cf.Bound {
	return cf.Bound{
		Kind:   cf.BoundNegInf,
		Closed: false,
	}
}

func posInfOpenBound() cf.Bound {
	return cf.Bound{
		Kind:   cf.BoundPosInf,
		Closed: false,
	}
}

// cf/range_bb_test.go v2
