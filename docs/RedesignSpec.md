# Spec.md

# goGCF / core Specification

## Status

This document is the current design contract for the `core` package and the planned companion package `named`.

When a design choice conflicts with Gosper / HAKMEM 101B, Gosper wins. The generalized-term convention used throughout is:

\[
X = p + \frac{q}{X'}
\]

where `(p,q)` is a generalized continued-fraction term and `X'` is the tail. Gosper explicitly treats `(p,q)` as the generalized term shape, with regular continued fractions as the `q = 1` special case. :contentReference[oaicite:0]{index=0}

The implementation uses exact arithmetic only. No floats anywhere in the code. Go `math/big` provides arbitrary-precision integers and rationals, and `Int.BitLen()` is available for configurable test/runtime limits. :contentReference[oaicite:1]{index=1}

---

## 1. Mission

Build a mathematically correct, testable Go library for arithmetic on generalized continued fractions in the style of Gosper / HAKMEM 101B.

Primary goals:

- mathematical correctness
- simplicity
- exact arithmetic only
- design for test
- top-down and bottom-up development together
- red/green TDD with realistic wrong-value stubs rather than “NYI” panics, except where the API explicitly reserves an NYI stub

Performance matters only after correctness and clarity.

---

## 2. Package layout

### 2.1 `core`
The `core` package contains:

- generalized source abstractions
- regular CF output stream abstractions
- BLFT-based transform engine
- range logic
- normalization/invariant enforcement
- client-facing arithmetic façade
- config / limits / policy

### 2.2 `named`
The companion package is `named`.

`named` contains:

- constructors from constants and `Rational`
- constructors from finite regular-CF streams
- named sources such as `Pi()`, `E()`, `Sqrt2()`, etc.
- other source generators that are not part of the minimal engine core

`named` does not weaken the `core` rule that engine ingestion is always in generalized `(p,q)` terms.

---

## 3. Core design principles

### 3.1 Exact arithmetic only
No floats anywhere in production code, tests, or helpers used for correctness decisions.

### 3.2 Simplicity
Prefer one clear algebraic path over multiple convenience paths.

### 3.3 DRY
ULFT, BLFT, and DLFT behavior should share helpers wherever algebraically sound.

### 3.4 One internal engine shape
Internally, use one BLFT-sized transform state whenever practical. Unary and diagonal cases are specialized views / constrained initializations of the same machinery.

### 3.5 Certified output only
The engine emits a regular continued-fraction term only when Gosper’s emission criterion says it is forced by the current range.

### 3.6 EOF is valid
EOF is not an error. EOF triggers algebraic collapse and computation continues if possible.

### 3.7 Invalid client input is reported
Malformed source terms coming from client-provided or companion-provided streams return an error/status rather than pretending to be EOF.

### 3.8 Internal invariant failure is fatal
Internal impossible states panic.

---

## 4. Mathematical conventions

## 4.1 Generalized source term convention
Every generalized source term uses:

\[
X = p + \frac{q}{X'}
\]

This is the exact convention everywhere in core.

## 4.2 Regular emitted output convention
Emitted output is a regular continued fraction term sequence.

- the first emitted regular term may be negative
- all later emitted regular terms are standard regular-CF terms

## 4.3 Input sign policy for generalized source terms
For v1, the project adopts this practical source-well-formedness rule:

- `p` may be negative on any term
- `q` must be nonzero on every term
- the first term may have negative `q`
- after the first term, negative `q` is invalid input

This is a source validation policy, not a claim that Gosper algebra would fail on signed later `q`. It is chosen to simplify v1 reasoning and debugging.

## 4.4 No internal `q == 1` shortcut path
Regular CFs are represented as a special case of generalized terms, but core logic does not fork to a distinct “regular only” arithmetic path.

---

## 5. Numeric types

## 5.1 Big integers
Use `*big.Int` internally for:

- `PQTerm.p`
- `PQTerm.q`
- `RCFTerm`
- BLFT coefficients
- Rational numerator/denominator

## 5.2 Rational
`Rational` is exact and BigInt-based.

Expected operations include:

- normalization
- comparison
- sign
- arithmetic needed for endpoint evaluation
- exact conversion from integer

## 5.3 Int64 usage
Use `int64` only where overflow is impossible by construction.

No correctness-critical internal algorithm may rely on `int64` if overflow is possible.

---

## 6. Conceptual architecture

There are three conceptual pieces.

## 6.1 `PQSource`
A generalized source stream.

Responsibilities:

- provide the next generalized term `(p,q)`
- return the tail source `X'`
- provide a range enclosure for the current source value

Conceptual contract:

- `Next()` returns `(PQTerm, tail, ok, err)`
- `ok=false` means EOF
- EOF is valid and not an error
- malformed client-visible source data returns `err != nil`
- live sources provide `Range()`

## 6.2 `RCFStream`
A regular continued-fraction output stream.

Responsibilities:

- provide the next regular term
- return the tail stream `Z'`
- provide a range enclosure for the current value

Conceptual contract:

- `NextRCF()` returns `(RCFTerm, tail, ok, err)`
- `ok=false` means no more emitted terms
- malformed client-visible source data or configured policy breach may return error
- live streams provide `Range()`

## 6.3 `GCF`
A client-facing arithmetic façade / pipeline object.

A `GCF` is created from:

- an initial transform/plugboard configuration
- zero, one, or two `PQSource` inputs
- config/policy

The client calls `.Next()` on the created `GCF` to obtain regular output terms.

Interpretation:

- sources feed the engine
- the engine ingests, emits, or collapses
- the client sees a stream of regular CF terms

### Note
This document treats `PQSource`, `RCFStream`, and `GCF` as conceptual roles. Whether they become separate public interfaces, structs, or a smaller public API with private implementations remains an open question.

---

## 7. Public API direction

This section records current intent, not final Go signatures.

## 7.1 Source-side direction

- `Next() (PQTerm, PQSource, bool, error)`
- `Range() Range`

Meaning:

- `bool=false` => EOF
- `error!=nil` => malformed source or other client-visible source error
- nil receiver use is invalid

## 7.2 Output-side direction

- `NextRCF() (RCFTerm, RCFStream, bool, error)`
- `Range() Range`

Meaning:

- `bool=false` => no more output terms
- `error!=nil` => client-visible error or policy breach

## 7.3 Client-facing `GCF`
The client creates a `GCF` with:

- an initial transform configuration
- zero, one, or two input streams
- a config

Then calls:

- `Next()`
- `Range()`
- `Cmp(other GCF)` or equivalent
- `NextDecimalDigit()`

### `NextDecimalDigit()`
This API exists now as a reserved stub.

Behavior for v1:

- returns `(0, ErrNYI)`

This is intentional and explicit.

## 7.4 `HasNext()`
No `HasNext()` API is part of v1.

Reason:
determining whether a next term exists may require actual engine work, ingestion, collapse, or emission logic. A separate `HasNext()` risks duplicating logic or creating misleading semantics.

---

## 8. Range model

## 8.1 Range kinds
A `Range` is one of:

- `InsideInterval`
- `OutsideInterval`

## 8.2 Endpoint type
Endpoints are exact `Rational`s.

Each endpoint has its own openness/closedness flag.

## 8.3 Inside interval meaning
An inside interval encloses values within bounded endpoints.

Recommended invariant:

- `lo <= hi`

Interpretation:

- value lies in the interval bounded by `lo` and `hi`
- if `lo == hi`, the value is known exactly

## 8.4 Outside interval meaning
An outside interval denotes the complement-style enclosure:

\[
(-\infty, lo] \cup [hi, +\infty)
\]

with independently open/closed finite endpoints.

Exact canonical ordering for outside intervals is still open. See Open Questions.

## 8.5 Range quality ordering
Range comparison uses the Gosper-aligned quality order:

- inside narrow
- inside wide
- outside wide
- outside narrow

from best/smallest to worst/widest.

This drives ingest decisions.

## 8.6 Range() validity
`Range()` is defined only on live stream/engine values.

No error return is planned for `Range()`.

Constructors return errors.
Internal misuse panics.
EOF is represented by `ok=false` and a nil tail, not by a broken live object.

---

## 9. BLFT core engine

## 9.1 Primary internal state
The primary internal transform state is BLFT-shaped:

\[
Z(X,Y) = \frac{aXY + bX + cY + d}{eXY + fX + gY + h}
\]

with BigInt coefficients.

ULFT and DLFT are represented as constrained/specialized cases of the same internal machinery where practical.

## 9.2 Core engine actions
At each step the engine does exactly one of:

1. ingest one `(p,q)` term from one input
2. emit one regular CF term
3. collapse on EOF from an input

The difficult decision is between ingest and emit. Collapse is triggered by EOF and is not an error.

---

## 10. Gosper ingest formulas

If `X` produces `(p,q)` with

\[
X = p + \frac{q}{X'}
\]

then ingesting `X` updates BLFT coefficients to:

\[
(a,b,c,d;e,f,g,h)
\mapsto
(pa+c,\ pb+d,\ qa,\ qb;\ pe+g,\ pf+h,\ qe,\ qf)
\]

If `Y` produces `(r,s)` with

\[
Y = r + \frac{s}{Y'}
\]

then ingesting `Y` updates BLFT coefficients to:

\[
(a,b,c,d;e,f,g,h)
\mapsto
(ra+b,\ sa,\ rc+d,\ sc;\ re+f,\ se,\ rg+h,\ sg)
\]

These are the governing algebraic rewrites for binary ingestion, following Gosper’s HAKMEM arithmetic. :contentReference[oaicite:2]{index=2}

---

## 11. Gosper emission formula

To emit a regular continued-fraction term, the engine applies the Gosper emission rewrite and updates the transform to the tail stream.

For BLFT, the generic emission rewrite for emitting `(t,u)` is:

\[
(a,b,c,d;e,f,g,h)
\mapsto
(ue,\ uf,\ ug,\ uh;\ a-te,\ b-tf,\ c-tg,\ d-th)
\]

Regular emitted output specializes to the regular CF case.

The public API emits `RCFTerm`, not generalized output terms.

:contentReference[oaicite:3]{index=3}

---

## 12. Collapse rules

## 12.1 EOF is not an error
If an input stream is exhausted, the engine collapses algebraically and continues if possible.

## 12.2 BLFT collapse when `X` hits EOF
BLFT collapses to the induced unary transform in `Y`.

## 12.3 BLFT collapse when `Y` hits EOF
BLFT collapses to the induced unary transform in `X`.

## 12.4 BLFT collapse when both hit EOF
The state becomes terminal / rational / constant and emits the rest directly.

## 12.5 ULFT collapse on EOF
When a ULFT loses its input due to EOF, it collapses directly and emits the rest.

## 12.6 DLFT collapse
DLFT collapses to the corresponding unary state when appropriate.

---

## 13. Emit-vs-ingest decision rule

## 13.1 Primary rule
Follow Gosper’s current widest-corner-range ingest rule exactly.

Do not replace it with speculative “try ingest X and Y, compare future ranges” logic.

## 13.2 Current range evaluation
At any point, the engine evaluates the transform using the current ranges of the live inputs.

For BLFT, use the current endpoint/corner analysis on `X.Range()` and `Y.Range()`.

## 13.3 Emit only when forced
Emit only when the current range proves the next regular term is uniquely determined.

## 13.4 Otherwise ingest
If emission is not yet forced, ingest from the variable associated with the widest current corner range.

## 13.5 Tie-break
If a tie-break rule is needed and Gosper does not force one, use deterministic left-bias:
prefer `X`.

This is current project intent and should be verified during implementation.

---

## 14. `CanEmitRCFTerm`
The project will isolate emission proof logic in a helper conceptually equivalent to:

- `CanEmitRCFTerm(r Range) (term RCFTerm, ok bool)`

This helper decides whether the current exact range forces a unique next regular CF term.

Full details for outside intervals and integer-boundary openness remain open.

---

## 15. ULFT and DLFT

## 15.1 ULFT
Unary transforms are conceptually:

\[
Z(X)=\frac{aX+b}{cX+d}
\]

Identity matrix:

\[
\begin{pmatrix}
1 & 0 \\
0 & 1
\end{pmatrix}
\]

Reciprocal is also a required unary transform.

## 15.2 DLFT
Diagonal/quadratic form begins as “BLFT but diagonal” for v1.

This is a pragmatic starting point, not a proof that all diagonal certification subtleties are permanently solved.

---

## 16. Initial operation plugboards

The client-facing arithmetic façade chooses initial BLFT/ULFT/DLFT plugboards.

Current intended milestone set:

1. identity ULFT passes through `(1,q)`-style sources
2. generalized source support for `p != 1`
3. `square(X)` via DLFT
4. `X + Y`
5. `X * Y`
6. `X - Y`
7. `X / Y`
8. `sqrt(X)` using Ouroboros / Newton feedback on emitted terms
9. target expression:
   `sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))`

Binary BLFT initializations include the standard Gosper forms for addition, subtraction, multiplication, and division. :contentReference[oaicite:4]{index=4}

---

## 17. Normalization and invariants

## 17.1 Transform normalization
After every state mutation:

- ingest
- emit
- collapse

normalize transform coefficients by gcd where possible.

## 17.2 Sign normalization
Apply a conservative canonical sign placement whenever algebraically sound.

## 17.3 Exactness
Normalization must preserve exact semantics.

## 17.4 Preconditions
Validate client-visible constructor/source input early when practical.

## 17.5 Postconditions
Use postcondition/invariant checks to catch engine bugs early.

## 17.6 Panic policy
Panic on internal invariant failure or impossible states.
Do not panic for normal EOF.

---

## 18. Config / limits / policy

## 18.1 Client-supplied config
The client may supply a config object controlling operational limits.

## 18.2 Bit-length checks
Bit-length limits are optional and default off.

Implementation direction:

- inspect all active BLFT coefficients with `BitLen()`
- perform the check just before or after gcd normalization
- intended primarily for testing and diagnostics

## 18.3 Time limit
Time limits are optional and default off.

Intended primarily for testing and diagnostics.

## 18.4 Policy response
Configurable responses may include:

- no-op
- warning hook
- error/status
- panic in stricter test modes

Final API shape remains open.

---

## 19. Error model

## 19.1 EOF
EOF is represented by:

- `ok=false`
- nil tail

EOF is not an error.

## 19.2 Malformed source term
Malformed source input returns an error/status from the source-facing or stream-facing API.

Example malformed term classes include:

- `q == 0`
- negative `q` after the first term
- any other source validity violation defined by the constructor/source contract

## 19.3 Internal invariant failure
Internal algebraic impossibility or violated invariant panics.

## 19.4 NYI
Reserved APIs may return `ErrNYI` explicitly, e.g. `NextDecimalDigit()` in v1.

---

## 20. Proposed wording for first-term restrictions

## 20.1 Input restriction wording
Proposed spec wording:

> A generalized source term is written as `X = p + q/X'`. For v1 source validation, `q` must be nonzero on every term. The first term may have any sign on `p` and `q`. After the first term, `q` must be positive; a later negative `q` is malformed input and causes a source-visible error.

Proposed error wording:

- `ErrZeroQ`
- `ErrNegativeQAfterFirst`
- or one wrapped validation error with those cases as messages/codes

## 20.2 Emission guarantee wording
Proposed spec wording:

> Emitted output terms are regular continued-fraction terms. The first emitted term may be negative. Every later emitted term is emitted only when certified by the current exact range and therefore is a valid regular continued-fraction term in the project’s standard emitted form.

---

## 21. Testing and development process

## 21.1 TDD
Work red/green throughout.

Expected pattern for unimplemented behavior:

- write failing tests first
- create type-correct but wrong-value production stubs where appropriate
- move to green by replacing the stub with correct logic

Do not use “Not Yet Implemented” panics except for APIs explicitly reserved to return `ErrNYI`.

## 21.2 White-box vs black-box tests
For `core`:

- WB tests live in package `core`
- BB tests live in package `core_test`

Naming conventions:

- WB test names begin with `TestWB_`
- BB test names begin with `TestBB_`
- WB filenames include `_wb_`
- BB filenames include `_bb_`

Apply the analogous convention in `named`.

## 21.3 Property testing
Develop WB property/invariant tests in parallel with invariant enforcement.

## 21.4 Repo truth
After each green, commit and push.
Pushed code is the authoritative implementation state.

---

## 22. Open questions

The following questions remain open and should be tracked explicitly.

### 22.1 Final public API shapes
Whether `PQSource`, `RCFStream`, and `GCF` become separate public interfaces/types or a smaller public façade with private implementations.

### 22.2 Exact Go method signatures
Especially:

- `Next()` shape
- `NextRCF()` shape
- whether status codes are preferable to plain `error`
- exact config injection pattern

### 22.3 Exact `Range` struct layout
Including:

- endpoint structs
- inside/outside encoding
- canonical representation for outside intervals
- exact openness rules at integer boundaries

### 22.4 Exact `CanEmitRCFTerm` logic
Especially for:

- outside intervals
- endpoint openness at integer boundaries
- exact proof rule for first vs later emitted terms

### 22.5 Exact tie-break semantics
If Gosper leaves a tie ambiguous, confirm the deterministic left-bias rule.

### 22.6 Exact limit-breach behavior
Whether bit/time policy breaches return error, status, callback signal, panic, or a configurable mix.

### 22.7 Final DLFT guarantees
Whether “BLFT but diagonal” remains sufficient or must be strengthened with additional certified range logic.

### 22.8 Named package scope
Whether `named` should include only named constants/sources or also all finite wrappers and source constructors.

### 22.9 `Cmp`
Exact semantics and API for `GCF.Cmp()`.

### 22.10 Decimal-digit API after v1
How decimal-digit generation should coexist with `NextRCF()` once implemented.

---