// cf/rcf_prefix_helpers.go v1
package cf

import "math/big"

func convergentFromTerms(terms []*big.Int, index int) (Rational, error) {
	if index < 0 || index >= len(terms) {
		return Rational{}, ErrUndefined
	}

	num := new(big.Int).Set(terms[len(terms)-1])
	den := big.NewInt(1)

	for i := len(terms) - 2; i >= index; i-- {
		nextNum := new(big.Int).Mul(terms[i], num)
		nextNum.Add(nextNum, den)
		num, den = nextNum, num
	}

	return Rational{
		Num: num,
		Den: den,
	}, nil
}

func rcfTermsFromRat64(num, den int64) []*big.Int {
	if den == 0 {
		return nil
	}

	n := num
	d := den
	if d < 0 {
		n = -n
		d = -d
	}

	terms := make([]*big.Int, 0, 8)
	for d != 0 {
		q := floorDivInt64(n, d)
		terms = append(terms, big.NewInt(q))

		r := n - q*d
		if r == 0 {
			break
		}

		n, d = d, r
	}

	return terms
}

func floorDivInt64(n, d int64) int64 {
	q := n / d
	r := n % d
	if r != 0 && ((r > 0) != (d > 0)) {
		q--
	}
	return q
}

// EOF cf/rcf_prefix_helpers.go v1
