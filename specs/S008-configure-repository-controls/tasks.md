# Tasks: Verified Repository Controls

**Input**: Design documents from `specs/S008-configure-repository-controls/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: S008 changes delivery configuration rather than shipped product behavior. Every external mutation has an immediate live read-back task, and existing local plus hosted regression gates remain mandatory.

**Organization**: Tasks are grouped by user story. The official pull request is a foundational dependency because hosted check identity cannot be proven before publication.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because the task owns a different file or independent read surface
- **[Story]**: Maps the task to the corresponding specification user story
- Every task names the file that records its durable result

## Phase 1: Setup and Baseline

**Purpose**: Establish the active slice, immutable baseline evidence, and source-controlled record before configuration changes.

- [x] T001 Verify issue #10 dependencies, parent, milestone, Slice `S008`, and Project Stage `Specced`, recording planning state in specs/S008-configure-repository-controls/plan.md
- [x] T002 Capture pre-mutation repository, merge, Actions, security, CodeQL, classic-protection, repository-ruleset, and effective organization-ruleset state in docs/repository-controls.md
- [x] T003 Prove the organization-owned ruleset identity and baseline update timestamp in docs/repository-controls.md
- [x] T004 Validate all 25 functional requirements and 8 success criteria against specs/S008-configure-repository-controls/checklists/requirements.md

---

## Phase 2: Foundational Documentation and Publication

**Purpose**: Create the durable desired-state and verification record, validate it locally, and publish the official pull request required for hosted evidence.

**CRITICAL**: Hosted configuration work cannot start until the official pull request exists and targets `main`.

- [x] T005 [P] Document effective GitHub rules, desired repository settings, security controls, mutation order, and truthful limitation handling in docs/repository-controls.md
- [x] T006 [P] Add repository-control and required-check architecture boundaries to docs/architecture.md
- [x] T007 [P] Add ruleset, check-evidence, recovery-bypass, and post-merge branch-deletion workflow to docs/project-management.md
- [x] T008 [P] Make private-reporting availability language conditional on verified live configuration in SECURITY.md
- [x] T009 [P] Record S008 repository-control additions and architecture decisions under `[Unreleased]` in CHANGELOG.md
- [x] T010 Verify Spec Kit artifacts, repository UTF-8 without BOM, line endings, mojibake absence, publication formatting, whitespace, root tests, race tests, vet, and nested script tests using specs/S008-configure-repository-controls/quickstart.md
- [x] T011 Run the blocking cross-artifact analysis gate over spec.md, plan.md, and tasks.md and resolve every critical or high-severity finding in specs/S008-configure-repository-controls/
- [ ] T012 Commit the validated pre-publication slice, push S008-configure-repository-controls, publish the official pull request with `Closes #10`, and read its formatted body back into the verification record in docs/repository-controls.md
- [ ] T013 Move issue #10 to Project Stage `In progress`, preserve Slice `S008`, and clear the default Status field, recording read-back in docs/repository-controls.md

**Checkpoint**: The official S008 pull request exists, is publication-safe, and can supply hosted evidence without granting administrative credentials to pull-request code.

---

## Phase 3: User Story 2 - Establish Secure Repository Defaults (Priority: P1)

**Goal**: Reduce workflow and merge privilege, enable supported security facilities, and preserve versioned CodeQL.

**Independent Test**: Every repository setting and security capability has an authoritative before value, explicit desired value, mutation result, and immediate after value in docs/repository-controls.md.

- [ ] T014 [US2] Read the initial S008 hosted CI, CodeQL, and policy results before permission changes and record their current head in docs/repository-controls.md
- [ ] T015 [US2] Restrict repository Actions to GitHub-owned actions and full-SHA references, immediately read the policy back, and record it in docs/repository-controls.md
- [ ] T016 [US2] Reduce default workflow permissions to read with pull-request approval disabled, immediately read the defaults back, and record them in docs/repository-controls.md
- [ ] T017 [US2] Enable Dependabot security updates with dependency alerts retained, immediately read both states back, and record them in docs/repository-controls.md
- [ ] T018 [US2] Enable secret scanning and push protection when supported, immediately read both states back, and record any limitation in docs/repository-controls.md
- [ ] T019 [US2] Enable private vulnerability reporting when supported, immediately read its state back, and record any limitation in docs/repository-controls.md
- [ ] T020 [US2] Preserve the versioned CodeQL workflow and read back that default setup remains unconfigured with current CodeQL checks healthy in docs/repository-controls.md
- [ ] T021 [US2] Make squash the only merge method, keep auto-merge disabled, enable automatic merged-head deletion, immediately read repository settings back, and record them in docs/repository-controls.md
- [ ] T022 [US2] Dispatch a new S008 head update under the reduced defaults and require all current-head hosted workflows to finish successfully, recording results in docs/repository-controls.md

**Checkpoint**: Secure repository defaults are verified independently, and existing workflows still operate with their explicit minimum permissions.

---

## Phase 4: User Story 1 - Protect the Default Branch with Proven Gates (Priority: P1)

**Goal**: Create repository-scoped default-branch protection using only successful current-head evidence.

**Independent Test**: The complete repository-owned ruleset reads back with the expected target, bypass, pull-request rules, ref protections, and exact proven check list.

- [ ] T023 [US1] Query the current S008 head's check runs and map all 17 S006 contract contexts to successful GitHub Actions evidence in docs/repository-controls.md
- [ ] T024 [US1] Query current-head combined statuses and record exact source evidence for both S007 policy contexts in docs/repository-controls.md
- [ ] T025 [US1] Evaluate the S007 Actions-bot second-round proof gate and either admit both policy contexts or defer both with an explicit reason in docs/repository-controls.md
- [ ] T026 [US1] Construct the strict required-check list from only admitted hosted evidence and record exact contexts plus integration identifiers in docs/repository-controls.md
- [ ] T027 [US1] Create the active repository-owned `cueson verified default branch` ruleset with one OrganizationAdmin recovery bypass and record its returned identifier in docs/repository-controls.md
- [ ] T028 [US1] Read the complete repository ruleset back and verify default-branch targeting, deletion and non-fast-forward protection, pull-request requirement, resolved conversations, squash-only policy, strict checks, and the sole bypass in docs/repository-controls.md
- [ ] T029 [US1] Read effective branch rules and confirm the repository and organization rules combine without changing organization ruleset `20478126`, recording the comparison in docs/repository-controls.md

**Checkpoint**: The default branch is protected by an active repository-owned ruleset whose check list is fully traceable to current-head evidence.

---

## Phase 5: User Story 3 - Preserve Scope and Recovery Boundaries (Priority: P1)

**Goal**: Prove that hardening remained repository-scoped and recoverable without consuming merge or release authority.

**Independent Test**: Live API state shows one S008-created repository ruleset, one organization-administrator recovery category, no organization-rule update, and no release, tag, schema, domain, or unrelated-repository mutation.

- [ ] T030 [US3] Re-query organization ruleset `20478126` and compare its identifier, source, rules, bypass actors, and update timestamp with the baseline in docs/repository-controls.md
- [ ] T031 [US3] Verify the S008-created ruleset has no bot, deploy-key, contributor-role, unrelated-team, or second bypass entry in docs/repository-controls.md
- [ ] T032 [US3] Verify no auto-merge, merge, tag, release, schema publication, production-domain, organization-policy, or unrelated-repository action occurred and record the boundary audit in docs/repository-controls.md
- [ ] T033 [US3] Exercise read-only effective-rule and mergeability inspection for the S008 pull request without using the administrative bypass, recording the result in docs/repository-controls.md

**Checkpoint**: Scope and recovery boundaries are explicit and independently auditable.

---

## Phase 6: User Story 4 - Produce Auditable Configuration Evidence (Priority: P2)

**Goal**: Finish a chronological evidence record that a maintainer can reproduce and review.

**Independent Test**: The quickstart, source-controlled evidence, issue body, Project fields, pull-request comments, final head, ruleset, and hosted results agree without placeholders or unsupported claims.

- [ ] T034 [US4] Reconcile docs/repository-controls.md, docs/architecture.md, docs/project-management.md, SECURITY.md, CHANGELOG.md, and all S008 artifacts with the final observed external state
- [ ] T035 [US4] Rerun the complete local verification suite and final Spec Kit analysis plus convergence gates, appending and completing any remediation tasks in specs/S008-configure-repository-controls/tasks.md
- [ ] T036 [US4] Commit and push the final evidence update, then wait for all current-head CI, CodeQL, Codex, security-bot, and policy results without starting more than one authorized second Codex round
- [ ] T037 [US4] Address every review finding with regression or evidence changes, reply with commit evidence, resolve only satisfied threads, and verify no unresolved review remains on the S008 pull request
- [ ] T038 [US4] Update issue #10 acceptance and verification evidence through the publication formatter, immediately read the body back, and preserve its open state until merge
- [ ] T039 [US4] Move issue #10 to Project Stage `PR review`, preserve Slice `S008`, clear default Status, and read every field back into docs/repository-controls.md
- [ ] T040 [US4] Publish and read back one formatted final-head evidence comment containing required and deferred contexts, successful run URLs, policy status sources, repository ruleset identifier, and final read-back time
- [ ] T041 [US4] Perform the final clean-worktree, branch, remote, issue, milestone, Project, ruleset, security, workflow, review, and CI audit, then stop for the operator's final review and merge ritual

**Checkpoint**: S008 is complete and verified but unmerged, with the human operator holding the final decision.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1**: Starts from clean synchronized `main` and issue #10's unblocked state.
- **Phase 2**: Depends on Phase 1 and blocks every live mutation because the official pull request supplies hosted identity evidence.
- **Phase 3 / US2**: Depends on the official pull request and verifies tightened defaults before required checks are activated.
- **Phase 4 / US1**: Depends on Phase 3's successful workflow rerun so required checks are proven under final permission defaults.
- **Phase 5 / US3**: Depends on the repository-owned ruleset created in Phase 4.
- **Phase 6 / US4**: Depends on all controls and scope checks; it owns final review remediation and evidence publication.

### User Story Dependencies

- **User Story 2 (P1)**: Begins after publication and can be verified from settings read-back without branch protection.
- **User Story 1 (P1)**: Depends on User Story 2 only to prove workflows under final least-privilege defaults before making their contexts required.
- **User Story 3 (P1)**: Depends on the User Story 1 ruleset and independently proves scope plus recovery boundaries.
- **User Story 4 (P2)**: Consolidates the independently verified outcomes after all P1 stories.

### Parallel Opportunities

- T005 through T009 edit independent documentation files and can proceed in parallel.
- T017 through T020 touch independent security API surfaces but remain sequential in execution so each read-back is isolated and partial failure is obvious.
- T023 and T024 read independent hosted evidence surfaces but feed the sequential admission decision in T025.
- T030 through T032 are independent read-only boundary audits after ruleset creation.

---

## Parallel Example: Foundational Documentation

```text
Task: "Document desired controls and evidence in docs/repository-controls.md"
Task: "Update protection architecture in docs/architecture.md"
Task: "Update delivery process in docs/project-management.md"
Task: "Correct private-reporting language in SECURITY.md"
Task: "Record the S008 decision in CHANGELOG.md"
```

---

## Implementation Strategy

### MVP First

1. Complete baseline and source-controlled contracts.
2. Publish the official S008 pull request.
3. Reduce repository and workflow privilege with immediate read-back.
4. Prove stable hosted checks on the tightened configuration.
5. Add repository-scoped default-branch protection using only that evidence.

### Incremental Delivery

1. Source artifacts establish intent without external mutation.
2. Secure defaults can be verified independently of ruleset creation.
3. Required checks are admitted one evidence-backed context at a time.
4. The final evidence pass reconciles the source diff with external GitHub state.

### Autopilot Review Boundary

The native Codex integration owns round one. If round one has findings, S008 resolves every finding and permits the trusted workflow or configured operator to request exactly one second review. A clean first round produces no second request. Second-round findings are resolved without a third automated request. Final merge always remains with the operator.

## Notes

- GitHub-bound Markdown is formatted before publication and read back immediately.
- External state is never inferred from a successful mutation response alone.
- Head-specific evidence belongs in the issue or pull-request record because a commit cannot durably name its own final hash.
- Existing organization policy is evidence input, never an S008 mutation target.
