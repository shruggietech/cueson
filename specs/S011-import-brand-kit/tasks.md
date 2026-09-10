# Tasks: Import the Official Cueson Brand Kit

**Input**: Design documents from `specs/S011-import-brand-kit/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/import-manifest.schema.json`

**Tests**: Integrity, formatter, documentation, integration, and hosted verification are mandatory because the specification requires deterministic byte-preservation and CI evidence.

**Organization**: Tasks are grouped by setup, blocking foundations, and independently testable user stories.

## Phase 1: Setup

**Purpose**: Align GitHub planning and establish the exact S011 acquisition boundary.

- [x] T001 Revise GitHub issue #23 to the operator-approved complete-import outcome using the repository publication formatter and verify the rendered body after publication
- [x] T002 Set issue #23 to Slice `S011`, Stage `In progress`, and an empty default Status in the `cueson Delivery` Project
- [x] T003 [P] Record the two-download acquisition and safety audit in `specs/S011-import-brand-kit/research.md`
- [x] T004 [P] Verify `.gitattributes` and `.editorconfig` preserve `brand/**` while retaining CRLF for `*.ps1`

---

## Phase 2: Foundational Byte-Integrity Controls

**Purpose**: Prevent repository tooling from altering the kit and establish the verifier contract before importing payload bytes.

**CRITICAL**: No retained payload or consumer integration is complete until these controls pass.

- [x] T005 Run the blocking Spec Kit analysis gate and remediate every critical, high, or constitution conflict before implementation begins
- [x] T006 [P] Add protected-tree regression tests for check and fix behavior in `scripts/github-format/main_test.go`
- [x] T007 Implement byte-protected directory exclusion for `brand/`, `docs/assets/brand/`, `docs/assets/favicons/`, and `docs/assets/fonts/` in `scripts/github-format/main.go`
- [x] T008 Create the dependency-free Go 1.25 module skeleton in `scripts/brand-verify/go.mod` and `scripts/brand-verify/main.go`
- [x] T009 [P] Add strict manifest, safe-path, malicious-ZIP, extraction-drift, and reference-verification tests in `scripts/brand-verify/verify_test.go`
- [x] T010 Implement strict import-manifest decoding, portable path validation, bounded ZIP inspection, extraction bijection, hashing, and repository-reference checks in `scripts/brand-verify/verify.go`
- [x] T011 Run focused formatter and brand-verifier tests and confirm malicious and drift fixtures fail deterministically

**Checkpoint**: Protected paths cannot be rewritten and the offline verifier contract is ready for real content.

---

## Phase 3: User Story 1 - Retain the Authoritative Brand Kit (Priority: P1)

**Goal**: Commit the official archive and every safe payload byte with complete independent acquisition evidence.

**Independent Test**: Run the brand verifier against the retained archive, manifest, and extraction and receive a successful 265-file bijection result.

- [x] T012 [US1] Copy the independently verified ZIP to `brand/cueson/1.0.0/archive/cueson-brand-1.0.0.zip`
- [x] T013 [US1] Extract all 265 audited regular files without byte changes beneath `brand/cueson/1.0.0/kit/`
- [x] T014 [US1] Generate the path-sorted 265-entry acquisition inventory and metadata in `brand/cueson/1.0.0/import-manifest.json`
- [x] T015 [US1] Verify archive size and SHA-256, safe catalog, extraction file set, entry sizes, and entry digests with `scripts/brand-verify`
- [x] T016 [US1] Verify representative Git attributes for ZIP, SVG, JSON, CSS, XML, TXT, font, and PowerShell paths and confirm no retained byte changed after Git staging

**Checkpoint**: User Story 1 is independently complete and the full official kit is retained exactly.

---

## Phase 4: User Story 2 - Use Official Assets in Repository Surfaces (Priority: P1)

**Goal**: Adopt the retained identity in primary repository documentation while remaining offline and truthful.

**Independent Test**: Render the README in light and dark modes and the media guide offline, then verify every declared asset reference resolves into the retained kit.

- [x] T017 [P] [US2] Add responsive official light and dark horizontal identity assets, exact slogan, canonical description, and brand links to `README.md`
- [x] T018 [P] [US2] Replace the media guide's remote fonts, ad hoc lockup, palette, and endorsement with retained official assets and roles in `docs/cueson-media-format-guide.html`
- [x] T019 [US2] Add the README and media-guide asset references to `brand/cueson/1.0.0/import-manifest.json`
- [x] T020 [US2] Verify direct asset references, local font loading, favicon loading, light and dark readability, narrow layout, print behavior, and unchanged capability disclaimer

**Checkpoint**: User Story 2 is independently complete with official local identity and no network dependency.

---

## Phase 5: User Story 3 - Review a Safe, Reproducible Import (Priority: P2)

**Goal**: Document, enforce, publish, and review the complete S011 outcome without crossing merge, release, or production boundaries.

**Independent Test**: Run every local and hosted gate, inspect GitHub issue and Project state, and confirm the pull request is green, fully reviewed, and configured to close issue #23.

- [x] T021 [P] [US3] Add complete acquisition, usage, provenance, licensing, verification, update, and release-boundary guidance in `docs/brand.md`
- [x] T022 [P] [US3] Add `docs/brand.md` to the required inventory and update expected counts in `scripts/docs-verify/verify.go` and `scripts/docs-verify/verify_test.go`
- [x] T023 [P] [US3] Document immutable brand ownership and verification boundaries in `docs/architecture.md` and `CONTRIBUTING.md`
- [x] T024 [P] [US3] Reconcile brand state and the v1 brand-kit gate in `docs/Cueson-Project-Specification-v0.0.0.md`
- [x] T025 [P] [US3] Record retained-kit attribution in `NOTICE` and S011 additions, changes, and decisions in `CHANGELOG.md`
- [x] T026 [US3] Add the brand-verifier module cache, tests, and offline run to `.github/workflows/ci.yml`
- [x] T027 [US3] Run all quickstart verification including formatter, brand, documentation, product, policy, release, workflow, race, vet, whitespace, encoding, and mojibake checks
- [x] T028 [US3] Run Spec Kit convergence and append and execute any required remediation tasks until convergence is clean
- [x] T029 [US3] Commit and push the verified S011 branch, publish the formatted official pull request with `Closes #23`, and verify its rendered body
- [x] T030 [US3] Move issue #23 to Stage `PR review`, retain Slice `S011`, and verify default Status remains empty
- [ ] T031 [US3] Monitor every current-head CI, CodeQL, policy, Codex, and security result; address every round-one finding and resolve threads only after remediation
- [ ] T032 [US3] Request exactly one `@codex review` second round only if round one reports findings, then address that round without requesting a third
- [ ] T033 [US3] Confirm every required check is green, all reviews are satisfied, the pull request is mergeable, issue #23 remains open pending merge, and no merge, tag, release, schema publication, milestone closure, or production mutation occurred

**Checkpoint**: S011 is ready for the operator's final review and merge ritual.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Starts immediately from the approved specification and operator authority.
- **Foundational controls (Phase 2)**: Depends on Setup and blocks payload import completion.
- **User Story 1 (Phase 3)**: Depends on the verifier and protected-path controls.
- **User Story 2 (Phase 4)**: Depends on the retained payload paths from User Story 1.
- **User Story 3 (Phase 5)**: Documentation tasks can begin after the retained layout is fixed; publication and hosted review depend on all local implementation and verification.

### User Story Dependencies

- **User Story 1**: Independently proves full retained-kit integrity.
- **User Story 2**: Uses User Story 1 assets directly but remains independently testable as a rendering and reference outcome.
- **User Story 3**: Integrates enforcement and evidence for User Stories 1 and 2 and owns GitHub publication and review.

### Parallel Opportunities

- T003 and T004 can run together.
- T005 and T008 can be authored in parallel before their implementations.
- T016 and T017 touch separate documentation surfaces.
- T020 through T024 touch separate maintained documents and verifier files.

## Implementation Strategy

1. Simplify the active issue and establish governed Project state.
2. Add formatter protection and a failing verifier test suite before importing source bytes.
3. Retain and prove the complete acquisition as the minimum viable outcome.
4. Adopt only the assets needed by current repository surfaces through direct references.
5. Complete documentation and CI integration, run analysis and convergence, then publish.
6. Remediate review findings, use at most one second Codex request, and stop only at the green fully reviewed boundary.
