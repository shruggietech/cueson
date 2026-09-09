# Tasks: Pull-Request Policy Automation

**Input**: Design documents from `specs/S007-enforce-pr-review-automation/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Tests are required by FR-019 and precede their corresponding implementation.

**Organization**: Tasks are grouped by independently testable user story. Shared GitHub transport and workflow wiring follow the pure policy stories.

## Phase 1: Setup

**Purpose**: Establish the isolated automation module and fixture vocabulary.

- [x] T001 Run and record the blocking pre-implementation Spec Kit analysis across specs/S007-enforce-pr-review-automation/spec.md, plan.md, and tasks.md
- [x] T002 Create the standalone Go 1.25 module and command/package skeletons in scripts/pr-policy/go.mod, scripts/pr-policy/main.go, scripts/pr-policy/policy.go, and scripts/pr-policy/github.go
- [x] T003 [P] Create shared snapshot and expected-result fixture types in scripts/pr-policy/testdata/schema.json and scripts/pr-policy/policy_test.go
- [x] T004 Add scripts/pr-policy/go.mod to the CI dependency-cache inputs in .github/workflows/ci.yml

---

## Phase 2: Foundational Policy and Transport

**Purpose**: Define deterministic types, narrow identity rules, complete collection behavior, and fail-closed transport before story behavior.

**Critical**: User-story work starts only after this phase passes.

- [x] T005 Write failing identity, summary-parsing, pagination, partial-error, and head-race tests in scripts/pr-policy/policy_test.go and scripts/pr-policy/github_test.go
- [x] T006 Implement snapshot, identity, summary, thread, exception, status, and mutation-decision types in scripts/pr-policy/policy.go
- [x] T007 Implement the standard-library GitHub client, REST Link pagination, GraphQL cursor pagination, partial-error rejection, and bounded head-stability retry in scripts/pr-policy/github.go
- [x] T008 Implement exact Codex, Dependabot, operator, summary-marker, round-two-marker, and directive parsing in scripts/pr-policy/policy.go
- [x] T009 Implement JSON result output, flag validation, token handling, dry-run behavior, one-PR mode, and all-open sweep mode in scripts/pr-policy/main.go
- [x] T010 Run gofmt, `go -C scripts/pr-policy test -count=1 ./...`, and `go -C scripts/pr-policy vet ./...` for the foundational module

**Checkpoint**: Complete snapshots and deterministic evidence primitives are independently testable.

---

## Phase 3: User Story 1 - Require Traceable Pull Requests (Priority: P1)

**Goal**: Publish a current-head issue-link result based on GitHub-resolved closing issues or an explicit authorized exception.

**Independent Test**: Replay valid local and cross-repository targets, no target, operator exception, non-operator exception, and Dependabot fixtures.

### Tests for User Story 1

- [x] T011 [US1] Write failing issue-link decision tests and fixtures in scripts/pr-policy/policy_test.go and scripts/pr-policy/testdata/links/
- [x] T012 [US1] Write failing closingIssuesReferences and status read-back transport tests in scripts/pr-policy/github_test.go

### Implementation for User Story 1

- [x] T013 [US1] Implement GitHub-resolved closing-reference, operator exception, Dependabot exception, and failure decisions in scripts/pr-policy/policy.go
- [x] T014 [US1] Implement current-head `Cueson PR policy / Issue link` no-op-aware status publication and read-back in scripts/pr-policy/github.go
- [x] T015 [US1] Verify every User Story 1 fixture and transport test passes in scripts/pr-policy/

**Checkpoint**: Issue-link policy is independently complete.

---

## Phase 4: User Story 2 - Reconcile the Automatic First Review (Priority: P1)

**Goal**: Reconcile draft, pending, acknowledged, clean, stale, finding-bearing, and failed native first-round evidence without issuing a duplicate request.

**Independent Test**: Replay empirical and synthetic first-round snapshots for current and stale heads with reactions, summaries, reviews, and threads.

### Tests for User Story 2

- [x] T016 [US2] Write failing first-round state and current-head attribution fixtures in scripts/pr-policy/policy_test.go and scripts/pr-policy/testdata/reviews/
- [x] T017 [US2] Write failing reaction, review, review-thread, and Codex identity collection tests in scripts/pr-policy/github_test.go

### Implementation for User Story 2

- [x] T018 [US2] Implement draft, pending, eyes, current-head completion, stale-head, failure, and unresolved-finding decisions in scripts/pr-policy/policy.go
- [x] T019 [US2] Implement complete reaction, review, and review-thread collection with exact Codex attribution in scripts/pr-policy/github.go
- [x] T020 [US2] Implement current-head `Cueson PR policy / Codex review` no-op-aware status publication and read-back in scripts/pr-policy/github.go
- [x] T021 [US2] Verify User Story 2 emits zero first-round review comments across all fixtures in scripts/pr-policy/

**Checkpoint**: Native round one is independently reconciled without duplication.

---

## Phase 5: User Story 3 - Bound the Second Review (Priority: P1)

**Goal**: Reserve, publish, and verify exactly one second review after resolved first-round findings, with no automatic third review under any state.

**Independent Test**: Replay unresolved, resolved, duplicate, concurrent, deleted-marker, ambiguous-post, clean-second, finding-second, failed-second, stale-head, and protocol-violation states ten times each.

### Tests for User Story 3

- [x] T022 [US3] Write failing round-attribution, reservation, tenfold replay, and no-third fixtures in scripts/pr-policy/policy_test.go and scripts/pr-policy/testdata/reviews/
- [x] T023 [US3] Write failing reservation status, comment post, ambiguous response, complete refetch, and mutation read-back tests in scripts/pr-policy/github_test.go

### Implementation for User Story 3

- [x] T024 [US3] Implement first-findings, resolved-first, second-ready, second-pending, second-clean, second-blocked, and protocol-violation decisions in scripts/pr-policy/policy.go
- [x] T025 [US3] Implement reviewed-head ancestry comparison, current-head CI validation, status-history reservation, pre-mutation refetch, and consumed-round detection in scripts/pr-policy/github.go
- [x] T026 [US3] Implement the publication-safe marked second-round comment, exact-once POST, full comment read-back, and fail-closed ambiguity behavior in scripts/pr-policy/github.go
- [x] T027 [US3] Prove tenfold replay, authenticated-reservation interleaving, and globally serialized workflow delivery create at most one marked second request and no third request in scripts/pr-policy/

**Checkpoint**: The shared two-round ceiling is independently proven.

---

## Phase 6: User Story 4 - Respect Exclusions and Operator Overrides (Priority: P2)

**Goal**: Make Dependabot exclusions and exact operator waivers observable without allowing association-based or self-authored bypass.

**Independent Test**: Replay configured operator, other owner/member/collaborator, PR author, missing reason, malformed directive, and Dependabot identities.

### Tests for User Story 4

- [x] T028 [US4] Write failing operator-only waiver and spoof-resistance fixtures in scripts/pr-policy/policy_test.go and scripts/pr-policy/testdata/overrides/

### Implementation for User Story 4

- [x] T029 [US4] Implement exact operator-only issue-link and Codex-review waivers with reason evidence in scripts/pr-policy/policy.go
- [x] T030 [US4] Implement explicit Dependabot success descriptions for both policy contexts in scripts/pr-policy/policy.go
- [x] T031 [US4] Verify every non-operator and body-text waiver attempt remains ineffective in scripts/pr-policy/

**Checkpoint**: Exclusions and override authority are independently complete.

---

## Phase 7: Trusted Workflow Integration

**Purpose**: Activate the tested adapter only through trusted default-branch code and bounded recovery.

- [x] T032 Write the globally serialized pull_request_target, issue_comment, and scheduled workflow in .github/workflows/pr-policy.yml
- [x] T033 Configure exact main checkout, persisted-credential disablement, no cache, minimal job permissions, event routing, and no PR-ref execution in .github/workflows/pr-policy.yml
- [x] T034 Add standalone automation tests to the CI Repository text job and Go formatting coverage in .github/workflows/ci.yml
- [x] T035 Add workflow-security assertions for trigger, permission, and fetched-base-ref guard drift in scripts/pr-policy/github_test.go
- [x] T036 Run actionlint and the complete standalone module suite against .github/workflows/pr-policy.yml

---

## Phase 8: Documentation, Analysis, and Controlled Proof

**Purpose**: Align governing documents, pass the blocking analysis gate, and gather local and GitHub evidence without crossing the merge boundary.

- [x] T037 Document S007 policy states, security boundaries, recovery, and post-merge activation limits in docs/architecture.md and docs/project-management.md
- [x] T038 Record the repository-automation architecture decision and user-visible policy gates in CHANGELOG.md
- [x] T039 Revalidate spec quality and run post-implementation Spec Kit convergence and cross-artifact analysis across spec.md, plan.md, and tasks.md
- [x] T040 Run formatter, standalone tests and vet, root vet and tests, race tests, actionlint, and whitespace checks from specs/S007-enforce-pr-review-automation/quickstart.md
- [ ] T041 Commit and push the complete S007 tree, publish a formatted pull request that closes #9, and read its body back from GitHub
- [ ] T042 Run the S007 adapter in dry-run mode against the live pull request and compare its JSON with GitHub evidence
- [x] T043 Record the trusted-default-branch activation limit and require hosted mutation-source proof on the first post-merge pull request before issue #10 configures required contexts

### External delivery gates

These gates are recorded in pull-request evidence rather than checked off in the branch after review, because changing the reviewed head merely to update this file would invalidate otherwise clean review evidence.

- T044 Resolve every first-round review finding, rerun verification, and request at most one second round only when repository protocol permits.
- T045 Confirm CI, CodeQL, security scanning, all review threads, review-round evidence, and Project `Stage = PR review` are ready for the human merge ritual.

---

## Dependencies and Execution Order

- Phase 1 starts immediately; Phase 2 depends on its module skeleton.
- User Story 1 depends on complete snapshot and status primitives from Phase 2.
- User Story 2 depends on identity, summary, and review-evidence primitives from Phase 2 but is otherwise independently testable from User Story 1.
- User Story 3 depends on User Story 2's round-one classification and current-head status transport.
- User Story 4 depends only on Phase 2 identity and directive parsing and can proceed in parallel with User Stories 1 and 2.
- Workflow integration depends on all four policy stories.
- Documentation, publication, live proof, and review closure follow implementation and the blocking analysis gate in chronological order.

## Parallel Opportunities

- T003 can proceed alongside the command skeleton after T001.
- T011 and T016 cover separate fixture directories and may proceed in parallel after Phase 2.
- T028 can proceed alongside User Stories 1 and 2 because it owns a separate fixture directory.
- T037 and T038 affect different documents and may proceed in parallel after workflow behavior stabilizes.
- Publication and live GitHub evidence tasks remain sequential because each depends on the exact current head and prior review state.

## Implementation Strategy

1. Establish the pure snapshot and adapter boundary.
2. Deliver issue-link and round-one policy as the minimum independently useful gates.
3. Add the exact-once second-round transition and permanent ceiling.
4. Add operator-only exceptions and Dependabot behavior.
5. Integrate only through trusted default-branch events.
6. Converge documentation and verification before publication.
7. Use the live S007 pull request for current-head API and status read-back evidence, then stop at human merge.
