# Feature Specification: Prepare the v1.0.0 Release Candidate

**Feature Branch**: `codex/S019-prepare-v1-release-candidate`

**Created**: 2026-09-11

**Status**: Draft

**Input**: User description: "Use Spec Kit and the repository autopilot protocol to prepare and prove the complete v1.0.0 release candidate, push the changes, publish the official pull request, address every review through no more than two Codex rounds, and stop only when the pull request is green and fully reviewed."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Review the Stable v1 Release Record (Priority: P1)

A maintainer can review one complete release record in which the executable identity, Cue JSON schema identity, immutable schema copy, detailed history, concise release notes, and stable capability declarations all agree on v1.0.0.

**Why this priority**: Cueson cannot make a trustworthy first stable release while its permanent contract bytes and public capability record remain development-only or inconsistent.

**Independent Test**: Compare all repository and emitted version surfaces, compare the canonical and versioned schemas byte for byte, inspect the dated changelog and release notes, and confirm that every maintained compatibility statement describes the implemented stable v1 contract.

**Acceptance Scenarios**:

1. **Given** the prepared v1 candidate, **When** a maintainer checks executable, schema, and intended tag identities, **Then** every surface reports `1.0.0` and the immutable schema copy matches the canonical schema byte for byte.
2. **Given** the complete implemented v1 history, **When** a maintainer reads the changelog and release notes, **Then** the changelog provides the detailed dated record, the notes provide highlights only, and the notes end with the exact tagged changelog link.
3. **Given** the implemented SubRip, WebVTT, conversion, validation, inspection, restoration, and completion workflows, **When** maintained documentation is inspected, **Then** those v1 capabilities are described consistently as stable without claiming unimplemented production hosting, signatures, or attestations.

---

### User Story 2 - Prove the Exact v1 Candidate (Priority: P1)

An operator can identify one exact clean source revision and audit a complete non-publishing candidate bundle proving all supported archives, checksums, schemas, executable identities, software bills of materials, and source provenance for that revision.

**Why this priority**: Publication authority must apply to a concrete, reproducible candidate rather than a moving branch, a development snapshot, or an informal green-check claim.

**Independent Test**: Build the supported six-target matrix from a clean exact commit and verify the resulting evidence against that commit, version `1.0.0`, the immutable schema digest, the archive and checksum inventory, target-bound software bills of materials, native host execution, and the non-published state.

**Acceptance Scenarios**:

1. **Given** a clean source revision, **When** candidate proof runs, **Then** it accepts exactly six documented archives, six matching software bills of materials, six checksum entries, and one byte-identical schema identity across repository, binary, and archive surfaces.
2. **Given** any missing, additional, duplicated, unsafe, version-mismatched, source-mismatched, or byte-different candidate member, **When** verification runs, **Then** the candidate is rejected and no accepted evidence is produced.
3. **Given** a host-compatible archive, **When** the packaged executable is smoke-tested, **Then** its version and schema commands report `1.0.0` and its emitted schema matches the immutable release schema byte for byte.
4. **Given** a later squash merge, **When** the same proof runs on the resulting default-branch commit, **Then** that new commit receives its own evidence and no pull-request head is predicted as the release target.

---

### User Story 3 - Hand Off a Governed v1 Publication Decision (Priority: P2)

The operator receives a green, fully reviewed pull request that closes the release-candidate issue and leaves the exact post-merge tag and release decision as the only remaining v1 delivery step.

**Why this priority**: The fastest safe path to v1 is to finish every candidate gate now while preserving the explicit authority boundary around public publication.

**Independent Test**: Inspect the issue, Delivery Project item, milestone, official pull request, current-head checks, review record, and repository state, then confirm the candidate is ready to merge while no tag, release, production schema, milestone closure, or production configuration was changed.

**Acceptance Scenarios**:

1. **Given** issue #37, **When** its delivery metadata is inspected, **Then** it appears exactly once in `cueson Delivery`, uses Slice `S019`, advances through the governed Stage values, leaves default Status empty, remains a child of the v1 epic, and continues to block publication issue #38 until merge.
2. **Given** the official S019 pull request, **When** its body and review state are inspected, **Then** it contains `Closes #37`, every finding has a recorded disposition, no more than two Codex rounds occurred, and every current-head check is green.
3. **Given** a green and fully reviewed S019 pull request, **When** autopilot reaches its halt, **Then** merge, tag creation, GitHub Release publication, milestone closure, production schema hosting, and production configuration remain untouched pending specific operator authorization.

### Edge Cases

- A schema copy can be semantically equivalent but byte-different because of whitespace, encoding, line endings, or key ordering.
- A release note can be well formed while omitting the exact tagged changelog URL, duplicating the changelog, or claiming production hosting and signing that do not exist.
- A development-mode switch can accidentally remain active and permit a candidate without the immutable v1 schema.
- A pull-request head cannot safely be recorded as the final release target because the configured squash merge creates a different default-branch commit.
- A workflow can pass on the pull request while no proof exists for the eventual default-branch merge commit.
- A dirty working tree can produce artifacts whose build metadata does not match a reviewable source commit.
- Build-tool output can contain an incomplete, unexpected, duplicated, source-mismatched, or target-mismatched artifact catalog.
- Software-bill-of-material timestamps and document identifiers can differ across equivalent builds even when their semantic target and source binding is correct.
- Version promotion can leave stale `0.1.0` development claims in maintained documentation, examples, tests, workflow names, or artifact paths.
- Candidate preparation can accidentally add a tag trigger, write permission, credential, signing claim, production deployment, or hidden publication path.
- Review remediation can change the pull-request head after hosted checks or automated conclusions were recorded.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The executable version, embedded schema version, canonical schema identifier, immutable release-schema copy, packaged schema, and intended tag identity MUST all equal `1.0.0`.
- **FR-002**: The immutable `1.0.0` release schema MUST be byte-identical to the canonical embedded schema and MUST be admitted without altering any previously released schema.
- **FR-003**: Current format capability declarations MUST change from development `experimental` to stable v1 only where the implemented and tested behavior supports that claim.
- **FR-004**: The detailed accumulated v1 history and architecture decisions MUST move into a dated `1.0.0` changelog section beneath a fresh empty `[Unreleased]` section.
- **FR-005**: Candidate release notes MUST contain concise highlights only and MUST end exactly with `Full changelog: https://github.com/shruggietech/cueson/blob/v1.0.0/CHANGELOG.md`.
- **FR-006**: Candidate construction MUST produce exactly six documented Windows, macOS, and Linux archives for amd64 and arm64, each containing only the target executable, immutable v1 schema, license, and notice.
- **FR-007**: Candidate construction MUST produce one checksum entry and one target-bound software bill of materials for each archive, and every archive MUST contain byte-identical repository `LICENSE` and `NOTICE` material with the approved non-executable mode.
- **FR-008**: Verification MUST reject missing, additional, duplicated, unsafe, version-mismatched, source-mismatched, dirty, mode-different, or byte-different candidate content before accepting evidence.
- **FR-009**: Accepted evidence MUST record the exact full source revision, version, intended tag, immutable-schema digest, repository legal-file digests, complete archive and software-bill-of-material inventory and digests, checksum count, compatible-host execution result, and `published: false`; candidate verification MUST run with development mode disabled.
- **FR-010**: Candidate verification MUST execute the same accepted candidate bundle on available native Linux, Windows, and macOS runners and MUST verify exact version output, exact schema-version output, and byte-identical emitted schema without pretending that unsupported foreign or arm64 binaries executed.
- **FR-011**: The same non-publishing proof MUST run on pull requests, explicit manual dispatch, and pushes to `main` so the eventual squash-merge commit receives independent candidate evidence.
- **FR-012**: The checked-in release path MUST remain read-only, secretless, credential-free after checkout, publication-disabled, and without tag triggers, identity-token authority, signing, deployment, or production mutation.
- **FR-013**: Maintained documentation MUST describe the actual v1 candidate commands, stable format support, release boundaries, artifact matrix, verification procedure, and known limitations without stale development claims.
- **FR-014**: Every v1 implementation, compatibility, documentation, security, review, Project, milestone, and repository-control gate MUST have current evidence before the candidate is considered ready.
- **FR-015**: A clean checkout of the exact candidate commit MUST reproduce the complete candidate proof without publishing or mutating protected state.
- **FR-016**: Repository-authored text changed by S019 MUST use UTF-8 without BOM, follow `.gitattributes`, avoid mojibake, and use one source line per Markdown paragraph or list item.
- **FR-017**: S019 MUST address every review and security finding, resolve threads only after the concern is handled, request at most one second Codex review when required, and never request a third round automatically.
- **FR-018**: S019 MUST stop for the operator's final review and merge ritual and MUST NOT merge or auto-merge the pull request, create or move a tag, publish a GitHub Release or release asset, close the v1 milestone, publish a schema, or mutate production configuration.

### Key Entities

- **Stable Release Identity**: The single `1.0.0` identity shared by the executable, schema, immutable repository copy, packaged artifacts, evidence, release notes, and intended tag.
- **Immutable v1 Schema**: The versioned repository copy whose bytes match the canonical embedded schema and every packaged copy.
- **Candidate Revision**: A full clean Git commit associated with one verified artifact set; after squash merge, the exact `main` merge commit becomes the publication candidate only after its own proof succeeds.
- **Candidate Evidence**: The deterministic verification record binding version, candidate mode, schema digest, archives, checksums, software bills of materials, host execution, publication state, and source revision.
- **Candidate Release Notes**: The concise, publication-ready highlights document with the exact tagged changelog link and truthful boundaries.
- **Publication Decision Package**: The reviewed repository and hosted evidence presented to the operator before any protected tag or release action.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The canonical, immutable repository, embedded binary, emitted, and six packaged schema copies achieve 100% byte identity and all report version `1.0.0`.
- **SC-002**: The verifier accepts exactly six archives, six software bills of materials, six checksum entries, and four members per archive, with zero missing, additional, duplicated, unsafe, dirty, or revision-mismatched items.
- **SC-003**: Accepted candidate evidence contains one 40-character lowercase source revision, intended tag `v1.0.0`, three 64-character lowercase digests for the immutable schema, license, and notice, `published: false`, one native host execution identity when requested, and complete evidence for all six targets; the candidate proof does not use development mode.
- **SC-004**: The changelog contains one dated `1.0.0` section and one fresh empty `[Unreleased]` section, while the release notes end with the exact required link and remain substantially shorter than the detailed changelog.
- **SC-005**: Maintained current-version and format-capability assertions contain zero stale `0.1.0` development identities or `experimental` labels where stable v1 behavior is being declared.
- **SC-006**: A clean exact commit reproduces the full six-target candidate proof within each hosted job's 20-minute timeout, smoke-tests the accepted bundle on Linux, Windows, and macOS, and performs zero protected publication actions.
- **SC-007**: The final pull-request head has zero failed or pending checks, zero unresolved review threads, every finding addressed, and no more than two Codex review rounds.
- **SC-008**: The S019 pull request contains `Closes #37`; the issue Project entry uses Slice `S019`, Stage `PR review`, and an empty default Status while awaiting merge.
- **SC-009**: S019 performs zero merges, tag creations or movements, GitHub Release or asset publications, milestone closures, production schema publications, and production configuration changes.

## Assumptions

- Version `1.0.0`, intended tag `v1.0.0`, release date `2026-09-11`, the six-target matrix, the four-member archive contract, and the currently pinned release tools remain approved for candidate preparation.
- S018 froze the intended v1 CLI and Cue JSON behavior, so S019 promotes the reviewed contract identity and capability status without adding unrelated product features.
- The final publication target cannot be known until S019 is squash-merged; the resulting `main` commit becomes eligible only after its own default-branch proof succeeds.
- Equivalent software-bill-of-material documents are verified for stable semantic identity and source binding rather than byte-for-byte reproducibility.
- Push and pull-request authorization in this request does not authorize merge, tag creation, GitHub Release publication, milestone closure, schema hosting, or production changes.
