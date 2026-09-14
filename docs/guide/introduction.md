# Introduction to Vidonex

**Vidonex** is an open-source, high-performance video composition engine and non-linear video editing workstation engineered natively in Go and Vue 3 / Nuxt UI.

## The Problem with Traditional Video Tooling

When building automated video pipelines, media servers, or desktop video software, engineers typically encounter two frustrating extremes:

1. **Raw FFmpeg Scripts & Strings**:
   - Manually constructing 1,000-character `-filter_complex` strings with Bash or string formatting.
   - Brittle pad reuse errors (`Filter x has no output pad named [v1]`).
   - Floating-point timing inaccuracies that cause audio/video desync over long renders.
   - High cognitive load and virtually untestable logic.

2. **Heavy Node.js / Python Frameworks**:
   - Libraries like Remotion (React/Node) rely on headless Chromium (Puppeteer), consuming gigabytes of RAM and leading to slow frame-by-frame HTML canvas captures.
   - Python libraries like MoviePy suffer from Python's Global Interpreter Lock (GIL) and lack type safety and native desktop UI capabilities.

## The Vidonex Solution

Vidonex introduces a compiler architecture inspired by modern compilers (like LLVM):

```
Declarative Timeline AST (Go / YAML)
       ↓
Filtergraph Intermediate Representation (DAG)
       ↓
Graph Optimization Passes (AutoSplit, DCE, Format Coercion)
       ↓
Strict Target Emitter (FFmpeg CLI / Mermaid Flowcharts)
```

### Core Value Propositions

- **Sub-millisecond Compilation**: The compiler builds and validates the DAG in less than 2 milliseconds.
- **Zero Pad-Reuse Errors**: The `AutoSplitPass` statically calculates how many downstream filters consume each audio or video pad, automatically inserting `split` or `asplit` nodes.
- **Rational Math**: Frame numbers, durations, and timestamps are calculated with `types.Rational` fractions, eliminating cumulative float drift.
- **Zero Electron Overhead**: The visual desktop studio is built with Wails v2 and Nuxt UI, running smoothly with less than 40MB of baseline memory.
- **Full Cloud & CI/CD Readiness**: Vidonex compiles to a standalone, zero-dependency static Go binary.
