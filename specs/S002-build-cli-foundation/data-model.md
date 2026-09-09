# Data Model: Build CLI Foundation

S002 persists no domain data. Its internal command model consists of the following transient concepts.

## Invocation

- **Arguments**: Ordered literal tokens supplied after the executable name.
- **Command**: Either root help or the shipped `version` command.
- **Global policy**: Quiet, silent, and no-color selections collected before execution.
- **Option state**: Whether option parsing remains active or has ended after `--`.
- **Help target**: Root or version usage selected by the location of `-h` or `--help`.

### Validation rules

- A non-option token selects at most one command.
- Only `version` is a valid command token.
- `version` accepts no operand.
- Global options are idempotent and may appear on either side of the command while option parsing is active.
- Once `--` is seen, every later token is literal.
- Root-level `--force` and every deferred command are invalid in S002.

## Diagnostic Policy

- **Severity**: Success, informational, warning, or error.
- **Suppression**: Normal, quiet, or silent.
- **Color eligibility**: True only for a terminal diagnostic destination when neither explicit nor environmental color suppression is present.

### Visibility rules

- Normal mode permits every severity.
- Quiet mode suppresses success and informational diagnostics.
- Silent mode suppresses every non-error diagnostic.
- Errors are never suppressed.
- Requested stdout payload is independent from diagnostic suppression.

## Command Result

- **Stdout payload**: Help or the exact version line on success; empty on failure.
- **Stderr diagnostic**: Empty on success; error plus relevant usage for invocation failures; error without usage for a runtime cancellation.
- **Exit status**: 0 for success, 1 for failure after accepting a valid command, or 2 for invocation and pre-execution failure.

## State transitions

1. Parse literal tokens and global policy.
2. On parse failure, emit an invocation diagnostic and relevant usage, then return 2.
3. On help selection, emit the selected help payload and return 0.
4. On a valid command, check cancellation before emitting a success payload.
5. On cancellation, emit a runtime diagnostic and return 1.
6. On `version`, emit the authoritative version and one LF, then return 0.
