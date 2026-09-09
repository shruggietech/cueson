# Implementation Plan: Build CLI Foundation

**Branch**: `S002-build-cli-foundation` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/S002-build-cli-foundation/spec.md`

**Note**: This plan is produced through the Spec Kit planning workflow.

## Summary

Create the first shipped Cueson executable as a dependency-free Go command with a minimal operating-system entry point, an internal command runner, and one internal version source. The slice provides truthful root and version help, exact `0.0.0` output, stable stream and exit-code behavior, suppression and color policy foundations, and focused in-process command tests while deliberately omitting every later command.

## Technical Context

**Language/Version**: Go 1.24.0 minimum; verified with Go 1.24.2

**Primary Dependencies**: Go standard library only

**Storage**: N/A

**Testing**: Go `testing`; in-process command invocation with injected streams; foreground build and smoke checks

**Target Platform**: Windows, macOS, and Linux command-line environments with `CGO_ENABLED=0`

**Project Type**: Command-line application

**Performance Goals**: Version and help perform bounded work with no input-size-dependent allocation

**Constraints**: Pure Go; no public Go API; stdout payload isolation; diagnostic stderr; exit codes limited to 0, 1, and 2; UTF-8 without BOM; LF output; no placeholder commands; no native child-process launch

**Scale/Scope**: Root help, `version`, global quiet/silent/no-color parsing, error handling, and their tests only

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Lossless Source Preservation**: Pass. S002 reads or writes no subtitle or Cue JSON source asset.
- **II. Schema and Official Software Discipline**: Pass. One software-version source reports `0.0.0`; schema identity remains owned by issue #5.
- **III. Common Model Plus Native Fidelity**: Pass. S002 adds no common model or format interpretation.
- **IV. No Silent Loss**: Pass. Invalid invocation produces explicit diagnostics and stable failure status.
- **V. Test-First Format Work**: Pass. No format work occurs; CLI behavior receives focused tests before or alongside implementation.
- **VI. Portable and Secure Operation**: Pass. The executable uses only the standard library and builds with `CGO_ENABLED=0`.
- **VII. Documentation and Delivery Authority**: Pass. S002 uses Spec Kit, issue #4, the Delivery Project, analysis, bounded review rounds, and operator-only final merge.
- **Contract boundaries**: Pass. `cmd/cueson` remains an operating-system adapter; command policy and version identity remain under `internal/`.
- **Text and delivery rules**: Pass. Repository text follows `.gitattributes`; Markdown remains one logical line per paragraph and list item; the authorized push and pull-request publication do not extend to merge, tag, release, or production changes.

Post-design re-check: Pass. The design artifacts retain the same boundaries and introduce no constitutional exception.

## Project Structure

### Documentation (this feature)

```text
specs/S002-build-cli-foundation/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── cli-foundation.md
├── checklists/
│   └── requirements.md
└── tasks.md
```

### Source Code (repository root)

```text
go.mod
cmd/
└── cueson/
    └── main.go
internal/
├── cli/
│   ├── cli.go
│   └── cli_test.go
└── version/
    ├── version.go
    └── version_test.go
```

**Structure Decision**: Use the ratified single-command layout. `cmd/cueson` supplies only context, process streams, arguments, and the returned exit status. `internal/cli` owns parsing, help, diagnostics, suppression, color decisions, and status mapping. `internal/version` owns the sole software-version value. Tests remain adjacent to these internal boundaries because no public package or separate black-box harness is justified for this slice.
