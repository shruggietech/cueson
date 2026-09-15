# Feature Specification: Harden scripted CLI conformance

**Feature Branch**: `codex/S027-harden-scripted-cli-conformance`

**Created**: 2026-09-15

**Status**: Specified

**Input**: S027 completes issues #62 and #63 together under the installed Spec Kit autopilot workflow, with explicit push/official PR authority and human final merge.

## User Scenarios & Testing

### User Story 1 - Discover and inspect scripted content consistently (Priority: P1)

A user validates and inspects ASS/SSA native files or Cue JSON and receives the same truthful format, schema, capability, and semantic decisions used by encoding, rendering, and conversion. They discover the approved vocabulary through help and all four shell completions.

**Why this priority**: Reliable discovery and inspection make the completed native/conversion workflows usable without privacy or compatibility surprises.

**Independent Test**: Exercise native and Cue JSON ASS/SSA, historical 1.0.0 documents, misleading filenames, aliases, help/completions, and output failures; compare command decisions and bounded summaries.

**Acceptance Scenarios**:

1. **Given** accepted native ASS/SSA and corresponding Cue JSON, **When** validate and inspect run, **Then** selection, schema, integrity, semantic, and installed/declared capabilities agree with other workflows.
2. **Given** scripted documents with retained sections, styles, events, overrides, drawings, and attachments, **When** inspect runs, **Then** useful deterministic counts appear without raw source fields, text, resource names, runtime paths, or local identities.
3. **Given** existing text and historical input contracts, **When** existing commands/help/completions and failure cases run, **Then** established fields, streams, options, and exit classes remain compatible.

### User Story 2 - Reject hostile content at its owning boundary (Priority: P1)

A user supplies complex or malformed scripted content and receives complete bounded rejection or accepted faithful output, with no silent truncation, partial publication, integrity bypass, or source mutation.

**Why this priority**: Source fidelity and safe processing are prerequisites for every supported format.

**Independent Test**: Execute accepted exact-restore/native cycles and adversarial ownership, integrity, unsafe-reference, privacy, limit, and cancellation regressions.

**Acceptance Scenarios**:

1. **Given** accepted governed fixtures, **When** native cycles and restoration run, **Then** every original asset restores byte-for-byte and the established text corpus remains unchanged.
2. **Given** model/source mismatches, corrupt envelopes, unsafe names/references, or path-bearing active content, **When** validation/render/conversion/restoration crosses the owning boundary, **Then** the operation fails with its established error class and no partial output.
3. **Given** content at or beyond collection, field, override, diagnostic, and loss ceilings, **When** processing runs, **Then** permitted content is complete and excess content rejects rather than truncates.

### User Story 3 - Prove the complete development contract (Priority: P2)

A maintainer can trace required scripted contract rows to executable evidence, exercise meaningful bounded fuzz cycles without writing files, and inspect native-platform and six-target build proof.

**Why this priority**: Development conformance is the evidence needed for the later stable-support and release decision.

**Independent Test**: Run the governed matrix/corpus, fixed-work detection/parser/override/field/model-render-reparse/conversion fuzzing, native CI, race/static/security checks, and non-publishing release proof.

**Acceptance Scenarios**:

1. **Given** the #53/#57 required matrix, **When** evidence is audited and executed, **Then** all in-scope rows have truthful runnable evidence and portable inapplicable versus future release gates remain explicit.
2. **Given** bounded fuzz input, **When** each callback runs, **Then** meaningful parser/projection/cycle invariants execute without filesystem writes or unbounded allocation.
3. **Given** the completed slice, **When** CI and external reviews finish, **Then** every required check is green and every actionable finding is answered, fixed where necessary, and resolved within at most two Codex rounds.

### Edge Cases

- Empty documents and native collections, drawing-only/whitespace events, duplicate/reordered declarations, reset/override transitions, malformed attachment spans, and dialect mismatches.
- Misleading extensions and explicit aliases, schema-only observations with installed codecs, exact historical/unknown schema identities, and unsupported encoding/capability.
- Cancellation, nil contexts, collection/line/output ceilings, complete diagnostic/loss ceilings, short stdout writes, diagnostic-write failures, and strict/overwrite refusal.
- Active content or free-form source messages attempting to leak paths, actors, styles, attachment names, source bytes, or local identity into summaries.

## Clarifications

### Session 2026-09-15

- Q: Should scripted counts change old inspection reports? → A: Add an optional counts-only scripted summary for ASS/SSA and preserve existing text report fields and bytes.
- Q: Should known incompatible encoding failures be reclassified? → A: Preserve established owning-command classes (encode/validate/inspect runtime refusal, convert explicit preflight invocation refusal) and test their common capability outcome.
- Q: Does development conformance imply stable release readiness? → A: Execute all development gates while retaining #64/#65 stable identity/documentation and candidate proof ownership.
- Q: How should reviewer-owned checklist markers be handled under autopilot? → A: Generate them unchecked, retain their ownership, record the independent requirements audit, and use the operator's authorized autopilot continuation without changing their markers.

## Requirements

### Functional Requirements

- **FR-001**: Validation and inspection MUST use the shared validated-input and capability decisions for native/Cue JSON ASS/SSA, historical 1.0.0, content-first selection, aliases, semantic checks, and integrity.
- **FR-002**: Inspection MUST expose deterministic useful scripted aggregate counts and truthful capabilities while excluding raw source/text/native values, runtime paths, and local identity.
- **FR-003**: Help and all four completion definitions MUST derive approved input/target vocabulary from one catalogue; existing text/historical report fields, commands, streams, and exit classes MUST remain compatible.
- **FR-004**: Every accepted governed fixture MUST retain exact restoration and meaningful native/common cycles, with the existing text corpus and all twelve conversion directions unchanged.
- **FR-005**: Unsafe references, corrupt envelopes, mismatched native/model ownership, and privacy violations MUST reject at their owning boundaries without altering original source or publishing partial output.
- **FR-006**: Processing MUST preserve complete structures, diagnostics, and losses within existing documented limits and reject excess work rather than truncate; long traversals MUST remain cancellable where their contract requires it.
- **FR-007**: Meaningful bounded fuzz boundaries MUST cover format selection, both parsers, overrides/declared fields, model-render-reparse and conversion cycles without filesystem writes in callbacks.
- **FR-008**: All required #53/#57 development conformance rows MUST have executable evidence; legitimate inapplicable portable claims and separately owned stable/release gates MUST remain truthful.
- **FR-009**: Windows/macOS/Linux native tests, six pure-Go targets, required race/static/security checks, full CLI/corpus tests, documentation/site checks and non-publishing release proof MUST pass.
- **FR-010**: Delivery MUST include complete Spec Kit artifacts, blocking analysis and convergence, formatted/read-back official PR closing #62/#63, current-head green CI, and responses to every external finding with at most two Codex rounds.

### Key Entities

- **Validated input**: A selected document with exact schema, semantic, integrity, and capability evidence shared by commands.
- **Inspection summary**: Existing privacy-bounded report plus optional scripted counts, with no source-value disclosure.
- **Conformance row**: One required capability linked to genuine fixture/test/fuzz/native evidence or an explicitly owned future gate.
- **Fuzz cycle**: Bounded in-memory input and invariants across parsing, model projection, rendering, conversion and target reparse.
- **Delivery evidence**: Task/checklist coverage, exact revision checks, review responses, and Project lifecycle state.

## Success Criteria

### Measurable Outcomes

- **SC-001**: All approved ASS/SSA input tokens and four shell completions agree, with zero failures in the historical/text command compatibility matrix.
- **SC-002**: Scripted inspection counts match accepted document structure, with zero prohibited source/runtime/local-identity disclosures in privacy regressions.
- **SC-003**: Every accepted corpus fixture restores exact bytes, all twelve conversion directions retain outcomes, and every required in-scope matrix row executes.
- **SC-004**: All fuzz callbacks perform zero filesystem writes and each required fixed-work target completes its budget without panic, unbounded accepted work, source mutation, or partial output.
- **SC-005**: All current-head required CI/security/platform/build/release-proof gates pass and zero actionable review findings remain unresolved.

## Assumptions

- #60/#61 closed through merged PR #71 at `771240800aa9ebec35ba3fab19e054cc10a81e74`; a fresh kickoff sees no new active issues.
- #62's closed #61 blocker permits work; #63's #62 prerequisite is satisfied by integrating US1 before final US3 acceptance in this same PR.
- Existing ceilings and privacy boundaries control; no new runtime dependency or format profile is necessary.
- ASS/SSA stay experimental and software/schema remain 1.1.0-dev. #64/#65 own stable/documentation freeze and exact release identity; #66/#67 own separately authorized release/public hosting.
- Additive scripted inspection metadata is present only for scripted inputs so existing SubRip/WebVTT reports retain their established shape.
