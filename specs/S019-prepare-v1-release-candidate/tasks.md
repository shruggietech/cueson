# Tasks: Prepare the v1.0.0 Release Candidate

**Input**: Design documents from `specs/S019-prepare-v1-release-candidate/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/release-evidence.schema.json`, `quickstart.md`

**Tests**: Version, schema identity, stable capability, governed fixture, immutable-copy, legal-file, evidence, artifact, workflow-authority, documentation, release-note, native-smoke, and complete candidate tests are mandatory because the specification requires exact independently reviewable release proof.

**Organization**: Tasks are grouped by setup, blocking foundations, and independently testable user stories.

## Phase 1: Setup

**Purpose**: Establish the governed S019 delivery record and freeze the candidate boundary.

- [x] T001 Verify GitHub issue #37 remains the single S019 work item under milestone v1.0.0, child epic #29, and blocker of #38
- [x] T002 Set issue #37 to Slice `S019` and Stage `In progress`, and keep the default Status field empty in the `cueson Delivery` Project
- [x] T003 [P] Record the version promotion, immutable-schema, evidence, legal-file, native-smoke, working-spec correction, and publication-boundary decisions in `specs/S019-prepare-v1-release-candidate/research.md`
- [x] T004 [P] Complete and validate the specification, quality checklist, implementation plan, data model, release-evidence contract, and quickstart under `specs/S019-prepare-v1-release-candidate/`

---

## Phase 2: Foundational Candidate Controls

**Purpose**: Establish failing coverage and complete the blocking analysis gate before changing stable product or release identity.

**CRITICAL**: No release identity or immutable v1 schema is admitted until the analysis gate and contract coverage are complete.

- [x] T005 Run the blocking Spec Kit analysis gate against `specs/S019-prepare-v1-release-candidate/spec.md`, `plan.md`, and `tasks.md`, then remediate every critical, high, or constitution conflict before implementation
- [x] T006 [P] Add failing `1.0.0` executable, exact schema-identity, immutable v0 preservation, immutable v1 admission, and stable capability tests in `internal/version/version_test.go`, `internal/model/model_test.go`, `internal/schema/schema_test.go`, and `internal/schema/annotations_test.go`
- [x] T007 [P] Add failing stable producer and current-document tests in `internal/cli/cli_test.go`, `internal/cli/documentation_test.go`, `internal/convert/convert_test.go`, and current Cue JSON fixtures under `internal/schema/testdata/` and `testdata/`
- [x] T008 [P] Add failing v1 candidate, exact schema authority, legal-file byte and mode, intended-tag, distinct evidence-path, packaging-source, workflow-authority, and native-smoke policy tests in `scripts/release-verify/verify_test.go` and `scripts/release-verify/policy_test.go`
- [x] T009 [P] Add failing v1 release-note inventory, exact suffix, current-state documentation, and working-spec gate expectations in `scripts/docs-verify/verify.go` and `scripts/docs-verify/verify_test.go`

**Checkpoint**: The analysis is clean and focused tests fail only for the intended missing v1 candidate implementation.

---

## Phase 3: User Story 1 - Review the Stable v1 Release Record (Priority: P1)

**Goal**: Admit one coherent stable identity, immutable schema, detailed history, concise notes, and truthful compatibility record.

**Independent Test**: Compare every current identity and stable capability assertion, validate governed fixtures, compare canonical and immutable schema bytes, and inspect the changelog, notes, compatibility, and release boundary.

- [x] T010 [US1] Promote executable and schema constants to `1.0.0` in `internal/version/version.go`, `internal/schema/schema.go`, and their tests
- [x] T011 [US1] Promote canonical schema identity, annotations, examples, and SubRip/WebVTT support to stable in `internal/schema/cueson.schema.json`
- [x] T012 [US1] Promote capability validation and produced documents to stable in `internal/model/model.go`, `internal/cli/workflows.go`, `internal/convert/project.go`, and their tests
- [x] T013 [US1] Migrate representative and governed current-version Cue JSON fixtures to `1.0.0` stable identity in `internal/schema/testdata/` and `testdata/`, then regenerate exact `testdata/manifest.json` lengths and SHA-256 values without changing each fixture's intended acceptance layer
- [x] T014 [US1] Recompute and review canonical schema normative and full digests, preserve the pinned v0.0.0 digest, and update candidate digest tests in `internal/schema/annotations_test.go`
- [x] T015 [US1] Create `schema/releases/v1.0.0/cueson.schema.json` as an exact finalized copy of `internal/schema/cueson.schema.json` and prove regular-file, UTF-8, line-ending, identity, and byte-equality requirements
- [x] T016 [US1] Move complete accumulated history into dated `1.0.0` sections, retain a fresh empty `[Unreleased]`, and add comparison links in `CHANGELOG.md`
- [x] T017 [P] [US1] Add concise publication-ready v1 highlights and truthful boundaries ending with the exact tagged changelog link in `docs/releases/v1.0.0.md`
- [x] T018 [US1] Reconcile v1 candidate identity, stable format support, compatibility policy, release boundaries, and current Project evidence in `README.md`, `CONTRIBUTING.md`, `SECURITY.md`, `docs/architecture.md`, `docs/cli.md`, `docs/compatibility.md`, `docs/conversion.md`, `docs/formats/srt.md`, `docs/formats/webvtt.md`, `docs/schema.md`, `docs/project-management.md`, `docs/release-process.md`, and `docs/release-verification.md`
- [x] T019 [US1] Correct the contradictory pre-release public-schema and verification gates in `docs/Cueson-Project-Specification-v0.0.0.md` while keeping production schema hosting separately authorized, and record the rationale in `CHANGELOG.md`
- [x] T020 [US1] Add `docs/releases/v1.0.0.md` and v1 candidate assertions to the maintained documentation contract in `scripts/docs-verify/verify.go` and `scripts/docs-verify/verify_test.go`
- [x] T021 [US1] Run focused version, schema, model, producer, fixture, conformance, documentation, formatter, JSON, encoding, mojibake, and byte-identity checks from `specs/S019-prepare-v1-release-candidate/quickstart.md`

**Checkpoint**: User Story 1 is independently complete and every permanent v1 candidate record is reviewable without claiming publication.

---

## Phase 4: User Story 2 - Prove the Exact v1 Candidate (Priority: P1)

**Goal**: Bind one clean exact source revision to the complete non-publishing artifact, schema, legal-file, provenance, checksum, software-bill-of-material, and native-execution record.

**Independent Test**: Build the six-target candidate from a clean exact commit, verify every target structurally, execute the accepted bundle on hosted Linux, Windows, and macOS runners, and inspect evidence for the intended tag, exact revision, immutable schema and legal digests, inventories, and non-publication.

- [x] T022 [US2] Extend exact schema authority, repository legal-file loading, archive byte and mode comparison, intended-tag evidence, and configurable evidence output in `scripts/release-verify/main.go`, `scripts/release-verify/verify.go`, and platform launch adapters
- [x] T023 [US2] Update and complete v1 candidate verification coverage in `scripts/release-verify/verify_test.go`, `scripts/release-verify/policy_test.go`, and `specs/S019-prepare-v1-release-candidate/contracts/release-evidence.schema.json`
- [x] T024 [US2] Fix candidate version `1.0.0` and immutable schema packaging while retaining the six-target, four-member, pure-Go, snapshot-only, publication-disabled contract in `.goreleaser.yaml`
- [x] T025 [US2] Convert `.github/workflows/release-proof.yml` to v1 candidate mode, upload one accepted bundle, and fan that exact bundle to read-only native packaged-binary smoke jobs on Ubuntu, Windows, and macOS without tag, credential, secret, signing, release, identity-token, or production authority
- [x] T026 [US2] Run release-verifier tests, race, vet, policy checks, actionlint, GoReleaser validation, and exact repository schema and legal-file identity checks against the final candidate configuration
- [ ] T027 [US2] Commit the complete implementation locally, build the six-target candidate with pinned GoReleaser and Syft from that clean exact commit, and run `scripts/release-verify` without development mode including host execution
- [ ] T028 [US2] Inspect `dist/release-evidence.json` and confirm exact revision, version, intended tag, schema and legal digests, six archives, six software bills of materials, six checksums, complete target digests, host execution, and `published: false`
- [ ] T029 [US2] Reproduce the release proof from a second clean checkout of the same exact commit, compare stable archive, checksum, schema, legal, and target semantics, and retain truthful non-byte-reproducibility boundaries for Syft documents

**Checkpoint**: User Story 2 is independently complete with exact clean-commit evidence and no protected publication action.

---

## Phase 5: User Story 3 - Hand Off a Governed v1 Publication Decision (Priority: P2)

**Goal**: Publish and review the official S019 pull request while preserving every later operator boundary.

**Independent Test**: Run every local and hosted gate, inspect issue and Project state, and confirm the pull request is green, fully reviewed, configured to close issue #37, and otherwise non-publishing.

- [ ] T030 [US3] Run the complete local quickstart including formatter, documentation, product, policy, release, workflow, race, vet, vulnerability, brand-integrity, whitespace, encoding, mojibake, fixture, and clean-checkout candidate proof from `specs/S019-prepare-v1-release-candidate/quickstart.md`
- [ ] T031 [US3] Run Spec Kit convergence against `specs/S019-prepare-v1-release-candidate/` and append and execute any required remediation tasks until convergence is clean
- [ ] T032 [US3] Push `codex/S019-prepare-v1-release-candidate`, publish the formatted official pull request with `Closes #37`, and verify the rendered GitHub body
- [ ] T033 [US3] Move issue #37 to Stage `PR review`, retain Slice `S019`, and verify default Status remains empty in `cueson Delivery`
- [ ] T034 [US3] Monitor every current-head CI, CodeQL, release-proof, native packaged-smoke, pull-request-policy, Codex, security, and review result; address every round-one finding and resolve threads only after remediation
- [ ] T035 [US3] Request exactly one `@codex review` second round only if round one reports findings or a repository-documented condition requires it, then address that round without requesting a third
- [ ] T036 [US3] Re-run convergence and the complete local verification suite after all remediation, push the final descendant head, and wait for every hosted gate to return green
- [ ] T037 [US3] Confirm all reviews are satisfied, the pull request is mergeable, issue #37 remains open pending merge, issue #38 remains blocked, and no merge, tag, release, milestone closure, schema publication, or production mutation occurred

**Checkpoint**: S019 is ready for the operator's final review and merge ritual.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Starts from the approved S019 direction and explicit push and pull-request authorization.
- **Foundational controls (Phase 2)**: Depends on Setup and blocks stable identity changes.
- **User Story 1 (Phase 3)**: Depends on failing contract coverage and the clean analysis gate.
- **User Story 2 (Phase 4)**: Depends on the finalized stable identity and immutable schema.
- **User Story 3 (Phase 5)**: Depends on all local implementation, clean-checkout candidate verification, and Spec Kit convergence.

### User Story Dependencies

- **User Story 1**: Independently establishes the permanent stable release record.
- **User Story 2**: Uses User Story 1's immutable schema and release identity but is independently testable as artifact evidence.
- **User Story 3**: Integrates the release record and candidate proof into the governed GitHub delivery flow.

### Parallel Opportunities

- T003 and T004 inspect and author separate planning surfaces.
- T006, T007, T008, and T009 modify separate test and contract areas.
- T017 can proceed while current-state documentation is reconciled.
- Native smoke runners execute the same accepted bundle independently on three operating systems.

## Implementation Strategy

1. Establish and verify the issue and Project lifecycle record.
2. Complete specification, planning, task generation, and the blocking analysis gate.
3. Write coordinated stable-identity, immutable-schema, fixture, evidence, legal-file, workflow, documentation, and notes coverage before or alongside implementation.
4. Admit and prove the stable release record as the minimum viable outcome.
5. Switch the existing non-publishing pipeline to candidate mode and add native packaged-binary evidence without adding publication authority.
6. Commit locally, prove the exact clean candidate twice, complete convergence and the repository verification suite, then publish the official pull request.
7. Remediate every review finding with at most one second Codex request and stop only at the green fully reviewed boundary.
