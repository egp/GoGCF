// sqrt2_source.go v2
package cf

import "math/big"

type sqrt2Source struct {
	index int
}

func Sqrt2() GCFSource {
	return sqrt2Source{index: 0}
}

func (s sqrt2Source) NextPQ() (PQTerm, GCFSource, error) {
	q := int64(2)
	if s.index == 0 {
		q = 1
	}

	term := PQTerm{
		Kind: PQValue,
		P:    big.NewInt(1),
		Q:    big.NewInt(q),
	}

	return term, sqrt2Source{index: s.index + 1}, nil
}

// sqrt2_source.go v2
