// cf/unary_strategy_ingest.go v1
package cf

import "math/big"

func (s unaryLFTStrategy) IngestChild(term *big.Int) (unaryStrategy, error) {
	if term == nil {
		return nil, ErrUndefined
	}
	return unaryLFTStrategy{op: s.op.ingest(term)}, nil
}

func (s diagLFTStrategy) IngestChild(term *big.Int) (unaryStrategy, error) {
	if term == nil {
		return nil, ErrUndefined
	}
	return diagLFTStrategy{op: s.op.ingest(term)}, nil
}

// cf/unary_strategy_ingest.go v1
