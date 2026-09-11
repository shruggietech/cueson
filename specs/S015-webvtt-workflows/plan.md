# Implementation Plan: Deliver Native WebVTT Workflows

**Branch**: `codex/S015-webvtt-workflows` | **Date**: 2026-09-11 | **Spec**: [spec.md](spec.md)

## Summary

Implement a bounded native WebVTT codec that retains exact source authority, parses the complete documented header, block, cue, setting, markup, entity, and inline-timing surface into common and native model views, and renders deterministic UTF-8 WebVTT through the existing generic CLI workflow. Refine the unreleased `0.1.0` model and schema for raw-line fidelity, ordered setting occurrences, parsed regions, and total cue/non-cue source ordering while preserving SubRip behavior and immutable v0.0.0 release artifacts.

## Technical Context

**Language/Version**: Go 1.25

**Primary Dependencies**: Go standard library, `github.com/santhosh-tekuri/jsonschema/v6`, `golang.org/x/text`

**Storage**: Local WebVTT source, Cue JSON, and rendered WebVTT files

**Testing**: Go unit, integration, conformance, fuzz-seed, golden, and cross-platform tests through `go test ./...`

**Target Platform**: Windows, macOS, and Linux with `CGO_ENABLED=0`

**Project Type**: Single cross-platform CLI with internal packages

**Performance Goals**: One bounded exact source capture, iterative parsing and markup scanning, no work superlinear in the 64 MiB source ceiling, and no malformed-input hangs

**Constraints**: UTF-8-only WebVTT semantic decoding, exact source preservation, total source order, no local path leakage, no partial output, deterministic LF rendering, immutable released v0.0.0 artifacts

**Scale/Scope**: One native text codec, existing encode/render commands, one unreleased schema refinement, complete documented WebVTT matrix, and issue #32

## Constitution Check

*GATE: Passed before Phase 0 research and after Phase 1 design.*

- **Lossless preservation**: PASS. Encode retains the exact captured source asset; semantic NUL replacement, markup parsing, and rendering never rewrite the source envelope.
- **Schema/software discipline**: PASS. The evolving canonical schema remains `0.1.0`; immutable v0.0.0 release files are untouched.
- **Common model plus native fidelity**: PASS. Normalized cue fields coexist with raw header, block, timing, setting, and payload data.
- **No silent loss**: PASS. Ordered setting occurrences, raw non-cue blocks, and stable diagnostics account for tolerated input; strict render refuses unsafe preserved errors.
- **Test-first format work**: PASS. Focused fixtures and tests precede or accompany each codec and integration change.
- **Portable and secure operation**: PASS. The design is pure Go, bounded, iterative, output-safe, path-free, and deterministic across supported platforms.
- **Documentation and delivery authority**: PASS. Spec Kit, issue #32, Project metadata, CI, and the bounded review protocol govern delivery; merge, tags, releases, and production remain excluded.

## Project Structure

```text
specs/S015-webvtt-workflows/
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
│   └── webvtt/
├── conformance/
├── model/
├── schema/
└── source/

testdata/
├── fixtures/webvtt/
├── malformed/webvtt/
└── fuzz/webvtt/

docs/formats/webvtt.md
```

**Structure Decision**: Extend the existing single-module CLI architecture. Keep WebVTT grammar and derived-text logic under a dedicated internal codec, reuse the registry and source boundaries, refine existing internal model/schema types, and avoid a generic subtitle-parser abstraction or public Go API.

## Complexity Tracking

No constitutional violation or exceptional architecture is required.
