# Feature Specification: Build Schema Foundation

**Feature Branch**: `S003-build-schema-foundation`

**Created**: 2026-09-09

**Status**: Implemented (local verification complete)

**Input**: User description: "Use Spec Kit to define and complete work slice S003 as the canonical v0.0.0 Cue JSON schema foundation from issue #5, then publish it for bounded automated review."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Validate a portable Cue JSON document (Priority: P1)

As a Cue JSON producer or consumer, I can use one canonical contract to validate a document that exposes normalized cues and preserves a portable multi-asset source envelope.

**Why this priority**: Every later restoration and format-codec slice depends on a stable, machine-checkable document contract.

**Independent Test**: Validate the canonical representative document and a set of deliberately invalid variants against the canonical schema and semantic rules, confirming that only the conforming document succeeds.

**Acceptance Scenarios**:

1. **Given** the canonical v0.0.0 schema and representative Cue JSON document, **When** the document is validated, **Then** Draft 2020-12 structural validation and Cueson semantic validation both succeed.
2. **Given** a document whose cue omits timing, raw text, plain text, logical lines, or OCR observations, **When** the document is validated, **Then** validation fails with no loss or silent correction.
3. **Given** a source envelope with multiple ordered assets, **When** the document is validated, **Then** one asset is identified as primary and every asset retains its identifier, role, safe basename, media description, exact-size declaration, lowercase SHA-256 declaration, standard-base64 source data, timestamps, and encoding observation.
4. **Given** a document containing a filesystem path, unsafe basename, malformed digest, malformed base64, unresolved primary asset, inconsistent cue duration, duplicate identity, or inconsistent source order, **When** validation runs, **Then** the applicable structural or semantic layer rejects it deterministically.

---

### User Story 2 - Retrieve the exact embedded schema (Priority: P2)

As a command-line user, I can print the executable's canonical schema, report its version, or save it to an explicitly selected path without an implicit overwrite.

**Why this priority**: An embedded schema makes the contract available offline and proves that the executable and schema ship together.

**Independent Test**: Invoke every supported `schema` form and compare its payload with the canonical embedded artifact, including overwrite and output-failure cases.

**Acceptance Scenarios**:

1. **Given** the S003 executable, **When** a user invokes `cueson schema`, **Then** stdout contains the canonical schema bytes, stderr is empty, and the process exits 0.
2. **Given** the S003 executable, **When** a user invokes `cueson schema --version`, **Then** stdout is exactly `0.0.0` followed by one LF and the process exits 0.
3. **Given** a nonexistent destination, **When** a user invokes `cueson schema --output PATH`, **Then** the canonical schema is written as UTF-8 without BOM and LF line endings, no payload is written to stdout, and the process exits 0.
4. **Given** an existing destination, **When** a user omits `--force`, **Then** the command refuses replacement before accepting the output operation and exits 2; **When** the user supplies `--force`, **Then** the destination is replaced with the exact canonical schema or the prior destination remains intact if replacement fails.

---

### User Story 3 - Detect contract drift before release (Priority: P3)

As a maintainer, I can rely on automated checks to reject schema/software version drift, nonconforming project-owned key names, and capability claims that exceed the shipped executable.

**Why this priority**: Lockstep and naming checks prevent accidental incompatible artifacts and false support claims from becoming release inputs.

**Independent Test**: Run the schema conformance and version-lockstep tests against valid and mutated artifacts and confirm each prohibited drift fails.

**Acceptance Scenarios**:

1. **Given** the executable and embedded schema, **When** their version identities are compared, **Then** the check succeeds only when both are exactly `0.0.0`.
2. **Given** a project-owned schema property or enum value that is not lowercase `snake_case`, **When** conformance checks inspect the canonical artifact, **Then** the check fails while standards-defined JSON Schema keywords remain exempt.
3. **Given** S003 has no public restoration command and no native subtitle codec, **When** official format capability is validated, **Then** `subrip` and `webvtt` are schema-recognized and `envelope_only`, while ingest, render, restore, and OCR-required flags all remain false.

### Edge Cases

- Empty asset and cue collections are rejected where the fixed contract requires content, while empty observation arrays remain valid.
- Unknown root or fixed-object properties are rejected; explicitly namespaced extension maps remain the only intentional extensibility boundary.
- Safe basenames reject separators, traversal tokens, control characters, Windows-invalid punctuation and trailing characters, URI or drive syntax, and case-insensitive reserved device stems including extensions.
- Asset and cue identifiers must be unique within their document-local collections, and the selected primary asset must exist and carry the primary role.
- Cue timing uses non-negative integer milliseconds, requires end not before start, and requires duration to equal end minus start.
- Timestamp pairs must describe one instant, and unavailable creation provenance must not accompany a claimed creation timestamp.
- Schema output failures, canceled execution, and malformed invocation never produce a partial stdout payload or a falsely successful status.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The project MUST provide one canonical Cue JSON schema identified as `https://cueson.io/schema/v0.0.0/cueson.schema.json` and using JSON Schema Draft 2020-12.
- **FR-002**: The canonical schema MUST define a fixed root containing `$schema`, `schema_version`, `format`, `format_support`, `producer`, `source`, `metadata`, `document`, `cues`, `format_data`, `diagnostics`, and `stats`, with unknown properties rejected at stable object boundaries.
- **FR-003**: Cueson-owned property names and enum values MUST use lowercase `snake_case`; standards-defined JSON Schema keywords are exempt.
- **FR-004**: The canonical format identifiers MUST be `subrip` and `webvtt`, and each document's format-specific cue and root data MUST use exactly the key matching its root format.
- **FR-005**: Official S003 capability data for both canonical formats MUST use status `envelope_only` and MUST declare ingest, render, restore, and OCR-required capability false.
- **FR-006**: Every common cue MUST expose a document-local identity, ordinal, total source order, nullable source identifier, non-negative integer-millisecond timing, raw text, plain text, ordered logical lines, speakers, tokens, required OCR observations, nullable placement, and matching format-specific data.
- **FR-007**: Every non-empty OCR observation MUST be independently identified, marked as derived, identify an OCR engine and source asset, and permit optional engine version, model, language, text, lines, confidence, alternatives, source event or image identity, regions, and namespaced processing options without replacing source truth.
- **FR-008**: The source envelope MUST contain a resolvable primary asset identifier and one or more ordered assets, supporting both primary and companion roles from the first schema.
- **FR-009**: Every source asset MUST declare a document-local identifier, role, portable safe basename, nullable media type, exact byte count, optional display size, lowercase SHA-256, standard-base64 original bytes, timestamp observations and provenance, and nullable encoding observations.
- **FR-010**: Structural validation MUST reject filesystem paths and unsafe basenames, invalid version or format values, malformed SHA-256 or base64 syntax, missing common cue fields, and mismatched format-specific object shapes.
- **FR-011**: Semantic validation MUST reject unresolved or multiply designated primary assets, duplicate asset or cue identities, inconsistent cue duration, duplicate or non-monotonic source order, mismatched document and summary counts, inconsistent timestamp pairs or provenance, and capability data inconsistent with S003.
- **FR-012**: Decoded-byte length, digest verification, portable Unicode basename collision enforcement, source timestamp capture, exact restoration, and restoration metadata application MUST remain assigned to source-foundation issue #6 and MUST NOT be claimed as shipped S003 operations.
- **FR-013**: The repository MUST include at least one representative v0.0.0 document that passes both structural and semantic validation and exercises the common cue and multi-asset-capable source envelope.
- **FR-014**: The executable MUST embed the canonical schema bytes and expose internal access to its bytes, identifier, version, structural validation, semantic validation, and software/schema lockstep checks without creating a public Go API.
- **FR-015**: `cueson schema` MUST emit the exact embedded canonical schema to stdout, followed by no commentary, and return exit code 0 when the write succeeds.
- **FR-016**: `cueson schema --version` MUST emit exactly `0.0.0` followed by one LF to stdout and return exit code 0.
- **FR-017**: `cueson schema --output PATH` MUST write the exact canonical schema to the literal selected path without a stdout payload, MUST refuse existing output unless `--force` is supplied, and MUST preserve the prior destination if a forced replacement cannot complete.
- **FR-018**: Schema help MUST document only implemented schema behavior, and invalid schema option combinations or failed overwrite preconditions MUST produce exit code 2 with stderr usage while runtime validation or I/O failures MUST produce exit code 1.
- **FR-019**: Focused tests MUST cover metaschema compilation, representative validation, structural rejection, semantic invariants, schema-key conformance, capability truth, embedded-byte identity, version lockstep, schema CLI output, help, overwrite behavior, and write failures.
- **FR-020**: S003 MUST NOT register restoration, native SRT or WebVTT ingest, model-driven render, conversion, validation, inspection, completion, release, or public-domain capabilities.
- **FR-021**: Repository documentation and the unreleased changelog MUST accurately describe the available schema foundation and its deferred source-integrity and codec boundaries.
- **FR-022**: The completed implementation MUST remain portable and pure Go with native dependencies disabled.

### Key Entities

- **Cue JSON document**: The versioned root contract tying official format capability, producer identity, source assets, common semantic cues, native fidelity, diagnostics, and derived statistics together.
- **Format support**: The release-specific maturity and boolean capability declaration for the document's canonical source format.
- **Source envelope and asset**: The ordered, path-free preservation declaration containing a primary reference and each asset's portable identity, integrity metadata, timestamp observations, encoding observations, and base64 source data.
- **Common cue**: A normalized timed semantic unit containing native-preserving text, consumer-friendly text, logical lines, derived observations, and exactly one matching format-specific extension object.
- **OCR observation**: A derived, independently identified and provenanced semantic observation linked to a source asset and optionally to a source event, source image, and regions.
- **Diagnostic and statistics summary**: Deterministic document data that surfaces defects and summarizes cue, timing, and diagnostic counts.
- **Embedded schema**: The byte-identical canonical contract carried by the executable with its own identity and version, checked against the software version.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The canonical schema compiles under Draft 2020-12 and validates 100 percent of the checked-in representative document set.
- **SC-002**: Every tested invalid variant for required shape, unsafe basename, digest syntax, base64 syntax, format matching, primary resolution, cue timing, ordering, counts, timestamps, naming, and capability truth is rejected by its assigned validation layer.
- **SC-003**: One hundred percent of successful schema-output invocations return bytes identical to the embedded canonical artifact, and schema-version invocations return the exact 6-byte payload `0.0.0\n`.
- **SC-004**: Version-lockstep checks report exactly zero difference among the executable version, embedded schema version, schema identifier version, and canonical instance version.
- **SC-005**: Root and schema help expose exactly the two shipped commands, `version` and `schema`, and expose zero deferred commands or unavailable capability claims.
- **SC-006**: Every schema and CLI test passes with and without race detection where supported, and a clean build succeeds with native dependencies disabled.
- **SC-007**: The Spec Kit checklist, analysis, convergence, repository verification, continuous integration, and both authorized automated review rounds report zero unresolved blocking item before the final merge ritual.

## Assumptions

- Issue #5 is the only implementation issue in S003; issue #6 retains source byte-integrity execution, portable collision keys, capture ordering, and exact restoration.
- Structural basename rules in the schema cover portable single-name safety; issue #6 adds the Unicode canonical caseless collision plan before any destination is opened.
- The representative document is externally authored contract evidence because S003 deliberately provides no native subtitle encoder.
- Schema output uses a durable replacement strategy so `--force` does not destroy a previously valid destination when the replacement write fails.
- The canonical schema URI remains an identifier only; activating `cueson.io` or publishing a release artifact is outside this slice.
- The operator has explicitly authorized this slice's push, pull-request publication, and at most one second Codex review request, but has retained final merge authority.
