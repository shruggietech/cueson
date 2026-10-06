# Feature Specification: Consumer speaker attribution and media timing

**Feature Branch**: `codex/S033-consumer-speaker-attribution`

**Created**: 2026-10-06

**Status**: Specified for S033; implementation pending analysis gate

**Input**: [Issue #82](https://github.com/shruggietech/cueson/issues/82), owner-approved opaque consumer identifiers, independent optional assignments, bounded strings and optional declared media timing. Explicit push and official PR authority; human final merge.

## User Scenarios & Testing

### User Story 1 - Retain consumer speaker results (Priority: P1)

A downstream application identifies speakers and stores its chosen identifiers against cues without changing source subtitle labels.

**Why this priority**: Consumers need a reliable storage contract for their own results.

**Independent Test**: Add human-name, UUID and international consumer IDs beside different native labels; validate, serialize and restore source bytes without changing either role.

**Acceptance Scenarios**:

1. **Given** a native label Roger, **When** a consumer assigns interviewer, **Then** both remain independent and the consumer ID survives exactly.
2. **Given** omitted or empty assignments, **When** the document is loaded, **Then** no speaker ID is generated.
3. **Given** repeated and overlapping assignments, **When** the document is saved, **Then** order, spelling and interval values remain intact.

### User Story 2 - Detect invalid speaker timing (Priority: P1)

A consumer supplies intervals and, when known, a finite media duration and subtitle timeline alignment.

**Why this priority**: Catch mistaken intervals without inventing media measurements.

**Independent Test**: A cue ending at620000 and an assignment ending610000 with declared media end600000 reject the assignment; an outlying source cue alone remains restorable with a warning.

**Acceptance Scenarios**:

1. **Given** duration600000 and offset30000, **When** an assignment ends630000 within its cue, **Then** the media check succeeds.
2. **Given** no duration, **When** an otherwise valid timed assignment is loaded, **Then** cue checks apply and media evaluation is unavailable.
3. **Given** an untimed assignment, **When** inspection runs, **Then** it does not claim an evaluated speaker interval.
4. **Given** a cue/media conflict without invalid consumer intervals, **When** validation/restore runs, **Then** the conflict is reported without modifying or losing source bytes.

### User Story 3 - Export without hidden annotation loss (Priority: P2)

A consumer renders or converts an annotated document and receives truthful loss accounting.

**Why this priority**: Native formats cannot encode the new arbitrary consumer identity and declaration surface faithfully.

**Independent Test**: Every matching render and twelve directed conversions warn for omitted assignments/declarations; strict refuses before destination publication.

**Acceptance Scenarios**:

1. **Given** consumer annotations, **When** a native format is requested, **Then** no native source label is fabricated and every omission is reported.
2. **Given** strict export, **When** any annotation is omitted, **Then** no output is created or replaced.
3. **Given** historical documents, **When** the new software loads them, **Then** their exact identity and established behavior remain unchanged.

### Edge Cases

Decoded length boundaries; supplementary Unicode and composed/decomposed differences; all excluded controls escaped/literal; prototype-like and markup/URI/Base64-looking IDs as inert values; boundary whitespace; null versus missing arrays; repeated entries;1024/1025 entries; partial or invalid intervals; signed offsets/checked endpoint overflow; zero media/cue duration; false media verification; stale interval edits; source-only conflict restoration; loss/diagnostic ceilings; exact historical identities.

## Requirements

### Functional Requirements

- **FR-001**: Allow optional consumer-assigned speaker_attributions per common cue without replacing native/heuristic speakers; missing/empty arrays mean no recorded assignment, null is invalid.
- **FR-002**: Treat speaker_id as opaque consumer-governed text; attach no global/universal identity meaning, enforce no per-document scope or registry uniqueness, and generate no IDs.
- **FR-003**: Accept only valid Unicode scalar strings of 1..256 decoded code points; preserve case, combining sequences, supplementary characters, caller order and repeated entries exactly.
- **FR-004**: Reject controls U+0000..001F/U+007F..009F, U+2028/2029, bidi formatting U+061C/U+200E/200F/U+202A..202E/U+2066..2069, U+FEFF and boundary Unicode White_Space; validate decoded escapes with no coercion, repair or hostile value leakage.
- **FR-005**: Close attribution objects to speaker_id plus optional paired start_milliseconds/end_milliseconds; cap each cue at 1024 entries under existing document/input bounds.
- **FR-006**: Support untimed participation without claiming continuous speech; timed intervals use nonnegative checked integer absolute milliseconds, half-open positive intervals within the cue, including overlap and repeated IDs.
- **FR-007**: Reject partial, fractional, reversed, empty, negative, overflowing or out-of-cue intervals and edits invalidating assignments before any output publication.
- **FR-008**: Provide optional closed root media_timing with required nonnegative integer duration_milliseconds and optional signed integer timeline_start_milliseconds, default zero for evaluation while preserving omission.
- **FR-009**: Never infer duration from cues, subtitle bytes, timestamps or unknown headers; supplied duration is a consumer declaration, not independent verification of media.
- **FR-010**: Compute media endpoints with checked arithmetic and require every timed assignment within both cue and declared media intervals; equal end is valid and offset mapping is explicit.
- **FR-011**: Without media_timing keep cue checks and mark media check unavailable; untimed assignments have no evaluated interval; zero duration is known zero, not unknown.
- **FR-012**: Produce non-destructive bounded runtime cue/media conflict diagnostics without altering source cues, persisted document diagnostics/statistics or preventing valid exact restore.
- **FR-013**: Report counts/state only through an optional extension of inspect report v1; no speaker IDs, raw content, or local machine identifiers enter inspection.
- **FR-014**: Preserve all original asset bytes, cue payloads, native labels, tokens and native provenance when assigning consumers; consumer IDs are content values, never automatically executed, fetched, interpreted, decoded or resolved.
- **FR-015**: Account deterministically for each consumer assignment and media declaration omitted by matching native render and all twelve conversions; strict refuses known loss and diagnostics stay bounded.
- **FR-016**: Retain historical1.0.0 and1.1.0 identities and exact released schemas while staging current1.2.0-dev schema/executable lockstep; old versions reject the new fields and identities.
- **FR-017**: Keep native support/capability truth, source integrity, source-aware private conversion targets and current version-specific scripted rules correct under repeated typed revalidation.
- **FR-018**: Publish self-contained schema/model/CLI/conversion/compatibility/architecture guidance and valid annotated examples without overclaiming release or identity/media verification.
- **FR-019**: Provide test-first boundary, Unicode/security, historical, all-format/native/twelve-conversion, source-restoration and meaningful fuzz evidence plus CI-parity/native hosted checks.
- **FR-020**: Drive this complete independently closeable #82 outcome through installed Spec Kit blocking analysis, authorized push/official PR and at most two review rounds; retain human final merge/release/production boundaries.

### Key Entities

- **Consumer attribution**: Consumer ID, optional paired cue-bounded media-timeline interval; independent from source speaker observations.
- **Media timing declaration**: Consumer-supplied finite duration and optional signed alignment of media start to cue timeline.
- **Consumer timing assessment**: Bounded non-content counts, evaluation state and nonfatal cue conflict observations; not persisted source truth.
- **Native omission observation**: Deterministic runtime report for exported consumer data not faithfully represented by the target.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Every accepted assignment in representative four-format documents retains its decoded value/order; exact original asset restoration remains byte-identical.
- **SC-002**: All specified invalid Unicode/interval/type/collection cases reject without repair or partial native publication.
- **SC-003**: Every one of twelve conversion directions and four matching native formats reports consumer omissions; strict creates/replaces no destination.
- **SC-004**: Historical1.0.0/1.1.0 inputs preserve identity and released schema bytes while new output selects exact staged identity.
- **SC-005**: Inspection exposes only counts/states and distinguishes absent or untimed media checks from evaluated timed intervals.
- **SC-006**: Required foreground and hosted CI/security/native checks and review protocol pass before operator merge handoff.

## Assumptions

- The consumer supplies IDs and any real media duration for the correct resource; Cueson does not probe audio or assert identity accuracy.
- Common cue timing is the existing nonnegative absolute millisecond timeline. Media start may be signed to align a different origin.
- No native target faithfully preserves this complete consumer contract; omission reporting is preferable to fabricated native labels.
- Source builds stage1.2.0-dev rather than prematurely publishing a stable version; public immutable routes/releases remain unchanged.
- No new public library API, diarization engine, CLI assignment/import command, registry or uncontrolled metadata container is required.

## Clarifications

### Session 2026-10-06

- Q: Are consumer identifiers restricted to a document? A: No. Consumers own scope; Cueson attaches no universal meaning.
- Q: What if actual media duration is unknown? A: Accept cue-valid assignments and report unavailable media check; never substitute subtitle coverage.
- Q: How do native exports retain arbitrary IDs/media declarations? A: All current targets omit with bounded deterministic reports; strict refuses, restoration remains exact.
- Q: How are supplied duration conflicts with source cues handled? A: Runtime warnings preserve valid source envelopes; invalid consumer intervals fail.
- Q: How is compatibility staged? A: Exact1.2.0-dev current plus local immutable historical1.1.0/1.0.0; no tag/release/production publishing.

