# Feature Specification: Publish and Verify v0.0.0

**Feature Branch**: `S013-publish-v0-release`
**Created**: 2026-09-10
**Status**: Draft
**Input**: Publish the operator-authorized Cueson v0.0.0 tag and GitHub Release from the exact verified post-S012 candidate, independently verify the public result, and update repository documentation through a reviewed pull request while preserving separately governed boundaries.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Obtain the Authorized Release (Priority: P1)

A Cueson user can open the official v0.0.0 GitHub Release, understand its current envelope-only capability, and download the correct archive, checksum manifest, and matching software bill of materials for a supported platform.

**Why this priority**: Publication is the authorized outcome. A release that is incomplete, mislabeled, or ambiguous cannot serve users safely.

**Independent Test**: Inspect the public tag and release, confirm that the tag targets the authorized revision, the release is public and final, the reviewed notes render correctly, and exactly thirteen approved assets are present.

**Acceptance Scenarios**:

1. **Given** the authorized candidate revision, **When** the v0.0.0 tag is inspected locally and remotely, **Then** the annotated tag peels to exactly `b294a6952c8bd041d852c502f5d7206c0b58edd6`.
2. **Given** the reviewed release notes, **When** the public GitHub Release is opened, **Then** it is neither draft nor prerelease, uses tag `v0.0.0`, presents the reviewed highlights, states the envelope-only limitations, and ends with the tagged changelog link.
3. **Given** the accepted candidate artifact, **When** a user inspects release downloads, **Then** exactly six platform archives, six corresponding SPDX JSON SBOMs, and one checksum manifest are available with no internal build metadata or CI evidence file.

---

### User Story 2 - Independently Verify Publication (Priority: P1)

A maintainer can prove that every public download is the exact authorized candidate artifact and that publication changed no separately protected repository, milestone, schema-hosting, or production state.

**Why this priority**: An API success is insufficient. Public bytes and GitHub state must be independently compared with the accepted evidence before publication is considered complete.

**Independent Test**: Download all thirteen public assets into a clean temporary directory, verify the checksum manifest, inspect every archive and SBOM against the authorized revision and schema digest, and read back protected GitHub state.

**Acceptance Scenarios**:

1. **Given** the public assets, **When** the checksum manifest is applied, **Then** all six archive hashes match and the manifest contains no missing, duplicated, or unexpected entry.
2. **Given** each published archive, **When** its members and executable are inspected, **Then** it contains exactly the platform executable, `cueson.schema.json`, `LICENSE`, and `NOTICE`; every binary contains the verified `0.0.0` release marker and expected target build identity; the host-compatible binary reports version `0.0.0`; and the schema matches the immutable tagged repository schema byte for byte.
3. **Given** each published SBOM, **When** its stable semantics are inspected, **Then** it identifies the matching platform, version `0.0.0`, and source revision `b294a6952c8bd041d852c502f5d7206c0b58edd6`.
4. **Given** completed publication, **When** protected boundaries are audited, **Then** milestone v0.0.0 remains open, production `cueson.io` is unchanged, and no signature, attestation, or public schema endpoint was created.

---

### User Story 3 - Record the Released State (Priority: P2)

A maintainer or contributor can read the repository and see that v0.0.0 is publicly released, where to obtain it, what it supports, how publication was verified, and which future capabilities remain unavailable.

**Why this priority**: Repository documentation must agree with public release state without obscuring the narrow capability boundary or implying that later roadmap work has shipped.

**Independent Test**: Review the S013 documentation diff and official pull request, run offline documentation and repository-text checks, and verify the issue and Project lifecycle record.

**Acceptance Scenarios**:

1. **Given** a successful public release, **When** repository status documentation is read, **Then** it links to v0.0.0, identifies the exact tagged revision and verified publication, and no longer describes the project as unreleased.
2. **Given** the v0.0.0 capability boundary, **When** release-facing documentation is read, **Then** it continues to state that native SRT and WebVTT ingest, rendering, and conversion are unavailable.
3. **Given** the official S013 pull request, **When** its delivery record is inspected, **Then** it closes issue #27 on merge, has every current-head check green, resolves every review finding, uses no more than two Codex rounds, and stops for the operator's final merge ritual.

### Edge Cases

- The tag or release may already exist because of a partial retry; S013 must stop instead of moving, overwriting, or duplicating published state without fresh operator judgment.
- An uploaded asset can have the correct name but different bytes; public verification must compare content rather than names alone.
- GitHub can create a release successfully while one or more assets fail to upload; partial publication must remain visibly incomplete and block dependent actions.
- The release body can be altered by transport or incorrect source formatting; the published body must be read back and compared structurally with the reviewed notes.
- An annotated tag has a tag-object identifier distinct from its peeled commit; verification must prove the peeled target rather than comparing only the tag object.
- SBOM document timestamps and namespace identifiers can vary upstream, so verification must use the exact accepted public files and their stable target, version, and source semantics without claiming reproducible SBOM bytes.
- The retained workflow artifact expires; if the accepted asset source becomes unavailable before publication completes, S013 must stop rather than silently rebuild or substitute artifacts.
- Repository automation can populate the default Project Status field; S013 must clear it and use only the governed Stage field.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: S013 MUST have one actionable GitHub issue in milestone v0.0.0, represented exactly once in `cueson Delivery` with Slice `S013`, lifecycle-appropriate Stage, and empty default Status.
- **FR-002**: Publication MUST use tag name `v0.0.0` and the tag MUST peel to exactly `b294a6952c8bd041d852c502f5d7206c0b58edd6`.
- **FR-003**: The tag MUST be an unsigned annotated tag and MUST NOT be moved or replaced after publication.
- **FR-004**: The GitHub Release MUST be public, non-draft, non-prerelease, associated with tag `v0.0.0`, and use the reviewed `docs/releases/v0.0.0.md` content as its body.
- **FR-005**: The public Release MUST contain exactly thirteen assets: six verified platform archives, six corresponding SPDX JSON SBOMs, and `cueson_0.0.0_checksums.txt`.
- **FR-006**: Every public asset MUST come from accepted default-branch release-proof run `34543376814`; S013 MUST NOT rebuild, regenerate, rename, or substitute any release asset.
- **FR-007**: The public Release MUST NOT include `release-evidence.json`, GoReleaser metadata, GoReleaser configuration, extracted binary directories, or any unexpected file.
- **FR-008**: Post-publication verification MUST read back the remote tag reference, peel the annotated tag, and prove its exact authorized commit.
- **FR-009**: Post-publication verification MUST read back the release tag, target, draft state, prerelease state, body, and complete asset metadata directly from GitHub.
- **FR-010**: Post-publication verification MUST download all thirteen public assets into a clean location and verify the six-entry archive checksum bijection.
- **FR-011**: Post-publication verification MUST prove every archive has exactly four safe regular members, uses the correct platform executable name, contains a binary with the verified `0.0.0` release marker and expected target build identity, and packages schema bytes matching the tagged immutable repository schema with SHA-256 `d15c7fa5227156109dd6be3d39b711aca3503794bb862169dfca96ee80adb975`; the host-compatible public binary MUST also report version `0.0.0` at runtime.
- **FR-012**: Post-publication verification MUST prove each SBOM's stable semantics identify its corresponding platform, version `0.0.0`, and source revision `b294a6952c8bd041d852c502f5d7206c0b58edd6`.
- **FR-013**: Repository documentation MUST identify v0.0.0 as publicly released, link the release and immutable tag, record the exact verified revision and evidence, and retain truthful envelope-only limitations.
- **FR-014**: Repository-authored text MUST use UTF-8 without BOM, follow `.gitattributes`, avoid mojibake, and retain one source line per Markdown paragraph or list item.
- **FR-015**: The official S013 pull request MUST contain `Closes #27`, pass the complete local and hosted verification suite, address every Codex, security, and review finding, request at most one second Codex round when required, and never request a third automatically.
- **FR-016**: S013 MUST stop for the operator's final review and merge ritual and MUST NOT merge or auto-merge its pull request.
- **FR-017**: S013 MUST NOT close milestone v0.0.0, publish a schema endpoint, mutate production `cueson.io`, add signatures or attestations, or claim native SRT/WebVTT ingest, rendering, or conversion.

### Key Entities

- **Authorized Candidate**: The exact default-branch commit and accepted release-proof artifact approved for public publication.
- **Release Tag**: The immutable annotated `v0.0.0` reference whose peeled target identifies the authorized candidate.
- **GitHub Release**: The public, non-draft, non-prerelease release record associated with the tag and reviewed notes.
- **Public Asset Set**: The exact six archives, six target-bound SBOMs, and one checksum manifest exposed to users.
- **Publication Evidence**: The read-back and independent verification record proving public GitHub state and bytes match the authorized candidate.
- **Release Documentation Record**: Repository-authored status, process, verification, schema, and architecture text updated after successful publication.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: One immutable annotated tag named `v0.0.0` peels to the authorized 40-character commit with zero target mismatch.
- **SC-002**: One public GitHub Release is associated with `v0.0.0`, is neither draft nor prerelease, and its body matches the reviewed release notes with the exact tagged changelog suffix.
- **SC-003**: The public asset inventory contains exactly thirteen files: six archives, six SBOMs, and one checksum manifest, with zero missing, duplicated, substituted, or unexpected item.
- **SC-004**: All six archive checksums pass, all six archives contain exactly four approved members, all six binaries contain the verified `0.0.0` marker and expected target build identity, the host-compatible binary reports `0.0.0`, and all six packaged schemas match the tagged schema at 100% byte identity.
- **SC-005**: All six published SBOMs identify the correct platform, version, and authorized source revision with zero semantic mismatch.
- **SC-006**: Repository documentation contains zero unreleased-project status claims, zero native-codec capability overclaims, and working links to the public release, immutable tag, and verification evidence.
- **SC-007**: The final S013 pull-request head has zero failed or pending checks, zero unresolved review threads, every review finding addressed, and no more than two Codex review rounds.
- **SC-008**: Issue #27 appears exactly once in `cueson Delivery` with Slice `S013`, Stage `PR review`, and empty default Status while awaiting merge.
- **SC-009**: S013 performs zero pull-request merges, milestone closures, production-domain changes, public schema hosting actions, signatures, attestations, rebuilt asset substitutions, and native codec claims.

## Assumptions

- The operator authorization in this work slice specifically covers creating and pushing unsigned annotated tag `v0.0.0` at `b294a6952c8bd041d852c502f5d7206c0b58edd6` and publishing the reviewed GitHub Release with the exact accepted thirteen-file asset set.
- Milestone closure remains a distinct later authorization even after publication succeeds.
- The accepted workflow artifact remains available long enough to complete publication; expiration requires a new decision rather than an inferred rebuild.
- The existing release notes are the approved public highlights and require no editorial change before publication.
- Post-publication repository documentation is reviewable follow-up state and does not change the already-authorized tag target or release asset bytes.
