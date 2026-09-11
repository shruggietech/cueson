# Implementation Plan: Deliver Bidirectional Subtitle Conversion

**Branch**: `codex/S016-bidirectional-conversion` | **Date**: 2026-09-11 | **Spec**: [spec.md](spec.md)

## Summary

Add source-to-source and Cue-JSON-to-source conversion between SubRip and WebVTT through a new internal conversion boundary. The converter constructs a private target-native document, performs target-aware payload translation, computes a complete stable ordered loss report before rendering, rejects any loss under strict policy, and returns deterministic target bytes without performing file I/O. The CLI reuses bounded source ingest, schema decoding, installed renderers, diagnostics, and transactional publication while leaving Cue JSON schema `0.1.0` unchanged.

## Technical Context

**Language/Version**: Go 1.25

**Primary Dependencies**: Go standard library, existing internal codec/model/schema/source packages, `github.com/santhosh-tekuri/jsonschema/v6`, `golang.org/x/text`

**Storage**: Local SubRip, WebVTT, Cue JSON, and converted subtitle files; no database or network service

**Testing**: Go unit, integration, conformance, fuzz-seed, golden, transaction, race, and cross-platform tests through `go test ./...`

**Target Platform**: Windows, macOS, and Linux with `CGO_ENABLED=0`

**Project Type**: Single cross-platform CLI with internal packages

**Performance Goals**: One bounded source acquisition, linear cue/block and payload scans, deterministic report sorting, and no work superlinear in the existing 64 MiB source ceiling

**Constraints**: Exact input source preservation, no source-envelope restoration during conversion, complete pre-publication loss analysis, target parser acceptance, UTF-8 LF output, no local path leakage, no partial file publication, immutable released v0.0.0 artifacts

**Scale/Scope**: Two conversion directions, two input classes, one new CLI command, a stable internal loss model, complete compatibility matrix, and issue #33

## Constitution Check

*GATE: Passed before Phase 0 research and after Phase 1 design.*

- **Lossless preservation**: PASS. Conversion never mutates the input document or source envelope and never substitutes restoration for semantic conversion.
- **Schema/software discipline**: PASS. Target-specific loss data stays outside Cue JSON, development schema and executable remain `0.1.0`, and immutable v0.0.0 release files remain unchanged.
- **Common model plus native fidelity**: PASS. Common fields drive the semantic bridge while source-native fields determine explicit mapping and loss decisions.
- **No silent loss**: PASS. Every known omission or degradation has one atomic stable loss entry, permissive mode warns, and strict mode rejects before rendering or publication.
- **Test-first format work**: PASS. Compatibility, report, markup, CLI, transaction, parser-cycle, malformed, and fuzz evidence precedes or accompanies implementation.
- **Portable and secure operation**: PASS. The design is pure Go, bounded, non-interactive, target-syntax-aware, output-safe, and deterministic across supported platforms.
- **Documentation and delivery authority**: PASS. Spec Kit, issue #33, Project metadata, current CI, and two-round review limits govern delivery; final merge, tags, releases, and production remain excluded.

## Project Structure

```text
specs/S016-bidirectional-conversion/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
├── checklists/
└── tasks.md

internal/
├── cli/
├── codec/
├── conformance/
├── convert/
├── model/
├── schema/
└── source/

testdata/
├── fixtures/conversion/
├── malformed/conversion/
└── fuzz/conversion/

docs/
├── architecture.md
├── cli.md
├── schema.md
└── formats/
```

**Structure Decision**: Extend the existing single-module CLI architecture with `internal/convert` as the sole owner of cross-format model projection and loss accounting. Reuse concrete native renderers behind the existing codec registry, expose read-only source-envelope integrity validation, and refactor shared CLI source loading only where it prevents divergent encode/convert behavior. Keep conversion results internal and avoid a public Go API or Cue JSON schema expansion.

## Complexity Tracking

No constitutional violation or exceptional architecture is required.
