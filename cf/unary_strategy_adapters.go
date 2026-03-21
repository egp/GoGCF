// cf/unary_strategy_adapters.go v2
package cf

import "math/big"

type unaryLFTStrategy struct {
	op unaryLFT
}

func (s unaryLFTStrategy) RangeFromOperand(xr Range) (Range, error) {
	return s.op.rangeFromXRange(xr)
}

func (s unaryLFTStrategy) EmitCandidateFromOperand(xr Range) (*big.Int, bool, error) {
	return s.op.emitCandidateFromRange(xr)
}

func (s unaryLFTStrategy) Emit(term *big.Int) unaryStrategy {
	return unaryLFTStrategy{op: s.op.emit(term)}
}

func (s unaryLFTStrategy) ExactEval(x Rational) (Rational, bool, error) {
	r, err := s.op.evalAt(x)
	if err != nil {
		return Rational{}, false, err
	}
	return r, true, nil
}

type diagLFTStrategy struct {
	op diagLFT
}

func (s diagLFTStrategy) RangeFromOperand(xr Range) (Range, error) {
	return s.op.rangeFromXRange(xr)
}

func (s diagLFTStrategy) EmitCandidateFromOperand(xr Range) (*big.Int, bool, error) {
	return s.op.emitCandidateFromRange(xr)
}

func (s diagLFTStrategy) Emit(term *big.Int) unaryStrategy {
	return diagLFTStrategy{op: s.op.emit(term)}
}

func (s diagLFTStrategy) ExactEval(x Rational) (Rational, bool, error) {
	r, err := s.op.evalAt(x)
	if err != nil {
		return Rational{}, false, err
	}
	return r, true, nil
}

// cf/unary_strategy_adapters.go v2
