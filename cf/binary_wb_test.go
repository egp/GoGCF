// cf/binary_wb_test.go v1
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

func TestWB_BinaryGCF_IngestLeftStep_UsesChildEvaluatorOutput(t *testing.T) {
	g := Add(FromSource(Int64(3)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Add() concrete type = %T, want binaryGCF", g)
	}

	next, err := bg.ingestLeftStep()
	if err != nil {
		t.Fatalf("ingestLeftStep() error: %v", err)
	}

	assertBigIntString(t, "op.a", next.op.a, "1")
	assertBigIntString(t, "op.b", next.op.b, "3")
	assertBigIntString(t, "op.c", next.op.c, "0")
	assertBigIntString(t, "op.d", next.op.d, "1")
	assertBigIntString(t, "op.e", next.op.e, "0")
	assertBigIntString(t, "op.f", next.op.f, "1")
	assertBigIntString(t, "op.g", next.op.g, "0")
	assertBigIntString(t, "op.h", next.op.h, "0")

	leftTerm, _, err := next.left.Next()
	if err != nil {
		t.Fatalf("next.left.Next() error: %v", err)
	}
	if !leftTerm.IsEOF() {
		t.Fatalf("left child should have advanced to EOF after one ingest")
	}
}

func TestWB_BinaryGCF_IngestRightStep_UsesChildEvaluatorOutput(t *testing.T) {
	g := Add(FromSource(Int64(3)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Add() concrete type = %T, want binaryGCF", g)
	}

	next, err := bg.ingestRightStep()
	if err != nil {
		t.Fatalf("ingestRightStep() error: %v", err)
	}

	assertBigIntString(t, "op.a", next.op.a, "1")
	assertBigIntString(t, "op.b", next.op.b, "0")
	assertBigIntString(t, "op.c", next.op.c, "5")
	assertBigIntString(t, "op.d", next.op.d, "1")
	assertBigIntString(t, "op.e", next.op.e, "0")
	assertBigIntString(t, "op.f", next.op.f, "0")
	assertBigIntString(t, "op.g", next.op.g, "1")
	assertBigIntString(t, "op.h", next.op.h, "0")

	rightTerm, _, err := next.right.Next()
	if err != nil {
		t.Fatalf("next.right.Next() error: %v", err)
	}
	if !rightTerm.IsEOF() {
		t.Fatalf("right child should have advanced to EOF after one ingest")
	}
}

func TestWB_BinaryGCF_IngestLeftStep_DoesNotMutateOriginalNode(t *testing.T) {
	g := Add(FromSource(Int64(3)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Add() concrete type = %T, want binaryGCF", g)
	}

	_, err := bg.ingestLeftStep()
	if err != nil {
		t.Fatalf("ingestLeftStep() error: %v", err)
	}

	assertBigIntString(t, "op.a", bg.op.a, "0")
	assertBigIntString(t, "op.b", bg.op.b, "1")
	assertBigIntString(t, "op.c", bg.op.c, "1")
	assertBigIntString(t, "op.d", bg.op.d, "0")
	assertBigIntString(t, "op.e", bg.op.e, "0")
	assertBigIntString(t, "op.f", bg.op.f, "0")
	assertBigIntString(t, "op.g", bg.op.g, "0")
	assertBigIntString(t, "op.h", bg.op.h, "1")

	leftTerm, _, err := bg.left.Next()
	if err != nil {
		t.Fatalf("original left.Next() error: %v", err)
	}
	if !leftTerm.IsValue() {
		t.Fatalf("original left child should remain unconsumed")
	}
}

// cf/binary_wb_test.go v1
