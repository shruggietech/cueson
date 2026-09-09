# Tasks: Establish Source Integrity and Restoration

**Input**: Design documents from `/specs/S004-establish-source-integrity/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`, and `quickstart.md`

**Tests**: Required by the feature specification and constitution. Each story's focused tests are written before or with its implementation and must fail for the missing behavior before they pass.

**Organization**: Tasks are chronological and grouped by user story. Story labels preserve requirement traceability, and `[P]` marks work on independent files that can proceed concurrently after its phase prerequisites.

## Phase 1: Setup

**Purpose**: Establish the package and dependency boundary without changing product behavior.

- [x] T001 Add the Go 1.24-compatible `golang.org/x/text` and `golang.org/x/sys` direct dependencies in `go.mod` and `go.sum`
- [x] T002 Create the internal source package API, result types, metadata modes, and typed precondition errors in `internal/source/source.go`

---

## Phase 2: Foundational Validation

**Purpose**: Complete the typed document boundary and shared model capability before source or CLI behavior.

- [x] T003 Add invalid-UTF-8 and typed decode tests in `internal/schema/schema_test.go`
- [x] T004 Implement the validated typed `Decode` boundary and delegate `Validate` to it in `internal/schema/schema.go`
- [x] T005 Update model capability invariants and tests for completed generic restoration in `internal/model/model.go` and `internal/model/model_test.go`

**Checkpoint**: Source work receives only structurally and semantically valid typed documents.

---

## Phase 3: User Story 1 - Reject unsafe or corrupted source envelopes (Priority: P1)

**Goal**: Validate every source asset and the complete destination plan before output begins.

**Independent Test**: Corrupt base64, length, digest, basename, bundle collision, destination collision, and non-regular-target cases fail without creating or changing output.

### Tests for User Story 1

- [x] T006 [P] [US1] Add canonical base64, length, digest, and full-bundle preflight tests in `internal/source/integrity_test.go`
- [x] T007 [P] [US1] Add safe-basename, normalization, case-folding, and portable-collision tests in `internal/source/basename_test.go`
- [x] T008 [P] [US1] Add destination-mode, complete-plan, race-boundary, and non-regular-target tests in `internal/source/plan_test.go`

### Implementation for User Story 1

- [x] T009 [P] [US1] Implement safe basename validation and NFD/default-fold/NFD portable identities in `internal/source/basename.go`
- [x] T010 [P] [US1] Implement canonical streaming base64, byte-count, and SHA-256 preflight validation in `internal/source/integrity.go`
- [x] T011 [US1] Implement complete destination planning and no-follow existing-target validation in `internal/source/plan.go`
- [x] T012 [US1] Run the focused User Story 1 tests and prove rejected inputs cause no output mutation

**Checkpoint**: Every source and destination decision is complete before staging.

---

## Phase 4: User Story 2 - Restore exact source bytes without a codec (Priority: P2)

**Goal**: Restore one or many assets through a rollback-safe staged transaction and expose it through the public command.

**Independent Test**: Each destination mode reproduces declared bytes and hash, unsafe overwrites fail, force replaces only approved regular files, injected failures roll back, and successful stdout remains empty.

### Tests for User Story 2

- [x] T013 [P] [US2] Add staged-write, publish, force-backup, rollback, cleanup, cancellation, and post-write verification tests in `internal/source/restore_test.go`
- [x] T014 [P] [US2] Add restore parsing, streams, diagnostics, help, option conflict, and status-code tests in `internal/cli/cli_test.go`

### Implementation for User Story 2

- [x] T015 [US2] Implement same-directory staging, exclusive publication, force backups, identity-safe rollback, cleanup warnings, and final integrity verification in `internal/source/restore.go`
- [x] T016 [US2] Add `restore` parsing, help, validated document loading, diagnostic filtering, and status mapping in `internal/cli/cli.go`
- [x] T017 [US2] Run focused transaction and command tests, including all documented destination modes and failure cleanup

**Checkpoint**: Exact codec-independent restoration works end to end with safe transaction semantics.

---

## Phase 5: User Story 3 - Preserve and report filesystem timestamps truthfully (Priority: P3)

**Goal**: Capture timestamps before content reads and apply only platform-supported metadata with explicit outcomes.

**Independent Test**: Native Windows tests prove capture ordering, restoration, and readback; platform-selected Linux and macOS sources compile and define truthful birth/creation behavior; default, strict, and no-metadata modes behave distinctly.

### Tests for User Story 3

- [x] T018 [P] [US3] Add source capture ordering, safe-name, byte-integrity, and metadata-mode tests in `internal/source/capture_test.go` and `internal/source/restore_test.go`
- [x] T019 [P] [US3] Add Windows handle timestamp conversion, capture, application, and exact-readback tests in `internal/source/timestamp_windows_test.go`
- [x] T020 [P] [US3] Add Linux descriptor timestamp and unsupported birth-setting tests in `internal/source/timestamp_linux_test.go`
- [x] T021 [P] [US3] Add macOS descriptor timestamp and truthful creation-setting tests in `internal/source/timestamp_darwin_test.go`

### Implementation for User Story 3

- [x] T022 [US3] Define common timestamp capture, application, readback, and result orchestration in `internal/source/timestamp.go`
- [x] T023 [P] [US3] Implement no-follow Windows capture and FILETIME restoration in `internal/source/timestamp_windows.go`
- [x] T024 [P] [US3] Implement no-follow Linux capture and descriptor-based atime/mtime restoration in `internal/source/timestamp_linux.go`
- [x] T025 [P] [US3] Implement no-follow macOS capture and descriptor-based atime/mtime restoration in `internal/source/timestamp_darwin.go`
- [x] T026 [P] [US3] Add explicit unsupported-platform behavior in `internal/source/timestamp_other.go`
- [x] T027 [US3] Implement capture-before-read and integrate metadata application, verification, warning, strict rollback, and no-metadata behavior in `internal/source/capture.go` and `internal/source/restore.go`
- [x] T028 [US3] Run native Windows timestamp and metadata-mode tests, then compile Linux and macOS platform-selected test sources with CGO disabled

**Checkpoint**: Every requested timestamp has one truthful result and byte fidelity remains independently verifiable.

---

## Phase 6: Capability, Documentation, and Full Verification

**Purpose**: Publish truthful v0.0.0 capability state, document the command, and converge the complete work slice.

- [x] T029 Update `restore_supported` to true and keep all other capability flags false in `internal/schema/cueson.schema.json`, `internal/schema/testdata/representative.cueson.json`, and schema tests
- [x] T030 [P] Document source boundaries, restore usage, metadata limits, and issue #7/#8 exclusions in `docs/architecture.md`, `docs/cli.md`, and `docs/schema.md`
- [x] T031 [P] Add the unreleased S004 entry to `CHANGELOG.md`
- [x] T032 Run formatting, mojibake, diff, unit, race, vet, native build, schema, command smoke, and foreign-platform compile verification from `quickstart.md`
- [x] T033 Run Spec Kit convergence, record any findings as append-only remediation tasks, resolve them, and repeat verification until clean

---

## Dependencies and Execution Order

- Phase 1 precedes all implementation.
- Phase 2 precedes all source preparation and restoration work.
- User Story 1 precedes User Story 2 because restoration consumes its validated assets and destination plan.
- User Story 2 precedes final User Story 3 integration because timestamp application occurs within the restoration transaction.
- Platform adapter implementations T023-T026 can proceed in parallel after common timestamp contract T022.
- Capability and documentation changes follow completed behavior so they cannot overclaim support.
- Final convergence follows every implementation and verification task.

## Requirement Coverage

- FR-001: T003-T004
- FR-002-FR-005: T006-T012
- FR-006-FR-016: T008, T013-T017
- FR-017-FR-023: T018-T028
- FR-024: T005, T029
- FR-025: T030-T031
- SC-001-SC-007: T006-T031
- SC-008: T028, T032
- SC-009: T033 plus the authorized pull-request review protocol

## Implementation Strategy

Complete each chronological checkpoint before advancing. Keep commits aligned to coherent foundations, integrity and planning, transactional restore, timestamp behavior, and final documentation. Do not add issue #7 fixture infrastructure, issue #8 workflows, native subtitle codecs, releases, tags, or production-domain changes.

## Phase 7: Convergence

- [x] T034 Add deterministic transaction-operation fault seams and prove write, sync, close, publication, replacement, final-verification, rollback, and cleanup behavior per SC-004 (partial)
- [x] T035 Reject distinct planned paths that resolve to the same existing native file identity and add hard-link alias coverage per FR-005 and FR-010 (partial)

## Phase 8: Convergence

- [x] T036 Classify missing or non-regular inputs, output directories, and output parents as deterministic pre-execution failures with command coverage per FR-015 (partial)
- [x] T037 Make staging and final-verification streams context-aware and prove cancellation cleans partial transaction artifacts per FR-011 and the cancellation edge case (partial)

## Phase 9: Convergence

- [x] T038 Replace the README's unavailable-restoration claim with truthful S004 status, usage, and remaining codec or CI boundaries per FR-025 and Constitution VII (contradicts)

## Phase 10: Convergence

- [x] T039 Mark the S004 feature specification implemented after local verification per the Spec Kit lifecycle record (partial)

## Phase 11: Round-One Review Remediation

- [x] T040 Reject unpaired UTF-16 surrogate escapes before typed JSON decoding and add paired or unpaired coverage per FR-001 and review thread `PRRT_kwDOUTHZE86grJj2`
- [x] T041 Join failed-transaction staging cleanup errors into the returned failure and add deterministic cleanup-fault coverage per FR-011, SC-004, and review thread `PRRT_kwDOUTHZE86grJj-`
- [x] T042 Classify destinations that appear or change after planning as runtime transaction failures and preserve raced-in paths per FR-015 and review thread `PRRT_kwDOUTHZE86grJkC`
- [x] T043 Downgrade equal macOS birth and change timestamps to explicit `ctime_fallback` provenance with platform-selected coverage per FR-018 and review thread `PRRT_kwDOUTHZE86grJkF`
- [x] T044 Reject whitespace-only literal output paths and output directories before resolution or filesystem mutation per the path edge case and review thread `PRRT_kwDOUTHZE86grJkL`
- [x] T045 Reject Windows superscript COM and LPT device aliases in schema and source validation with portable coverage per FR-004 and review thread `PRRT_kwDOUTHZE86grJkP`
