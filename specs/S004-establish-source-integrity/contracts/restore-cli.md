# `cueson restore` Contract

## Invocation

```text
cueson [global-options] restore [restore-options] INPUT
```

`INPUT` is one existing Cue JSON file. Standard input is not an S004 input mode.

## Restore Options

| Option | Meaning |
|---|---|
| `--output PATH` | Restore a one-asset document to the literal path. |
| `--output-dir DIR` | Restore one or more assets directly beneath an existing directory using stored basenames. |
| `--force` | Permit replacement of preflighted existing regular files only. |
| `--strict-metadata` | Require every captured requested timestamp to be reproduced or roll back. |
| `--no-metadata` | Intentionally skip all timestamp application and claims. |
| `--help` | Show command help without performing work. |

`--output` and `--output-dir` are mutually exclusive. `--strict-metadata` and `--no-metadata` are mutually exclusive. A multi-asset document requires `--output-dir`. A one-asset document without a destination option uses its stored basename in the current directory.

Existing global `--quiet`, `--silent`, `--no-color`, and help behavior remain available according to the root CLI contract.

## Successful Behavior

Success writes no payload to stdout. Exact bytes are restored directly from `source.assets[].data_base64`, then byte count and SHA-256 are verified. Stderr is empty unless an eligible warning is produced, such as default-mode unsupported metadata.

Exit status is 0.

## Failure Behavior

Invocation errors and deterministic option or destination-mode preconditions return status 2. Input parsing, schema validation, semantic validation, integrity mismatch, restoration, metadata, cancellation, and runtime I/O failures return status 1.

Failures emit one stable `error:` diagnostic to stderr subject to the existing silent-mode rules. Stdout remains empty. No failure may leave an unauthorized destination change or an uncommitted staging artifact.

## Metadata Behavior

Default mode attempts available metadata and warns for unsupported captured values. Strict mode rolls back if any requested captured value is not restored. No-metadata mode skips timestamp application and makes no metadata-fidelity claim.

## Help Behavior

Root help lists `restore`. Restore help documents destination cardinality, overwrite safety, metadata modes, streams, and status codes. Help does not read the input or touch destinations.
