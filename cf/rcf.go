// cf/rcf.go v2
package cf

import "math/big"

func prefixGCFfromRational(r Rational) (GCF, error) {
	terms, err := rcfTermsFromRational(r.Num, r.Den)
	if err != nil {
		return nil, err
	}

	return rcfPrefixGCF{
		terms: terms,
		index: 0,
	}, nil
}

func rcfTermsFromRational(num, den *big.Int) ([]*big.Int, error) {
	if num == nil || den == nil || den.Sign() == 0 {
		return nil, ErrUndefined
	}

	n := new(big.Int).Set(num)
	d := new(big.Int).Set(den)
	if d.Sign() < 0 {
		n.Neg(n)
		d.Neg(d)
	}

	terms := make([]*big.Int, 0, 8)
	zero := big.NewInt(0)

	for d.Sign() != 0 {
		q, r := floorDivModBig(n, d)
		terms = append(terms, q)

		if r.Cmp(zero) == 0 {
			break
		}

		n, d = d, r
	}

	return terms, nil
}

func floorDivModBig(n, d *big.Int) (*big.Int, *big.Int) {
	q := new(big.Int)
	r := new(big.Int)
	q.QuoRem(n, d, r)

	if r.Sign() != 0 && r.Sign() != d.Sign() {
		q.Sub(q, big.NewInt(1))
		r.Add(r, d)
	}

	return q, r
}

// cf/rcf.go v2
