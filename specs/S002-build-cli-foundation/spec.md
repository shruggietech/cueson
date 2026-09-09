# Feature Specification: Build CLI Foundation

**Feature Branch**: `S002-build-cli-foundation`

**Created**: 2026-09-09

**Status**: Implemented (local verification complete)

**Input**: User description: "Use Spec Kit to define and complete work slice S002 as the bounded v0.0.0 command-line foundation from issue #4, automatically publish its pull request, satisfy no more than two external review rounds, and return the final merge decision to the operator after continuous integration is green."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Read the exact executable version (Priority: P1)

As an operator or script author, I can ask Cueson for its version and receive a stable payload suitable for direct comparison.

**Why this priority**: A single truthful version is the smallest independently useful executable capability and is required before the schema can be versioned with the software.

**Independent Test**: Invoke `cueson version` and confirm that stdout contains exactly `0.0.0` followed by one line feed, stderr is empty, and the process succeeds.

**Acceptance Scenarios**:

1. **Given** the v0.0.0 executable, **When** a user invokes `cueson version`, **Then** stdout is exactly `0.0.0` followed by one LF, stderr is empty, and the exit code is 0.
2. **Given** output-suppression options, **When** a user invokes the version command with quiet or silent behavior selected, **Then** the explicitly requested version payload is not suppressed or decorated.

---

### User Story 2 - Discover only available behavior (Priority: P2)

As a new user, I can read concise help that lists only behavior the executable actually provides.

**Why this priority**: Truthful discovery prevents users from relying on schema, restoration, or format commands that belong to later work slices.

**Independent Test**: Invoke the executable with no arguments, with root help, and with version help; confirm each succeeds on stdout, documents the available behavior, and contains no unavailable command.

**Acceptance Scenarios**:

1. **Given** no command-line arguments, **When** a user starts Cueson, **Then** root help is printed to stdout and the process exits 0.
2. **Given** `-h` or `--help`, **When** a user requests root help, **Then** the same root help is printed to stdout and the process exits 0.
3. **Given** `version -h` or `version --help`, **When** a user requests command help, **Then** version-specific help and an example are printed to stdout and the process exits 0.
4. **Given** the S002 executable, **When** a user reads root help, **Then** only `version` is listed and later commands such as `schema`, `restore`, `encode`, `render`, `convert`, `validate`, `inspect`, and `completion` are absent.

---

### User Story 3 - Diagnose invalid invocation predictably (Priority: P3)

As a script author, I can distinguish invocation mistakes from successful output through stable streams and exit codes.

**Why this priority**: A deterministic failure boundary lets future commands inherit one command-line contract without changing basic automation behavior.

**Independent Test**: Invoke unknown commands, unknown options, and extra version operands; confirm each returns exit code 2, writes no stdout payload, and writes a plain diagnostic plus relevant usage to stderr.

**Acceptance Scenarios**:

1. **Given** an unknown command, **When** Cueson parses the invocation, **Then** stdout remains empty, stderr identifies the unknown command and shows root usage, and the process exits 2.
2. **Given** an unknown option, **When** Cueson parses the invocation, **Then** stdout remains empty, stderr identifies the option and shows relevant usage, and the process exits 2.
3. **Given** an extra operand for `version`, **When** Cueson parses the invocation, **Then** stdout remains empty, stderr explains that no operand is accepted and shows version usage, and the process exits 2.
4. **Given** `--` before the command name, **When** Cueson parses the invocation, **Then** option processing ends and the following token is interpreted literally as the command name.

### Edge Cases

- Global quiet, silent, and no-color options may appear before or after the shipped command but must not consume or alter its requested payload.
- `NO_COLOR` may be empty or contain a value; its presence disables diagnostic color.
- Diagnostics written to a non-terminal destination contain no terminal-control sequences.
- An option after `--` is treated as a literal command name or operand rather than as an option.
- Repeated global options remain deterministic and do not change the version payload.
- Root-level `--force` is rejected because S002 ships no output-producing command; future commands own their overwrite option when they implement actual output behavior.
- Cancellation before execution completes returns a runtime failure without printing a successful payload.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The project MUST provide a buildable `cueson` executable at version `0.0.0`.
- **FR-002**: The executable MUST maintain one authoritative software-version value used by every version-reporting path.
- **FR-003**: `cueson version` MUST print exactly `0.0.0` followed by one LF to stdout, print nothing to stderr, and exit 0.
- **FR-004**: Starting Cueson without arguments and requesting root help with `-h` or `--help` MUST print root help to stdout and exit 0.
- **FR-005**: Root help MUST provide a concise description, usage, global options, and worked examples while listing only the implemented `version` command.
- **FR-006**: `version -h` and `version --help` MUST print version-specific help to stdout and exit 0.
- **FR-007**: Unknown commands, unknown options, and unexpected operands MUST print no stdout payload, write an error and relevant usage to stderr, and exit 2.
- **FR-008**: The command line MUST reserve exit code 0 for success, 1 for failures after a valid shipped command is accepted, and 2 for invocation or pre-execution failures.
- **FR-009**: The command line MUST reserve stdout for requested payload or help and stderr for diagnostics and error-associated usage.
- **FR-010**: Quiet behavior MUST suppress informational and success diagnostics while retaining warnings and errors; silent behavior MUST suppress every non-error diagnostic; neither behavior may suppress an explicitly requested stdout payload.
- **FR-011**: `--no-color`, the presence of `NO_COLOR`, and non-terminal diagnostic output MUST disable terminal color, and no command output may contain emoji or decorative banners.
- **FR-012**: `--` MUST end option processing, and the command line MUST not expand user-supplied glob, tilde, or environment-variable syntax.
- **FR-013**: S002 MUST NOT register or imply schema, restoration, subtitle-codec, rendering, conversion, validation, inspection, or completion capabilities.
- **FR-014**: S002 MUST NOT expose a root-level overwrite option or add unused overwrite helpers; output-producing commands added later MUST own explicit force behavior and MUST never replace existing output implicitly.
- **FR-015**: Unit and command-level tests MUST cover version identity, help, valid invocation, invalid invocation, stream separation, suppression behavior, option termination, color disabling, and the absence of deferred commands.
- **FR-016**: Repository documentation and the unreleased changelog MUST describe the executable foundation without claiming schema, restoration, native format, continuous integration, or release availability.
- **FR-017**: The completed implementation MUST remain portable, require no native runtime dependency, and avoid exposing a public library contract.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: One hundred percent of successful version invocations produce the exact 6-byte payload `0.0.0\n` and no diagnostic bytes.
- **SC-002**: All documented help entry points succeed, include at least one worked example, and expose zero deferred commands.
- **SC-003**: Every tested invalid invocation produces zero stdout bytes, a non-empty plain stderr diagnostic with relevant usage, and exit code 2.
- **SC-004**: Every executable-foundation test passes with and without race detection where supported, and a clean build succeeds with native dependencies disabled.
- **SC-005**: The Spec Kit checklist and analysis gates report no unresolved blocking item, and repository verification reports zero failure.

## Assumptions

- Issue #4 is the only implementation issue in S002; canonical schema work remains in issue #5 and source restoration remains in issue #6.
- Root invocation without arguments is help, not an invocation error, because no required operation is implied.
- Global suppression and color options may be placed on either side of the shipped command for predictable command-line use.
- Repeating an idempotent global option is accepted.
- S002 can establish the overwrite contract by rejecting irrelevant root-level force options and avoiding dead code until an output-producing command is shipped.
- Continuous integration or review automation that is absent or inapplicable is reported truthfully rather than simulated.
- The operator has explicitly authorized this slice's push, pull-request publication, and at most one second Codex review request, but has retained final merge authority.
