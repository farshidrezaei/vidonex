# Contributing to Vidonex

Thank you for contributing to **Vidonex**! We welcome contributions from developers, researchers, and video engineers worldwide.

---

## 1. Code Standards & Engineering Principles

Vidonex follows strict production-grade Go conventions:

1. **Pure Standard Library Core**: Do not add heavy external dependencies to core packages (`types`, `timeline`, `filtergraph`, `compiler`).
2. **Modern Go Idioms**:
   - Use Go iterators (`iter.Seq`, `iter.Seq2`) for graph and slice traversals.
   - Use `log/slog` for structured logging.
   - Use `errors.Join` for multi-error collection.
   - Use `types.Rational` instead of raw `float64` for video timestamps and frame rates.
3. **Table-Driven Tests**: All new features, filters, and passes must include table-driven unit tests covering happy paths, edge cases, and error conditions.
4. **Race Detection**: All tests must pass with the `-race` detector enabled.
5. **Documentation**: Every exported struct, interface, function, and constant must have a clear, descriptive docstring.

---

## 2. Development Workflow

### Prerequisites
- Go 1.22+ (or Go 1.23+)
- Optional: `ffmpeg` binary (for manual end-to-end rendering tests)

### Running Tests
Run the entire test suite with race detection:
```bash
go test -race -v ./...
```

### Running Examples
```bash
go run ./examples/01_simple_cut/main.go
go run ./examples/02_picture_in_picture/main.go
```

---

## 3. Submitting Pull Requests

1. **Fork and Branch**: Create a feature branch with a descriptive name (e.g. `feat/chroma-key-enhancement` or `fix/audio-desync-pad`).
2. **Implement & Test**: Ensure all existing and new unit tests pass cleanly.
3. **Commit Messages**: Follow [Conventional Commits](https://www.conventionalcommits.org/) format (e.g. `feat: add drawbox effect`, `fix: handle negative source trims`).
4. **Open PR**: Submit your pull request with a clear summary of changes and diagram outputs if modifying the compiler or filtergraph IR.
