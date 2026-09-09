# Tasks: Build CLI Foundation

**Input**: Design documents from `/specs/S002-build-cli-foundation/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: S002 requires test-first unit and command-level coverage.

**Organization**: Tasks are grouped by independently testable user story. Publication, external review handling, and final merge authority remain delivery-protocol actions outside the product implementation checklist so review completion does not require a checkbox-only code push.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it owns a different file and has no incomplete dependency.
- **[Story]**: Maps the task to a user story in `spec.md`.
- Every task names an exact repository path or authoritative GitHub record.

## Phase 1: Setup

**Purpose**: Establish the product module without coupling it to bootstrap tooling.

- [X] T001 Create the root Go 1.24.0 module for `github.com/shruggietech/cueson` in `go.mod`
- [X] T002 Verify product build artifacts and nested publication-tooling output remain excluded in `.gitignore`

---

## Phase 2: Foundational Version Boundary

**Purpose**: Create the single version source required by every shipped command.

- [X] T003 Write failing tests for the authoritative `0.0.0` value and stable accessor in `internal/version/version_test.go`
- [X] T004 Implement the linker-replaceable single version source in `internal/version/version.go`

**Checkpoint**: Version identity is independently verified and ready for command use.

---

## Phase 3: User Story 1 - Read the exact executable version (Priority: P1) MVP

**Goal**: Provide the exact scriptable version payload through the executable boundary.

**Independent Test**: The in-process command runner returns 0, writes exactly `0.0.0\n` to stdout, and writes no stderr for every valid version invocation.

### Tests for User Story 1

- [X] T005 [US1] Write failing exact-output, suppression, repeated-option, and canceled-context command tests in `internal/cli/cli_test.go`

### Implementation for User Story 1

- [X] T006 [US1] Implement the command runner, version dispatch, stream injection, and exit-status constants in `internal/cli/cli.go`
- [X] T007 [US1] Add the minimal context, argument, stream, and process-exit adapter in `cmd/cueson/main.go`

**Checkpoint**: `cueson version` is independently buildable and testable.

---

## Phase 4: User Story 2 - Discover only available behavior (Priority: P2)

**Goal**: Provide truthful root and version help without exposing deferred commands.

**Independent Test**: No-argument, root-help, and version-help invocations succeed on stdout; help contains worked examples and only the shipped command.

### Tests for User Story 2

- [X] T008 [US2] Write failing root-help, version-help, example, and deferred-command exclusion tests in `internal/cli/cli_test.go`

### Implementation for User Story 2

- [X] T009 [US2] Implement stable root and version help payloads and help dispatch in `internal/cli/cli.go`

**Checkpoint**: Every supported discovery path is truthful and independently testable.

---

## Phase 5: User Story 3 - Diagnose invalid invocation predictably (Priority: P3)

**Goal**: Establish stable option parsing, diagnostics, suppression, color, and invocation failure behavior.

**Independent Test**: Unknown commands, unknown options, literal post-`--` tokens, and unexpected operands return 2 with empty stdout and relevant stderr, while an already-canceled valid command returns 1.

### Tests for User Story 3

- [X] T010 [US3] Write failing invalid-invocation, option-termination, option-placement, stream-separation, suppression, and color-policy tests in `internal/cli/cli_test.go`

### Implementation for User Story 3

- [X] T011 [US3] Implement literal option parsing, diagnostic filtering, terminal-color eligibility, invocation errors, and usage selection in `internal/cli/cli.go`

**Checkpoint**: The complete S002 command contract is independently testable.

---

## Phase 6: Documentation and Verification

**Purpose**: Align public status and execute every local completion gate.

- [X] T012 [P] Update executable availability, current commands, and development commands in `README.md`
- [X] T013 [P] Record the S002 executable foundation under Unreleased in `CHANGELOG.md`
- [X] T014 Run every command and repository gate in `specs/S002-build-cli-foundation/quickstart.md`, including formatting, vet, unit tests, race tests, clean `CGO_ENABLED=0` build, smoke invocations, publication-tool tests, whitespace, encoding, mojibake, and Spec Kit readiness

---

## Dependencies & Execution Order

### Phase Dependencies

- Setup has no dependency.
- Foundational Version Boundary depends on Setup and blocks all user stories.
- User Story 1 depends on the version boundary and is the MVP.
- User Story 2 builds on the shared command runner from User Story 1.
- User Story 3 builds on the runner and help surfaces so its errors can select relevant usage.
- Documentation and Verification depends on all user stories.

### User Story Dependencies

- **User Story 1 (P1)**: Depends only on the foundational version boundary.
- **User Story 2 (P2)**: Depends on the User Story 1 runner but remains independently verifiable through help invocations.
- **User Story 3 (P3)**: Depends on the runner and help payloads but remains independently verifiable through failure invocations.

### Parallel Opportunities

- README and changelog updates can proceed in parallel after behavior is stable.
- Version tests and implementation are intentionally sequential for TDD.
- User-story tests and implementations are intentionally sequential because they share `internal/cli/cli_test.go` and `internal/cli/cli.go`.

## Parallel Example: Documentation

```text
Task T012: Update README.md with truthful executable and development status.
Task T013: Update CHANGELOG.md with the Unreleased S002 addition.
```

## Implementation Strategy

### MVP First

1. Complete Setup and the Foundational Version Boundary.
2. Write the User Story 1 tests and observe their expected failure.
3. Implement the command runner and operating-system adapter.
4. Validate the exact version payload independently.

### Incremental Delivery

1. Add truthful help after the version MVP.
2. Add the complete invalid-invocation and diagnostic policy.
3. Align documentation and run all local gates.
4. Continue through the separately authorized GitHub publication and bounded review protocol without changing final merge authority.

## Notes

- Completed tasks are marked `[X]` during the implementation phase.
- Tests precede the corresponding implementation within every behavior phase.
- S002 adds no schema, source, codec, rendering, conversion, validation, inspection, completion, CI, release, or production-domain implementation.
