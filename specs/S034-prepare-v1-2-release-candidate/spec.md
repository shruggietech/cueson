# Feature Specification: Prepare the v1.2.0 release candidate

**Feature Branch**: `codex/S034-prepare-v1-2-release-candidate`

**Created**: 2026-10-06

**Status**: Analysis and local verification passed; final-head hosted checks and review pending

**Input**: The owner requests S034 as recommended, using installed Spec Kit and unattended delivery through the public release authorization boundary.

## User Scenarios & Testing

### User Story 1 - Adopt the reviewed speaker contract (Priority: P1)

A downstream application author can target an exact stable 1.2.0 candidate containing S033's optional consumer speaker attribution and declared-media timing, with established native workflows and source preservation intact.

**Why this priority**: Applications need an unambiguous contract and matching executable before adopting the feature.

**Independent Test**: A candidate emits its matching schema, accepts annotated current documents and supported historical inputs, preserves consumer values, and restores source bytes exactly.

**Acceptance Scenarios**:

1. **Given** reviewed S033, **When** a candidate is built, **Then** executable, emitted documents, packaged schema and immutable candidate schema identify exact 1.2.0.
2. **Given** valid 1.0.0 and 1.1.0 inputs, **When** the candidate processes them, **Then** their identity and original bytes remain unchanged.
3. **Given** the published v1.1.0 consumer, **When** every applicable consuming operation receives new documents, **Then** it rejects before creating or replacing output.

### User Story 2 - Assess one verified candidate bundle (Priority: P1)

A maintainer can assess six candidate packages tied to one exact source revision, with same-bundle execution on Windows, Linux and macOS amd64.

**Why this priority**: A version change alone cannot prove that downstream packages enforce the contract.

**Independent Test**: Verify six archives, six SBOMs and one checksum manifest; execute compatible packages; reject tampering, mismatched identities and stale evidence.

**Acceptance Scenarios**:

1. **Given** one accepted bundle, **When** compatible native hosts verify it, **Then** each proves four-format workflows, twelve conversions, strict loss refusal, consumer storage/timing and historical compatibility.
2. **Given** altered package/schema/revision/old-consumer binding, **When** proof runs, **Then** it fails without weakening checks.
3. **Given** arm64 packages, **When** structural proof passes, **Then** evidence does not claim native execution.

### User Story 3 - Make a concrete publication decision (Priority: P2)

The owner receives concise release highlights and a decision procedure distinguishing reviewed preparation from public availability and production hosting.

**Why this priority**: Public prose must distinguish actual availability from candidate readiness.

**Independent Test**: Changelog, notes, roadmap and publication contract agree on identity, scope, evidence and authority.

**Acceptance Scenarios**:

1. **Given** candidate preparation, **When** documentation is reviewed, **Then** v1.1.0 remains the latest published release.
2. **Given** a future accepted main revision, **When** publication is proposed, **Then** the decision binds actual post-merge main, fresh artifact/run IDs and expiry, thirteen asset sizes/digests and three native proofs without predicting a squash revision.
3. **Given** newly arrived issues, **When** kickoff establishes scope, **Then** native planning records document inclusions, exclusions and dependencies.

### Edge Cases

- Former 1.2.0-dev identity is not admitted as a historical public contract.
- Historical schemas, published asset contracts and dated planning snapshots remain unchanged.
- Optional consumer values retain S033's exact ownership, Unicode, timing, bounds, privacy and loss rules.
- An unavailable old-consumer download or expired artifact cannot be reported as completed proof.
- Old consumers reject real current documents and minimal identity probes without payload or occupied-output mutation.
- Cue coverage never becomes inferred media duration; unavailable declaration checks are never described as passed.
- Changed source, notes, tools, artifact bytes, reviews or expiry invalidate the pending publication decision.
- Public download metadata and production schema routes remain on verified releases.

## Requirements

### Functional Requirements

- **FR-001**: Promote only the reviewed current contract to exact 1.2.0, preserving S033 behavior and existing stable native capabilities.
- **FR-002**: Keep executable, embedded/emitted schema, documents and six packages in exact lockstep with a byte-identical immutable 1.2.0 candidate schema.
- **FR-003**: Preserve released schemas/contracts; retain exact 1.0.0/1.1.0 input support and reject former development identities.
- **FR-004**: Packaged native proof covers current consumer storage, inert identifiers, timing checks, unavailable/evaluated media states, exact restoration and deterministic omission/strict refusal.
- **FR-005**: Verify six historical input paths across validation, count-only inspection, restoration, matching rendering and strict conversion, preserving identity.
- **FR-006**: Verify the checksummed published v1.1.0 consumer on each compatible amd64 host; all applicable paths reject current documents and identity probes without publication or mutation. Preserve existing published v1.0.0 proof.
- **FR-007**: Verify one non-publishing bundle of six archives, six target-bound SBOMs and one checksum manifest bound to the exact source revision and immutable schema.
- **FR-008**: Execute that same bundle on three amd64 native hosts; preserve pure-Go builds and structural arm64 proof without claiming arm64 execution.
- **FR-009**: Finalize a dated 1.2.0 changelog section and substantially shorter release highlights, preserving released history and one fresh Unreleased section.
- **FR-010**: Reconcile current-state documentation with S030-S033 completion while retaining dated snapshots and distinguishing candidate, public release and production hosting.
- **FR-011**: Prepare the publication decision contract requiring actual reviewed main, fresh proof, asset/native/notes bindings, expiry and explicit authority; altered or expired evidence cannot be substituted silently.
- **FR-012**: Establish native v1.2.0 milestone/issues with independent outcomes, native dependencies and exactly one Project item per issue using Stage and Slice only.
- **FR-013**: Complete installed Spec Kit analysis, local/hosted gates and review remediation with at most two automated Codex rounds and truthful verification records.
- **FR-014**: Do not create/push tags, publish release/assets, add production schema routes, change public download identity or deploy production within preparation.

### Key Entities

- **Stable candidate**: Exact 1.2.0 contract and packages, distinct from public availability.
- **Candidate bundle**: Thirteen intended public files and revision-bound verification evidence.
- **Native proof**: Same-bundle host execution, current/historical workflows and published-consumer refusal records.
- **Publication decision**: Exact later tag/release state, immutable bindings, expiry, recovery and authority.
- **Delivery issue**: Independently closeable candidate or later publication/hosting outcome represented once in native planning.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Six packages agree on exact 1.2.0 and matching immutable schema bytes.
- **SC-002**: Three native hosts prove the same bundle's current and historical workflows and all old-consumer refusal/probe paths without mutation.
- **SC-003**: Fourteen requirements have complete task coverage and no unresolved critical/high analysis finding.
- **SC-004**: Public prose agrees with candidate and existing public v1.1.0 state; preparation performs no release/production publication.
- **SC-005**: Final handoff identifies the concrete candidate head, executed checks/reviews, limits and remaining publication authorization boundary.

## Assumptions

- S033 merged through PR #83; issue #82 is complete. Kickoff found no open issues or new arrivals.
- This minor release introduces no new format family, diarization/probing engine or public Go API.
- The owner's unattended request carries repository preparation forward; final merge and public actions remain subject to applicable action-specific governance.
- Network/native proof failures remain incomplete evidence rather than fabricated success.

## Clarifications

### Session 2026-10-06

- Autopilot: stable 1.2.0 is candidate identity, not public availability. Preserve public v1.1.0 downloads until authorized publication.
- Autopilot: retain frozen historical proof and add separate v1.2 records; authenticate the published v1.1 consumer for new-version refusal.
- Autopilot: publication and hosting are later independent issues with native dependencies and separately granted action authority.
- Autopilot: do not admit former 1.2.0-dev as historical input; it was staged development rather than immutable released identity.
