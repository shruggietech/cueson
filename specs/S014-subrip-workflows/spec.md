# Feature Specification: Deliver Usable SubRip Workflows

**Feature Branch**: `codex/S014-subrip-workflows`

**Created**: 2026-09-11

**Status**: Implemented, pending pull-request verification

**Input**: Deliver issues #30 and #31 as one Spec Kit work slice: establish the shared v1 text-codec, model, detection, decoding, diagnostic, and capability foundation, then provide complete native SubRip encode and render workflows while preserving exact restoration and deferring WebVTT and cross-format conversion.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Encode SubRip as Cue JSON (Priority: P1)

A user can give Cueson a supported SubRip file and receive a structurally valid, semantically valid Cue JSON document that exposes useful cue text and timing while retaining the exact original source bytes.

**Why this priority**: Native encoding is the first missing workflow that turns the released repository foundation into usable subtitle software.

**Independent Test**: Encode each accepted SubRip fixture, validate the resulting Cue JSON, inspect its common and SubRip-specific fields, and restore the original source to prove byte identity.

**Acceptance Scenarios**:

1. **Given** a supported UTF-8 SubRip file, **When** the user runs `cueson encode`, **Then** Cueson writes valid Cue JSON containing ordered cues, normalized timing, raw text, plain text, lines, native SubRip details, source integrity metadata, and the original bytes.
2. **Given** an accepted source using a documented line-ending, BOM, timing, coordinate, sequence, tag, or encoding variation, **When** it is encoded, **Then** the semantic model is correct and the source envelope retains the exact input bytes.
3. **Given** an encoded document whose source payload has not been altered, **When** the user restores it, **Then** the restored file has the same byte length and SHA-256 digest as the input.
4. **Given** a source path containing user or machine-specific directories, **When** the file is encoded, **Then** the Cue JSON contains only the portable source basename and no original path or local machine identifier.

---

### User Story 2 - Render SubRip from the Structured Model (Priority: P1)

A user can render a valid SubRip Cue JSON document into deterministic SubRip syntax after editing its structured cue model, without confusing rendered output with exact source restoration.

**Why this priority**: Encoding alone is archival. Rendering makes the common model actionable and proves that normalized content is a real interchange surface.

**Independent Test**: Render representative valid Cue JSON documents, compare output with canonical SubRip golden files, re-encode the output, and compare the resulting normalized cue model.

**Acceptance Scenarios**:

1. **Given** valid SubRip Cue JSON, **When** the user runs `cueson render --to srt`, **Then** Cueson writes deterministic SubRip with ordered cues, canonical timestamps, preserved multiline content, and applicable coordinates.
2. **Given** a structured cue model whose semantic values were edited, **When** it is rendered and re-encoded, **Then** the new normalized model is semantically equivalent to the edited model.
3. **Given** an existing destination, **When** the user renders without `--force`, **Then** Cueson refuses to overwrite it and leaves the existing file unchanged.
4. **Given** a document for a schema-recognized format whose renderer is unavailable, **When** rendering is requested, **Then** Cueson reports the missing capability clearly and creates no output.

---

### User Story 3 - Diagnose Format and Source Variations (Priority: P2)

A user receives deterministic, actionable results when format evidence, text encoding, or SubRip structure is ambiguous, inconsistent, tolerated, or malformed.

**Why this priority**: Real-world SubRip files vary substantially. Truthful detection and diagnostics prevent silent corruption while allowing documented dialects.

**Independent Test**: Exercise content and extension disagreements, explicit format and encoding overrides, supported legacy encodings, malformed blocks, ambiguous inputs, and unsupported formats while comparing ordered diagnostics and output-creation behavior.

**Acceptance Scenarios**:

1. **Given** content that satisfies the documented SubRip grammar but has a disagreeing extension, **When** automatic detection runs, **Then** content evidence wins and a deterministic warning identifies the disagreement.
2. **Given** input whose format cannot be distinguished safely, **When** no explicit format is supplied, **Then** encoding fails with an actionable diagnostic and creates no Cue JSON.
3. **Given** invalid UTF-8 with no Unicode BOM, **When** no explicit legacy encoding is supplied, **Then** Cueson refuses to guess between legacy encodings and tells the user to provide `--encoding`.
4. **Given** a supported explicitly named legacy encoding, **When** the source is encoded, **Then** decoded cue fields are correct while exact original bytes remain authoritative.
5. **Given** malformed but preservable source content, **When** Cueson can still identify valid cues safely, **Then** it retains all source bytes, emits deterministic diagnostics for each unmodeled construct, and never silently presents the result as lossless semantic interpretation.

### Edge Cases

- Empty input, whitespace-only input, and files containing no valid cue must fail without producing a partial output document.
- A Unicode BOM that conflicts with an explicit encoding override must fail rather than choosing one silently.
- UTF-16 input with an odd trailing byte, invalid surrogate sequence, or unsupported byte order must fail safely.
- Input larger than the documented maximum must be rejected before an unbounded allocation or output write.
- Mixed CRLF, LF, and lone-CR separators must retain exact source bytes and report a mixed line-ending observation without changing decoded cue lines.
- Missing, non-integer, duplicated, decreasing, or unusually large sequence lines must not cause destructive renumbering of native source data.
- A timecode with period milliseconds or one-to-three millisecond digits must parse deterministically while retaining its native timing line.
- Reversed times, overflow, negative time values, missing arrows, or incomplete coordinate groups must produce deterministic errors without panic.
- Empty payload lines inside a cue and trailing whitespace in payload text must not disappear from the raw derived view.
- Formatting tags and speaker prefixes must not be removed from `payload.raw_text` or `payload.lines`; plain text and speaker observations are derived separately.
- Unknown tags must remain in raw text and be removed from plain text only when the documented tag scanner can do so without inventing semantics.
- Standard output mode must contain only Cue JSON or rendered SubRip bytes, never informational diagnostics.
- Any failure after a valid invocation must avoid leaving a newly created partial destination or replacing an existing destination.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: S014 MUST deliver issues #30 and #31 through one Spec Kit record while preserving each issue's independently verifiable acceptance criteria.
- **FR-002**: The common document model MUST continue to expose normalized timing, `payload.raw_text`, `payload.plain_text`, `payload.lines`, speakers, tokens, OCR extension space, placement, format-native data, ordered diagnostics, statistics, and the authoritative source envelope.
- **FR-003**: The shared native-codec boundary MUST provide format identity, content and extension evidence, decode capability, render capability, and deterministic capability lookup without owning exact source restoration.
- **FR-004**: Automatic format detection MUST prefer content signature or grammar over extension, warn on extension disagreement, and reject ambiguity unless an explicit format is supplied.
- **FR-005**: The explicit format selector MUST accept `auto`, `srt`, and `vtt`; S014 MUST implement native SubRip selection and report WebVTT as schema-recognized but codec unavailable.
- **FR-006**: Input processing MUST enforce a documented maximum source size of 64 MiB before retaining or decoding source bytes.
- **FR-007**: The decoder MUST honor UTF-8 and UTF-16 LE or BE BOMs, attempt strict UTF-8 when no BOM exists, recognize BOM-less UTF-16 only from strong and tested byte-pattern evidence, and reject malformed Unicode.
- **FR-008**: The supported explicit encoding names MUST be `utf-8`, `utf-16le`, `utf-16be`, `windows-1252`, and `iso-8859-1`, with stable aliases documented in the CLI and SubRip format contract.
- **FR-009**: Without a BOM or strict UTF-8 success, automatic decoding MUST NOT guess between supported single-byte legacy encodings and MUST require `--encoding`.
- **FR-010**: An explicit encoding that conflicts with a present Unicode BOM MUST fail with an actionable diagnostic.
- **FR-011**: Every successful native encode MUST capture filesystem timestamp metadata before reading source bytes and MUST retain the original bytes, byte length, lowercase SHA-256, media type, safe basename, encoding observations, and line-ending observations in one primary source asset.
- **FR-012**: Generated Cue JSON MUST NOT contain the original filesystem path, current working directory, username, hostname, drive or mount identity, temporary path, or another local machine identifier.
- **FR-013**: Native SubRip parsing MUST use an explicit deterministic block and line state model rather than destructive whole-file blank-line normalization.
- **FR-014**: The parser MUST support integer sequence lines and preserve missing, duplicated, decreasing, non-integer, and irregular source sequence distinctions without silently renumbering native data.
- **FR-015**: The parser MUST support comma and tolerated period millisecond separators, one-to-three millisecond digits, and non-negative hour values that fit the common millisecond range.
- **FR-016**: The parser MUST preserve the complete native timing line and parse an optional complete `X1`, `X2`, `Y1`, `Y2` coordinate suffix without claiming WebVTT placement equivalence.
- **FR-017**: The parser MUST preserve ordered multiline cue payloads, empty payload lines, decoded trailing whitespace, inline formatting tags, and source block order in derived fields without changing the source envelope.
- **FR-018**: `payload.raw_text` MUST join decoded payload lines with LF without stripping content, `payload.lines` MUST contain those decoded lines without line terminators, and `payload.plain_text` MUST remove documented SubRip formatting tags while retaining their textual content.
- **FR-019**: A deterministic `Name:` prefix heuristic MAY add one `heuristic` speaker observation but MUST NOT alter raw text, plain text, lines, or source bytes.
- **FR-020**: Malformed or unmodeled source constructs MUST remain recoverable through the source envelope and MUST generate ordered diagnostics; no successful encode may silently discard a known interpretation failure.
- **FR-021**: A successful native SubRip document MUST declare `format: "subrip"`, `format_support.status: "experimental"`, native ingest and render support true, restore support true, and OCR requirement false until the v1 stability gate is completed.
- **FR-022**: The model and schema MUST continue to accept valid envelope-only WebVTT documents while distinguishing them from native SubRip documents with incompatible format data rejected.
- **FR-023**: `cueson encode INPUT` MUST default to `INPUT.cueson.json`, and MUST support `--output`, `--force`, `--format`, `--encoding`, `--pretty`, `--stdout`, `--no-speaker-detection`, quiet, silent, no-color, and help behavior consistent with the CLI contract.
- **FR-024**: Encode standard-output mode MUST emit only one valid Cue JSON document to stdout and must send permitted diagnostics to stderr.
- **FR-025**: `cueson render INPUT.cueson.json --to srt` MUST render the structured model rather than restoring source bytes and MUST support `--output`, `--force`, `--strict`, quiet, silent, no-color, and help behavior consistent with the CLI contract.
- **FR-026**: SubRip rendering MUST emit deterministic cue order, canonical integer sequence lines, canonical `HH:MM:SS,mmm` timestamps, applicable complete coordinates, the structured raw payload, one blank line between cues, and a final platform-independent LF newline.
- **FR-027**: Rendered output MUST be accepted by the native parser and re-encoding it MUST produce a semantically equivalent normalized cue model for all representable SubRip data.
- **FR-028**: `restore` MUST remain the only operation that claims byte identity; every accepted SubRip encode fixture MUST restore to the original byte length and SHA-256 independently of rendered output.
- **FR-029**: Invalid invocation or conflicting options MUST exit 2, while I/O, detection, decoding, parsing, model, schema, integrity, or rendering failures after a valid command MUST exit 1; all failures MUST avoid partial output publication.
- **FR-030**: Parser, timing, detection, decoder, renderer, model, schema, source integration, CLI, malformed-input, and round-trip behavior MUST be covered by focused tests and provenance-recorded fixtures before the implementation is considered complete.
- **FR-031**: Every documented SubRip grammar and tolerated-variant row in `docs/formats/srt.md` MUST map to at least one parser fixture and every applicable rendering row MUST map to a renderer fixture.
- **FR-032**: Fuzz seeds and tests MUST cover format detection, SubRip block parsing, time parsing, tag scanning, text decoding, and renderer/parser cycles without panic, path leakage, or silent accepted-data loss.
- **FR-033**: Existing schema, source restoration, conformance, cross-platform, release-proof, repository-format, security, and CI behavior MUST remain green.
- **FR-034**: Documentation and runtime help MUST describe only the implemented S014 capability, exact encoding names and aliases, size limit, output behavior, diagnostics, experimental support state, and the distinction between render and restore.
- **FR-035**: S014 MUST NOT implement native WebVTT parsing or rendering, cross-format conversion, `validate`, `inspect`, shell completion, version 1.0.0 publication, milestone closure, public schema hosting, production-domain changes, or pull-request merge.
- **FR-036**: The official pull request MUST close both #30 and #31, pass all current-head checks, address every external review finding, request no more than one second Codex round, and stop for the operator's final merge ritual.
- **FR-037**: The evolving canonical schema and executable MUST advance together from the immutable released v0.0.0 identity to development version `0.1.0`; the immutable `schema/releases/v0.0.0/cueson.schema.json`, tag, and release assets MUST remain unchanged.

### Key Entities

- **Codec Registry**: The authoritative collection of installed native format capabilities and their format aliases.
- **Detection Result**: Ordered content, structure, and extension evidence with a selected format, confidence, and diagnostics.
- **Decode Options**: Explicit format, optional encoding override, source-size boundary, and semantic derivation controls for native ingest.
- **Render Options**: Target format, strictness, and deterministic output controls for model-driven serialization.
- **SubRip Document**: Cue JSON with a common cue model, SubRip-native document and cue data, diagnostics, statistics, and one exact primary source asset.
- **SubRip Cue**: One source-ordered cue with native sequence and timing lines, normalized timing, raw and plain payload views, line structure, optional coordinates, and optional derived speaker observations.
- **Encoding Observation**: The selected encoding, BOM state, line-ending classification, and confidence retained with the primary source asset.
- **Diagnostic**: A deterministic severity, code, message, and optional cue or source-order association describing tolerated, ambiguous, malformed, or unsupported input.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can encode a representative UTF-8 SubRip file, produce valid Cue JSON, render it, re-encode the rendering, and restore the original source through one documented workflow with every command returning the expected status.
- **SC-002**: One hundred percent of documented SubRip grammar and tolerated-variant rows have traceable accepted or malformed parser fixtures, and one hundred percent of applicable renderer rows have golden output coverage.
- **SC-003**: Every accepted SubRip fixture restores with 100% byte identity and an identical SHA-256 digest, including BOM, legacy-encoding, CRLF, LF, lone-CR, and mixed-line-ending fixtures.
- **SC-004**: Every generated Cue JSON document validates structurally and semantically, exposes cue text without base64 decoding, and contains zero original filesystem paths or local machine identifiers.
- **SC-005**: Every supported encoding name produces deterministic decoded fields, and every ambiguous or conflicting encoding case fails with exactly one actionable primary diagnostic and no output file.
- **SC-006**: All canonical SubRip render fixtures are byte-for-byte deterministic across Windows, macOS, and Linux, and every rendered fixture is accepted by the native parser.
- **SC-007**: Malformed corpus and bounded fuzz execution produce zero panics, zero hangs, zero writes outside caller-selected destinations, and zero silently accepted interpretation gaps.
- **SC-008**: The final S014 pull-request head has zero failed or pending checks, zero unresolved review threads, every external finding addressed, and no more than two Codex review rounds.

## Assumptions

- S014 preserves the existing cross-platform, dependency-light product boundary and does not introduce a public library contract.
- The canonical schema on unreleased `main` advances to development version `0.1.0` with the executable so released schema identity is never reused for changed bytes; final v1 release preparation owns the stable `1.0.0` transition.
- Single-byte legacy input cannot be detected reliably from bytes alone, so explicit user selection is safer than a confidence guess.
- Speaker observation uses a conservative deterministic prefix heuristic and remains derived; false positives do not justify source mutation.
- WebVTT remains schema-recognized and envelope-only during S014.
- The v1 documentation and release-candidate issues own final stable support declarations and the `1.0.0` version transition.
