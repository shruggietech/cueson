# CLI Foundation Contract

This slice implements the S002 subset of [`docs/cli.md`](../../../docs/cli.md). That ratified document remains authoritative.

## Root invocation

```text
cueson
cueson -h
cueson --help
```

Each form writes root help to stdout, writes nothing to stderr, and exits 0. Root help contains a concise description, usage, the `version` command, the four shipped global options, and worked examples. It contains no deferred command name.

## Version invocation

```text
cueson version
```

The command writes the six bytes represented by `0.0.0\n` to stdout, writes nothing to stderr, and exits 0.

```text
cueson version -h
cueson version --help
```

Each help form writes version-specific usage and a worked example to stdout, writes nothing to stderr, and exits 0.

## Global options

```text
-q, --quiet
--silent
--no-color
-h, --help
```

Quiet, silent, and no-color options may appear before or after `version` while option processing is active. Repeating them is harmless. `--` ends option processing. Help resolves to the root when requested before a command and to the version command when requested after `version`.

## Invalid invocation

Unknown commands, unknown options, and unexpected version operands write no stdout bytes, write a plain error plus relevant usage to stderr, and exit 2. `help`, `schema`, `restore`, `encode`, `render`, `convert`, `validate`, `inspect`, `completion`, and root-level `--force` are unregistered.

## Runtime cancellation

A valid `version` command whose context is already canceled writes no stdout payload, writes a cancellation error to stderr, and exits 1.

## Text and decoration

Generated text is UTF-8 without BOM and uses LF. Non-terminal diagnostics contain no terminal-control sequences. No output contains emoji or a decorative banner.
