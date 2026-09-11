# Implementation Plan: Complete CLI Workflows

**Branch**: `codex/S017-complete-cli-workflows` | **Date**: 2026-09-11 | **Spec**: [spec.md](spec.md)

## Summary

Complete the v1 CLI command surface by adding assertion-only validation, privacy-bounded human and JSON inspection, and static completion for Bash, Zsh, Fish, and PowerShell. A single internal validated-input loader reuses current bounded capture, Cue JSON precedence, schema and semantic validation, source-integrity validation, codec selection, and native decoding. An ordered CLI surface catalogue drives help and completion vocabulary without replacing the existing semantic parser. Inspection remains a private CLI report contract and does not change Cue JSON or publish source content.

## Technical Context

**Language/Version**: Go 1.25

**Primary Dependencies**: Go standard library, existing internal codec/model/schema/source packages, `github.com/santhosh-tekuri/jsonschema/v6`, `golang.org/x/text`

**Storage**: Read-only local Cue JSON, SubRip, and WebVTT inputs; completion and inspection output use stdout; no database or network service

**Testing**: Go unit, golden, process, conformance, fuzz-seed, race, cross-platform, and completion-syntax smoke tests through `go test ./...`

**Target Platform**: Windows, macOS, and Linux with `CGO_ENABLED=0`; generated completion targets Bash 3.2+, Zsh 5.x, Fish 3.x, Windows PowerShell 5.1, and PowerShell 7+

**Project Type**: Single cross-platform CLI with internal packages

**Performance Goals**: One bounded source acquisition, linear validation and inspection passes within the existing 64 MiB source ceiling, deterministic ordered output, and no subprocess for normal completion generation

**Constraints**: No source bytes or local identifiers in inspection output, no validation output artifacts, no duplicate format logic, fixed snake_case JSON shape, static non-interactive completion, accurate current-command help, immutable v0.0.0 release artifacts

**Scale/Scope**: Three new commands, three input classes, four completion shells, nine complete command help surfaces, one inspection report version, and issue #34

## Constitution Check

*GATE: Passed before Phase 0 research and after Phase 1 design.*

- **Lossless preservation**: PASS. Validation and inspection are read-only, use the retained source envelope only for integrity proof, and never alter or republish source bytes.
- **Schema/software discipline**: PASS. Cue JSON validation uses the embedded schema and lockstep check; report and completion contracts remain outside Cue JSON and version `0.1.0` is unchanged.
- **Common model plus native fidelity**: PASS. Native parsing derives the existing common and native model, while inspection reports only bounded structural facts and never replaces source truth.
- **No silent loss**: PASS. Validation surfaces ordered native diagnostics, inspection explicitly labels conversion loss `not_evaluated`, and neither command makes a lossless-conversion claim.
- **Test-first format work**: PASS. Loader, validation, report, privacy, completion, help, process, and cross-platform tests precede or accompany implementation.
- **Portable and secure operation**: PASS. The design is pure Go, bounded, path-free in report output, deterministic, non-interactive, and contains no completion-time subprocess or network behavior.
- **Documentation and delivery authority**: PASS. Spec Kit, issue #34, Project metadata, current CI, and the two-round review limit govern delivery; merge, version stability, schema publication, tags, releases, and production remain excluded.

## Project Structure

```text
specs/S017-complete-cli-workflows/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
├── checklists/
└── tasks.md

internal/cli/
├── cli.go
├── workflows.go
├── input.go
├── validate.go
├── inspect.go
├── surface.go
├── help.go
├── completion.go
├── *_test.go
└── testdata/

docs/
├── architecture.md
├── cli.md
├── schema.md
└── formats/
```

**Structure Decision**: Keep every new public contract in `internal/cli`. Extract one neutral validated-input loader from the existing conversion workflow, create a versioned private inspection report and deterministic human renderer, and introduce one ordered CLI surface catalogue shared by help and completion generation. Retain the existing hand-written semantic parser because replacing it with a framework would be disproportionate; catalogue invariant tests prevent parser, help, and completion drift. No `internal/model` or Cue JSON schema changes are required.

## Complexity Tracking

No constitutional violation or exceptional architecture is required.
