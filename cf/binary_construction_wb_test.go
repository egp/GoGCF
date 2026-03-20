// cf/binary_construction_wb_test.go v1
package cf

import (
	"math/big"
	"testing"
)

func TestWB_PQFromRCFTerm_ValueSevenMapsToPQ_1_7(t *testing.T) {
	pq, err := pqFromRCFTerm(RCFTerm{
		Kind:  RCFValue,
		Value: big.NewInt(7),
	})
	if err != nil {
		t.Fatalf("pqFromRCFTerm() error: %v", err)
	}
	if !pq.IsValue() {
		t.Fatalf("pqFromRCFTerm(value) should produce PQValue")
	}
	assertBigIntString(t, "pq.P", pq.P, "1")
	assertBigIntString(t, "pq.Q", pq.Q, "7")
}

func TestWB_PQFromRCFTerm_EOFMapsToPQEOF(t *testing.T) {
	pq, err := pqFromRCFTerm(RCFTerm{Kind: RCFEOF})
	if err != nil {
		t.Fatalf("pqFromRCFTerm() error: %v", err)
	}
	if !pq.IsEOF() {
		t.Fatalf("pqFromRCFTerm(EOF) should produce PQEOF")
	}
}

func TestWB_Add_ConstructsBinaryGCF_WithAdditionLFT(t *testing.T) {
	g := Add(FromSource(Int64(1)), FromSource(Int64(2)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Add() concrete type = %T, want binaryGCF", g)
	}

	assertBigIntString(t, "op.a", bg.op.a, "0")
	assertBigIntString(t, "op.b", bg.op.b, "1")
	assertBigIntString(t, "op.c", bg.op.c, "1")
	assertBigIntString(t, "op.d", bg.op.d, "0")
	assertBigIntString(t, "op.e", bg.op.e, "0")
	assertBigIntString(t, "op.f", bg.op.f, "0")
	assertBigIntString(t, "op.g", bg.op.g, "0")
	assertBigIntString(t, "op.h", bg.op.h, "1")
}

func TestWB_Sub_ConstructsBinaryGCF_WithSubtractionLFT(t *testing.T) {
	g := Sub(FromSource(Int64(1)), FromSource(Int64(2)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Sub() concrete type = %T, want binaryGCF", g)
	}

	assertBigIntString(t, "op.a", bg.op.a, "0")
	assertBigIntString(t, "op.b", bg.op.b, "1")
	assertBigIntString(t, "op.c", bg.op.c, "-1")
	assertBigIntString(t, "op.d", bg.op.d, "0")
	assertBigIntString(t, "op.e", bg.op.e, "0")
	assertBigIntString(t, "op.f", bg.op.f, "0")
	assertBigIntString(t, "op.g", bg.op.g, "0")
	assertBigIntString(t, "op.h", bg.op.h, "1")
}

func TestWB_Mul_ConstructsBinaryGCF_WithMultiplicationLFT(t *testing.T) {
	g := Mul(FromSource(Int64(1)), FromSource(Int64(2)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Mul() concrete type = %T, want binaryGCF", g)
	}

	assertBigIntString(t, "op.a", bg.op.a, "1")
	assertBigIntString(t, "op.b", bg.op.b, "0")
	assertBigIntString(t, "op.c", bg.op.c, "0")
	assertBigIntString(t, "op.d", bg.op.d, "0")
	assertBigIntString(t, "op.e", bg.op.e, "0")
	assertBigIntString(t, "op.f", bg.op.f, "0")
	assertBigIntString(t, "op.g", bg.op.g, "0")
	assertBigIntString(t, "op.h", bg.op.h, "1")
}

func TestWB_Div_ConstructsBinaryGCF_WithDivisionLFT(t *testing.T) {
	g := Div(FromSource(Int64(1)), FromSource(Int64(2)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Div() concrete type = %T, want binaryGCF", g)
	}

	assertBigIntString(t, "op.a", bg.op.a, "0")
	assertBigIntString(t, "op.b", bg.op.b, "1")
	assertBigIntString(t, "op.c", bg.op.c, "0")
	assertBigIntString(t, "op.d", bg.op.d, "0")
	assertBigIntString(t, "op.e", bg.op.e, "0")
	assertBigIntString(t, "op.f", bg.op.f, "0")
	assertBigIntString(t, "op.g", bg.op.g, "1")
	assertBigIntString(t, "op.h", bg.op.h, "0")
}

// cf/binary_construction_wb_test.go v1
