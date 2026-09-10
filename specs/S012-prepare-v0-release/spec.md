# Feature Specification: Prepare the v0.0.0 Release

**Feature Branch**: `S012-prepare-v0-release`

**Created**: 2026-09-10

**Status**: Draft

**Input**: User description: "Use Spec Kit and the repository autopilot protocol to prepare the v0.0.0 release state, complete implementation and verification, push the branch, publish the official pull request, address no more than two Codex review rounds, and stop only when the pull request is green and fully reviewed."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Review a Complete Release Record (Priority: P1)

A maintainer can review the complete v0.0.0 history, immutable schema candidate, and concise publication notes together before authorizing any protected release action.

**Why this priority**: The release cannot be considered ready while its permanent contract bytes and human-readable history remain provisional or scattered.

**Independent Test**: Compare the versioned schema with the canonical embedded schema byte for byte, inspect the dated changelog and release notes, and verify that version identity and documented limitations agree across all three surfaces.

**Acceptance Scenarios**:

1. **Given** the reviewed v0.0.0 candidate, **When** a maintainer compares the versioned release schema with the canonical schema source, **Then** the files are byte-identical and both declare version `0.0.0`.
2. **Given** the accumulated candidate history, **When** a maintainer reads the changelog, **Then** the v0.0.0 changes appear in a dated release section beneath a fresh empty `[Unreleased]` section.
3. **Given** the prepared GitHub release notes, **When** a maintainer inspects them, **Then** they contain concise highlights, state the envelope-only limitation, and end with the exact tagged changelog link.

---

### User Story 2 - Prove the Exact Release Candidate (Priority: P1)

An operator can identify one exact source revision and audit a deterministic evidence file proving the complete six-target artifact set, checksums, schema identity, executable identity, and SBOM binding for that revision.

**Why this priority**: Publication authority must apply to a concrete source and artifact set, not to a moving branch or an informal claim that CI passed.

**Independent Test**: Run the non-publishing release proof from a clean exact commit and verify that its evidence records the full source revision, release-schema digest, six archives, six SBOMs, six checksum entries, version `0.0.0`, and `published: false`.

**Acceptance Scenarios**:

1. **Given** a clean source revision, **When** the release proof builds and verifies the candidate, **Then** every binary and archive schema reports `0.0.0`, every archive checksum matches, every SBOM binds to that revision, and the evidence names the same full revision.
2. **Given** the repository schema copies, **When** the release verifier runs, **Then** it rejects a missing or byte-different `schema/releases/v0.0.0/cueson.schema.json` before accepting artifacts.
3. **Given** the release-preparation pull request is squash-merged, **When** default-branch release proof runs, **Then** the merge commit itself becomes the only proposed publication target and receives its own retained non-publishing evidence.

---

### User Story 3 - Hand Off a Governed Publication Decision (Priority: P2)

An operator receives one green, fully reviewed pull request and a precise post-merge ritual that preserves separate authorization for merge, tag creation, GitHub Release publication, milestone closure, and production-domain work.

**Why this priority**: Release readiness is useful only when the next decision is clear and no preparation step silently crosses a protected boundary.

**Independent Test**: Inspect the GitHub issue, Delivery Project item, official pull request, required checks, review state, and release documentation, then confirm the issue closes only on merge and no tag, release, milestone closure, published schema, or production mutation occurred.

**Acceptance Scenarios**:

1. **Given** the S012 work item, **When** its GitHub metadata is inspected, **Then** it appears exactly once in `cueson Delivery`, uses Slice `S012`, advances through the governed Stage values, leaves default Status empty, and belongs to milestone v0.0.0.
2. **Given** the official S012 pull request, **When** its body and review record are inspected, **Then** it contains a complete closing reference, every finding has a recorded disposition, no more than two Codex rounds occurred, and every current-head check is green.
3. **Given** a green and fully reviewed S012 pull request, **When** autopilot reaches its halt, **Then** merge, tag creation, GitHub Release publication, milestone closure, production schema hosting, and production configuration remain untouched pending specific operator authorization.

### Edge Cases

- A versioned release schema can be semantically equivalent but byte-different because of whitespace, encoding, line endings, or key ordering.
- A release note can be well formed but omit the exact final changelog URL or claim native SRT and WebVTT codec support that does not exist.
- A pull-request head commit cannot safely be recorded as the final release target because squash merge creates a different default-branch commit.
- A workflow can pass on the pull request while no proof exists for the eventual default-branch merge commit.
- A dirty working tree can make Go build metadata report a modified revision or create artifacts that do not correspond to a reviewable commit.
- GoReleaser or Syft output can contain an incomplete, unexpected, duplicated, or revision-mismatched artifact catalog.
- Upstream SBOM timestamps and document identifiers can differ without invalidating their stable semantic binding.
- Release preparation can accidentally add a tag trigger, write permission, publication token, signing claim, production deployment, or hidden publish path.
- The v0.0.0 milestone can be closed before public release and post-publication read-back are actually complete.
- Review remediation can change the pull-request head after checks or bot conclusions were recorded.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: S012 MUST create one actionable GitHub issue under milestone v0.0.0 with the repository-required sections, independently testable acceptance criteria, and a pull-request closing boundary.
- **FR-002**: The issue MUST appear exactly once in `cueson Delivery`, use Slice `S012`, use the lifecycle-appropriate Stage, and leave the default Status field unused.
- **FR-003**: The repository MUST contain `schema/releases/v0.0.0/cueson.schema.json` as the reviewed release-schema candidate and it MUST be byte-identical to `internal/schema/cueson.schema.json`.
- **FR-004**: Automated verification MUST reject a missing or byte-different versioned release schema and MUST verify its schema identifier and `schema_version` constant against software version `0.0.0`.
- **FR-005**: GoReleaser archives MUST package the byte-identical versioned release schema without changing the exact six-target, four-member archive contract.
- **FR-006**: `CHANGELOG.md` MUST retain a fresh `[Unreleased]` section and move the accumulated candidate history into `## [0.0.0] - 2026-09-10` with valid Keep a Changelog comparison links.
- **FR-007**: The repository MUST contain concise v0.0.0 GitHub release notes that describe highlights only, state the absence of native SRT and WebVTT ingest, render, and conversion, and end exactly with `Full changelog: https://github.com/shruggietech/cueson/blob/v0.0.0/CHANGELOG.md`.
- **FR-008**: Release verification evidence MUST record the exact full source revision, release version, immutable-schema SHA-256, complete archive and SBOM inventory, checksum count, host execution result, and `published: false`.
- **FR-009**: Release verification MUST continue to prove six expected archives, six target-bound SBOMs, six archive checksums, four exact archive members, pure-Go target metadata, clean VCS state, trimmed build paths, embedded schema bytes, and release-version markers.
- **FR-010**: The release-proof workflow MUST run for pull requests, manual dispatch, and pushes to `main` so the eventual squash-merge commit can receive independent candidate proof.
- **FR-011**: The release-proof workflow MUST remain non-publishing, read-only, secretless, credential-free after checkout, without tag triggers, identity-token permission, signing, production deployment, or a release-capable GoReleaser configuration.
- **FR-012**: Documentation MUST distinguish pull-request-head evidence from the exact default-branch evidence required before tag or release authorization and MUST describe how the merge commit becomes the proposed publication target.
- **FR-013**: The release process MUST present the operator with the proposed full commit, version, changelog, release-schema digest, artifact inventory, checksums, SBOM evidence, release notes, green checks, and known limitations before a publication decision.
- **FR-014**: Repository-authored text added or changed by S012 MUST use UTF-8 without BOM, follow `.gitattributes`, avoid mojibake, and use one source line per Markdown paragraph or list item.
- **FR-015**: Existing formatter, documentation, product, policy, release-verifier, workflow, race, vet, vulnerability, CodeQL, and hosted CI gates MUST pass for the final pull-request head.
- **FR-016**: S012 MUST address every Codex, security, and review finding, resolve a thread only after its concern is handled, request at most one second Codex round when required, and never request a third round automatically.
- **FR-017**: S012 MUST stop for the operator's final review and merge ritual and MUST NOT merge or auto-merge the pull request, create or move a tag, publish a GitHub Release or release asset, close the v0.0.0 milestone, publish a schema, or mutate production configuration.

### Key Entities

- **Release Schema Candidate**: The versioned repository copy of the v0.0.0 schema whose bytes must match the canonical embedded source and every packaged copy.
- **Release Notes**: The concise, publication-ready highlights document with the exact tagged changelog link and truthful capability limitations.
- **Candidate Revision**: A full clean Git commit associated with a verified artifact set; after squash merge, the exact `main` merge commit is the proposed publication target.
- **Release Evidence**: The deterministic verifier output binding schema digest, archives, checksums, SBOMs, host execution, version, publication state, and source revision.
- **Publication Decision Package**: The reviewed repository and hosted evidence presented to the operator before any protected tag or release action.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The canonical, versioned repository, embedded binary, and six packaged schema copies match at 100% byte identity and all report version `0.0.0`.
- **SC-002**: The verifier accepts exactly six archives, six SBOMs, six checksum entries, and four members per archive, with zero missing, additional, duplicated, or revision-mismatched items.
- **SC-003**: The release evidence contains one 40-character lowercase source revision, one 64-character lowercase release-schema SHA-256, `published: false`, and complete target evidence for all six platform and architecture pairs.
- **SC-004**: The changelog contains one dated v0.0.0 section and one fresh `[Unreleased]` section, while the release notes end with the exact required URL and contain zero native-codec capability claims.
- **SC-005**: A push of the eventual S012 squash-merge commit to `main` automatically runs the same non-publishing release proof used on the pull request without gaining write, secret, tag, signing, release, or production authority.
- **SC-006**: The final pull-request head has zero failed or pending required checks, zero unresolved review threads, every review finding addressed, and no more than two Codex review rounds.
- **SC-007**: The S012 pull request contains a complete closing reference and the issue's Project entry uses Slice `S012`, Stage `PR review`, and an empty default Status field while awaiting merge.
- **SC-008**: S012 performs zero merges, tag creations or movements, GitHub Release or asset publications, milestone closures, production schema publications, and production configuration changes.

## Assumptions

- Version `0.0.0`, release date `2026-09-10`, the six-target matrix, GoReleaser v2.18.1, Syft v1.51.1, and the four-member archive contract remain approved for this preparation slice.
- The existing canonical schema and executable version are already `0.0.0`; S012 freezes and proves those bytes rather than changing the public contract.
- The final release target cannot be known until the release-preparation pull request is squash-merged. The exact resulting `main` commit becomes the candidate only after its own default-branch proof succeeds.
- Upstream SBOM files are verified for stable semantic identity and source binding, not byte-for-byte reproducibility.
- Pull-request publication authority in this request does not authorize merge, tag creation, GitHub Release publication, milestone closure, or production changes.
