# Tasks: Establish CI and Cross-Platform Build Gates

**Input**: Design documents from `specs/S006-establish-ci-gates/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`, `quickstart.md`

**Tests**: S006 requires local command parity, workflow linting, a controlled failing hosted run, a corrected green hosted run, native three-platform execution, and the bounded review protocol.

**Organization**: Tasks are grouped by user story so quality gates, portability evidence, and security evidence remain independently reviewable.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it owns a different file or independent verification surface.
- **[Story]**: Maps implementation tasks to the specification's user stories.
- Every task names its principal file or evidence path.

## Phase 1: Setup and Specification

**Purpose**: Establish the governed S006 branch, scope, research, and design artifacts.

- [X] T001 Create and select `S006-establish-ci-gates`, persist `specs/S006-establish-ci-gates` in `.specify/feature.json`, and verify the branch starts from clean `main`.
- [X] T002 Write and validate `specs/S006-establish-ci-gates/spec.md` and `specs/S006-establish-ci-gates/checklists/requirements.md` with issue #8 acceptance criteria, measurable outcomes, and explicit exclusions.
- [X] T003 [P] Resolve current immutable action commits, compatible analysis-tool versions, runner labels, permission guidance, and failure-proof strategy in `specs/S006-establish-ci-gates/research.md`.
- [X] T004 [P] Complete the technical design, constitution gates, and repository structure in `specs/S006-establish-ci-gates/plan.md` and `specs/S006-establish-ci-gates/data-model.md`.
- [X] T005 [P] Define stable checks, permissions/events, and controlled failure evidence in `specs/S006-establish-ci-gates/contracts/` and runnable local parity in `specs/S006-establish-ci-gates/quickstart.md`.

**Checkpoint**: S006 is fully specified with no unresolved clarification marker.

---

## Phase 2: Foundational Delivery Contract

**Purpose**: Establish the shared automation rules that every user story consumes.

- [X] T006 Verify `.gitignore` already covers Go outputs, runner artifacts, local environments, and editor state without adding unrelated ignore patterns.
- [X] T007 Record workflow names, job names, explicit runner generations, Go compatibility version, tool pins, action commits, and concurrency keys consistently across `specs/S006-establish-ci-gates/plan.md`, `specs/S006-establish-ci-gates/contracts/check-contract.md`, and `specs/S006-establish-ci-gates/contracts/permissions-and-events.md`.
- [X] T008 Validate the generated Spec Kit artifacts with `go run ./scripts/github-format/main.go .`, the built-in checklist in `specs/S006-establish-ci-gates/checklists/requirements.md`, and the blocking Spec Kit analysis gate before implementation.

**Checkpoint**: Stable names and authority boundaries are fixed before workflow implementation.

---

## Phase 3: User Story 1 - Trust Every Proposed Change (Priority: P1) MVP

**Goal**: Every pull request and `main` update receives stable, fail-closed quality evidence for all implemented repository contracts.

**Independent Test**: Lint `.github/workflows/ci.yml`, run every mapped command locally, then prove a draft pull-request revision fails on the controlled marker and the corrected revision passes under the same stable names.

### Implementation and Verification for User Story 1

- [X] T009 [US1] Create `.github/workflows/ci.yml` with `CI` display name, `pull_request` and `main` push triggers, optional manual trigger, `contents: read`, explicit concurrency grouping, and pull-request-only cancellation.
- [X] T010 [US1] Add the stable `Formatting` job to `.github/workflows/ci.yml` for Go formatting and actionlint v1.7.12 workflow validation.
- [X] T011 [US1] Add the stable `Repository text` job to `.github/workflows/ci.yml` for publication-formatter repository checks, nested formatter-module tests, `git diff --check`, and rejection of `.github/ci-failure-probe`.
- [X] T012 [US1] Add the stable `Vet` job to `.github/workflows/ci.yml` for root-module vetting.
- [X] T013 [US1] Add the stable `Schema and conformance` job to `.github/workflows/ci.yml` for schema, lockstep, fixture, path-leak, and conformance packages.
- [X] T014 [US1] Add the stable `Static analysis` job to `.github/workflows/ci.yml` using Go 1.25-compatible Staticcheck v0.7.0 with an explicit version, remove the unused `internal/source` inspection wrapper, and annotate only helpers that are build-selectively used outside Windows.
- [X] T015 [US1] Add the stable `Vulnerability scan` job to `.github/workflows/ci.yml` using govulncheck v1.8.0 with an explicit version and a current stable analysis runtime, then remediate reachable findings by raising `go.mod` to Go 1.25 and upgrading `golang.org/x/text` to v0.39.0.
- [X] T016 [US1] Add the stable `Race detection` job to `.github/workflows/ci.yml` and run the complete root suite with the race detector on Linux amd64.
- [X] T017 [US1] Run actionlint and every local User Story 1 command from `specs/S006-establish-ci-gates/quickstart.md`, then reconcile any mismatch in `.github/workflows/ci.yml` or the check contract.

**Checkpoint**: General CI is locally valid, least-privilege, stable-named, and ready for the hosted red/green proof.

---

## Phase 4: User Story 2 - Prove Portable Native Behavior and Builds (Priority: P2)

**Goal**: Native behavior executes on Windows, macOS, and Linux, while temporary pure-Go builds prove all six supported target pairs.

**Independent Test**: Inspect one corrected hosted run for three native root-suite jobs and six successful `CGO_ENABLED=0` build jobs, with no artifact upload.

### Implementation and Verification for User Story 2

- [X] T018 [US2] Add a closed native-test include matrix to `.github/workflows/ci.yml` with stable Linux, Windows, and macOS job names on `ubuntu-24.04`, `windows-2025`, and `macos-15`.
- [X] T019 [US2] Configure each native-test matrix entry in `.github/workflows/ci.yml` to use the latest Go 1.25 patch and run the complete root-module suite in the foreground.
- [X] T020 [US2] Add a closed six-target pure-Go build include matrix to `.github/workflows/ci.yml` for `linux`, `windows`, and `darwin` on `amd64` and `arm64` with stable job names.
- [X] T021 [US2] Configure `.github/workflows/ci.yml` cross-build outputs beneath the runner temporary directory with `CGO_ENABLED=0`, target-specific executable suffixes, `-trimpath`, and no upload or publication step.
- [X] T022 [US2] Run local Windows tests and all six non-publishing cross-build commands from `specs/S006-establish-ci-gates/quickstart.md`, then confirm no output enters the repository.

**Checkpoint**: Native and build claims are distinct, complete, and locally reproducible where the current host permits.

---

## Phase 5: User Story 3 - Receive Independent Security Evidence (Priority: P3)

**Goal**: CodeQL supplies independently visible security evidence with narrowly scoped reporting authority.

**Independent Test**: Lint `.github/workflows/codeql.yml`, read back its events and permissions, and verify the corrected pull-request run reports `CodeQL / Analyze Go` successfully.

### Implementation and Verification for User Story 3

- [X] T023 [US3] Create `.github/workflows/codeql.yml` with `CodeQL` display name, pull-request and `main` push triggers, a weekly schedule, isolated concurrency, and no `pull_request_target` or secrets.
- [X] T024 [US3] Declare only `contents: read` and `security-events: write` in `.github/workflows/codeql.yml` and use the stable `Analyze Go` job name.
- [X] T025 [US3] Pin CodeQL initialization and analysis in `.github/workflows/codeql.yml` to verified v4.38.0 commit `b96794f015dfd88f77b49b1c93e0fa7110f94c63` and configure Go analysis without artifact or repository mutation.
- [X] T026 [US3] Run actionlint against `.github/workflows/codeql.yml` and audit every `uses:`, permission, event, and concurrency declaration against `specs/S006-establish-ci-gates/contracts/permissions-and-events.md`.

**Checkpoint**: Independent CodeQL automation is syntactically valid and no broader permission exists.

---

## Phase 6: Documentation, Analysis, and Local Convergence

**Purpose**: Align public status, architecture, release history, Spec Kit evidence, and the intended final tree before publication.

- [X] T027 [P] Replace the planned CI badge and hosted-CI status prose with the stable workflow badge and truthful S006 boundary in `README.md`.
- [X] T028 [P] Add the stable hosted automation responsibility, native-platform evidence, pure-Go build proof, security boundary, and deferred checks to `docs/architecture.md`, then align the staged gate list in `docs/Cueson-Project-Specification-v0.0.0.md`.
- [X] T029 [P] Record the S006 automation addition and draft-roadmap refinement under `[Unreleased]` in `CHANGELOG.md`.
- [X] T030 Re-run Spec Kit analysis against `specs/S006-establish-ci-gates/spec.md`, `specs/S006-establish-ci-gates/plan.md`, and `specs/S006-establish-ci-gates/tasks.md`; resolve every blocking inconsistency before implementation completion.
- [X] T031 Run the complete foreground validation sequence from `specs/S006-establish-ci-gates/quickstart.md`, including root tests, race, vet, formatter checks, actionlint, Staticcheck, govulncheck, six cross-builds, UTF-8/BOM/mojibake checks, line-ending checks, and `git diff --check`.
- [X] T032 Run Spec Kit convergence against the intended final tree and append any remaining work only to `specs/S006-establish-ci-gates/tasks.md`; implement appended tasks and repeat until clean.
- [X] T033 Update issue #8 to Project `Slice: S006` and `Stage: In progress`, keep default `Status` empty, and record the read-back evidence in `specs/S006-establish-ci-gates/tasks.md`.

**Checkpoint**: The intended final tree is locally green and converged before the controlled hosted failure proof.

---

## Phase 7: Authorized Publication and Hosted Red/Green Proof

**Purpose**: Use the operator's explicit authorization to publish S006, prove real failure behavior, correct it, and expose only the final green tree to automated reviewers.

- [ ] T034 Commit the intended final implementation and Spec Kit evidence on `S006-establish-ci-gates`, then add `.github/ci-failure-probe` in a separate controlled-failure commit.
- [ ] T035 Push `S006-establish-ci-gates`, publish the official pull request as a draft with a formatter-verified body containing `Closes #8`, and read the body back from GitHub.
- [ ] T036 Verify the hosted `CI / Repository text` check and overall `CI` workflow fail specifically because `.github/ci-failure-probe` exists, then retain the failed run URL in `specs/S006-establish-ci-gates/tasks.md`.
- [ ] T037 Remove `.github/ci-failure-probe`, commit and push the correction, and verify the final pull-request diff contains no probe.
- [ ] T038 Wait for every corrected CI and CodeQL check in `specs/S006-establish-ci-gates/contracts/check-contract.md` to complete successfully within the 20-minute success-criterion window, record elapsed time, investigate every failure, and push fixes until green.
- [ ] T039 Mark the pull request ready for review only after the corrected checks are green so third-party Codex and security bots evaluate the intended tree.

**Checkpoint**: Hosted failure and recovery are proven, the final tree is green, and automated review may begin.

---

## Phase 8: Bounded Review Closure and Final Handoff

**Purpose**: Address every external finding under the recorded maximum-two-round protocol and stop before human merge.

- [ ] T040 Inspect every first-round review, inline thread, security-bot result, pull-request reaction, and CI check; respond to every actionable finding and record remediation tasks in `specs/S006-establish-ci-gates/tasks.md`.
- [ ] T041 Implement every valid first-round remediation, rerun the complete relevant foreground verification, push the fixes, reply with commit evidence, and resolve threads only after the concern is handled.
- [ ] T042 If round one had findings, post exactly one formatter-verified `@codex review` request, read it back, and do not request any third automated review; if round one was clean, record why no second request was needed in `specs/S006-establish-ci-gates/tasks.md`.
- [ ] T043 Inspect and resolve every second-round finding and security result with the same fix, verification, evidence, push, reply, and thread-resolution discipline in `specs/S006-establish-ci-gates/tasks.md`.
- [ ] T044 Re-run final Spec Kit convergence and the complete applicable verification suite, mark all completed tasks in `specs/S006-establish-ci-gates/tasks.md`, and push final evidence updates.
- [ ] T045 Set issue #8 Project `Stage: PR review` with `Slice: S006` and empty default `Status`, verify all checks and review threads are clear, and notify the operator for the final review and merge ritual without merging.

---

## Dependencies and Execution Order

### Phase Dependencies

- **Phase 1** establishes the Spec Kit source of intent.
- **Phase 2** fixes shared names and authority boundaries before workflows are written.
- **User Story 1** supplies the MVP quality workflow and blocks hosted red/green proof.
- **User Story 2** and **User Story 3** can proceed in parallel after Phase 2 because they own distinct workflow sections/files, but both must finish before publication.
- **Phase 6** integrates documentation and local verification after all three stories.
- **Phase 7** depends on a clean local convergence result and the operator's explicit push/PR authorization.
- **Phase 8** begins only after the corrected hosted checks are green and stops before final merge.

### User Story Dependencies

- **User Story 1 (P1)**: Depends only on the shared delivery contract and provides the failure-proof gate.
- **User Story 2 (P2)**: Depends only on the shared delivery contract and extends `ci.yml`; coordinate edits sequentially with User Story 1 despite conceptual independence.
- **User Story 3 (P3)**: Depends only on the shared delivery contract and owns `codeql.yml`, so it can be implemented independently of `ci.yml` edits.

### Parallel Opportunities

- T003, T004, and T005 own independent Spec Kit artifacts.
- T023 through T026 own `codeql.yml` and can proceed alongside the `ci.yml` work when separate agents use explicit file ownership.
- T027, T028, and T029 own separate documentation files.
- Native hosted jobs and six build jobs execute in parallel after workflow publication.

---

## Parallel Example: User Stories 2 and 3

```text
Task: "Implement native and cross-build matrices in .github/workflows/ci.yml"
Task: "Implement independent least-privilege CodeQL workflow in .github/workflows/codeql.yml"
```

---

## Implementation Strategy

### MVP First

1. Complete the Spec Kit and shared-contract phases.
2. Implement and locally validate User Story 1.
3. Confirm every implemented repository invariant has a stable, fail-closed gate.
4. Add native/build and CodeQL evidence without weakening the MVP.

### Incremental Delivery

1. Establish stable CI events, permissions, names, and concurrency.
2. Add repository quality and analysis jobs.
3. Add native operating-system and pure-Go build matrices.
4. Add independent CodeQL scanning.
5. Align public documentation and converge locally.
6. Publish a draft, prove red then green, expose the corrected tree to reviews, and stop for human merge.

### Review Discipline

- The draft failure-probe revision is delivery evidence, not a review target.
- Every first-round finding is handled before the sole permitted second request.
- A clean first round receives no second request.
- A second round, when needed, is the final automated round.
- No AI merge, auto-merge, release, tag, ruleset, or production mutation occurs.

## Notes

- S006 closes issue #8 only.
- Issue #9 owns issue-link and Codex-review automation, issue #10 owns repository settings and required checks, issue #11 owns release proof, and later codec slices own native subtitle and conversion checks.
- The live badge is structurally validated in S006 and visually rechecked on `main` during post-merge housekeeping.

## Evidence

- 2026-09-09 local verification passed repository formatting and text checks, nested formatter tests, vet, root tests, race tests, focused schema/conformance tests, Staticcheck v0.7.0, govulncheck v1.8.0 with zero affecting symbol-level vulnerabilities, actionlint v1.7.12, all six `CGO_ENABLED=0` cross-builds, and staged plus unstaged whitespace checks.
- 2026-09-09 independent specification, workflow, and security audits were reconciled by removing stale issue #8 future-tense architecture prose, making quickstart verification fail closed, enumerating all six local cross-builds, measuring hosted duration, upgrading the analyzer pins, disabling checkout credential persistence, and recording durable vulnerability evidence.
- 2026-09-09 Spec Kit convergence found all 17 functional requirements and 8 success criteria covered by the chronological task plan, with no unresolved clarification marker, constitution conflict, or missing pre-publication work. T034 through T045 remain intentionally pending because they own authorized publication, hosted proof, automated review, and final handoff.
- 2026-09-09 Project read-back confirmed issue #8 item `PVTI_lADOBpohEc4Bi59Izg6FU8k` has `Slice: S006`, `Stage: In progress`, and an empty default `Status` field.
