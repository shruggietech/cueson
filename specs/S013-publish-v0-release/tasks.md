# Tasks: Publish and Verify v0.0.0

**Input**: Design documents from `specs/S013-publish-v0-release/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/release-publication-contract.json`

**Tests**: Exact tag, release body, public asset inventory, byte identity, checksum, compatible-host execution, protected-state, documentation, local verification, and hosted review checks are mandatory because the specification requires independently verified publication.

**Organization**: Tasks are grouped by setup, blocking foundations, and independently testable user stories. Remote mutations remain sequential so every protected state transition is read back before the next begins.

## Phase 1: Setup

**Purpose**: Establish the governed S013 delivery record and freeze the operator-authorized inputs.

- [x] T001 Create GitHub issue #27 with the repository-required sections, v0.0.0 milestone, governed labels, and independently testable acceptance criteria
- [x] T002 Add issue #27 exactly once to `cueson Delivery`, set Slice `S013` and Stage `Ready`, and clear the default Status field
- [x] T003 [P] Record tag, release, asset, verification, documentation, and protected-boundary decisions in `specs/S013-publish-v0-release/research.md`
- [x] T004 [P] Record the exact candidate, release entities, public asset inventory, and digest contract in `specs/S013-publish-v0-release/data-model.md` and `specs/S013-publish-v0-release/contracts/release-publication-contract.json`
- [x] T005 [P] Verify the accepted run, artifact, evidence, source revision, schema digest, current GitHub issue inventory, Project fields, milestone state, and absence of an existing v0.0.0 tag or release

---

## Phase 2: Foundational Publication Controls

**Purpose**: Complete the blocking analysis gate and validate all immutable inputs before remote mutation.

**CRITICAL**: No tag or GitHub Release may be created until every task in this phase passes.

- [x] T006 Run the blocking Spec Kit analysis gate against `specs/S013-publish-v0-release/spec.md`, `plan.md`, and `tasks.md` and remediate every critical, high, or constitution conflict before implementation
- [x] T007 Validate `specs/S013-publish-v0-release/contracts/release-publication-contract.json` syntax, uniqueness, kind counts, exact filenames, exact SHA-256 digests, source revision, schema digest, and accepted evidence alignment
- [x] T008 Format-check `docs/releases/v0.0.0.md` through `scripts/github-format`, prove the formatter changes no bytes, verify its exact final changelog link, and preserve the file unchanged
- [x] T009 Reconfirm accepted default-branch CI, CodeQL, and release-proof success for `b294a6952c8bd041d852c502f5d7206c0b58edd6` and confirm the retained artifact is still available
- [x] T010 Reconfirm no local or remote `v0.0.0` tag and no GitHub Release exists immediately before publication; stop on any conflicting or partial state

**Checkpoint**: Analysis is clean, all authorized input bytes and GitHub evidence match, and publication preflight passes.

---

## Phase 3: User Story 1 - Obtain the Authorized Release (Priority: P1)

**Goal**: Publish one immutable tag and one complete user-facing v0.0.0 GitHub Release.

**Independent Test**: Read the tag and release directly from GitHub and compare the tag target, release state, body, and thirteen asset names with the authorized contract.

- [x] T011 [US1] Create unsigned annotated tag `v0.0.0` at `b294a6952c8bd041d852c502f5d7206c0b58edd6` and prove its local object type and peeled target before push
- [x] T012 [US1] Push only `refs/tags/v0.0.0`, then read back the remote annotated tag and prove `refs/tags/v0.0.0^{}` equals the authorized revision
- [x] T013 [US1] Publish `Cueson v0.0.0` as a public non-draft non-prerelease GitHub Release using unchanged `docs/releases/v0.0.0.md` and the thirteen explicit accepted asset paths
- [x] T014 [US1] Immediately read the GitHub Release back and verify its tag, target, name, state, body structure, final changelog link, and exact thirteen-name asset inventory with zero internal evidence or metadata files

**Checkpoint**: User Story 1 is publicly visible and exactly matches the authorized release identity and inventory.

---

## Phase 4: User Story 2 - Independently Verify Publication (Priority: P1)

**Goal**: Prove that every public byte is the accepted candidate and that separately protected state remains unchanged.

**Independent Test**: Download the public release to a new directory, compare thirteen digests, apply the archive checksums, execute the compatible binary, compare schema bytes, and inspect protected boundaries.

- [x] T015 [US2] Download all public v0.0.0 assets into a newly created clean temporary directory and prove the downloaded filename set equals the thirteen-file contract
- [x] T016 [US2] Compare all thirteen downloaded SHA-256 values with `release-publication-contract.json` and apply the six-entry checksum manifest as an exact archive bijection
- [x] T017 [US2] Extract the public Windows amd64 archive safely, run `version`, `schema --version`, and schema emission non-interactively, and prove exact `0.0.0` output plus byte identity with the tagged versioned schema
- [x] T018 [US2] Reconcile exact public byte identity with accepted `release-evidence.json` to prove all six archives and SBOMs retain their verified members, target, version, source revision, and schema semantics
- [x] T019 [US2] Verify milestone v0.0.0 remains open and no production `cueson.io`, public schema-hosting, signature, attestation, or other excluded action occurred

**Checkpoint**: User Story 2 is independently complete and the public release is verified without crossing another protected boundary.

---

## Phase 5: User Story 3 - Record the Released State (Priority: P2)

**Goal**: Make the repository and official S013 pull request truthfully describe and govern the completed publication.

**Independent Test**: Inspect all release-facing prose, run offline documentation and repository checks, and verify the issue, Project, pull request, hosted checks, and reviews.

- [x] T020 [P] [US3] Update current release identity, download, status, capability, and verification links in `README.md` while retaining the envelope-only limitation
- [x] T021 [P] [US3] Update released-state, immutable-tag, verified-asset, protected-boundary, and milestone-lifecycle prose in `docs/release-process.md`, `docs/release-verification.md`, and `docs/schema.md`
- [x] T022 [P] [US3] Update release topology, exact publication ownership, and current delivery status in `docs/architecture.md`, `docs/cli.md`, and `docs/project-management.md`
- [x] T023 [US3] Record the post-tag documentation and release-verification update under `[Unreleased]` in `CHANGELOG.md` without modifying the tagged v0.0.0 notes or versioned schema
- [x] T024 [US3] Search maintained documentation for current unreleased-project claims and native-codec overclaims, then remediate every false current-state statement while preserving explicit historical context
- [x] T025 [US3] Run all quickstart verification including formatter, documentation, product, policy, release, workflow, race, vet, vulnerability, brand-integrity, whitespace, encoding, mojibake, and tagged-file-immutability checks
- [x] T026 [US3] Run Spec Kit convergence against `specs/S013-publish-v0-release/` and append and execute any required remediation tasks until convergence is clean
- [ ] T027 [US3] Push `S013-publish-v0-release`, publish the formatted official pull request with `Closes #27`, and verify the rendered GitHub body
- [ ] T028 [US3] Move issue #27 to Stage `PR review`, retain Slice `S013`, and verify default Status remains empty in `cueson Delivery`
- [ ] T029 [US3] Monitor every current-head CI, CodeQL, release-proof, pull-request policy, Codex, security, and review result; address every round-one finding and resolve threads only after remediation
- [ ] T030 [US3] Request exactly one `@codex review` second round only if round one reports findings or a repository-documented condition requires it, then address that round without requesting a third
- [ ] T031 [US3] Re-run convergence and the complete local verification suite after all remediation, push the final descendant head, and wait for every hosted gate to return green
- [ ] T032 [US3] Confirm all reviews are satisfied, the pull request is mergeable, issue #27 remains open pending merge, and no PR merge, milestone closure, schema publication, production mutation, signature, attestation, or native codec claim occurred

**Checkpoint**: S013 is ready for the operator's final review and merge ritual.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Starts immediately from the operator's exact publication and GitHub delivery authorization.
- **Foundational controls (Phase 2)**: Depends on Setup and blocks every remote mutation.
- **User Story 1 (Phase 3)**: Depends on complete preflight; its tag and release mutations are strictly sequential.
- **User Story 2 (Phase 4)**: Depends on the public release and blocks released-state documentation claims.
- **User Story 3 (Phase 5)**: Depends on verified publication and concludes with the official reviewed pull request.

### User Story Dependencies

- **User Story 1**: Independently produces the authorized public tag, notes, and downloads.
- **User Story 2**: Uses User Story 1's public bytes but independently proves identity and protected-state containment.
- **User Story 3**: Uses verified publication evidence to reconcile repository and GitHub delivery records.

### Parallel Opportunities

- T003, T004, and T005 inspect or record separate preparation surfaces.
- T020, T021, and T022 modify separate documentation files after publication verification.
- Remote tag creation, tag push, release creation, release read-back, and public download verification are intentionally not parallel.

## Implementation Strategy

1. Establish issue and Project ownership and freeze the exact authorized candidate and asset contract.
2. Complete specification, planning, task generation, and the blocking analysis gate.
3. Revalidate every accepted input and prove absence of conflicting public state.
4. Create and verify the immutable tag, then publish and read back the complete release.
5. Download all public assets, prove byte identity and compatible execution, and audit excluded boundaries.
6. Reconcile repository documentation without altering tagged history.
7. Complete convergence and the full local suite, then publish the official pull request.
8. Remediate every review finding with at most one second Codex request and stop only at the green fully reviewed boundary.
