// cf/unary_strategy_wb_test.go v1
package cf

import (
	"math/big"
	"testing"
)

type scriptedUnaryStrategy struct {
	r       Range
	q       *big.Int
	canEmit bool
	next    unaryStrategy
	exact   Rational
	exactOK bool
}

func (s scriptedUnaryStrategy) RangeFromOperand(xr Range) (Range, error) {
	return s.r, nil
}

func (s scriptedUnaryStrategy) EmitCandidateFromOperand(xr Range) (*big.Int, bool, error) {
	if !s.canEmit || s.q == nil {
		return nil, false, nil
	}
	return new(big.Int).Set(s.q), true, nil
}

func (s scriptedUnaryStrategy) Emit(term *big.Int) unaryStrategy {
	if s.next != nil {
		return s.next
	}
	return scriptedUnaryStrategy{}
}

func (s scriptedUnaryStrategy) ExactEval(x Rational) (Rational, bool, error) {
	if !s.exactOK {
		return Rational{}, false, nil
	}
	return s.exact, true, nil
}

func TestWB_StrategyUnaryGCF_Range_DelegatesToStrategy(t *testing.T) {
	g := strategyUnaryGCF{
		strategy: scriptedUnaryStrategy{
			r: exactRangeFromRational(Rational{
				Num: big.NewInt(5),
				Den: big.NewInt(2),
			}),
		},
		child: FromSource(Rat64(99, 1)),
	}

	r := g.Range()
	if !r.IsExact() {
		t.Fatalf("Range() should be exact")
	}

	assertBigIntString(t, "r.Lo.Value.Num", r.Lo.Value.Num, "5")
	assertBigIntString(t, "r.Lo.Value.Den", r.Lo.Value.Den, "2")
	assertBigIntString(t, "r.Hi.Value.Num", r.Hi.Value.Num, "5")
	assertBigIntString(t, "r.Hi.Value.Den", r.Hi.Value.Den, "2")
}

func TestWB_StrategyUnaryGCF_Next_EmitsStrategyCandidateAndLeavesStepwiseRemainder(t *testing.T) {
	nextStrategy := scriptedUnaryStrategy{
		r: exactRangeFromRational(Rational{
			Num: big.NewInt(5),
			Den: big.NewInt(2),
		}),
	}

	g := strategyUnaryGCF{
		strategy: scriptedUnaryStrategy{
			q:       big.NewInt(2),
			canEmit: true,
			next:    nextStrategy,
		},
		child: FromSource(Rat64(12, 5)),
	}

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
	if got, want := value.String(), "2"; got != want {
		t.Fatalf("first term = %s, want %s", got, want)
	}

	nextSG, ok := rest.(strategyUnaryGCF)
	if !ok {
		t.Fatalf("remainder concrete type = %T, want strategyUnaryGCF", rest)
	}
	if nextSG.resolved != nil {
		t.Fatalf("remainder should stay stepwise, not resolved")
	}

	r := nextSG.Range()
	if !r.IsExact() {
		t.Fatalf("remainder Range() should be exact")
	}

	assertBigIntString(t, "r.Lo.Value.Num", r.Lo.Value.Num, "5")
	assertBigIntString(t, "r.Lo.Value.Den", r.Lo.Value.Den, "2")
	assertBigIntString(t, "r.Hi.Value.Num", r.Hi.Value.Num, "5")
	assertBigIntString(t, "r.Hi.Value.Den", r.Hi.Value.Den, "2")
}

func TestWB_StrategyUnaryGCF_Next_DoesNotMutateOriginalNode(t *testing.T) {
	nextStrategy := scriptedUnaryStrategy{
		r: exactRangeFromRational(Rational{
			Num: big.NewInt(5),
			Den: big.NewInt(2),
		}),
	}

	g := strategyUnaryGCF{
		strategy: scriptedUnaryStrategy{
			q:       big.NewInt(2),
			canEmit: true,
			next:    nextStrategy,
		},
		child: FromSource(Rat64(12, 5)),
	}

	_, _, err := g.Next()
	if err != nil {
		t.Fatalf("Next() error: %v", err)
	}

	if g.resolved != nil {
		t.Fatalf("original strategy unary node should remain unresolved")
	}

	r := g.child.Range()
	if !r.IsExact() {
		t.Fatalf("original child should remain unchanged")
	}
	assertBigIntString(t, "r.Lo.Value.Num", r.Lo.Value.Num, "12")
	assertBigIntString(t, "r.Lo.Value.Den", r.Lo.Value.Den, "5")
}

// cf/unary_strategy_wb_test.go v1
