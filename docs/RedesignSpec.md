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