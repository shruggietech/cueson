# Implementation Plan: Deliver Usable SubRip Workflows

**Branch**: `codex/S014-subrip-workflows` | **Date**: 2026-09-11 | **Spec**: [spec.md](spec.md)

## Summary

Deliver the first usable Cueson workflow through a capability-based codec registry, bounded exact source capture, deterministic detection and text decoding, documented SubRip parsing and rendering, and `encode` and `render` commands. Advance the evolving executable and canonical schema to development `0.1.0` while preserving released v0.0.0 artifacts unchanged.

## Technical Context

**Language/Version**: Go 1.25

**Primary Dependencies**: Go standard library, `github.com/santhosh-tekuri/jsonschema/v6`, `golang.org/x/text`

**Storage**: Local source, Cue JSON, and rendered SubRip files

**Testing**: Go unit, integration, conformance, fuzz-seed, golden, and cross-platform tests through `go test ./...`

**Target Platform**: Windows, macOS, and Linux with `CGO_ENABLED=0`

**Project Type**: Single cross-platform CLI with internal packages

**Performance Goals**: Single streamed source capture; reject over 64 MiB; no unbounded parser allocation or malformed-input hangs

**Constraints**: Exact source authority, no local path leakage, no partial output, deterministic LF render output, immutable released v0.0.0 artifacts

**Scale/Scope**: One text codec, two commands, shared foundations, canonical schema evolution, documented corpus, and issues #30/#31

## Constitution Check

*GATE: Passed before Phase 0 research and after Phase 1 design.*

- **Lossless preservation**: PASS. Encode captures exact bounded bytes before parsing; restore stays separate from render.
- **Schema/software discipline**: PASS WITH RECORDED DEVIATION. Changed canonical schema and executable advance to `0.1.0`; immutable v0.0.0 material is untouched.
- **Common model plus native fidelity**: PASS. Common text/timing and SubRip-native observations coexist.
- **No silent loss**: PASS. Source bytes and ordered diagnostics account for tolerated or unmodeled input.
- **Test-first format work**: PASS. Focused and corpus tests precede or accompany implementation.
- **Portable and secure operation**: PASS. Pure Go, limits, safe names, no-follow access, strict Unicode, and safe publication are explicit.
- **Delivery authority**: PASS. Spec Kit and issues govern; push/PR are authorized; merge/release/production remain excluded.

## Project Structure

```text
specs/S014-subrip-workflows/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
├── checklists/
└── tasks.md

cmd/cueson/
internal/
├── cli/
├── codec/
│   └── subrip/
├── model/
├── schema/
├── source/
└── version/
testdata/subrip/
docs/formats/srt.md
```

**Structure Decision**: Preserve the existing single-module CLI architecture. Add internal codec primitives and SubRip without creating a public Go API.

## Complexity Tracking

No constitutional violation or exceptional architecture is required.
