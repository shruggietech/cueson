# Feature Specification: Publish and Verify v1.0.0

**Feature Branch**: `codex/S020-publish-v1-release`
**Created**: 2026-09-11
**Status**: Draft
**Input**: Prepare the exact verified Cueson v1.0.0 candidate for an operator-authorized public release, publish it only after that separate authorization, independently verify the public result, and reconcile repository and GitHub delivery records without crossing production or merge boundaries.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Obtain the Authorized Stable Release (Priority: P1)

A Cueson user can open the official v1.0.0 GitHub Release, understand the stable SubRip, WebVTT, Cue JSON, and command-line capabilities, and download the correct archive, checksum manifest, and matching software bill of materials for a supported platform.

**Why this priority**: The public v1 release is the outcome that turns the verified candidate into usable software. An incomplete, mislabeled, or unauthorized publication cannot satisfy the project goal.

**Independent Test**: Inspect the public tag and release, confirm that the annotated tag targets the exact authorized revision, the release is public and final, the approved notes render correctly, and exactly thirteen approved assets are present.

**Acceptance Scenarios**:

1. **Given** the exact candidate and a separate operator authorization naming it, **When** the v1.0.0 tag is inspected locally and remotely, **Then** the annotated tag peels to exactly `2cad4c816340404289b4d1d87179a4071713bb46`.
2. **Given** the approved release notes, **When** the public GitHub Release is opened, **Then** it is neither draft nor prerelease, uses tag `v1.0.0`, presents the stable v1 capabilities and boundaries, and ends with the tagged changelog link.
3. **Given** the accepted S019 candidate artifact, **When** a user inspects release downloads, **Then** exactly six platform archives, six corresponding SPDX JSON SBOMs, and one checksum manifest are available with no internal build metadata or evidence file.

---

### User Story 2 - Independently Verify Publication (Priority: P1)

A maintainer can prove that every public download is the exact accepted candidate artifact and that publication changed no separately protected production, schema-hosting, milestone, or pull-request state.

**Why this priority**: Successful upload calls do not prove that the public bytes, release metadata, or immutable tag are correct. Independent verification is required before v1 can be called released.

**Independent Test**: Download all thirteen public assets into a clean temporary directory, compare their digests with the accepted candidate, inspect every archive and SBOM, execute the host-compatible package, and read protected GitHub state back.

**Acceptance Scenarios**:

1. **Given** the public assets, **When** the checksum manifest is applied, **Then** all six archive hashes match and the manifest contains no missing, duplicated, or unexpected entry.
2. **Given** each published archive, **When** its members and executable are inspected, **Then** it contains exactly the platform executable, `cueson.schema.json`, `LICENSE`, and `NOTICE`; every binary contains the verified `1.0.0` release marker and expected target build identity; the host-compatible binary reports version `1.0.0`; and the schema matches the immutable tagged repository schema byte for byte.
3. **Given** each published SBOM, **When** its stable semantics are inspected, **Then** it identifies the matching platform, version `1.0.0`, and source revision `2cad4c816340404289b4d1d87179a4071713bb46`.
4. **Given** completed publication, **When** protected boundaries are audited, **Then** the v1.0.0 milestone remains open pending separate closure authority, production `cueson.io` remains unchanged, and no signature, attestation, or public schema endpoint was created.

---

### User Story 3 - Record and Close the v1 Delivery State (Priority: P2)

A maintainer or contributor can read the repository and GitHub planning surfaces and see that v1.0.0 is publicly released, where to obtain it, what was verified, and which separately governed capabilities remain deferred.

**Why this priority**: The repository, issue, epic, milestone, and Project must agree with the verified public state so future work starts from a truthful baseline.

**Independent Test**: Review the S020 documentation diff and official pull request, run the complete local and hosted gates, and inspect issue #38, epic #29, milestone v1.0.0, and Project lifecycle data.

**Acceptance Scenarios**:

1. **Given** a verified public release, **When** current repository status documentation is read, **Then** it links to v1.0.0, identifies the exact tagged revision and evidence, and no longer describes v1 as an unpublished candidate.
2. **Given** the v1 capability boundary, **When** release-facing documentation is read, **Then** it describes SubRip, WebVTT, Cue JSON, and CLI support accurately without claiming public schema hosting, signatures, attestations, or external Go API stability.
3. **Given** the official S020 pull request, **When** its delivery record is inspected, **Then** it closes issue #38 on merge, has every current-head check green, resolves every review finding, uses no more than two Codex rounds, and stops for the operator's final merge ritual.
4. **Given** issue #38 is closed by the merged pull request and publication evidence is complete, **When** separately authorized planning reconciliation occurs, **Then** epic #29 and milestone v1.0.0 close only if all their governed work is complete.

### Edge Cases

- The tag or release may already exist because of a partial or external attempt; S020 must stop instead of moving, overwriting, or duplicating public state without fresh operator judgment.
- The accepted workflow artifact expires on 2026-09-14; if it becomes unavailable before publication, S020 must stop rather than silently rebuild or substitute artifacts.
- An uploaded asset can have the correct name but different bytes; public verification must compare content rather than names alone.
- GitHub can create a release successfully while one or more assets fail to upload; partial publication must remain visibly incomplete and block dependent actions.
- The release body can be altered by transport or incorrect source formatting; the published body must be read back and compared structurally with the approved notes.
- An annotated tag has a tag-object identifier distinct from its peeled commit; verification must prove the peeled target rather than comparing only the tag object.
- SBOM timestamps and namespace identifiers can vary between builds, so S020 verifies the exact accepted files and their stable target, version, and source semantics without claiming reproducible SBOM bytes.
- Repository automation can populate the default Project Status field; S020 must clear it and use only the governed Stage field.
- The operator's instruction to kick off S020 authorizes specification and non-mutating preflight, but does not by itself authorize tag creation, tag push, release publication, milestone closure, branch push, pull-request publication, or pull-request merge.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: S020 MUST use GitHub issue #38 as its single actionable issue, preserve milestone v1.0.0 and the v1 epic relationship, represent the issue exactly once in `cueson Delivery`, set Slice `S020`, use the lifecycle-appropriate governed Stage, and leave default Status empty.
- **FR-002**: Publication MUST use tag name `v1.0.0` and the tag MUST peel to exactly `2cad4c816340404289b4d1d87179a4071713bb46`.
- **FR-003**: The tag MUST be an unsigned annotated tag and MUST NOT be created, pushed, moved, or replaced without explicit operator authorization for this exact candidate.
- **FR-004**: The GitHub Release MUST be public, non-draft, non-prerelease, associated with tag `v1.0.0`, named `Cueson v1.0.0`, and use unchanged `docs/releases/v1.0.0.md` content as its body.
- **FR-005**: The public Release MUST contain exactly thirteen assets: six verified platform archives, six corresponding SPDX JSON SBOMs, and `cueson_1.0.0_checksums.txt`.
- **FR-006**: Every public asset MUST come from accepted default-branch release-proof run `34621429626`, artifact `10273380044`; S020 MUST NOT rebuild, regenerate, rename, or substitute any release asset.
- **FR-007**: The public Release MUST NOT include `release-evidence.json`, `artifacts.json`, `metadata.json`, `config.yaml`, extracted binary directories, or any unexpected file.
- **FR-008**: Before publication, S020 MUST prove that the accepted artifact is still available, its evidence names the exact candidate revision and schema digest, the release notes pass the GitHub formatter unchanged, and neither a local nor remote tag nor a GitHub Release conflicts with `v1.0.0`.
- **FR-009**: Post-publication verification MUST read back the remote tag reference, peel the annotated tag, and prove its exact authorized commit.
- **FR-010**: Post-publication verification MUST read back the release tag, target, name, draft state, prerelease state, body, and complete asset metadata directly from GitHub.
- **FR-011**: Post-publication verification MUST download all thirteen public assets into a new clean location, compare every SHA-256 digest with the accepted candidate contract, and verify the six-entry archive checksum bijection.
- **FR-012**: Post-publication verification MUST prove every archive has exactly four safe regular root members, uses the correct platform executable name, contains a binary with the verified `1.0.0` release marker and expected target build identity, and packages schema bytes matching the tagged immutable repository schema with SHA-256 `1aad14567033d7e14d9beb78985e18007aefb5345095370b11b6b887df7ec541`; the host-compatible public binary MUST report version `1.0.0` at runtime.
- **FR-013**: Post-publication verification MUST prove each SBOM's stable semantics identify its corresponding platform, version `1.0.0`, and source revision `2cad4c816340404289b4d1d87179a4071713bb46`.
- **FR-014**: Repository documentation MUST identify v1.0.0 as publicly released, link the release and immutable tag, record the exact verified revision and evidence, retain truthful capability boundaries, and leave the tagged release notes and immutable versioned schema unchanged.
- **FR-015**: Repository-authored text MUST use UTF-8 without BOM, follow `.gitattributes`, avoid mojibake, and retain one source line per Markdown paragraph or list item.
- **FR-016**: The official S020 pull request MUST contain `Closes #38`, pass the complete local and hosted verification suite, address every Codex, security, and review finding, request at most one second Codex round when required, and never request a third automatically.
- **FR-017**: S020 MUST stop for the operator's final review and merge ritual and MUST NOT merge or auto-merge its pull request.
- **FR-018**: S020 MUST NOT close epic #29 or milestone v1.0.0, publish a schema endpoint, mutate production `cueson.io`, add signatures or attestations, create or publish the tag or release, push its branch, or publish its pull request without the specific authority required for each protected action.

### Key Entities

- **Authorized Candidate**: Exact default-branch commit `2cad4c816340404289b4d1d87179a4071713bb46` and accepted release-proof artifact `10273380044` from run `34621429626`.
- **Release Tag**: Immutable annotated `v1.0.0` reference whose peeled target identifies the Authorized Candidate.
- **GitHub Release**: Public, non-draft, non-prerelease release record associated with the tag, approved notes, and exactly thirteen user-facing assets.
- **Public Asset Set**: Exact six archives, six target-bound SBOMs, and one checksum manifest exposed to users.
- **Publication Evidence**: GitHub read-back and independent download verification proving public state and bytes match the Authorized Candidate.
- **Release Documentation Record**: Repository-authored status, process, verification, schema, and architecture text updated after successful publication.
- **Delivery Reconciliation**: Issue #38, epic #29, milestone v1.0.0, dependencies, and Project stages aligned only after evidence and applicable operator authority.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: One immutable annotated tag named `v1.0.0` peels to the authorized 40-character commit with zero target mismatch.
- **SC-002**: One public GitHub Release is associated with `v1.0.0`, is neither draft nor prerelease, and its body matches the approved release notes with the exact tagged changelog suffix.
- **SC-003**: The public asset inventory contains exactly thirteen files: six archives, six SBOMs, and one checksum manifest, with zero missing, duplicated, substituted, or unexpected item.
- **SC-004**: All thirteen public digests match the accepted candidate, all six archive checksums pass, all six archives contain exactly four approved members, the host-compatible binary reports `1.0.0`, and all six packaged schemas match the tagged schema at 100% byte identity.
- **SC-005**: All six published SBOMs identify the correct platform, version, and authorized source revision with zero semantic mismatch.
- **SC-006**: Repository documentation contains zero current unpublished-v1 claims, zero capability overclaims, and working links to the public release, immutable tag, and verification evidence.
- **SC-007**: The final S020 pull-request head has zero failed or pending checks, zero unresolved review threads, every review finding addressed, and no more than two Codex review rounds.
- **SC-008**: Issue #38 appears exactly once in `cueson Delivery` with Slice `S020`, lifecycle-appropriate Stage, and empty default Status throughout delivery.
- **SC-009**: Before separate publication authority, S020 performs zero tag, release, branch-push, pull-request-publication, milestone-closure, epic-closure, production-domain, public-schema-hosting, signature, attestation, or merge mutations.

## Assumptions

- A later operator authorization must identify the exact candidate commit and accepted artifact before the tag and GitHub Release may be created or published.
- Branch push and pull-request publication require explicit authority in a later instruction; final pull-request merge always remains a separate single-use human authorization.
- Milestone and epic closure remain distinct protected reconciliation actions after the pull request closes issue #38 and all publication evidence is complete.
- The accepted workflow artifact remains available long enough to complete publication; expiration requires a new decision rather than an inferred rebuild.
- The existing release notes are the approved public highlights and require no editorial change before publication.
- Post-publication repository documentation is reviewable follow-up state and does not change the already-authorized tag target or release asset bytes.
