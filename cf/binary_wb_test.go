// cf/binary_wb_test.go v2
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
	assertBigIntString(t, "op.b", bg.op.b, "-1")
	assertBigIntString(t, "op.c", bg.op.c, "1")
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
	assertBigIntString(t, "op.b", bg.op.b, "0")
	assertBigIntString(t, "op.c", bg.op.c, "1")
	assertBigIntString(t, "op.d", bg.op.d, "0")
	assertBigIntString(t, "op.e", bg.op.e, "0")
	assertBigIntString(t, "op.f", bg.op.f, "1")
	assertBigIntString(t, "op.g", bg.op.g, "0")
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
		t.Fatalf("right child should have advanced to EOF after right-ingest step")
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

func TestWB_BinaryGCF_Step_PrefersLeftWhenBothChildrenCanIngest(t *testing.T) {
	g := Add(FromSource(Int64(3)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Add() concrete type = %T, want binaryGCF", g)
	}

	action, next, err := bg.step()
	if err != nil {
		t.Fatalf("step() error: %v", err)
	}
	if action != decisionIngestLeft {
		t.Fatalf("step() action = %v, want %v", action, decisionIngestLeft)
	}

	assertBigIntString(t, "op.a", next.op.a, "1")
	assertBigIntString(t, "op.b", next.op.b, "3")
	assertBigIntString(t, "op.c", next.op.c, "0")
	assertBigIntString(t, "op.d", next.op.d, "1")
}

func TestWB_BinaryGCF_Step_ChoosesRightWhenLeftIsExhausted(t *testing.T) {
	g := Add(FromSource(Int64(3)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Add() concrete type = %T, want binaryGCF", g)
	}

	_, next1, err := bg.step()
	if err != nil {
		t.Fatalf("first step() error: %v", err)
	}

	action2, next2, err := next1.step()
	if err != nil {
		t.Fatalf("second step() error: %v", err)
	}
	if action2 != decisionIngestRight {
		t.Fatalf("second step() action = %v, want %v", action2, decisionIngestRight)
	}

	rightTerm, _, err := next2.right.Next()
	if err != nil {
		t.Fatalf("next2.right.Next() error: %v", err)
	}
	if !rightTerm.IsEOF() {
		t.Fatalf("right child should be exhausted after right-ingest step")
	}
}

func TestWB_BinaryGCF_Step_ChoosesEmitWhenNeitherChildCanIngestAndExactEmitIsAvailable(t *testing.T) {
	g := Add(FromSource(Int64(3)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Add() concrete type = %T, want binaryGCF", g)
	}

	_, next1, err := bg.step()
	if err != nil {
		t.Fatalf("first step() error: %v", err)
	}

	_, next2, err := next1.step()
	if err != nil {
		t.Fatalf("second step() error: %v", err)
	}

	action3, _, err := next2.step()
	if err != nil {
		t.Fatalf("third step() error: %v", err)
	}
	if action3 != decisionEmitOutput {
		t.Fatalf("third step() action = %v, want %v", action3, decisionEmitOutput)
	}
}

func TestWB_BinaryGCF_Step_DoesNotMutateOriginalNode(t *testing.T) {
	g := Add(FromSource(Int64(3)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Add() concrete type = %T, want binaryGCF", g)
	}

	_, _, err := bg.step()
	if err != nil {
		t.Fatalf("step() error: %v", err)
	}

	leftTerm, _, err := bg.left.Next()
	if err != nil {
		t.Fatalf("original left.Next() error: %v", err)
	}
	if !leftTerm.IsValue() {
		t.Fatalf("original left child should remain unconsumed")
	}

	rightTerm, _, err := bg.right.Next()
	if err != nil {
		t.Fatalf("original right.Next() error: %v", err)
	}
	if !rightTerm.IsValue() {
		t.Fatalf("original right child should remain unconsumed")
	}

	assertBigIntString(t, "op.a", bg.op.a, "0")
	assertBigIntString(t, "op.b", bg.op.b, "1")
	assertBigIntString(t, "op.c", bg.op.c, "1")
	assertBigIntString(t, "op.d", bg.op.d, "0")
}

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

func TestWB_BinaryGCF_Next_FirstCallReturnsBinaryGCFRemainderWithResolvedState(t *testing.T) {
	g := Add(FromSource(Int64(3)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Add() concrete type = %T, want binaryGCF", g)
	}

	term, rest, err := bg.Next()
	if err != nil {
		t.Fatalf("Next() error: %v", err)
	}
	if !term.IsValue() {
		t.Fatalf("first term should be a value")
	}
	value, ok := term.BigInt()
	if !ok {
		t.Fatalf("first term should expose a BigInt")
	}
	if got, want := value.String(), "8"; got != want {
		t.Fatalf("first term = %s, want %s", got, want)
	}

	nextBG, ok := rest.(binaryGCF)
	if !ok {
		t.Fatalf("remainder concrete type = %T, want binaryGCF", rest)
	}
	if nextBG.resolved == nil {
		t.Fatalf("remainder binaryGCF should retain resolved continuation")
	}
}

func TestWB_BinaryGCF_Next_OnResolvedNodeWrapsResolvedRemainderBackIntoBinaryGCF(t *testing.T) {
	prefix, err := prefixGCFfromRational(Rational{
		Num: big.NewInt(12),
		Den: big.NewInt(5),
	})
	if err != nil {
		t.Fatalf("prefixGCFfromRational() error: %v", err)
	}

	bg := binaryGCF{resolved: prefix}

	term, rest, err := bg.Next()
	if err != nil {
		t.Fatalf("Next() error: %v", err)
	}
	if !term.IsValue() {
		t.Fatalf("first term should be a value")
	}
	value, ok := term.BigInt()
	if !ok {
		t.Fatalf("first term should expose a BigInt")
	}
	if got, want := value.String(), "2"; got != want {
		t.Fatalf("first term = %s, want %s", got, want)
	}

	nextBG, ok := rest.(binaryGCF)
	if !ok {
		t.Fatalf("remainder concrete type = %T, want binaryGCF", rest)
	}
	if nextBG.resolved == nil {
		t.Fatalf("resolved remainder should stay wrapped in binaryGCF")
	}
}

func TestWB_BinaryGCF_Next_DoesNotMutateOriginalNodeWhenResolving(t *testing.T) {
	g := Add(FromSource(Int64(3)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Add() concrete type = %T, want binaryGCF", g)
	}

	_, _, err := bg.Next()
	if err != nil {
		t.Fatalf("Next() error: %v", err)
	}

	if bg.resolved != nil {
		t.Fatalf("original node should remain immutable and unresolved")
	}

	leftTerm, _, err := bg.left.Next()
	if err != nil {
		t.Fatalf("original left.Next() error: %v", err)
	}
	if !leftTerm.IsValue() {
		t.Fatalf("original left child should remain unconsumed")
	}

	rightTerm, _, err := bg.right.Next()
	if err != nil {
		t.Fatalf("original right.Next() error: %v", err)
	}
	if !rightTerm.IsValue() {
		t.Fatalf("original right child should remain unconsumed")
	}
}

func TestWB_BinaryGCF_EmitStep_AfterBothIngestsOnDivTwelveByFive_EmitsTwoAndLeavesFiveOverTwo(t *testing.T) {
	g := Div(FromSource(Int64(12)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Div() concrete type = %T, want binaryGCF", g)
	}

	_, bg1, err := bg.step()
	if err != nil {
		t.Fatalf("first step() error: %v", err)
	}

	_, bg2, err := bg1.step()
	if err != nil {
		t.Fatalf("second step() error: %v", err)
	}

	term, next, err := bg2.emitStep()
	if err != nil {
		t.Fatalf("emitStep() error: %v", err)
	}
	if !term.IsValue() {
		t.Fatalf("emitStep() should emit a value term")
	}

	value, ok := term.BigInt()
	if !ok {
		t.Fatalf("emitStep() value term should expose a BigInt")
	}
	if got, want := value.String(), "2"; got != want {
		t.Fatalf("emitStep() term = %s, want %s", got, want)
	}

	r, ok := next.op.exactQuotient()
	if !ok {
		t.Fatalf("remainder exactQuotient() should succeed after emit")
	}
	assertBigIntString(t, "r.Num", r.Num, "5")
	assertBigIntString(t, "r.Den", r.Den, "2")
}

func TestWB_BinaryGCF_Step_ChoosesEmitWhenExactEmitIsAvailable(t *testing.T) {
	g := Div(FromSource(Int64(12)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Div() concrete type = %T, want binaryGCF", g)
	}

	_, bg1, err := bg.step()
	if err != nil {
		t.Fatalf("first step() error: %v", err)
	}

	_, bg2, err := bg1.step()
	if err != nil {
		t.Fatalf("second step() error: %v", err)
	}

	action, next, err := bg2.step()
	if err != nil {
		t.Fatalf("third step() error: %v", err)
	}
	if action != decisionEmitOutput {
		t.Fatalf("third step() action = %v, want %v", action, decisionEmitOutput)
	}

	r, ok := next.op.exactQuotient()
	if !ok {
		t.Fatalf("remainder exactQuotient() should succeed after emit")
	}
	assertBigIntString(t, "r.Num", r.Num, "5")
	assertBigIntString(t, "r.Den", r.Den, "2")
}

// cf/binary_wb_test.go v2
