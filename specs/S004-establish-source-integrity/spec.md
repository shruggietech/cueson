# Feature Specification: Establish Source Integrity and Restoration

**Feature Branch**: `S004-establish-source-integrity`

**Created**: 2026-09-09

**Status**: Implemented (local verification complete)

**Input**: User description: "Implement issue #6 as the focused S004 source-integrity and timestamp-boundary slice, including exact generic restoration but excluding the broader fixture, CI, codec, release, and production-domain work assigned elsewhere."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Reject an unsafe or corrupted source envelope before output (Priority: P1)

An operator can submit an already-valid Cue JSON document for restoration and receive deterministic rejection before any destination is opened when source bytes, integrity declarations, stored basenames, or the complete destination plan are invalid.

**Why this priority**: Exact restoration is trustworthy only when every asset and destination is proven safe before filesystem mutation begins.

**Independent Test**: Submit valid and deliberately corrupted single-asset and multi-asset documents, then verify that valid documents produce a complete plan while invalid base64, length, hash, basename, portable-name collision, destination collision, or unsafe existing-target cases fail without creating or changing output.

**Acceptance Scenarios**:

1. **Given** an asset whose encoded bytes match its declared length and SHA-256, **When** its source envelope is validated, **Then** the asset is accepted for restoration.
2. **Given** any asset with malformed encoded data, a decoded-length mismatch, or a SHA-256 mismatch, **When** the document is prepared for restoration, **Then** the operation fails before any destination is opened.
3. **Given** two stored basenames that differ textually but collide under the portable canonical-caseless rule, **When** a multi-asset restoration is planned, **Then** the entire plan is rejected before output begins.
4. **Given** a destination that is a symbolic link, directory, or other non-regular entry, **When** restoration is requested with or without force, **Then** the operation refuses to replace it.

---

### User Story 2 - Restore exact source bytes without a codec (Priority: P2)

An operator can use the public restore command to recreate one or every source asset directly from the Cue JSON source envelope, selecting the contractually supported destination mode and explicit overwrite policy without invoking a subtitle codec.

**Why this priority**: Byte-exact restoration is the first executable proof that the preservation envelope is meaningful independently of native format support.

**Independent Test**: Restore representative single-asset and multi-asset documents through each destination mode and compare every resulting byte count and SHA-256 with the source declaration while verifying stdout, stderr, exit status, overwrite refusal, force behavior, and cleanup after failure.

**Acceptance Scenarios**:

1. **Given** a valid single-asset document and no destination option, **When** restore runs, **Then** the exact bytes are written to the stored safe basename in the current directory.
2. **Given** a valid document and an output directory, **When** restore runs, **Then** every asset is written beneath that directory using its stored safe basename.
3. **Given** a single-asset document and an explicit output path, **When** restore runs, **Then** the exact bytes are written to that literal path even when its basename differs from the stored name.
4. **Given** an existing regular destination and no force option, **When** restore runs, **Then** it fails before changing that destination; with force, it replaces only the caller-approved destination.
5. **Given** multiple assets and no output directory, or multiple assets with a single-file output option, **When** restore runs, **Then** it reports a deterministic pre-execution failure without writing output.

---

### User Story 3 - Preserve and report filesystem timestamps truthfully (Priority: P3)

A source-preservation caller can capture observable source timestamps before content reads, and an operator restoring bytes receives an explicit result for each requested timestamp instead of a portable overclaim.

**Why this priority**: Timestamp fidelity matters to archival restoration, but platform and filesystem limits must never be disguised as success.

**Independent Test**: On each supported operating system, capture a source whose timestamps are known, prove capture occurs before the content read, restore it with default, strict, and no-metadata policies, and verify each timestamp result as restored, unsupported, unavailable, or failed at the platform's effective precision.

**Acceptance Scenarios**:

1. **Given** a source file, **When** source preservation captures it, **Then** timestamp values and provenance are obtained before any content read that could update access time.
2. **Given** restorable modification and access timestamps, **When** restore completes with metadata enabled, **Then** the destination is restated and the effective values are verified after byte integrity succeeds.
3. **Given** a captured timestamp the destination cannot reproduce, **When** default restore runs, **Then** exact-byte restoration succeeds with a warning and an unsupported result.
4. **Given** the same unsupported timestamp and strict metadata mode, **When** restore runs, **Then** the operation fails and leaves no partially accepted restored artifact.
5. **Given** no-metadata mode, **When** restore runs, **Then** timestamp application is intentionally skipped without claiming timestamp equality.

### Edge Cases

- Empty and whitespace-only paths, nonexistent inputs, directories supplied as input files, and malformed JSON fail without output.
- Encoded source data with legal syntax but non-canonical trailing bits is rejected rather than accepted as an alternative representation.
- Declared sizes at integer boundaries and encoded payloads whose decoded size cannot be represented safely fail without uncontrolled allocation.
- Duplicate asset identifiers, missing primary assets, unsafe stored basenames, and embedded local paths remain rejected by the existing schema and semantic layers.
- Portable collision checks cover canonical Unicode equivalence, default case folding, Windows device stems with extensions, and destinations that normalize to the same path.
- Existing targets that change between planning and commit never authorize replacement of a non-regular entry and never allow one bundle asset to replace another.
- Write, close, rename, post-write integrity, timestamp application, verification, and cancellation failures clean staged artifacts and report failure deterministically.
- Timestamp values outside the destination platform's representable range are reported as failed or unsupported according to the platform contract and never wrapped or truncated silently.
- `--strict-metadata` and `--no-metadata` are mutually exclusive.
- Quiet and silent modes preserve the existing payload and error rules while filtering only eligible diagnostics.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The restore operation MUST parse and validate the complete Cue JSON document through the existing structural and semantic contract before source-specific validation begins.
- **FR-002**: Every source asset MUST decode from canonical standard base64, and its decoded byte count and lowercase SHA-256 MUST exactly equal the declared values.
- **FR-003**: Source-specific validation MUST complete for every asset before any destination is opened or modified.
- **FR-004**: Every stored basename MUST pass the portable safe-basename rules and MUST be unique within the bundle under NFD normalization, default Unicode case folding, and a final NFD normalization.
- **FR-005**: The complete destination plan MUST be built before output, and no two assets may resolve to the same destination under the portable collision rule or platform path equivalence.
- **FR-006**: A single asset without a destination option MUST target its stored basename in the current directory.
- **FR-007**: An output-directory selection MUST place one or more assets directly beneath that existing directory using stored basenames.
- **FR-008**: A single-file output selection MUST be accepted only for a one-asset document, MUST be a separately validated literal runtime path, and MAY rename the restored asset.
- **FR-009**: The single-file output and output-directory selections MUST be mutually exclusive, and a multi-asset document without an output directory MUST fail before output.
- **FR-010**: Existing destinations MUST be refused unless force is explicit; force MUST apply only to existing regular files and MUST NOT authorize symbolic-link, directory, device, or intra-bundle replacement.
- **FR-011**: Every output MUST be completely staged and byte-verified before becoming an accepted destination, and failed operations MUST clean uncommitted staging artifacts.
- **FR-012**: Restoration MUST recompute byte count and SHA-256 from the staged or written output before reporting exact-byte success.
- **FR-013**: Exact restoration MUST operate directly from the source envelope and MUST NOT require or invoke a native format codec.
- **FR-014**: The public command MUST support `--output`, `--output-dir`, `--force`, `--strict-metadata`, and `--no-metadata`, plus the existing global diagnostic controls and help behavior.
- **FR-015**: Invocation and deterministic pre-execution failures MUST return status 2; parsing, validation, integrity, restoration, metadata, cancellation, and runtime I/O failures MUST return status 1; success MUST return status 0.
- **FR-016**: Successful restoration MUST keep stdout empty and use stderr only for diagnostics; default unsupported metadata MUST produce a warning, while successful exact-byte restoration without warnings MUST be silent.
- **FR-017**: Source capture MUST obtain observable timestamps and their provenance before opening or reading source content, and MUST store only a safe basename rather than an input path.
- **FR-018**: Timestamp capture MUST preserve modification and access instants where observable, MUST identify true creation or birth time separately from a fallback, and MUST never label Unix change time as creation time.
- **FR-019**: Metadata-enabled restoration MUST apply timestamps only after source bytes, byte count, SHA-256, and the final destination name are established, then MUST read back every claimed supported value at the platform's effective precision.
- **FR-020**: Each timestamp restoration result MUST be classified as restored, unsupported, unavailable, or failed, with exact-byte restoration reported separately from metadata fidelity.
- **FR-021**: Default metadata mode MUST warn rather than fail for unsupported values; strict metadata mode MUST fail if any captured requested timestamp is not reproduced and MUST leave no partially accepted artifact; no-metadata mode MUST intentionally skip timestamp application.
- **FR-022**: Windows behavior MUST attempt creation, modification, and access restoration where supported; macOS behavior MUST attempt modification and access restoration and report creation support truthfully; Linux behavior MUST attempt modification and access restoration, capture birth time when available, and report setting birth time as unsupported rather than faking success.
- **FR-023**: Platform-specific timestamp behavior MUST have native-selectable tests for Windows, macOS, and Linux, while execution of the complete hosted native matrix remains assigned to issue #8 so S004 does not create CI workflows.
- **FR-024**: Completing S004 MUST change the official v0.0.0 restoration capability from false to true for recognized envelope-only formats while leaving native ingest, render, and OCR-required capabilities false.
- **FR-025**: Documentation and the unreleased changelog MUST describe exact restoration, metadata limitations, and the remaining issue #7 fixture and issue #8 CI boundaries without claiming native subtitle codecs.

### Key Entities

- **Validated source asset**: One ordered source-envelope asset whose canonical decoded bytes, declared length, declared SHA-256, safe basename, and timestamp observations have passed validation.
- **Destination plan**: The complete ordered mapping from validated assets to caller-selected runtime destinations, including portable collision identities and existing-target decisions.
- **Restoration result**: The per-asset exact-byte outcome, destination, verified length and digest, and separate timestamp result set.
- **Timestamp observation**: A captured instant paired with provenance that distinguishes creation or birth time, platform creation time, fallback observation, and unavailability.
- **Timestamp result**: One requested timestamp's restored, unsupported, unavailable, or failed state and the effective precision used for verification.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: One hundred percent of accepted test assets restore to the declared byte length and SHA-256 without invoking a native format codec.
- **SC-002**: One hundred percent of tested corrupt base64, length, digest, unsafe-name, portable-collision, destination-collision, and non-regular-target cases fail before an unauthorized destination change.
- **SC-003**: Every multi-asset test proves the complete bundle is validated and planned before the first output is opened.
- **SC-004**: Every tested failure path leaves zero uncommitted staging files and strict metadata failures leave zero partially accepted new artifacts or preserve every pre-existing destination.
- **SC-005**: Every observed timestamp in native platform tests receives exactly one truthful restored, unsupported, unavailable, or failed result, with zero creation-time claims derived from Unix change time.
- **SC-006**: The public command's destination modes, overwrite policy, streams, diagnostics, and exit statuses match all documented acceptance scenarios.
- **SC-007**: Software/schema capability checks report restoration support true and ingest, render, and OCR-required support false for both recognized v0.0.0 formats.
- **SC-008**: All S004 tests pass with race detection where supported, native Windows tests pass in the current environment, Linux and macOS platform-selected test sources compile, and the executable builds with native dependencies disabled.
- **SC-009**: Spec Kit analysis and convergence, repository verification, every configured continuous-integration and security check, and both permitted Codex review rounds leave zero unresolved blocking item before the final merge ritual.

## Assumptions

- Issue #6 is the only implementation issue in S004; issue #7 retains broad fixture provenance, golden helpers, malformed corpus conventions, and fuzz infrastructure.
- Hosted Windows, macOS, and Linux matrix execution remains owned by issue #8. S004 supplies platform-selected behavior and tests, runs the current Windows path natively, and cross-compiles the other supported paths without claiming that cross-compilation proves native behavior.
- Destination directories and parent directories must already exist; S004 does not recursively create caller-selected directory trees.
- Existing symbolic links and non-regular destinations are always refused, including under force.
- Concurrent external mutation of source or destination paths is treated as a runtime failure; Cueson validates observable state again at commit boundaries but does not claim cross-process transaction isolation.
- The operator authorized S004 branch publication and an official pull request, plus at most one agent-triggered second Codex review, but retained final merge, tag, release, and production-domain authority.
