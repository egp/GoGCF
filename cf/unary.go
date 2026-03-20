// cf/unary.go v2
package cf

import (
	"io"
	"math/big"
)

type unaryGCF struct {
	op       unaryLFT
	child    GCF
	resolved GCF
}

func (u unaryLFT) evalAt(x Rational) (Rational, error) {
	if u.a == nil || u.b == nil || u.c == nil || u.d == nil {
		return Rational{}, ErrUndefined
	}

	xr, err := normalizedRational(x)
	if err != nil {
		return Rational{}, err
	}

	xn := xr.Num
	xd := xr.Den

	num := new(big.Int).Mul(new(big.Int).Set(u.a), new(big.Int).Set(xn))
	num.Add(num, new(big.Int).Mul(new(big.Int).Set(u.b), new(big.Int).Set(xd)))

	den := new(big.Int).Mul(new(big.Int).Set(u.c), new(big.Int).Set(xn))
	den.Add(den, new(big.Int).Mul(new(big.Int).Set(u.d), new(big.Int).Set(xd)))

	if den.Sign() == 0 {
		return Rational{}, ErrUndefined
	}

	return normalizedRational(Rational{
		Num: num,
		Den: den,
	})
}

func (u unaryLFT) emit(term *big.Int) unaryLFT {
	if term == nil {
		return u
	}

	return unaryLFT{
		a: new(big.Int).Set(u.c),
		b: new(big.Int).Set(u.d),
		c: new(big.Int).Sub(new(big.Int).Set(u.a), new(big.Int).Mul(new(big.Int).Set(term), u.c)),
		d: new(big.Int).Sub(new(big.Int).Set(u.b), new(big.Int).Mul(new(big.Int).Set(term), u.d)),
	}
}

func (u unaryLFT) rangeFromXRange(xr Range) (Range, error) {
	x, exact := exactFiniteRangeValue(xr)
	if exact {
		z, err := u.evalAt(x)
		if err != nil {
			return Range{}, err
		}
		return exactRangeFromRational(z), nil
	}

	xLo, xHi, ok := finiteInsideRangeEndpoints(xr)
	if !ok {
		return Range{}, nil
	}

	values := make([]Rational, 0, 2)
	for _, x := range []Rational{xLo, xHi} {
		v, err := u.evalAt(x)
		if err != nil {
			continue
		}
		values = append(values, v)
	}

	if len(values) == 0 {
		return Range{}, nil
	}

	lo := values[0]
	hi := values[0]
	for _, v := range values[1:] {
		if rationalCmp(v, lo) < 0 {
			lo = v
		}
		if rationalCmp(v, hi) > 0 {
			hi = v
		}
	}

	allClosed := xr.Lo.Closed && xr.Hi.Closed
	return finiteInsideRange(lo, hi, allClosed), nil
}

func (u unaryLFT) emitCandidateFromRange(xr Range) (*big.Int, bool, error) {
	zr, err := u.rangeFromXRange(xr)
	if err != nil {
		return nil, false, err
	}
	if !rangeWellFormed(zr) {
		return nil, false, nil
	}
	if zr.Outside {
		return nil, false, nil
	}
	if zr.Lo.Kind != BoundFinite || zr.Hi.Kind != BoundFinite {
		return nil, false, nil
	}
	if !rationalWellFormed(zr.Lo.Value) || !rationalWellFormed(zr.Hi.Value) {
		return nil, false, nil
	}

	lo, err := normalizedRational(zr.Lo.Value)
	if err != nil {
		return nil, false, err
	}
	hi, err := normalizedRational(zr.Hi.Value)
	if err != nil {
		return nil, false, err
	}

	qLo, _ := floorDivModBig(lo.Num, lo.Den)
	qHi, hiRem := floorDivModBig(hi.Num, hi.Den)
	if !zr.Hi.Closed && hiRem.Sign() == 0 {
		qHi.Sub(qHi, big.NewInt(1))
	}

	if qLo.Cmp(qHi) != 0 {
		return nil, false, nil
	}

	return qLo, true, nil
}

func (g unaryGCF) Next() (RCFTerm, GCF, error) {
	if g.resolved != nil {
		term, rest, err := g.resolved.Next()
		if err != nil {
			return RCFTerm{}, g, err
		}
		return term, unaryGCF{
			op:       g.op,
			child:    g.child,
			resolved: rest,
		}, nil
	}

	xr := g.child.Range()
	q, ok, err := g.op.emitCandidateFromRange(xr)
	if err != nil {
		return RCFTerm{}, g, err
	}
	if !ok {
		return RCFTerm{}, g, ErrUndefined
	}

	term := RCFTerm{
		Kind:  RCFValue,
		Value: new(big.Int).Set(q),
	}

	if x, exact := exactFiniteRangeValue(xr); exact {
		z, err := g.op.evalAt(x)
		if err != nil {
			return RCFTerm{}, g, err
		}
		_, rem := floorDivModBig(z.Num, z.Den)
		if rem.Sign() == 0 {
			return term, unaryGCF{
				op:       g.op,
				child:    g.child,
				resolved: rcfPrefixGCF{terms: nil, index: 0},
			}, nil
		}
	}

	return term, unaryGCF{
		op:    g.op.emit(q),
		child: g.child,
	}, nil
}

func (g unaryGCF) Range() Range {
	if g.resolved != nil {
		return g.resolved.Range()
	}

	zr, err := g.op.rangeFromXRange(g.child.Range())
	if err != nil {
		return Range{}
	}
	return zr
}

func (g unaryGCF) Take(n int) (GCF, error) {
	if n < 0 {
		return nil, ErrUndefined
	}

	terms := make([]*big.Int, 0, n)
	cur := GCF(g)

	for len(terms) < n {
		term, rest, err := cur.Next()
		if err != nil {
			return nil, err
		}

		if term.IsEOF() {
			return rcfPrefixGCF{terms: terms, index: 0}, io.EOF
		}

		value, ok := term.BigInt()
		if !ok {
			return nil, ErrUndefined
		}

		terms = append(terms, new(big.Int).Set(value))
		cur = rest
	}

	return rcfPrefixGCF{terms: terms, index: 0}, nil
}

func (g unaryGCF) Convergent() (Rational, error) {
	return Rational{}, ErrUndefined
}

// cf/unary.go v2
