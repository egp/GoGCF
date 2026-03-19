<!-- RequirementsSpec.md v2 -->
# RequirementsSpec.md

# Continued Fraction Arithmetic Library Requirements Specification

## 1. Purpose

This document defines the requirements for a library that performs arithmetic on continued fractions using Gosper-style methods.

The library is intended to support mathematically correct, testable, demand-driven arithmetic on continued fractions, especially in the spirit of HAKMEM items 101A and 101B.

This specification is focused on requirements, not implementation details. It describes what the library must do, what guarantees it must provide, and what behaviors must be observable by callers and tests.

The Requirements Specification is the north-pole document for design and API decisions. It is expected to evolve as new questions arise.

---

## 2. Background and Context

Continued fractions provide an alternative numeric representation with useful arithmetic properties, especially for exact and incremental computation.

The target library is intended to support:

- exact arithmetic on continued fractions
- demand-driven evaluation
- regular continued fractions (RCF) and generalized continued fractions (GCF)
- composable unary and binary arithmetic operators
- certified partial output when full evaluation is not yet available

The initial architecture should be aligned with Gosper-style homographic and bihomographic transformation machinery, while remaining testable and implementation-language-agnostic at the requirements level.

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
- construction of immutable `GCF` values from `GCFSource`
- demand-driven procedural infinite GCF sources
- exact comparison where possible
- unary and binary arithmetic over continued fractions
- demand-driven RCF term emission
- convergent and interval/range inspection
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

This section defines the terms used throughout the specification.

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

- the next generalized input term
- a new immutable `GCFSource` representing the remainder
- an error, including `io.EOF` when exhausted

### 6.7 GCF

An immutable arithmetic evaluator/value over continued fractions.

Calling `Next()` on a `GCF` returns:

- the next emitted RCF output term or out-of-band output signal
- a new immutable `GCF` representing the remainder
- an error for exceptional conditions

### 6.8 Finite GCF

A GCF whose underlying source has been exhausted.

### 6.9 Infinite GCF Source

A source of `(p, q)` terms that is assumed infinite until exhausted.

### 6.10 Convergent

A finite rational approximation induced by a finite prefix of a continued fraction.

### 6.11 Certified Output

An emitted term, interval, comparison result, or other observable result that is mathematically justified by the information consumed so far.

### 6.12 Range

An interval-like object describing uncertainty about the final value of a GCF.

### 6.13 Inside Range

A range where the value lies inside the interval or set described by the bounds.

### 6.14 Outside Range

A range where the value lies outside the interval or set described by the bounds.

### 6.15 Exact Point

A range denoting a single exact value.

### 6.16 Homographic Transform / Unary LFT

A unary transform of the form:

\[
z(x)=\frac{ax+b}{cx+d}
\]

### 6.17 Bihomographic Transform / Binary LFT

A binary transform of the form:

\[
z(x,y)=\frac{axy+bx+cy+d}{exy+fx+gy+h}
\]

### 6.18 Diagonal LFT

A special case of a Binary LFT where the two operand variables are equal.

---

## 7. Supported Mathematical Objects

The library shall support the following categories of values.

### 7.1 Required

- `BigInt`
- `Rational`
- finite `RCF`
- infinite `RCF`
- immutable infinite `GCFSource` values
- immutable `GCF` values with associated `Range`

### 7.2 External Representation Policy

External representation is either:

- `GCFSource` for ingestion
- `RCFTerm` for emission
- `GCF` for arithmetic evaluation and observation

Only infinite GCF sources may be ingested as a series of `(p, q)` terms.

Only RCF output terms or explicit out-of-band output signals shall be emitted by the arithmetic core.

### 7.3 Input Assumption Policy

All continued-fraction inputs shall be treated as infinite until exhausted.

Early development shall assume infinite GCF input by default rather than designing first around finite continued fractions.

### 7.4 Optional / Deferred

- decimal digit emission
- arbitrary radix digit emission
- a unary operator that ingests a GCF and emits a `Rational` after consuming a specified number of input terms
- additional specialized named generators beyond the initial constant set

---

## 8. Representation Requirements

### 8.1 External Representation

The library shall define public representations for:

- `BigInt`
- `Rational`
- immutable `GCFSource`
- immutable `GCF`
- emitted `RCFTerm`
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

- emitted RCF terms shall follow the chosen regular continued-fraction conventions
- positive and negative infinity shall be treated separately
- internal transform coefficients shall be normalized by dividing by the `GCD` where appropriate and safe
- normalization shall preserve semantics exactly

### 8.4 Immutable Stepping Requirement

Neither `GCFSource` nor `GCF` shall require in-place public mutation.

Advancing either abstraction shall return a new immutable remainder object.

### 8.5 In-Memory GCF Range Requirement

Every in-memory `GCF` object shall provide a `Range()` function that returns a `Range` object.

The returned `Range` shall represent the uncertainty interval or uncertainty set containing, or excluding, the actual final value of the GCF.

---

## 9. Range Requirements

### 9.1 Bound Semantics

A `Bound` shall represent:

- finite or infinite kind
- exact value when finite
- open or closed endpoint status

### 9.2 Range Semantics

A `Range` shall contain two bounds, `Lo` and `Hi`.

### 9.3 Inside / Outside Semantics

A `Range` shall provide `IsInside()` and `IsOutside()` semantics derived from bound ordering together with open/closed endpoint conventions.

Inside/outside shall not require a separate stored flag in the public abstraction.

### 9.4 Exactness

A range is exact if and only if it denotes a single exact value.

Examples:

- `[x, x]` is exact
- `(x, x)` is not exact

### 9.5 Infinity Support

`Bound` and `Range` shall support positive and negative infinity.

Positive and negative infinity are also useful in range comparisons and in arithmetic results.

### 9.6 Comparison Ordering

`Range` shall define a comparison relation suitable for choosing among competing uncertainty descriptions.

The intended ordering is:

1. inside narrow
2. inside wide
3. outside wide
4. outside narrow

The exact formalization of “narrow” and “wide” shall be defined in the design and API documents, but the ordering semantics above are required.

---

## 10. Construction Requirements

The library shall provide ways to construct continued-fraction objects from the following inputs.

### 10.1 Integers

The library shall construct exact immutable `GCFSource` values from integers represented as `BigInt`.

### 10.2 Rationals

The library shall construct exact immutable `GCFSource` values from `Rational` inputs.

### 10.3 Procedural GCF Sources

The library shall accept demand-driven sources that generate `(p, q)` terms lazily, using immutable remainder-returning semantics.

### 10.4 Named Constant Sources

The library shall support specialized named generators for important constants, including at least:

- `pi`
- `e`

### 10.5 First Named Infinite Source

Early development shall include one named infinite procedural source for `sqrt(2)`.

This source will serve as the first simple infinite `GCFSource` used for black-box and white-box tests.

### 10.6 Validation

The library shall detect malformed inputs and report errors or invalid-state results according to the error model defined later in this specification.

---

## 11. Core Operation Requirements

### 11.1 Unary Operations

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

### 11.2 Binary Operations

The library shall support, at minimum:

- addition
- subtraction
- multiplication
- division
- comparison

### 11.3 Exactness

For exact inputs and mathematically defined operations, the library shall not silently degrade to inexact floating-point arithmetic.

### 11.4 Demand-Driven Behavior

Operators shall consume source terms incrementally and only as needed to justify emitted output or decision progress.

### 11.5 Emission Policy

The arithmetic core shall emit tagged `RCFTerm` values rather than bare `BigInt` values.

### 11.6 Division by Zero Semantics

Division by exact zero is not automatically an error.

When mathematically certified, division by exact zero may emit positive or negative infinity rather than fail.

Truly undefined forms such as `0/0` remain errors.

---

## 12. Transform Engine Requirements

The arithmetic core shall be expressible in terms of transform machinery compatible with Gosper-style methods.

### 12.1 Required Major Objects

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

### 12.2 Unary Transform Support

The library shall support homographic transforms for unary arithmetic pipelines.

### 12.3 Binary Transform Support

The library shall support bihomographic transforms for binary arithmetic pipelines.

### 12.4 Diagonal LFT Support

The design shall support a `DiagonalLFT` as a degenerate `BinaryLFT` where `X` and `Y` are equal.

### 12.5 Coefficient State

The implementation shall maintain explicit transform state, even if exposed through immutable remainder values rather than mutable objects.

### 12.6 Normalization

Transform coefficients shall be eligible for normalization using `GCD` reduction where such reduction preserves semantics and aids stability, inspection, or debugging.

### 12.7 Initial Development Starting Point

Early development shall begin from identity transforms, including at least:

- Unary LFT identity: `(1,0)/(0,1)`
- Binary LFT identity-style initial form: `(1,0,0,0)/(0,0,0,1)`

### 12.8 Degenerate States

The implementation shall define behavior for singular, degenerate, or otherwise undefined transform states.

---

## 13. Output Requirements

The library shall support one or more of the following output forms.

### 13.1 Continued-Fraction Terms

The library shall be able to emit result RCF terms incrementally.

### 13.2 Output Signals

The emitted output term abstraction shall be able to represent at least:

- ordinary RCF term
- EOF
- positive infinity
- negative infinity

### 13.3 Finite Materialization

The library shall be able to materialize finite prefixes.

### 13.4 Convergents

The library shall be able to produce convergents from finite evaluators, especially those returned by `Take(n)`.

### 13.5 Bounds / Ranges

The library shall be able to report current certified `Range` information.

### 13.6 Future Digit / Radix Output

Emission of decimal digits or digits in other radices is a long-term goal and shall remain architecturally feasible.

---

## 14. Correctness Requirements

### 14.1 Term Correctness

Every emitted RCF term shall be mathematically justified by the already-consumed operand information and the valid unread-tail assumptions defined by the model.

### 14.2 Finite Input Correctness

For inputs that become finite through source exhaustion, the library shall produce the exact result when the operation is defined.

### 14.3 Equivalence

Equivalent representations of the same value shall behave equivalently under supported operations, modulo canonicalization policy.

### 14.4 Comparison Correctness

Comparison results shall not be reported unless they are mathematically justified.

### 14.5 No Silent Corruption

The library shall not emit known-incorrect terms or silently substitute approximate arithmetic in exact modes.

### 14.6 Infinity Handling

Positive infinity and negative infinity shall be modeled distinctly and shall not be conflated.

---

## 15. Progress and Termination Requirements

### 15.1 Infinite-First Assumption

Early development shall assume infinite GCF inputs unless and until a source is exhausted.

### 15.2 Finite by Exhaustion

A source that becomes exhausted shall thereafter be treated as finite, and downstream logic shall handle the transition correctly.

### 15.3 Infinite Inputs

For infinite inputs, the library shall support ongoing incremental progress when mathematically possible.

### 15.4 Blocked States

The implementation shall define what it means internally for an operator to be blocked waiting for more input.

Ordinary blocked waiting shall not normally surface as a public emission result.

### 15.5 Stalled States

The library shall define what it means for evaluation to stall or fail to make progress.

### 15.6 Observability

The caller and tests shall be able to distinguish among:

- emitted output term
- EOF
- positive infinity
- negative infinity
- mathematically undefined state
- implementation-detected stuck or non-progress state

### 15.7 Bounded Work Modes

Optional bounded-step or bounded-resource execution modes may be provided for diagnostics and testing.

---

## 16. Error and Exceptional Behavior

The specification shall define the error model for:

- malformed input terms
- invalid generalized terms
- truly undefined forms such as `0/0`
- singular transforms
- exhausted finite sources in illegal contexts
- source-protocol violations
- internal invariant failures
- implementation-detected stuck / non-progress states

Pre-condition and post-condition checks for invariants are explicitly permitted and encouraged when they improve debugging and correctness assurance.

The exact split among ordinary errors, test-time assertions, and internal-fault conditions shall be defined by the API and design documents.

---

## 17. Diagnostics and Introspection

The library shall provide observability sufficient for testing and debugging.

### 17.1 Required Diagnostic Capabilities

- inspect current operator state
- inspect current transform coefficients or equivalent
- inspect emitted-result prefix
- inspect current convergents and/or ranges
- inspect why output is not currently possible
- inspect whether a state is finite or infinite where meaningful

### 17.2 Trace Support

Optional tracing should allow step-by-step examination of:

- input decisions
- output decisions
- transform rewrites
- simplifications
- range updates

### 17.3 Invariant Checking

Debug and test modes may enforce stronger internal invariant checks.

### 17.4 Non-Progress Testing Support

Test cases involving implementation-detected stuck or non-progress states shall include an expiration timer or equivalent bounded termination mechanism.

---

## 18. API Requirements

### 18.1 Public API Goal

An early phase of development shall define the public API, and that API shall be as small as possible while still being sufficient to compute the target formula.

### 18.2 Public API Style

The library shall expose a clean public API with clear separation between:

- immutable source construction
- immutable arithmetic evaluation
- diagnostics/testing hooks
- lower-level transform machinery

### 18.3 Immutable Stepping API

A caller shall be able to advance both `GCFSource` and `GCF` without mutating the original value.

Each step shall return a remainder object.

### 18.4 Prefix Extraction API

The API shall support extracting a finite prefix from an evaluator via `Take(n)` or equivalent semantics, returning both the finite prefix and the remaining evaluator.

### 18.5 Convergent API

`Convergent()` shall produce a `Rational` from a finite evaluator, especially one returned by `Take(n)`.

### 18.6 Determinism

Given the same inputs and the same evaluation strategy, the library shall behave deterministically.

### 18.7 Concurrency

Thread-safety requirements are deferred.

---

## 19. Performance Requirements

### 19.1 Priority

Performance is secondary to correctness, testability, and debuggability.

### 19.2 Arbitrary Precision

The library shall support arbitrary-precision integer arithmetic where needed to preserve exactness.

### 19.3 Laziness

The implementation shall avoid unnecessary full collapse to rationals when incremental progress is sufficient.

### 19.4 Resource Growth

The implementation shall expose or document possible growth in:

- coefficient size
- memory usage
- work per emitted term

### 19.5 Safeguards

Optional safeguards may include:

- maximum steps
- maximum coefficient size
- maximum source pulls
- maximum emitted terms
- timeout support in tests

---

## 20. Testability Requirements

The design shall support iterative TDD and precise verification.

### 20.1 TDD Workflow

Development shall normally proceed in two phases:

1. create failing tests, allowing production stubs that return an incorrect value of the correct type
2. modify production code until those tests pass

The expected workflow is red, then green, then commit, then repeat.

### 20.2 Public vs Private Testing

The public interface shall be tested with black-box tests.

Private interfaces and internal machinery shall be tested with white-box tests.

### 20.3 Package-Level Test Access

In the initial Go implementation, tests will live in the same package and therefore may access private functions and data as needed for white-box verification.

### 20.4 Unit Tests

The library shall support unit testing of:

- source ingestion
- emitted term stepping
- transform updates
- canonicalization
- range behavior
- error handling

### 20.5 Property Tests

Where practical, operations shall be testable against rational arithmetic or equivalent reference models.

### 20.6 Golden Tests

The test suite shall include at least:

- test cases from Gosper’s articles
- additional tests derived from newer insights about GCF behavior
- regression tests for discovered edge cases

### 20.7 Stall / Progress Regression Tests

The library shall support regression tests for blocked, stalled, or historically problematic evaluation paths.

### 20.8 Inspection Hooks

The API and implementation shall not hide essential state needed for correctness-oriented testing.

---

## 21. Documentation Requirements

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

## 22. Compatibility and Portability

This specification is intended to be language-agnostic at the mathematical level.

Language-specific implementations may differ in:

- naming
- packaging
- iterator/stream conventions
- error-handling idioms
- numeric backend libraries

However, implementations shall preserve the required mathematical behavior described here.

The initial implementation is expected to target Go.

---

## 23. Robustness and Security Considerations

Robustness is important and includes:

- bounded non-progress testing
- invariant checks
- handling malformed source behavior
- handling pathological coefficient growth
- handling exhausted sources correctly

Security concerns are deferred in the early phases and are not a first-pass requirement.

---

## 24. Versioning and Evolution

The specification shall support staged delivery.

### 24.1 Phase 1

- requirements specification
- high-level design
- major object model
- minimal public API design
- `BigInt`, `Rational`, `Bound`, `Range`
- `GCFSource`, `GCF`, `PQTerm`, `RCFTerm`
- identity-transform startup path

### 24.2 Phase 2

- red/green TDD harness
- exact source constructors
- infinite-GCF ingestion
- immutable stepping semantics
- RCF emission
- convergents
- diagnostics foundation

### 24.3 Phase 3

- exact binary arithmetic core
- exact unary arithmetic substrate
- finite-by-exhaustion correctness
- stronger progress diagnostics

### 24.4 Phase 4

- specialized constant generators including `pi` and `e`
- named infinite `sqrt(2)` source
- target-formula unary operators: `sqrt`, `sin`, `tanh`
- minimal public API sufficient to compute the target formula

### 24.5 Phase 5

- improved certification and range behavior
- bounded rational-collapse unary operator
- decimal/radix emission
- broader unary and transcendental support

---

## 25. Open Design Questions

The following issues remain to be decided.

### 25.1 Canonicalization Policy

How aggressively should finite and emitted RCFs be normalized?

### 25.2 Output Signal Surface

Should `RCFTerm` eventually represent additional out-of-band signals beyond EOF and `±Inf`?

### 25.3 Progress Strategy

What policy should govern source selection and output attempts in ambiguous operator states?

### 25.4 Range Ordering Formalization

How exactly should “inside narrow”, “inside wide”, “outside wide”, and “outside narrow” be measured and compared?

### 25.5 Stuck-State Semantics

How should the library report cases where evaluation appears to fail to make progress?

### 25.6 Finite Prefix Metadata

Should `Take(n)` preserve metadata distinguishing termination by count limit from early source exhaustion beyond the returned error?

### 25.7 Bounded Rational Collapse

What should be the exact semantics of the future unary operator that consumes a specified number of input terms and emits a `Rational`?

---

## 26. Acceptance Criteria for First Implementable Milestone

The first serious milestone shall be considered complete when all of the following are true:

1. The requirements specification and high-level design exist and are coherent.
2. The public API exists in minimal form and is sufficient in principle to express the target formula.
3. The library can construct immutable `GCFSource` values from `BigInt` and `Rational`.
4. The library can lift `GCFSource` values into immutable `GCF` evaluators.
5. The library can ingest infinite GCF sources of `(p, q)` terms.
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
<!-- RequirementsSpec.md v2 -->