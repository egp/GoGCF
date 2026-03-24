# Spec.md

# goGCF / core Specification

## Status

This document is the current design contract for the `core` package and the companion package `named`.

When a design choice conflicts with Gosper / HAKMEM 101B, Gosper wins.

For GitHub rendering, all displayed formulas in this file use fenced `math` blocks rather than bracketed pseudo-math.

---

## 1. Mission

Build a mathematically correct, testable Go library for arithmetic on generalized continued fractions in the style of Gosper / HAKMEM 101B.

Primary goals:

- mathematical correctness
- simplicity
- exact arithmetic only
- DRY implementation
- design for test
- top-down and bottom-up development together
- red/green TDD with realistic wrong-value stubs where appropriate

Performance matters only after correctness and clarity.

### Target formula

```math
\frac{\sqrt{\frac{3}{\pi^2} + e}}{\tanh(\sqrt{5}) - \sin(69^\circ)}
```

This is a future celebratory WB test target.

---

## 2. Package layout

### 2.1 `core`

The `core` package contains:

- generalized source-stream abstractions
- regular CF output-stream abstractions
- BLFT-centered transform engine
- range logic
- normalization and invariant enforcement
- client-facing arithmetic façade
- config / limits / status codes

### 2.2 `named`

The companion package is `named`.

`named` contains:

- named streams such as `Pi()`, `E()`, `Sqrt2()`, etc.
- constructors from `Rational`
- constructors from finite or procedural source descriptions
- wrappers that help clients build streams for testing and use

`named` does not weaken the `core` rule that engine ingestion is always in generalized `(p,q)` terms.

---

## 3. Core design principles

### 3.1 Exact arithmetic only

No floats anywhere in production code, tests, or correctness helpers.

### 3.2 Simplicity

Prefer one clear algebraic path over multiple convenience paths.

### 3.3 DRY

ULFT, BLFT, and DLFT must share helpers wherever algebraically sound.

### 3.4 One internal engine shape

Internally, prefer one BLFT-sized transform state. Unary and diagonal cases are specialized initializations or constrained views of the same machinery whenever practical.

### 3.5 Certified output only

The engine emits a regular continued-fraction term only when Gosper’s emission criterion says it is forced by the current range.

### 3.6 EOF is valid

EOF is not an error. EOF triggers algebraic collapse and computation continues if possible.

### 3.7 Malformed client-visible input is reported

Malformed source terms supplied by client-visible streams return a failure status.

### 3.8 Internal invariant failure is fatal

Internal impossible states panic.

### 3.9 Invalid use should be unrepresentable where possible

APIs should prefer states and signatures that make misuse hard or impossible.

---

## 4. Mathematical conventions

### 4.1 Generalized source-term convention

Every generalized source term uses:

```math
X = p + \frac{q}{X'}
```

This is the exact convention everywhere in `core`.

### 4.2 Emission convention

Emission uses Gosper’s regular tail rewrite:

```math
Z = t + \frac{u}{Z'} \quad\Longleftrightarrow\quad Z' = \frac{u}{Z - t}
```

For regular emitted output, `u = 1`.

### 4.3 Regular emitted output convention

Emitted output is a regular continued-fraction term sequence.

- the first emitted regular term may be negative
- every later emitted regular term is a standard regular CF term

### 4.4 Input sign policy for generalized source terms

For v1 source validation:

- `p` may be negative on any term
- `q` must be nonzero on every term
- the first term may have negative `q`
- after the first term, negative `q` is malformed input

This is a practical v1 source-validity rule. It is stricter than pure algebraic permissiveness.

### 4.5 No internal `q == 1` shortcut path

Regular CFs are the special case `q = 1`, but core logic does not fork into a separate arithmetic path for regular CFs.

### 4.6 Exact value vs integer value

If a range has `lo == hi`, then the value is known exactly.

That does **not** imply the value is an integer. It may be any exact rational.

---

## 5. Numeric types

### 5.1 Big integers

Use `*big.Int` internally for:

- `PQTerm.p`
- `PQTerm.q`
- `RCFTerm`
- BLFT coefficients
- Rational numerator
- Rational denominator

### 5.2 Rational

`Rational` is exact and BigInt-based.

Expected operations include:

- normalization
- comparison
- sign
- arithmetic for endpoint evaluation
- exact conversion from integer

### 5.3 Exact conversion from integer

“Exact conversion from integer” means:

- integer `n` converts to rational `n / 1`
- no rounding
- no approximation
- no float intermediary

### 5.4 Int64 usage

Use `int64` only where overflow is impossible by construction.

No correctness-critical internal algorithm may rely on `int64` if overflow is possible.

---

## 6. Status codes

The project prefers status codes over plain `error` returns for operational APIs.

### 6.1 Status direction

Representative statuses include:

- `StatusOK`
- `StatusEOF`
- `StatusInvalidInput`
- `StatusBitLenExceeded`
- `StatusTimeout`
- `StatusNYI`
- `StatusIndeterminate`

Exact names may change during cleanup.

### 6.2 Meaning

- `StatusOK`: a term was produced successfully
- `StatusEOF`: normal end of stream
- `StatusInvalidInput`: malformed client-visible source term or invalid constructor input
- `StatusBitLenExceeded`: configured bigint limit was exceeded
- `StatusTimeout`: configured runtime limit was exceeded
- `StatusNYI`: reserved API exists but is not implemented yet
- `StatusIndeterminate`: computation did not prove the requested result before a policy limit or logical stop

### 6.3 Why status instead of `(ok, err)`

EOF is not an error, and malformed input is not EOF. Status codes carry that distinction directly and are easier to inspect during library and client debugging.

---

## 7. Conceptual architecture

There are three conceptual pieces.

### 7.1 Source stream

A generalized source stream provides:

- the next generalized term `(p,q)`
- the tail source `X'`
- the current range enclosure

Conceptual direction:

- `Next() -> (PQTerm, tail, Status)`
- `Range() -> Range`

### 7.2 Output stream

A regular output stream provides:

- the next regular term
- the tail output stream `Z'`
- the current range enclosure

Conceptual direction:

- `NextRCF() -> (RCFTerm, tail, Status)`
- `Range() -> Range`

### 7.3 `GCF`

`GCF` is the client-facing attachment point.

A client creates a `GCF` from:

- an initial transform plugboard
- zero, one, or two input source streams
- an optional config, otherwise defaults

Then calls:

- `Next()` to obtain the next emitted regular term
- `Range()` to inspect the current enclosure
- `Cmp()` to compare values
- `NextDecimalDigit()` for the reserved decimal API

### 7.4 No public `Emit()` method

There is no public `Emit()` method.

Emission is an internal event. To the client, emission simply appears as a successful `Next()` / `NextRCF()` result.

---

## 8. Client-facing operations

The client should have access to arithmetic by selecting predefined plugboards.

Representative direction:

- `Add(X, Y)`
- `Sub(X, Y)`
- `Mul(X, Y)`
- `Div(X, Y)`
- `Reciprocal(X)`
- `Square(X)`
- later `Sin(X)`, `Tanh(X)`, `Sqrt(X)`

Operationally, these morph into choosing the proper transform initialization and then running the common engine.

There is no special core arithmetic path for division, addition, etc., beyond plugboard selection plus common engine logic.

---

## 9. Reserved decimal API

`NextDecimalDigit()` exists in v1 only as a reserved stub.

Behavior:

- returns `(0, StatusNYI)`

Mixing future decimal and regular-term APIs on the same evolving object is currently unspecified and postponed.

---

## 10. Range model

### 10.1 Range kinds

A `Range` is one of:

- `InsideInterval`
- `OutsideInterval`

### 10.2 Endpoint type

Endpoints are exact `Rational`s.

Each endpoint has its own openness/closedness flag.

### 10.3 Inside interval meaning

An inside interval encloses values within bounded endpoints.

Invariant:

- `lo <= hi`

Interpretation:

- the value lies within the interval
- if `lo == hi`, the exact value is known

### 10.4 Outside interval meaning

An outside interval denotes the enclosure:

```math
(-\infty, lo] \cup [hi, +\infty)
```

with independently open/closed finite endpoints.

Canonical ordering details for outside intervals remain open.

### 10.5 Range comparison ordering

Use Gosper’s uncertainty ordering:

- inside narrow
- inside wide
- outside wide
- outside narrow

The intent is:

- earlier in the list means more certain / better / smaller uncertainty
- later in the list means less certain / worse / larger uncertainty

This ordering need not expose a numeric width to clients.

### 10.6 Range validity

`Range()` is defined only on live objects.

It does not return a status.

Constructors and step methods carry status.
Invalid internal use panics.

---

## 11. Exact openness rules at integer boundaries

These rules matter for emission.

### 11.1 Inside intervals

For an inside interval, the next regular term is forced only if **every** value in the interval has the same floor.

Define:

- `lowerTerm = floor(lo)`
- `upperTerm = floor(hi)`, except:
  - if `hi` is open and `hi` is exactly an integer `n`, then `upperTerm = n - 1`

Then an inside interval certifies a unique next regular term iff:

- `lowerTerm == upperTerm`

Examples:

- `[1.2, 1.9]` emits `1`
- `[1.2, 2)` emits `1`
- `[1, 2)` emits `1`
- `[1, 2]` does **not** emit
- `(2, 2.5]` emits `2`

### 11.2 Outside intervals

For v1, `OutsideInterval` does not directly certify a unique regular output term.

If the current range is outside, the engine ingests rather than emitting.

This matches Gosper’s stated emission rule, which is phrased for inside intervals whose endpoint floors agree.

---

## 12. BLFT core engine

### 12.1 Primary internal state

The primary internal transform state is BLFT-shaped:

```math
Z(X,Y) = \frac{aXY + bX + cY + d}{eXY + fX + gY + h}
```

with BigInt coefficients.

ULFT and DLFT begin as constrained/specialized cases of the same internal machinery where practical.

### 12.2 Core engine actions

At each step the engine does exactly one of:

1. ingest one `(p,q)` term from one input
2. emit one regular CF term
3. collapse on EOF from an input

The difficult decision is between ingest and emit.

---

## 13. Gosper ingest formulas

If `X` produces `(p,q)` with

```math
X = p + \frac{q}{X'}
```

then ingesting `X` updates BLFT coefficients to:

```math
(a,b,c,d;e,f,g,h)
\mapsto
(pa+c,\ pb+d,\ qa,\ qb;\ pe+g,\ pf+h,\ qe,\ qf)
```

If `Y` produces `(r,s)` with

```math
Y = r + \frac{s}{Y'}
```

then ingesting `Y` updates BLFT coefficients to:

```math
(a,b,c,d;e,f,g,h)
\mapsto
(ra+b,\ sa,\ rc+d,\ sc;\ re+f,\ se,\ rg+h,\ sg)
```

---

## 14. Gosper emission formula

To emit `(t,u)` so that

```math
Z = t + \frac{u}{Z'}
```

the BLFT coefficients update to:

```math
(a,b,c,d;e,f,g,h)
\mapsto
(ue,\ uf,\ ug,\ uh;\ a-te,\ b-tf,\ c-tg,\ d-th)
```

For regular emitted output, `u = 1`.

---

## 15. Collapse rules

### 15.1 EOF is not an error

EOF is represented by `StatusEOF` and causes algebraic collapse.

### 15.2 BLFT collapse when `X` hits EOF

BLFT collapses to the induced unary transform in `Y`.

### 15.3 BLFT collapse when `Y` hits EOF

BLFT collapses to the induced unary transform in `X`.

### 15.4 BLFT collapse when both hit EOF

The state becomes terminal / rational / constant and emits the rest directly.

### 15.5 ULFT collapse on EOF

When a ULFT loses its input due to EOF, it collapses directly and emits the rest.

### 15.6 DLFT

“BLFT but diagonal” is postponed as a final guarantee question.
For now, keep the decision explicitly deferred.

---

## 16. Emit-vs-ingest decision rule

### 16.1 Primary rule

Follow Gosper’s current widest-corner-range ingest rule exactly.

Do not replace it with speculative “try ingest X and Y, compare future ranges” logic.

### 16.2 Emit only when forced

Emit only when the current range proves the next regular term is uniquely determined.

### 16.3 Otherwise ingest

If emission is not forced, ingest from the variable associated with the widest current corner range.

### 16.4 Tie-break

If a tie-break is needed, tie goes to `X`.

---

## 17. `CanEmitRCFTerm`

Model `CanEmitRCFTerm` as a predicate/helper, not the primary public API.

Conceptual direction:

- `CanEmitRCFTerm(r Range) -> (RCFTerm, bool)`

For v1:

- only `InsideInterval` can certify emission
- use the integer-boundary openness rules from Section 11
- `OutsideInterval` returns `false`

A state-transition diagram may be added later as explanatory documentation, but the core rule is the predicate.

---

## 18. Normalization and invariants

### 18.1 Transform normalization

After every state mutation:

- ingest
- emit
- collapse

normalize transform coefficients by gcd where possible.

### 18.2 Sign normalization

Apply a conservative canonical sign placement whenever algebraically sound.

### 18.3 Exactness

Normalization must preserve exact semantics.

### 18.4 Preconditions

Validate client-visible input early when practical.

### 18.5 Postconditions

Use invariant checks to catch engine bugs early.

### 18.6 Bit-length checks

If enabled by config, inspect all eight BLFT coefficients using `BitLen()` around each operation, just before or just after gcd normalization.

### 18.7 Panic policy

Panic on:

- internal invariant failure
- impossible internal algebraic states
- configured timeout breach
- configured excessive bit-length breach

Do not panic for normal EOF.

---

## 19. Config

The client can supply config either:

- at creation time as an optional parameter, or
- via a set/get config object

Defaults must exist either way.

Representative config controls include:

- bit-length limit, default off
- timeout limit, default off
- panic/warn behavior if that ever becomes configurable beyond v1

---

## 20. `Cmp`

### 20.1 Intent

`GCF.Cmp()` is important and should follow Gosper-style range reasoning.

### 20.2 Direction

Recommended semantics:

- compare `X` and `Y` by comparing the sign of `X - Y`
- operationally, construct the subtraction transform and refine until:
  - range proves negative
  - range proves zero exactly
  - range proves positive
  - or a configured limit/status intervenes

### 20.3 Proposed result style

Direction:

- return a comparison code such as `-1`, `0`, `+1`
- plus a `Status`

Representative behavior:

- `(-1, StatusOK)` if `X < Y`
- `(0, StatusOK)` if equality is proven
- `(+1, StatusOK)` if `X > Y`
- `(0, StatusIndeterminate)` if a policy limit fires before proof
- other non-OK statuses as appropriate

### 20.4 Exact zero

If subtraction yields a range whose exact value is proven to be zero, `Cmp` returns equality.

---

## 21. Milestones

Current intended milestone set:

1. identity unary pass-through for regular-style input where `q = 1`
2. generalized source support for `q != 1`
3. `square(X)`
4. `X + Y`
5. `X * Y`
6. `X - Y`
7. `X / Y`
8. `sqrt(X)` using Ouroboros / Newton feedback
9. target formula

---

## 22. Testing and development process

### 22.1 TDD

Work red/green throughout.

Typical pattern:

- write failing tests first
- use type-correct but wrong-value stubs where appropriate
- move to green by replacing the stub with correct logic

### 22.2 WB / BB conventions

For `core`:

- WB tests live in package `core`
- BB tests live in package `core_test`

Naming conventions:

- WB test names begin with `TestWB_`
- BB test names begin with `TestBB_`
- WB filenames include `_wb_`
- BB filenames include `_bb_`

Apply the analogous convention in `named`.

### 22.3 Property testing

Co-develop WB property tests in parallel with invariant enforcement.

Each important struct should have an explicit list of invariants and corresponding WB checks.

### 22.4 Repo truth

After each green, commit and push.

Pushed code is the authoritative implementation state and can be re-read from the GitHub repo in later chats, reducing the need to cache old file contents except for work in progress.

---

## 23. Sources

HAKMEM / Gosper:
http://w3.pppl.gov/~hammett/work/2009/AIM-239-ocr.pdf

GitHub math rendering docs:
http://docs.github.com/en/get-started/writing-on-github/working-with-advanced-formatting/writing-mathematical-expressions

Go `math/big` docs:
http://pkg.go.dev/math/big