# Tasks: Build Schema Foundation

**Input**: Design documents from `/specs/S003-build-schema-foundation/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: S003 requires test-first schema, semantic, lockstep, and command-level coverage.

**Organization**: Tasks are grouped by independently testable user story. Publication, review handling, and final merge authority remain delivery-protocol actions outside the implementation checklist.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it owns a different file and has no incomplete dependency.
- **[Story]**: Maps the task to a user story in `spec.md`.
- Every task names an exact repository path.

## Phase 1: Setup

**Purpose**: Establish the standards-compliant validation dependency and schema artifact locations.

- [X] T001 Pin `github.com/santhosh-tekuri/jsonschema/v6` v6.0.3 and record its checksums in `go.mod` and `go.sum`
- [X] T002 Verify Go and universal build-output exclusions remain complete in `.gitignore`

---

## Phase 2: Foundational Canonical Artifacts

**Purpose**: Define the contract files required by every story.

- [X] T003 Write failing Draft 2020-12 compilation, representative-document, invalid-variant, and embedded-byte tests in `internal/schema/schema_test.go`
- [X] T004 Create the canonical fixed v0.0.0 contract in `internal/schema/cueson.schema.json`
- [X] T005 Create a conforming externally authored example in `internal/schema/testdata/representative.cueson.json`

**Checkpoint**: The canonical artifact and representative evidence exist for model and service implementation.

---

## Phase 3: User Story 1 - Validate a portable Cue JSON document (Priority: P1) MVP

**Goal**: Validate the canonical common cue and multi-asset envelope structurally and semantically.

**Independent Test**: The representative document succeeds while focused invalid structural and cross-field variants fail deterministically.

### Tests for User Story 1

- [X] T006 [US1] Write failing model invariant tests for primary assets, identities, timing, ordering, counts, timestamps, format matching, OCR references, and S003 capability truth in `internal/model/model_test.go`

### Implementation for User Story 1

- [X] T007 [US1] Define Cue JSON root, source, cue, observation, format-data, diagnostic, and statistics types in `internal/model/model.go`
- [X] T008 [US1] Implement deterministic format-neutral semantic validation in `internal/model/model.go`
- [X] T009 [US1] Implement embedded bytes, offline Draft 2020-12 compilation, raw JSON parsing, structural validation, and semantic dispatch in `internal/schema/schema.go`

**Checkpoint**: User Story 1 validates the representative contract without codec or filesystem behavior.

---

## Phase 4: User Story 2 - Retrieve the exact embedded schema (Priority: P2)

**Goal**: Add the shipped schema command, exact schema version output, and failure-safe file output.

**Independent Test**: Stdout, version, help, new-file, refusal, force-replacement, cancellation, and output-failure cases produce exact payloads and statuses.

### Tests for User Story 2

- [X] T010 [US2] Write failing schema command parsing, help, exact-output, version, file-output, overwrite, cancellation, and I/O failure tests in `internal/cli/cli_test.go`

### Implementation for User Story 2

- [X] T011 [US2] Extend command parsing, truthful root/schema help, schema dispatch, and exit classification in `internal/cli/cli.go`
- [X] T012 [US2] Implement staged schema file output with explicit overwrite and rollback behavior in `internal/cli/cli.go`

**Checkpoint**: User Story 2 exposes the embedded artifact without advertising any deferred command.

---

## Phase 5: User Story 3 - Detect contract drift before release (Priority: P3)

**Goal**: Make naming, version identity, capability truth, and embedded-byte drift fail automatically.

**Independent Test**: Conformance helpers accept the canonical artifact and reject controlled mutations for every protected identity or naming rule.

### Tests for User Story 3

- [X] T013 [US3] Add canonical `$id`, instance `$schema`, schema-version, executable-version, and embedded-byte lockstep mutation tests in `internal/schema/schema_test.go`
- [X] T014 [US3] Add recursive project-owned property and enum lowercase-snake-case conformance tests in `internal/schema/schema_test.go`
- [X] T015 [US3] Add interim `subrip` and `webvtt` capability-truth mutation tests in `internal/model/model_test.go`

### Implementation for User Story 3

- [X] T016 [US3] Expose stable internal schema identity/version access and software lockstep validation in `internal/schema/schema.go`

**Checkpoint**: User Story 3 prevents protected contract drift before release inputs are produced.

---

## Phase 6: Documentation and Verification

**Purpose**: Align public status and execute every local completion gate.

- [X] T017 [P] Update schema availability, commands, and deferred capability boundaries in `README.md`
- [X] T018 [P] Record the S003 canonical schema foundation under Unreleased in `CHANGELOG.md`
- [X] T019 Run every command and repository gate in `specs/S003-build-schema-foundation/quickstart.md`, including formatting, vet, unit tests, race tests, clean `CGO_ENABLED=0` build, CLI smoke checks, publication-tool tests, whitespace, encoding, mojibake, and Spec Kit readiness

---

## Dependencies & Execution Order

### Phase Dependencies

- Setup has no dependency.
- Foundational Canonical Artifacts depends on Setup and blocks all user stories.
- User Story 1 depends on the canonical artifacts and is the MVP.
- User Story 2 depends on the embedded schema service from User Story 1.
- User Story 3 builds on the model and schema service but remains independently verifiable through controlled drift.
- Documentation and Verification depends on all user stories.

### User Story Dependencies

- **User Story 1 (P1)**: Depends only on the canonical artifacts and delivers document validation.
- **User Story 2 (P2)**: Depends on embedded schema access from User Story 1 and delivers the public retrieval command.
- **User Story 3 (P3)**: Depends on the canonical artifact and internal accessors, then independently verifies protected drift.

### Parallel Opportunities

- The final README and changelog updates can proceed in parallel after behavior is stable.
- Schema artifact and representative fixture authoring are sequential because the fixture must target the completed shape.
- Tests and implementations that share a package are intentionally sequential for TDD and review clarity.

## Parallel Example: Documentation

```text
Task T017: Update README.md with schema command availability and deferred boundaries.
Task T018: Record the S003 schema foundation in CHANGELOG.md.
```

## Implementation Strategy

### MVP First

1. Pin the validator and define the canonical artifacts.
2. Write model invariant tests and observe their expected failures.
3. Implement the typed model and semantic validation.
4. Implement embedded structural validation and validate User Story 1 independently.

### Incremental Delivery

1. Add exact schema retrieval and safe explicit output after validation works.
2. Add drift and naming conformance gates after stable accessors exist.
3. Align documentation and execute every local verification gate.
4. Run convergence, then continue through the authorized push, PR, CI, and maximum two-round review protocol.

## Notes

- Completed tasks are marked `[X]` during implementation.
- Tests precede their corresponding implementation within each behavior phase.
- S003 adds no source-integrity execution, restoration, native format ingest, rendering, conversion, public validation command, completion, release, or production-domain work.

## Phase 7: Review Remediation

**Purpose**: Resolve verified first-round review findings without expanding S003 scope.

- [X] T020 Reject RFC 3339 timestamp fractions that exceed the `unix_ns` representation and add regression coverage in `internal/model/model.go` and `internal/model/model_test.go`
- [X] T021 Treat post-commit backup cleanup as a surfaced warning rather than a failed replacement and add status regression coverage in `internal/cli/cli.go` and `internal/cli/cli_test.go`
