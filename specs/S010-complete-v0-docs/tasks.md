# Tasks: v0.0.0 Documentation and Milestone Verification

**Input**: Design documents from `specs/S010-complete-v0-docs/`

**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/](contracts/), and [quickstart.md](quickstart.md)

**Tests**: S010 requires test-first deterministic documentation inventory and local-link validation plus full repository, hosted, and delivery-evidence verification.

**Organization**: Tasks are grouped by independently testable user story. Documentation verification remains a standalone repository tool and does not enter the shipped product dependency graph.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it owns a different file and does not depend on incomplete work
- **[Story]**: Maps the task to a user story in [spec.md](spec.md)
- Every task names its implementation, documentation, evidence, or verification path

## Phase 1: Setup

**Purpose**: Establish the isolated documentation-verification boundary and confirm repository exclusions.

- [x] T001 Create the dependency-free standalone documentation-verifier module in `scripts/docs-verify/go.mod`
- [x] T002 Verify Go build-output exclusions remain complete for the new standalone module in `.gitignore`

---

## Phase 2: Foundational Documentation Model

**Purpose**: Define the required-document inventory, audited scope, repository path index, and command behavior used by every story.

- [x] T003 Write failing tests for required-document presence, regular-file and symlink rejection, audited root/docs scope, strict argument handling, exit codes, and deterministic diagnostics in `scripts/docs-verify/verify_test.go`
- [x] T004 Implement strict read-only argument parsing and stable foreground reporting in `scripts/docs-verify/main.go`
- [x] T005 Implement the canonical twelve-document inventory, maintained documentation discovery, exact-case repository index, file-size bound, and deterministic violation model in `scripts/docs-verify/verify.go`

**Checkpoint**: The standalone module builds and its inventory and command-contract tests pass without parsing documentation links.

---

## Phase 3: User Story 1 - Trust the v0.0.0 Documentation (Priority: P1) 🎯 MVP

**Goal**: Reject missing, unsafe, malformed, or unresolved local documentation references and eliminate unsupported current-state claims.

**Independent Test**: Run the verifier against accepted and rejected temporary repositories, then audit current documents against executable, schema, automation, and GitHub evidence.

### Tests for User Story 1

- [x] T006 [US1] Add failing extraction tests for inline Markdown links and images, reference definitions and uses, autolinks, HTML `href`/`src`, and ignored fenced code, inline code, comments, and CSS imports in `scripts/docs-verify/verify_test.go`
- [x] T007 [US1] Add failing resolution tests for relative, parent, repository-root, directory, percent-encoded, fragment-only, duplicate/setext-heading, explicit-HTML-anchor, traversal, wrong-case, symlink, control, backslash, unsupported-scheme, missing-target, and missing-fragment cases in `scripts/docs-verify/verify_test.go`

### Implementation for User Story 1

- [x] T008 [US1] Implement one-line Markdown destination parsing, reference resolution, autolink extraction, and code/comment exclusion in `scripts/docs-verify/verify.go`
- [x] T009 [US1] Implement HTML attribute and explicit-anchor extraction plus entity decoding without treating CSS URLs as documentation links in `scripts/docs-verify/verify.go`
- [x] T010 [US1] Implement repository-confined exact-case path resolution, external-link classification without network access, GitHub-style heading slugs, duplicate anchors, and sorted line-specific diagnostics in `scripts/docs-verify/verify.go`
- [x] T011 [US1] Add documentation-verifier cache input, tests, and repository acceptance to the existing `Repository text` job in `.github/workflows/ci.yml`
- [x] T012 [P] [US1] Update repository status, source-only usage, protection state, and contribution posture in `README.md` and `CONTRIBUTING.md`
- [x] T013 [P] [US1] Reconcile implemented, reserved, and future boundaries plus format-document links in `docs/architecture.md`, `docs/schema.md`, and `docs/cli.md`
- [x] T014 [P] [US1] Replace stale bootstrap and pre-merge statements with durable current evidence in `docs/project-management.md`, `docs/repository-controls.md`, and `docs/release-verification.md`
- [x] T015 [P] [US1] Add an unmistakable v0.0.0 envelope-only disclaimer and recast present-tense format capability claims as roadmap intent in `docs/cueson-media-format-guide.html`

**Checkpoint**: User Story 1 independently validates the maintained documentation graph and leaves zero unsupported current capability or delivery claims.

---

## Phase 4: User Story 2 - Navigate the Required Documentation Set (Priority: P1)

**Goal**: Supply the missing canonical format and release-process pages with clear current, future, and protected-action boundaries.

**Independent Test**: Navigate from the README to each new page, resolve every local file and heading link, and compare current capability tables with the schema, CLI, architecture, and release verifier.

### Implementation for User Story 2

- [x] T016 [P] [US2] Create the SubRip envelope-only status, planned v1 grammar/fidelity matrix, diagnostics, fixtures boundary, and explicit exclusions in `docs/formats/srt.md`
- [x] T017 [P] [US2] Create the WebVTT envelope-only status, planned v1 structure/rolling-caption/fidelity matrix, diagnostics, fixtures boundary, and explicit exclusions in `docs/formats/webvtt.md`
- [x] T018 [P] [US2] Create the candidate, authorization, tag/release, immutable-schema, verification, milestone, and post-release lifecycle with no executable publishing path in `docs/release-process.md`
- [x] T019 [US2] Add direct format, release-process, and documentation-verifier navigation without false installation or support claims in `README.md`, `docs/architecture.md`, `docs/schema.md`, and `docs/cli.md`
- [x] T020 [US2] Run the standalone verifier against the complete canonical documentation set and correct every reported inventory, file, and anchor violation through `scripts/docs-verify/main.go`

**Checkpoint**: User Story 2 independently provides all twelve canonical documents and a complete offline-resolving navigation path.

---

## Phase 5: User Story 3 - Review the Completed Foundation Milestone (Priority: P2)

**Goal**: Present durable evidence for every foundation outcome while keeping merge-complete and release-gated state truthful.

**Independent Test**: Audit issues #1 through #12, epic relationships, Project fields, milestone membership, version lockstep, current-head checks, and non-publication state against the milestone contract.

### Implementation for User Story 3

- [x] T021 [P] [US3] Make `[Unreleased]` the truthful v0.0.0 candidate record, remove premature release/tag links, and record S010 documentation and verification decisions in `CHANGELOG.md`
- [x] T022 [P] [US3] Mark only proven v0.0.0 foundation criteria complete, add evidence references, correct the candidate-changelog requirement, and leave release/v1 gates incomplete in `docs/Cueson-Project-Specification-v0.0.0.md`
- [x] T023 [US3] Format, publish, and read back durable post-merge acceptance evidence for closed issues #6, #8, #9, #10, and #11 using `scripts/github-format/main.go` and `specs/S010-complete-v0-docs/contracts/milestone-contract.md`
- [x] T024 [US3] Reconcile issues #12 and #2 to Slice `S010`, governed Stage `In progress`, and unused default Status using `specs/S010-complete-v0-docs/contracts/milestone-contract.md`
- [x] T025 [US3] Verify the issue #1-#12 child graph, dependencies, closing pull requests, milestone inventory, Project uniqueness, software/schema lockstep, post-merge CI, release proof, and absence of tags/releases against `specs/S010-complete-v0-docs/contracts/milestone-contract.md`
- [x] T026 [US3] Update only already-proven acceptance checkboxes and verification evidence on issues #12 and #2 while leaving merge-dependent child closure and milestone completion pending in `specs/S010-complete-v0-docs/contracts/milestone-contract.md`

**Checkpoint**: User Story 3 provides operator-readable foundation evidence while issues #12/#2 and the milestone remain open until their separately governed transitions.

---

## Phase 6: Polish and Cross-Cutting Verification

**Purpose**: Reconcile Spec Kit artifacts, repository quality, GitHub publication, and bounded review.

- [x] T027 Run Go formatting, repository text and documentation checks, root and all standalone module tests, race detection, vet, workflow linting, CLI/schema lockstep probes, `git diff --check`, UTF-8/BOM/line-ending checks, and mojibake scans using `specs/S010-complete-v0-docs/quickstart.md`
- [x] T028 Run Spec Kit convergence against `spec.md`, `plan.md`, and `tasks.md`; append and implement any remaining work in `specs/S010-complete-v0-docs/tasks.md`
- [x] T029 Confirm issues #12 and #2 remain open at Stage `In progress` with Slice `S010`, default Status unused, milestone `v0.0.0` open, and no tag/release/schema-copy/production mutation using `specs/S010-complete-v0-docs/contracts/milestone-contract.md`
- [ ] T030 Format the official pull-request body through `scripts/github-format`, publish it with `Closes #12` and `Closes #2`, read it back, verify rendering, and move both Project items to `PR review` using `specs/S010-complete-v0-docs/contracts/milestone-contract.md`
- [ ] T031 Wait for all hosted checks and first-round third-party reviews, address every finding, request at most one `@codex review` second round only when required, and stop for the operator's final merge ritual after all checks and reviews are satisfied using `specs/S010-complete-v0-docs/contracts/milestone-contract.md`

---

## Dependencies and Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Starts after approved S010 specification and planning artifacts exist.
- **Foundational (Phase 2)**: Depends on setup and blocks link validation and repository integration.
- **User Story 1 (Phase 3)**: Depends on the verifier foundation and creates the trusted documentation graph.
- **User Story 2 (Phase 4)**: Can draft the three missing documents after planning, then requires User Story 1 for final navigation acceptance.
- **User Story 3 (Phase 5)**: Requires truthful documentation and verifier results before milestone evidence can be declared complete.
- **Polish (Phase 6)**: Depends on all user stories and owns convergence, publication, and review monitoring.

### User Story Dependencies

- **User Story 1 (P1)**: Provides the independently testable documentation-trust MVP.
- **User Story 2 (P1)**: Its files can be drafted in parallel, but complete link acceptance uses User Story 1's verifier.
- **User Story 3 (P2)**: Depends on both P1 stories so issue and epic evidence cannot get ahead of repository truth.

### Within Each User Story

- Verifier rejection tests are written and observed failing before their implementation.
- Repository discovery precedes link extraction; extraction precedes resolution; resolution precedes CI integration.
- Format and release pages precede cross-document navigation updates.
- Local documentation and product verification precede mutable GitHub issue and Project reconciliation.
- All current-head checks and resolved first-round findings precede any second-round request.

### Parallel Opportunities

- T012 through T015 own distinct documentation surfaces after verifier behavior stabilizes.
- T016 through T018 create independent canonical pages.
- T021 and T022 own separate release-candidate records.
- Read-only issue evidence, document truth, and link-pattern research can proceed in parallel; the coordinating agent owns all writes and final verification.

## Parallel Example: User Story 2

```text
Task: "Create the SubRip capability and planned-coverage page in docs/formats/srt.md"
Task: "Create the WebVTT capability and planned-coverage page in docs/formats/webvtt.md"
Task: "Create the protected release lifecycle in docs/release-process.md"
```

## Implementation Strategy

### MVP First

1. Establish the verifier module and required-document inventory.
2. Write failing link and heading-resolution tests.
3. Implement offline extraction and repository-confined resolution.
4. Audit and correct current-state documentation.
5. Validate User Story 1 independently before adding missing navigation surfaces.

### Incremental Delivery

1. Deterministic documentation inventory and local-link safety.
2. Current implementation and delivery truth reconciliation.
3. Dedicated SRT, WebVTT, and release-process documentation.
4. Changelog, working-spec, issue, Project, and milestone evidence.
5. Full local/hosted verification, bounded review, and operator handoff.

## Notes

- `[P]` tasks modify separate files and have no unmet dependency.
- `scripts/docs-verify` is read-only, performs no network requests, and stays outside the shipped product module.
- Historical Spec Kit artifacts remain time-bound evidence and are not part of the maintained documentation-link audit.
- S010 never creates a tag, release, immutable release-schema copy, public schema, milestone closure, or production-domain state.

## Phase 7: Convergence

- [x] T032 Add direct verifier tests for required-path non-regular and symlink rejection, directory targets, control characters, protocol-relative URLs, shortcut references, and CLI violation/read-failure exit codes; implement any exposed shortcut-reference gap in `scripts/docs-verify/verify_test.go` and `scripts/docs-verify/verify.go` per T003, T006, T007, and the verification contract (partial)
