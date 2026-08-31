# CLAUDE.md: Coding & Development Guide

This repository contains **Vidonyx**, a declarative video composition engine and FFmpeg filtergraph compiler in Go.

---

## RTK Token-Optimized CLI Commands

When using terminal commands, prefer using `rtk` where applicable:
```bash
rtk gain              # Show token savings analytics
rtk gain --history    # Show command history with token savings
```

---

## Core Build & Test Commands

```bash
# Run all unit tests with race detection (mandatory for all changes)
go test -race -v ./...

# Run tests for specific packages
go test -race -v ./types/...
go test -race -v ./timeline/...
go test -race -v ./filtergraph/...
go test -race -v ./compiler/...
go test -race -v ./executor/...
go test -race -v ./composer/...

# Run example recipes
go run ./examples/01_simple_cut/main.go
go run ./examples/02_picture_in_picture/main.go
```

---

## Key Coding Conventions

1. **Idiomatic Go**: Use Go 1.22+ iterators, `log/slog`, `errors.Join`, and table-driven testing.
2. **Zero String Concatenation for Filters**: Always model transformations as DAG nodes and pads in `filtergraph/`.
3. **Rational Arithmetic**: Use `types.Rational` for temporal precision and frame conversions.
4. **Hermetic Testing**: Never require external `ffmpeg` binaries in unit tests; use `executor.NewMockExecutor()`.
5. **Clean Docstrings**: Ensure every exported symbol has a standard Go doc comment.
