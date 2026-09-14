# Versioning Policy & Stability

Vidonex strictly follows [Semantic Versioning (SemVer 2.0.0)](https://semver.org/).

Given a version number `MAJOR.MINOR.PATCH`, increments indicate:

- **`MAJOR`**: Incompatible API or timeline schema changes.
- **`MINOR`**: Backwards-compatible new features (new effects, compiler optimizations, CLI commands).
- **`PATCH`**: Backwards-compatible bug fixes and performance patches.

---

## Stability Guarantees

### 1. Timeline Specification Schema (`project.yaml`)
- The declarative YAML/JSON format adheres to strict schema versioning (`version: "1.0"`).
- Any project authored for `v1.0.0` is guaranteed to compile cleanly on all future `v1.x.y` releases.
- If a future property is deprecated, a migration warning is issued at compile time while preserving backwards compatibility.

### 2. Go SDK Public APIs
- The core interfaces and exported functions in `types`, `timeline`, `compiler`, and `filtergraph` will not introduce breaking changes without bumping the `MAJOR` version.
- Internal implementation details remain inside private packages or unexported fields.

---

## Release Schedule

| Release Type | Cadence | Scope |
| :--- | :--- | :--- |
| **Patch (`1.0.x`)** | Bi-weekly / As needed | Bug fixes, platform edge cases, dependency updates |
| **Minor (`1.x.0`)** | Every 4–6 weeks | New video effects, UI workstation features, compiler passes |
| **Major (`2.0.0`)** | Long-term | Architectural paradigm shifts |

---

## Migration Guides & Changelog

When a new version is released, breaking changes or deprecations are highlighted in the [Changelog](/changelog) with explicit migration steps and before/after code snippets.
