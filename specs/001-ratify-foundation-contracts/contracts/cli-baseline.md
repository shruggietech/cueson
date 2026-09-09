# CLI Baseline Contract

## Availability

- Help and command listings expose only implemented commands.
- The executable-foundation slice exposes help and `version`.
- The schema slice adds `schema`, `schema --version`, and schema output behavior.
- The source-foundation slice adds generic `restore` for exact assets from a valid source envelope.
- Encode, render, convert, validate, inspect, and completion enter help only when their implementation slice provides truthful behavior.

## Streams and output

- Explicit help prints to stdout and exits successfully.
- Stdout otherwise contains payload data only.
- Invocation errors, runtime diagnostics, and error-associated usage use stderr.
- Structured output contains no decoration.
- Quiet, silent, color, literal-path, option-termination, and explicit-overwrite conventions apply consistently to shipped commands.

## Failure classification

| Situation | Classification | Exit code |
|---|---|---:|
| Successful command | Success | `0` |
| Unknown or unregistered command | Invocation | `2` |
| Invalid option or missing required operand | Invocation | `2` |
| Missing pre-execution environment requirement | Precondition | `2` |
| Shipped command lacks a required codec | Runtime capability | `1` |
| Unknown input format | Runtime | `1` |
| Strict mode blocks known loss | Runtime | `1` |
| Parse, validation, integrity, render, conversion, or I/O failure | Runtime | `1` |

Diagnostics must distinguish unknown input format, schema-recognized format with a missing required codec, and strict loss rejection.
