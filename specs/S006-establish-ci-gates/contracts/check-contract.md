# Stable Check Contract

Issue #10 may later require only names proven by successful S006 runs. S006 does not itself configure required checks.

## CI workflow

| Stable check | Evidence |
|---|---|
| `CI / Formatting` | Go formatting is clean and workflow syntax passes actionlint |
| `CI / Repository text` | Publication formatter, UTF-8/BOM/mojibake/line-ending policy, nested formatter tests, whitespace checks, and controlled failure marker absence pass |
| `CI / Vet` | `go vet` passes for the root module |
| `CI / Schema and conformance` | Schema, version lockstep, fixture, path-leak, and cross-package conformance tests pass |
| `CI / Static analysis` | Pinned Staticcheck passes for the root module |
| `CI / Vulnerability scan` | Pinned govulncheck completes without an affecting vulnerability |
| `CI / Native tests (Linux)` | Complete root-module tests pass on `ubuntu-24.04` |
| `CI / Native tests (Windows)` | Complete root-module tests pass on `windows-2025` |
| `CI / Native tests (macOS)` | Complete root-module tests pass on `macos-15` |
| `CI / Race detection` | Root-module tests pass with the race detector on Linux amd64 |
| `CI / Pure-Go build (linux-amd64)` | The executable builds for Linux amd64 with CGO disabled |
| `CI / Pure-Go build (linux-arm64)` | The executable builds for Linux arm64 with CGO disabled |
| `CI / Pure-Go build (windows-amd64)` | The executable builds for Windows amd64 with CGO disabled |
| `CI / Pure-Go build (windows-arm64)` | The executable builds for Windows arm64 with CGO disabled |
| `CI / Pure-Go build (darwin-amd64)` | The executable builds for macOS amd64 with CGO disabled |
| `CI / Pure-Go build (darwin-arm64)` | The executable builds for macOS arm64 with CGO disabled |

## CodeQL workflow

| Stable check | Evidence |
|---|---|
| `CodeQL / Analyze Go` | CodeQL initializes, analyzes Go, and publishes its result |

## Deferred checks

S006 does not create empty check names for native SRT/WebVTT parsing, rendering, cross-format conversion, Codex review automation, release snapshots, or repository controls. Their owning issues establish those checks after implementation exists.
