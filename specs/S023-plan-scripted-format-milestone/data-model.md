# S023 planning entities and native relationships

This slice models planning/contracts, not new shipped Go/JSON types. Native fact records are read-back snapshots; GitHub remains the planning authority.

## Milestone

Native title v1.1.0; scope ASS/SSA native workflows plus v1 historical compatibility, verified official release and immutable public hosting. Remains open until every independently verified child completes. No due date is invented.

## Epic and atomic outcomes

One coordination epic has sixteen native sub-issues. P01-P16 are traceability aliases resolved to actual issue numbers/IDs after publication. Every child has exactly the six required Markdown sections, native milestone and governed labels, individual verification and one Project item. The epic's Slice remains empty because no single execution slice owns the whole release.

## Hard blocker graph

| Child | Direct native blocked-by outcome |
|---|---|
| P01 | Completed issue #49 |
| P02 | P01 |
| P03 | P01 |
| P04 | P03 |
| P05 | P02, P03 |
| P06 | P05 |
| P07 | P04, P06 |
| P08 | P07 |
| P09 | P08 |
| P10 | P09 |
| P11 | P10 |
| P12 | P11 |
| P13 | P12 |
| P14 | P13 |
| P15 | P14 |
| P16 | P15 |

This acyclic graph contains eighteen direct edges, including the completed S022 prerequisite. Same-slice outcomes can be implemented and closed together; an open native prerequisite still prevents advertising a downstream execution slice Ready before the grouped PR merges. P02/P03 and P04/P05 coordinate without reciprocal edges. Transitive prerequisite coverage is preserved.

## Project state

Children have owning Slice text S023 through S030 according to the approved roadmap; later codes are provisional until their kickoff. S023 children use In progress then PR review while this slice runs; future children remain Backlog while native blockers are open. The epic uses In progress and an empty Slice. Default Status is explicitly cleared after automation writes it. All preexisting Done items are retained without duplicate membership.

## Contract acceptance row

A row has a unique local identifier, primary authority/source reference, source classification, required native retention, common derived semantics or explicit rejection, deterministic diagnostic/failure behavior and future issue owner. Future row evidence is not falsely marked passed in S023.

## Version matrix

Exact current/historical identities select their local immutable schemas and matching semantic validators. Producer version is independent. Processing historical input preserves its identity/provenance; new encode output uses the current contract. Old exact-version consumers rejecting new identities is documented rather than silently relaxed.
