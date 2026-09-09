# Feature Specification: Pull-Request Policy Automation

**Feature Branch**: `S007-enforce-pr-review-automation`

**Created**: 2026-09-09

**Status**: Draft

**Input**: User description: "Implement issue-linked pull-request enforcement and bounded Codex review automation for work slice S007 under the repository autopilot protocol."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Require Traceable Pull Requests (Priority: P1)

A maintainer opening or updating a normal pull request receives an unambiguous policy result that accepts at least one complete closing reference and rejects an incomplete or absent reference unless a trusted, documented exception applies.

**Why this priority**: Closing references are the link between implementation evidence and the issue acceptance criteria they claim to satisfy.

**Independent Test**: Evaluate GitHub-resolved closing-reference lists for valid local and cross-repository references, malformed placeholders, ordinary prose references, operator exception directives, untrusted exception directives, and Dependabot authorship.

**Acceptance Scenarios**:

1. **Given** a normal pull request with `Closes #123`, **When** policy is evaluated, **Then** the issue-link result passes and identifies the closing reference.
2. **Given** a normal pull request without a complete closing reference, **When** policy is evaluated, **Then** the issue-link result fails with corrective guidance.
3. **Given** the configured human operator records a documented issue-link exception, **When** policy is reevaluated, **Then** the result passes while preserving the exception rationale as evidence.
4. **Given** an untrusted participant supplies the same exception text, **When** policy is evaluated, **Then** the exception is ignored and normal validation remains in force.
5. **Given** a Dependabot pull request has no closing reference, **When** policy is evaluated, **Then** the issue-link result passes as an explicit automated-source exception.

---

### User Story 2 - Reconcile the Automatic First Review (Priority: P1)

An eligible pull request becomes reviewable and the repository records a single first-round Codex review state without duplicating the native Codex integration's automatic request.

**Why this priority**: Duplicate requests waste review capacity and make the two-round limit impossible to audit.

**Independent Test**: Replay reviewable, draft, acknowledged, clean, finding-bearing, and failed first-round states and verify one deterministic policy result for each state without emitting a second first-round request.

**Acceptance Scenarios**:

1. **Given** an eligible pull request has just become reviewable, **When** review state is reconciled, **Then** the policy records that round one is pending and relies on the configured native trigger.
2. **Given** Codex has acknowledged round one with an eyes reaction, **When** review state is reconciled, **Then** acknowledgement remains pending rather than being treated as approval.
3. **Given** Codex completes round one with a thumbs-up and no findings, **When** review state is reconciled, **Then** the policy passes and emits no second-round request.
4. **Given** Codex reports a failed first-round attempt without a review result, **When** review state is reconciled, **Then** the policy reports a terminal review failure requiring maintainer judgment rather than silently passing or consuming another round.
5. **Given** Codex completes through a reaction without another supported event, **When** scheduled or manual recovery runs, **Then** the current-head policy state is reconciled without requesting another review.
6. **Given** the pull-request head changes after a clean round one, **When** state is reconciled, **Then** the earlier review becomes stale and the policy remains blocking until the human operator explicitly requests or waives another review.
7. **Given** round-one findings were reviewed on head A and a remediation produces descendant head B, **When** all A findings are resolved and B's required CI succeeds, **Then** the policy may classify B as ready for the one second review.

---

### User Story 3 - Bound the Second Review (Priority: P1)

After every first-round finding has been addressed and its thread resolved, the repository may request exactly one second Codex review and can never request a third automatically.

**Why this priority**: The second round proves remediations while the hard upper bound prevents event loops and uncontrolled external work.

**Independent Test**: Replay unresolved findings, resolved findings, duplicate delivery events, an existing marked second-round request, clean second-round completion, new second-round findings, and failed second-round completion.

**Acceptance Scenarios**:

1. **Given** at least one first-round Codex thread remains unresolved, **When** state is reconciled, **Then** the policy fails and no second review is requested.
2. **Given** first-round findings exist and every associated thread is resolved, **When** state is reconciled for the first time, **Then** exactly one marked `@codex review` request is emitted.
3. **Given** the same resolved state is delivered repeatedly, **When** each event is reconciled, **Then** the existing marked request is recognized and no duplicate is emitted.
4. **Given** round two completes cleanly for the current head, **When** state is reconciled, **Then** the policy passes and emits no further request.
5. **Given** round two fails, produces unresolved findings, or becomes stale after a head change, **When** state is reconciled, **Then** the policy remains blocking for human judgment and never emits a third automated request.

---

### User Story 4 - Respect Exclusions and Trusted Overrides (Priority: P2)

Dependabot pull requests and pull requests carrying a pull-request-specific waiver from the configured human operator do not receive automatic Codex requests, while the resulting exemption remains visible and auditable.

**Why this priority**: Automation must honor the repository's bot exclusion and human override authority without allowing pull-request authors to bypass review themselves.

**Independent Test**: Evaluate Dependabot identities, configured-operator waiver comments with reasons, missing reasons, comments from every other identity, self-waiver attempts by pull-request authors, and ordinary eligible authors.

**Acceptance Scenarios**:

1. **Given** a Dependabot-authored pull request, **When** review state is reconciled, **Then** the policy records the default exclusion and emits no Codex request.
2. **Given** the configured human operator records a pull-request-specific review waiver with a reason, **When** state is reconciled, **Then** the policy passes with that waiver as evidence.
3. **Given** any other identity records a waiver, a pull-request author attempts self-waiver, or a waiver lacks a reason, **When** state is reconciled, **Then** the waiver is ignored.

### Edge Cases

- Pull-request edits, synchronizations, reopened events, review-comment changes, and duplicate event deliveries can arrive in any order.
- An eyes acknowledgement without a later terminal result must remain pending.
- A thumbs-up does not establish a clean review when actionable Codex findings remain unresolved.
- A stale thumbs-up or failure from an earlier round must not complete a later round.
- A changed head invalidates earlier clean or second-round terminal results. The sole exception is the explicit remediation bridge from a finding-bearing reviewed ancestor to its descendant current head with all findings resolved and current-head CI green.
- Deleted or edited comments must cause the next reconciliation to derive state from current GitHub evidence rather than cached event assumptions.
- Fork-originated and untrusted pull requests must not execute repository-controlled pull-request code with a write-capable credential.
- API failures and incomplete pagination must fail closed without posting a review request.
- Pull requests may resolve multiple valid closing references or none; policy consumes GitHub's current parsed result rather than maintaining a competing Markdown grammar.
- A top-level Codex finding with no resolvable thread remains blocking until the human operator records a waiver or the integration supplies unambiguous resolution evidence.
- A native review that remains pending for more than the configured recovery interval remains pending; the schedule preserves liveness but never invents terminal evidence.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The policy MUST fail a normal pull request that lacks at least one syntactically complete closing reference.
- **FR-002**: The policy MUST use GitHub's current parsed closing-issue references as the authority for supported local and cross-repository closing syntax.
- **FR-003**: The policy MUST reject placeholders, bare mentions, examples, inaccessible targets, and other prose whenever GitHub reports no complete closing-issue reference.
- **FR-004**: A documented issue-link exception MUST require a pull-request-specific reason from an exactly configured human-operator login and MUST remain observable in the policy result.
- **FR-005**: The review policy MUST exclude Dependabot by default and MUST support a pull-request-specific waiver only from an exactly configured human-operator login without treating association, authorship, or body text as authority.
- **FR-006**: An eligible reviewable pull request MUST receive exactly one automatic first-round trigger or reconciliation action, with the configured native integration treated as the first-round trigger authority.
- **FR-007**: An eyes reaction MUST be interpreted only as acknowledgement or work in progress.
- **FR-008**: A thumbs-up MUST be interpreted as a clean review only when it is attributable to the configured Codex integration, belongs to the current review round, and no actionable Codex finding remains unresolved.
- **FR-009**: A failed Codex attempt without a completed review result MUST fail closed and require maintainer judgment.
- **FR-010**: The policy MUST request round two only when first-round findings exist, every associated finding thread is resolved, and no earlier second-round request exists.
- **FR-022**: When remediation changes the head after first-round findings, the reviewed finding head MUST be an ancestor of the current head, the current head MUST be newer, and every configured current-head CI gate MUST succeed before round two is requested; force-pushed or ambiguous ancestry MUST fail closed.
- **FR-011**: An automated round-two request MUST contain exactly one `@codex review` invocation and a stable machine-readable marker attributable to the exact GitHub Actions bot identity; a matching marker from any other actor MUST fail closed without establishing a round boundary.
- **FR-012**: Repeated or concurrent event delivery MUST NOT create duplicate second-round requests.
- **FR-013**: Once a second-round request exists, the policy MUST NOT generate a third automated review request under any later state; an operator-issued request MUST cite exactly one backticked commit prefix that uniquely resolves to a collected pull-request head before its result can be attributed, MUST use the comment's latest edit time as its round boundary, and MUST create an authenticated durable reservation when observed so editing or deleting the comment cannot restore the allowance, while an unbound request still consumes the allowance and fails closed.
- **FR-014**: Second-round failure, findings, or a later head change MUST remain blocking until current-head clean evidence or an explicit human-operator waiver exists; automation MUST NOT request another review.
- **FR-015**: Policy results MUST apply to the pull request's current head revision and expose stable, independently selectable issue-link and Codex-review contexts.
- **FR-016**: The automation MUST reconcile relevant pull-request, comment, and review changes, run scheduled recovery often enough to observe reaction-only completion, and provide a safe manual reconciliation path.
- **FR-017**: The automation MUST operate against trusted default-branch code, MUST NOT check out or execute pull-request-controlled content with write authority, and MUST use only the minimum repository permissions needed to read evidence, publish policy state, and emit the permitted review request.
- **FR-023**: Every event path that cannot filter by base branch before invocation MUST fetch the pull request and refuse mutation unless its current base branch is `main`.
- **FR-018**: API errors, unsupported evidence shapes, missing pages, and ambiguous round attribution MUST fail closed with actionable diagnostics and no review-request mutation.
- **FR-019**: Tests MUST cover valid, invalid, exception, clean, findings, resolved, duplicate-event, stale-evidence, Dependabot, failure, and documented-override scenarios.
- **FR-020**: S007 MUST NOT merge pull requests, enable auto-merge, configure required checks or repository rules, alter release state, or claim product-format capability.
- **FR-021**: Every status or review-request mutation MUST use publication-safe UTF-8 content, be preceded by a current-evidence duplicate check, and be followed by immediate GitHub read-back verification.

### Key Entities

- **Pull Request Policy Snapshot**: Current head revision, author identity and trust, draft state, body, comments, reactions, Codex findings, and resolution state used for one evidence-based evaluation.
- **Closing Reference**: A GitHub closing keyword paired with one or more complete local or cross-repository issue identifiers outside example-only content.
- **Trusted Exception**: A pull-request-specific comment from an exactly configured human-operator login naming the skipped policy and containing a non-empty reason.
- **Review Round**: The automatic first review or the single permitted second review, including request time, acknowledgement, terminal result, and findings.
- **Policy Result**: A stable context, state, summary, and current pull-request head revision suitable for maintainers and later repository protection.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The complete validation corpus produces the expected issue-link result for 100% of valid, invalid, and exception fixtures.
- **SC-002**: The complete review-state corpus produces the expected pending, passing, or failing result and mutation decision for 100% of defined scenarios.
- **SC-003**: Replaying any supported event scenario ten times results in at most one marked second-round request and zero third-round requests.
- **SC-004**: Every policy mutation is followed by a read-back or duplicate-prevention check against current evidence before the run reports success.
- **SC-005**: No workflow path executes pull-request-controlled repository content with a write-capable credential.
- **SC-006**: Maintainers can identify the blocking policy, affected head revision, evidence state, and corrective action from one policy result without inspecting workflow source.
- **SC-007**: A controlled live pull-request audit agrees with the locally replayed policy result for closing-reference and review-round evidence.

## Assumptions

- The repository's configured native Codex integration remains the authority that automatically starts round one when an eligible pull request becomes reviewable.
- The Codex GitHub App identity remains `chatgpt-codex-connector[bot]`; identity is a declared automation input so a future rename fails visibly rather than widening trust.
- The human operator login is configured explicitly as `h8rt3rmin8r`; repository association alone never grants waiver authority, and pull-request body text never grants an exception.
- Dependabot is identified by the canonical `dependabot[bot]` login and is excluded from Codex review by default.
- Repository rules and required-check configuration remain owned by issue #10 after S007 proves stable contexts.
- Because a newly added trusted-default-branch workflow cannot govern the pull request that introduces it, S007 limits its own pull request to a read-only live adapter audit plus fixture-backed mutation tests. The first post-merge pull request must verify hosted status and Actions-bot comment behavior before issue #10 makes either context required; no unmerged code receives a write credential.
