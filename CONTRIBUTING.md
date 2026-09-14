# Contributing to Vidonex

First off, thank you for taking the time to contribute! 🎉

Vidonex is an open-source, high-performance video composition engine and modern desktop/web workstation written in Go, Vue 3, and Nuxt UI. We welcome contributions of all kinds: bug reports, documentation improvements, feature requests, new video/audio filters, and architectural enhancements.

---

## Code of Conduct

All contributors and maintainers are expected to follow our [Code of Conduct](CODE_OF_CONDUCT.md). Please report any violations or abusive behavior to [security@vidonex.io](mailto:security@vidonex.io).

---

## How Can I Contribute?

- **Reporting Bugs**: Check existing issues first. If not reported, submit a detailed report using the [Bug Report Template](.github/ISSUE_TEMPLATE/bug_report.yml).
- **Proposing Features**: Open a discussion or use the [Feature Request Template](.github/ISSUE_TEMPLATE/feature_request.yml) explaining the use case, API design, and expected FFmpeg filtergraph output.
- **Improving Documentation**: Fix typos, add recipes in `examples/`, or improve guides in `docs/`.
- **Submitting Code**: Pick up an issue marked `good first issue` or `help wanted`, or propose an optimization.

---

## Development Setup

### Prerequisites
- **Go**: 1.22+ (Go 1.25/1.26 recommended)
- **Node.js**: 20+ & `pnpm` (if working on the Web Studio in `ui/`)
- **FFmpeg & FFprobe**: Version 5.x, 6.x, or 7.x installed in system `$PATH` (for E2E tests and rendering)
- **golangci-lint**: `v1.59+`

### Repository Setup
```bash
# Clone the repository
git clone https://github.com/farshidrezaei/vidonex.git
cd vidonex

# Download Go modules
go mod download

# Build CLI binary
make build
```

---

## Engineering Contracts & Guidelines

Vidonex enforces strict architectural contracts documented in [CONTRACTS.md](CONTRACTS.md):

1. **Pure Standard Library Core**:
   - Packages `types`, `timeline`, `filtergraph`, and `compiler` MUST NOT import external third-party dependencies.
2. **No String Concatenation for Filtergraphs**:
   - NEVER construct raw FFmpeg filter strings manually with `fmt.Sprintf` or string concatenation. Build a `filtergraph.Graph`, add `filtergraph.Node`s with typed pads, and connect them.
3. **No Floating-Point for Core Timing**:
   - Always use `types.Rational` or `time.Duration` for timestamps, frame rates, and durations to prevent float precision drift over time.
4. **Stream Type Safety**:
   - Video pads (`StreamTypeVideo`) can only connect to video pads. Audio pads (`StreamTypeAudio`) can only connect to audio pads.
5. **Table-Driven Tests**:
   - Every unit test must be table-driven (`Test...Table`), covering happy paths, edge cases, zero-values, and error branches.

---

## Quality Checks & Verification

Before submitting your pull request, ensure all linters and tests pass cleanly:

```bash
# 1. Run static analysis & linter
golangci-lint run ./...

# 2. Run unit tests with race detection
go test -race -v ./...

# 3. Run real-FFmpeg end-to-end tests (requires ffmpeg installed)
go test -race -v ./tests/e2e/...
```

---

## Pull Request Guidelines

1. **Branch Naming**:
   - `feat/feature-name`
   - `fix/bug-description`
   - `docs/topic-update`
   - `refactor/subsystem-name`
2. **Conventional Commits**:
   - Format: `<type>(<scope>): <short summary>`
   - Examples:
     - `feat(ducking): add dynamic attack and release curves`
     - `fix(filtergraph): prevent dangling output pad in auto-split pass`
     - `docs(readme): add comparison matrix with Remotion`
3. **Keep PRs Focused**: A pull request should address a single concern. If you have multiple independent improvements, split them into separate PRs.

Thank you for helping make Vidonex the premier programmatic video engine! 🚀
