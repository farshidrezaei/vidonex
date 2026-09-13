# CLAUDE.md: Coding & Development Guide

This repository contains **Vidonex**, a declarative video composition engine and FFmpeg filtergraph compiler in Go.

---

## RTK Token-Optimized CLI Commands

When using terminal commands, prefer using `rtk` where applicable:
```bash
rtk gain              # Show token savings analytics
rtk gain --history    # Show command history with token savings
```

---

## Core Build, Test & Lint Commands

```bash
# Run all unit and real-FFmpeg E2E tests with race detection (mandatory for all changes)
go test -race -v ./...

# Run static analysis and linting (mandatory 0 issues)
golangci-lint run ./...

# Run Real-FFmpeg E2E test suite specifically
go test -race -v ./tests/e2e/...

# Run all 10 example recipes
go run ./examples/01_simple_cut/main.go
go run ./examples/02_picture_in_picture/main.go
go run ./examples/03_transitions/main.go
go run ./examples/04_animated_motion/main.go
go run ./examples/05_animated_subtitles/main.go
go run ./examples/06_audio_ducking/main.go
go run ./examples/07_platform_presets/main.go
go run ./examples/08_podcast_audio_waveform/main.go
go run ./examples/09_green_screen_studio/main.go
go run ./examples/10_declarative_yaml_project/main.go

# Run Web Studio & UI Frontend
cd ui && pnpm dev     # Run Nuxt 4 development workstation
cd ui && pnpm build   # Production static build & validation

# Run CLI Tool
./vidonex --help
```

---

## Mandatory Engineering Contracts

1. **Descriptive Naming**: No cryptic abbreviations (`compositionTimeline`, `durationSeconds`, `sourceOutputPad`, `waveformOptions`).
2. **Zero String Concatenation for Filters**: Always model transformations as DAG nodes and pads in `filtergraph/` and connect with `graph.Connect()`.
3. **Rational Arithmetic**: Use `types.Rational` or `time.Duration` for temporal precision and frame conversions to prevent floating-point drift.
4. **Table-Driven Tests & Real E2E Verification**: All unit tests must be table-driven (`Test...Table`). Real FFmpeg tests live in `tests/e2e/`.
5. **Continuous Documentation Synchronization**: With every feature addition, refactor, or code modification, ALL documentation (`.md` files, diagrams, contracts, and example recipes) MUST be updated immediately to keep code and documentation 100% in sync.
6. **Standard Library First**: Core packages must remain 100% standard library only.
