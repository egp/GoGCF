<!-- HighLevelDesign.md v1 -->
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

```text
+------------------+        lift/open         +------------------+
|    GCFSource     | -----------------------> |       GCF        |
| immutable input  |                          | immutable eval    |
| emits PQTerm     |                          | emits RCFTerm     |
+------------------+                          +------------------+
         |                                              |
         | built by                                      | observed by
         |                                               |
         v                                               v
+------------------+                          +------------------+
| constructors and |                          | Range, Take,     |
| named sources    |                          | Convergent       |
+------------------+                          +------------------+
                                                       |
                                                       | implemented by
                                                       v
                                             +------------------+
                                             | ULFT / BLFT /    |
                                             | DiagonalLFT      |
                                             | internal engine  |
                                             +------------------+