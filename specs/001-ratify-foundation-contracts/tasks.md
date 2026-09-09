# Tasks: Ratify Foundation Contracts

**Input**: Design documents from `/specs/001-ratify-foundation-contracts/`

**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/](contracts/), [quickstart.md](quickstart.md)

**Tests**: Contract and documentation assertions run before the corresponding edits and again during final verification. No shipped product tests or code are added.

**Organization**: Tasks are grouped by user story so the ratified baseline, public README status, and GitHub tracking state remain independently reviewable.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it owns a different file or external record.
- **[Story]**: Maps the task to the corresponding prioritized user story.
- Every task names its repository file or GitHub record.

## Phase 1: Setup

**Purpose**: Establish the feature context and repeatable pre-implementation evidence.

- [X] T001 Confirm branch and feature resolution through `.specify/feature.json` and `specs/001-ratify-foundation-contracts/spec.md`
- [X] T002 Confirm the scope and authority boundaries from `.specify/memory/constitution.md`, `AGENTS.md`, and GitHub issues #1, #3, and #6

---

## Phase 2: Foundational Decisions

**Purpose**: Complete blocking decision records and prove the current repository does not yet satisfy the slice.

- [X] T003 Reconcile all research decisions with `specs/001-ratify-foundation-contracts/contracts/architecture-baseline.md`, `schema-baseline.md`, and `cli-baseline.md`
- [X] T004 Run pre-implementation assertions from `specs/001-ratify-foundation-contracts/quickstart.md` and record that `docs/architecture.md`, `docs/schema.md`, `docs/cli.md`, and truthful badge targets are initially absent

**Checkpoint**: Contract decisions are stable and the expected failing baseline is demonstrated before publication-document edits.

---

## Phase 3: User Story 1 - Implement from one ratified baseline (Priority: P1) MVP

**Goal**: Publish concise, mutually consistent architecture, schema, and CLI contracts that issues #4 through #7 can implement directly.

**Independent Test**: Trace every acceptance criterion in issues #4 through #7 to a normative section in the three ratified documents and find no unresolved alternative for the decisions named by issue #3.

- [X] T005 [P] [US1] Author package responsibilities, dependency direction, source authority, and deferred ownership in `docs/architecture.md`
- [X] T006 [P] [US1] Author schema identity, lifecycle, common/source model, format support, and OCR cardinality in `docs/schema.md`
- [X] T007 [P] [US1] Author implemented-command visibility, stream policy, exit codes, and capability diagnostics in `docs/cli.md`
- [X] T008 [US1] Link the ratified documents and resolve touched provisional contradictions in `docs/Cueson-Project-Specification-v0.0.0.md`
- [X] T009 [US1] Amend GitHub issue #6 so its Scope and acceptance criteria own the generic exact-restoration command required by `docs/schema.md` and `docs/cli.md`
- [X] T010 [US1] Run the contract traceability assertions in `specs/001-ratify-foundation-contracts/quickstart.md` against `docs/architecture.md`, `docs/schema.md`, and `docs/cli.md`

**Checkpoint**: The contract baseline is independently usable without consulting chat history or provisional alternatives.

---

## Phase 4: User Story 2 - See truthful repository status (Priority: P2)

**Goal**: Make public status badges truthful before CI and release surfaces exist.

**Independent Test**: Inspect `README.md` and confirm the CI and release images are static state badges linked to issue #8 and milestone #1, with no reference to missing `ci.yml` or `releases/latest` resources.

- [X] T011 [US2] Replace unavailable CI and release dynamic badge resources with planned and unreleased static badges in `README.md`
- [X] T012 [US2] Verify the four badge image and link contracts plus the pre-release status statement in `README.md`

**Checkpoint**: The README communicates current status without missing-resource badge requests.

---

## Phase 5: User Story 3 - Trust delivery status (Priority: P3)

**Goal**: Reconcile evidence-backed GitHub state without claiming an unmerged change is complete or public.

**Independent Test**: Read issues #1, #3, and #6 plus Project 3 and confirm the pending bootstrap criterion, active slice, amended restore ownership, and absence of a public pull request.

- [X] T013 [P] [US3] Update GitHub issue #1 with checked completed criteria, verification evidence, and one explicitly pending README merge criterion while leaving it open and In progress
- [X] T014 [P] [US3] Set GitHub issue #3 Project Stage to In progress and Slice to `001-ratify-foundation-contracts` while leaving the issue open
- [X] T015 [US3] Read back issues #1, #3, and #6 plus Project 3 and open pull requests to validate GitHub state against `specs/001-ratify-foundation-contracts/quickstart.md`

**Checkpoint**: Remote planning state matches evidence and all protected publication boundaries remain intact.

---

## Phase 6: Polish and Cross-Cutting Verification

**Purpose**: Record decisions, complete traceability, and prove local readiness for review.

- [X] T016 Record the contract, restore-boundary, badge, and bootstrap-closure decisions under Unreleased in `CHANGELOG.md`
- [X] T017 Re-run formatter tests, repository formatting, whitespace, encoding, documentation-link, and Spec Kit checks from `specs/001-ratify-foundation-contracts/quickstart.md`
- [X] T018 Confirm every required task is complete, every issue #3 acceptance criterion has local evidence, and no push, tag, release, production-domain change, or public pull request occurred

---

## Phase 7: Round-One Review Resolution

**Purpose**: Resolve every actionable first-round Codex finding before the single permitted second-round request.

- [X] T019 Inspect and validate every first-round Codex finding against the ratified contracts and working project specification
- [X] T020 Amend the portable safe-basename contract in `docs/architecture.md`, `docs/schema.md`, the working project specification, and the Spec Kit decision artifacts
- [X] T021 Move generic exact restoration into the v0.0.0 roadmap while retaining codec-integrated round trips in the 0.x series
- [X] T022 Amend GitHub issue #6 with the portable safe-basename acceptance boundary and read the published body back
- [X] T023 Rerun the complete verification suite, prepare the review-resolution commit, and confirm both findings are addressed before the single permitted second-round request

---

## Dependencies & Execution Order

### Phase Dependencies

- Setup has no dependencies.
- Foundational decisions depend on Setup and block all user stories.
- User Story 1, User Story 2, and the initial portions of User Story 3 can proceed independently after the foundation.
- T015 depends on T009, T013, and T014 because it validates their combined remote state.
- Polish and cross-cutting verification depends on all selected user stories.

### User Story Dependencies

- **User Story 1 (P1)**: No dependency on another story; it is the MVP and issue #3's primary outcome.
- **User Story 2 (P2)**: No dependency on User Story 1; it independently corrects public status.
- **User Story 3 (P3)**: Remote reconciliation can begin independently, but issue #1 intentionally remains pending until User Story 2 eventually merges to `main`.

### Parallel Opportunities

- T005, T006, and T007 own separate ratified documents and can be authored in parallel after T003 and T004.
- T011 can proceed independently of T005 through T010.
- T013 and T014 own separate GitHub issue records and can proceed in parallel after the foundation.
- T016 waits for all decisions so the changelog records final rather than provisional choices.

## Parallel Example: User Story 1

```text
Task T005: Author docs/architecture.md from contracts/architecture-baseline.md
Task T006: Author docs/schema.md from contracts/schema-baseline.md
Task T007: Author docs/cli.md from contracts/cli-baseline.md
```

## Implementation Strategy

### MVP First

1. Complete Setup and Foundational Decisions.
2. Complete User Story 1 and its traceability check.
3. Validate that issues #4 through #7 have one authoritative baseline.

### Incremental Delivery

1. Add truthful badge status through User Story 2.
2. Reconcile remote issue and Project evidence through User Story 3.
3. Complete cross-cutting verification and create one local conventional commit.
4. Halt before pushing the feature branch or opening a public pull request.

## Notes

- Every task uses an exact file path or GitHub record.
- No task adds shipped Go code or prematurely implements issue #4 or later product work.
- Completed tasks are marked `[X]` during `/speckit-implement`.
- The exact push command is presented only at the autopilot halt.
