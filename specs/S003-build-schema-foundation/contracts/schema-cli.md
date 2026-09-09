# Schema CLI Contract

## Supported invocations

```text
cueson schema
cueson schema --version
cueson schema --output PATH
cueson schema --output PATH --force
```

Global quiet, silent, and no-color options retain the S002 placement and payload rules. `-h` and `--help` select schema help. `-o` aliases `--output`; `-f` aliases `--force`.

## Payload behavior

`schema` writes the exact embedded canonical schema bytes to stdout. `schema --version` writes exactly `0.0.0` and one LF. `schema --output PATH` writes the same canonical bytes to the literal path and leaves stdout empty. Successful commands emit no diagnostic.

## Invocation and precondition failures

Missing output values, extra operands, unknown options, `--version` combined with output or force, force without output, and an existing destination without force write an error plus schema usage to stderr, leave stdout empty, and return 2.

## Runtime failures

Canceled execution, stdout failure, temporary-file creation or write failure, close failure, and replacement commit failure write an error to stderr, avoid success output, and return 1. A failed forced replacement retains or restores the prior destination.

Failure to remove a backup after the replacement has already committed produces a warning and returns 0 because the requested destination contains the complete schema; the warning identifies the retained backup for operator cleanup.

## Help truth

Root help lists `version` and `schema` only. Schema help describes its three behaviors and explicit replacement option. Restore, encode, render, convert, validate, inspect, and completion remain absent.
