# Tasks: Prepare the v0.0.0 Release

**Input**: Design documents from `specs/S012-prepare-v0-release/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/release-evidence.schema.json`

**Tests**: Schema identity, evidence, artifact, workflow-authority, documentation, release-note, and complete candidate tests are mandatory because the specification requires exact and independently reviewable release proof.

**Organization**: Tasks are grouped by setup, blocking foundations, and independently testable user stories.

## Phase 1: Setup

**Purpose**: Establish the governed S012 delivery record and confirm the release boundary.

- [x] T001 Create GitHub issue #25 with the repository-required sections, v0.0.0 milestone, governed labels, and independently testable acceptance criteria
- [x] T002 Add issue #25 exactly once to `cueson Delivery`, set Slice `S012` and Stage `In progress`, and clear the default Status field
- [x] T003 [P] Record release-candidate, schema-admission, evidence, workflow, and publication-boundary decisions in `specs/S012-prepare-v0-release/research.md`
- [x] T004 [P] Verify the current release configuration, release verifier, release documentation, issue inventory, Project fields, milestone state, and protected authority boundaries in `specs/S012-prepare-v0-release/plan.md`

---

## Phase 2: Foundational Release Controls

**Purpose**: Establish failing contract coverage and complete the blocking analysis gate before release records are admitted.

**CRITICAL**: No release schema, notes, or candidate proof is complete until these controls pass.

- [x] T005 Run the blocking Spec Kit analysis gate against `specs/S012-prepare-v0-release/spec.md`, `plan.md`, and `tasks.md` and remediate every critical, high, or constitution conflict before implementation
- [x] T006 [P] Add failing release-schema identity, missing-schema, byte-drift, evidence-digest, and deterministic-evidence tests in `scripts/release-verify/verify_test.go`
- [x] T007 [P] Add failing packaging-source, `main` push trigger, exact release-note suffix, and non-publishing authority tests in `scripts/release-verify/policy_test.go`
- [x] T008 [P] Add failing required-document and release-note link expectations for `docs/releases/v0.0.0.md` in `scripts/docs-verify/verify_test.go`

**Checkpoint**: Analysis is clean and the release contract tests fail for the expected missing implementation.

---

## Phase 3: User Story 1 - Review a Complete Release Record (Priority: P1)

**Goal**: Admit the exact release schema, dated history, and publication-ready highlights as one coherent record.

**Independent Test**: Compare canonical and versioned schema bytes, validate version identity, inspect the changelog transition, and verify the release notes and their exact final link.

- [x] T009 [US1] Copy `internal/schema/cueson.schema.json` byte for byte to `schema/releases/v0.0.0/cueson.schema.json`
- [x] T010 [US1] Change the schema source in `.goreleaser.yaml` to `schema/releases/v0.0.0/cueson.schema.json` while preserving the four-member archive contract
- [x] T011 [US1] Add repository-schema loading, byte-identity validation, version validation, and release-schema SHA-256 evidence to `scripts/release-verify/verify.go`
- [x] T012 [US1] Move accumulated history into `## [0.0.0] - 2026-09-10`, retain a fresh `[Unreleased]` section, and add comparison links in `CHANGELOG.md`
- [x] T013 [P] [US1] Add concise, limitation-aware, publication-ready highlights with the exact final changelog link in `docs/releases/v0.0.0.md`
- [x] T014 [US1] Add `docs/releases/v0.0.0.md` to the maintained documentation contract and update expected inventory in `scripts/docs-verify/verify.go` and `scripts/docs-verify/verify_test.go`
- [x] T015 [US1] Run focused release-verifier, documentation-verifier, formatter, byte-comparison, JSON, and whitespace checks for the completed release record

**Checkpoint**: User Story 1 is independently complete and every permanent v0.0.0 record is reviewable.

---

## Phase 4: User Story 2 - Prove the Exact Release Candidate (Priority: P1)

**Goal**: Bind one clean source revision to the complete non-publishing artifact and schema evidence.

**Independent Test**: Build the six-target snapshot from a clean commit and inspect deterministic evidence for exact revision, schema digest, target inventory, checksums, SBOMs, host execution, and non-publication.

- [x] T016 [US2] Extend `.github/workflows/release-proof.yml` to run the same proof on pushes to `main` without adding tag triggers, write authority, secrets, credentials, signing, deployment, or publication
- [x] T017 [P] [US2] Reconcile pull-request and post-squash candidate binding, evidence fields, exact operator decision package, and protected actions in `docs/release-process.md`
- [x] T018 [P] [US2] Document versioned-schema verification, enriched evidence, and `main` push proof in `docs/release-verification.md`
- [x] T019 [P] [US2] Reconcile current candidate and release-preparation status links in `README.md` without claiming that v0.0.0 is public
- [x] T020 [US2] Run release policy tests, actionlint, GoReleaser configuration validation, and repository schema identity checks against the final configuration
- [x] T021 [US2] Commit the release-preparation implementation locally, build the complete snapshot with GoReleaser v2.18.1 and Syft v1.51.1 from that clean commit, and verify it with `scripts/release-verify` including compatible host execution
- [x] T022 [US2] Inspect `dist/release-evidence.json` and confirm exact revision, version, schema digest, six archives, six SBOMs, six checksums, complete target digests, host execution, and `published: false`

**Checkpoint**: User Story 2 is independently complete with exact clean-commit evidence and no publication action.

---

## Phase 5: User Story 3 - Hand Off a Governed Publication Decision (Priority: P2)

**Goal**: Publish and review the official S012 pull request while preserving every later operator boundary.

**Independent Test**: Run every local and hosted gate, inspect issue and Project state, and confirm the pull request is green, fully reviewed, configured to close issue #25, and otherwise non-publishing.

- [x] T023 [US3] Run all quickstart verification including formatter, documentation, product, policy, release, workflow, race, vet, vulnerability, brand-integrity, whitespace, encoding, and mojibake checks from `specs/S012-prepare-v0-release/quickstart.md`
- [x] T024 [US3] Run Spec Kit convergence against `specs/S012-prepare-v0-release/` and append and execute any required remediation tasks until convergence is clean
- [ ] T025 [US3] Push `S012-prepare-v0-release`, publish the formatted official pull request with `Closes #25`, and verify the rendered GitHub body
- [ ] T026 [US3] Move issue #25 to Stage `PR review`, retain Slice `S012`, and verify default Status remains empty in `cueson Delivery`
- [ ] T027 [US3] Monitor every current-head CI, CodeQL, release-proof, pull-request policy, Codex, security, and review result; address every round-one finding and resolve threads only after remediation
- [ ] T028 [US3] Request exactly one `@codex review` second round only if round one reports findings or a repository-documented condition requires it, then address that round without requesting a third
- [ ] T029 [US3] Re-run convergence and the complete local verification suite after all remediation, push the final descendant head, and wait for every hosted gate to return green
- [ ] T030 [US3] Confirm all reviews are satisfied, the pull request is mergeable, issue #25 remains open pending merge, and no merge, tag, release, milestone closure, schema publication, or production mutation occurred

**Checkpoint**: S012 is ready for the operator's final review and merge ritual.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Starts immediately from the approved S012 direction and push/PR authorization.
- **Foundational controls (Phase 2)**: Depends on Setup and blocks release-record implementation.
- **User Story 1 (Phase 3)**: Depends on failing contract coverage and the clean analysis gate.
- **User Story 2 (Phase 4)**: Depends on the admitted versioned schema and verifier evidence fields.
- **User Story 3 (Phase 5)**: Depends on all local implementation and exact candidate verification.

### User Story Dependencies

- **User Story 1**: Independently establishes the permanent human-readable and schema release record.
- **User Story 2**: Uses User Story 1's versioned schema and notes but is independently testable as artifact evidence.
- **User Story 3**: Integrates the release record and candidate proof into the governed GitHub delivery flow.

### Parallel Opportunities

- T003 and T004 inspect separate authority surfaces.
- T006, T007, and T008 modify separate test files.
- T013 can proceed while T011 implements release evidence.
- T017, T018, and T019 modify separate maintained documents.

## Implementation Strategy

1. Establish the issue and Project lifecycle record.
2. Complete specification, planning, task generation, and the blocking analysis gate.
3. Write schema, workflow, notes, documentation, and evidence tests before or alongside implementation.
4. Admit and prove the permanent release record as the minimum viable outcome.
5. Add exact post-squash candidate binding and run the full non-publishing artifact proof from a clean local commit.
6. Complete convergence and the repository verification suite, then publish the official pull request.
7. Remediate review findings with at most one second Codex request and stop only at the green fully reviewed boundary.
