<!-- MasterPlan.md v1 -->
# goGCF Master Plan

>
## Mission
Build a mathematically correct, testable Go library for continued-fraction arithmetic in the spirit of Gosper/HAKMEM, using immutable generalized-CF sources and immutable evaluators that emit regular CF terms.

## Current status
The project has moved beyond repo/type scaffolding into a real evaluator path for exact finite arithmetic cases, while still preserving the intended infinite-first architecture.

### Working now
- Core public types exist: `PQTerm`, `RCFTerm`, `Rational`, `Bound`, `Range`, `GCFSource`, `GCF`.
- Immutable sources work: `Int64`, `Rat64`, `Sqrt2`.
- `FromSource` pass-through evaluator works for RCF-compatible sources.
- `Take(n)` and `Convergent()` work for finite `rcfPrefixGCF`.
- Rational-to-RCF conversion works.
- Internal ULFT/BLFT identity initializers exist.
- Binary decision enum and choose logic exist.
- BLFT ingest-left / ingest-right transforms exist and are tested.
- EOF collapse exists: `binaryLFT -> unaryLFT -> Rational`.
- Private adapter `RCFTerm -> PQTerm` exists.
- `binaryGCF` exists for `Add`, `Sub`, `Mul`, `Div`.
- `binaryGCF` can step, ingest child output, collapse to exact rational, and emit correct exact-finite public results for current tested cases.
- `binaryGCF.Next()` now preserves continuation state by wrapping the resolved remainder back into a `binaryGCF`.


## Critical path
1. Make `binaryGCF.Next()` genuinely stepwise and emission-driven.
2. Define and test the first correct emit-certification rule for `binaryLFT`.
3. Add arithmetic `Range()` behavior for binary/unary evaluators.
4. Add unary evaluator shell and unary stepping.
5. Implement `Neg` / `Inv`.
6. Implement `sqrt`.
7. Implement named sources `Pi()` and `E()`.
8. Implement `sin` and `tanh`.
9. Compose and run the target formula in `main.go`.
10. Expand certification/progress logic for nontrivial infinite arithmetic.

## Deferred / future work
- Decimal or radix-digit emission.
- Additional named constant generators beyond the minimum target-formula set.
- DLFT only when needed for `sqrt(x)` or other diagonal-use cases.
- Broader certification logic and richer progress/blocked diagnostics.
- More complete `Range` ordering/comparison semantics as arithmetic matures.

## Known risks / unresolved design questions
- Exact emit-certification rule for public stepwise emission from `binaryLFT` is still unresolved.
- Current EOF handling via collapse is coherent for finite exact cases, but needs extension for general incremental arithmetic.
- Unary operator architecture remains open enough that early binary decisions should avoid constraining it unnecessarily.
- Must preserve immutability and avoid accidental mutation/leakage of resolved/internal continuation state.

<!-- MasterPlan.md v1 -->