# CLI Contract: `convert`

## Invocation

```text
cueson [global options] convert [options] INPUT --to srt|vtt
```

## Options

```text
--to srt|vtt
--from auto|cueson|srt|vtt
--encoding NAME
-o, --output PATH
-f, --force
--strict
--no-speaker-detection
-q, --quiet
--silent
--no-color
-h, --help
```

`--from` defaults to `auto`. `cueson`, `json`, and `cue-json` select Cue JSON; `subrip` and `webvtt` are accepted aliases for `srt` and `vtt`. Native encoding names and aliases match `encode`. An encoding is invalid with explicit Cue JSON input, and WebVTT accepts only its existing UTF-8 aliases.

## Input classification

Explicit `--from` controls classification. Auto mode first recognizes valid Cue JSON, treats `.json` or `.cueson.json` and JSON-looking invalid content as Cue JSON failures, then applies the existing content-first native registry selection with extension disagreement diagnostics. One bounded no-follow source acquisition supplies native input bytes.

## Output

Without `--output`, and with `--output -`, stdout contains only canonical target subtitle bytes. A real output path uses the established safe transactional publisher. `--force` is valid only with a real file path and authorizes replacement of an existing regular file, not links, directories, devices, missing parents, or unsafe entries.

## Loss and strict behavior

Normal mode emits each ordered loss as a stable-code warning on stderr and publishes parser-valid target bytes. Quiet retains warnings and errors. Silent suppresses informational and loss-warning diagnostics but retains errors. Strict mode computes the complete report and returns exit 1 before rendering or publication when any loss exists; its error identifies the deterministic first loss.

## Exit codes

| Code | Meaning |
|---:|---|
| 0 | Conversion and requested output publication succeeded |
| 1 | Detection, decoding, parsing, Cue JSON validation, source integrity, same-format request, missing capability, strict loss, fatal projection, rendering, or runtime I/O failure |
| 2 | Invalid invocation, missing operand, unknown or conflicting option, invalid explicit selector or encoding, encoding with Cue JSON, incompatible explicit WebVTT encoding, output precondition, or `--force` without a real path |

Same-format conversion returns exit 1 and directs the user to `render`. Conversion never invokes exact restoration.

## Diagnostic order

Source detection and parser diagnostics appear first, followed by the complete ordered conversion-loss warnings, followed by target renderer diagnostics. Standard output contains no commentary. Loss and result diagnostics never include input paths, source bytes, local identifiers, or unbounded payload excerpts; a CLI I/O error may identify the literal caller-supplied path involved in that failure.
