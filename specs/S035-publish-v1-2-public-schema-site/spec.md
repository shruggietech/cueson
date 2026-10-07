# Feature Specification: Prepare v1.2.0 public schema and site

**Feature Branch**: `codex/S035-publish-v1-2-public-schema-site`

**Created**: 2026-10-06

**Status**: Draft

**Input**: Owner kickoff for S035 under Spec Kit/autopilot with explicit push and official PR authority, at most two Codex review rounds, and a final human review/merge boundary.

## Clarifications

### Session 2026-10-06

- Q: Which outcome does this PR close? → A: Preparation #88; live activation #86 remains open, preserving its full acceptance criteria.
- Q: Does the published release justify serving production immediately? → A: Prepare and review now; the owner subsequently authorized deployment/live verification as part of S035 immediately after their final merge.
- Q: Should old release records be rewritten? → A: Update maintained current-state prose and preserve dated history, frozen slices and historical download links.
- Q: Does DNS resolution prove both address families can serve the site? → A: Later verification must separately demonstrate HTTPS connectivity over IPv4 and IPv6.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Discover the published release and speaker contract (Priority: P1)

A downstream developer can find the released v1.2.0 schema, download the official software and read how optional consumer speaker assignments and declared media bounds work.

**Why this priority**: Published software needs matching accessible documentation and schema discovery.

**Independent Test**: Inspect the reviewed artifact's release links, schema bytes, release page and speaker guide.

**Acceptance Scenarios**:

1. **Given** independently verified public v1.2.0 assets, **When** the artifact is generated, **Then** seven primary downloads reference their exact official URLs and the release page describes shipped behavior.
2. **Given** a consumer speaker identifier, **When** a developer reads the guide, **Then** its consumer-governed scope, validation limits and separation from native source observations are explicit.
3. **Given** optional declared media duration, **When** a developer reads the guide, **Then** assignment checks and source-cue warnings are distinguished and no last-cue timestamp is presented as proven audio duration.

### User Story 2 - Preserve existing consumers (Priority: P1)

Existing consumers continue using prior schema and release routes while new consumers discover the v1.2.0 contract.

**Why this priority**: Released contracts are immutable and older consumers retain valid entry points.

**Independent Test**: Compare all historical schemas with approved sizes/hashes and verify prior release links/routes.

**Acceptance Scenarios**:

1. **Given** released 0.0.0, 1.0.0 and 1.1.0 schemas, **When** adding 1.2.0, **Then** their bytes and public paths remain identical.
2. **Given** prior release documentation, **When** primary downloads advance, **Then** the old release pages retain working version-specific links.

### User Story 3 - Authorize an exact production deployment (Priority: P2)

The operator can review the finished PR and complete the final merge, after which the agent performs the already authorized production deployment with evidence bound to actual merged main.

**Why this priority**: Review, merge and production activation are distinct outcomes.

**Independent Test**: Verify the handoff and negative verification checks without production mutation.

**Acceptance Scenarios**:

1. **Given** PR preparation, **When** checks and reviews pass, **Then** the agent requests final human review/merge and leaves live activation pending.
2. **Given** a later merge and explicit production authorization, **When** production work starts, **Then** fresh current-main proof, deployment metadata, schema/documentation/download checks and separate IPv4/IPv6 reachability determine completion.
3. **Given** wrong revision, altered schema/metadata, incomplete activation or missing authority, **When** completion is assessed, **Then** production issue #86 and milestone completion remain pending.

### Edge Cases

- A published release is available while production still serves its predecessor: record both states distinctly.
- An old release URL remains intentionally historical, including examples and frozen specifications: do not globally rewrite history.
- Missing, altered or unexpected schema/download metadata fails verification.
- A stale main revision or failed address family prevents production completion.
- Production deployment credentials stay out of preparation and proof steps; use established protected workflow boundaries.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The prepared site MUST identify independently verified public v1.2.0 as its current release and bind seven primary downloads to authentic official assets.
- **FR-002**: The artifact MUST include immutable v1.2.0 schema bytes at `/schema/v1.2.0/cueson.schema.json`, matching its published package contract.
- **FR-003**: The three historical immutable schema byte sequences and public paths MUST remain unchanged; no mutable latest-schema alias may be introduced.
- **FR-004**: The artifact MUST include `/docs/releases/v1.2.0/` and `/docs/consumer-speakers/` with navigation, links and existing site presentation/accessibility behavior.
- **FR-005**: Prior release pages and their version-specific download links MUST remain available.
- **FR-006**: Maintained public prose MUST distinguish shipped v1.2.0 software, prepared site source and pending production activation, while preserving dated historical evidence.
- **FR-007**: Speaker prose MUST describe optional consumer-governed identifiers, established string limits, native observation separation and optional declared-media checks without universal identity or inferred audio-duration claims.
- **FR-008**: Independent inventory, browser, schema-byte, metadata and deployment tests MUST verify the expanded artifact and retain relevant negative checks.
- **FR-009**: Protected deployment MUST retain manual dispatch, exact current-main binding, step-scoped credentials, unrelated configuration preservation and independent live verification.
- **FR-010**: Following the final owner merge, S035 MUST execute fresh actual post-merge main proof, protected deployment under recorded production authority, four live schema identities, documentation/download/revision/header/redirect verification and separate IPv4/IPv6 HTTPS evidence. No second kickoff or approval is required under the owner steering instruction.
- **FR-011**: Preparation MUST have an independently closeable issue; production #86 and milestone v1.2.0 MUST remain open until live completion, with native dependencies and governed Project fields.
- **FR-012**: This run MUST complete local verification, push, official PR and every review finding within at most two Codex rounds, then stop for final human review/merge. Following confirmed merge, the same slice MUST deploy and independently verify production under the operator instruction without another kickoff or approval request.

### Key Entities

- **Release declaration**: Exact published version/tag/release URL and seven official primary downloads.
- **Immutable schema inventory**: Four versioned paths with positive byte lengths and SHA-256 values.
- **Document inventory**: Maintained source documents mapped to public routes and navigation.
- **Deployment evidence**: Exact reviewed/merged revision, proof, authority and later independent live results.
- **Planning outcome**: Separately testable preparation and production activation issues with native milestone/dependency state.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All four prepared schema byte lengths/hashes match approved immutable sources, including all three historical identities.
- **SC-002**: All seven primary downloads target authentic v1.2.0 release assets, and all four release pages and the speaker guide are reachable in the artifact.
- **SC-003**: The artifact exposes exactly 28 declared routes (24 HTML, four schemas), with all applicable static/browser/accessibility checks passing.
- **SC-004**: Current maintained documents contain no stale claim that v1.2.0 software is unpublished; they distinguish pending production where relevant.
- **SC-005**: The official PR has green required checks and handled review findings with no more than two Codex rounds. After the owner merges it, production completion requires exact-main deployment and passing independent live evidence; S035 remains incomplete until then.

## Assumptions

- Completed #85 supplies independently verified publication evidence; read back current release identity and exact relevant assets during preparation.
- Existing static generation, Worker and protected deployment architecture is retained.
- The owner authorized push/PR and subsequently instructed S035 to include production deployment/live verification immediately after their final PR merge. Final merge remains human-owned.
- Product runtime behavior and released schema/package/tag bytes are outside this preparation change.
