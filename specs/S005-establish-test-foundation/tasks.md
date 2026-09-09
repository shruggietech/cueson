# Tasks: Establish Fixture and Conformance-Test Infrastructure

**Input**: Design documents from `/specs/S005-establish-test-foundation/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`

**Tests**: Test-first work is required by the specification and constitution. Each implementation group begins with a failing focused test or an equivalent deterministic verification.

**Organization**: Tasks are grouped by user story so each story has an independently testable outcome while sharing the smallest domain-neutral foundation.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it owns different files and has no unresolved dependency.
- **[Story]**: Maps the task to the user story whose outcome it proves.
- Every task names the file or directory it changes.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Establish repository rules and the documented root for byte-sensitive test assets.

- [X] T001 Update `.gitattributes` and `.editorconfig` so repository-authored test metadata stays UTF-8/LF while authoritative source, expected-byte, malformed-input, and future fuzz-input subtrees bypass text and whitespace normalization.
- [X] T002 Create `testdata/README.md` with the governed directory layout, central-manifest authority, provenance and redistribution obligations, contributor workflow, and the explicit boundary excluding native codec and CI coverage.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Build the portable manifest, inventory, integrity, and path-safety layer required by every story.

**Critical gate**: No user-story implementation begins until this phase passes focused tests.

- [X] T003 Add failing table-driven tests for strict manifest shape, provenance, redistribution, portable path identity, symlink and non-regular rejection, inventory completeness, exact size/hash checks, and declared byte characteristics in `internal/testutil/fixtures_test.go` and `internal/testutil/paths_test.go`.
- [X] T004 Implement portable relative-path validation, Unicode canonical-caseless collision keys, safe component checks, and path-free logical labels in `internal/testutil/paths.go`.
- [X] T005 Implement strict versioned manifest loading plus deterministic provenance, redistribution, inventory, regular-file, byte-characteristic, size, and SHA-256 verification in `internal/testutil/fixtures.go`.
- [X] T006 Run `go test ./internal/testutil` and confirm all foundational tests pass without mutating the temporary fixture trees.

**Checkpoint**: The shared fixture verifier is ready and rejects invalid corpus state before conformance assertions.

---

## Phase 3: User Story 1 - Add Trustworthy Fixtures (Priority: P1) MVP

**Goal**: Commit a small, byte-stable, fully attributed seed corpus whose inventory and integrity are machine-verifiable.

**Independent Test**: Verify the committed root manifest, then perturb copies of its metadata and payload bytes to prove all governed failure modes reject deterministically.

### Implementation and verification

- [X] T007 [US1] Add the synthetic LF source-envelope seed and normalized expected model and diagnostic artifacts beneath `testdata/fixtures/source-envelope/basic-lf/`.
- [X] T008 [US1] Create `testdata/manifest.json` with ordered version-1 records, complete origin and redistribution decisions, byte contracts, exact lengths, lowercase SHA-256 digests, and the accepted expectation for the seed corpus.
- [X] T009 [US1] Add committed-corpus verification and exact source-envelope byte identity coverage to `internal/conformance/conformance_test.go`.
- [X] T010 [US1] Verify `git check-attr`, manifest inventory, byte lengths, and SHA-256 values for every committed fixture artifact, prove the same hashes from a clean Git-materialized copy with checkout normalization applied, and record the portable commands in `testdata/README.md`.

**Checkpoint**: Every committed fixture payload is declared once, byte-stable, redistributable, and independently verifiable.

---

## Phase 4: User Story 2 - Reuse One Conformance Vocabulary (Priority: P2)

**Goal**: Supply path-free comparison helpers and prove one accepted case across schema, model, diagnostics, source bytes, restoration, integrity, and timestamps.

**Independent Test**: Exercise each comparison surface once with equal values and once with a deliberate mismatch, then compare equivalent expectations from two different roots.

### Tests for User Story 2

- [X] T011 [US2] Add failing tests for semantic JSON/model equality, ordered diagnostic equality, exact bytes with stable offset evidence, byte length, lowercase SHA-256, and restored/unsupported/unavailable/failed timestamp outcomes in `internal/testutil/golden_test.go`.
- [X] T012 [US2] Add failing tests for explicit forbidden sentinel variants, including native, slash, backslash, and JSON-escaped forms, and equivalent portable output from distinct roots in `internal/testutil/paths_test.go`.

### Implementation for User Story 2

- [X] T013 [US2] Implement generic JSON, diagnostics, bytes, integrity, and timestamp comparisons with fixture-ID and logical-surface diagnostics in `internal/testutil/golden.go`.
- [X] T014 [US2] Implement explicit forbidden-token expansion and path-free portable-output checks in `internal/testutil/paths.go`.
- [X] T015 [US2] Extend `internal/conformance/conformance_test.go` to prove accepted Cue JSON decoding, typed model semantics, ordered diagnostics, exact source restoration and hash identity, explicit timestamp outcomes, distinct-root equivalence, and zero forbidden local identifiers.
- [X] T016 [US2] Run `go test ./internal/testutil ./internal/conformance` and confirm every comparison surface has a passing and deliberately mismatching assertion.

**Checkpoint**: Future codecs can share one deterministic, path-free conformance vocabulary without importing domain packages into `internal/testutil`.

---

## Phase 5: User Story 3 - Retain Malformed and Fuzz Regressions (Priority: P3)

**Goal**: Preserve deterministic boundary regressions and establish bounded, side-effect-free fuzz entry points for implemented Cue JSON and source-envelope behavior.

**Independent Test**: Run every malformed case twice and run each fuzz target for a fixed mutation count with no panic, destination write, network access, external process, or native codec claim.

### Implementation and verification

- [X] T017 [US3] Add minimal parse, structure, semantics, and integrity rejection inputs beneath `testdata/malformed/` and register their exact bytes, provenance, rejection stage, and stable diagnostic fragment in `testdata/manifest.json`.
- [X] T018 [US3] Extend `internal/conformance/conformance_test.go` to execute every registered malformed case twice, verify the declared stage and fragment, preserve input bytes, and prove integrity failure creates no output.
- [X] T019 [P] [US3] Add `FuzzDecodeCueJSON` with representative, empty, malformed, structurally invalid, and semantically invalid seeds plus a 64 KiB harness guard and accepted-document revalidation in `internal/schema/fuzz_test.go`.
- [X] T020 [P] [US3] Add `FuzzInspectEncoded` and `FuzzValidateSafeBasename` with valid and corrupt envelope seeds, deterministic repeat checks, 64 KiB and 1 KiB harness guards, and no restoration call in `internal/source/fuzz_test.go`.
- [X] T021 [US3] Run the deterministic malformed suite twice and execute each fuzz target in the foreground with `-fuzztime=1000x -parallel=1`, confirming bounded completion and no persisted fuzz artifacts.

**Checkpoint**: Implemented schema and envelope boundaries have permanent deterministic regressions and bounded fuzz smoke coverage, with native subtitle grammars still out of scope.

---

## Phase 6: Polish and Cross-Cutting Verification

**Purpose**: Reconcile documentation, delivery evidence, and the full repository verification surface.

- [X] T022 Update `docs/architecture.md` and finish `testdata/README.md` with ownership boundaries, the embedded representative exception, contributor promotion procedure, path-leak rules, verification commands, and explicit deferral of CI and native codec coverage.
- [X] T023 Update the `Unreleased` section of `CHANGELOG.md` with the S005 fixture, conformance, malformed-regression, and fuzz-boundary additions without claiming codec support.
- [X] T024 Run formatting, `go test ./...`, `go test -race ./...`, `go vet ./...`, repository text-integrity checks, fixture verification, malformed repeat tests, bounded fuzz smoke runs, `git diff --check`, and Spec Kit convergence; resolve every in-scope failure and reconcile issue #7 plus the S005 Project fields before publication.

---

## Dependencies and Execution Order

### Phase dependencies

- **Setup (Phase 1)** starts immediately.
- **Foundational (Phase 2)** depends on Setup and blocks all user stories.
- **User Story 1 (Phase 3)** depends on the foundational verifier.
- **User Story 2 (Phase 4)** depends on the accepted seed from User Story 1.
- **User Story 3 (Phase 5)** depends on the manifest and generic helper contracts but its two fuzz files can be implemented in parallel after those boundaries are stable.
- **Polish (Phase 6)** depends on all selected stories.

### Within each story

- Tests or deterministic failing checks precede their implementation.
- Manifest verification precedes conformance assertions.
- Generic helpers precede cross-package projections.
- Full verification and convergence follow all code and documentation changes.

### Parallel opportunities

- T019 and T020 own separate package-local fuzz files and can run in parallel after T018.
- Documentation review can occur alongside focused test execution once implementation semantics stop changing.

## Implementation Strategy

1. Establish normalization and corpus rules.
2. Build and verify the shared manifest foundation.
3. Commit and validate the smallest accepted seed.
4. Add reusable comparison semantics and the accepted cross-package proof.
5. Add deterministic malformed cases and bounded fuzz targets.
6. Reconcile documentation and run the complete foreground verification sequence.

## Notes

- S005 closes issue #7 only. Issue #8 owns hosted CI, and later codec issues own native SRT and WebVTT grammars.
- The existing `internal/schema/testdata/representative.cueson.json` remains a canonical embedded schema example outside the governed root corpus.
- The schema and source packages currently apply different 255-unit basename interpretations. S005 uses the stricter portable byte-safe fixture rule and does not redefine either product contract.
- No task authorizes a push, pull-request merge, release, tag, or production-domain mutation. The operator has separately authorized S005 push and pull-request publication; final merge remains human-controlled.
