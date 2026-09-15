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

After a finding-bearing first review, automation may post exactly one marked `@codex review` request only when every resolvable finding is resolved, the reviewed commit is an ancestor of the remediation head, and the configured CI gates pass on that head. Only a marker attributed to the exact GitHub Actions bot or a configured-operator second request consumes the allowance permanently; marker-shaped comments from other actors fail closed. An operator request must cite exactly one backticked commit prefix that resolves uniquely in the pull-request head history before review evidence can be attributed to that request, its latest edit time establishes the boundary, and first observation records an authenticated reservation so deletion cannot restore another request. A second-round finding may be remediated without a third review only when its thread is resolved, its reviewed head is a proven ancestor of the remediation head, and every configured CI gate succeeds on that newer head. Missing ancestry, incomplete CI, unresolved findings, or a second-round failure remains blocking. Only the configured human operator may waive the link or review policy, and every exception is pull-request-specific with a non-empty reason.

The trusted workflow uses current GitHub evidence for every reconciliation and writes only current-head statuses plus the single permitted review comment. Pull-request and comment events provide normal wakeups. A non-hourly schedule recovers reaction-only completion and review-thread changes, while a configured-operator `/cueson reconcile` comment requests immediate recovery. Unknown evidence and incomplete API pagination fail closed without mutation.

The workflow became active when S007 merged through pull request [#19](https://github.com/shruggietech/cueson/pull/19). S008 and S009 subsequently demonstrated trusted current-head status publication from default-branch code. The S009 automatic second-round path reserved its one request but GitHub rejected the Actions-authored comment with HTTP 403, so that capability remains unproven and the two pull-request-policy contexts remain outside the required-check ruleset. Pull request [#21](https://github.com/shruggietech/cueson/pull/21) completed through the documented operator exception after its one manual second round and all remediation-head gates succeeded.

## Repository controls

The 17 stable S006 CI and CodeQL contexts became required only after successful current-head runs were read back with their provider identities during S008. The S007 issue-link and Codex-review contexts have an additional activation gate: neither becomes required until one eligible pull request proves both trusted status publication and the GitHub Actions-authored second-round request path. S009 proved trusted status publication but exposed the HTTP 403 limitation on Actions-authored review comments, so that pair remains deferred without weakening the 17 independently proven gates.

Cueson-specific protection uses a repository-owned ruleset and does not modify the organization-owned baseline. The repository rule targets only `main`, requires pull requests and resolved conversations, blocks deletion and non-fast-forward updates, and contains one organization-administrator recovery bypass. Required checks use strict current-base policy and the exact evidence-backed contexts recorded in [repository controls](repository-controls.md).

Repository configuration is changed in dependency order and each write is followed by a separate authoritative read. External-state evidence records both successful controls and plan or policy limitations. Automatic branch deletion was behaviorally verified after the operator merged pull requests [#20](https://github.com/shruggietech/cueson/pull/20) and [#21](https://github.com/shruggietech/cueson/pull/21); both remote head references are absent.

## Delivery state

Work slices S001 through S024 are present on `main`. S023 planning/contracts merged through PR [#68](https://github.com/shruggietech/cueson/pull/68), followed by S024 through PR [#69](https://github.com/shruggietech/cueson/pull/69) at `86ea7ecf8cd9166f1603a528f21e4e1c214f9fa0`, closing #55 and #56. S025 now implements #57/#58/#59 under unreleased `1.1.0-dev`. The foundation, published releases, stable SubRip/WebVTT and S021 public hosting remain established; production remains the reviewed S021 revision and does not automatically deploy on merge.

| Outcome | Delivery evidence on `main` |
|---|---|
| Issue [#1](https://github.com/shruggietech/cueson/issues/1), repository bootstrap | Bootstrap commit `1188b5e`, completed by pull request [#13](https://github.com/shruggietech/cueson/pull/13) at `ae796dd` |
| Issue [#3](https://github.com/shruggietech/cueson/issues/3), ratified foundation contracts | Pull request [#13](https://github.com/shruggietech/cueson/pull/13), merge commit `ae796dd` |
| Issue [#4](https://github.com/shruggietech/cueson/issues/4), CLI foundation | Pull request [#14](https://github.com/shruggietech/cueson/pull/14), merge commit `addcecb` |
| Issue [#5](https://github.com/shruggietech/cueson/issues/5), schema foundation | Pull request [#15](https://github.com/shruggietech/cueson/pull/15), merge commit `0b60c52` |
| Issue [#6](https://github.com/shruggietech/cueson/issues/6), source integrity and restoration | Pull request [#16](https://github.com/shruggietech/cueson/pull/16), merge commit `4c8513f` |
| Issue [#7](https://github.com/shruggietech/cueson/issues/7), fixture and conformance foundation | Pull request [#17](https://github.com/shruggietech/cueson/pull/17), merge commit `2903824` |
| Issue [#8](https://github.com/shruggietech/cueson/issues/8), CI and cross-platform gates | Pull request [#18](https://github.com/shruggietech/cueson/pull/18), merge commit `d66d6a0` |
| Issue [#9](https://github.com/shruggietech/cueson/issues/9), pull-request policy automation | Pull request [#19](https://github.com/shruggietech/cueson/pull/19), merge commit `bce3ebe` |
| Issue [#10](https://github.com/shruggietech/cueson/issues/10), repository controls | Pull request [#20](https://github.com/shruggietech/cueson/pull/20), merge commit `cc2bac0` |
| Issue [#11](https://github.com/shruggietech/cueson/issues/11), non-publishing release proof | Pull request [#21](https://github.com/shruggietech/cueson/pull/21), merge commit `3da0a4b` |
| Issue [#12](https://github.com/shruggietech/cueson/issues/12), documentation and milestone verification | Pull request [#22](https://github.com/shruggietech/cueson/pull/22), merge commit `730b649` |
| Issue [#23](https://github.com/shruggietech/cueson/issues/23), official brand-kit integration | Pull request [#24](https://github.com/shruggietech/cueson/pull/24), merge commit `7756309` |
| Issue [#25](https://github.com/shruggietech/cueson/issues/25), protected release preparation | Pull request [#26](https://github.com/shruggietech/cueson/pull/26), merge commit `b294a69` |
| Issue [#27](https://github.com/shruggietech/cueson/issues/27), public v0.0.0 publication and verification | Annotated tag [`v0.0.0`](https://github.com/shruggietech/cueson/tree/v0.0.0), verified [GitHub Release](https://github.com/shruggietech/cueson/releases/tag/v0.0.0), and pull request [#28](https://github.com/shruggietech/cueson/pull/28), merge commit `a558a3a` |
| Issue [#30](https://github.com/shruggietech/cueson/issues/30), codec and model foundation | Pull request [#39](https://github.com/shruggietech/cueson/pull/39), merge commit `691cb29` |
| Issue [#31](https://github.com/shruggietech/cueson/issues/31), native SubRip workflows | Pull request [#39](https://github.com/shruggietech/cueson/pull/39), merge commit `691cb29` |
| Issue [#32](https://github.com/shruggietech/cueson/issues/32), native WebVTT workflows | Pull request [#40](https://github.com/shruggietech/cueson/pull/40), merge commit `f4efe02` |
| Issue [#33](https://github.com/shruggietech/cueson/issues/33), bidirectional conversion | Pull request [#42](https://github.com/shruggietech/cueson/pull/42), merge commit `fb5e136` |
| Issue [#34](https://github.com/shruggietech/cueson/issues/34), validation, inspection, and completion | Pull request [#43](https://github.com/shruggietech/cueson/pull/43), merge commit `e0d9a1c` |
| Issues [#35](https://github.com/shruggietech/cueson/issues/35), [#36](https://github.com/shruggietech/cueson/issues/36), and [#41](https://github.com/shruggietech/cueson/issues/41), v1 contract freeze | Pull request [#44](https://github.com/shruggietech/cueson/pull/44), merge commit `f5f2cb2` |
| Issue [#37](https://github.com/shruggietech/cueson/issues/37), exact v1 release candidate | Pull request [#45](https://github.com/shruggietech/cueson/pull/45), commit `2cad4c8` |
| Issue [#38](https://github.com/shruggietech/cueson/issues/38), official v1 release verification and reconciliation | Pull request [#46](https://github.com/shruggietech/cueson/pull/46), merge commit `8fc80a2` |
| Issue [#47](https://github.com/shruggietech/cueson/issues/47), public site and immutable schema hosting | Pull request [#48](https://github.com/shruggietech/cueson/pull/48), merge/deployment commit `46838fd` |
| Issue [#49](https://github.com/shruggietech/cueson/issues/49), independent DNS address-family verification | Pull request [#50](https://github.com/shruggietech/cueson/pull/50), merge commit `5c75dee` |
| Issues [#52](https://github.com/shruggietech/cueson/issues/52), [#53](https://github.com/shruggietech/cueson/issues/53), and [#54](https://github.com/shruggietech/cueson/issues/54), scripted milestone contracts | Pull request [#68](https://github.com/shruggietech/cueson/pull/68), merge commit `0b1922f` |
| Issues [#55](https://github.com/shruggietech/cueson/issues/55) and [#56](https://github.com/shruggietech/cueson/issues/56), historical compatibility and scripted models | Pull request [#69](https://github.com/shruggietech/cueson/pull/69), merge commit `86ea7ec` |

Both `v0.0.0` and `v1.0.0` milestones are closed after independently verified publication and repository reconciliation. Their epics and all prior issues, including #38, #47 and #49, are closed; the twenty-eight historical Project items use Stage `Done` and leave default Status unused. Completed hosting and milestone closure remain recorded outcomes, not future authorization for another release or production change.

The v1 release binds version/schema lockstep, stable SubRip/WebVTT declarations, immutable schema bytes, concise release notes, six archives, six target-bound SPDX JSON SBOMs and one checksum manifest to exact revision `2cad4c816340404289b4d1d87179a4071713bb46`. S020 used explicit tag/Release authority; the operator subsequently merged its PR and reconciled the v1 epic/milestone. S021 used separate explicit public-hosting/deployment authority. Signatures, attestations, future tag/releases, production changes and each final PR merge remain protected transitions.

## Next milestone and active slice

The [S023 roadmap](roadmap.md) maintains [milestone v1.1.0](https://github.com/shruggietech/cueson/milestone/3), [epic #51](https://github.com/shruggietech/cueson/issues/51) and sixteen atomic children. S023 issues #52/#53/#54 and S024 issues #55/#56 are closed. S025 owns #57/#58/#59, using In progress and Slice S025 at kickoff. Native GitHub stages and dependencies control; no additional active arrivals were present at the 2026-09-15 post-S024 assessment.

At S023 publication read-back, `cueson Delivery` contains forty-five unique issue items: twenty-eight historical Done items and seventeen new milestone items. Every child has the native milestone/parent/blockers and governed labels, owning Slice text, actual Stage and unused default Status. The spanning epic uses In progress and no single Slice. S023 issues proceed through PR review until human merge, while downstream issues stay Backlog until dependencies clear. Native GitHub state controls after subsequent transitions; this dated read-back is not a permanent assertion that the queue never changes.
