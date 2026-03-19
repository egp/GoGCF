<!-- HighLevelDesign.md v3 -->
# HighLevelDesign.md

# Continued Fraction Arithmetic Library High-Level Design

## 1. Purpose

This document describes the high-level design of the GoGCF project.

It is guided by the Requirements Specification and is intended to:

- identify the major objects
- describe their responsibilities
- describe how data flows through the system
- separate immutable public abstractions from internal transform machinery
- support iterative TDD and correctness-first implementation

This document is deliberately high level. It does not attempt to settle every implementation detail.

---

## 2. Design Priorities

The design is driven by the following priorities:

1. mathematical correctness
2. exact arithmetic using arbitrary precision integers
3. immutable remainder-returning stepping semantics
4. strong black-box and white-box testability
5. small public API
6. extensibility toward Gosper-style arithmetic and later unary/transcendental work
7. idiomatic Go structure and naming

Performance is secondary.

---

## 3. Architectural Overview

The design separates the system into four layers:

1. **Public immutable source layer**
2. **Public immutable evaluator layer**
3. **Internal transform/arithmetic layer**
4. **Support/value layer**

At a high level:

- `GCFSource` represents an immutable generalized continued-fraction input source
- `GCF` represents an immutable arithmetic evaluator/value that emits `RCFTerm`
- internal LFT machinery performs Gosper-style arithmetic
- support types such as `Rational`, `Bound`, and `Range` provide exact metadata and observability

---

## 4. Simple Diagram

    +------------------+      FromSource       +------------------+
    |    GCFSource     | --------------------> |       GCF        |
    | immutable input  |                       | immutable eval   |
    | emits PQTerm     |                       | emits RCFTerm    |
    +------------------+                       +------------------+
             |                                           |
             | built by                                   | observed by
             v                                           v
    +------------------+                       +------------------+
    | constructors and |                       | Range, Take,     |
    | named sources    |                       | Convergent       |
    +------------------+                       +------------------+
                                                        |
                                                        | implemented by
                                                        v
                                              +------------------+
                                              | ULFT / BLFT /    |
                                              | DiagonalLFT      |
                                              | internal engine  |
                                              +------------------+

---

## 5. Major Public Objects

### 5.1 GCFSource

#### Purpose

`GCFSource` is the immutable generalized continued-fraction input abstraction.

It represents the input side of the system.

#### Responsibilities

- provide the next generalized term as a `PQTerm`
- return a new immutable remainder source
- represent exact finite or infinite sources
- support named sources such as `Pi()`, `E()`, and `Sqrt2()`
- support source construction from integers and rationals

#### Public behavior

Conceptually:

    NextPQ() -> (PQTerm, GCFSource remainder, error)

#### Notes

- `GCFSource` is not itself an arithmetic evaluator
- it does not expose `Range()` or `Convergent()`
- source exhaustion is represented as `PQTerm` EOF rather than as a required exceptional condition

---

### 5.2 GCF

#### Purpose

`GCF` is the immutable arithmetic evaluator/value abstraction.

It represents the output side of the system.

#### Responsibilities

- emit the next `RCFTerm`
- return a new immutable remainder evaluator
- expose the current `Range`
- support `Take(n)` for finite prefix extraction
- support `Convergent()` for finite evaluators

#### Public behavior

Conceptually:

    Next() -> (RCFTerm, GCF remainder, error)
    Range() -> Range
    Take(n) -> (prefix GCF, error)
    Convergent() -> (Rational, error)

#### Notes

- `GCF` is immutable at the public API level
- all operational state is carried forward into the returned remainder evaluator
- no cloning is required in this model

---

### 5.3 PQTerm

#### Purpose

`PQTerm` is the public exact input-term representation for generalized continued fractions.

#### Responsibilities

- represent one ingested generalized term
- carry exact `BigInt` numerator/denominator-like components
- represent source EOF as an ordinary tagged term
- remain simple and transparent for tests

#### Notes

This is an ingestion type, not an emission type.

---

### 5.4 RCFTerm

#### Purpose

`RCFTerm` is the public emitted output-term abstraction.

#### Responsibilities

- represent ordinary emitted RCF terms
- represent output-side EOF

#### Notes

This avoids overloading `*big.Int` with non-numeric stream signals while keeping the initial public surface intentionally small.

---

### 5.5 Rational

#### Purpose

`Rational` is the exact finite numeric type used for convergents and bounds.

#### Responsibilities

- represent exact finite rational values
- normalize sign
- support comparison and formatting
- serve as the finite value carrier for `Bound`

---

### 5.6 Bound

#### Purpose

`Bound` represents one endpoint of a `Range`.

#### Responsibilities

- represent finite, `-Inf`, or `+Inf`
- represent open or closed endpoint state
- support endpoint comparison rules

#### Notes

A `Bound` is not itself a range. It is one half of a `Range`.

---

### 5.7 Range

#### Purpose

`Range` represents the current certified uncertainty set for a `GCF`.

#### Responsibilities

- represent uncertainty using `Lo`, `Hi`, and explicit `Outside`
- support `IsInside()`
- support `IsOutside()`
- support `IsExact()`
- support informativeness comparison against another range

#### Notes

Range is an observability and correctness tool, not merely a convenience.

---

## 6. Major Internal Objects

The following types are expected to exist internally even if they are not all public.

### 6.1 UnaryLFT

#### Purpose

Represents a unary linear fractional transform:

\[
z(x)=\frac{ax+b}{cx+d}
\]

#### Responsibilities

- hold transform coefficients
- ingest terms from one operand stream
- determine when output is certified
- emit output terms or defer ingestion
- support normalization and invariant checks

#### Required initializer

The ULFT identity matrix shall be available:

    (1,0)/(0,1)

This identity will be used to initialize selected evaluator states.

---

### 6.2 BinaryLFT

#### Purpose

Represents a binary linear fractional transform:

\[
z(x,y)=\frac{axy+bx+cy+d}{exy+fx+gy+h}
\]

#### Responsibilities

- hold bihomographic coefficients
- ingest terms from either operand
- decide whether to ingest from left, ingest from right, or emit output
- support exact arithmetic for `+`, `-`, `*`, and `/`
- support range and progress logic

#### Required initializer

The BLFT identity-style initializer shall be available:

    (1,0,0,0)/(0,0,0,1)

This identity will be used to initialize selected evaluator states.

#### Core decision machine

The most important part of the project is the internal decision machine whose binary-step action space is exactly:

- ingest left
- ingest right
- emit output

Nothing else is a primary binary evaluation action.

That decision machine should be represented explicitly in internal code and targeted directly by white-box tests.

---

### 6.3 DiagonalLFT

#### Purpose

Represents the special case of a `BinaryLFT` where the two operands are equal.

#### Responsibilities

- support efficient or simplified evaluation when `x == y`
- serve as a natural internal form for selected unary or algebraic operations
- remain mathematically equivalent to the corresponding binary form

---

### 6.4 Internal Evaluator Nodes

A `GCF` value will likely be backed internally by one of several evaluator-node kinds.

Possible examples:

- source-backed evaluator
- unary-op evaluator
- binary-op evaluator
- finite-prefix evaluator
- exhausted evaluator
- error/stuck sentinel evaluator

The exact concrete node representation is up to the implementation.

---

## 7. Data Flow

The normal data flow is:

1. build a `GCFSource`
2. lift it into a `GCF`
3. combine `GCF` values with unary/binary operators
4. repeatedly call `Next()`
5. use `Range()` and `Take(n)` for observation and finite analysis

Example conceptual flow:

    Int64 / Rat64 / Pi / E / Sqrt2
                |
                v
            GCFSource
                |
                v
            FromSource
                |
                v
               GCF
                |
       +--------+--------+
       |                 |
     unary ops        binary ops
       |                 |
       +--------+--------+
                |
                v
          Next / Take / Range
                |
                v
       RCFTerm / Rational / Range

---

## 8. Why Immutable Remainders

The design uses immutable remainder-returning semantics rather than public mutation.

### 8.1 Benefits

- easier reasoning
- easier black-box testing
- easier white-box debugging
- no cloning requirement
- natural match for infinite streams
- old states remain available if the caller keeps references

### 8.2 Cost

- more object creation conceptually
- implementation may later want structural sharing or memoization internally

This cost is acceptable because correctness and clarity come first.

---

## 9. Finite and Infinite Values

The system treats all sources as infinite until exhaustion.

### Source side

- `GCFSource.NextPQ()` emits a `PQTerm`
- source exhaustion appears as EOF in `PQTerm`

### Evaluator side

- `GCF.Next()` emits an `RCFTerm`
- output-side end of stream is represented by EOF in `RCFTerm`

This symmetry is intentional.

---

## 10. Arithmetic Construction Model

The intended construction model is:

- constructors create `GCFSource`
- `FromSource` creates `GCF`
- arithmetic operators consume and return `GCF`

Examples:

- `Int64(3)` returns a `GCFSource`
- `FromSource(Int64(3))` returns a `GCF`
- `Add(x, y)` returns a `GCF`

This keeps raw source construction separate from arithmetic evaluation.

---

## 11. First Named Infinite Source: sqrt(2)

The first named infinite procedural source will be `sqrt(2)`.

### Why `sqrt(2)` first

- mathematically simple
- well known continued fraction
- infinite but easy to verify
- excellent first black-box fixture
- useful before more difficult named constants are available

### Initial intended representation

The regular continued fraction for `sqrt(2)` is:

\[
[1; 2,2,2,2,\dots]
\]

This makes it a good first procedural source for early API and stepping tests.

### Design note

Even though the project is GCF-first internally, an early `sqrt(2)` source may be implemented in the simplest mathematically valid way available, as long as it conforms to the public `GCFSource` contract.

---

## 12. Testing Strategy at the Design Level

The design is intended to support both black-box and white-box tests.

### 12.1 Black-box tests

Use only the public API to verify:

- source construction
- evaluator lifting
- emitted term stepping
- range behavior
- finite prefix extraction
- convergent calculation
- target formula expression shape

Recommended conventions:

- package `cf_test`
- file suffix `*_bb_test.go`
- test names prefixed `TestBB_`

### 12.2 White-box tests

Use same-package test access to verify:

- internal coefficient transitions
- normalization rules
- invariant checks
- stuck/non-progress detection
- transform scheduling choices
- explicit three-action decision-machine behavior

Recommended conventions:

- package `cf`
- file suffix `*_wb_test.go`
- test names prefixed `TestWB_`

### 12.3 Pending tests

Pending or intentionally not-yet-green tests should remain in the same package directory and use opt-in build tags such as:

    //go:build pending

### 12.4 TDD pattern

The normal pattern remains:

1. write failing tests
2. allow simple stubs of correct type
3. verify red
4. implement until green

---

## 13. Proposed Initial File Responsibilities

The first coding slice can be organized around these responsibilities.

### `api_types.go`

Public interfaces and small public type declarations.

### `pq_term.go`

Definition and helpers for `PQTerm`.

### `rcf_term.go`

Definition and helpers for `RCFTerm`.

### `rational.go`

Exact finite rational type and helpers.

### `bound.go`

Bound endpoint representation and comparison helpers.

### `range.go`

Range logic and comparison ordering.

### `source.go`

Core `GCFSource` abstractions and `FromSource`.

### `gcf.go`

Core `GCF` interface and shared evaluator helpers.

### `sqrt2_source.go`

First named infinite procedural `GCFSource`.

### `api_bb_test.go`

Initial black-box API tests.

### `sqrt2_source_bb_test.go`

Initial black-box tests for the first infinite named source.

### `decision_machine_wb_test.go`

Early white-box tests for the three-action binary decision machine.

---

## 14. Resolved Design Decisions

The following decisions are considered settled for the current phase:

- immutable `GCFSource`
- immutable `GCF`
- constructors create `GCFSource`
- arithmetic operates on `GCF`
- no public cloning
- `Take(n)` returns only the finite prefix and an error
- `PQTerm` includes EOF
- `RCFTerm` initially includes only value and EOF
- `Range` stores `Outside` explicitly
- ULFT and BLFT identity initializers are required
- BB and WB tests stay in the same package directory but are distinguished by package, file suffix, and test-name prefix

---

## 15. Immediate Next Steps

After this design document, the recommended next implementation steps are:

1. create the directory structure
2. create the public API skeleton types
3. write initial black-box tests for:
   - `GCFSource`
   - `GCF`
   - `PQTerm`
   - `RCFTerm`
   - `Range`
4. implement the first named infinite source: `sqrt(2)`
5. add production stubs for not-yet-implemented arithmetic pieces
6. keep arithmetic operators simple or stubbed until the first tests are in place

This keeps development aligned with the Requirements Specification and the TDD workflow.

<!-- HighLevelDesign.md v3 -->