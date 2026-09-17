# Feature Specification: Prepare v1.1.0 public schema and site

**Feature Branch**: `codex/S030-publish-v1-1-public-schema-site`

**Created**: 2026-09-16

**Status**: Specified

**Input**: Operator S030 autopilot with explicit branch push and official pull-request authorization after independently verified v1.1.0 publication. Preparation issue #76 is a sub-issue and blocker of production outcome #67. Production mutation and final merge retain separate authority.

## Clarifications

### Session 2026-09-16

- Q: Can the official pull request close #67? A: No. It closes independently testable preparation #76 and references #67, which remains open for post-merge exact-main production deployment, live verification, epic reconciliation and milestone closure.
- Q: What is the current public release identity? A: Final non-draft GitHub Release `v1.1.0` targets `7ff45c1d8cd8df377e1fb568b9785286b649fd7c` with six archives, six target-specific SPDX JSON SBOMs and one checksum manifest; #66 independently verified and closed that outcome.
- Q: Which schema routes must exist? A: Add exact immutable v1.1.0 beside unchanged v0.0.0 and v1.0.0 routes, and continue to omit every mutable `latest` schema alias.
- Q: May this slice deploy production? A: No. It must make the reviewed repository artifact ready and strengthen exact-main verification, but cueson.io mutation requires a separate explicit authorization after human merge identifies the actual reviewed main revision.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Find and use the current release (Priority: P1)

A visitor can identify v1.1.0 as the current stable release, understand four-format support and compatibility limits, and reach the exact package, checksum, SBOM, schema and release details from the public artifact.

**Why this priority**: The public site still presents v1.0.0 and candidate-era guidance after v1.1.0 became the verified official release.

**Independent Test**: Build the public artifact, follow the landing and maintained documentation routes, and verify the displayed release identity, format state, links and compatibility statements against the official release.

**Acceptance Scenarios**:

1. **Given** independently verified v1.1.0 exists, **When** a visitor opens the landing page, **Then** v1.1.0 is identified as current, ASS/SSA and SubRip/WebVTT support is described truthfully, and install/release/schema actions lead to exact versioned targets.
2. **Given** a visitor needs a platform package, **When** they follow a primary download, **Then** the target is one of the six exact v1.1.0 archives or the exact checksum manifest and no primary action silently falls back to v1.0.0.
3. **Given** older released documentation remains part of the public record, **When** the artifact is generated, **Then** v0.0.0 and v1.0.0 release pages and historical statements remain available without being rewritten as current behavior.

### User Story 2 - Resolve every immutable schema (Priority: P1)

A consumer can resolve the canonical schema URI emitted by Cueson v0.0.0, v1.0.0 or v1.1.0 and receive the exact immutable released bytes.

**Why this priority**: v1.1.0 output already identifies its versioned public URI, while public hosting still lacks that route.

**Independent Test**: Generate the deployable artifact, retrieve all three versioned schema paths, compare bytes and SHA-256 values with repository authorities, and prove mutable aliases remain absent.

**Acceptance Scenarios**:

1. **Given** a consumer requests `/schema/v1.1.0/cueson.schema.json`, **When** the artifact serves the path, **Then** the response is exactly 185,641 bytes with SHA-256 `223b61cbcf6337167039268576b2c739564585fff10e6cc6076e47a03526a0f7`.
2. **Given** a consumer requests either older canonical schema route, **When** the artifact serves it, **Then** its body remains byte-identical to the corresponding immutable repository copy.
3. **Given** a consumer requests `/schema/latest/cueson.schema.json`, **When** the artifact handles the request, **Then** no schema document is returned.

### User Story 3 - Review an exact deployable update (Priority: P2)

The operator receives one official, fully verified pull request whose artifact and deployment record can later be reproduced from the exact merged main revision without granting the pull request production authority.

**Why this priority**: Production should use the reviewed source and complete declared inventory, while deployment credentials and mutation remain owner-controlled.

**Independent Test**: Audit generated inventories and deployment metadata, reproduce the artifact from the pull-request head, run every local and hosted gate, and verify that deployment can select only the exact current main revision after merge.

**Acceptance Scenarios**:

1. **Given** repository content or release metadata changes, **When** generation is omitted or inconsistent, **Then** drift, inventory, link, metadata or schema verification fails with an actionable result.
2. **Given** an official pull request, **When** hosted checks run, **Then** they receive no production credentials and perform no production mutation.
3. **Given** the human later merges the reviewed pull request, **When** a separately authorized deployment is selected, **Then** the workflow requires the selected revision to equal current default-branch head and verifies public metadata against the locally generated declared inventory.

### Edge Cases

- A release link uses the correct tag but an unapproved asset name, omits a declared download, or points to a mutable page where an exact asset is required.
- Candidate-era prose survives in a maintained current-state document or current release page after publication.
- An older release page or schema route is dropped while adding v1.1.0.
- The immutable v1.1.0 schema copy is normalized, regenerated, truncated, or differs by line ending.
- Generated route, download, schema or deployment metadata describes content that the artifact does not contain.
- A remote deployment self-reports an incomplete inventory that a verifier would accept without comparing it to the local reviewed declaration.
- A stale default-branch ancestor is selected for deployment after newer main changes exist.
- A pull-request or unreviewed revision can reach credentials or a production mutation path.
- A mutable `latest` schema path becomes reachable through generated content, redirects or worker fallback.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Public-facing current-state content MUST identify independently verified v1.1.0 as the current stable release and describe stable SubRip, WebVTT, ASS and SSA workflows with truthful compatibility and fidelity limits.
- **FR-002**: The public artifact MUST include a maintained v1.1.0 release page and preserve the existing v0.0.0 and v1.0.0 release material.
- **FR-003**: Primary install and download actions MUST use exact official v1.1.0 release or asset targets, including six platform archives and the checksum manifest, without inventing mutable aliases.
- **FR-004**: Current release identity MUST have one validated maintained source from which generated page metadata, deployment metadata and user-facing release actions agree.
- **FR-005**: The artifact MUST serve the exact immutable v1.1.0 schema bytes at their canonical versioned route while preserving byte-identical v0.0.0 and v1.0.0 routes.
- **FR-006**: The artifact MUST NOT provide a mutable `latest` schema document or route any unknown schema version to a released schema.
- **FR-007**: Generated content and inventories MUST include every declared route, download and schema exactly once and MUST reject stale or divergent generated adaptations.
- **FR-008**: Release and schema metadata MUST record enough version, route, length and digest identity to prove the generated artifact against the maintained authorities.
- **FR-009**: Artifact and browser verification MUST cover navigation, links, fragments, metadata, accessibility, responsive layouts, all declared routes, all declared downloads and all three schema identities.
- **FR-010**: Public verification logic MUST compare remote deployment metadata with the locally reviewed declaration and MUST independently fetch and hash the public content manifest instead of trusting a remote self-declared subset.
- **FR-011**: Production deployment selection MUST require the selected revision to equal the current default-branch head immediately before proof and mutation; an older ancestor MUST be rejected.
- **FR-012**: Pull-request checks MUST remain read-only, receive no production credentials and contain no event path that performs production deployment.
- **FR-013**: Repository documentation and detailed changelog MUST describe the final released state, exact historical compatibility, public route/download contract and separate production authority without rewriting historical release evidence or immutable schema bytes.
- **FR-014**: Installed Spec Kit design, blocking analysis, implementation, independent convergence and complete foreground verification MUST pass before publication of the official pull request.
- **FR-015**: The official pull request MUST close #76, reference #67, use intentional formatted/read-back publication bodies, reach green exact-head CI/security checks and address every review finding with no more than two Codex rounds.
- **FR-016**: Completion of this preparation MUST leave #67, epic #51 and milestone v1.1.0 open; Project Stage and Slice MUST remain truthful and the default Status field MUST remain unused.
- **FR-017**: This preparation MUST NOT mutate production cueson.io, create or move a tag, replace release assets, publish another release, or perform the final pull-request merge.

### Key Entities

- **Current release identity**: Exact stable version, tag, official release URL and approved download inventory shared by generated site surfaces.
- **Released schema**: Immutable version, canonical route, repository source, exact byte length, SHA-256 digest and media type.
- **Content inventory**: Complete declared public pages, downloads, schemas and generated source identities expected in one artifact.
- **Deployment record**: Exact source revision, release identity, route/download/schema inventory and content-manifest digest reproduced from reviewed source.
- **Preparation evidence**: Spec Kit gates, local artifact verification, official pull-request checks and terminal review results owned by #76.
- **Production outcome**: Separately authorized exact-main deployment, Cloudflare read-back, public DNS/TLS/HTTP/content proof and milestone reconciliation owned by #67.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All six #76 acceptance criteria have explicit requirement, task and verification coverage with zero material analysis or convergence findings.
- **SC-002**: Every maintained current-state surface names v1.1.0 consistently, seven primary download targets resolve to exact approved v1.1.0 assets, and zero candidate-availability claims remain outside historical preparation records.
- **SC-003**: The artifact declares and serves 22 HTML routes plus three immutable schema routes, with zero missing, duplicate or unexpected declared routes and no mutable schema alias.
- **SC-004**: All three public schema copies match their immutable repository authorities byte for byte and by SHA-256; v1.1.0 is exactly 185,641 bytes with digest `223b61cbcf6337167039268576b2c739564585fff10e6cc6076e47a03526a0f7`.
- **SC-005**: Complete site, documentation, repository quality, artifact and deployment dry-run checks pass with zero broken links, serious accessibility findings, responsive overflow, metadata defects or generation drift.
- **SC-006**: Official pull-request publication renders correctly, all exact-head checks are green, and no actionable review remains within the two-round Codex limit.
- **SC-007**: Zero production, tag, release, asset-replacement or final-merge transitions occur; #67, #51 and milestone v1.1.0 remain open with truthful native and Project relationships.

## Assumptions

- #66 independently verified the final v1.1.0 release, its exact reviewed revision, six archives, six SBOMs, checksum manifest, schema/legal identity and native behavior before S030 specification began.
- #76 is the repository-preparation child of #67 and blocks #67; closing #76 does not assert live public behavior.
- Existing v0.0.0 and v1.0.0 release pages, schemas and public routes remain immutable historical authorities.
- Root maintained documents and generated adaptations retain their established authority boundary; generated site content is not independently edited.
- The approved Cueson brand, apex-domain canonical origin, redirect-only `www` behavior and Cloudflare deployment architecture remain unchanged.
- Branch push and official pull-request publication are explicitly authorized. Production deployment and final merge are not authorized by this slice instruction.
