# Feature Specification: Extend versioned schema and scripted model

**Feature Branch**: `codex/S024-extend-versioned-schema-model`

**Created**: 2026-09-15

**Status**: Approved scope; blocking analysis passed before implementation

**Input**: Operator kickoff adopts the S023 roadmap group #55 and #56 and authorizes Spec Kit/autopilot implementation, push, official PR, all review remediation, one optional second review, and human final merge handoff.

## Context and Scope

S023 is ratified by merged PR #68 at `0b1922f8a50cdfda766db63012f0d9f7d259146c`. Issues #55 and #56 are Ready with all prerequisites closed. S024 delivers exact historical v1.0.0 input compatibility and the unreleased annotated ASS/SSA schema/model together. Maintained documentation records S023 ratification and the current development contract. Native codecs, corpus execution, scripted conversion, stable promotion, immutable new release copies, tags/releases, production, and final merge remain downstream outcomes.

## Clarifications

### Session 2026-09-15

- Q: Which current identity is appropriate before final candidate promotion? A: Exact 1.1.0-dev with URI https://cueson.io/schema/v1.1.0-dev/cueson.schema.json and matching software version; #65 owns final 1.1.0 promotion.
- Q: How is unimplemented scripted support declared? A: schema_only, ingest/render false, generic restore true, OCR false; installed codecs remain separately unavailable.
- Q: Are nonempty scripted documents without dialogue valid? A: Yes, zero cues with null media summaries and no word timing; historical SubRip/WebVTT minimum-one behavior is unchanged.
- Q: Where are declarations stored and when can references be omitted? A: Declaration fields belong to format_declaration records; only newly constructed recognized style/event objects without captured declarations may omit declaration_id and raw capture observations.
- Q: May pinned packaging/process assumptions change for development proof? A: Only current candidate/version/schema-source assumptions needed for 1.1.0-dev verification, with dated decisions; historical artifacts/evidence and protected publication remain unchanged.

## User Scenarios & Testing

### User Story 1 - Continue using historical documents (Priority: P1)

A user can validate, inspect, restore, render, and convert existing valid v1.0.0 SubRip/WebVTT Cue JSON with the developing executable, without rewriting its identity, producer attribution, or source content.

**Why this priority**: Historical compatibility is a release prerequisite and protects working subtitle workflows before the current model evolves.

**Independent Test**: Run every promised command with tagged historical fixtures and compare accepted/rejected outcomes, identities, native outputs, and restored bytes against the existing contract.

**Acceptance Scenarios**:

1. **Given** a valid v1.0.0 document, **When** each promised command processes it, **Then** exact historical validation runs offline, existing behavior is preserved, and input identity/provenance remain unchanged.
2. **Given** missing, mismatched, malformed, v0.0.0, or unknown identity, **When** Cue JSON is loaded, **Then** it fails before payload publication without native fallback or schema retrieval.
3. **Given** historical invalid structure, semantics, or corrupt source data, **When** validation or output is requested, **Then** the owning boundary rejects it with no destination replacement, including force mode.

### User Story 2 - Express complete scripted native models (Priority: P2)

A schema-aware producer can describe ASS v4+ and SSA v4 native structures and derived common dialogue under the unreleased current contract, with usable annotations and explicit relationships, without being told native codecs already exist.

**Why this priority**: A complete model and schema must precede native ingest/render implementation.

**Independent Test**: Validate positive and negative examples for both native branches, retained/constructed records, dialogue projection, metadata safety, references, and complexity bounds.

**Acceptance Scenarios**:

1. **Given** a conforming current ASS or SSA model, **When** structure and semantics are evaluated, **Then** native sections/records/fields/styles/events/overrides/attachments and common provenance survive typed decoding.
2. **Given** constructed recognized sections/styles/dialogue without capture observations, **When** they are validated, **Then** structured owners are accepted without fabricating raw headers, lines, or timestamps.
3. **Given** wrong native branch, invalid ownership, contradictory projection, unsafe known/unknown metadata, or exceeded bounds, **When** it is evaluated, **Then** validation rejects safely without source disclosure.
4. **Given** a recognized scripted document, **When** native render or conversion is requested, **Then** missing codec capability is reported truthfully; exact restoration uses only verified generic source assets.

### User Story 3 - Discover the actual current contract (Priority: P3)

A user can discover and emit the current output schema and software identity while understanding historical input support, unimplemented scripted codecs, and release status.

**Why this priority**: Exact-version negotiation and truthful support avoid accidental forward-compatibility or stable-release claims.

**Independent Test**: Compare schema command output with canonical bytes and executable version; check new existing-format encodes, maintained documentation, generated site content, and immutable older schema hashes.

**Acceptance Scenarios**:

1. **Given** historical support is installed, **When** schema discovery or new native encoding runs, **Then** it emits the exact current development contract, independently of the loaded historical version.
2. **Given** the new schema/model, **When** consumers review descriptions/examples or support claims, **Then** recognition is provisional, current/released versions are distinct, and historical artifacts and production remain unchanged.

### Edge Cases

- JSON-looking content and .json filenames cannot fall through to native decoding after identity failure.
- Unknown URI references, malformed identity types, duplicate native IDs, orphaned or mixed native branches, attachment-range overlap, integer overflow, excessive occurrences, and unknown metadata safety ambiguity reject.
- Empty dialogue and drawing-only dialogue retain a common cue with one empty logical line; nonempty scripted documents with no dialogue have zero cues and null media summaries, without inventing units.
- Edited common timing does not rewrite capture-time timestamp observations; logical hard/soft-break projections and constructed capture omission follow the ratified contract.
- Historical third-party producer versions do not select a schema or require software lockstep.
- Cancellation, strict losses, overwrite preflight, stderr/stdout separation, quiet/silent modes, and source integrity preserve established behavior.

## Requirements

### Functional Requirements

- **FR-001**: Recognize historical v1.0.0 and one documented current development identity by complete exact URI/version pairs only.
- **FR-002**: Historical structural validation must use locally bundled byte-identical immutable v1.0.0 authority, without network schema loading.
- **FR-003**: Historical semantics retain every released accepted field, capability, native/common relationship, collection rule, and established command behavior.
- **FR-004**: Reject missing, malformed, mismatched, v0.0.0, and unknown identities without relabeling or native fallback.
- **FR-005**: Preserve historical identity, producer, source envelope, and native content through all promised commands.
- **FR-006**: Apply existing source integrity, acquisition bounds, cancellation, safe destination, and no-partial-publication rules to both supported contracts.
- **FR-007**: Current schema discovery, new encode output, official software identity, and embedded current authority remain in lockstep under an explicitly unreleased development identity.
- **FR-008**: Admit only matching ass/ssa native branches with lowercase snake_case fields and complete ordered section/record/declaration/style/event/attachment relationships.
- **FR-009**: Validate unique document-local IDs, contiguous physical ordering, exactly resolving ownership, cue/event/style links, and bounded nonoverlapping attachment record ranges.
- **FR-010**: Preserve native lexical information and unknown accepted content independently of common semantics; constructed recognized objects omit unobserved capture-only fields.
- **FR-011**: Common timing, logical lines, raw_text, plain_text, speaker/token observations, override/drawing spans, and native-derived provenance must follow the ratified ownership contract.
- **FR-012**: Explicit native encoding/timing/Timer, field declarations, metadata privacy/active-content, and model complexity requirements must have positive/negative model evidence.
- **FR-013**: Known and unknown unsafe or unclassifiable authoring metadata rejects the whole operation; explicit dialogue/font-family/embedded-content roles remain content contexts.
- **FR-014**: Every new public field has meaningful machine-readable descriptions and valid portable examples, with recursive annotation coverage.
- **FR-015**: Scripted schema recognition and generic restore must not imply native ingest/render/convert or stable support; unavailable operations fail before publication.
- **FR-016**: Existing SubRip/WebVTT semantics and current capabilities remain stable; scripted zero-dialogue documents may have zero cues without changing historical contracts.
- **FR-017**: Older released schema bytes, hashes, and identities remain unchanged; no current artifact is copied into a released historical identity.
- **FR-018**: Reconcile S023 merge state and describe current output negotiation, historical support, unimplemented codecs, and downstream gates in maintained docs and generated content.
- **FR-019**: Run blocking Spec Kit analysis, test-first behavioral/security validation, complete appropriate foreground verification, and all hosted review/CI gates before handoff.
- **FR-020**: Publish an official issue-closing PR under existing authorization, handle every finding, consume at most one second review, and leave specific human merge and protected release/production actions pending.

### Key Entities

- Exact input contract: complete identity pair, immutable local structural authority, and version-specific semantic authority.
- Current scripted native document: ordered physical occurrences, declared fields, native styles/events/assets, typed lexical values, and references.
- Common dialogue cue: owned timing, deterministic readable logical lines, exact native event payload, and derived provenance.
- Capture observation: optional original header/line/time view, never fabricated for constructed content.
- Capability declaration: supported document operations, distinct from installed native codecs and stable-release evidence.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Both historical native formats pass all five promised command families with unchanged valid/invalid outcomes and byte-exact accepted restoration.
- **SC-002**: All unsupported identity classes reject offline before destination or structured payload publication.
- **SC-003**: Both scripted native shapes pass complete positive/negative relationship, construction, projection, privacy, and bound scenarios without invented source information.
- **SC-004**: All new public fields have descriptions and validating portable examples; both released schema hashes remain unchanged and current schema/software discovery agrees.
- **SC-005**: Verification and review close with no unresolved blocking finding; one official PR independently closes #55 and #56 after human merge.

## Assumptions

- Ratified S023 native/compatibility contracts control the implementation, including second-review privacy narrowing.
- Historical semantics can be isolated from current scripted rules without duplicating a released executable.
- Pure-Go existing toolchain, shared input/output/integrity infrastructure, fixture manifest, and feature branch delivery remain authoritative.
- Authentication and tenancy are absent in this offline CLI; input/privacy/output safety tests are required.
