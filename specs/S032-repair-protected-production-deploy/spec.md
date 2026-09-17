# Feature Specification: Repair protected production deployment

**Feature Branch**: `codex/S032-repair-production-deploy`

**Created**: 2026-09-17

**Status**: Draft

**Input**: User description: "Use S032 to restore the protected default-branch production path end to end under the autopilot protocol."

## User Scenarios & Testing

### User Story 1 - Run the protected deployment command path (Priority: P1)

An authorized operator can dispatch the production workflow from the default branch at one exact current-main revision, and every verifier receives the arguments it declares.

**Why this priority**: The current workflow stops before mutation because its package-script boundary forwards an unsupported literal separator. Until that boundary is repaired, the durable protected path cannot deploy.

**Independent Test**: Exercise both verifier package scripts through the same command boundary used by the workflow and prove valid named arguments reach each parser without a literal separator.

**Acceptance Scenarios**:

1. **Given** an exact lowercase current-main revision and a default-branch dispatch, **When** the workflow runs its preflight and public verification commands, **Then** both parsers receive their declared named arguments and no standalone separator.
2. **Given** a dispatch from another ref or a revision different from the workflow execution revision, **When** validation runs, **Then** the workflow fails before checkout-controlled tools or credentials are used.
3. **Given** a correctly dispatched main workflow, **When** exact-main proof runs before artifact verification and again before mutation, **Then** any remote movement stops deployment.

---

### User Story 2 - Operate a durable least-privilege credential (Priority: P1)

A repository operator can provision, rotate, revoke, and audit one protected Cloudflare credential with only the account, zone, and permissions required by deployment and full state read-back.

**Why this priority**: The temporary recovery credential was removed. A durable protected path needs a reproducible credential contract without exposing secret material to unrelated workflow steps.

**Independent Test**: Review the credential contract and workflow environment boundaries, then prove the secret is exposed only to the three steps that call Cloudflare and that those steps cover deployment plus full pre/post state verification.

**Acceptance Scenarios**:

1. **Given** the documented token policy, **When** an operator creates the credential, **Then** its owner, resource scope, permissions, rotation interval, expiry policy, and emergency revocation procedure are unambiguous.
2. **Given** a protected workflow run, **When** checkout, setup, installation, tests, or local artifact proof execute, **Then** they cannot read the Cloudflare token.
3. **Given** deployment and state verification steps, **When** they execute, **Then** the token is available only to those steps and supports the complete declared read-back.

---

### User Story 3 - Prove the default-branch production path after merge (Priority: P2)

After human merge and durable credential provisioning, an authorized operator can run the default-branch workflow without a recovery branch and obtain matching GitHub and public deployment evidence for the exact main revision.

**Why this priority**: Source review alone cannot prove the protected environment, Cloudflare permissions, or public deployment record. The post-merge run closes the operational defect.

**Independent Test**: Dispatch the merged default-branch workflow with exact current main, then verify the GitHub deployment revision, public deployment record, Cloudflare before/after snapshots, all routes, downloads, redirects, metadata, manifest, TLS, DNS, and immutable schema bytes.

**Acceptance Scenarios**:

1. **Given** the merged workflow and protected credential, **When** the authorized production run completes, **Then** it deploys without an execution-only branch and all workflow checks pass.
2. **Given** a successful run, **When** deployment evidence is read back, **Then** GitHub environment metadata and public deployment metadata identify the same exact main revision.
3. **Given** preflight and post-deployment snapshots, **When** they are compared, **Then** intended bindings are correct and unrelated DNS, Workers, Custom Domains, and redirect rulesets are preserved.

### Edge Cases

- A syntactically valid revision is rejected when it differs from the workflow execution revision or freshly fetched `origin/main`.
- A workflow dispatched from a tag, pull-request ref, or recovery branch fails before any protected credential is made available.
- Missing, expired, revoked, under-scoped, or wrong-account credentials fail before mutation or during the explicit Cloudflare preflight with actionable output that does not reveal secret material.
- An independently moved default branch between the first and second exact-main checks stops the run before production mutation.
- Cloudflare read-back differences in unrelated state fail the run instead of being silently accepted.
- Public metadata that is internally consistent but identifies another revision fails exact local-authority comparison.

## Requirements

### Functional Requirements

- **FR-001**: The production workflow MUST pass Cloudflare preflight, Cloudflare read-back, and public verification arguments exactly as their parsers declare, without forwarding a standalone separator.
- **FR-002**: Automated tests MUST execute the actual package-script command boundary for both verifiers and distinguish accepted named arguments from malformed forwarding.
- **FR-003**: The workflow MUST accept only a full lowercase commit identifier that equals both the default-branch workflow execution revision and freshly fetched `origin/main`.
- **FR-004**: The workflow MUST repeat exact-main proof immediately before checkout-controlled artifact proof and immediately before production mutation.
- **FR-005**: The Cloudflare credential MUST be stored in the protected GitHub `production` environment and exposed only to preflight, deployment, and post-deployment read-back steps.
- **FR-006**: The credential contract MUST identify the organization-owned account and `cueson.io` zone resources plus the minimum permissions required to upload the Worker and read zone identity, DNS, Worker inventory, Custom Domains, and redirect rulesets.
- **FR-007**: The credential contract MUST document ownership, creation validation, rotation no less frequently than every 90 days, expiry handling, emergency revocation, and replacement read-back without storing secret material.
- **FR-008**: The workflow MUST publish a GitHub production deployment record whose revision is the same exact main revision recorded by the public deployment metadata.
- **FR-009**: Full preflight and post-deployment verification MUST prove zone identity, intended bindings, and preservation of unrelated DNS, Workers, Custom Domains, and redirect rulesets.
- **FR-010**: The reviewable pull request MUST close an independently testable implementation issue while retaining the production validation issue for the post-merge run.
- **FR-011**: After human merge, the explicitly authorized validation deployment MUST run from the default-branch workflow without a recovery branch and pass the complete production verifier.
- **FR-012**: Temporary recovery credentials, remote branches, and other mutable recovery state MUST remain absent after durable-path validation; any retained local branch must be reported and removed only when proven safe and clean.

### Key Entities

- **Deployment revision**: One full lowercase commit identifier that must equal workflow execution, freshly fetched default branch, GitHub deployment metadata, and public deployment metadata.
- **Protected production credential**: A secret token owned by the organization, restricted to the target account and zone, and available only inside Cloudflare-facing workflow steps.
- **Cloudflare state snapshot**: A before-or-after normalized record of zone identity, DNS, Workers, Custom Domains, and redirect rulesets used to prove intended change and unrelated-state preservation.
- **Implementation issue**: The independently closeable source, test, and documentation outcome reviewed in the pull request.
- **Production validation issue**: The operational outcome that remains open until the merged default-branch workflow and live public state are proven.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Both verifier package-script boundary tests reach their intended validation logic and zero tested invocations report a standalone separator as an argument.
- **SC-002**: All non-Cloudflare workflow steps have zero references to the protected token, and exactly three Cloudflare-facing steps receive it.
- **SC-003**: Invalid ref, workflow-revision mismatch, and moved-main cases all stop before production mutation in automated policy coverage.
- **SC-004**: One authorized post-merge run from the default branch completes with matching GitHub and public exact-main revision evidence and no recovery branch.
- **SC-005**: The post-merge production verifier proves 25 routes, seven downloads, three immutable schemas, DNS/TLS/redirect behavior, metadata, manifest integrity, and full Cloudflare state preservation.

## Assumptions

- Production is healthy at release v1.1.0 and exact main commit `926b062f7417580b8fac83d906b0fd34dfda3223` before this slice.
- The repository operator will provision or approve the durable account-owned token after reviewing the documented least-privilege contract; no token value belongs in source, logs, issues, or pull-request text.
- The pull request can prove source behavior and credential isolation, while the default-branch production run necessarily occurs after human merge.
- The user's standing authorization covers the necessary production validation after merge, but does not authorize this agent to merge the pull request.
