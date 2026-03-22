// cf/source_wb_test.go v3
package cf

import (
	"math/big"
	"testing"
)

func TestWB_FromSource_Sqrt2_TakeDoesNotMutateOriginalEvaluator(t *testing.T) {
	g := FromSource(Sqrt2())

	sb, ok := g.(sourceBackedGCF)
	if !ok {
		t.Fatalf("FromSource(Sqrt2()) concrete type = %T, want sourceBackedGCF", g)
	}

	src0, ok := sb.src.(sqrt2Source)
	if !ok {
		t.Fatalf("sourceBackedGCF src concrete type = %T, want sqrt2Source", sb.src)
	}
	if src0.index != 0 {
		t.Fatalf("initial sqrt2 index = %d, want 0", src0.index)
	}

	_, _ = g.Take(3)

	src1, ok := sb.src.(sqrt2Source)
	if !ok {
		t.Fatalf("sourceBackedGCF src concrete type after Take = %T, want sqrt2Source", sb.src)
	}
	if src1.index != 0 {
		t.Fatalf("Take mutated original evaluator source index = %d, want 0", src1.index)
	}

	term, rest, err := g.Next()
	if err != nil {
		t.Fatalf("original evaluator Next() error: %v", err)
	}
	if !term.IsValue() {
		t.Fatalf("original evaluator first term should still be a value")
	}
	value, ok := term.BigInt()
	if !ok {
		t.Fatalf("original evaluator first term should expose a BigInt")
	}
	if got, want := value.String(), "1"; got != want {
		t.Fatalf("original evaluator first term = %s, want %s", got, want)
	}

	sb1, ok := rest.(sourceBackedGCF)
	if !ok {
		t.Fatalf("remainder concrete type = %T, want sourceBackedGCF", rest)
	}
	srcAfterNext, ok := sb1.src.(sqrt2Source)
	if !ok {
		t.Fatalf("remainder src concrete type = %T, want sqrt2Source", sb1.src)
	}
	if srcAfterNext.index != 1 {
		t.Fatalf("remainder sqrt2 index = %d, want 1", srcAfterNext.index)
	}
}

func TestWB_RCFPrefixGCF_ConvergentOnConcretePrefix_SevenOverFive(t *testing.T) {
	prefix := rcfPrefixGCF{
		terms: []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(2)},
		index: 0,
	}

	r, err := prefix.Convergent()
	if err != nil {
		t.Fatalf("Convergent() error: %v", err)
	}

	if got, want := r.Num.String(), "7"; got != want {
		t.Fatalf("numerator = %s, want %s", got, want)
	}
	if got, want := r.Den.String(), "5"; got != want {
		t.Fatalf("denominator = %s, want %s", got, want)
	}
}

func TestWB_RCFPrefixGCF_ConvergentUsesRemainingTermsFromIndex(t *testing.T) {
	prefix := rcfPrefixGCF{
		terms: []*big.Int{big.NewInt(9), big.NewInt(1), big.NewInt(2), big.NewInt(2)},
		index: 1,
	}

	r, err := prefix.Convergent()
	if err != nil {
		t.Fatalf("Convergent() error: %v", err)
	}

	if got, want := r.Num.String(), "7"; got != want {
		t.Fatalf("numerator = %s, want %s", got, want)
	}
	if got, want := r.Den.String(), "5"; got != want {
		t.Fatalf("denominator = %s, want %s", got, want)
	}
}

func TestWB_BinaryLFT_IngestLeft_IdentityWithPQ_1_2(t *testing.T) {
	b := identityBinaryLFT()

	got := b.ingestLeft(big.NewInt(1), big.NewInt(2))

	assertBigIntString(t, "got.a", got.a, "2")
	assertBigIntString(t, "got.b", got.b, "0")
	assertBigIntString(t, "got.c", got.c, "1")
	assertBigIntString(t, "got.d", got.d, "0")
	assertBigIntString(t, "got.e", got.e, "0")
	assertBigIntString(t, "got.f", got.f, "1")
	assertBigIntString(t, "got.g", got.g, "0")
	assertBigIntString(t, "got.h", got.h, "0")
}

func TestWB_BinaryLFT_IngestRight_IdentityWithPQ_1_2(t *testing.T) {
	b := identityBinaryLFT()

	got := b.ingestRight(big.NewInt(1), big.NewInt(2))

	assertBigIntString(t, "got.a", got.a, "2")
	assertBigIntString(t, "got.b", got.b, "1")
	assertBigIntString(t, "got.c", got.c, "0")
	assertBigIntString(t, "got.d", got.d, "0")
	assertBigIntString(t, "got.e", got.e, "0")
	assertBigIntString(t, "got.f", got.f, "0")
	assertBigIntString(t, "got.g", got.g, "1")
	assertBigIntString(t, "got.h", got.h, "0")
}

func TestWB_BinaryLFT_IngestLeft_NonTrivialCoefficients(t *testing.T) {
	b := binaryLFT{
		a: big.NewInt(2),
		b: big.NewInt(3),
		c: big.NewInt(5),
		d: big.NewInt(7),
		e: big.NewInt(11),
		f: big.NewInt(13),
		g: big.NewInt(17),
		h: big.NewInt(19),
	}

	got := b.ingestLeft(big.NewInt(23), big.NewInt(29))

	assertBigIntString(t, "got.a", got.a, "63")
	assertBigIntString(t, "got.b", got.b, "94")
	assertBigIntString(t, "got.c", got.c, "46")
	assertBigIntString(t, "got.d", got.d, "69")
	assertBigIntString(t, "got.e", got.e, "336")
	assertBigIntString(t, "got.f", got.f, "396")
	assertBigIntString(t, "got.g", got.g, "253")
	assertBigIntString(t, "got.h", got.h, "299")
}

func TestWB_BinaryLFT_IngestRight_NonTrivialCoefficients(t *testing.T) {
	b := binaryLFT{
		a: big.NewInt(2),
		b: big.NewInt(3),
		c: big.NewInt(5),
		d: big.NewInt(7),
		e: big.NewInt(11),
		f: big.NewInt(13),
		g: big.NewInt(17),
		h: big.NewInt(19),
	}

	got := b.ingestRight(big.NewInt(23), big.NewInt(29))

	assertBigIntString(t, "got.a", got.a, "61")
	assertBigIntString(t, "got.b", got.b, "46")
	assertBigIntString(t, "got.c", got.c, "152")
	assertBigIntString(t, "got.d", got.d, "115")
	assertBigIntString(t, "got.e", got.e, "332")
	assertBigIntString(t, "got.f", got.f, "253")
	assertBigIntString(t, "got.g", got.g, "512")
	assertBigIntString(t, "got.h", got.h, "391")
}

func TestWB_BinaryLFT_IngestLeft_DoesNotMutateOriginal(t *testing.T) {
	b := identityBinaryLFT()

	_ = b.ingestLeft(big.NewInt(1), big.NewInt(2))

	assertBigIntString(t, "b.a", b.a, "1")
	assertBigIntString(t, "b.b", b.b, "0")
	assertBigIntString(t, "b.c", b.c, "0")
	assertBigIntString(t, "b.d", b.d, "0")
	assertBigIntString(t, "b.e", b.e, "0")
	assertBigIntString(t, "b.f", b.f, "0")
	assertBigIntString(t, "b.g", b.g, "0")
	assertBigIntString(t, "b.h", b.h, "1")
}

func TestWB_BinaryLFT_IngestRight_DoesNotMutateOriginal(t *testing.T) {
	b := identityBinaryLFT()

	_ = b.ingestRight(big.NewInt(1), big.NewInt(2))

	assertBigIntString(t, "b.a", b.a, "1")
	assertBigIntString(t, "b.b", b.b, "0")
	assertBigIntString(t, "b.c", b.c, "0")
	assertBigIntString(t, "b.d", b.d, "0")
	assertBigIntString(t, "b.e", b.e, "0")
	assertBigIntString(t, "b.f", b.f, "0")
	assertBigIntString(t, "b.g", b.g, "0")
	assertBigIntString(t, "b.h", b.h, "1")
	// EOF cf/source_wb_test.go V3
}
