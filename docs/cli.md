# Cueson CLI Contract

**Status:** Complete v1-bound command contract implemented by v0.1.0 development source, with a published v0.0.0 envelope-only baseline

This document is the maintained CLI authority. The public command names and behavior described here are intended to become stable at v1.0.0, but no v1 binary has been published. The exact generated help under `internal/cli/testdata/help/` and executable documentation tests are checked against this reference.

## Invocation, streams, and status

The root form is `cueson [global options] <command>`. Global options may appear before or after the command until the `--` delimiter ends option processing. User paths are literal; Cueson does not expand globs, tildes, or environment-variable syntax.

- `-q`, `--quiet`: suppress informational and success diagnostics while retaining warnings and errors.
- `--silent`: suppress every non-error diagnostic while retaining errors.
- `--no-color`: disable diagnostic color. `NO_COLOR` is honored, and non-terminal diagnostics are uncolored.
- `-h`, `--help`: print complete explicit help to stdout and exit successfully.
- `--`: end option processing so following operands remain literal.

Stdout contains only a requested payload or explicit help. Diagnostics, warnings, errors, and error-associated short usage use stderr. Quiet and silent never suppress an explicitly requested stdout payload. Cueson emits no decorative banner or emoji. Generated text is UTF-8 without BOM and LF-only unless exact restoration reproduces different preserved source bytes.

- Exit code `0`: the command or explicit help completed successfully; warnings may be present.
- Exit code `1`: an accepted operation failed during processing, cancellation, validation, integrity checking, strict rejection, capability handling, output delivery, or runtime I/O.
- Exit code `2`: command syntax, an operand, an option, or a deterministic pre-execution path or environment requirement was invalid.

Existing output is never replaced implicitly. A filesystem destination requires `--force` before Cueson will replace an existing regular file. `--force` does not permit replacing a directory, symbolic link, device, or other non-regular entry, and it is invalid when output is stdout.

The root help lists exactly nine commands: `encode`, `restore`, `render`, `convert`, `validate`, `inspect`, `schema`, `version`, and `completion`.

## `encode`

```text
cueson [global options] encode [options] INPUT
```

`encode` captures one bounded regular SubRip or WebVTT input, derives common and native cue data, validates the resulting document, and preserves the exact source bytes in Cue JSON.

- `-o`, `--output` `PATH`: write Cue JSON to `PATH`; without an output selection, the destination is `INPUT.cueson.json`. `--output -` selects stdout.
- `-f`, `--force`: replace an approved existing regular filesystem output.
- `--format` `FORMAT`: select `auto`, `srt`, or `vtt`; `subrip` aliases `srt`, and `webvtt` aliases `vtt`.
- `--encoding` `NAME`: select `utf-8`, `utf-8-bom`, `utf-16le`, `utf-16be`, `windows-1252`, or `iso-8859-1`. Accepted aliases are `utf8`; `utf8-bom` and `utf-8-sig`; `utf16le` and `utf-16-le`; `utf16be` and `utf-16-be`; `windows1252` and `cp1252`; and `iso8859-1`, `latin1`, and `latin-1`.
- `--pretty`: indent Cue JSON output.
- `--stdout`: write Cue JSON to stdout; it is mutually exclusive with a filesystem output.
- `--no-speaker-detection`: disable conservative derived `Name:` speaker observations.

WebVTT accepts UTF-8 only, including BOM-specific UTF-8 selection when the input has the matching BOM. SubRip accepts automatic UTF-8, BOM-marked UTF-16, strongly evidenced BOM-less UTF-16, and explicitly selected legacy single-byte encodings. Input is limited to 64 MiB.

For a filesystem output, success exits 0 with empty stdout and only source warnings on stderr. The executable README scenario is:

```text
cueson encode --pretty --output quickstart/document.cueson.json testdata/fixtures/conversion/srt-loss-free/source/input.srt
```

## `restore`

```text
cueson [global options] restore [options] INPUT.cueson.json
```

`restore` validates the source envelope and recreates its exact asset bytes without invoking a codec or rendering the structured model. Cue JSON input is bounded at 1 GiB, and `encode` refuses to publish a larger representation, so every document produced by the CLI remains within the exact-restore boundary.

- `-o`, `--output` `PATH`: restore a single asset to a separately validated literal path.
- `--output-dir` `DIR`: restore all assets beneath an existing directory using their safe stored basenames.
- `-f`, `--force`: replace approved existing regular files after the complete restoration plan passes preflight.
- `--strict-metadata`: require every captured timestamp to be restored; an unsupported value rolls back the operation and exits 1.
- `--no-metadata`: skip timestamp restoration. It is mutually exclusive with `--strict-metadata`.

Without a destination option, one asset uses its safe stored basename in the current directory. `--output` and `--output-dir` are mutually exclusive, and multi-asset documents require `--output-dir`. Parents must already exist. Restoration validates canonical base64, length, SHA-256, safe basenames, portable collisions, and the complete destination plan before publication. It does not claim crash atomicity or cross-process isolation.

Success writes no stdout. Stderr contains only metadata warnings allowed by the global filters. The executable README scenario is:

```text
cueson restore --no-metadata --output quickstart/restored.srt quickstart/document.cueson.json
```

## `render`

```text
cueson [global options] render [options] INPUT.cueson.json --to FORMAT
```

`render` validates Cue JSON and serializes its structured cue model in the document's matching native format. It is intentionally distinct from exact restoration.

- `--to` `FORMAT`: select `srt` or `vtt`; `subrip` aliases `srt`, and `webvtt` aliases `vtt`.
- `-o`, `--output` `PATH`: write native output to a filesystem path instead of stdout; `--output -` selects stdout.
- `-f`, `--force`: replace an approved existing regular filesystem output.
- `--strict`: reject known non-representable, ambiguous, or preserved nonconforming model content before publication.

The target must match the document's native format; use `convert` for another target. Success writes canonical native bytes to stdout unless a filesystem output is selected. Warnings and errors use stderr.

```text
cueson render --to srt --output quickstart/rendered.srt quickstart/document.cueson.json
```

## `convert`

```text
cueson [global options] convert [options] INPUT --to FORMAT
```

`convert` accepts Cue JSON, SubRip, or WebVTT and serializes a private target projection in the requested native format. Cue JSON recognition has precedence; JSON-looking or `.json`-named invalid content does not fall through to a native codec.

- `--to` `FORMAT`: select `srt` or `vtt`; `subrip` aliases `srt`, and `webvtt` aliases `vtt`.
- `--from` `FORMAT`: select `auto`, `cueson`, `srt`, or `vtt`; `json` and `cue-json` alias `cueson`, `subrip` aliases `srt`, and `webvtt` aliases `vtt`.
- `--encoding` `NAME`: for native input, select `utf-8`, `utf-8-bom`, `utf-16le`, `utf-16be`, `windows-1252`, or `iso-8859-1`. Accepted aliases are `utf8`; `utf8-bom` and `utf-8-sig`; `utf16le` and `utf-16-le`; `utf16be` and `utf-16-be`; `windows1252` and `cp1252`; and `iso8859-1`, `latin1`, and `latin-1`.
- `-o`, `--output` `PATH`: write native output to a filesystem path instead of stdout; `--output -` selects stdout.
- `-f`, `--force`: replace an approved existing regular filesystem output.
- `--strict`: compute the complete loss report and reject before rendering or publication when any known loss exists.
- `--no-speaker-detection`: disable derived speaker observations for native input.

Normal conversion writes all representable target bytes and reports every known omission, degradation, or ambiguity on stderr. Loss records have a stable code, severity, kind, source and target format, portable JSON Pointer, and bounded occurrence context; they never enter Cue JSON or expose source bytes or local identity. Fatal target incompatibility fails in normal and strict modes. The complete mapping is maintained in the [conversion contract](conversion.md).

```text
cueson convert --strict --no-speaker-detection --to vtt --output quickstart/converted.vtt testdata/fixtures/conversion/srt-loss-free/source/input.srt
cueson convert --strict --to srt --output quickstart/converted.srt testdata/fixtures/conversion/webvtt-loss-free/source/input.vtt
```

## `validate`

```text
cueson [global options] validate [options] INPUT
```

`validate` accepts Cue JSON, SubRip, or WebVTT and performs the complete applicable validation stack without creating or printing a payload.

- `--format` `FORMAT`: select `auto`, `cueson`, `srt`, or `vtt`; `json` and `cue-json` alias `cueson`, `subrip` aliases `srt`, and `webvtt` aliases `vtt`.
- `--encoding` `NAME`: for native input, select `utf-8`, `utf-8-bom`, `utf-16le`, `utf-16be`, `windows-1252`, or `iso-8859-1`. Accepted aliases are `utf8`; `utf8-bom` and `utf-8-sig`; `utf16le` and `utf-16-le`; `utf16be` and `utf-16-be`; `windows1252` and `cp1252`; and `iso8859-1`, `latin1`, and `latin-1`.

Cue JSON prohibits `--encoding`; WebVTT accepts only compatible UTF-8 selections. Auto mode gives valid Cue JSON precedence, then uses native content evidence and extension evidence. Success exits 0, writes empty stdout, and emits one success diagnostic plus any ordered warnings on stderr. Quiet suppresses success; silent also suppresses warnings.

```text
cueson validate testdata/fixtures/webvtt/minimal/source/minimal.vtt
```

## `inspect`

```text
cueson [global options] inspect [options] INPUT
```

`inspect` uses the same classification and validation depth as `validate`, then emits a privacy-bounded structural report.

- `--format` `FORMAT`: select `auto`, `cueson`, `srt`, or `vtt`; `json` and `cue-json` alias `cueson`, `subrip` aliases `srt`, and `webvtt` aliases `vtt`.
- `--encoding` `NAME`: for native input, select `utf-8`, `utf-8-bom`, `utf-16le`, `utf-16be`, `windows-1252`, or `iso-8859-1`. Accepted aliases are `utf8`; `utf8-bom` and `utf-8-sig`; `utf16le` and `utf-16-le`; `utf16be` and `utf-16-be`; `windows1252` and `cp1252`; and `iso8859-1`, `latin1`, and `latin-1`.
- `--json`: emit one compact deterministic inspection-report-version-1 JSON object plus LF instead of the human report.

Success places the report on stdout and source warnings on stderr. Both modes exclude preserved bytes, content hashes, asset and cue IDs, payload and annotation text, native raw fields, timestamps, free-form diagnostic messages, caller paths, and machine identifiers. Loss is `not_evaluated` with reason `target_format_required` until a target-specific operation is requested.

```text
cueson inspect testdata/fixtures/webvtt/minimal/source/minimal.vtt
cueson inspect --json testdata/fixtures/webvtt/minimal/source/minimal.vtt
```

## `schema`

```text
cueson [global options] schema [options]
```

`schema` emits the byte-identical canonical schema embedded in the executable.

- `--version`: print the embedded schema version followed by one LF.
- `-o`, `--output` `PATH`: write schema bytes to a filesystem path instead of stdout.
- `-f`, `--force`: replace an approved existing regular filesystem output.

`--version` cannot be combined with `--output` or `--force`; `--force` requires a filesystem output. Schema or version bytes use stdout unless a file is selected; errors and error-associated short usage use stderr.

```text
cueson schema
cueson schema --version
cueson schema --output quickstart/cueson.schema.json
```

## `version`

```text
cueson [global options] version
```

`version` accepts no local options or operands. Current development source writes exactly `0.1.0` plus LF to stdout and uses stderr only for errors. The published v0.0.0 binary continues to write `0.0.0`.

```text
cueson version
```

## `completion`

```text
cueson [global options] completion bash|zsh|fish|powershell
```

`completion` accepts exactly one case-sensitive selector: `bash`, `zsh`, `fish`, or `powershell`. It has no local options. Success writes one deterministic UTF-8 LF static definition to stdout and nothing to stderr. Missing, unsupported, or extra selectors exit 2 with empty stdout and an error plus short usage on stderr.

Generated definitions use shell-native registration and never invoke Cueson recursively, execute a subprocess, access a network, inspect subtitle content, infer a shell from environment values, interpret completion candidates as code, or modify a profile.

```text
cueson completion bash
cueson completion powershell
```

## Compatibility and release boundary

The command names, options, aliases, streams, and exit-code classes above form the intended v1 CLI compatibility surface. Go packages remain under `internal/` and are not public APIs. Current source and schema identity remains `0.1.0`; S018 does not create a v1 binary, immutable v1 schema, tag, release, public schema endpoint, or production deployment. See the [compatibility contract](compatibility.md) and [release process](release-process.md).
