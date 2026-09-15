# S023 native delivery and publication contract

## Scope

P01-P03 are this PR's three closing outcomes. Publish one new v1.1.0 milestone, one coordination epic and sixteen child issues after duplicate/new-arrival inspection. No native-codec implementation, version bump, release or deployment occurs.

## Publication transaction

1. Inspect current open/all matching issues, milestones and Project membership. Reuse matching state and preserve acceptance criteria.
2. Format all UTF-8 six-section bodies through `go run ./scripts/github-format/main.go -stdin`; use file-based publication and immediately compare each body read-back, six headings, checkbox count and plausible line count.
3. Create/reuse the milestone and governed ASS/SSA area labels, publish the epic and children, then resolve P aliases to verified native numbers and numeric IDs.
4. Assign native sub-issues and blocked-by relationships from the acyclic graph in [data-model.md](../data-model.md); never use reciprocal coordination edges or replace an unrelated parent.
5. Add/reuse exactly one Project item per issue, set actual Stage and child Slice, and clear default Status. Preserve all existing Done items.
6. Replace local draft references in published bodies with verified issue links and actual scope/compatibility decisions, format again and compare read-back. Commit traceability snapshots and maintain roadmap links.
7. Audit milestone membership, native parents, exact blocker sets, labels, Project cardinality/Stage/Slice/Status, and no new unassessed issue before publication completion.

## Recovery and acceptance

If a write fails or its response is uncertain, inspect authoritative state before retrying and reuse confirmed writes. No automatic broad deletion or compensating destructive mutation is permitted. A missing permission/native relationship is a failed gate, not a custom-field workaround. Every declared eighteen blocker edges and sixteen sub-issues must read back exactly.

The milestone description records a compatible-addition target and complete release/public-hosting closure, but the final schema compatibility gates govern whether that version may ship. Human PR merge ratifies repository contract prose; current external issues stay open with S023 Stage PR review until merge. Publication snapshots do not claim later implementation or protected transitions complete.

## PR and reviews

Push the S023 branch and create one normal official PR with complete closing references for the three verified S023 issue numbers. Format and read back its body. Inspect every review/comment/reaction and check on the exact head. First Codex review is automatic; eyes acknowledge only; a clean attributable thumbs-up is terminal only with no findings. Address findings with evidence and changes, resolve only handled threads, rerun checks, and request at most one second `@codex review`. No third automatic request, human merge, auto-merge or merge queue.
