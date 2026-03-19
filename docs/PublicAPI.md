<!-- PublicAPI.md  v3 -->
# PublicAPI.md

# Continued Fraction Arithmetic Library Public API Draft V3

## Status

This document is the third draft of the public API for the continued-fraction arithmetic library.

This version adopts the current design decisions:

- the Requirements Specification is the north-pole document
- `GCFSource` is immutable
- `GCF` is immutable
- evaluation proceeds by returning a remainder object, not by mutating state
- ingestion is GCF-oriented
- emission is RCF-oriented
- emitted terms are tagged objects, not bare `BigInt`
- division by exact zero may yield positive or negative infinity instead of an error
- cloning is not part of the public API
- memoization, if needed later, is an internal optimization only

The first implementation target is Go.

---

## 1. Package

Proposed package name:

    package cf

---

## 2. Public Design Principles

1. The Requirements Specification is authoritative when API decisions are ambiguous.
2. `GCFSource` is immutable.
3. `GCF` is immutable.
4. Every stepping operation returns:
   - the next emitted term
   - a new immutable remainder object
5. External ingestion is GCF-oriented.
6. External emission is RCF-oriented.
7. All inputs are assumed infinite until exhausted.
8. A finite GCF is a GCF whose underlying source has been exhausted.
9. The public API should be as small as possible while still supporting the target formula.
10. Performance is secondary to correctness, testability, and debuggability.

---

## 3. Core Abstractions

This API distinguishes between:

- `GCFSource`: immutable generalized-CF input source
- `GCF`: immutable arithmetic evaluator/value that emits RCF terms

These are related, but not identical.

### 3.1 GCFSource

`GCFSource` is an immutable description of a generalized continued-fraction input stream.

Examples:

- exact integer source
- exact rational source
- named constant source such as `Pi()` or `E()`
- procedural `(p,q)` source adapter

A `GCFSource` does not expose arithmetic emission methods such as `Range()` or `Convergent()`.

Proposed declaration:

    type GCFSource interface {
        NextPQ() (PQTerm, GCFSource, error)
    }

Semantics:

- `NextPQ()` returns one generalized input term and the remainder of the source
- the remainder is another immutable `GCFSource`
- `io.EOF` means the input source is exhausted
- sources are assumed infinite until `io.EOF` is observed

### 3.2 GCF

`GCF` is an immutable arithmetic evaluator/value.

A `GCF` may internally encode transform coefficients, range state, emitted-prefix state, and progress state, but none of that is mutated in place through the public API.

Proposed declaration:

    type GCF interface {
        Next() (RCFTerm, GCF, error)
        Range() Range
        Take(n int) (prefix GCF, rest GCF, err error)
        Convergent() (Rational, error)
    }

Semantics:

- `Next()` returns one emitted output term and the remainder evaluator
- `Range()` observes the current evaluator state
- `Take(n)` returns a finite prefix evaluator and the remainder evaluator
- `Convergent()` is intended for finite evaluators

---

## 4. Exact Numeric Types

### 4.1 Rational

`Rational` is the public exact rational type used in bounds, convergents, and exact finite results.

Proposed declaration:

    type Rational struct {
        Num *big.Int
        Den *big.Int
    }

Required invariants:

- `Den != 0`
- normalized sign convention
- reduced form preferred for public results

---

## 5. Input Term Type

### 5.1 PQTerm

`PQTerm` is the public generalized input-term type.

Proposed declaration:

    type PQTerm struct {
        P *big.Int
        Q *big.Int
    }

Semantics:

- `P` and `Q` are exact `BigInt` values
- this is the ingestion format for generalized continued fractions

Input-source exhaustion is reported by `io.EOF`, not by a special `PQTerm` value.

---

## 6. Output Term Type

Bare `BigInt` is not sufficient for emitted output because the output stream must also represent `EOF` and `±Inf`.

### 6.1 RCFTermKind

Proposed declaration:

    type RCFTermKind int

    const (
        TermRCF RCFTermKind = iota
        TermEOF
        TermPosInf
        TermNegInf
    )

### 6.2 RCFTerm

`RCFTerm` is the public tagged output object returned by `GCF.Next()`.

Proposed declaration:

    type RCFTerm struct {
        Kind  RCFTermKind
        Value *big.Int
    }

Semantics:

- if `Kind == TermRCF`, `Value` must be non-nil and is the emitted RCF term
- if `Kind != TermRCF`, `Value` must be nil

Required helper methods:

    func (t RCFTerm) IsRCF() bool
    func (t RCFTerm) IsEOF() bool
    func (t RCFTerm) IsPosInf() bool
    func (t RCFTerm) IsNegInf() bool
    func (t RCFTerm) BigInt() (*big.Int, bool)

`BigInt()` returns `(value, true)` only for `TermRCF`.

### 6.3 Error vs Term-Channel Split

The output-term channel is used for ordinary arithmetic stream results:

- emitted RCF term
- EOF
- positive infinity
- negative infinity

Go `error` values are reserved for exceptional conditions such as:

- malformed source behavior
- truly undefined operations such as `0/0`
- invariant failures
- implementation-detected stuck or non-progress states

---

## 7. Bounds and Ranges

### 7.1 BoundKind

A range endpoint may be finite or infinite.

Proposed declaration:

    type BoundKind int

    const (
        BoundFinite BoundKind = iota
        BoundNegInf
        BoundPosInf
    )

### 7.2 Bound

A `Bound` is an interval endpoint with:

- finite/infinite kind
- exact value when finite
- open/closed endpoint status

Proposed declaration:

    type Bound struct {
        Kind   BoundKind
        Value  Rational
        Closed bool
    }

Semantics:

- `Kind == BoundFinite` means `Value` is meaningful
- `Kind == BoundNegInf` means negative infinity
- `Kind == BoundPosInf` means positive infinity
- `Closed == true` means the endpoint is included
- `Closed == false` means the endpoint is excluded

### 7.3 Range

`Range` describes the current uncertainty set for a `GCF`.

Proposed declaration:

    type Range struct {
        Lo Bound
        Hi Bound
    }

This representation supports both inside and outside semantics.

Required methods:

    func (r Range) IsInside() bool
    func (r Range) IsOutside() bool
    func (r Range) IsExact() bool
    func (r Range) Cmp(other Range) int

Semantics:

- `IsInside()` means the value lies inside the interval/set described by `Lo` and `Hi`
- `IsOutside()` means the value lies outside the interval/set described by `Lo` and `Hi`
- `IsExact()` means the value is known exactly

Inside/outside is inferred from endpoint ordering together with open/closed endpoint semantics. It is not stored as a separate explicit flag.

### 7.4 Exactness

A range is exact when it denotes a single exact value.

Example:

- `[x, x]` is exact
- `(x, x)` is not exact

### 7.5 Range Comparison

`Range.Cmp(other)` orders ranges by informativeness.

The intended ordering remains:

1. inside narrow
2. inside wide
3. outside wide
4. outside narrow

The exact metric for “narrow” remains to be formalized.

---

## 8. Source Constructors

Special constructors create `GCFSource`, not `GCF`.

### 8.1 Exact Integer Sources

    func Int64(v int64) GCFSource
    func BigInt(v *big.Int) GCFSource

### 8.2 Exact Rational Sources

    func Rat64(num, den int64) GCFSource
    func RationalSource(num, den *big.Int) GCFSource

These represent exact rational values.

### 8.3 Procedural GCF Sources

For user-supplied generalized continued-fraction generators:

    func SourceFromFunc(fn func() (PQTerm, GCFSource, error)) GCFSource

This supports user-defined immutable sources directly in remainder-returning style.

### 8.4 Named Constant Sources

Initially required:

    func Pi() GCFSource
    func E() GCFSource

Each call returns a fresh immutable source.

### 8.5 Explicit RCF Source Constructor

An explicit RCF constructor is useful for tests and fixtures, but is not part of the public API in this draft.

It may exist privately inside the package for internal tests and helpers.

---

## 9. Lifting Source to Evaluator

A raw `GCFSource` must be lifted into a `GCF` before arithmetic emission begins.

Proposed declaration:

    func FromSource(src GCFSource) GCF

This function creates the initial evaluator for a source.

---

## 10. Arithmetic Combinators

The public arithmetic combinators operate on immutable `GCF` values and return immutable `GCF` values.

### 10.1 Unary Operators

    func Neg(x GCF) GCF
    func Inv(x GCF) GCF
    func Sqrt(x GCF) GCF
    func Sin(x GCF) GCF
    func Tanh(x GCF) GCF

Notes:

- `Sin` expects radians
- degree-based expressions should currently be written explicitly, for example:

      Sin(Mul(FromSource(Rat64(69, 180)), FromSource(Pi())))

### 10.2 Binary Operators

    func Add(x, y GCF) GCF
    func Sub(x, y GCF) GCF
    func Mul(x, y GCF) GCF
    func Div(x, y GCF) GCF

### 10.3 Division by Zero Semantics

Division by exact zero is not automatically an error.

When mathematically certified, division by zero may yield:

- `TermPosInf`
- `TermNegInf`

through ordinary `Next()` emission.

Examples:

- positive finite numerator divided by exact zero may yield `+Inf`
- negative finite numerator divided by exact zero may yield `-Inf`

Truly undefined forms such as `0/0` remain errors.

---

## 11. Evaluator Methods

### 11.1 Next

    func (z GCF) Next() (RCFTerm, GCF, error)

Behavior:

- returns the next emitted output term
- returns the remainder evaluator
- ordinary waiting for more source information is handled internally and is not a public “blocked” condition
- `TermEOF`, `TermPosInf`, and `TermNegInf` are returned in the term channel
- exceptional failures are returned through `error`

### 11.2 Take

    func (z GCF) Take(n int) (prefix GCF, rest GCF, err error)

Behavior:

- consumes up to `n` emitted RCF terms from `z`
- returns:
  - a finite RCF-backed prefix evaluator
  - the remainder evaluator
- if fewer than `n` terms are available because the source terminated, the prefix is still valid and `err == io.EOF`
- non-EOF errors indicate exceptional conditions

This is the preferred way to materialize a finite prefix from an infinite evaluation.

### 11.3 Convergent

    func (z GCF) Convergent() (Rational, error)

Semantics:

- `Convergent()` produces a `Rational` from a finite evaluator
- it is intended especially for finite RCF-backed values such as those returned by `Take(n)`
- calling `Convergent()` on a non-finite evaluator is an error unless a later design revision broadens the meaning

### 11.4 Range

    func (z GCF) Range() Range

Every `GCF` must provide its current certified range.

---

## 12. Comparison

Comparison remains required by the Requirements Spec.

Current draft candidate:

    type Ordering int

    const (
        Less Ordering = -1
        Equal Ordering = 0
        Greater Ordering = 1
    )

Candidate function:

    func Compare(x, y GCF) (Ordering, error)

Intended behavior:

- ingest as needed until a certified comparison is available
- leave `x` and `y` unchanged because they are immutable
- return an error if evaluation becomes truly undefined or stuck

Whether a bounded / timeout-aware comparison variant is needed remains open.

---

## 13. Errors

The public API should define a small set of sentinel errors.

Proposed sentinel errors:

    var ErrUndefined error
    var ErrStuck error

Intended meanings:

- `ErrUndefined`:
  mathematically undefined operation, such as `0/0`

- `ErrStuck`:
  implementation detected non-progress or a strategy failure unlikely to advance without intervention

Ordinary source exhaustion is not an error for `Next()`. It is represented as `TermEOF`.

Ordinary division-by-zero-to-infinity is not an error when the sign is mathematically certified.

---

## 14. Minimal Expression for the Target Formula

The public API must be sufficient to express:

    sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))

Using the current draft API:

    three := FromSource(Int64(3))
    pi1   := FromSource(Pi())
    pi2   := FromSource(Pi())
    e     := FromSource(E())
    five  := FromSource(Int64(5))
    deg69 := FromSource(Rat64(69, 180))
    pi3   := FromSource(Pi())

    numerator := Sqrt(
        Add(
            Div(
                three,
                Mul(pi1, pi2),
            ),
            e,
        ),
    )

    denominator := Sub(
        Tanh(
            Sqrt(five),
        ),
        Sin(
            Mul(deg69, pi3),
        ),
    )

    target := Div(numerator, denominator)

No cloning is required because all values are immutable.

---

## 15. Public API Exclusions for This Draft

The following are not public in V3:

- cloning
- memoization controls
- explicit public RCF constructor
- public ULFT / BLFT / DiagonalLFT types
- public transform-coefficient inspection
- decimal or radix digit emission
- public bounded rational-collapse unary operator
- concurrency guarantees

These may be added later if needed.

---

## 16. Black-Box Test Expectations

The public API should support black-box tests that verify:

- exact source construction from integers and rationals
- named constant source construction
- lifting from `GCFSource` to `GCF`
- expression-building for the target formula
- tagged term emission through `Next()`
- finite prefix extraction through `Take(n)`
- convergent computation from finite prefixes
- range behavior
- target-formula emitted-prefix checks
- comparison behavior
- stuck / undefined behavior through public errors
- infinity emission through ordinary term output

---

## 17. Likely Internal Types Not Exposed Publicly

The following are expected internally but are not part of the public API draft:

- `UnaryLFT`
- `BinaryLFT`
- `DiagonalLFT`
- coefficient-normalization helpers
- invariant-check helpers
- private explicit RCF constructor for tests
- white-box tracing hooks
- progress scheduler internals
- memoization caches, if later needed

---

## 18. Open Questions

1. Should `Compare` be public in this phase, or deferred until stuck/progress behavior is better characterized?
2. Should `RCFTermKind` eventually include additional non-numeric stream states beyond `EOF` and `±Inf`?
3. Should degree-based convenience helpers such as `SinDeg` exist publicly, or should radians-only remain the rule?
4. Should `Take(n)` preserve metadata indicating whether the prefix ended by count limit or by source exhaustion, beyond the returned `error`?
5. Should `Range.Cmp` be defined in terms of exact interval width, projective width, or another uncertainty metric?
6. Should `Convergent()` remain finite-only permanently, or should a later API allow it to mean “convergent of the currently emitted prefix”?

Current draft answers:

- `Compare` remains provisional
- `RCFTermKind` stays small in V3
- radians-only remains the rule
- `Take(n)` uses `io.EOF` to indicate early exhaustion
- `Range.Cmp` remains semantically specified but not fully formalized
- `Convergent()` is finite-only in V3
<!-- EOF PublicAPI.md  v3 -->