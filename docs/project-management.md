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

The pull-request policy consumes GitHub's resolved closing-issue references rather than maintaining a separate Markdown parser. A normal pull request with no resolved closing issue fails `Cueson PR policy / Issue link`. Dependabot passes through an explicit automated-source exception. Any other exception requires an exact `skip: issue-link - <reason>` comment from the configured human operator.

`Cueson PR policy / Codex review` reconciles the native integration's evidence against the current pull-request head. Draft, acknowledged, incomplete, or reaction-only work without current-head terminal evidence remains pending. Unresolved findings, failed or ambiguous evidence, stale terminal results, and protocol violations remain blocking. A clean first round passes without another request.

After a finding-bearing first review, automation may post exactly one marked `@codex review` request only when every resolvable finding is resolved, the reviewed commit is an ancestor of the remediation head, and the configured CI gates pass on that head. Only a marker attributed to the exact GitHub Actions bot or a configured-operator second request consumes the allowance permanently; marker-shaped comments from other actors fail closed. An operator request must cite exactly one backticked commit prefix that resolves uniquely in the pull-request head history before review evidence can be attributed to that request, its latest edit time establishes the boundary, and first observation records an authenticated reservation so deletion cannot restore another request. A second-round failure or finding never causes an automatic third request. Only the configured human operator may waive the link or review policy, and every exception is pull-request-specific with a non-empty reason.

The trusted workflow uses current GitHub evidence for every reconciliation and writes only current-head statuses plus the single permitted review comment. Pull-request and comment events provide normal wakeups. A non-hourly schedule recovers reaction-only completion and review-thread changes, while a configured-operator `/cueson reconcile` comment requests immediate recovery. Unknown evidence and incomplete API pagination fail closed without mutation.

The workflow becomes active only after its implementation reaches the default branch. Its introducing pull request receives read-only live API validation plus fixture-backed mutation proof. Hosted status and Actions-bot request behavior must be demonstrated on the first eligible post-merge pull request before repository rules make either context required.

## Repository controls

The stable S006 CI and CodeQL contexts may become required only after successful current-head runs are read back with their provider identities. The S007 issue-link and Codex-review contexts have an additional activation gate: neither becomes required until one eligible post-S007 pull request proves both trusted status publication and the GitHub Actions-authored second-round request path. A clean review leaves that pair deferred without weakening independently proven gates.

Cueson-specific protection uses a repository-owned ruleset and does not modify the organization-owned baseline. The repository rule targets only `main`, requires pull requests and resolved conversations, blocks deletion and non-fast-forward updates, and contains one organization-administrator recovery bypass. Required checks use strict current-base policy and the exact evidence-backed contexts recorded in [repository controls](repository-controls.md).

Repository configuration is changed in dependency order and each write is followed by a separate authoritative read. External-state evidence records both successful controls and plan or policy limitations. Automatic branch deletion is verified after the next operator merge; a successful setting read-back proves configuration, while the post-merge observation proves behavior.

## Bootstrap state

The initial planning structure is established before shipped product code begins. CI-required checks, repository rulesets, security automation, Codex review automation, and release dry runs are tracked outcomes that are activated only after their implementation exists.
