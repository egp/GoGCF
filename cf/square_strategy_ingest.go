// cf/square_strategy_ingest.go v1
package cf

import "math/big"

func (s squareStrategy) IngestChild(term *big.Int) (unaryStrategy, error) {
	if term == nil {
		return nil, ErrUndefined
	}
	return diagLFTStrategy{op: s.effectiveOp().ingest(term)}, nil
}

// cf/square_strategy_ingest.go v1
