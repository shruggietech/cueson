# Project management

GitHub is Cueson's planning and delivery source of truth. The canonical repository is `shruggietech/cueson`, and the organization Project is `cueson Delivery`.

## Issue contract

Every actionable issue owns one independently closeable and independently verifiable outcome. Its body contains Outcome, Context, Scope, Acceptance criteria, Dependencies, and Verification. Native assignees, milestones, parent/sub-issues, and dependencies own those facts; they are not duplicated in custom Project fields.

## Labels

Labels describe governed type, priority, effort, area, and exceptional gates. Families prefixed with `type:`, `priority:`, and `effort:` are mutually exclusive on an issue. Workflow state belongs only in the Project `Stage` field, and release identity belongs only in milestones.

## Project fields

`Stage` uses exactly Backlog, Ready, Specced, In progress, PR review, Release verification, and Done. `Slice` identifies a coherent cross-issue implementation batch. GitHub's default `Status` field remains unused.

Every in-scope repository issue appears exactly once in the Project. Closed issues use Done. Open issues awaiting release or platform evidence use Release verification. Automation must read mutations back and must not infer unknown planning values.

## Work slices

A work slice may close multiple atomic issues when they share a clear purpose, bounded change surface, review story, and verification boundary. Issue count alone does not define slice size. Each included issue retains its own acceptance criteria and closure evidence.

## Pull requests

Normal pull requests contain at least one complete closing reference. Eligible non-Dependabot pull requests use the two-round Codex review protocol in `AGENTS.md`. Final merge authority remains with the human operator unless a single-use override explicitly identifies the pull request.

## Bootstrap state

The initial planning structure is established before shipped product code begins. CI-required checks, repository rulesets, security automation, Codex review automation, and release dry runs are tracked outcomes that are activated only after their implementation exists.
