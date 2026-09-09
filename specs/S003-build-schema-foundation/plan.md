# Implementation Plan: Build Schema Foundation

**Branch**: `S003-build-schema-foundation` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/S003-build-schema-foundation/spec.md`

**Note**: This plan is produced through the Spec Kit planning workflow.

## Summary

Establish Cue JSON `0.0.0` as a real Draft 2020-12 contract, not draft prose. The slice adds a fixed canonical schema and representative document under the internal schema boundary, Go types and format-neutral semantic checks under the internal model boundary, schema compilation and software-version lockstep verification, and the public `schema` command with exact stdout/version/file behavior. The schema describes the complete preservation envelope while executable byte-integrity and restoration operations remain unavailable until issue #6.

## Technical Context

**Language/Version**: Go 1.24.0 minimum; local toolchain verification required

**Primary Dependencies**: Go standard library; `github.com/santhosh-tekuri/jsonschema/v6` v6.0.3 for Draft 2020-12 compilation and validation

**Storage**: Read-only embedded JSON artifacts; optional user-selected schema output file

**Testing**: Go `testing`; embedded fixtures; schema compilation and instance validation; focused model and CLI tests; foreground build and smoke checks

**Target Platform**: Windows, macOS, and Linux command-line environments with `CGO_ENABLED=0`

**Project Type**: Command-line application with an internal document model

**Performance Goals**: Compile the embedded schema once per process; validate representative documents deterministically; schema printing scales only with the fixed embedded artifact

**Constraints**: Pure Go; no public Go API; exact UTF-8/LF payloads; stable exit codes; no network resolution during schema compilation; no native codec or restoration claim; safe explicit overwrite; repository text UTF-8 without BOM

**Scale/Scope**: One schema version, two recognized format keys, one representative document, common model and semantic foundation, one new CLI command, and focused tests

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Lossless Source Preservation**: Pass. The schema requires authoritative base64 source bytes and portable asset metadata. S003 does not mutate the source envelope.
- **II. Schema and Official Software Discipline**: Pass. The canonical schema retains its own fixed identity and version while a lockstep check compares it with the executable's sole version source.
- **III. Common Model Plus Native Fidelity**: Pass. Common timing and payload fields coexist with required format-specific data and the source envelope.
- **IV. No Silent Loss**: Pass. Stable object boundaries reject unknown fields, diagnostics are document data, and semantic validation rejects contradictions rather than repairing them.
- **V. Test-First Format Work**: Pass. Schema, model, and CLI contract tests precede or accompany implementation; no native parser or renderer is introduced.
- **VI. Portable and Secure Operation**: Pass. The design is pure Go, treats JSON and output paths as untrusted, blocks implicit replacement, and adds no child process.
- **VII. Documentation and Delivery Authority**: Pass. S003 uses issue #5, Spec Kit, analysis and convergence gates, the Delivery Project, CI, and at most two Codex rounds. Final merge remains human-controlled.
- **Contract boundaries**: Pass. Model rules stay in `internal/model`, embedded-schema operations stay in `internal/schema`, CLI policy stays in `internal/cli`, and source-integrity/restoration work remains in issue #6.
- **Text and delivery rules**: Pass. Repository text follows `.gitattributes`; Markdown uses one logical source line per paragraph and list item; push and PR publication are authorized, while merge, tag, release, and production-domain changes are not.

Post-design re-check: Pass. Research and design artifacts preserve these boundaries. The deliberate interim deviation from the completed-milestone prose is `restore_supported: false`, because the ratified schema baseline explicitly requires that truthful value until issue #6 ships the public restore command.

## Project Structure

### Documentation (this feature)

```text
specs/S003-build-schema-foundation/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── cue-json-v0.0.0.md
│   └── schema-cli.md
├── checklists/
│   └── requirements.md
└── tasks.md
```

### Source Code (repository root)

```text
go.mod
go.sum
internal/
├── cli/
│   ├── cli.go
│   └── cli_test.go
├── model/
│   ├── model.go
│   └── model_test.go
├── schema/
│   ├── cueson.schema.json
│   ├── schema.go
│   ├── schema_test.go
│   └── testdata/
│       └── representative.cueson.json
└── version/
    └── version.go
```

**Structure Decision**: Keep every Go package under `internal/`. The canonical JSON artifact lives beside its embedding package so one tracked file is both the reviewable schema and the compiled payload. `internal/model` owns the typed common contract and cross-field invariants that do not require filesystem access. `internal/schema` owns embedded bytes, offline Draft 2020-12 compilation, raw-document validation, and version lockstep. `internal/cli` gains only the schema invocation and output transaction policy.

## Implementation Phases

### Phase 1: Canonical contract and dependency

Pin the validated Draft 2020-12 implementation, add the canonical schema and representative document, and prove the schema compiles without remote resolution.

### Phase 2: Model and semantic foundation

Define internal types matching the canonical contract and deterministic validation for primary references, identities, ordering, durations, counts, timestamps, format extensions, and interim capability truth.

### Phase 3: Embedded schema service

Embed the exact artifact, expose stable internal identity/version/bytes access, validate raw Cue JSON structurally and semantically, and enforce software/schema lockstep.

### Phase 4: Schema CLI

Extend command parsing and truthful help for `schema`, `schema --version`, and `schema --output PATH`, including explicit force semantics and failure-safe file replacement.

### Phase 5: Documentation and verification

Update public availability statements and the changelog, run the complete quickstart, analyze implementation coverage, and converge any remaining Spec Kit tasks before publication.

## Complexity Tracking

No constitutional violation or exceptional complexity is accepted. A single third-party validator is justified because implementing Draft 2020-12 correctly in-project would be both unsafe and out of scope.
