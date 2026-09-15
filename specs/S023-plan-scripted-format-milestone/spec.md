# Feature Specification: Plan the scripted-format milestone

**Feature Branch**: `codex/S023-plan-scripted-format-milestone`

**Created**: 2026-09-15

**Status**: Approved scope; blocking analysis and preparation verification passed

**Input**: Operator kickoff authorizes S023 through installed Spec Kit/autopilot, automatic push and official PR creation, complete review remediation with at most one second review, and a final human merge handoff.

## Context and Scope

S022 is merged as `5c75deee710e89349689789a5f854deb091069f4`; all existing issues and both prior milestones are closed. The kickoff adopts the proposed ASS/SSA milestone and P01-P16 draft backlog. S023 delivers P01 (governed roadmap and delivery reconciliation), P02 (scripted native contracts), and P03 (minor-release compatibility contract), assigned real independently closeable GitHub issues during publication.

Scope includes the next milestone/epic, sixteen atomic children, native parent/dependency/Project metadata, maintained roadmap/contracts, source-linked acceptance matrices and current-state documentation. Runtime implementation, a version bump, immutable-schema modification, tags/releases, production deployment and final merge are excluded.

## Clarifications

### Session 2026-09-15

- Q: Does kickoff authorize publishing the proposed native backlog? A: Yes, the S023 scope explicitly includes the milestone, epic and sixteen child outcomes; reuse any matching new issues.
- Q: Does v1.1.0 mean old exact-version consumers automatically accept new documents? A: No. Preserve promised historical input workflows on the new executable, identify new output exactly, and document old-consumer rejection of an unknown schema identity.
- Q: Does this slice implement or advertise ASS/SSA codecs? A: No. It delivers future contracts and planning; shipped version/capabilities and immutable artifacts stay unchanged.
- Q: May original-source fidelity override portable/privacy and external-asset safety? A: No. Define safe accepted-source retention and explicit whole-operation rejection where a source cannot be represented under the contract; never silently sanitize raw fields or fetch dependencies.
- Q: Do custom checklist markers authorize a stop despite autopilot? A: No. They stay reviewer-owned; the existing kickoff authorizes continuing after a recorded quality assessment and clean blocking analysis, without falsely checking reviewer markers.

These answers were selected under autopilot from operator scope, constitutional source/privacy rules and the current exact-version validation contract. No unanswered critical ambiguity remains; dialect details are resolved by primary-source planning research.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Select a complete next milestone (Priority: P1)

The operator can see every independently testable outcome through the next verified release and public hosting, with complete dependencies and no duplicate issues.

**Why this priority**: There is no open delivery queue, and a codec-only roadmap would omit compatibility and release closure.

**Independent Test**: Audit the native milestone, epic, sixteen children, acyclic dependency graph, Project membership and six-section bodies against the chronological roadmap.

**Acceptance Scenarios**:

1. **Given** closed prior milestones, **When** S023 publishes planning, **Then** one epic and sixteen atomic outcomes have native metadata and exactly one Project item each.
2. **Given** a matching or new issue arrives, **When** publication is reconciled, **Then** its acceptance criteria are retained and matching outcomes reused.
3. **Given** future families in the main draft, **When** the roadmap is reviewed, **Then** ASS/SSA prioritization is justified and XML, SAMI, broadcast and bitmap/OCR remain explicit later milestones.

### User Story 2 - Implement native fidelity from a clear contract (Priority: P1)

A future implementer can distinguish accepted ASS/SSA structures, retained unknown content, dialogue semantics, model rendering, exact restoration and unsupported visual behavior.

**Why this priority**: Plain-text flattening would violate native fidelity and misstate stable support.

**Independent Test**: Trace source-linked native acceptance rows to field ownership, common projection, security/privacy limits, diagnostics and future verification owners.

**Acceptance Scenarios**:

1. **Given** declared fields, dialogue commas, overrides, karaoke or drawings, **When** the contract is applied, **Then** native retention and supported common interpretation are explicit.
2. **Given** attachments, external references or source-path metadata, **When** retention is specified, **Then** original bytes remain preserved without arbitrary external fetch/execution or captured runtime identity leaks.
3. **Given** an edit conflicts with raw source views, **When** rendering is planned, **Then** field ownership and safe rejection prevent stale raw text overriding the edit.

### User Story 3 - Upgrade without losing historical behavior (Priority: P1)

A consumer can determine exact historical input support, promised CLI workflows and current output identity.

**Why this priority**: Changing a version constant does not prove minor compatibility.

**Independent Test**: Review the command/version matrix for historical inputs, current output, unsupported identities, provenance, integrity and existing CLI/loss behavior.

**Acceptance Scenarios**:

1. **Given** a valid v1.0.0 document, **When** a promised future workflow processes it, **Then** its original identity/provenance and owning immutable contract remain authoritative.
2. **Given** mismatched or unsupported identity, **When** input classification occurs, **Then** explicit rejection precedes publication.
3. **Given** an incompatible proposal, **When** the release gate is assessed, **Then** a new major-version decision precedes implementation.

### Edge Cases

- Partial native publication must inspect/reuse existing writes rather than duplicate issues/items.
- Unavailable native relationships must not be replaced by duplicate Project fields or falsely reported verified.
- New issues require explicit inclusion or exclusion at handoff.
- Historical summaries must reconcile current facts without erasing chronological evidence.
- Original source can contain path-like data; distinguish source truth from captured runtime identity without silent deletion.
- Schema recognition before codecs must not advertise native capability.
- Custom requirements-quality checklist markers belong to the reviewer and do not represent implementation completion.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Target compatible v1.1.0 scripted delivery; require a major-version decision for unavoidable public breaks.
- **FR-002**: Publish one coordination epic and sixteen independently closeable children, with P01-P03 assigned S023.
- **FR-003**: Every atomic issue MUST have the required six sections and preserve accepted draft/new issue criteria.
- **FR-004**: Use native milestone, parent/sub-issue, blockers and governed labels without duplicate custom metadata.
- **FR-005**: Each issue MUST appear exactly once in cueson Delivery with approved Stage, owning Slice and unused default Status.
- **FR-006**: Hard blockers MUST form an acyclic graph; coordination is not a reciprocal dependency and protected release/deployment gates remain distinct.
- **FR-007**: Published bodies MUST pass the publication formatter and immediate content/rendering read-back.
- **FR-008**: Inspect issues immediately before publication and final handoff, and account for new arrivals.
- **FR-009**: Sequence contracts, native workflows, conversion, hardening, candidate, release and public hosting chronologically.
- **FR-010**: Identify every future family in the main draft, retain OCR for complete bitmap semantics and justify ASS/SSA-first ordering.
- **FR-011**: Reconcile completed v1 release/milestone closure, S021 hosting and S022 verification in maintained current-state prose without erasing historical evidence.
- **FR-012**: Pin native dialects, keys/aliases, encoding/timing, ordered sections/declared fields, styles/events and unknown-occurrence retention.
- **FR-013**: Separate original assets, native interpretation, derived dialogue/speakers/karaoke, drawings and visual behavior.
- **FR-014**: Define edit ownership, deterministic rendering, reparsing equivalence and safe failure for ambiguous/invalid models.
- **FR-015**: Preserve source bytes, prohibit arbitrary external fetch/execution and prevent captured runtime paths/machine identifiers in Cue JSON.
- **FR-016**: Define deterministic malformed/unrepresentable behavior, diagnostics, complexity bounds and future verification owners.
- **FR-017**: Define exact historical identities, promised commands, local schema dispatch, version-aware semantics and mismatched/unsupported rejection.
- **FR-018**: Preserve historical identity/provenance unless an explicitly specified transformation produces a new document.
- **FR-019**: Separate official current software/schema lockstep, historical input support and third-party producer versions.
- **FR-020**: Distinguish schema recognition, ingest, model render, exact restore and stable gates; planning MUST NOT advertise installed scripted capability.
- **FR-021**: Plan twelve distinct-format conversion directions, complete losses, precision and strict/fatal publication safety, or explicitly revise scope before a narrower stable claim.
- **FR-022**: Source references, acceptance rows and examples MUST trace to future schema/corpus/codec/conversion/CLI verification outcomes.
- **FR-023**: Deliver UTF-8 without BOM, governed line endings, unwrapped prose and valid local links.
- **FR-024**: Pass blocking analysis and integrated verification before commit/publication success; handle all findings with at most two automated review rounds.
- **FR-025**: Leave runtime behavior, current version, immutable artifacts and production unchanged; return before human final merge.

### Key Entities

- **Milestone**: Native release coordination with independently verified closure gates.
- **Atomic outcome**: Six-section issue with acceptance, verification, native relationships and one Project item.
- **Work slice**: Coherent issue group sharing implementation/review/verification and an approved code.
- **Native contract**: Accepted grammar, retention/edit ownership, projection, diagnostics and limits.
- **Version matrix**: Exact document identity and command expectations distinct from producer software version.
- **Acceptance row**: Source-linked boundary classification and future verification owner.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: One epic and sixteen children are read back with correct native metadata, exactly one item each, zero dependency cycles and unused default Status.
- **SC-002**: All twenty-five requirements have task/verification coverage and zero unresolved critical/high analysis findings.
- **SC-003**: Every FR-012 through FR-022 boundary has an explicit acceptance row, authority/source reference and future issue owner.
- **SC-004**: Current summaries no longer describe completed v1 closure/hosting as pending; historical evidence stays chronological.
- **SC-005**: Published/repository/generated docs pass applicable formatting, encoding, links and consistency checks; required hosted checks and reviews pass before handoff.
- **SC-006**: Final diff contains no runtime implementation, version bump, immutable change or production mutation; human merge remains outstanding.

## Assumptions

- The S023 kickoff adopts the proposed milestone and backlog publication; human PR merge ratifies repository contracts.
- Constitution and maintained current contracts control over provisional draft prose.
- v0.0.0 support is not inferred merely because its immutable release exists; the compatibility contract must state its actual baseline.
- Later slice codes remain provisional and are reassessed after each merge.
