# Feature Specification: Complete scripted conversion

**Feature Branch**: `codex/S026-complete-scripted-conversion`

**Created**: 2026-09-15

**Status**: Draft

**Input**: Execute S026 under the installed Spec Kit autopilot protocol, closing GitHub #60 and #61 with ten new conversions and a verified twelve-direction matrix. Push and official pull request publication are expressly authorized; final merge remains with the operator.

## Scope and authority

S025 delivered bounded experimental ASS v4+ and SSA v4 ingest, native models, renderers, restoration, and corpus evidence. S026 completes conversion among SubRip, WebVTT, ASS, and SSA. It includes conversion command selection, help, and completion needed to use these directions. Release-wide CLI/catalog work (#62), broader hardening (#63), stable promotion (#64/#65), release publication (#66), and production site publication (#67) remain subsequent slices. ASS/SSA stay experimental. The constitution, architecture of record, native issue acceptance criteria, and existing released contracts control this slice.

## Clarifications

### Session 2026-09-15 (autopilot)

- Q: How should drawing-only and empty scripted text output behave? A: Fail both modes without inventing text or dropping cues; mixed drawing/readable spans have atomic omission losses; valid empty scripted variants may remain valid.
- Q: Which centisecond policy applies? A: Nearest centisecond with ties upward, checked quotient/remainder arithmetic, one loss per changed endpoint, and fatal collapse/overflow.
- Q: What source-free target defaults and literal control rules apply? A: Deterministic documented Default style and ScriptType/WrapStyle without inferred video dimensions; actual LF/NBSP become controls, literal braces and recognized control-looking pairs are fatal.
- Q: How are dialect differences handled? A: Preserve representable native owners, map alignment explicitly, and atomically account for non-equivalent Layer/Marked, color/alpha roles, and omitted observed dialect fields.
- Q: How can private targets preserve original source truth? A: Use an internal original-plus-target validation/rendering boundary with exact envelope equality; public source validation remains unchanged and retained captures stay source-checked.

## User Scenarios & Testing

### User Story 1 - Export scripted dialogue with an honest loss report (Priority: P1)

A user converts ASS or SSA to SubRip or WebVTT to obtain readable dialogue and usable timing, with a complete account of scripted presentation and other source information the target cannot retain.

**Why this priority**: This independently closes #60 and establishes the loss boundary for the remaining conversions.

**Independent Test**: Convert each scripted dialect to each text target, reparse the result, compare dialogue, timing, and supported emphasis, and compare complete ordered loss reports with portable references.

**Acceptance Scenarios**:

1. **Given** baseline dialogue with timing and shared emphasis, **When** any of the four outbound conversions runs, **Then** the target preserves representable semantics and reparses successfully.
2. **Given** styles, overrides, non-dialogue records, attachments, layout, or other unsupported native information, **When** permissive conversion succeeds, **Then** every known omission or degradation has an atomic stable loss entry.
3. **Given** drawing-only dialogue, **When** conversion is requested, **Then** the documented policy reports or rejects its omission without inventing readable text.
4. **Given** a known loss in strict mode or a fatal representation failure in either mode, **When** conversion runs, **Then** no target payload is published and any existing destination remains byte-identical.

### User Story 2 - Create deterministic scripted targets (Priority: P1)

A user converts SubRip or WebVTT to either scripted dialect, or converts between ASS and SSA, obtaining deterministic native output with explicit defaults and explicit treatment of timing precision and dialect differences.

**Why this priority**: These six directions independently close #61 and complete the four-format conversion graph.

**Independent Test**: Exercise all six directions with golden output, semantic reparsing, repeated conversion, precision boundaries, and native dialect differences.

**Acceptance Scenarios**:

1. **Given** baseline representable dialogue, **When** any of the six conversions runs, **Then** the target preserves expected text and timing and uses documented deterministic script and style defaults where the source has no corresponding information.
2. **Given** timing that cannot be expressed exactly in centiseconds, **When** permissive conversion succeeds, **Then** the documented deterministic quantization policy is applied with a loss for each affected endpoint; strict conversion refuses publication.
3. **Given** style, override, alignment, event, or variant differences, **When** dialect conversion runs, **Then** representable semantics survive and every remaining known loss is reported atomically or conversion fails explicitly.
4. **Given** an interval, text construct, or native value that the target cannot safely represent, **When** either mode runs, **Then** conversion fails before publishing a payload.

### User Story 3 - Rely on one complete, safe conversion matrix (Priority: P1)

A user selects any distinct pair of the four formats and gets a truthful supported, lossy, or fatal result, without regressions to established text conversions or source preservation.

**Why this priority**: Users must trust the advertised support graph and refusal behavior.

**Independent Test**: Run all twelve directions through the conversion interface and command, including strict destination and stdout checks, cancellation, bounds, privacy, historical input, and existing regressions.

**Acceptance Scenarios**:

1. **Given** one of the twelve distinct pairs, **When** conversion is requested, **Then** documentation and executable behavior agree on its baseline and lossy/fatal edges.
2. **Given** an existing SubRip/WebVTT conversion fixture, **When** conversion runs after S026, **Then** established output and loss code meanings remain compatible.
3. **Given** source assets, historical supported Cue JSON, or private caller paths, **When** conversion runs, **Then** original source bytes and integrity metadata remain unchanged and no private identity enters target reports or Cue JSON.

### Edge Cases

- Empty scripts, drawing-only cues, mixed readable/drawing content, comments, unknown records, and preserved malformed native records must have explicit outcomes.
- Centisecond ties, intervals that collapse after quantization, maximum timestamps, and integer overflow must be handled deterministically without silent truncation.
- Literal braces, backslashes, line breaks, markup, entities, and scripted controls must not silently acquire different meanings.
- Duplicate declarations, missing styles, attachments, actor/effect fields, dialect-only columns, styles, and override groups require complete accounting.
- Loss amplification, excessive records or payloads, cancellation, invalid source integrity, and renderer rejection must fail boundedly before publication.

## Requirements

### Functional Requirements

- **FR-001**: Support ASS to SubRip, ASS to WebVTT, SSA to SubRip, and SSA to WebVTT, preserving baseline readable dialogue, timing, and target-supported shared emphasis.
- **FR-002**: Support SubRip to ASS, SubRip to SSA, WebVTT to ASS, WebVTT to SSA, ASS to SSA, and SSA to ASS with deterministic target output.
- **FR-003**: Preserve established SubRip to WebVTT and WebVTT to SubRip behavior and existing stable loss code meanings.
- **FR-004**: Account atomically for every known omitted or degraded source feature using a stable code, portable structural reference, and deterministic ordering.
- **FR-005**: Document and implement an explicit drawing-only policy without inventing readable text or silently discarding events.
- **FR-006**: Document deterministic scripted script/style defaults and distinguish target construction from source-native provenance.
- **FR-007**: Apply a documented deterministic centisecond quantization policy, reporting each changed endpoint; reject invalid or unsafe target intervals.
- **FR-008**: Preserve representable native semantics across ASS/SSA and report or reject every known style, override, alignment, event, or dialect difference.
- **FR-009**: Strict mode must refuse every known lossy conversion before publication and return its complete report when analysis succeeds.
- **FR-010**: Fatal failures in either mode must publish no target payload, preserve existing destinations, and avoid partial stdout payloads.
- **FR-011**: Target output must reparse through its registered codec and agree with expected semantic projections.
- **FR-012**: Source assets, lengths, hashes, and restoration truth must remain unchanged; conversion must not claim restoration or losslessness for degraded output.
- **FR-013**: Selection, reports, diagnostics, and constructed model fields must not expose original filesystem paths or local machine identifiers.
- **FR-014**: Bound conversion work, target allocation, native record expansion, and report growth; cancellation and malformed inputs must fail without panic.
- **FR-015**: Preserve supported historical Cue JSON input and immutable released schema artifacts; use the current development contract when constructing private scripted targets.
- **FR-016**: Conversion selection, help, and completion must expose implemented directions while preserving same-format and unsupported-format error behavior.
- **FR-017**: Publish a verified twelve-direction matrix with baseline, lossy, and fatal expectations, including explicit unsupported edges within those directions.
- **FR-018**: Verify goldens, loss ordering, target reparsing, precision boundaries, defaults, strict destination/stdout safety, privacy, bounds, and old regressions.
- **FR-019**: Keep scripted codecs experimental and exact development schema/software identity unchanged; do not promote, tag, release, or publish production resources.
- **FR-020**: Complete blocking Spec Kit analysis, foreground verification, official PR publication, CI, and every review finding; request at most one second Codex review and leave merge to the operator.

### Key Entities

- **Conversion direction**: An ordered distinct source/target pair with documented baseline and lossy/fatal edges.
- **Loss entry**: One source omission or degradation with a stable code and portable structural reference, without source excerpts or private identity.
- **Constructed target**: Deterministic native output derived from the source model, preserving the immutable source envelope and distinguishing defaults from source observations.
- **Conversion evidence**: Expected bytes, reparsed semantics, ordered reports, and publication-refusal checks.

## Success Criteria

### Measurable Outcomes

- **SC-001**: All twelve distinct directions pass documented baseline semantic and target-reparse checks.
- **SC-002**: Every known omission or degradation in governed fixtures appears in its complete ordered report, with zero silent losses.
- **SC-003**: Every governed strict or fatal case publishes zero target payload bytes and leaves existing destination bytes unchanged.
- **SC-004**: Repeated scripted target baseline conversions produce byte-identical output and reports; all governed precision boundaries have explicit outcomes.
- **SC-005**: Existing text fixtures and released schema integrity checks pass unchanged, and all CI checks and external findings are satisfied before human merge handoff.

## Assumptions

- #59 is closed by merged S025; #61 depends on #60 and follows shared outbound analysis inside this coherent slice.
- Existing preservation, reporting, atomic publication, and registered codec contracts remain authoritative.
- Routine choices within these acceptance criteria are resolved under autopilot and recorded in clarification/design artifacts; authorized implementation, push, and PR publication need no additional approval.
- Target defaults are not lost source information when no corresponding observation exists; replacing observed values requires accounting.
