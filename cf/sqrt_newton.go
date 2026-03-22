// cf/sqrt_newton.go v4
package cf

import "math/big"

type sqrtOuroborosState struct {
	Operand    Rational
	Post       unaryLFT
	UpperBound Rational
}

type sqrtOuroborosResult struct {
	ImageRange Range
	Candidate  *big.Int
	UpperBound Rational
}

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

	qNum := new(big.Int).Mul(new(big.Int).Set(xr.Num), new(big.Int).Set(ur.Den))
	qDen := new(big.Int).Mul(new(big.Int).Set(xr.Den), new(big.Int).Set(ur.Num))

	sumNum := new(big.Int).Mul(new(big.Int).Set(ur.Num), new(big.Int).Set(qDen))
	sumNum.Add(sumNum, new(big.Int).Mul(new(big.Int).Set(qNum), new(big.Int).Set(ur.Den)))
	sumDen := new(big.Int).Mul(new(big.Int).Set(ur.Den), new(big.Int).Set(qDen))

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

func runSqrtOuroboros(state sqrtOuroborosState) (sqrtOuroborosResult, bool, error) {
	xr, err := normalizedRational(state.Operand)
	if err != nil {
		return sqrtOuroborosResult{}, false, err
	}
	if xr.Num.Sign() < 0 {
		return sqrtOuroborosResult{}, false, ErrUndefined
	}

	var upper Rational
	seeded := false
	if rationalWellFormed(state.UpperBound) {
		upper, err = normalizedRational(state.UpperBound)
		if err != nil {
			return sqrtOuroborosResult{}, false, err
		}
		seeded = true
	} else {
		upper, err = normalizedRational(Rational{
			Num: ceilBigIntSqrt(xr.Num),
			Den: floorBigIntSqrt(xr.Den),
		})
		if err != nil {
			return sqrtOuroborosResult{}, false, err
		}
	}

	for i := 0; i < 64; i++ {
		if seeded {
			nextUpper, err := newtonUpperSqrtBound(xr, upper)
			if err != nil {
				return sqrtOuroborosResult{}, false, err
			}
			upper = nextUpper
			seeded = false
		}

		refined, err := sqrtEnclosureFromUpperBound(xr, upper)
		if err != nil {
			return sqrtOuroborosResult{}, false, err
		}

		q, ok, err := state.Post.emitCandidateFromRange(refined)
		if err != nil {
			return sqrtOuroborosResult{}, false, err
		}
		if ok {
			zr, err := state.Post.rangeFromXRange(refined)
			if err != nil {
				return sqrtOuroborosResult{}, false, err
			}
			if rangeWellFormed(zr) {
				return sqrtOuroborosResult{
					ImageRange: zr,
					Candidate:  q,
					UpperBound: upper,
				}, true, nil
			}
		}

		nextUpper, err := newtonUpperSqrtBound(xr, upper)
		if err != nil {
			return sqrtOuroborosResult{}, false, err
		}
		if rationalCmp(nextUpper, upper) == 0 {
			break
		}
		upper = nextUpper
	}

	return sqrtOuroborosResult{}, false, nil
}

func exactPositiveSqrtNewtonFeedback(x Rational, post unaryLFT, seedUpper Rational) (sqrtOuroborosResult, bool, error) {
	return runSqrtOuroboros(sqrtOuroborosState{
		Operand:    x,
		Post:       post,
		UpperBound: seedUpper,
	})
}

// cf/sqrt_newton.go v4
