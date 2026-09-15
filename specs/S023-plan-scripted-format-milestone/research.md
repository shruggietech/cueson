# S023 research decisions

## Milestone and scope

**Decision**: Adopt v1.1.0 ASS/SSA as a compatible-addition target, with sixteen atomic children and three S023 planning outcomes.

**Rationale**: The operator adopted this proposed pathway at kickoff; shared scripted sections/styles/events give a coherent native preservation and conversion story. The main draft lists families without requiring XML-first ordering.

**Alternatives considered**: TTML/IMSC first introduces namespaces/profile/timing inheritance; bitmap first requires binary image structures and OCR. Those remain later milestones rather than being removed from the roadmap.

## Current compatibility baseline

**Decision**: Promise future 1.1.0 processing of current 1.0.0 documents with exact local historical schema dispatch and version-specific semantics. Keep 0.0.0 unsupported in the new executable; its historical executable remains available.

**Rationale**: Current code compiles one exact 1.0.0 schema and fixed current model semantics. It does not accept 0.0.0. Reviving its adapter is optional new scope, not necessary for the v1 promise.

**Alternatives considered**: Merely changing constants loses old inputs; loose semver acceptance ignores exact identities; network retrieval changes trust; adding 0.0.0 support expands this milestone without an existing compatibility requirement. Independent compatibility research confirms the actual baseline.

## Native grammar authority

**Decision**: Define a bounded Cueson ASS v4+ / SSA v4 acceptance contract using maintained Aegisub/libass documentation and pinned implementation sources, not an implicit promise to duplicate pixel-renderer quirks.

**Rationale**: Maintained libass guidance is generation guidance rather than exhaustive parsing authority. Ordered raw/native retention and derived dialogue interpretation need explicit acceptance rows. Timer other than 100 is rejected rather than inventing common timing or redesigning the stable cue shape.

**Alternatives considered**: Treat every variant as compatible, strip override content to plain text, or rely on inaccurate early grammar sketches. Independent native-contract research supplies source-linked boundary decisions.

## Privacy and field ownership

**Decision**: Retain accepted source bytes/native unknown occurrences; reject whole encode when known path/reference metadata cannot be represented under portable/privacy rules. Never silently sanitize raw source fields or fetch/execute attachments. Native editable event text owns its projection; parsed/common views must agree rather than stale raw records overriding edits.

**Rationale**: Constitution and AGENTS bind source fidelity and privacy simultaneously. Safe rejection is proportional and does not alter the source or create lossy Cue JSON.

**Alternatives considered**: Deleting metadata violates no-silent-loss; exporting filesystem paths violates privacy; opaque-only native content fails usable common semantics; arbitrary file/network resolution expands trust.

## Native GitHub relationships

**Decision**: Publish six-section file bodies, native milestone/parent/dependencies, governed labels and exactly one Project item per issue. Dependency edges use hard prerequisites only and are transitively reduced; coordination never introduces cycles.

**Rationale**: GitHub and the Project contract own those facts. Current primary REST docs provide sub-issue and blocked-by endpoints with numeric issue IDs and separate read-back.

**Alternatives considered**: Checklist-only epics and custom dependency/milestone fields duplicate native metadata. See [GitHub sub-issues API](https://docs.github.com/en/rest/issues/sub-issues) and [dependencies API](https://docs.github.com/en/rest/issues/issue-dependencies).

## Verification and authority

**Decision**: Use existing deterministic text/docs/Go/site checks and direct native-state audit; add no runtime test/tool changes to this documentation/planning slice. Run verification in the foreground with the previously verified hidden Windows launcher. Publish push/PR under explicit kickoff authority, address reviews, and stop before human merge.

**Rationale**: Changes affect maintained docs and external planning state. Existing behavioral and security checks remain intact. No tag/release or production authority is inferred.

**Alternatives considered**: Duplicate text-only tests mirror the prose without proving native state; skipping generated-site checks misses adaptation issues. Custom quality checklist markers remain reviewer-owned, with separate recorded assessment under the operator's autopilot authorization.
