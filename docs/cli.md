# Cueson CLI Contract

**Status:** v0.1.0 development contract with published v0.0.0 baseline

**Ratified:** 2026-09-09 through Spec Kit slice `001-ratify-foundation-contracts`

This document is the CLI authority for current development. Commands enter help and command listings only when the executable implements their documented behavior. The broader command surface in the [working project specification](Cueson-Project-Specification-v0.0.0.md) is a roadmap, not permission to register placeholders.

v0.0.0 is publicly available from the [official GitHub Release](https://github.com/shruggietech/cueson/releases/tag/v0.0.0) as six verified platform archives with checksums and matching SPDX JSON SBOMs. The implemented commands may also be exercised with `go run ./cmd/cueson ...` from a Go 1.25 source checkout. Snapshot archives produced by later non-publishing verification runs remain review evidence unless separately published through an authorized release.

## v0.0.0 command delivery

| Owning issue | Commands added | Capability boundary |
|---|---|---|
| [#4](https://github.com/shruggietech/cueson/issues/4) | Root help, `version` | Executable and CLI foundation |
| [#5](https://github.com/shruggietech/cueson/issues/5) | `schema`, `schema --version`, schema output | Embedded canonical schema |
| [#6](https://github.com/shruggietech/cueson/issues/6) | `restore` | Generic exact source-envelope restoration without a codec |

Current v0.1.0 source additionally ships `encode`, `render`, `convert`, `validate`, `inspect`, and `completion` for the experimental SubRip and WebVTT workflow.

## `encode`

```text
cueson [global options] encode [options] INPUT
```

`encode` detects or explicitly selects SubRip or WebVTT, decodes the exact bounded source bytes, derives common and native cue data, validates the resulting document, and writes Cue JSON to `INPUT.cueson.json` by default. It accepts `--output`, `--force`, `--format auto|srt|vtt`, `--encoding`, `--pretty`, `--stdout`, and `--no-speaker-detection`. `--output -` is equivalent to `--stdout`.

SubRip automatic decoding accepts UTF-8 and BOM-marked UTF-16. BOM-less UTF-16 requires strong byte-pattern evidence. Ambiguous single-byte SubRip input requires an explicit `--encoding windows-1252` or `--encoding iso-8859-1`; aliases shown by command help normalize to the same canonical observations. WebVTT accepts UTF-8 only, with an optional UTF-8 BOM, and rejects an incompatible explicit `--encoding` before execution. Input is limited to 64 MiB. Standard output contains only Cue JSON.

## `render`

```text
cueson [global options] render [options] INPUT.cueson.json --to srt|vtt
```

`render` validates Cue JSON and serializes its structured cue model rather than restoring captured bytes. Without `--output` it writes stdout; `--output -` is equivalent. `--force` applies only to real filesystem destinations. The `--to` value must match the document's native format; use `convert` for a different target format. Canonical SubRip uses ordered integer sequence lines, `HH:MM:SS,mmm`, complete coordinates when present, raw payload text, LF line endings, one blank line between cues, and a final LF. Canonical WebVTT uses a `WEBVTT` signature, preserved header metadata, the contiguous merged cue and non-cue source order, normalized dot-millisecond timestamps, retained native cue settings and payload, LF separators, and a final LF. Normal rendering warns when preserved content is known to be nonconforming; `--strict` rejects known ambiguity or nonconformance before publishing output.

## `convert`

```text
cueson [global options] convert [options] INPUT --to srt|vtt
```

`convert` accepts Cue JSON, SubRip, or WebVTT input and serializes the requested native target from a private target projection. Input format defaults to `auto` and may be selected explicitly with `--from auto|cueson|srt|vtt`. Cue JSON recognition has precedence over native detection; content presented as Cue JSON that fails JSON, schema, semantic, or source-integrity validation does not fall through to a native codec. Native input shares the bounded acquisition, `--encoding`, and `--no-speaker-detection` behavior of `encode`.

Without `--output`, conversion writes the target bytes to stdout; `--output -` is equivalent. `--force` applies only to a real filesystem destination. The converter retains cue order, integer-millisecond timing, overlap, multiline payloads, and shared `b`, `i`, and `u` emphasis. Target-incompatible metadata, placement, native blocks, identifiers, annotations, and markup produce deterministic stderr warnings. Conversion losses are runtime-only observations identified by stable code, severity, kind, source and target format, and a portable JSON Pointer path; they never enter Cue JSON or expose source bytes, input paths, or machine identity. `--strict` computes the complete loss report and fails before target rendering or destination publication when any loss is known. A non-positive SubRip cue duration cannot be represented as valid WebVTT and is fatal in every mode.

## General invocation rules

- `-h` and `--help` are reserved for help.
- Explicit help prints to stdout and exits 0.
- `--` ends option processing.
- User paths are literal. Cueson does not expand globs, tildes, or environment-variable syntax.
- Existing output is never replaced implicitly. `--force` is required.
- `--no-color` disables color, `NO_COLOR` is honored, and non-TTY diagnostics are uncolored.
- No emoji or decorative banner appears in command output.
- Command help provides meaningful descriptions and worked examples only for shipped behavior.

## Streams and suppression

Stdout contains payload data only. Diagnostics, progress, warnings, errors, and error-associated usage use stderr. Structured JSON on stdout contains no commentary or decoration.

`--quiet` suppresses informational and success diagnostics while retaining warnings and errors. `--silent` suppresses every non-error diagnostic. Neither option suppresses command payloads explicitly requested on stdout.

Text generated by Cueson uses UTF-8 without BOM and LF unless the operation is exact restoration of preserved source bytes.

## Exit codes

| Code | Meaning |
|---:|---|
| `0` | Success |
| `1` | Failure after a valid shipped command is accepted |
| `2` | Invocation or pre-execution environment failure |

Exit code 2 covers unknown or unregistered commands, unknown flags, missing required operands, invalid option combinations, and deterministic environment requirements checked before the operation begins.

Exit code 1 covers runtime I/O, parsing, schema or semantic validation, integrity, restoration, rendering, conversion, assertion, strict-loss, and missing-runtime-capability failures.

Higher exit codes require an explicit contract amendment and changelog decision.

## Capability diagnostics

Diagnostics must distinguish these conditions:

- **Unknown format**: the input or requested output does not name a schema-recognized format.
- **Schema-recognized format with missing codec**: the format exists in the contract, but the shipped command lacks native ingest or render capability.
- **Strict loss rejection**: the operation is understood but cannot represent known information without loss under strict policy.

A missing codec for a valid shipped command is a runtime capability failure with exit code 1. It is not an invocation error and does not make the format unknown.

## `version`

```text
cueson version
```

Current development `version` prints exactly `0.1.0` followed by one LF. The published v0.0.0 binary continues to print `0.0.0`.

## `schema`

```text
cueson schema
cueson schema --version
cueson schema --output PATH
```

`schema` prints the embedded canonical schema to stdout unless an output path is selected. Current development `schema --version` prints exactly `0.1.0` followed by one LF. Output replacement follows the explicit `--force` rule. The schema version equals the executable version.

## `restore`

```text
cueson [global options] restore [options] INPUT.cueson.json
```

`restore` accepts Cue JSON, validates its already-populated source envelope, and recreates its exact asset bytes without calling a format codec. It does not accept an SRT or WebVTT file as input, parse native subtitle syntax, derive semantic cues, or render the normalized model. With no destination option, a single asset uses its stored portable safe basename in the current directory. `--output-dir` places one or more assets beneath the caller-selected directory using their stored basenames. `--output` is valid only for a single-asset document and supplies a separately validated literal runtime destination that may rename the file; it is not required to match the stored basename. `--output` and `--output-dir` are mutually exclusive.

Before opening any output, restoration validates every stored basename, builds the full destination plan, and rejects portable basename-key or destination collisions. Overwrite checks apply to the completed plan, so `--force` never permits one source asset to replace another asset from the same bundle.

Baseline options are:

```text
-o, --output PATH
--output-dir DIRECTORY
-f, --force
--strict-metadata
--no-metadata
-q, --quiet
--silent
--no-color
-h, --help
```

Restoration verifies length and SHA-256 before acceptance. Timestamp application occurs after bytes, integrity, and final name are established. Without `--strict-metadata`, unsupported timestamp restoration produces a warning. With strict metadata enabled, a requested timestamp that cannot be reproduced fails with exit code 1 and must not leave a partially accepted artifact.

Successful restoration writes no stdout payload. It emits stderr only for warnings and errors according to the global diagnostic filters. Input parsing, schema validation, semantic validation, encoded-byte integrity, runtime I/O, and metadata application failures use exit code 1. Invalid option combinations and destination preconditions known before execution use exit code 2.

All caller-selected output directories and parents must already exist. Existing symbolic links, directories, devices, and other non-regular entries are refused even with `--force`. The command stages and verifies the complete bundle before publication, preserves forced regular-file destinations for rollback, and does not claim cross-process transaction isolation or crash atomicity.

## `validate`

```text
cueson [global options] validate [options] INPUT
```

`validate` accepts Cue JSON, SubRip, or WebVTT. `--format auto|cueson|srt|vtt` selects the input class and defaults to `auto`; `json` and `cue-json` alias `cueson`, while `subrip` and `webvtt` are native-format aliases. `--encoding` uses the same native text names as `encode`, is invalid with Cue JSON, and accepts only UTF-8 selections for explicitly selected WebVTT.

Auto mode recognizes valid Cue JSON first. JSON-looking or `.json`-named invalid content remains a Cue JSON failure rather than falling through to a native grammar. Cue JSON validation covers UTF-8 JSON parsing, the embedded schema, model semantics, source-envelope integrity, and executable/schema lockstep. Native validation covers bounded capture, content-first selection, decoding, the installed grammar, model semantics, and generated-envelope integrity.

Validation creates no file and writes no stdout payload. Ordered warnings and one successful-validation diagnostic use stderr. Quiet suppresses the success diagnostic, silent also suppresses warnings, and errors remain visible. Invalid invocation and deterministic input-path or size preconditions exit 2; accepted-input parsing, schema, semantic, integrity, capability, cancellation, and runtime failures exit 1.

## `inspect`

```text
cueson [global options] inspect [options] INPUT
cueson [global options] inspect [options] --json INPUT
```

Input classification, `--format`, `--encoding`, and validation depth match `validate`. Default output is one concise human report on stdout. `--json` emits one compact deterministic inspection-report-version-1 object followed by LF, with lowercase `snake_case` keys and no commentary on stdout.

Both modes report input classification, canonical format, schema compatibility, declared and installed capabilities, verified source-integrity status, safe source-asset facts, aggregate document state, structural cue and block summaries, and diagnostic codes and locations. They exclude preserved bytes, content hashes, asset and cue IDs, payload and annotation text, native raw fields, timestamps, free-form diagnostic messages, caller paths, and machine identifiers. Failure diagnostics use stable stage descriptions and do not echo source-integrity details, parser content, caller paths, underlying writer errors, or other potentially sensitive values. Because conversion loss depends on a target, inspection reports it as `not_evaluated` with reason `target_format_required`; this is not a lossless-conversion claim.

## `completion`

```text
cueson [global options] completion bash
cueson [global options] completion zsh
cueson [global options] completion fish
cueson [global options] completion powershell
```

Completion accepts exactly one case-sensitive supported shell selector and writes only one deterministic UTF-8 LF static script to stdout. Generation is non-interactive and does not modify a shell profile. The generated definitions use shell-native completion registration and never invoke Cueson recursively, execute subprocesses, access a network, inspect subtitle content, infer a shell from environment values, or interpret completion candidates as code.

Missing, unsupported, or extra selectors are invocation failures with exit 2 and empty stdout. A valid generation whose stdout write fails exits 1. Quiet, silent, and no-color do not alter the requested script payload.

## Complete development command contract

The v0.1.0 development executable now exposes the complete intended v1 command names: `encode`, `restore`, `render`, `convert`, `validate`, `inspect`, `schema`, `version`, and `completion`. Completion of the command surface does not by itself declare format stability or release readiness; the remaining v1 hardening, end-to-end documentation, release-candidate, publication, and release work retain their own acceptance gates.

Release packaging does not expand this command contract. See [release verification](release-verification.md) for the v0.0.0 artifact and publication proof and the [release process](release-process.md) for the separately authorized release lifecycle. The standalone documentation verifier checks that these maintained CLI and documentation links resolve offline; it is repository tooling, not a `cueson` subcommand.
