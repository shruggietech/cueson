# Implementation Plan: Establish Source Integrity and Restoration

**Branch**: `S004-establish-source-integrity` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/S004-establish-source-integrity/spec.md`

## Summary

Implement issue #6 as a focused source-foundation slice. Extend schema validation with a typed decode path, add `internal/source` for canonical base64, length, SHA-256, basename, collision, capture-ordering, destination-plan, staged transaction, rollback, and timestamp behavior, then expose exact codec-independent restoration through `cueson restore`. A hard-link rollback copy preserves every forced destination until a staged same-directory replacement commits. Native timestamp adapters use handle-based APIs and report restored, unsupported, unavailable, or failed without tolerance-based overclaims. S004 changes the v0.0.0 schema capability to `restore_supported: true` only when the complete command and contract pass.

## Technical Context

**Language/Version**: Go 1.24.0

**Primary Dependencies**: Go standard library; `github.com/santhosh-tekuri/jsonschema/v6` v6.0.3; `golang.org/x/text` v0.31.0 for Unicode NFD and default case folding; `golang.org/x/sys` v0.41.0 for native Windows, Linux, and macOS filesystem APIs

**Storage**: Caller-selected local files only; no database or persistent application state

**Testing**: Go unit, command, transaction-failure, and build-tagged native platform tests; race detection; CGO-disabled builds and cross-compilation

**Target Platform**: Windows, macOS, and Linux on amd64 and arm64, with current-host native Windows execution and issue #8 retaining the complete hosted native matrix

**Project Type**: Pure-Go command-line application with internal domain packages

**Performance Goals**: Validate and restore asset streams in linear time; avoid decoded whole-bundle retention; use bounded chunk buffers; complete every preflight before the first destination is opened

**Constraints**: UTF-8 Cue JSON only; no source paths stored in documents; no codec calls; no unsafe overwrite; no symlink or non-regular replacement; exact hash and byte-count agreement; deterministic diagnostics; no CGO; no CI, fixture-framework, release, tag, or production-domain work

**Scale/Scope**: One Cue JSON input containing one or more ordered assets; asset size is limited only by representable declared length, local filesystem capacity, and the already-loaded JSON document rather than an invented product limit

## Constitution Check

*GATE: Passed before Phase 0 research and after Phase 1 design.*

- **I. Lossless Source Preservation**: Pass. Restoration uses authoritative `data_base64`, validates byte count and SHA-256, stores no runtime path, and never derives bytes from normalized cues.
- **II. Schema and Official Software Discipline**: Pass. The canonical unreleased v0.0.0 schema and representative instance change capability truth in lockstep without changing the software/schema version.
- **III. Common Model Plus Native Fidelity**: Pass. Source operations consume the envelope directly and do not alter common or format-native cue data.
- **IV. No Silent Loss**: Pass. Every metadata request receives an explicit result; unsupported creation restoration is warned or rejected according to mode.
- **V. Test-First Format Work**: Pass. Integrity, restoration, transaction, capture-ordering, and native platform tests precede or accompany implementation. Broad fixture/fuzz infrastructure remains issue #7.
- **VI. Portable and Secure Operation**: Pass. The design is pure Go, treats documents and timestamps as untrusted, refuses links and non-regular targets, preflights the full bundle, stages safely, verifies hashes, and rolls back controlled failures.
- **VII. Documentation and Delivery Authority**: Pass. S004 uses the complete Spec Kit analysis/convergence path, issue #6 remains authoritative, push/PR authorization is explicit, and final merge/release/domain authority remains with the operator.

Post-design re-check: Pass. The deliberate refinement to issue #6 is that S004 creates native-selectable platform tests and runs Windows natively, while issue #8 performs the hosted three-OS matrix. Requiring issue #8's CI before the source and test foundations that unblock issue #8 would be circular. Cross-compilation in S004 proves selection and buildability but is not recorded as native behavioral proof.

## Project Structure

### Documentation (this feature)

```text
specs/S004-establish-source-integrity/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── restore-cli.md
│   └── source-service.md
└── tasks.md
```

### Source Code (repository root)

```text
cmd/cueson/
└── main.go

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
│   └── testdata/representative.cueson.json
├── source/
│   ├── basename.go
│   ├── basename_test.go
│   ├── capture.go
│   ├── capture_test.go
│   ├── integrity.go
│   ├── integrity_test.go
│   ├── plan.go
│   ├── plan_test.go
│   ├── restore.go
│   ├── restore_test.go
│   ├── timestamp.go
│   ├── timestamp_windows.go
│   ├── timestamp_windows_test.go
│   ├── timestamp_linux.go
│   ├── timestamp_linux_test.go
│   ├── timestamp_darwin.go
│   ├── timestamp_darwin_test.go
│   └── timestamp_other.go
└── version/
    └── version.go

docs/
├── architecture.md
├── cli.md
└── schema.md
```

**Structure Decision**: Preserve the ratified single-command repository layout. `internal/source` owns all source and filesystem behavior, `internal/schema` returns validated typed documents, `internal/model` retains format-neutral invariants, and `internal/cli` owns only parsing, streams, diagnostics, and status mapping. Platform-selected files isolate native timestamp APIs without introducing a public Go package.

## Complexity Tracking

No constitution violation requires justification.
