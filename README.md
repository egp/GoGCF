<!-- README.md v1 -->
# GoGCF

GoGCF is a Go library for exact arithmetic on continued fractions in the style of Gosper.

Current focus:
- immutable `GCFSource` and `GCF`
- exact `BigInt` / `Rational` support
- RCF emission from GCF arithmetic
- correctness-first TDD

Key docs:
- `docs/RequirementsSpec.md`
- `docs/PublicAPI.md`
- `docs/HighLevelDesign.md`

Early milestone:
compute `sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))`

<!-- README.md v1 -->