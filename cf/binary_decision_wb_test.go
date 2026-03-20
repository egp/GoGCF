// cf/binary_decision_wb_test.go v1
package cf

import (
	"math/big"
	"testing"
)

type scriptedRangeGCF struct {
	term RCFTerm
	rest GCF
	r    Range
}

func (g scriptedRangeGCF) Next() (RCFTerm, GCF, error) {
	if g.term.IsEOF() {
		return g.term, g, nil
	}
	if g.rest == nil {
		return g.term, scriptedRangeGCF{term: RCFTerm{Kind: RCFEOF}, r: Range{}}, nil
	}
	return g.term, g.rest, nil
}

func (g scriptedRangeGCF) Range() Range {
	return g.r
}

func (g scriptedRangeGCF) Take(n int) (GCF, error) {
	return nil, ErrUndefined
}

func (g scriptedRangeGCF) Convergent() (Rational, error) {
	return Rational{}, ErrUndefined
}

func TestWB_BinaryGCF_Step_ChoosesIngestXWhenPredictedOutputRangeIsBetter(t *testing.T) {
	xNow := finiteInsideRange(
		Rational{Num: big.NewInt(1), Den: big.NewInt(1)},
		Rational{Num: big.NewInt(3), Den: big.NewInt(1)},
		true,
	)
	yNow := finiteInsideRange(
		Rational{Num: big.NewInt(10), Den: big.NewInt(1)},
		Rational{Num: big.NewInt(20), Den: big.NewInt(1)},
		true,
	)

	x := scriptedRangeGCF{
		term: RCFTerm{Kind: RCFValue, Value: big.NewInt(1)},
		rest: FromSource(Rat64(2, 1)),
		r:    xNow,
	}
	y := scriptedRangeGCF{
		term: RCFTerm{Kind: RCFValue, Value: big.NewInt(7)},
		rest: FromSource(Rat64(99, 1)),
		r:    yNow,
	}

	bg := binaryGCF{
		op: binaryLFT{
			a: big.NewInt(0),
			b: big.NewInt(1),
			c: big.NewInt(0),
			d: big.NewInt(0),
			e: big.NewInt(0),
			f: big.NewInt(0),
			g: big.NewInt(0),
			h: big.NewInt(1),
		},
		left:  x,
		right: y,
	}

	action, next, err := bg.step()
	if err != nil {
		t.Fatalf("step() error: %v", err)
	}
	if action != decisionIngestLeft {
		t.Fatalf("step() action = %v, want %v", action, decisionIngestLeft)
	}

	r := next.Range()
	if !r.IsExact() {
		t.Fatalf("next.Range() should be exact after ingest-x")
	}
	assertBigIntString(t, "r.Lo.Value.Num", r.Lo.Value.Num, "3")
	assertBigIntString(t, "r.Lo.Value.Den", r.Lo.Value.Den, "2")
}

func TestWB_BinaryGCF_Step_ChoosesIngestYWhenPredictedOutputRangeIsBetter(t *testing.T) {
	xNow := finiteInsideRange(
		Rational{Num: big.NewInt(10), Den: big.NewInt(1)},
		Rational{Num: big.NewInt(20), Den: big.NewInt(1)},
		true,
	)
	yNow := finiteInsideRange(
		Rational{Num: big.NewInt(1), Den: big.NewInt(1)},
		Rational{Num: big.NewInt(3), Den: big.NewInt(1)},
		true,
	)

	x := scriptedRangeGCF{
		term: RCFTerm{Kind: RCFValue, Value: big.NewInt(7)},
		rest: FromSource(Rat64(99, 1)),
		r:    xNow,
	}
	y := scriptedRangeGCF{
		term: RCFTerm{Kind: RCFValue, Value: big.NewInt(1)},
		rest: FromSource(Rat64(2, 1)),
		r:    yNow,
	}

	bg := binaryGCF{
		op: binaryLFT{
			a: big.NewInt(0),
			b: big.NewInt(0),
			c: big.NewInt(1),
			d: big.NewInt(0),
			e: big.NewInt(0),
			f: big.NewInt(0),
			g: big.NewInt(0),
			h: big.NewInt(1),
		},
		left:  x,
		right: y,
	}

	action, next, err := bg.step()
	if err != nil {
		t.Fatalf("step() error: %v", err)
	}
	if action != decisionIngestRight {
		t.Fatalf("step() action = %v, want %v", action, decisionIngestRight)
	}

	r := next.Range()
	if !r.IsExact() {
		t.Fatalf("next.Range() should be exact after ingest-y")
	}
	assertBigIntString(t, "r.Lo.Value.Num", r.Lo.Value.Num, "3")
	assertBigIntString(t, "r.Lo.Value.Den", r.Lo.Value.Den, "2")
}
