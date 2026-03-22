// cf/unary_child_ingest.go v1
package cf

import "math/big"

func (u unaryLFT) ingest(term *big.Int) unaryLFT {
	if term == nil {
		return u
	}

	return unaryLFT{
		a: addMul(u.a, term, u.b),
		b: new(big.Int).Set(u.a),
		c: addMul(u.c, term, u.d),
		d: new(big.Int).Set(u.c),
	}
}

func (d diagLFT) ingest(term *big.Int) diagLFT {
	if term == nil {
		return d
	}

	q2 := new(big.Int).Mul(new(big.Int).Set(term), new(big.Int).Set(term))
	twoQ := new(big.Int).Mul(big.NewInt(2), new(big.Int).Set(term))

	return diagLFT{
		a: addMul(addMul(d.a, q2, new(big.Int).Mul(new(big.Int).Set(d.b), new(big.Int).Set(term))), big.NewInt(1), d.c),
		b: addMul(d.a, twoQ, d.b),
		c: new(big.Int).Set(d.a),
		d: addMul(addMul(d.d, q2, new(big.Int).Mul(new(big.Int).Set(d.e), new(big.Int).Set(term))), big.NewInt(1), d.f),
		e: addMul(d.d, twoQ, d.e),
		f: new(big.Int).Set(d.d),
	}
}

// cf/unary_child_ingest.go v1
