# Vidonex Development & Engineering Contracts

This document specifies the **Mandatory Development Contracts and Invariants** that every human contributor and AI agent MUST adhere to when authoring or modifying code in the **Vidonex** repository.

---

## 1. Descriptive Naming (No Cryptic Abbreviations)

1. **Function and Method Names**:
   - MUST clearly express their exact operation and intent (e.g., `InjectVideoNormalizer`, `BuildVideoCompositor`, `CalculateOffset`, `FormattedFilterComplex`, `ApplySidechainDucking`, `ApplyWaveformVisualizer`, `ApplyChromaKey`).
   - Avoid vague names like `Do()`, `Handle()`, `Process()`, or `Init()` unless explicitly scoped.

2. **Variable and Parameter Names**:
   - Variable and parameter names MUST be descriptive, meaningful, and self-documenting (e.g., `compositionTimeline`, `durationSeconds`, `sourceOutputPad`, `destinationInputPad`, `consumerCount`, `frameRate`, `waveformOptions`).
   - **DO NOT** use cryptic or arbitrary abbreviations (`tl`, `tr`, `dur`, `res`, `pct`, `sz`, `sc`, `pts`, `inPad`, `outPad`, `src`, `dst`, `fmt`).

---

## 2. Exhaustive Table-Driven Testing & Real E2E Verification

1. **Table-Driven Tests Across All Scenarios**:
   - All unit tests and feature/integration tests MUST follow Go's **Table-Driven Test** structure (`Test...Table`).
   - Test suites MUST exhaustively cover:
     - **Happy paths** and common use cases.
     - **Edge cases** and boundary conditions (zero values, extreme values, negative intervals).
     - **Error paths** and invalid parameter validation.
     - **Complex graph topologies** (e.g., diamond graphs, multi-consumer fan-outs, cycles).

2. **Deterministic Hermetic Unit Tests & Synthetic Real E2E Tests**:
   - Core package unit tests MUST NOT depend on a locally installed binary. Use `executor.NewMockExecutor()` and graph structure assertions.
   - The end-to-end suite (`tests/e2e/`) generates synthetic media with FFmpeg and executes real renders, probing output containers with `probe.FFprobeProber` to guarantee byte-accurate container validity.

---

## 3. Mandatory Quality Checks (`golangci-lint` & Race Detector)

Before finishing any task or opening a pull request, the following commands MUST pass with **zero warnings and zero errors**:

1. **Linter Inspection**:
   ```bash
   golangci-lint run ./...
   ```
   *Requirement: 0 issues reported.*

2. **Race Condition Detector & Test Suite**:
   ```bash
   go test -race -v ./...
   ```
   *Requirement: 100% tests pass with no race conditions.*

---

## 4. Modern Go Idioms & Cutting-Edge Standards

1. **Language Standards (Go 1.22+ / 1.23+)**:
   - Utilize standard Go iterators (`iter.Seq`, `iter.Seq2`) for zero-allocation slice and graph traversals.
   - Employ `log/slog` for structured logging.
   - Aggregate validation errors using `errors.Join`.
   - Leverage Go Generics where appropriate for type-safe data structures.

2. **Temporal Precision (Zero Floating-Point Drift)**:
   - NEVER use raw `float64` for core timestamps and frame rate calculations. Always use `types.Rational` or `time.Duration`.

3. **Graph-First Architecture (No String Concatenation)**:
   - NEVER construct raw FFmpeg filter strings manually with string formatting in domain logic.
   - Always construct a `filtergraph.Graph`, add `filtergraph.Node`s with `filtergraph.Pad`s, and link them using `graph.Connect()`.

4. **Standard Library First**:
   - Core packages (`types`, `timeline`, `filtergraph`, `compiler`) must remain 100% standard library only.

---

## 5. Mandatory Continuous Documentation Synchronization

1. **Keep Code, Tests, and Documentation 100% in Sync**:
   - With every feature addition, modification, or architectural refactor, ALL relevant documentation (`.md` files, architecture diagrams, contracts, and example recipes) MUST be updated immediately.
   - Never leave documentation outdated, incomplete, or referencing stale code signatures or removed concepts.
   - Any new package or feature must be represented in `README.md`, `CONTRACTS.md`, `AGENTS.md`, and `CLAUDE.md`, accompanied by a runnable example recipe in `examples/`.
