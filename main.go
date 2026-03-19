// main.go v1
package main

import (
	"fmt"
	"log"
	"math/big"

	"example.com/yourmodule/cf"
)

func main() {
	targetSrc := cf.Div(
		cf.Sqrt(
			cf.Add(
				cf.Div(
					cf.Int64(3),
					cf.Mul(cf.Pi(), cf.Pi()),
				),
				cf.E(),
			),
		),
		cf.Sub(
			cf.Tanh(
				cf.Sqrt(cf.Int64(5)),
			),
			cf.Sin(
				cf.Mul(
					cf.Rat64(69, 180),
					cf.Pi(),
				),
			),
		),
	)

	target := cf.New(targetSrc)

	fmt.Println("target formula:")
	fmt.Println("sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))")
	fmt.Println()

	fmt.Println("first 12 emitted RCF terms:")
	for i := 0; i < 12; i++ {
		term, err := target.Next()
		if err != nil {
			log.Fatalf("Next term %d failed: %v", i, err)
		}
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(term.String())
	}
	fmt.Println()
	fmt.Println()

	rng := target.Range()
	fmt.Printf("current range: [%s, %s]\n", rationalString(rng.Lo), rationalString(rng.Hi))
	fmt.Printf("inside=%v outside=%v exact=%v\n", rng.IsInside(), rng.IsOutside(), rng.IsExact())
	fmt.Println()

	conv := target.Convergent()
	fmt.Printf("current convergent: %s\n", rationalString(conv))
}

func rationalString(r cf.Rational) string {
	if r.Den == nil || r.Den.Sign() == 0 {
		return "<invalid>"
	}
	if r.Den.Cmp(big.NewInt(1)) == 0 {
		return r.Num.String()
	}
	return r.Num.String() + "/" + r.Den.String()
}

// main.go v1
