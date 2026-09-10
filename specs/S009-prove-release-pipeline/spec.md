# Feature Specification: Non-Publishing Release Proof

**Feature Branch**: `S009-prove-release-pipeline`

**Created**: 2026-09-09

**Status**: Draft

**Input**: User description: "Use Spec Kit and the repository autopilot protocol to prove the complete v0.0.0 release pipeline through a non-publishing dry run, automatically publish the official pull request, complete no more than two Codex review rounds, and stop for the operator's final review and merge ritual."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Produce the Complete Release Candidate Matrix (Priority: P1)

A maintainer can run one documented snapshot process and receive a complete v0.0.0 release-candidate artifact set for every supported operating-system and architecture combination without publishing anything.

**Why this priority**: The release foundation is useful only if it proves that every supported target can be packaged consistently before a public release is attempted.

**Independent Test**: Run the snapshot process from a clean checkout and verify that exactly one correctly named archive is produced for each supported Windows, macOS, and Linux target on amd64 and arm64, with no tag or remote release created.

**Acceptance Scenarios**:

1. **Given** a clean checkout at an identifiable revision, **When** a maintainer runs the documented snapshot process, **Then** it creates six target archives covering Windows, macOS, and Linux on amd64 and arm64.
2. **Given** the complete snapshot output, **When** archive names and formats are inspected, **Then** every name unambiguously identifies the product, version, operating system, and architecture and uses the platform-appropriate binary name.
3. **Given** the snapshot process runs on a pull request or by an explicitly invoked maintainer action, **When** it completes, **Then** it has not created or moved a tag, created a GitHub release, uploaded a release asset, or modified the production schema domain.

---

### User Story 2 - Verify Artifact Integrity and Contents (Priority: P1)

A maintainer can independently prove that every packaged archive is complete, internally consistent, checksum-covered, and free from developer-machine paths.

**Why this priority**: Building archives is insufficient evidence if their contents, schema identity, version relationship, or provenance can silently drift.

**Independent Test**: Inspect all six archives and their checksum manifest, compare the packaged schema with the canonical embedded schema, execute each host-compatible binary, and scan artifact metadata for local path leakage.

**Acceptance Scenarios**:

1. **Given** any produced archive, **When** its contents are listed, **Then** it contains exactly one target-appropriate Cueson executable plus the canonical Cue JSON schema, license, and notice at stable archive paths.
2. **Given** the complete artifact set and checksum manifest, **When** checksums are verified, **Then** every distributable archive has exactly one matching valid checksum entry and no unknown archive is accepted.
3. **Given** the packaged schema and a host-compatible packaged executable, **When** their identities are compared, **Then** the executable version, schema version, schema identifier, and requested release version all agree on `0.0.0`.
4. **Given** artifacts and release metadata from any supported target, **When** they are scanned for workspace, drive, username, and checkout identifiers, **Then** no local machine path or identifier is present.
5. **Given** the generated software-bill-of-materials and build metadata, **When** maintainers inspect them, **Then** each distributable is traceable to its target and source revision without claiming a signature, attestation, or publication that did not occur.

---

### User Story 3 - Reproduce the Proof Locally and in Hosted Automation (Priority: P2)

A maintainer or reviewer can follow one versioned validation contract to reproduce the release proof locally and confirm that hosted automation applies the same acceptance rules.

**Why this priority**: Durable release confidence requires reviewable validation rather than a one-time successful command whose artifact contents were never inspected.

**Independent Test**: Run the repository-owned release verification from a clean local checkout and compare its assertions and results with the pull-request-safe hosted run.

**Acceptance Scenarios**:

1. **Given** the documented prerequisites, **When** a maintainer follows the validation guide, **Then** the snapshot and all artifact assertions complete in the foreground without interactive prompts.
2. **Given** an arbitrary pull request, **When** hosted release verification runs, **Then** it uses read-only repository authority, receives no publication credential, and retains no path capable of publishing a release.
3. **Given** a release proof fails at any build, packaging, checksum, schema, version, or metadata assertion, **When** automation reports the result, **Then** the failure is visible and the run cannot be represented as a successful release candidate.

### Edge Cases

- A target build can succeed while its archive is missing required legal or schema material.
- Windows requires an `.exe` binary while non-Windows archives must not acquire that suffix.
- Archive formats can enumerate paths differently or introduce directory traversal and absolute-path entries.
- A checksum manifest can omit an archive, contain duplicates, include itself, or reference an unexpected file.
- Snapshot tooling can derive a non-release version label from repository state even though the artifact contract being proved is v0.0.0.
- The packaged schema can be byte-different from the canonical repository schema while still declaring the same version.
- Build metadata can contain an absolute checkout path even when archive member names are portable.
- A software-bill-of-materials generator can be unavailable, fail, or produce metadata for the wrong target.
- Pull-request automation can accidentally inherit a write-capable token even when the snapshot command itself does not publish.
- Running from a dirty checkout can make source revision and reproducibility claims ambiguous.
- A host cannot execute binaries built for the other operating systems or architectures, so execution checks must distinguish host-compatible validation from structural validation of the full matrix.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: S009 MUST provide one documented, deterministic snapshot entry point for the v0.0.0 release candidate.
- **FR-002**: The snapshot MUST build with native extensions disabled for Windows, macOS, and Linux on amd64 and arm64.
- **FR-003**: The snapshot MUST produce exactly one distributable archive for each of the six supported targets.
- **FR-004**: Archive names MUST follow `cueson_0.0.0_<goos>_<goarch>.<format>` so they identify Cueson, release version, target operating system, and target architecture without collisions.
- **FR-005**: Windows archives MUST contain `cueson.exe`; macOS and Linux archives MUST contain `cueson` without an executable suffix.
- **FR-006**: Every distributable archive MUST contain the canonical `cueson.schema.json`, `LICENSE`, and `NOTICE` at documented portable paths.
- **FR-007**: Packaged schema bytes MUST exactly match `internal/schema/cueson.schema.json` from the source revision used for the snapshot.
- **FR-008**: The snapshot MUST produce one checksum manifest that covers every distributable archive exactly once and excludes nondistributable working files.
- **FR-009**: Repository-owned verification MUST reject a missing, duplicate, unknown, or mismatched checksum entry.
- **FR-010**: Repository-owned verification MUST reject missing, duplicate, absolute, traversing, or unexpectedly named archive members.
- **FR-011**: Repository-owned verification MUST prove that the executable version, embedded schema version, packaged schema version, schema identifier, and requested release version are in lockstep at `0.0.0`.
- **FR-012**: Host-compatible packaged binaries MUST execute the public version and schema-version surfaces successfully during foreground verification.
- **FR-013**: The snapshot MUST produce a software bill of materials for each distributable target and identify the source revision and target without asserting nonexistent signing or publication.
- **FR-014**: Repository-owned verification MUST scan archive names, members, executable metadata, software bills of materials, and release metadata for local checkout paths and machine identifiers.
- **FR-015**: All release configuration, verification logic, and hosted workflow behavior MUST use pinned inputs, stable semantic assertions, non-interactive execution, and reviewable output. S009 MUST NOT claim byte-identical SBOM reproducibility while the pinned generator emits variable timestamps or document identifiers.
- **FR-016**: Pull-request release verification MUST run with read-only repository permission, without secrets or publication credentials, and without executing untrusted code under a privileged event model.
- **FR-017**: The snapshot and pull-request workflow MUST NOT create or move tags, create GitHub releases, upload release assets, sign artifacts, publish attestations, or modify `cueson.io`.
- **FR-018**: Any future publishing path MUST be separate from the S009 snapshot path and require a separately authorized operator action bound to the intended release.
- **FR-019**: Failures in target build, packaging, contents, checksums, lockstep, software-bill-of-materials generation, or path-leak validation MUST fail the overall proof visibly.
- **FR-020**: The final pull request MUST close issue #11, use publication-safe Markdown, complete no more than two Codex review rounds, and stop for human final review and merge.
- **FR-021**: S009 MUST NOT merge or auto-merge its pull request, create or move a tag, publish a release or schema, or mutate production configuration.

### Key Entities

- **Release Candidate**: The complete non-published v0.0.0 output associated with one source revision and one deterministic snapshot invocation.
- **Target**: One supported operating-system and architecture pair with a required binary name, archive name, and archive format.
- **Distributable Archive**: The portable package containing one target executable, the canonical schema, license, and notice.
- **Checksum Manifest**: The one-to-one digest inventory for all distributable archives.
- **Software Bill of Materials**: Target-specific dependency and build-component metadata tied to an archive and source revision.
- **Release Evidence**: The foreground validation record proving matrix completeness, artifact contents, checksums, lockstep, path hygiene, and non-publication.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: One snapshot invocation produces six and only six distributable archives, covering 100% of the approved operating-system and architecture matrix.
- **SC-002**: 100% of archives contain the expected executable, byte-identical canonical schema, license, and notice, with zero unsafe or unexpected member paths.
- **SC-003**: One checksum manifest verifies 100% of distributable archives with zero missing, duplicate, unknown, or mismatched entries.
- **SC-004**: Executable, embedded schema, packaged schema, and requested release identities agree on `0.0.0` in 100% of target evidence.
- **SC-005**: Every distributable target has one verified software bill of materials tied to the source revision, and zero inspected artifacts claim unsupported signing, attestation, or publication.
- **SC-006**: Automated scans find zero local paths, usernames, drive identifiers, or checkout-specific values in release artifacts and metadata.
- **SC-007**: Local foreground verification and hosted pull-request verification apply the same complete acceptance contract and both fail on every seeded invalid-artifact case.
- **SC-008**: The S009 workflow creates zero tags, GitHub releases, uploaded release assets, published schemas, or production-domain mutations.

## Assumptions

- The supported v0.0.0 release matrix remains Windows, macOS, and Linux on amd64 and arm64.
- `internal/version` remains the single executable build-version source, and the canonical embedded schema remains the schema identity authority.
- Snapshot artifacts are validation outputs and are not public releases, even when hosted automation retains them temporarily for reviewer inspection.
- Artifact signing, cryptographic attestations, tag-triggered publication, GitHub Release creation, and `cueson.io` publication require later specifications and separate operator authorization.
- The release proof may use repository-owned standalone tooling outside the shipped Cueson product when doing so keeps release validation deterministic and testable.
- Issue #12 owns final user documentation and milestone closure; S009 documents only the release validation contract required to prove issue #11.
