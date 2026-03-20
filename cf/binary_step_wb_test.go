// cf/binary_step_wb_test.go v1
package cf

import "testing"

func TestWB_BinaryGCF_Step_PrefersLeftWhenBothChildrenCanIngestAndNoEmitIsCertified(t *testing.T) {
	g := Add(FromSource(Int64(3)), FromSource(Sqrt2()))

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
func TestWB_BinaryGCF_Step_ChoosesRightWhenLeftIsExhaustedAndNoEmitIsCertified(t *testing.T) {
	g := Add(FromSource(Int64(3)), FromSource(Sqrt2()))

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
	if !rightTerm.IsValue() {
		t.Fatalf("right child should still be a value-producing stream after right-ingest step")
	}
}
func TestWB_BinaryGCF_Step_ChoosesEmitWhenNeitherChildCanIngestAndExactEmitIsAvailable(t *testing.T) {
	g := Add(FromSource(Int64(3)), FromSource(Int64(5)))

	bg, ok := g.(binaryGCF)
	if !ok {
		t.Fatalf("Add() concrete type = %T, want binaryGCF", g)
	}

	next1, err := bg.ingestLeftStep()
	if err != nil {
		t.Fatalf("ingestLeftStep() error: %v", err)
	}

	next2, err := next1.ingestRightStep()
	if err != nil {
		t.Fatalf("ingestRightStep() error: %v", err)
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

// cf/binary_step_wb_test.go v1
