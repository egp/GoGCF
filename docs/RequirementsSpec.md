<!-- RequirementsSpec.md v3 -->
# RequirementsSpec.md

# Continued Fraction Arithmetic Library Requirements Specification

## 1. Purpose

This document defines the requirements for a Go library that performs arithmetic on continued fractions using Gosper-style methods.

The library is intended to support mathematically correct, testable, demand-driven arithmetic on continued fractions, especially in the spirit of HAKMEM items 101A and 101B.

This specification describes required behavior, guarantees, and observable semantics. It is the north-pole document for API and design decisions, and it is expected to evolve as new questions arise.

---

## 2. Background and Context

Continued fractions provide an alternative numeric representation with useful arithmetic properties, especially for exact and incremental computation.

The target library is intended to support:

- exact arithmetic on continued fractions
- demand-driven evaluation
- regular continued fractions (RCF) and generalized continued fractions (GCF)
- composable unary and binary arithmetic operators
- certified partial output when full evaluation is not yet available

The initial architecture should be aligned with Gosper-style homographic and bihomographic transformation machinery, while remaining testable and implementation-language-aware at the requirements level.

A major project milestone is the ability to compute the target formula:

\[
\frac{\sqrt{\frac{3}{\pi^2} + e}}{\tanh(\sqrt{5}) - \sin(69^\circ)}
\]

This milestone strongly influences the minimum required operator set and the shape of the first public API.

---

## 3. Scope

### 3.1 In Scope

The first versions of the library shall address:

- representation of regular continued fractions (RCF)
- representation of generalized continued fractions (GCF)
- exact `BigInt` and `Rational` support
- construction of immutable `GCFSource` values from integers and rationals
- construction of immutable `GCF` evaluators from `GCFSource`
- demand-driven procedural infinite GCF sources
- exact comparison where possible
- unary and binary arithmetic over continued fractions
- demand-driven RCF term emission
- convergent and range inspection
- diagnostics and observability for blocked or stalled operator states
- specialized named generators for important constants, including at least `pi` and `e`
- at least one simple named infinite procedural source for early development, namely `sqrt(2)`
- minimum unary operators needed for the target formula, including `sqrt`, `sin`, and `tanh`

### 3.2 Out of Scope for Initial Release

The first release is not required to include:

- graphical tools
- distributed computation
- hardened security features
- decimal or other radix digit emission as a completed feature
- public non-finite RCF output terms such as `+Inf` or `-Inf`
- unary operators beyond those required for the target formula and immediate architectural needs

Robustness is important, but security concerns are explicitly deferred in the early phases.

---

## 4. Goals

The library shall prioritize the following goals:

1. **Mathematical correctness first**
2. **Exactness over floating approximation**
3. **Demand-driven incremental evaluation**
4. **Strong testability and observability**
5. **Composable operator architecture**
6. **Support for infinite GCF-first development**
7. **A minimal public API sufficient to compute the target formula**
8. **Persistent immutable stepping semantics**
9. **Extensibility toward more advanced operators later**
10. **Idiomatic Go structure and naming**

Performance is not a primary goal.

---

## 5. Non-Goals

The library is not required to optimize primarily for:

- fixed-width machine arithmetic
- approximate floating-point throughput
- minimal abstraction count
- hiding internal state from tests when observability improves debugging or correctness validation
- finite-CF-first design assumptions in early development

---

## 6. Terminology and Definitions

### 6.1 BigInt

An arbitrarily large integer.

### 6.2 Rational

An exact rational value represented as `BigInt/BigInt`.

### 6.3 Regular Continued Fraction (RCF)

A continued fraction in regular form whose emitted ordinary terms are `BigInt` values.

### 6.4 Generalized Continued Fraction (GCF)

A continued fraction whose ingested source terms are `(p, q)` pairs, where both `p` and `q` are `BigInt` values.

### 6.5 Relationship Among Core Types

The following containment and conversion relationships apply:

- any integer can be represented as a `BigInt`
- any `BigInt` can be represented as a `Rational`
- any `Rational` can be represented as a finite `RCF`
- any `RCF` is a `GCF`

### 6.6 GCFSource

An immutable generalized-continued-fraction input source.

Calling `NextPQ()` on a `GCFSource` returns:

- the next generalized input term as a `PQTerm`
- a new immutable `GCFSource` representing the remainder
- an error for exceptional conditions

### 6.7 GCF

An immutable arithmetic evaluator/value over continued fractions.

Calling `Next()` on a `GCF` returns:

- the next emitted RCF output term as an `RCFTerm`
- a new immutable `GCF` representing the remainder
- an error for exceptional conditions

### 6.8 Finite GCF

A GCF whose underlying source has been exhausted or whose finite emitted prefix has been materialized.

### 6.9 Infinite GCF Source

A source of `(p, q)` terms that is assumed infinite until exhausted.

### 6.10 Convergent

A finite rational approximation induced by a finite prefix of a continued fraction.

### 6.11 Certified Output

An emitted term, interval, comparison result, or other observable result that is mathematically justified by the information consumed so far.

### 6.12 Bound

An endpoint of a `Range`, which may be finite or infinite and may be open or closed.

### 6.13 Range

An object describing the current certified uncertainty set for a `GCF`.

### 6.14 Exact Point

A range denoting a single exact value.

### 6.15 Homographic Transform / Unary LFT

A unary transform of the form:

\[
z(x)=\frac{ax+b}{cx+d}
\]

### 6.16 Bihomographic Transform / Binary LFT

A binary transform of the form:

\[
z(x,y)=\frac{axy+bx+cy+d}{exy+fx+gy+h}
\]

### 6.17 Diagonal LFT

A special case of a Binary LFT where the two operand variables are equal.

---

## 7. Supported Mathematical Objects

### 7.1 Required

- `BigInt`
- `Rational`
- finite `RCF`
- infinite `RCF`
- immutable infinite `GCFSource` values
- immutable `GCF` evaluators with associated `Range`

### 7.2 External Representation Policy

External representation is either:

- `GCFSource` for ingestion
- `RCFTerm` for emission
- `GCF` for arithmetic evaluation and observation

Only GCF sources may be ingested as a series of generalized terms.

Public emitted output shall initially support only ordinary RCF terms and EOF.

### 7.3 Input Assumption Policy

All continued-fraction inputs shall be treated as infinite until exhausted.

Early development shall assume infinite GCF input by default rather than designing first around finite continued fractions.

### 7.4 Optional / Deferred

- decimal digit emission
- arbitrary radix digit emission
- a unary operator that ingests a GCF and emits a `Rational` after consuming a specified number of input terms
- additional specialized named generators beyond the initial constant set
- public non-finite output terms such as `+Inf` or `-Inf`

---

## 8. Representation Requirements

### 8.1 External Representation

The library shall define public representations for:

- `BigInt`
- `Rational`
- immutable `GCFSource`
- immutable `GCF`
- `PQTerm`
- `RCFTerm`
- `Range`
- `Bound`

### 8.2 Internal Representation

The implementation may use any internal representation, provided that:

- mathematical semantics are preserved
- exactness guarantees are not weakened
- testing and diagnostics remain possible

### 8.3 Canonicalization

The specification shall define which forms are canonical and where equivalent non-canonical forms are permitted.

At minimum:

- externally materialized finite RCF values shall follow the chosen regular continued-fraction conventions
- internal transform coefficients shall be normalized by dividing by the `GCD` where appropriate and safe
- normalization shall preserve semantics exactly

### 8.4 Immutable Stepping Requirement

Neither `GCFSource` nor `GCF` shall require in-place public mutation.

Advancing either abstraction shall return a new immutable remainder object.

### 8.5 In-Memory GCF Range Requirement

Every in-memory `GCF` object shall provide a `Range()` function that returns a `Range` object.

The returned `Range` shall represent the current certified uncertainty set for the final value.

---

## 9. Term Requirements

### 9.1 PQTerm Requirements

`PQTerm` shall be a tagged public abstraction that can represent at least:

- ordinary generalized input term
- EOF

Source exhaustion shall therefore be representable as an ordinary `PQTerm` outcome, not only as an `error`.

### 9.2 RCFTerm Requirements

`RCFTerm` shall be a tagged public abstraction that can represent at least:

- ordinary RCF output term
- EOF

Public non-finite output terms are intentionally deferred for the initial API, though the design should not preclude them later.

### 9.3 Error Separation

Ordinary stream outcomes such as value and EOF shall travel in the term objects.

Exceptional conditions shall use the error channel.

---

## 10. Range Requirements

### 10.1 Bound Semantics

A `Bound` shall represent:

- finite or infinite kind
- exact value when finite
- open or closed endpoint status

### 10.2 Range Structure

A `Range` shall contain:

- lower bound `Lo`
- upper bound `Hi`
- explicit boolean `Outside`

### 10.3 Inside / Outside Semantics

A `Range` shall provide `IsInside()` and `IsOutside()` semantics derived from the explicit `Outside` flag.

Outside-ness shall not be encoded indirectly by swapping bound order.

### 10.4 Exactness

A range is exact if and only if it denotes a single exact value.

`Range` shall provide `IsExact()`.

### 10.5 Infinity Support

`Bound` and `Range` shall support positive and negative infinity for range work and comparison.

### 10.6 Comparison Ordering

`Range` shall define a comparison relation suitable for choosing among competing uncertainty descriptions.

The ordering shall be:

1. exact
2. inside narrow
3. inside wide
4. outside narrow
5. outside wide

For ranges of the same class, narrower finite span is better than wider finite span, and any finite span is narrower than an infinite span.

---

## 11. Construction Requirements

### 11.1 Integers

The library shall construct exact immutable `GCFSource` values from integers represented as `BigInt`.

### 11.2 Rationals

The library shall construct exact immutable `GCFSource` values from `Rational` inputs.

### 11.3 Procedural GCF Sources

The library shall accept demand-driven sources that generate `PQTerm` values lazily, using immutable remainder-returning semantics.

### 11.4 Named Constant Sources

The library shall support specialized named generators for important constants, including at least:

- `pi`
- `e`

### 11.5 First Named Infinite Source

Early development shall include one named infinite procedural source for `sqrt(2)`.

This source will serve as the first simple infinite `GCFSource` used for black-box and white-box tests.

### 11.6 Validation

The library shall detect malformed inputs and report errors or invalid-state results according to the error model defined later in this specification.

---

## 12. Core Operation Requirements

### 12.1 Unary Operations

The library shall support, at minimum:

- identity
- negation
- reciprocal
- square root
- sine
- hyperbolic tangent

Deferred unary operations may include:

- absolute value
- additional algebraic operators
- a bounded-ingestion unary operator that emits a `Rational` after consuming a specified number of input terms
- broader transcendental operators

### 12.2 Binary Operations

The library shall support, at minimum:

- addition
- subtraction
- multiplication
- division
- comparison

### 12.3 Exactness

For exact inputs and mathematically defined operations, the library shall not silently degrade to inexact floating-point arithmetic.

### 12.4 Demand-Driven Behavior

Operators shall consume source terms incrementally and only as needed to justify emitted output or decision progress.

### 12.5 Emission Policy

The arithmetic core shall emit tagged `RCFTerm` values rather than bare `BigInt` values.

### 12.6 Division by Zero Semantics

The long-term design shall remain open to representing mathematically certified non-finite results.

The initial public API does not yet expose non-finite emitted `RCFTerm` values, so such cases may remain deferred or be reported through the current error model until that surface is added.

---

## 13. Transform Engine Requirements

The arithmetic core shall be expressible in terms of transform machinery compatible with Gosper-style methods.

### 13.1 Required Major Objects

Early design phases shall include a high-level design describing the major objects and their responsibilities, including at least:

- `GCFSource`
- `GCF`
- `PQTerm`
- `RCFTerm`
- `Rational`
- `Bound`
- `Range`
- `UnaryLFT`
- `BinaryLFT`
- `DiagonalLFT`

### 13.2 Unary Transform Support

The library shall support homographic transforms for unary arithmetic pipelines.

### 13.3 Binary Transform Support

The library shall support bihomographic transforms for binary arithmetic pipelines.

### 13.4 Diagonal LFT Support

The design shall support a `DiagonalLFT` as a degenerate `BinaryLFT` where `X` and `Y` are equal.

### 13.5 Coefficient State

The implementation shall maintain explicit transform state, even if exposed through immutable remainder values rather than mutable objects.

### 13.6 Normalization

Transform coefficients shall be eligible for normalization using `GCD` reduction where such reduction preserves semantics and aids stability, inspection, or debugging.

### 13.7 Identity Initializers

The implementation shall provide identity initializers for transform machinery, including at least:

- ULFT identity: `(1,0)/(0,1)`
- BLFT identity-style initializer: `(1,0,0,0)/(0,0,0,1)`

These identity matrices shall be used to initialize `GCF` evaluation objects where appropriate.

### 13.8 Core Decision Machine

The most important operational requirement in the project is the internal decision machine that decides exactly one of three actions at each binary-evaluation step:

- ingest from left
- ingest from right
- emit output

Nothing else is a primary binary evaluation action.

This decision machine shall be explicit in the implementation and testable through white-box tests.

### 13.9 Degenerate States

The implementation shall define behavior for singular, degenerate, or otherwise undefined transform states.

---

## 14. Output Requirements

### 14.1 Continued-Fraction Terms

The library shall be able to emit result RCF terms incrementally.

### 14.2 Output Signals

The emitted output term abstraction shall initially represent at least:

- ordinary RCF term
- EOF

### 14.3 Finite Materialization

The library shall be able to materialize finite prefixes.

### 14.4 Prefix Extraction

The API shall support extracting a finite prefix from an evaluator via `Take(n)` or equivalent semantics.

`Take(n)` shall return:

- a finite `GCF` containing up to `n` emitted `RCFTerm` values
- an error which may be `io.EOF` if the evaluator ended before `n` terms were available

It shall not return the remaining evaluator.

### 14.5 Convergents

The library shall be able to produce convergents from finite evaluators, especially those returned by `Take(n)`.

### 14.6 Bounds / Ranges

The library shall be able to report current certified `Range` information.

### 14.7 Future Digit / Radix Output

Emission of decimal digits or digits in other radices is a long-term goal and shall remain architecturally feasible.

---

## 15. Correctness Requirements

### 15.1 Term Correctness

Every emitted RCF term shall be mathematically justified by the already-consumed operand information and the valid unread-tail assumptions defined by the model.

### 15.2 Finite Input Correctness

For inputs that become finite through source exhaustion, the library shall produce the exact result when the operation is defined.

### 15.3 Equivalence

Equivalent representations of the same value shall behave equivalently under supported operations, modulo canonicalization policy.

### 15.4 Comparison Correctness

Comparison results shall not be reported unless they are mathematically justified.

### 15.5 No Silent Corruption

The library shall not emit known-incorrect terms or silently substitute approximate arithmetic in exact modes.

---

## 16. Progress and Termination Requirements

### 16.1 Infinite-First Assumption

Early development shall assume infinite GCF inputs unless and until a source is exhausted.

### 16.2 Finite by Exhaustion

A source that becomes exhausted shall thereafter be treated as finite, and downstream logic shall handle the transition correctly.

### 16.3 Infinite Inputs

For infinite inputs, the library shall support ongoing incremental progress when mathematically possible.

### 16.4 Blocked States

The implementation shall define what it means internally for an operator to be blocked waiting for more input.

Ordinary blocked waiting shall not normally surface as a public emission result.

### 16.5 Stalled States

The library shall define what it means for evaluation to stall or fail to make progress.

### 16.6 Observability

The caller and tests shall be able to distinguish among:

- emitted output term
- EOF
- mathematically undefined state
- implementation-detected stuck or non-progress state

### 16.7 Bounded Work Modes

Optional bounded-step or bounded-resource execution modes may be provided for diagnostics and testing.

---

## 17. Error and Exceptional Behavior

The specification shall define the error model for:

- malformed input terms
- invalid generalized terms
- singular transforms
- exhausted finite sources in illegal contexts
- source-protocol violations
- internal invariant failures
- implementation-detected stuck / non-progress states

Pre-condition and post-condition checks for invariants are explicitly permitted and encouraged when they improve debugging and correctness assurance.

The exact split among ordinary errors, test-time assertions, and internal-fault conditions shall be defined by the API and design documents.

---

## 18. Go Style and Implementation Conventions

The implementation shall conform to generally recognized Go style guides so that an experienced Go programmer feels at home in the codebase.

At minimum, the project shall aim to follow:

- idiomatic Go naming conventions
- short, lower-case package names
- standard export conventions
- `gofmt` formatting
- idiomatic import grouping
- idiomatic receiver naming
- common Go error-handling patterns

The project shall prefer the simplest Go design that preserves mathematical correctness and testability.

---

## 19. Diagnostics and Introspection

The library shall provide observability sufficient for testing and debugging.

### 19.1 Required Diagnostic Capabilities

- inspect current operator state
- inspect current transform coefficients or equivalent
- inspect emitted-result prefix
- inspect current convergents and/or ranges
- inspect why output is not currently possible
- inspect whether a state is finite or infinite where meaningful

### 19.2 Trace Support

Optional tracing should allow step-by-step examination of:

- input decisions
- output decisions
- transform rewrites
- simplifications
- range updates

### 19.3 Invariant Checking

Debug and test modes may enforce stronger internal invariant checks.

### 19.4 Non-Progress Testing Support

Test cases involving implementation-detected stuck or non-progress states shall include an expiration timer or equivalent bounded termination mechanism.

---

## 20. API Requirements

### 20.1 Public API Goal

An early phase of development shall define the public API, and that API shall be as small as possible while still being sufficient to compute the target formula.

### 20.2 Public API Style

The library shall expose a clean public API with clear separation between:

- immutable source construction
- immutable arithmetic evaluation
- diagnostics/testing hooks
- lower-level transform machinery

### 20.3 Immutable Stepping API

A caller shall be able to advance both `GCFSource` and `GCF` without mutating the original value.

Each step shall return a remainder object.

### 20.4 Prefix Extraction API

The API shall support extracting a finite prefix from an evaluator via `Take(n)` semantics that return only the finite prefix and an error.

### 20.5 Convergent API

`Convergent()` shall produce a `Rational` from a finite evaluator, especially one returned by `Take(n)`.

### 20.6 Determinism

Given the same inputs and the same evaluation strategy, the library shall behave deterministically.

### 20.7 Concurrency

Thread-safety requirements are deferred.

### 20.8 Cloning

The public API shall not require cloning.

If sharing or memoization is later useful, it shall be treated as an internal optimization rather than a public requirement.

---

## 21. Testing Requirements

The design shall support iterative TDD and precise verification.

### 21.1 TDD Workflow

Development shall normally proceed in two phases:

1. create failing tests, allowing production stubs that return an incorrect value of the correct type
2. modify production code until those tests pass

The expected workflow is red, then green, then commit, then repeat.

### 21.2 Public vs Private Testing

The public interface shall be tested with black-box tests.

Private interfaces and internal machinery shall be tested with white-box tests.

### 21.3 Package-Level Test Access

In the initial Go implementation, white-box tests will live in the same package and therefore may access private functions and data as needed.

### 21.4 BB / WB Test Convention

The project shall adopt a test-file naming convention that distinguishes black-box and white-box tests.

Recommended convention:

- black-box files: `*_bb_test.go`
- white-box files: `*_wb_test.go`

Recommended package convention:

- black-box tests in package `cf_test`
- white-box tests in package `cf`

Recommended test-name convention:

- black-box tests prefixed `TestBB_`
- white-box tests prefixed `TestWB_`

This convention should make it easy to run BB and WB tests separately.

### 21.5 Pending Tests

Pending or intentionally not-yet-green tests shall not be kept in separate source directories.

Instead, they should use opt-in build tags such as `//go:build pending`.

### 21.6 Unit Tests

The library shall support unit testing of:

- source ingestion
- emitted term stepping
- transform updates
- canonicalization
- range behavior
- error handling

### 21.7 Property Tests

Where practical, operations shall be testable against rational arithmetic or equivalent reference models.

### 21.8 Golden Tests

The test suite shall include at least:

- test cases from Gosper’s articles
- additional tests derived from newer insights about GCF behavior
- regression tests for discovered edge cases

### 21.9 Stall / Progress Regression Tests

The library shall support regression tests for blocked, stalled, or historically problematic evaluation paths.

### 21.10 Inspection Hooks

The API and implementation shall not hide essential state needed for correctness-oriented testing.

---

## 22. Documentation Requirements

The project documentation shall include:

- mathematical overview
- requirements specification
- high-level design
- user guide
- API reference
- examples for RCF and GCF usage
- examples of exact arithmetic
- examples of streaming term production
- explanation of guarantees and limitations
- glossary of terminology

---

## 23. Compatibility and Portability

This specification is intended to be language-aware but mathematically portable.

Language-specific implementations may differ in:

- naming
- packaging
- iterator/stream conventions
- error-handling idioms
- numeric backend libraries

However, implementations shall preserve the required mathematical behavior described here.

The initial implementation is expected to target Go.

---

## 24. Robustness and Security Considerations

Robustness is important and includes:

- bounded non-progress testing
- invariant checks
- handling malformed source behavior
- handling pathological coefficient growth
- handling exhausted sources correctly

Security concerns are deferred in the early phases and are not a first-pass requirement.

---

## 25. Versioning and Evolution

The specification shall support staged delivery.

### 25.1 Phase 1

- requirements specification
- high-level design
- major object model
- minimal public API design
- `BigInt`, `Rational`, `Bound`, `Range`
- `GCFSource`, `GCF`, `PQTerm`, `RCFTerm`
- identity-transform startup path

### 25.2 Phase 2

- red/green TDD harness
- exact source constructors
- infinite-GCF ingestion
- immutable stepping semantics
- RCF emission
- convergents
- diagnostics foundation

### 25.3 Phase 3

- exact binary arithmetic core
- exact unary arithmetic substrate
- finite-by-exhaustion correctness
- stronger progress diagnostics

### 25.4 Phase 4

- specialized constant generators including `pi` and `e`
- named infinite `sqrt(2)` source
- target-formula unary operators: `sqrt`, `sin`, `tanh`
- minimal public API sufficient to compute the target formula

### 25.5 Phase 5

- improved certification and range behavior
- bounded rational-collapse unary operator
- decimal/radix emission
- broader unary and transcendental support
- possible public non-finite output terms if needed

---

## 26. Acceptance Criteria for First Implementable Milestone

The first serious milestone shall be considered complete when all of the following are true:

1. The requirements specification and high-level design exist and are coherent.
2. The public API exists in minimal form and is sufficient in principle to express the target formula.
3. The library can construct immutable `GCFSource` values from `BigInt` and `Rational`.
4. The library can lift `GCFSource` values into immutable `GCF` evaluators.
5. The library can ingest infinite GCF sources of generalized terms.
6. The library can emit tagged RCF output terms incrementally.
7. The library exposes `Range` semantics for in-memory `GCF` values.
8. The implementation exposes enough diagnostic state to test and explain blocked or stalled behavior.
9. The implementation supports red/green TDD with white-box access to private machinery.
10. The implementation is backed by automated tests including Gosper-derived examples and newer GCF regression cases.
11. The implementation includes one named infinite procedural source for `sqrt(2)`.

---

## 27. Acceptance Criteria for Major Project Milestone

A major milestone shall be considered complete when the library can compute the target formula:

\[
\frac{\sqrt{\frac{3}{\pi^2} + e}}{\tanh(\sqrt{5}) - \sin(69^\circ)}
\]

with exact continued-fraction machinery, using named constant generators and the required unary and binary operators, without silently degrading to floating-point arithmetic.

---

## 28. Placeholder Appendices

### Appendix A: Mathematical Conventions

TBD

### Appendix B: Canonical Examples

TBD

### Appendix C: Error Taxonomy

TBD

### Appendix D: Public API Sketch

TBD

### Appendix E: Test Matrix

TBD
<!-- RequirementsSpec.md v3 -->