# Feature Specification: Establish Fixture and Conformance-Test Infrastructure

**Feature Branch**: `S005-establish-test-foundation`

**Created**: 2026-09-09

**Status**: Implemented (local verification complete)

**Input**: User description: "Deliver issue #7 as S005 by establishing byte-stable fixture provenance, reusable golden and conformance helpers, path-leak prevention, deterministic malformed cases, and initial fuzz boundaries without claiming native subtitle grammar coverage."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Add trustworthy fixtures (Priority: P1)

A maintainer can add a small test fixture with explicit provenance, redistribution status, expected byte identity, and intentional text characteristics so reviewers can determine exactly what the repository contains and why.

**Why this priority**: Every later parser, renderer, restoration, and conversion claim depends on trustworthy source material whose bytes cannot drift silently between checkouts.

**Independent Test**: Add valid and deliberately invalid fixture records, then verify that complete records with matching files and digests pass while missing, duplicated, unsafe, unlicensed, or byte-drifted records fail deterministically.

**Acceptance Scenarios**:

1. **Given** a redistributable fixture and a complete provenance record, **When** fixture verification runs, **Then** every declared file exists beneath the fixture root and matches its recorded byte length and SHA-256.
2. **Given** a fixture with intentional non-UTF-8 bytes or mixed line endings, **When** it is checked out and verified, **Then** its authoritative bytes remain unchanged rather than being normalized as repository text.
3. **Given** a missing redistribution decision, local-machine path, duplicate fixture identifier, unsafe relative path, unlisted file, or digest mismatch, **When** fixture verification runs, **Then** the repository test fails with a deterministic diagnostic identifying the record and violated rule.

---

### User Story 2 - Reuse one conformance vocabulary (Priority: P2)

A codec contributor can use shared golden-test helpers to compare the common model, ordered diagnostics, rendered output, exact bytes and hashes, and platform-supported timestamp results without inventing incompatible assertions for each format.

**Why this priority**: Consistent comparisons are required before SRT and WebVTT development can make compatible, reviewable fidelity claims.

**Independent Test**: Exercise every comparison surface with one passing seed and one intentional mismatch, verify stable failure detail, and prove generated expectations contain no original filesystem path or local machine identifier.

**Acceptance Scenarios**:

1. **Given** matching expected and observed models, diagnostics, rendered output, source bytes, hashes, and supported timestamps, **When** the shared comparisons run, **Then** every requested comparison passes without depending on an absolute checkout path.
2. **Given** a mismatch on any supported comparison surface, **When** the shared comparison runs, **Then** it fails at that surface with deterministic evidence that does not expose a local path.
3. **Given** expectations produced from equivalent fixtures located in different temporary directories, **When** their portable representations are compared, **Then** they are identical and contain neither source directory.

---

### User Story 3 - Retain malformed and fuzz regressions (Priority: P3)

A maintainer can promote malformed-input and fuzz discoveries into small permanent regression cases that remain deterministic, panic-safe, path-free, and clearly bounded to implemented validation behavior.

**Why this priority**: The initial corpus must harden schema and source-envelope boundaries now while avoiding premature claims about unimplemented SRT or WebVTT grammars.

**Independent Test**: Run all malformed cases twice and execute bounded fuzz smoke runs from the accepted seeds, then verify stable outcomes, zero panics, no filesystem mutation, and no native codec coverage claim.

**Acceptance Scenarios**:

1. **Given** a registered malformed Cue JSON or source-envelope case, **When** the regression harness runs repeatedly, **Then** it returns the same rejection class and stable diagnostic fragment without panic.
2. **Given** valid and malformed seed inputs, **When** each initial fuzz boundary executes for a bounded foreground interval, **Then** it either accepts or rejects safely without panic, unbounded allocation, destination mutation, or local-path output.
3. **Given** a new fuzz-found regression, **When** it is committed permanently, **Then** it is represented by a manifest-governed minimal fixture with provenance and a deterministic expected outcome.

### Edge Cases

- Empty manifests, empty fixture identifiers, duplicate identifiers, duplicate declared paths, and case- or normalization-colliding paths fail deterministically.
- Fixture paths are portable relative paths beneath the governed testdata root; absolute paths, traversal, drive prefixes, URI-like prefixes, backslash separators, and symbolic links are rejected.
- Every regular fixture file is declared exactly once; directories and repository bookkeeping files are not interpreted as fixture payloads.
- Intentionally empty files are valid when their zero length and empty-content digest are declared.
- Non-UTF-8, byte-order marks, CRLF, bare carriage returns, mixed line endings, and missing final newlines remain authoritative bytes when explicitly described.
- Provenance records distinguish original, generated, synthetic, and derived material and require an explicit redistribution decision rather than inferring permission from origin.
- Timestamp comparisons require only observations the active platform reports as restored; unsupported and unavailable observations remain explicit rather than being treated as equality failures.
- Fuzz targets accept arbitrary bytes, return normally for rejected input, and never write restoration destinations.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The repository MUST define one governed fixture root with documented directories for valid seeds, malformed regressions, expected artifacts, and future format-specific corpora.
- **FR-002**: Every committed fixture payload MUST be represented exactly once in a machine-verifiable provenance manifest with a stable identifier, portable relative path, purpose, origin classification, source reference or synthetic recipe, redistribution decision, byte length, and lowercase SHA-256.
- **FR-003**: Fixture verification MUST reject missing or extra payload files, duplicate identifiers or paths, unsafe or non-portable paths, symbolic links, non-regular payloads, incomplete provenance, and integrity mismatches.
- **FR-004**: Byte-sensitive fixture payloads MUST be excluded from checkout text normalization and whitespace rewriting, including fixtures that happen to contain valid UTF-8 text.
- **FR-005**: Repository-authored manifest and expected text MUST remain UTF-8 without a byte-order mark and use the repository's governed line endings unless a manifest record explicitly classifies the file as an intentional byte fixture.
- **FR-006**: Fixture verification MUST preserve authoritative source bytes and MUST NOT decode, normalize, or rewrite a fixture as part of validation.
- **FR-007**: The fixture contract MUST require an explicit redistribution status and MUST reject committed third-party material whose redistribution permission is unknown or denied.
- **FR-008**: Shared comparison helpers MUST support typed common-model equality, ordered diagnostics, rendered bytes, exact source bytes, byte length, SHA-256, and timestamp results.
- **FR-009**: Shared comparisons MUST produce deterministic mismatch evidence based on portable fixture identifiers and logical field names rather than absolute filesystem locations.
- **FR-010**: Timestamp comparisons MUST distinguish restored, unsupported, unavailable, and failed outcomes and require equality only where the platform truthfully claims restoration.
- **FR-011**: The harness MUST prove that portable expectations generated from equivalent content in different source directories are identical and contain no original path, checkout path, username, hostname, or temporary-directory identifier.
- **FR-012**: The harness MUST provide reusable structural and semantic validation scenarios for valid and invalid Cue JSON without weakening the canonical schema or model rules.
- **FR-013**: Malformed regression records MUST declare a stable rejection class and diagnostic fragment, run deterministically, preserve their input bytes, and never panic.
- **FR-014**: Initial fuzz boundaries MUST cover complete Cue JSON decoding and source-envelope integrity preparation without invoking restoration, native subtitle codecs, network access, or external processes.
- **FR-015**: Fuzz seeds MUST include valid representative input, empty input, malformed JSON, structurally invalid Cue JSON, semantically invalid Cue JSON, and corrupt source-envelope declarations.
- **FR-016**: A fuzz-found regression retained in the repository MUST be minimized where practical, added to the governed manifest, and exercised by a deterministic regression test.
- **FR-017**: Fixture, malformed, conformance, and fuzz-smoke verification MUST run in the foreground through documented repository commands and terminate within a bounded operator-selected duration.
- **FR-018**: Completing S005 MUST NOT add or advertise native SRT or WebVTT parsing, rendering, conversion, or grammar coverage and MUST NOT create CI workflows owned by issue #8.
- **FR-019**: Documentation and the unreleased changelog MUST describe the fixture contract, contributor workflow, verification commands, provenance obligations, and the boundary between S005 infrastructure and later codec or CI work.

### Key Entities

- **Fixture record**: One stable manifest entry governing a payload's identity, portable location, purpose, origin, redistribution decision, integrity values, and intentional byte characteristics.
- **Fixture payload**: One regular file beneath the governed root whose exact bytes are authoritative and whose path and digest are declared once.
- **Golden expectation**: A portable expected model, diagnostic sequence, rendered byte sequence, integrity value, or timestamp result associated with a fixture identifier.
- **Malformed regression**: A minimal governed input paired with a stable rejection class and diagnostic fragment.
- **Fuzz boundary**: A panic-safe, side-effect-free validation entry point supplied with permanent representative and malformed seeds.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: One hundred percent of committed fixture payloads have exactly one complete provenance record and pass byte-length and SHA-256 verification.
- **SC-002**: One hundred percent of tested fixture drift, undeclared files, unsafe paths, symbolic links, duplicate identities, and disallowed redistribution states are rejected before a conformance assertion begins.
- **SC-003**: Every supported golden comparison surface has at least one passing and one deliberately mismatching test with stable path-free evidence.
- **SC-004**: Equivalent portable expectations produced from two distinct local directory roots are byte-identical and contain zero root, username, hostname, or temporary-path fragments.
- **SC-005**: Every registered malformed case produces the same declared rejection class and diagnostic fragment across repeated runs with zero panics.
- **SC-006**: Each initial fuzz boundary completes a bounded smoke run from all required seeds with zero panics, hangs, destination writes, network access, or external-process launches.
- **SC-007**: Fixture hashes remain identical after a fresh checkout applies repository line-ending rules, including all intentionally nonstandard byte fixtures.
- **SC-008**: Full repository tests, race detection, static analysis, text-integrity checks, fixture verification, malformed regressions, and bounded fuzz smoke verification pass before publication.
- **SC-009**: The specification, task plan, issue acceptance criteria, documentation, and Project tracking agree that S005 supplies test infrastructure only and leaves CI workflows and native codec grammar coverage unimplemented.

## Assumptions

- Issue #7 is the only implementation issue in S005 because issue #9 has a separate GitHub-automation review surface and issue #8 depends on this fixture foundation.
- The initial committed corpus remains small, synthetic or clearly redistributable, and focused on canonical Cue JSON and source-envelope validation rather than real-world subtitle grammar breadth.
- Future codec slices will add format-specific fixtures and expectations under this contract without redesigning its identity, provenance, integrity, or path-safety rules.
- The current canonical schema, model validation, and source preparation boundaries remain authoritative and are exercised rather than duplicated.
- Hosted cross-platform execution remains issue #8; S005 may make its helpers portable and run current-host verification without claiming the future CI matrix.
