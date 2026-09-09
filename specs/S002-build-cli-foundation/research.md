# Research: Build CLI Foundation

## Dependency strategy

**Decision**: Implement S002 with the Go standard library and no third-party command framework.

**Rationale**: The slice has one command and a small fixed global-option set. A focused parser can meet the ratified behavior without adding dependency, completion, hidden-command, or generated-help surfaces that are outside S002. The package boundary permits a framework to be evaluated later if the command surface materially grows.

**Alternatives considered**: Cobra or another command framework could reduce future registration work, but it would add a direct dependency and default behaviors that require suppression or customization before the project has more than one command.

## Invocation and help behavior

**Decision**: Treat no arguments, `-h`, and `--help` as successful root-help requests; treat `version -h` and `version --help` as successful command-help requests; reject `help` as an unregistered command.

**Rationale**: This keeps help discoverable while exposing only the forms ratified in `docs/cli.md`. It avoids registering an undocumented help command.

**Alternatives considered**: Returning exit code 2 for no arguments would be defensible for commands that require an operation, but the one-command foundation benefits from a successful discovery path and the specification makes that behavior explicit.

## Global-option placement

**Decision**: Accept `-q`, `--quiet`, `--silent`, and `--no-color` before or after `version`, accept repeats idempotently, and honor `--` as the end of option processing.

**Rationale**: Position-independent global policy is predictable for scripts and future subcommands. Once `--` appears, later dash-prefixed tokens remain literal and can be diagnosed as command names or operands.

**Alternatives considered**: Restricting global options to a prefix simplifies parsing slightly but creates avoidable ordering sensitivity.

## Diagnostic policy

**Decision**: Keep diagnostics behind an internal severity-aware writer. Quiet suppresses informational and success messages, silent suppresses every non-error message, and errors remain visible. Color is enabled only when the diagnostic stream is terminal-like, `--no-color` is absent, and `NO_COLOR` is absent.

**Rationale**: This establishes the ratified policy with a small testable boundary. The shipped `version` path emits no informational diagnostic, so the requested payload remains unchanged.

**Alternatives considered**: Leaving suppression and color as undocumented future work would fail issue #4. Emitting no color under any condition would avoid control sequences but would not establish a meaningful color-policy foundation.

## Overwrite policy

**Decision**: Do not add `--force` or a file-overwrite helper in S002. Reject root-level `--force` as an unknown option and retain the documented rule that a future output-producing command owns explicit force behavior.

**Rationale**: S002 performs no file output. Dead overwrite code would have no valid command path or end-to-end test and could prematurely constrain schema and restoration work.

**Alternatives considered**: A generic output helper could be added speculatively, but issue #5 and issue #6 have different output planning and validation needs.

## Version source

**Decision**: Keep an unexported build-version variable initialized to `0.0.0` in `internal/version` and expose it through a function.

**Rationale**: The package provides one authoritative value and permits future release tooling to inject build identity without duplicating the version across commands.

**Alternatives considered**: A public constant is simpler but cannot be set by later release builds; duplicating a literal in the CLI would violate the single-source requirement.

## Test boundary

**Decision**: Test command behavior in process by calling the internal runner with injected arguments and streams, and verify the thin process adapter through foreground build and smoke commands.

**Rationale**: In-process tests capture exact bytes and exit status deterministically without launching child consoles. This also complies with the Windows no-visible-console rule.

**Alternatives considered**: Subprocess tests exercise `os.Exit` directly but add platform process-launch complexity and could violate the repository's Windows launcher requirements.
