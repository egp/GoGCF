// cf/binary_complete_wb_test.go v1
package cf

import "testing"

func TestWB_BinaryGCF_CompleteToRational_AddThreeAndFive_GivesEightOverOne(t *testing.T) {
	g := Add(FromSource(Int64(3)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Add() concrete type = %T, want binaryGCF", g)
	}

	r, err := bg.completeToRational()
	if err != nil {
		t.Fatalf("completeToRational() error: %v", err)
	}

	assertBigIntString(t, "r.Num", r.Num, "8")
	assertBigIntString(t, "r.Den", r.Den, "1")
}

func TestWB_BinaryGCF_CompleteToRational_DivTwelveByFive_GivesTwelveOverFive(t *testing.T) {
	g := Div(FromSource(Int64(12)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Div() concrete type = %T, want binaryGCF", g)
	}

	r, err := bg.completeToRational()
	if err != nil {
		t.Fatalf("completeToRational() error: %v", err)
	}

	assertBigIntString(t, "r.Num", r.Num, "12")
	assertBigIntString(t, "r.Den", r.Den, "5")
}

func TestWB_BinaryGCF_CompleteToRational_SubThreeMinusEight_GivesMinusFiveOverOne(t *testing.T) {
	g := Sub(FromSource(Int64(3)), FromSource(Int64(8)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Sub() concrete type = %T, want binaryGCF", g)
	}

	r, err := bg.completeToRational()
	if err != nil {
		t.Fatalf("completeToRational() error: %v", err)
	}

	assertBigIntString(t, "r.Num", r.Num, "-5")
	assertBigIntString(t, "r.Den", r.Den, "1")
}

func TestWB_BinaryGCF_CompleteToRational_DivFiveByTwelve_GivesFiveOverTwelve(t *testing.T) {
	g := Div(FromSource(Int64(5)), FromSource(Int64(12)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Div() concrete type = %T, want binaryGCF", g)
	}

	r, err := bg.completeToRational()
	if err != nil {
		t.Fatalf("completeToRational() error: %v", err)
	}

	assertBigIntString(t, "r.Num", r.Num, "5")
	assertBigIntString(t, "r.Den", r.Den, "12")
}

// cf/binary_complete_wb_test.go v1
