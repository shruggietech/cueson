# CLI Contract: Validation, Inspection, Completion, and Help

## Validation

```text
cueson [global options] validate [options] INPUT
```

Options are `--format auto|cueson|srt|vtt`, `--encoding NAME`, and help. `--format` defaults to `auto`; `json` and `cue-json` alias `cueson`, `subrip` aliases `srt`, and `webvtt` aliases `vtt`. Encoding names and aliases match `encode`. Cue JSON prohibits `--encoding`; explicit WebVTT accepts only UTF-8 selections.

Validation performs one bounded no-follow acquisition. Cue JSON passes JSON parsing, canonical schema, semantic model, source-integrity, and executable/schema lockstep checks. Native input passes content-first format selection, decoding, grammar parsing, generated model semantics, and generated source-envelope integrity. Native validation disables heuristic speaker derivation because it does not affect grammar validity.

Validation never writes stdout or a file. Ordered warnings and one success diagnostic use stderr. Quiet suppresses success, silent suppresses success and warnings, and errors are never suppressed.

## Inspection

```text
cueson [global options] inspect [options] INPUT
```

Options are `--format auto|cueson|srt|vtt`, `--encoding NAME`, `--json`, and help. Classification and validation match `validate`. Native inspection retains default derived speaker observations only as structural counts.

Default output is one deterministic human report on stdout. `--json` outputs one compact deterministic report-version-1 JSON object plus one LF. Neither mode emits input paths, source bytes, raw fields, content hashes, user-controlled identifiers or messages, or machine identity. Parser warnings remain on stderr according to global filters.

## Completion

```text
cueson [global options] completion SHELL
```

`SHELL` is exactly one case-sensitive value: `bash`, `zsh`, `fish`, or `powershell`. Success writes only one static UTF-8 LF completion script to stdout. It does not modify a profile. Unsupported, missing, or extra values are invocation failures and write no stdout payload.

Generated scripts use only shell-native completion registration and already-tokenized command state. They never recursively execute Cueson, run a subprocess, access a network, inspect an input file, infer a shell from environment values, interpret a candidate as code, or mutate persistent configuration.

## Input classification

Explicit format selection is authoritative and never falls through. Auto mode accepts valid Cue JSON first. JSON-looking or `.json`-named invalid content remains a Cue JSON failure. Otherwise the installed registry uses content evidence first, registered extension evidence second, and reports deterministic disagreement warnings.

## Shared global options

```text
-q, --quiet
--silent
--no-color
-h, --help
--
```

Global options may appear before or after the command until `--` ends option processing. Paths are literal. `NO_COLOR` and non-terminal diagnostics disable color.

## Streams

| Operation | Standard output | Standard error |
|---|---|---|
| Explicit help | Complete help text | Empty |
| Validation success | Empty | Success and warnings allowed by filters |
| Inspection success | Human report or one JSON object | Source warnings allowed by filters |
| Completion success | Static script | Empty |
| Invalid invocation | Empty | Error and short usage |
| Runtime failure | Empty unless a stream failed after partial write | Error diagnostic when writable |

Existing command payload rules remain unchanged. Quiet and silent never suppress an explicitly requested stdout payload.

## Exit codes

| Code | Meaning |
|---:|---|
| 0 | Command or explicit help completed successfully; warnings may be present |
| 1 | Accepted operation failed during classification, parsing, schema, semantics, integrity, capability, cancellation, output-stream, or runtime I/O |
| 2 | Command syntax, selector, option, operand, or deterministic pre-execution path requirement is invalid |

## Help completeness

Root help lists all nine shipped commands. Every command help page contains its description, usage, complete local options, the shared global options, stream behavior, exit-code behavior, and at least one valid example. Short error usage remains on stderr and does not replace full explicit help.
