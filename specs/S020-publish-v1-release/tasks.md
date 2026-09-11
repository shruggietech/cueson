# Tasks: Publish and Verify v1.0.0

**Input**: Design documents from `specs/S020-publish-v1-release/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/release-publication-contract.json`

**Tests**: Exact candidate, tag, release body, public asset inventory, byte identity, checksum, archive, SBOM, compatible-host execution, protected-state, documentation, local verification, and hosted review checks are mandatory because the specification requires independently verified publication.

**Organization**: Tasks are grouped by setup, blocking controls, and independently testable user stories. Protected remote mutations remain sequential and require the exact operator authority stated at their checkpoints.

## Phase 1: Setup

**Purpose**: Establish the governed S020 delivery record and freeze exact accepted inputs.

- [x] T001 Confirm GitHub issue #38 has the repository-required sections, v1.0.0 milestone, governed labels, parent epic relationship, and independently testable acceptance criteria
- [x] T002 Add issue #38 exactly once to `cueson Delivery`, set Slice `S020` and Stage `Specced`, and clear the default Status field
- [x] T003 [P] Record publication, verification, documentation, authority, and protected-boundary decisions in `specs/S020-publish-v1-release/research.md`
- [x] T004 [P] Record the exact candidate, release entities, public inventory, sizes, and digest contract in `specs/S020-publish-v1-release/data-model.md` and `specs/S020-publish-v1-release/contracts/release-publication-contract.json`
- [x] T005 [P] Download artifact `10273380044` from accepted run `34621429626` to a clean temporary directory and record its evidence, source revision, schema digest, public-file sizes, and public-file SHA-256 values

---

## Phase 2: Foundational Publication Controls

**Purpose**: Complete the blocking analysis gate and validate immutable inputs before requesting publication authority.

**CRITICAL**: No tag or GitHub Release may be created until every task in this phase passes and the operator grants exact publication authority.

- [x] T006 Run the blocking Spec Kit analysis gate against `specs/S020-publish-v1-release/spec.md`, `plan.md`, and `tasks.md` and remediate every critical, high, or constitution conflict before implementation
- [x] T007 Validate `specs/S020-publish-v1-release/contracts/release-publication-contract.json` syntax, uniqueness, kind counts, exact filenames, exact sizes, exact SHA-256 digests, source revision, schema digest, and accepted evidence alignment
- [x] T008 Format-check `docs/releases/v1.0.0.md` through `scripts/github-format`, prove the formatter changes no bytes, verify its exact final changelog link, and preserve the file unchanged
- [x] T009 Reconfirm accepted default-branch CI, CodeQL, release proof, and three native packaged smoke results for `2cad4c816340404289b4d1d87179a4071713bb46` and confirm artifact `10273380044` remains available
- [x] T010 Reconfirm no local or remote `v1.0.0` tag and no GitHub Release exists immediately before the authority checkpoint; stop on any conflicting or partial state
- [ ] T011 Present the exact frozen transaction to the operator and obtain explicit authorization to create and push annotated tag `v1.0.0` and publish `Cueson v1.0.0` from artifact `10273380044`

**Checkpoint**: Analysis and preflight are clean. Halt here unless exact tag and release publication authority has been granted.

---

## Phase 3: User Story 1 - Obtain the Authorized Stable Release (Priority: P1)

**Goal**: Publish one immutable tag and one complete user-facing v1.0.0 GitHub Release.

**Independent Test**: Read the tag and release directly from GitHub and compare tag target, release state, body, and thirteen asset names with the frozen contract.

- [ ] T012 [US1] After exact authority, move issue #38 to Stage `In progress`, clear default Status, create unsigned annotated tag `v1.0.0` at `2cad4c816340404289b4d1d87179a4071713bb46`, and prove its local object type and peeled target before push
- [ ] T013 [US1] Push only `refs/tags/v1.0.0`, then read back the remote annotated tag and prove its peeled target equals the authorized revision
- [ ] T014 [US1] Publish `Cueson v1.0.0` as a public non-draft non-prerelease GitHub Release using unchanged `docs/releases/v1.0.0.md` and the thirteen explicit accepted asset paths
- [ ] T015 [US1] Immediately read the GitHub Release back and verify its tag, target, name, state, body structure, final changelog link, and exact thirteen-name asset inventory with zero internal evidence or metadata files

**Checkpoint**: User Story 1 is publicly visible and exactly matches the authorized release identity and inventory.

---

## Phase 4: User Story 2 - Independently Verify Publication (Priority: P1)

**Goal**: Prove every public byte is the accepted candidate and separately protected state remains unchanged.

**Independent Test**: Download the public release to a new directory, compare thirteen digests, apply archive checksums, inspect archive and SBOM semantics, execute the compatible binary, and audit protected boundaries.

- [ ] T016 [US2] Move issue #38 to Stage `Release verification`, clear default Status, download all public v1.0.0 assets into a newly created clean temporary directory, and prove the downloaded filename set and sizes equal the thirteen-file contract
- [ ] T017 [US2] Compare all thirteen downloaded SHA-256 values with `release-publication-contract.json` and apply the six-entry checksum manifest as an exact archive bijection
- [ ] T018 [US2] Inspect all six public archives for safe four-member structure, executable identity, schema bytes, legal bytes and modes, release marker, and expected target build identity
- [ ] T019 [US2] Inspect all six public SBOMs for matching platform, version `1.0.0`, source revision, catalog shape, and target association
- [ ] T020 [US2] Extract the public Windows amd64 archive safely, run `version`, `schema --version`, and schema emission non-interactively, and prove exact `1.0.0` output plus byte identity with the tagged immutable schema
- [ ] T021 [US2] Verify milestone v1.0.0 and epic #29 remain open and no production `cueson.io`, public schema-hosting, signature, attestation, or other excluded action occurred

**Checkpoint**: User Story 2 is independently complete and the public release is verified without crossing another protected boundary.

---

## Phase 5: User Story 3 - Record and Close the v1 Delivery State (Priority: P2)

**Goal**: Make the repository and official S020 pull request truthfully describe and govern the completed publication.

**Independent Test**: Inspect all release-facing prose, run offline documentation and repository checks, and verify the issue, Project, pull request, hosted checks, and reviews.

- [ ] T022 [P] [US3] Update current release identity, download, status, capability, and verification links in `README.md` while retaining explicit deferred boundaries
- [ ] T023 [P] [US3] Update released-state, immutable-tag, verified-asset, protected-boundary, and milestone-lifecycle prose in `docs/release-process.md`, `docs/release-verification.md`, and `docs/schema.md`
- [ ] T024 [P] [US3] Update release topology, exact publication ownership, and current delivery status in `docs/architecture.md`, `docs/cli.md`, and `docs/project-management.md`
- [ ] T025 [US3] Record the post-tag documentation and release-verification update under `[Unreleased]` in `CHANGELOG.md` without modifying the tagged v1.0.0 notes or versioned schema
- [ ] T026 [US3] Search maintained documentation for current candidate or unpublished-v1 claims and capability overclaims, then remediate false current-state statements while preserving explicit historical context
- [ ] T027 [US3] Run all quickstart verification including formatter, documentation, product, policy, release, workflow, race, vet, vulnerability, brand-integrity, whitespace, encoding, mojibake, and tagged-file-immutability checks
- [ ] T028 [US3] Run Spec Kit convergence against `specs/S020-publish-v1-release/` and append and execute required remediation tasks until convergence is clean
- [ ] T029 [US3] Obtain explicit branch push and official pull-request publication authority if it has not already been granted
- [ ] T030 [US3] Push `codex/S020-publish-v1-release`, publish the formatted official pull request with `Closes #38`, and verify the rendered GitHub body
- [ ] T031 [US3] Move issue #38 to Stage `PR review`, retain Slice `S020`, and verify default Status remains empty in `cueson Delivery`
- [ ] T032 [US3] Monitor every current-head CI, CodeQL, release-proof, pull-request policy, Codex, security, and review result; address every round-one finding and resolve threads only after remediation
- [ ] T033 [US3] Request exactly one `@codex review` second round only if round one reports findings or a repository-documented condition requires it, then address that round without requesting a third
- [ ] T034 [US3] Re-run convergence and the complete local verification suite after remediation, push the final descendant head, and wait for every hosted gate to return green
- [ ] T035 [US3] Confirm all reviews are satisfied, the pull request is mergeable, issue #38 remains open pending merge, and no pull-request merge, milestone closure, epic closure, production mutation, public schema publication, signature, or attestation occurred

**Checkpoint**: S020 is ready for the operator's final review and merge ritual.

---

## Phase 6: Post-Merge Reconciliation

**Purpose**: Reconcile v1 delivery only after the operator confirms the S020 pull request merged.

- [ ] T036 Verify the merged pull request, fetch and prune, fast-forward a clean non-divergent `main`, and remove only stale clean slice state proven safe to delete
- [ ] T037 Verify post-merge CI, CodeQL, release proof, issue #38 closure, and Project Stage `Done` with Slice `S020` and empty default Status
- [ ] T038 Obtain explicit authority to close epic #29 and milestone v1.0.0, then close only if all linked work and publication evidence remain complete
- [ ] T039 Read back issue, epic, milestone, dependency, Project, release, and tag state and report the final v1.0.0 delivery record

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Starts from the operator's S020 kickoff.
- **Foundational controls (Phase 2)**: Depends on Setup and blocks every publication mutation.
- **User Story 1 (Phase 3)**: Depends on exact publication authority; tag and release mutations are strictly sequential.
- **User Story 2 (Phase 4)**: Depends on the public release and blocks released-state documentation claims.
- **User Story 3 (Phase 5)**: Depends on verified publication and separate publication authority for branch and pull-request mutations.
- **Post-merge reconciliation (Phase 6)**: Depends on a human-confirmed merge and distinct epic and milestone closure authority.

### User Story Dependencies

- **User Story 1**: Produces the authorized public tag, notes, and downloads.
- **User Story 2**: Uses User Story 1's public bytes but independently proves identity and protected-state containment.
- **User Story 3**: Uses verified publication evidence to reconcile repository and GitHub delivery records.

### Parallel Opportunities

- T003, T004, and T005 inspect or record separate preparation surfaces.
- T022, T023, and T024 modify separate documentation files after publication verification.
- Remote tag creation, tag push, release creation, release read-back, and public download verification are intentionally sequential.

## Implementation Strategy

1. Establish issue and Project ownership and freeze the exact candidate and asset contract.
2. Complete specification, planning, task generation, and the blocking analysis gate.
3. Revalidate every accepted input and prove absence of conflicting public state.
4. Halt for exact tag and release publication authority.
5. After authority, create and verify the immutable tag, then publish and read back the complete release.
6. Download all public assets, prove byte identity and compatible execution, and audit excluded boundaries.
7. Reconcile repository documentation without altering tagged history.
8. After push and pull-request authority, complete hosted checks and at most two Codex rounds.
9. Stop at the green, fully reviewed pull request for the operator's final merge ritual.
