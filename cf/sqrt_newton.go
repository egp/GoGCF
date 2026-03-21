// cf/sqrt_newton.go v2
package cf

import "math/big"

func newtonUpperSqrtBound(x Rational, upper Rational) (Rational, error) {
	xr, err := normalizedRational(x)
	if err != nil {
		return Rational{}, err
	}
	ur, err := normalizedRational(upper)
	if err != nil {
		return Rational{}, err
	}

	if xr.Num.Sign() < 0 {
		return Rational{}, ErrUndefined
	}
	if ur.Num.Sign() <= 0 {
		return Rational{}, ErrUndefined
	}

	// q = x / upper
	qNum := new(big.Int).Mul(new(big.Int).Set(xr.Num), new(big.Int).Set(ur.Den))
	qDen := new(big.Int).Mul(new(big.Int).Set(xr.Den), new(big.Int).Set(ur.Num))

	// sum = upper + q
	sumNum := new(big.Int).Mul(new(big.Int).Set(ur.Num), new(big.Int).Set(qDen))
	sumNum.Add(sumNum, new(big.Int).Mul(new(big.Int).Set(qNum), new(big.Int).Set(ur.Den)))
	sumDen := new(big.Int).Mul(new(big.Int).Set(ur.Den), new(big.Int).Set(qDen))

	// next = sum / 2
	return normalizedRational(Rational{
		Num: sumNum,
		Den: new(big.Int).Mul(sumDen, big.NewInt(2)),
	})
}

func sqrtEnclosureFromUpperBound(x Rational, upper Rational) (Range, error) {
	xr, err := normalizedRational(x)
	if err != nil {
		return Range{}, err
	}
	ur, err := normalizedRational(upper)
	if err != nil {
		return Range{}, err
	}

	if xr.Num.Sign() < 0 {
		return Range{}, ErrUndefined
	}
	if ur.Num.Sign() <= 0 {
		return Range{}, ErrUndefined
	}

	// lower = x / upper
	lo, err := normalizedRational(Rational{
		Num: new(big.Int).Mul(new(big.Int).Set(xr.Num), new(big.Int).Set(ur.Den)),
		Den: new(big.Int).Mul(new(big.Int).Set(xr.Den), new(big.Int).Set(ur.Num)),
	})
	if err != nil {
		return Range{}, err
	}

	return Range{
		Lo: Bound{
			Kind:   BoundFinite,
			Value:  lo,
			Closed: false,
		},
		Hi: Bound{
			Kind:   BoundFinite,
			Value:  ur,
			Closed: false,
		},
		Outside: false,
	}, nil
}

// cf/sqrt_newton.go v2
