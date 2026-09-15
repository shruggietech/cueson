# Feature Specification: S025 ASS/SSA native workflows

**Feature Branch**: `codex/S025-deliver-ass-ssa-native-workflows`

**Created**: 2026-09-15

**Status**: Implemented; PR/CI/review delivery in progress

**Input**: Implement #57, #58 and #59 together under Spec Kit autopilot, preserving their acceptance criteria and completing review/CI delivery with explicit push and PR authority.

## Scope and authority

Deliver the provenance corpus, content-first native ingest and model-driven native rendering for the ratified bounded ASS v4+/SSA v4 profile. Historical 1.0.0 behavior and immutable released schema bytes remain intact. Current software/schema stays 1.1.0-dev. This slice does not implement cross-format conversion (#60/#61), complete CLI discovery/inspection enhancements (#62), release-wide conformance closure (#63), stable promotion (#64/#65), publication (#66) or production hosting (#67). Necessary existing encode/render/shared-input integration and truthful codec reporting are included so native operations work end-to-end.

## Clarifications

### Session 2026-09-15

- Q: What maturity should implemented native codecs report before the stable gate? A: Experimental with ingest/render/restore true and OCR false; existing schema_only models remain accepted as truthful input observations. Stable is deferred to S028. This avoids rewriting older model declarations merely because the executable now has codecs.
- Q: How should ambiguous or unsupported scripted content interact with filename fallback? A: Reject recognized scripted candidates explicitly, including explicit mismatched dialect requests, instead of selecting another codec by extension.
- Q: How should source content before the first section be represented? A: Reject preamble content rather than fabricate section/capture identity; leading empty physical lines are also rejected under this closed profile.
- Q: How should safely retained malformed declarations or non-dialogue records render? A: Ingest may retain bounded inert malformed content with diagnostics; every renderer mode refuses such preservation-only models. Dialogue ambiguity is fatal at ingest.
- Q: How should constructed record declaration transitions behave? A: Insert an accepted Format matching complete structured owner order when needed, and restore the required prior declaration before following captured owners; reject unresolved competing owners or unsupported transitions.

## User Scenarios & Testing

### User Story 1 - Trust the scripted fidelity corpus (Priority: P1)

Maintainers can trace every scripted acceptance row to a provenance-verified fixture and an honest implemented or pending assertion.

**Why this priority**: Independent byte and semantic expectations prevent two mutually incorrect codecs from passing a cycle.

**Independent Test**: Verify the shared inventory, byte contracts, hashes, redistribution evidence and all scripted row references without executing future conversion assertions.

**Acceptance Scenarios**:

1. **Given** an authored accepted ASS/SSA fixture, **When** the inventory is verified, **Then** exactly one manifest entry accounts for each payload with byte length/hash, encoding, line endings and redistribution evidence.
2. **Given** the acceptance matrix, **When** row evidence is resolved, **Then** native retention, common projection, diagnostics and applicable cycles have distinct expectations; downstream conversion/release assertions remain explicitly pending.

### User Story 2 - Encode native scripts without losing source truth (Priority: P1)

Users encode accepted scripts into valid Cue JSON with readable timed dialogue, retained native content and exact original assets.

**Why this priority**: Ingest establishes both usable semantics and authoritative restoration independently of rendering.

**Independent Test**: Select each dialect from content, encode corpus inputs, validate structure/semantics/integrity and restore accepted inputs byte-for-byte.

**Acceptance Scenarios**:

1. **Given** valid content with a conflicting filename extension, **When** auto encoding selects a format, **Then** content selects the matching dialect; explicit aliases are tested without bypassing dialect checks.
2. **Given** reordered/extended declarations, repeated sections, dialogue commas, styles, non-dialogue records, drawings, overrides, karaoke or embedded assets, **When** accepted encoding completes, **Then** ordered native content and source observations remain faithful and common dialogue/speaker/token projections follow the ratified rules.
3. **Given** unsafe metadata, active scripts, unsupported encodings/dialects, ambiguous dialogue or exceeded bounds, **When** encoding runs, **Then** the operation rejects safely with deterministic bounded diagnostics and publishes no payload.
4. **Given** accepted empty or comment-only scripted content, **When** encoding completes, **Then** its nonempty native document and zero common cues remain valid and restorable without invented dialogue.

### User Story 3 - Render model edits deterministically (Priority: P1)

Users render accepted validated models into reparsable native syntax, with structured edits controlling output and restoration retaining original bytes.

**Why this priority**: Replaying source bytes would silently ignore edits and cannot establish a native codec.

**Independent Test**: Render both dialects repeatedly, re-ingest to compare required semantics/retention and check edited and newly constructed models independently of exact restoration.

**Acceptance Scenarios**:

1. **Given** an unedited renderable accepted model, **When** rendering/re-ingest completes, **Then** dialect, native occurrence order, declared fields, styles/events/assets and common timing/text/speakers/tokens remain equivalent under the canonicalization contract.
2. **Given** consistent edits to common timing or owning native fields/Text, **When** rendered, **Then** output contains those edits while capture-only observations and the original source envelope remain unchanged.
3. **Given** newly constructed recognized sections/styles/dialogue without capture fields, **When** rendered, **Then** complete structured owners produce deterministic valid native syntax without fabricated observations.
4. **Given** malformed retained syntax, unresolved ownership, unsafe references, incompatible precision or exceeded output bounds, **When** permissive or strict rendering runs, **Then** fatal cases refuse before output publication; strict additionally refuses every known conformance ambiguity.

### Edge Cases

- Missing/mixed/duplicate contradictory ScriptType or style dialect signatures; rejected candidates must never fall through to another codec by filename.
- UTF-8 BOM and LF/CRLF accepted source bytes versus canonical BOM-free LF output; invalid UTF-8, NUL and bare CR reject.
- Complete reordered declarations, duplicate unknown fields, nonfinal or duplicate Text, malformed retained comments/styles and invalid mandatory scalar fields.
- Empty text, drawing-only dialogue, overlapping cues, native positive intervals and centisecond precision/overflow.
- Relative/absolute/URI/identity metadata in known and unknown sections or fields; content roles are established by valid grammar before privacy exemptions.
- Encoded attachments, duplicate safe names, invalid encoded data, record ownership/ranges and decoded bounds; no extraction or resource execution.
- Native/common edits with stale source captures, derived-only edits, inserted separators, missing capture-only fields in constructed owners and deterministic unknown-override retention.
- Acquisition, record/line/declaration/aggregate/span/depth/diagnostic/attachment/output limits reject without truncation; no partial destination/stdout on fatal or strict refusal.

## Requirements

### Functional Requirements

- **FR-001**: Corpus payloads MUST reuse the existing manifest with unique inventory ownership, provenance, redistribution, byte/hash and encoding contracts.
- **FR-002**: Every scripted acceptance row MUST link fixture IDs to separate native/common/diagnostic and applicable operation assertions; future conversion/release gates MUST remain pending.
- **FR-003**: Content-first detection MUST distinguish ASS v4+ and SSA v4, test canonical/alias/explicit selectors and extension conflicts, and reject ambiguous/unsupported scripted candidates without fallback.
- **FR-004**: Ingest MUST use one bounded exact source acquisition and preserve original source bytes/metadata/integrity without serializing runtime paths or local identity.
- **FR-005**: Ingest MUST retain ordered sections, declarations, fields, style/event occurrences, unknown safe records, overrides, drawings, embedded assets and observed native timing.
- **FR-006**: Common dialogue text/lines, native speaker observations and eligible karaoke tokens MUST derive from retained native owners using the ratified projection; drawing spans MUST never become readable text.
- **FR-007**: Valid documents MUST pass current structural, semantic and source-integrity checks with truthful experimental codec capabilities; historical input semantics MUST remain unchanged.
- **FR-008**: Privacy/active-content classification MUST cover original and edited views, unknown metadata/fields and grammar-dependent content roles, reject unsafe whole operations and never fetch/execute/extract resources.
- **FR-009**: Diagnostics MUST have closed identifiers, severity, stable encounter order and safe positions, preserve all applicable findings within existing ceilings and never echo unsafe values.
- **FR-010**: Every accepted source fixture MUST restore byte-for-byte independently of native rendering, including diagnosed retained models that cannot render.
- **FR-011**: Rendering MUST serialize structured owners and common timing, honor consistent edits, retain uninterpreted safe content and never replay stale captures over edits.
- **FR-012**: Constructed recognized content MUST render from complete structured owners without fabricated capture fields; ambiguous or inconsistent ownership MUST reject.
- **FR-013**: Canonical output MUST be deterministic UTF-8 without BOM, LF with final LF, accepted field/order retention and exact representable centisecond timing; incompatible precision MUST reject.
- **FR-014**: Render/reparse equivalence MUST compare required native/common semantics and retention using position-based identity remapping while excluding canonical whitespace/numeric/capture-only differences.
- **FR-015**: Malformed retained records/attachments, ambiguous declarations, unsafe models and fatal limits MUST refuse rendering in all modes; strict additionally MUST refuse known conformance ambiguity before publication.
- **FR-016**: Existing acquisition/input/collection/diagnostic limits and scripted line/declaration/aggregate/span/depth/attachment limits MUST remain enforced; canonical output MUST not exceed 64 MiB and no ceiling permits truncation.
- **FR-017**: Existing encode/render/shared-input command boundaries MUST expose the implemented native workflows with established streams, overwrite behavior, error classes and publication safety, without advertising unimplemented conversion/stable support.
- **FR-018**: Existing SubRip/WebVTT and historical 1.0.0 command behavior, schema hashes, six pure-Go targets and native platform gates MUST remain valid.
- **FR-019**: Documentation and schema annotations/examples MUST match experimental native capability and edit/fidelity restrictions; dated decisions and slice evidence MUST record departures from S024 schema-only state.
- **FR-020**: Blocking analysis, meaningful native/security regression tests, foreground CI parity and every review finding MUST converge before final human review; automated Codex review MUST not exceed two rounds.

### Key Entities

- **Corpus fixture**: Unique source/model/diagnostic/byte artifacts with independent expected observations and portable provenance.
- **Scripted native document**: Dialect, ordered sections/physical records, active declarations, styles/events, assets and capture-only source observations.
- **Common dialogue**: Editable positive timing plus consistently derived native Text, lines, speaker provenance and karaoke tokens.
- **Native operation result**: Validated model or canonical bytes with bounded ordered diagnostics, independent from exact source restoration.

## Success Criteria

### Measurable Outcomes

- **SC-001**: 100% of required scripted acceptance rows have resolvable fixture/operation ownership with no speculative future assertion marked passed.
- **SC-002**: 100% of accepted corpus sources validate and restore byte-identically; both variants independently match their native/common expected projections.
- **SC-003**: 100% of renderable native baseline/edit/construction fixtures reparse to their expected required semantics and deterministic bytes.
- **SC-004**: All declared malformed, privacy, active-content, precision and ceiling boundaries produce their expected safe refusal/diagnostics without partial publication.
- **SC-005**: Existing format/historical regressions and required local/hosted verification pass, with every actionable review resolved and no more than two Codex rounds.

## Assumptions

- The S023 ratified ASS/SSA profile and S024 typed model are the starting authority; discovered contradictions are resolved proportionally and documented.
- Software/schema identity remains 1.1.0-dev; native ASS/SSA support is experimental until the later release-wide stable gate.
- #57 precedes #58 and #58 precedes #59 within this same slice; no reciprocal dependency or premature issue closure is introduced.
- Push and official PR publication are explicitly authorized for S025. The human retains final merge; tag/release/production actions are separately protected.
