# Feature Specification: Deliver Bidirectional Subtitle Conversion

**Feature Branch**: `codex/S016-bidirectional-conversion`

**Created**: 2026-09-11

**Status**: Implemented, pending pull-request verification

**Input**: Deliver issue #33 through Spec Kit as work slice S016, providing SubRip-to-WebVTT and WebVTT-to-SubRip conversion through the common model with complete deterministic loss accounting and strict prevention of lossy output.

## User Scenarios & Testing

### User Story 1 - Convert between SubRip and WebVTT (Priority: P1)

A user supplies a supported SubRip or WebVTT source, or a valid Cue JSON document derived from either format, and receives deterministic output in the other subtitle format through the common semantic model.

**Why this priority**: Bidirectional conversion is the largest remaining user-facing capability on the v1 critical path and unlocks the rest of the complete CLI, hardening, and documentation work.

**Independent Test**: Convert representative SubRip and WebVTT sources and Cue JSON documents in both directions, parse every result with the target codec, and compare all representable common semantics with the input model.

**Acceptance Scenarios**:

1. **Given** a supported SubRip source, **When** the user converts it to WebVTT, **Then** the output is valid deterministic WebVTT and retains every representable cue timing and payload semantic.
2. **Given** a supported WebVTT source, **When** the user converts it to SubRip, **Then** the output is valid deterministic SubRip and retains every representable cue timing and payload semantic.
3. **Given** a valid SubRip or WebVTT Cue JSON document, **When** the user targets the other format, **Then** conversion uses the structured model rather than restoring or reparsing its source envelope.
4. **Given** ambiguous source detection or an unsupported source or target format, **When** conversion is attempted, **Then** the command fails with a truthful typed diagnostic and publishes no output.

---

### User Story 2 - Understand and prevent conversion loss (Priority: P2)

A user can see every known source semantic that cannot be represented in the target format and can require conversion to stop before output whenever any known loss exists.

**Why this priority**: Cross-format output is trustworthy only when format-specific information cannot disappear silently and automated callers can enforce a no-loss policy.

**Independent Test**: Exercise every compatibility-matrix row in both directions, compare ordered loss entries and warnings with golden expectations, and prove strict mode leaves new and existing destinations unchanged for every lossy case.

**Acceptance Scenarios**:

1. **Given** a conversion containing only representable semantics, **When** it runs normally or strictly, **Then** it succeeds without a loss warning.
2. **Given** a WebVTT document containing NOTE, STYLE, REGION, cue settings, inline timing, or format-specific markup that SubRip cannot represent, **When** normal conversion runs, **Then** output is produced and every known loss appears in a stable deterministic report and warning stream.
3. **Given** SubRip coordinate semantics without enough viewport information for a proven WebVTT mapping, **When** normal conversion runs, **Then** the coordinates are not presented as equivalent and their omission is explicitly reported.
4. **Given** any known loss, **When** strict conversion runs, **Then** the command exits with runtime failure and creates no target or change to an existing target.

---

### User Story 3 - Automate conversion safely (Priority: P3)

A user or script can select input and target formats, choose stdout or a file, control replacement and diagnostics, and receive stable results without partial publication or machine-specific data leakage.

**Why this priority**: Conversion must compose reliably in command pipelines and file workflows across supported platforms before its v1 contract can stabilize.

**Independent Test**: Run the complete invocation, stream, output, overwrite, suppression, encoding, strict, malformed-input, and injected-write-failure matrix across supported platforms and compare statuses, streams, destination bytes, and diagnostics.

**Acceptance Scenarios**:

1. **Given** a valid conversion with no file destination, **When** the command succeeds, **Then** stdout contains only target subtitle bytes and diagnostics use stderr.
2. **Given** an explicit new file destination, **When** conversion succeeds, **Then** the complete target is published once with no temporary artifact left behind.
3. **Given** an existing regular-file destination, **When** conversion runs without replacement authorization or fails before completion, **Then** the original file remains byte-identical.
4. **Given** quiet or silent diagnostic filtering, **When** conversion emits information, loss warnings, or errors, **Then** filtering matches the established command contract without suppressing explicitly requested stdout payload bytes.

### Edge Cases

- The input extension disagrees with strong SubRip or WebVTT content evidence.
- A source-format override disagrees with the detected grammar, or an encoding override is incompatible with WebVTT.
- Cue JSON declares one native format while the user requests the same target format, which belongs to model rendering rather than cross-format conversion.
- The input contains zero cues, duplicate or overlapping cues, zero-duration cues, very long valid timestamps, empty or whitespace-only payload lines, or native source-order gaps.
- A payload contains markup that is valid in the source format but unsafe, ambiguous, or semantically different in the target grammar.
- WebVTT includes ordered non-cue blocks, duplicate settings, unknown settings, inline timestamps, or speaker observations with no direct SubRip representation.
- SubRip includes sequence numbers, coordinate suffixes, or lexical constructs with no proven WebVTT equivalent.
- Multiple losses apply to the same cue or document and must remain distinct, stable, and deterministically ordered.
- A destination exists, is a link or non-regular entry, has a missing parent, races into existence, or encounters a short or failed write.
- Standard output or standard error fails partway through a write.
- The source or Cue JSON exceeds established bounded-input limits, is malformed, fails schema or semantic validation, or contains inconsistent native and common fields.

## Requirements

### Functional Requirements

- **FR-001**: The product MUST expose `cueson convert INPUT --to srt|vtt` only when the complete S016 behavior is implemented.
- **FR-002**: Conversion MUST accept a supported SubRip or WebVTT source and a valid Cue JSON document whose native format is SubRip or WebVTT.
- **FR-003**: Source conversion MUST use the same bounded exact acquisition, content-first detection, explicit format selection, decoding, parsing, and diagnostic contracts as native encoding.
- **FR-004**: Cue JSON conversion MUST validate JSON syntax, the active schema, common semantic invariants, source-envelope integrity, and software/schema compatibility before producing output.
- **FR-005**: Cue JSON conversion MUST derive target output from validated structured common and native fields and MUST NOT restore or reparse preserved source bytes as a substitute for conversion.
- **FR-006**: The target format MUST be explicitly selected, MUST have an installed renderer, and MUST differ from the input document's native format; same-format model serialization remains the responsibility of `render`.
- **FR-007**: Conversion MUST map every common cue's normalized start, end, payload text, line structure, speaker observations, and token timing deliberately according to the documented compatibility matrix.
- **FR-008**: Conversion MUST preserve cue order, overlapping-cue order, duplicate payloads, and zero-duration cues whenever the target grammar accepts those values, without sorting, deduplication, or inferred timing changes.
- **FR-009**: The compatibility matrix MUST classify every supported SubRip and WebVTT common and native semantic as represented, transformed, or lost in each conversion direction.
- **FR-010**: Every known transformed or non-representable semantic that changes meaning or disappears MUST produce one or more stable machine-testable loss entries; exact source restoration capability MUST NOT be used to label conversion as lossless.
- **FR-011**: Each loss entry MUST provide a stable code, severity, human-readable message, source format, target format, and deterministic document, block, cue, or field association when applicable.
- **FR-012**: Loss entries MUST use a single documented deterministic ordering independent of map iteration, filesystem behavior, platform, or diagnostic filtering.
- **FR-013**: Normal conversion MUST publish valid target output when all encountered losses are permitted and MUST emit a clear warning for each loss according to the established diagnostic filters.
- **FR-014**: Strict conversion MUST reject the operation before publication when any known loss exists, return exit code 1, and identify the first loss deterministically while retaining the complete loss result for tests and internal callers.
- **FR-015**: Strict rejection, parser failure, validation failure, rendering failure, or output failure MUST leave a new destination absent and an existing destination byte-identical.
- **FR-016**: WebVTT NOTE, STYLE, REGION, cue identifiers, cue settings, inline timestamps, header metadata, ordered unknown blocks, and format-specific markup MUST each receive explicit documented conversion treatment.
- **FR-017**: SubRip sequence numbers, coordinate semantics, legacy source encoding observations, and format-specific markup MUST each receive explicit documented conversion treatment.
- **FR-018**: SubRip coordinates MUST NOT be represented as equivalent WebVTT placement unless all information required for a deterministic meaning-preserving mapping is available; otherwise their omission MUST be reported as loss.
- **FR-019**: Markup conversion MUST preserve readable payload meaning, MUST escape target-syntax injection, MUST map only proven equivalent constructs, and MUST report every removed or semantically weakened construct.
- **FR-020**: Speaker and token-timing observations MUST remain explicitly derived; conversion MUST NOT invent speakers, token boundaries, timestamps, viewport dimensions, styles, regions, or placement values absent from the validated input model.
- **FR-021**: Target output MUST be deterministic UTF-8 without BOM, use LF line endings, end with one LF, and satisfy the target codec's canonical grammar.
- **FR-022**: Every successful conversion result MUST be accepted by the target native parser and re-encoding that result MUST retain all semantics classified as represented or transformed.
- **FR-023**: The conversion library boundary MUST return target bytes and the complete ordered loss report without printing or writing files directly.
- **FR-024**: `convert` MUST support `--to`, `--from auto|cueson|srt|vtt`, `--output`, `--force`, `--strict`, source `--encoding`, `--no-speaker-detection`, `--quiet`, `--silent`, `--no-color`, and help consistently with the existing command contracts.
- **FR-025**: Without `--output`, and with `--output -`, successful conversion MUST emit exactly the target subtitle bytes to stdout; permitted diagnostics MUST use stderr.
- **FR-026**: With a real file destination, conversion MUST require an existing parent, reject unsafe or non-regular destinations even with `--force`, refuse implicit replacement, and publish complete output transactionally.
- **FR-027**: Invalid invocation, missing operands, unknown options, conflicting options, invalid format or encoding names, `--force` without a real file destination, an encoding supplied for Cue JSON, and an encoding incompatible with an explicit WebVTT source MUST exit 2.
- **FR-028**: Detection, decoding, parsing, Cue JSON validation, source integrity, same-format conversion, missing capability, strict-loss, rendering, or runtime I/O failure after a valid invocation MUST exit 1; same-format failure MUST direct the user to `render`.
- **FR-029**: Quiet mode MUST suppress informational diagnostics while retaining loss warnings and errors; silent mode MUST suppress informational and loss-warning diagnostics while retaining errors; neither mode may suppress requested stdout payload bytes.
- **FR-030**: Loss entries and target-render diagnostics retained by the conversion result MUST contain no original filesystem path, local machine identifier, source bytes, or unsafe unbounded payload excerpt; CLI invocation and I/O errors MAY identify the literal caller-supplied path needed to correct the failure.
- **FR-031**: Focused tests and provenance-recorded fixtures MUST cover every compatibility-matrix row, both input classes, both directions, strict and permissive modes, parser cycles, malformed inputs, and output transactions before completion.
- **FR-032**: Bounded fuzz or generative tests MUST exercise loss classification, markup conversion, renderer/parser cycles, and hostile structured input without panic, hang, path leakage, or silent semantic disappearance.
- **FR-033**: Existing native encode, render, exact restore, schema, conformance, cross-platform, release-proof, repository-format, and security behavior MUST remain green.
- **FR-034**: Maintained documentation, runtime help, architecture, schema guidance, format matrices, and changelog MUST distinguish conversion from rendering and restoration and describe only implemented S016 behavior.
- **FR-035**: S016 MUST preserve development version `0.1.0`, experimental SubRip and WebVTT support states, and every immutable v0.0.0 release artifact.
- **FR-036**: S016 MUST NOT implement `validate`, `inspect`, shell completion, stable v1 support declarations, version `1.0.0`, release-candidate preparation, tag or release publication, public schema hosting, production-domain changes, or pull-request merge.
- **FR-037**: The official pull request MUST close issue #33, pass every current-head check, address every external review finding, request no more than one second Codex round, and stop for the operator's final merge ritual.

### Key Entities

- **Conversion Request**: A validated source or Cue JSON input, its source-format evidence, explicit target format, conversion policy, and destination policy.
- **Conversion Result**: Deterministic target subtitle bytes paired with the complete ordered loss report and non-loss diagnostics, independent of filesystem publication.
- **Loss Entry**: A stable description of one known semantic transformation or omission, including code, severity, formats, message, and optional document, block, cue, or field association.
- **Compatibility Matrix Row**: The bidirectional decision for one common or native semantic, classifying it as represented, transformed, or lost and naming its required evidence.
- **Target Document**: A transient target-format model containing only semantics deliberately mapped from the validated source model before canonical rendering.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Users can convert representative SubRip and WebVTT sources and Cue JSON documents in both directions through one documented command, and 100 percent of produced files are accepted by the target native parser.
- **SC-002**: One hundred percent of compatibility-matrix rows in both directions have traceable conversion and loss-accounting evidence, with zero undocumented drops of supported common or native semantics.
- **SC-003**: Every loss-bearing fixture produces byte-for-byte identical ordered loss results and warning order across Windows, macOS, and Linux.
- **SC-004**: Strict mode prevents 100 percent of known lossy cases from creating or changing a destination, including forced replacement and injected output failures.
- **SC-005**: Every loss-free fixture preserves all common cue timings, payload meaning, cue order, speaker observations, and token timing that the target format can represent.
- **SC-006**: Every generated output is deterministic UTF-8 LF text with one final newline, contains zero local machine identifiers, and is byte-identical across repeated runs.
- **SC-007**: Malformed corpus and bounded fuzz execution produce zero panics, zero hangs, zero writes outside caller-selected destinations, and zero silent accepted-data loss.
- **SC-008**: The final S016 pull-request head has zero failed or pending checks, zero unresolved review threads, every external finding addressed, and no more than two Codex review rounds.

## Assumptions

- SubRip and WebVTT remain experimental native capabilities during S016; later v1 documentation and release-candidate gates own stable support declarations.
- The existing common cue model is the semantic bridge. Format-native fields inform explicit mapping and loss decisions but never replace source truth.
- Cross-format conversion targets only SubRip and WebVTT in S016. Reserved or future formats fail truthfully as unknown or missing-capability cases.
- The conversion result's ordered loss report is an internal execution contract for S016 and CLI diagnostics; persistent public Cue JSON additions occur only if planning proves they are required for a stable v1 consumer contract.
- Source input behavior reuses native encode detection and decoding rules, while Cue JSON input behavior reuses active schema, semantic, and source-integrity validation.
- Conversion writes target subtitle bytes to stdout by default, matching model-driven render behavior and avoiding inferred destination names.
- No post-v1 website or production-domain work is part of this slice.
