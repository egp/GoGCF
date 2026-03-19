<!-- PublicAPI.md v4 -->
# PublicAPI.md

# Continued Fraction Arithmetic Library Public API Draft

## Status

This document is the current public API draft for the continued-fraction arithmetic library.

This version adopts the current design decisions:

- the Requirements Specification is the north-pole document
- `GCFSource` is immutable
- `GCF` is immutable
- evaluation proceeds by returning a remainder object, not by mutating state
- constructors produce `GCFSource`
- arithmetic operates on `GCF`
- emitted terms are tagged objects, not bare `BigInt`
- public `RCFTerm` currently exposes only value and EOF
- `Range` stores `Outside` explicitly
- `Take(n)` returns only the finite prefix and an error
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
8. A finite GCF is a GCF whose source has been exhausted or whose finite prefix has been materialized.
9. The public API should be as small as possible while still supporting the target formula.
10. Performance is secondary to correctness, testability, and debuggability.

---

## 3. Core Abstractions

This API distinguishes between:

- `GCFSource`: immutable generalized-CF input source
- `GCF`: immutable arithmetic evaluator/value that emits `RCFTerm`

These are related, but not identical.

### 3.1 GCFSource

`GCFSource` is an immutable description and remainder-carrier for a generalized continued-fraction input stream.

Proposed declaration:

    type GCFSource interface {
        NextPQ() (PQTerm, GCFSource, error)
    }

Semantics:

- `NextPQ()` returns one generalized input term and the remainder of the source
- the remainder is another immutable `GCFSource`
- ordinary source EOF is represented by `PQTerm`
- `error` is reserved for exceptional source behavior

### 3.2 GCF

`GCF` is an immutable arithmetic evaluator/value.

Proposed declaration:

    type GCF interface {
        Next() (RCFTerm, GCF, error)
        Range() Range
        Take(n int) (GCF, error)
        Convergent() (Rational, error)
    }

Semantics:

- `Next()` returns one emitted output term and the remainder evaluator
- `Range()` observes the current evaluator state
- `Take(n)` returns a finite prefix evaluator and an error
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

### 5.1 PQTermKind

Proposed declaration:

    type PQTermKind int

    const (
        PQValue PQTermKind = iota
        PQEOF
    )

### 5.2 PQTerm

`PQTerm` is the public generalized input-term type.

Proposed declaration:

    type PQTerm struct {
        Kind PQTermKind
        P    *big.Int
        Q    *big.Int
    }

Semantics:

- if `Kind == PQValue`, `P` and `Q` are meaningful
- if `Kind == PQEOF`, `P` and `Q` are nil

Required helper methods:

    func (t PQTerm) IsValue() bool
    func (t PQTerm) IsEOF() bool

---

## 6. Output Term Type

Bare `BigInt` is not sufficient for emitted output because the output stream must also represent EOF.

### 6.1 RCFTermKind

Proposed declaration:

    type RCFTermKind int

    const (
        RCFValue RCFTermKind = iota
        RCFEOF
    )

### 6.2 RCFTerm

`RCFTerm` is the public tagged output object returned by `GCF.Next()`.

Proposed declaration:

    type RCFTerm struct {
        Kind  RCFTermKind
        Value *big.Int
    }

Semantics:

- if `Kind == RCFValue`, `Value` must be non-nil and is the emitted RCF term
- if `Kind == RCFEOF`, `Value` must be nil

Required helper methods:

    func (t RCFTerm) IsValue() bool
    func (t RCFTerm) IsEOF() bool
    func (t RCFTerm) BigInt() (*big.Int, bool)

### 6.3 Error vs Term-Channel Split

The term channel is intended for ordinary stream results:

- emitted value term
- EOF

Go `error` values are reserved for exceptional conditions such as:

- malformed source behavior
- invariant failures
- implementation-detected stuck or non-progress states

The public API intentionally does not yet expose non-finite emitted terms such as `+Inf` or `-Inf`.

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
        Lo      Bound
        Hi      Bound
        Outside bool
    }

Required methods:

    func (r Range) IsInside() bool
    func (r Range) IsOutside() bool
    func (r Range) IsExact() bool
    func (r Range) Cmp(other Range) int

Semantics:

- `IsInside()` means `Outside == false`
- `IsOutside()` means `Outside == true`
- `IsExact()` means the represented set is a single exact value

### 7.4 Exactness

A range is exact when it denotes a single exact value, for example `[x, x]`.

### 7.5 Range Comparison

`Range.Cmp(other)` orders ranges by informativeness.

The intended ordering is:

1. exact
2. inside narrow
3. inside wide
4. outside narrow
5. outside wide

Among same-class ranges, narrower finite span is better than wider finite span.

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
    func Sqrt2() GCFSource

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

### 10.3 Comparison

    type Ordering int

    const (
        Less Ordering = -1
        Equal Ordering = 0
        Greater Ordering = 1
    )

    func Compare(x, y GCF) (Ordering, error)

Intended behavior:

- ingest as needed until a certified comparison is available
- leave `x` and `y` unchanged because they are immutable
- return an error if evaluation becomes undefined or stuck

---

## 11. Evaluator Methods

### 11.1 Next

    func (z GCF) Next() (RCFTerm, GCF, error)

Behavior:

- returns the next emitted output term
- returns the remainder evaluator
- ordinary waiting for more source information is handled internally and is not a public “blocked” condition
- `RCFEOF` is returned in the term channel
- exceptional failures are returned through `error`

### 11.2 Take

    func (z GCF) Take(n int) (GCF, error)

Behavior:

- consumes up to `n` emitted RCF terms from `z`
- returns a finite RCF-backed prefix evaluator
- if fewer than `n` terms are available because the evaluator ended, the returned prefix is still valid and `err == io.EOF`
- non-EOF errors indicate exceptional conditions

This is the preferred way to materialize a finite prefix from an infinite evaluation.

### 11.3 Convergent

    func (z GCF) Convergent() (Rational, error)

Semantics:

- `Convergent()` produces a `Rational` from a finite evaluator
- it is intended especially for finite RCF-backed values such as those returned by `Take(n)`
- calling `Convergent()` on a non-finite evaluator is an error

### 11.4 Range

    func (z GCF) Range() Range

Every `GCF` must provide its current certified range.

---

## 12. Identity Initializers for Internal Machinery

The public API does not need to expose ULFT and BLFT directly yet, but the design assumes internal identity initializers for evaluator construction.

Required internal identities:

- ULFT identity: `(1,0)/(0,1)`
- BLFT identity-style initializer: `(1,0,0,0)/(0,0,0,1)`

These are part of the internal design contract.

---

## 13. Error Model

The public API should define a small set of sentinel errors.

Proposed sentinel errors:

    var ErrUndefined error
    var ErrStuck error

Intended meanings:

- `ErrUndefined`:
  mathematically undefined or currently unsupported result under the initial public surface

- `ErrStuck`:
  implementation detected non-progress or a strategy failure unlikely to advance without intervention

Ordinary source or evaluator exhaustion is not an error for `Next()`. It is represented as EOF in the term object.

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

The following are not public in this draft:

- cloning
- memoization controls
- explicit public RCF constructor
- public ULFT / BLFT / DiagonalLFT types
- public transform-coefficient inspection
- decimal or radix digit emission
- public bounded rational-collapse unary operator
- public non-finite emitted output terms
- concurrency guarantees

These may be added later if needed.

---

## 16. Test Conventions

Recommended conventions for test organization:

### 16.1 Black-box tests

- package `cf_test`
- file suffix `*_bb_test.go`
- test names prefixed `TestBB_`

### 16.2 White-box tests

- package `cf`
- file suffix `*_wb_test.go`
- test names prefixed `TestWB_`

### 16.3 Pending tests

Use opt-in build tags such as:

    //go:build pending

Pending tests should remain in the same package directory.

---

## 17. Resolved API Decisions

The following decisions are considered settled for the current phase:

- constructors create `GCFSource`
- arithmetic operates on `GCF`
- `GCFSource` is immutable
- `GCF` is immutable
- `PQTerm` includes value and EOF
- `RCFTerm` includes value and EOF
- `Range` stores `Outside` explicitly
- `Take(n)` returns only the finite prefix and an error
- cloning is not needed publicly
- public non-finite term output is deferred

---

## 18. Immediate Next Implementation Slice

The first code slice should define the public skeleton and the first infinite source:

- `api_types.go`
- `pq_term.go`
- `rcf_term.go`
- `rational.go`
- `bound.go`
- `range.go`
- `source.go`
- `gcf.go`
- `sqrt2_source.go`

The first tests should focus on:

- `GCFSource`
- `PQTerm`
- `RCFTerm`
- `Range`
- `FromSource`
- `Sqrt2()`

<!-- PublicAPI.md v4 -->