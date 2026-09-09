# Feature Specification: Verified Repository Controls

**Feature Branch**: `S008-configure-repository-controls`

**Created**: 2026-09-09

**Status**: Draft

**Input**: User description: "Configure verified repository quality and security controls for work slice S008 under the repository autopilot protocol, publish the official pull request automatically, complete at most two Codex review rounds, and stop for the operator's final merge ritual."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Protect the Default Branch with Proven Gates (Priority: P1)

A maintainer can trust that changes reach the default branch through pull requests whose conversations are resolved and whose required quality gates have already demonstrated stable identities on the exact revision under review.

**Why this priority**: The default branch is the delivery authority. Requiring unproven or ambiguous checks could either admit unverified work or lock maintainers out of legitimate recovery.

**Independent Test**: Compare the effective default-branch rules with successful hosted evidence for the S008 head, then demonstrate that deletion, force pushes, direct unprivileged updates, unresolved conversations, and missing required gates are blocked while the configured operator recovery path remains available.

**Acceptance Scenarios**:

1. **Given** a stable check has succeeded on the current S008 head from its expected source, **When** repository protection is configured, **Then** the exact proven check may be required for future default-branch updates.
2. **Given** a check name or source has not been proven on the current hosted pull request, **When** repository protection is configured, **Then** that check is not made required and the evidence gap is recorded.
3. **Given** an unprivileged actor attempts to delete, force-update, or directly update the default branch, **When** repository rules are evaluated, **Then** the operation is blocked.
4. **Given** a pull request has an unresolved conversation or a missing required check, **When** it is evaluated for merge, **Then** the merge remains blocked.
5. **Given** repository protection is active, **When** the configured human operator needs to recover from a broken gate, **Then** a narrowly identified administrative bypass remains available and auditable.

---

### User Story 2 - Establish Secure Repository Defaults (Priority: P1)

A maintainer can rely on repository-wide merge, workflow, branch-cleanup, dependency, code-scanning, secret-scanning, and private-reporting defaults that minimize unnecessary privilege and expose unavailable features honestly.

**Why this priority**: Secure defaults reduce the number of individual workflows and contributors that must compensate for permissive repository settings.

**Independent Test**: Read each supported repository setting before and after configuration and compare its final value with the approved desired state, including truthful limitation evidence for every unavailable feature.

**Acceptance Scenarios**:

1. **Given** repository workflows do not request broader permissions, **When** the default workflow permission is configured, **Then** it grants read access by default and does not allow workflows to approve pull requests.
2. **Given** contributors merge an approved pull request, **When** merge options are presented, **Then** squash merge is the only enabled merge strategy and merged head branches are deleted automatically when GitHub can delete them.
3. **Given** a supported security feature is available to this public repository, **When** controls are applied, **Then** the feature is enabled and its state is read back.
4. **Given** a requested security feature is unavailable because of plan, organization, licensing, or API constraints, **When** verification completes, **Then** the limitation is recorded with the observed response and no success is claimed.

---

### User Story 3 - Preserve Scope and Recovery Boundaries (Priority: P1)

The human operator retains recovery authority without S008 altering organization-wide policy, release state, or unrelated repositories.

**Why this priority**: Repository hardening must not create an organization-wide blast radius or consume authorities the operator did not grant.

**Independent Test**: Compare organization and repository rule sources before and after S008, inspect all bypass actors and mutation targets, and prove that only `shruggietech/cueson` repository settings changed.

**Acceptance Scenarios**:

1. **Given** an existing organization-owned ruleset applies to the repository, **When** S008 adds Cueson-specific requirements, **Then** the organization ruleset remains byte-for-byte unchanged and a repository-owned ruleset carries the new controls.
2. **Given** the operator authorized S008 publication but not merging or releasing, **When** the slice completes, **Then** no agent merges the pull request, enables auto-merge, creates or moves a tag, publishes a release, or changes production `cueson.io` configuration.
3. **Given** repository settings are applied in stages, **When** a later mutation fails, **Then** the verified earlier state and failed operation are reported without concealing partial progress or widening bypass authority.

---

### User Story 4 - Produce Auditable Configuration Evidence (Priority: P2)

A maintainer reviewing the pull request can reconstruct what changed, why it changed, which hosted evidence justified each required gate, and what remains deferred.

**Why this priority**: GitHub configuration is external state and does not appear completely in a source diff, so durable read-back evidence is required for review and future maintenance.

**Independent Test**: Follow the S008 validation guide from a clean checkout and compare the recorded desired state, mutation sequence, hosted check inventory, final API read-back, and documented exceptions.

**Acceptance Scenarios**:

1. **Given** a repository setting is mutated, **When** the operation succeeds, **Then** its value is immediately read back and recorded against the intended value.
2. **Given** a required-check candidate is evaluated, **When** evidence is recorded, **Then** the record identifies the exact context, current head revision, result, and source needed to justify or reject protection.
3. **Given** the S008 pull request receives a clean first Codex review without findings, **When** policy protection is finalized, **Then** the lack of a GitHub Actions-authored second-round request is recorded and both PR-policy contexts remain non-required.
4. **Given** S008 receives first-round findings and the trusted workflow emits the one permitted second-round request, **When** the request and policy statuses are read back successfully, **Then** the two PR-policy contexts may be required only after the complete hosted proof is recorded.

### Edge Cases

- Required-check names can collide across providers, change capitalization, or appear on a stale revision.
- A check can report success on the pull request merge commit while remaining absent from the pull request head.
- A policy status can be written by an unexpected actor even when its context text matches.
- GitHub commit-status list and create responses can omit a redundant commit SHA even though the request endpoint is bound to that SHA.
- The first post-S007 pull request can complete cleanly without exercising the automated second-round comment path.
- Repository rules can combine cumulatively with an organization-owned ruleset even when classic branch protection reports no configuration.
- An administrative bypass can be too broad, absent, or unavailable for repository-owned rules.
- Enabling secret scanning or push protection can fail because of repository visibility, plan, organization policy, or API restrictions.
- Dependency alerts can be enabled while automated security updates remain paused or unavailable.
- Automatic branch deletion cannot delete a protected branch or a branch for which GitHub lacks permission.
- Reducing default workflow permissions must not remove explicit narrow permissions required by existing CI, CodeQL, or pull-request policy workflows.
- Applying a required check before it succeeds on the active pull-request head can make the branch unmergeable until operator recovery.
- Concurrent operator or organization-policy changes can make the post-mutation read-back differ from the requested state.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: S008 MUST inventory the effective repository, merge, workflow-permission, security, ruleset, branch-protection, and required-check state before any mutation.
- **FR-002**: Every mutable control MUST have one explicit desired value, a mutation result, and an immediate authoritative read-back result.
- **FR-003**: S008 MUST keep squash merge enabled and MUST disable merge-commit and rebase-merge methods so squash is the only normal merge strategy.
- **FR-004**: S008 MUST enable automatic merged-head branch deletion and MUST keep automatic merge disabled.
- **FR-005**: Default workflow permissions MUST be reduced to read access and workflows MUST remain unable to approve pull-request reviews by default.
- **FR-006**: Existing workflows MUST retain only their explicit minimum permissions and MUST continue to complete after the repository default is reduced.
- **FR-007**: Applicable dependency graph, vulnerability alert, automated dependency security update, code-scanning, secret-scanning, secret push-protection, and private vulnerability reporting features MUST be enabled when supported.
- **FR-008**: Any unavailable or externally constrained security feature MUST be recorded with the authoritative limitation response and MUST NOT be represented as enabled.
- **FR-009**: S008 MUST NOT replace, edit, disable, or widen the bypass of the existing organization-owned default-branch ruleset.
- **FR-010**: Cueson-specific protections MUST be expressed through a repository-owned ruleset that targets only the default branch.
- **FR-011**: The repository-owned ruleset MUST block default-branch deletion and non-fast-forward updates.
- **FR-012**: The repository-owned ruleset MUST require pull requests and resolution of every review conversation before merge.
- **FR-013**: The repository-owned ruleset MUST require only check contexts whose stable names and successful current-head executions have been read back from GitHub.
- **FR-014**: Required check evidence MUST identify the current pull-request head revision, exact context, terminal success, and provider identity when GitHub exposes one.
- **FR-015**: A successful check on a stale head, merge-only revision, unrelated branch, or unexpected provider MUST NOT justify a required-check rule.
- **FR-016**: The S006 CI and CodeQL contracts MUST be compared with current hosted results before their contexts are selected for protection.
- **FR-017**: The S008 pull request MUST demonstrate the hosted `Cueson PR policy / Issue link` and `Cueson PR policy / Codex review` statuses before either is considered for protection.
- **FR-018**: Neither PR-policy context MAY become required unless S008 also demonstrates the GitHub Actions-authored second-round request behavior required by the S007 activation contract.
- **FR-019**: If S008 completes without a finding-bearing first review, both PR-policy contexts MUST remain non-required and the deferred proof MUST be documented without blocking other proven controls.
- **FR-020**: The repository-owned ruleset MUST retain a narrowly identified bypass for the configured human operator or repository administration role sufficient for recovery from a broken required gate.
- **FR-021**: S008 MUST inspect the combined effective behavior of organization-owned and repository-owned rules without claiming the classic branch-protection endpoint is the sole authority.
- **FR-022**: S008 MUST preserve a chronological, replayable record of desired controls, evidence gates, attempted mutations, read-back results, limitations, and deferred items.
- **FR-023**: A failed or partial mutation MUST fail visibly, preserve the last verified state in the record, and stop dependent mutations until the failure is understood.
- **FR-024**: The final pull request MUST close issue #10, contain publication-safe Markdown, and stop for human final review and merge after all available reviews are resolved and all required checks are green.
- **FR-025**: S008 MUST NOT merge or auto-merge its pull request, modify organization-wide rules, alter unrelated repositories, create or move a tag, publish a release, publish a schema, or mutate production `cueson.io` configuration.
- **FR-026**: Hosted status mutation verification MUST compare the accepted identifier, context, state, description, target, and trusted creator against the requested commit endpoint without requiring a redundant field that GitHub omits from its status representation.

### Key Entities

- **Repository Control**: One repository setting or feature with a scope, current value, desired value, mutation capability, and verification result.
- **Hosted Check Evidence**: An exact check or status context observed on a particular pull-request head, including state, source, and completion time.
- **Repository Ruleset**: A repository-owned default-branch policy containing targeting conditions, required behaviors, required checks, enforcement state, and recovery bypass.
- **Security Capability**: A requested GitHub security feature with availability, observed state, desired state, mutation result, and any limitation evidence.
- **Control Evidence Record**: The chronological source-controlled account that maps each requested control to its before state, action, after state, and supporting hosted evidence.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of in-scope mutable settings have a recorded before value, desired value, mutation result, and authoritative after value.
- **SC-002**: 100% of required check contexts match a successful result on the final S008 pull-request head and include provider identity whenever GitHub exposes it.
- **SC-003**: Zero PR-policy contexts are required without the complete hosted activation evidence mandated by S007.
- **SC-004**: The effective default-branch policy blocks deletion, non-fast-forward updates, direct unprivileged updates, unresolved conversations, and every selected missing or failing required check.
- **SC-005**: Exactly one recovery bypass path remains for the configured operator authority, with no S008-created bypass for contributors, automation, or unrelated teams.
- **SC-006**: 100% of requested security capabilities are either read back as enabled or recorded with a specific authoritative limitation.
- **SC-007**: Existing CI, CodeQL, and pull-request policy workflows complete under the reduced default workflow permission without acquiring broader explicit permission.
- **SC-008**: The final repository evidence contains zero mutations to organization-owned rules, unrelated repositories, release state, tags, schemas, or production domain configuration.
- **SC-009**: The first mutation of each PR-policy context reads back successfully in the same run when GitHub returns its documented commit-status representation.

## Assumptions

- The existing organization-owned `default-branch PR gate` remains an independent baseline and is outside S008 mutation authority.
- The configured human operator is `h8rt3rmin8r`, and organization or repository administration is an acceptable recovery identity only when it does not grant a new bypass to ordinary contributors or automation.
- The exact required-check set will be selected from successful hosted S006 and S007 evidence rather than from names copied only from documentation.
- CodeQL's versioned workflow remains the code-scanning implementation; S008 does not enable a competing default CodeQL setup.
- GitHub plan and organization policy can legitimately make some security capabilities unavailable; a precise limitation record satisfies the truthful-delivery requirement for those capabilities.
- Operational repository mutations occur only after the S008 pull request supplies the hosted evidence required by their dependency order.
