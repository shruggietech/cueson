# Feature Specification: Deliver Native WebVTT Workflows

**Feature Branch**: `codex/S015-webvtt-workflows`

**Created**: 2026-09-11

**Status**: Implemented, pending pull-request verification

**Input**: Deliver issue #32 through Spec Kit as work slice S015, providing native WebVTT encode and model-driven render workflows while preserving exact source bytes and every ordered native block.

## User Scenarios & Testing

### User Story 1 - Encode WebVTT as Cue JSON (Priority: P1)

A user supplies a conforming or documented real-world WebVTT file and receives valid Cue JSON containing accessible normalized cues, complete WebVTT-native structure, deterministic diagnostics, and an authoritative source envelope.

**Why this priority**: Native WebVTT ingest is the missing prerequisite for every remaining conversion and complete-CLI outcome on the v1 critical path.

**Independent Test**: Encode representative accepted WebVTT files containing every documented header, block, cue, setting, markup, entity, inline-timestamp, whitespace, ordering, and rolling-caption construct; validate the documents and restore every original byte.

**Acceptance Scenarios**:

1. **Given** a UTF-8 WebVTT file with a valid signature and cues, **When** the user runs `cueson encode`, **Then** the command writes valid Cue JSON with normalized cue timing and text, WebVTT-native fields, exact source bytes, and native ingest and render capability declarations.
2. **Given** a file containing NOTE, STYLE, REGION, cue identifiers, settings, markup, entities, inline timestamps, adjacent blocks, and overlapping cues, **When** the user encodes it, **Then** every block retains total source order and every source construct is either interpreted or preserved with a deterministic diagnostic.
3. **Given** a WebVTT file whose `.vtt` extension disagrees with its content, **When** automatic detection runs, **Then** content evidence wins and the disagreement is reported without losing source information.
4. **Given** malformed Unicode, an invalid signature, invalid cue timing, an over-limit source, or an unpreservable malformed construct, **When** the user encodes it, **Then** the command fails deterministically without panic or partial output.

---

### User Story 2 - Render Cue JSON as WebVTT (Priority: P2)

A user renders a valid WebVTT Cue JSON document into deterministic, conforming WebVTT derived from the structured model rather than receiving a disguised restoration of the original source.

**Why this priority**: Model-driven WebVTT output makes Cue JSON useful as an editable interchange format and completes the native codec needed before cross-format conversion can begin.

**Independent Test**: Render representative WebVTT documents to files and standard output, compare canonical bytes across platforms, re-encode each result, and prove semantic equivalence for every representable field.

**Acceptance Scenarios**:

1. **Given** valid WebVTT Cue JSON, **When** the user runs `cueson render --to vtt`, **Then** the command emits a valid signature, document blocks and cues in total model order, deterministic timing and settings, preserved native markup, and a final LF newline.
2. **Given** structured fields that no longer agree with reusable native lexical fields, **When** rendering occurs, **Then** structured semantics win and the affected syntax is regenerated canonically rather than emitting stale lexical content.
3. **Given** valid source-envelope bytes alongside edited structured cues, **When** the user renders and separately restores, **Then** render reflects the edited model while restore reproduces the original bytes exactly.
4. **Given** invalid or unsupported WebVTT model content, **When** rendering occurs, **Then** the command reports a deterministic failure and does not create or replace a destination.

---

### User Story 3 - Preserve and Explain Native Fidelity (Priority: P3)

A user working with unusual but preservable WebVTT content can trust that Cueson keeps it in source order, does not deduplicate rolling cues, and clearly distinguishes raw native syntax from derived consumer views.

**Why this priority**: WebVTT carries document structure and timed-text semantics that cannot safely be flattened into a SubRip-shaped cue list.

**Independent Test**: Exercise malformed, unknown, overlapping, adjacent-block, token-timing, markup, and entity fixtures and assert ordered preservation, deterministic diagnostics, stable derived fields, and absence of silent loss.

**Acceptance Scenarios**:

1. **Given** overlapping rolling-caption cues with similar payloads, **When** the file is encoded, **Then** each cue remains a separate source-ordered cue and block without deduplication.
2. **Given** voice, class, language, ruby, entity, and inline-timestamp syntax, **When** the file is encoded, **Then** raw payload syntax remains available while plain text, speaker observations, and token timing are derived separately.
3. **Given** an unknown or misplaced but preservable block, setting, tag, or entity, **When** it is accepted, **Then** its raw content and source order remain available and a stable diagnostic identifies the interpretation limit.

### Edge Cases

- Empty input, whitespace-only input, a missing signature, and a file containing no valid cues fail without publishing partial Cue JSON.
- A UTF-8 BOM is accepted only at the start of the file and its source-byte presence remains observable without entering derived text.
- Non-UTF-8 input, malformed UTF-8, and an explicit non-UTF-8 encoding override fail because WebVTT is UTF-8 only; embedded NUL remains in the source envelope, becomes U+FFFD only in derived text, and produces a diagnostic.
- The signature line, optional header text, header metadata, and the blank line ending the header remain distinct from document blocks.
- `NOTE`, `STYLE`, and `REGION` text resembling cue timing does not become a cue, and a cue identifier resembling a block keyword remains distinguishable by surrounding grammar.
- Missing, duplicated, or unusual cue identifiers do not replace Cueson's document-local cue identity.
- Hours-optional timestamps, long hours, millisecond precision, boundary values, reversed intervals, and inline timestamps outside their cue interval produce deterministic accepted or rejected behavior.
- Repeated and unknown cue settings retain lexical order; recognized settings expose structured values without silently discarding invalid or duplicate occurrences.
- Whitespace-only payload lines, leading or trailing payload whitespace, adjacent blocks, and final files with or without a line terminator remain distinguishable in native data and source bytes.
- Markup nesting errors, unsupported tags, malformed entities, and invalid inline timestamps never cause a panic or destructive text normalization.
- Standard-output mode contains only Cue JSON or WebVTT bytes; informational and warning diagnostics use standard error.
- Any failure after valid invocation leaves a new destination absent and an existing destination unchanged.

## Requirements

### Functional Requirements

- **FR-001**: S015 MUST deliver issue #32 through one complete Spec Kit record and retain traceability to each issue acceptance criterion.
- **FR-002**: Native WebVTT encode MUST accept only valid UTF-8, optionally preceded by one UTF-8 BOM, MUST reject non-UTF-8 encodings and malformed Unicode, and MUST replace embedded NUL with U+FFFD only in derived text while preserving source bytes and emitting a diagnostic.
- **FR-003**: Automatic detection MUST recognize the WebVTT signature after permissible BOM handling, prefer that content evidence over extension, and report extension disagreement deterministically.
- **FR-004**: The parser MUST require the `WEBVTT` signature, preserve its raw line, preserve optional header text and metadata in order, and keep header termination distinct from subsequent blocks.
- **FR-005**: The parser MUST represent every cue and non-cue block in one total source order across the complete document, and the union of their source-order values MUST be unique and contiguous so an omitted body unit is detectable.
- **FR-006**: NOTE, STYLE, and REGION blocks MUST retain their complete raw lines and ordering; REGION blocks MUST additionally retain ordered raw setting occurrences and expose valid effective `id`, `width`, `lines`, `regionanchor`, `viewportanchor`, and `scroll` values without replacing raw content.
- **FR-007**: Misplaced, unknown, malformed, or unsupported but safely preservable blocks MUST retain raw content and total order with stable diagnostics; constructs that cannot be bounded or preserved safely MUST fail.
- **FR-008**: Cue identifiers MUST remain independent from document-local cue IDs and MUST preserve absence, duplication, whitespace, and native text subject to grammar validity.
- **FR-009**: Cue timing MUST support valid hours-optional WebVTT timestamps, expose non-negative normalized milliseconds, retain the native timing line, and reject invalid precision, overflow, equal endpoints, or reversed intervals; decreasing cue start times MUST remain in source order with a conformance diagnostic rather than being silently sorted.
- **FR-010**: Cue settings MUST preserve the complete raw settings text and ordered setting occurrences while exposing valid recognized values for vertical, line, position, size, align, and region semantics.
- **FR-011**: Unknown, invalid, and duplicate cue settings MUST be retained when the cue remains safely parseable and MUST produce ordered diagnostics rather than disappearing.
- **FR-012**: Cue payload parsing MUST preserve ordered decoded lines, whitespace-only lines, leading and trailing whitespace, raw markup, entities, and inline timestamps independently from derived views.
- **FR-013**: `payload.raw_text` MUST join decoded payload lines with LF without stripping content, and `payload.lines` MUST contain decoded lines without line terminators.
- **FR-014**: `payload.plain_text` MUST derive readable text through a deterministic WebVTT markup and entity scanner while leaving `payload.raw_text`, `payload.lines`, format-native data, and source bytes unchanged.
- **FR-015**: Voice markup MAY create deterministic speaker observations, and valid inline timestamps MAY create ordered token-timing observations, but neither derivation may invent source truth or mutate raw payload data.
- **FR-016**: Overlapping and rolling-caption cues MUST remain one-to-one in source order without content-based deduplication or destructive merging.
- **FR-017**: Every successful encode MUST use the existing 64 MiB bounded exact-source capture, retain the original bytes and integrity metadata, and preserve source restoration as a codec-independent operation.
- **FR-018**: Generated Cue JSON MUST NOT contain the original input path or another local machine identifier.
- **FR-019**: A successfully encoded document MUST use `format: "webvtt"`, compatible WebVTT-native format data, and truthful experimental native ingest, render, and restore capability declarations without claiming conversion support.
- **FR-020**: The common model and schema MUST continue to accept valid SubRip documents and reject incompatible native format data without regressing existing behavior.
- **FR-021**: `cueson encode INPUT` MUST support automatic and explicit `vtt` selection, reject WebVTT encoding overrides other than UTF-8 aliases, and preserve existing output, force, pretty, stdout, quiet, silent, no-color, and help contracts.
- **FR-022**: Encode standard-output mode MUST emit exactly one valid Cue JSON document to stdout and route permitted diagnostics to stderr.
- **FR-023**: `cueson render INPUT.cueson.json --to vtt` MUST render structured WebVTT semantics rather than restore source bytes and MUST preserve existing output, force, strict, quiet, silent, no-color, and help contracts.
- **FR-024**: Canonical rendering MUST emit a valid `WEBVTT` signature, deterministic LF line endings, blocks in total model order, deterministic timestamps and recognized settings, one grammar-valid separator between blocks, and a final LF newline.
- **FR-025**: Rendering MAY reuse native lexical content only when validation proves it remains consistent with the authoritative structured semantics; otherwise it MUST regenerate affected syntax canonically or fail if safe regeneration is impossible.
- **FR-026**: NOTE, STYLE, REGION, unknown preservable blocks, payload markup, and other native content without a complete structured equivalent MUST remain available to permissive same-format rendering with deterministic diagnostics unless the model explicitly removes them; strict rendering MUST reject known preserved conformance errors or unsafe ambiguities without output.
- **FR-027**: Rendered output MUST be accepted by the native WebVTT parser, and re-encoding it MUST produce a semantically equivalent normalized and native model for all representable content.
- **FR-028**: Restore MUST remain the only byte-identity operation, and every accepted WebVTT source fixture MUST restore to its original byte length and SHA-256 independently from rendering.
- **FR-029**: Invalid invocation and conflicting options MUST exit 2; detection, decoding, parsing, model, schema, integrity, rendering, or I/O failure after valid invocation MUST exit 1; all failures MUST avoid partial output publication.
- **FR-030**: Every planned structure row in `docs/formats/webvtt.md` MUST map to accepted or malformed parser evidence, and every applicable rendering behavior MUST map to deterministic golden evidence.
- **FR-031**: Tests MUST cover signature detection, UTF-8 handling, headers, every native block kind, identifiers, timing, settings, markup, entities, inline timestamps, whitespace, ordering, rolling cues, rendering, exact restore, malformed input, and CLI behavior before completion.
- **FR-032**: Fuzz seeds and bounded tests MUST cover detection, block parsing, timestamps, cue settings, markup and entity scanning, source-envelope integration, and renderer/parser cycles without panic, hang, path leakage, or silent accepted-data loss.
- **FR-033**: Existing schema, SubRip, restoration, conformance, cross-platform, release-proof, repository-format, security, and hosted CI behavior MUST remain green.
- **FR-034**: Documentation and runtime help MUST describe only the implemented S015 capability, UTF-8 restriction, size and output behavior, diagnostics, experimental support state, and render-versus-restore distinction.
- **FR-035**: S015 MUST NOT implement cross-format conversion, `validate`, `inspect`, shell completion, stable v1 support declarations, version `1.0.0`, tag or release publication, public schema hosting, production-domain changes, or pull-request merge.
- **FR-036**: The official pull request MUST close issue #32, pass all current-head checks, address every external review finding, request no more than one second Codex round, and stop for the operator's final merge ritual.

### Key Entities

- **WebVTT Document**: Cue JSON containing common cues, a total ordered native block sequence, document header data, diagnostics, statistics, capability declarations, and one authoritative source asset.
- **WebVTT Block**: A source-ordered cue, NOTE, STYLE, REGION, or preservable unknown unit with native raw content and any recognized structured semantics.
- **WebVTT Cue**: One source-ordered timed payload with an optional native identifier, normalized timing, native timing and setting data, raw and plain payload views, optional placement, speakers, and timed tokens.
- **Cue Setting Occurrence**: One ordered native setting with raw name and value plus recognized structured interpretation when valid.
- **Markup Observation**: A bounded interpretation of voice, class, language, ruby, entity, or timestamp syntax used to derive consumer views without replacing raw payload syntax.
- **Diagnostic**: A deterministic severity, code, message, and optional block, cue, line, or source-order association for tolerated or rejected input.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A user can encode a representative WebVTT file, produce valid Cue JSON, render it, re-encode the rendering, and restore the original source through one documented workflow with every command returning the expected status.
- **SC-002**: One hundred percent of the WebVTT structure rows documented for S015 have traceable accepted or malformed parser fixtures, and one hundred percent of applicable renderer rows have golden output coverage.
- **SC-003**: Every accepted WebVTT source fixture restores with 100 percent byte identity and an identical SHA-256 digest, including BOM, CRLF, LF, mixed-line-ending, block-order, markup, and rolling-caption fixtures.
- **SC-004**: Every generated Cue JSON document validates structurally and semantically, exposes readable cue text without source-byte decoding, retains total block order, and contains zero original filesystem paths or local machine identifiers.
- **SC-005**: Canonical WebVTT output is byte-for-byte deterministic across Windows, macOS, and Linux, and every renderer fixture is accepted by the native parser with equivalent representable semantics.
- **SC-006**: Every accepted source construct is interpreted or receives at least one deterministic diagnostic; no preservable block, setting, payload line, markup token, entity, or rolling cue silently disappears.
- **SC-007**: Malformed corpus and bounded fuzz execution produce zero panics, zero hangs, zero writes outside caller-selected destinations, and zero partial destination publications.
- **SC-008**: The final S015 pull-request head has zero failed or pending checks, zero unresolved review threads, every external finding addressed, and no more than two Codex review rounds.

## Assumptions

- WebVTT remains an experimental native capability during S015; issue #36 and the v1 release-candidate gate own the stable support declaration.
- The executable and evolving canonical schema remain at development version `0.1.0`; S015 does not alter immutable v0.0.0 release material.
- WebVTT source decoding is UTF-8 only. Existing generic explicit-encoding support remains available for SubRip and is not generalized into a false WebVTT legacy-encoding claim.
- Same-format rendering preserves ordered native blocks and raw native content while structured timing and recognized setting semantics remain authoritative when edited.
- Conversion loss accounting belongs exclusively to issue #33 and begins only after this native codec is merged.
- Existing codec, model, schema, source, CLI, fixture, and conformance boundaries are extended proportionally rather than replaced.
