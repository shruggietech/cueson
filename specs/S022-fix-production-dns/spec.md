# Feature Specification: Fix production DNS verification

**Feature Branch**: `codex/S022-fix-production-dns`

**Created**: 2026-09-15

**Status**: Approved for implementation after blocking analysis

**Input**: Operator kickoff: "Spec out S022, then drive it through autopilot, automatically push and open the official pull request, address every external review, allow at most one second review, and return for human final review and merge."

## Context and Scope

Issue [#49](https://github.com/shruggietech/cueson/issues/49) owns the complete outcome. Its S021 deployment dependency [#47](https://github.com/shruggietech/cueson/issues/47) is complete. The current verifier can reject a healthy IPv6-only hostname because an IPv4 lookup fails before IPv6 is attempted, and its independent DNS-over-HTTPS check accepts only IPv4 queries. This slice corrects that operational false failure and makes address evidence explicit.

Scope includes independent address-family checks for both production hostnames through the system DNS resolver and the configured DNS-over-HTTPS resolver, deterministic regression coverage, site-suite integration, and narrowly relevant deployment documentation and changelog updates. TLS, redirects, routes, downloads, deployment metadata, schema verification, release artifacts, and the owner-controlled deployment workflow retain their existing behavior. DNS configuration changes, production deployment, software releases, and merge are outside this slice's implementation authority.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Verify a healthy hostname using either address family (Priority: P1)

A maintainer can verify the apex and redirect hostname when IPv4, IPv6, or both are available, without requiring an unused family to resolve.

**Why this priority**: The current production deployment includes a healthy hostname that the verifier can falsely reject, blocking trustworthy deployment acceptance.

**Independent Test**: Supply controlled IPv4-only, IPv6-only, and dual-stack responses independently to each resolver and prove each resolver accepts usable addresses, including when the other family fails.

**Acceptance Scenarios**:

1. **Given** a hostname resolves only through IPv4, **When** each resolver verifies it, **Then** each accepts its non-empty IPv4 address evidence.
2. **Given** a hostname resolves only through IPv6 and IPv4 produces no answer or an error, **When** each resolver verifies it, **Then** IPv6 is still queried and each accepts its non-empty IPv6 address evidence.
3. **Given** a hostname is dual-stack, **When** responses arrive in different orders or contain duplicates, **Then** the verification decision and reported address inventory are deterministic.

### User Story 2 - Diagnose a hostname without usable address evidence (Priority: P1)

A maintainer receives a failure identifying the hostname, resolver, and both family outcomes when neither family supplies usable addresses. An alias or malformed response cannot falsely prove address resolution.

**Why this priority**: Removing false failures must preserve the verifier's ability to reject unresolved hostnames and invalid resolver evidence.

**Independent Test**: Supply empty, rejected, alias-only, wrong-family, malformed, and unsuccessful responses and require a useful diagnostic unless the other family supplies valid address evidence.

**Acceptance Scenarios**:

1. **Given** both families lack addresses, **When** a resolver verifies a hostname, **Then** verification fails with the hostname, resolver, and outcomes for IPv4 and IPv6.
2. **Given** an answer contains aliases or records without a correctly typed address, **When** DNS-over-HTTPS evidence is evaluated, **Then** that answer does not count as usable address evidence.
3. **Given** a family encounters a transport, protocol, or malformed-response error, **When** the other family supplies a usable address, **Then** that resolver succeeds; otherwise both family outcomes remain visible in the failure.

### User Story 3 - Complete production acceptance without changing deployment authority (Priority: P1)

A maintainer can run the complete site suite and the repaired verifier against the existing S021 deployment while preserving all independent acceptance checks and human delivery boundaries.

**Why this priority**: Issue closure requires real deployment verification, rather than only isolated resolver tests.

**Independent Test**: Run the site suite and full public verifier against the recorded deployed revision, inspect the unchanged manual deployment workflow, and bind the official pull request to #49.

**Acceptance Scenarios**:

1. **Given** the existing S021 artifact is deployed, **When** the full verifier checks its expected revision, **Then** DNS, TLS, redirects, routes, downloads, metadata, and both immutable schemas all pass.
2. **Given** a pull request or arbitrary working branch, **When** validation runs, **Then** it has no additional production mutation authority and the existing exact-main manual deployment boundary remains intact.
3. **Given** implementation and local verification are complete, **When** authorized publication and reviews finish, **Then** the maintainer receives a green reviewed pull request for human final review and merge.

### Edge Cases

- One family rejects before or after the other family completes.
- Successful address sets contain duplicate addresses or arrive in different orders.
- DNS-over-HTTPS returns a successful HTTP response with unsuccessful DNS status, no Answer field, empty answers, aliases only, or wrong-family records.
- Response bodies or answer collections are malformed, and a transport request fails or times out.
- System DNS and DNS-over-HTTPS disagree; each still independently needs usable address evidence.
- The live deployment revision differs from the recorded expected revision; metadata verification must fail truthfully.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Both `cueson.io` and `www.cueson.io` MUST be checked through the system DNS resolver and the configured DNS-over-HTTPS resolver, with each resolver independently requiring usable address evidence.
- **FR-002**: Each resolver MUST query IPv4 and IPv6 independently and accept a hostname when at least one family returns non-empty usable address evidence, including when the other family fails.
- **FR-003**: DNS-over-HTTPS evidence MUST come from a successful HTTP response and successful DNS status and include an actual address record of the requested family with a valid address value; aliases alone, wrong-family records, and malformed evidence MUST NOT count.
- **FR-004**: When neither family supplies usable evidence, failure MUST identify the hostname, resolver, and IPv4 and IPv6 outcomes, distinguishing absent addresses from operational or response errors.
- **FR-005**: Address inventories and failure ordering MUST be deterministic and MUST NOT depend on query completion order or duplicate addresses.
- **FR-006**: Deterministic tests MUST cover IPv4-only, IPv6-only, dual-stack, absent addresses, partial-family failures, invalid evidence, and diagnostic behavior, and MUST run in the existing complete site suite.
- **FR-007**: Existing TLS, redirect, route, download, metadata, immutable-schema, missing-alias, and deployment-authority checks MUST retain their behavior and MUST NOT be bypassed to claim acceptance.
- **FR-008**: Completion MUST include the complete site suite, a Wrangler dry run, and full public verification of the existing S021 revision `46838fd5cc888b299a09b89da676db0005b00f16`, without a production mutation.
- **FR-009**: Slice documentation and changelog MUST explain the corrected dual-family rule and preserve every acceptance criterion owned by #49.
- **FR-010**: The authorized official pull request MUST close #49, pass hosted checks, address every actionable review, and use at most two total Codex review rounds before human final review and merge.

### Key Entities

- **Hostname**: The apex or redirect hostname independently checked by both resolver paths.
- **Family outcome**: Usable addresses, absent addresses, or a query/response error for one requested address family.
- **Resolver outcome**: A deterministic usable address inventory or a diagnostic combining both failed family outcomes.
- **Deployment evidence**: Existing production revision, routes, downloads, schema identities, DNS, TLS, and redirect checks required for operational acceptance.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All IPv4-only, IPv6-only, and dual-stack acceptance scenarios pass independently for both resolver paths.
- **SC-002**: Every no-usable-address scenario fails with both family outcomes, and alias-only or invalid responses never create a false pass.
- **SC-003**: The complete site suite and full public verifier pass against the recorded production revision, with both hostnames and all existing non-DNS acceptance checks covered.
- **SC-004**: One official pull request owns #49, all hosted gates and actionable reviews are satisfied, and no production mutation or final merge is performed by the agent.

## Assumptions

- The existing Google DNS-over-HTTPS resolver remains the configured independent public resolver; switching providers is unnecessary.
- The existing site toolchain and test runner are sufficient; no runtime dependency or public product contract change is needed.
- #49 is the only active issue; its dependency #47 is complete and no release milestone is assigned.
- Operator authorization covers S022 branch pushes, official pull-request publication, and necessary review responses, including at most one second Codex request. Production mutation and final merge remain outside that authorization.
