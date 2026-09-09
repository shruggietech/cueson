# Repository controls

**Repository:** `shruggietech/cueson`

**Default branch:** `main`

**Control owner:** Work slice `S008`, issue [#10](https://github.com/shruggietech/cueson/issues/10)

This document records the repository's delivery-security baseline, desired state, mutation order, and authoritative GitHub read-back. It does not grant merge, tag, release, schema-publication, production-domain, organization-policy, or unrelated-repository authority.

## Control principles

- Repository settings are external state. A successful mutation response is provisional until a separate read returns the intended value.
- Required checks are selected only from successful hosted results on the exact pull-request head under review.
- Repository rules add to the existing organization baseline. S008 never edits the organization-owned ruleset.
- A single organization-administrator bypass category preserves operator recovery from a broken gate. It is not the normal merge path.
- Unsupported or plan-limited security features remain explicit limitations, never implied successes.

## Before S008

The baseline was read through GitHub's repository, Actions, code-security, ruleset, and branch-protection APIs on 2026-09-09 before any S008 mutation.

### Repository and merge settings

| Control | Baseline |
|---|---|
| Default branch | `main` |
| Squash merge | Enabled |
| Merge commits | Enabled |
| Rebase merge | Enabled |
| Auto-merge | Disabled |
| Delete merged head branches | Disabled |

### Actions settings

| Control | Baseline |
|---|---|
| Actions | Enabled |
| Allowed action sources | All actions |
| Full-SHA pinning required | No |
| Default workflow permission | Write |
| Workflows may approve pull requests | No |

### Security settings

| Control | Baseline |
|---|---|
| Dependency graph and vulnerability alerts | Enabled (`204 No Content` from the enablement probe) |
| Dependabot security updates | Disabled |
| Versioned CodeQL workflow | Enabled and successful |
| GitHub default CodeQL setup | Not configured |
| Open code-scanning alerts | Zero |
| Secret scanning | Disabled |
| Secret push protection | Disabled |
| Private vulnerability reporting | Disabled |

### Effective branch policy

GitHub's classic branch-protection endpoint returned `404 Branch not protected`, but that endpoint is not the effective-policy authority because an organization ruleset applies to the repository.

| Property | Baseline |
|---|---|
| Ruleset ID | `20478126` |
| Name | `default-branch PR gate` |
| Source | Organization `shruggietech` |
| Enforcement | Active |
| Target | Default branch |
| Rules | Block deletion and non-fast-forward updates; require pull requests |
| Conversation resolution | Not required |
| Required status checks | None |
| Bypasses | Organization administrators and repository role `5` |
| Updated | `2026-08-11T18:09:16.051-04:00` |

S008 treats the ruleset identifier, source, rules, bypass actors, and update timestamp as immutable comparison evidence.

## Desired state

### Repository and Actions

| Control | Desired state |
|---|---|
| Merge strategy | Squash only |
| Auto-merge | Disabled |
| Delete merged head branches | Enabled |
| Allowed action sources | GitHub-owned actions only |
| Full-SHA pinning | Required |
| Default workflow permission | Read |
| Workflow pull-request approval | Disabled |

### Security

| Control | Desired state |
|---|---|
| Dependency alerts | Enabled |
| Dependabot security updates | Enabled when available |
| Versioned CodeQL workflow | Retained and successful |
| Default CodeQL setup | Not configured |
| Secret scanning | Enabled when available |
| Secret push protection | Enabled when available |
| Private vulnerability reporting | Enabled when available |

### Repository-owned rules

S008 will add one active branch ruleset named `cueson verified default branch`. It targets only the default branch, blocks deletion and non-fast-forward updates, requires pull requests and resolved review conversations, allows only squash merging, applies strict required status checks, and contains one `OrganizationAdmin` recovery bypass in `always` mode.

The required-check candidates are the 17 stable S006 contexts recorded in `specs/S006-establish-ci-gates/contracts/check-contract.md`. Each candidate must succeed on the current S008 head from the expected GitHub Actions integration before admission.

The `Cueson PR policy / Issue link` and `Cueson PR policy / Codex review` statuses are conditional candidates. S007 requires both hosted status-source evidence and GitHub Actions-authored second-round request evidence before either becomes required. A clean first review correctly produces no second request and therefore leaves both contexts deferred.

## Mutation order

1. Publish the official S008 pull request with `Closes #10` and read its Markdown body back.
2. Observe initial current-head CI, CodeQL, and PR-policy results without changing repository settings.
3. Restrict action sources and require immutable action references, then read the Actions policy back.
4. Reduce default workflow permission to read and keep workflow approvals disabled, then read the defaults back.
5. Enable supported security facilities one at a time, reading each capability back before continuing.
6. Make squash the only merge method and enable automatic merged-head deletion, then read repository settings back.
7. Push the evidence update and require workflows to succeed under the tightened defaults.
8. Build the required-check list from that exact head and create the repository-owned ruleset.
9. Read the complete ruleset and effective rules back, compare the organization rule with its baseline, and publish final head-specific evidence.

## After S008

### Hosted activation observation

The official pull request is [#20](https://github.com/shruggietech/cueson/pull/20). Its initial head was `83d4f58bde8d9df5c1871fd1a5ce69498acbfc7f`, its formatted 39-line body read back successfully, and GitHub resolved `Closes #10`.

The first trusted reconciliation run, [34400932303](https://github.com/shruggietech/cueson/actions/runs/34400932303), published `Cueson PR policy / Issue link` successfully but then failed its local read-back comparison. GitHub's actual create and list representations omit the redundant `sha` field while binding the request to the commit-status endpoint for the exact SHA. A comment-triggered retry published `Cueson PR policy / Codex review` and failed for the same reason. This is a dependency defect, not a failed policy decision: the two statuses exist on the intended head with IDs `53854883221` and `53854918611`, exact context, expected state and description, run target, and creator `github-actions[bot]` ID `41898282`.

S008 stopped all administrative mutations at this failure, amended its Spec Kit requirements, and added test-first remediation before continuing. The adapter now validates the requested commit-status endpoint, exact status fields, and GitHub Actions creator without requiring a redundant `sha` member that GitHub does not return. Regression coverage also proves that an otherwise matching status from an untrusted creator is replaced rather than accepted as a no-op.

### Verified repository and Actions settings

The first control group completed at 2026-09-09T20:34:36Z. Every mutation received a separate successful read-back.

| Control | Before | Verified after |
|---|---|---|
| Actions | Enabled | Enabled |
| Allowed action sources | All actions | GitHub-owned actions only |
| Full-SHA pinning required | No | Yes |
| Default workflow permission | Write | Read |
| Workflows may approve pull requests | No | No |
| Squash merge | Enabled | Enabled |
| Merge commits | Enabled | Disabled |
| Rebase merge | Enabled | Disabled |
| Auto-merge | Disabled | Disabled |
| Delete merged head branches | Disabled | Enabled |

The selected-action read-back returned `github_owned_allowed: true`, `verified_allowed: false`, and an empty additional-pattern list. Every existing workflow action is GitHub-owned and already uses a full commit SHA.

### Verified security settings

| Control | Before | Verified after |
|---|---|---|
| Dependency graph and vulnerability alerts | Enabled | Enabled |
| Dependabot security updates | Disabled | Enabled and not paused |
| Versioned CodeQL workflow | Successful | Successful on S008 head `9cd5666a8bebdfb1bd86e5669cb22cd91ee1644b` in run [34401628498](https://github.com/shruggietech/cueson/actions/runs/34401628498) |
| GitHub default CodeQL setup | Not configured | Not configured |
| Open code-scanning alerts | Zero | Zero |
| Secret scanning | Disabled | Enabled |
| Secret push protection | Disabled | Enabled |
| Open secret-scanning alerts | Endpoint unavailable while disabled | Zero after enablement |
| Private vulnerability reporting | Disabled | Enabled |

GitHub also reports optional non-provider patterns, AI detection, validity checks, delegated alert dismissal, and delegated bypass as disabled. S008 did not enable those separately licensed or experimental subfeatures because they are outside the approved baseline and no implementation depends on them.

### Hosted check evidence and admission

The reduced Actions defaults were exercised by S008 head `6568fc02d0dcff585a73b46ccc06419abea70b37`. CI run [34402009634](https://github.com/shruggietech/cueson/actions/runs/34402009634) and CodeQL run [34402009797](https://github.com/shruggietech/cueson/actions/runs/34402009797) completed successfully. Every admitted check run was created by the GitHub Actions integration, app ID `15368`. This proves the S006 CI and CodeQL contract under the reduced defaults; it does not claim that the separate trusted-base pull-request-policy adapter completed successfully before its S008 fix reached `main`.

GitHub's required-status API uses the raw check-run `name`, while the S006 contract also records the workflow-qualified display name. The following 17 raw names are therefore the exact required contexts:

- `Analyze Go`
- `Formatting`
- `Repository text`
- `Vet`
- `Schema and conformance`
- `Static analysis`
- `Vulnerability scan`
- `Native tests (Linux)`
- `Native tests (Windows)`
- `Native tests (macOS)`
- `Race detection`
- `Pure-Go build (linux-amd64)`
- `Pure-Go build (linux-arm64)`
- `Pure-Go build (windows-amd64)`
- `Pure-Go build (windows-arm64)`
- `Pure-Go build (darwin-amd64)`
- `Pure-Go build (darwin-arm64)`

The security summary check named `CodeQL`, created by GitHub Code Scanning app ID `57789`, also succeeded. It is not part of the stable S006 contract and is not required by the repository ruleset.

On the same head, combined-status read-back contained `Cueson PR policy / Issue link` with state `success`, description `Closes shruggietech/cueson#10`, target run [34402007765](https://github.com/shruggietech/cueson/actions/runs/34402007765), and creator `github-actions[bot]` ID `41898282`. The `Cueson PR policy / Codex review` status was absent. The policy workflow's `Reconcile` check failed because `pull_request_target` correctly ran the still-trusted adapter from `main`, which cannot include S008's fix before merge. Running the pull-request copy with administrative credentials would violate the trust boundary, so S008 records this as a pre-merge bootstrap limitation instead of treating policy completion as a successful hosted gate. In addition, the only clean first-round Codex review was bound to the earlier head `83d4f58bde8d9df5c1871fd1a5ce69498acbfc7f`, so the S007 Actions-bot second-round proof gate was not satisfied on this admission head. Both policy contexts remain deferred together. Neither is a required check. Post-merge housekeeping must exercise the corrected adapter from trusted `main` before either context is reconsidered.

### Verified repository rules

At 2026-09-09T20:39:08Z, GitHub created repository ruleset `22685253`, `cueson verified default branch`. Its complete read-back at 2026-09-09T20:41:47Z showed:

| Property | Verified value |
|---|---|
| Source | Repository `shruggietech/cueson` |
| Enforcement | Active |
| Target | Branches matching `~DEFAULT_BRANCH` |
| Ref protection | Deletion and non-fast-forward updates blocked |
| Pull requests | Required |
| Conversation resolution | Required |
| Allowed merge methods | Squash only |
| Required approvals | Zero |
| Required checks | The 17 contexts above, each bound to GitHub Actions app ID `15368` |
| Strict checks | Enabled |
| Enforcement on branch creation | Disabled |
| Bypass | One `OrganizationAdmin` actor in `always` mode |

No bot, deploy key, repository role, team, or second bypass entry exists in the repository ruleset. Effective-rule read-back for `main` showed the organization and repository rules composing together, including the repository's stricter conversation-resolution, squash-only, and required-check controls.

Organization ruleset `20478126` retained its original organization source, active enforcement, deletion and non-fast-forward rules, pull-request rule, two bypass actors, and update timestamp `2026-08-11T18:09:16.051-04:00`. S008 did not mutate it.

### Boundary and mergeability audit

Read-only inspection reported pull request #20 as mergeable without using the organization-administrator bypass. Its merge state remained unstable only because the non-required trusted-main `Reconcile` check exposed the adapter defect described above; all 17 admitted checks and the CodeQL security summary were successful on the inspected head.

S008 did not merge or enable auto-merge, create or move a tag, publish a release or schema, alter production-domain configuration, modify an organization policy, or touch another repository. Automatic merged-head deletion is configured but cannot be behaviorally observed until the operator performs the separately authorized final merge.

### Delivery record

Issue #10 remains open and assigned to milestone `v0.0.0`. Its five acceptance checkboxes are complete, and the official pull request retains the closing reference that delegates issue closure to the final merge. Immediate Project read-back showed the issue exactly once in `cueson Delivery`, Stage `PR review`, Slice `S008`, and no value in the unused default Status field.

The single authorized second Codex review reported one P2 finding: the original T025 and SC-007 wording claimed the trusted-base policy workflow had completed even though its adapter fix cannot run from `main` until merge. Commit `86085eb` narrowed the hosted success claim to the 17 S006 CI and CodeQL contexts, preserved the explicit policy limitation, received a commit-specific reply, and resolved the review thread. No third review was requested. The final head-specific check and review inventory is published on pull request #20 because a source commit cannot name its own hash.

That remediation exposed a second policy defect: the evaluator rejected every head newer than the second-review request before considering whether the second-round finding itself had been resolved. This made the required remediation commit permanently red while a third review remained forbidden. S008 adds test-first handling for the only safe success path: the request head must be a proven ancestor, every second-round thread must be resolved, and every S006 check must succeed on the current remediation head. Missing ancestry, incomplete CI, a clean review followed by unrelated changes, unresolved findings, and failed reviews continue to fail closed. The trusted workflow cannot exercise this correction until it reaches `main`, so clearing the existing red commit status before merge remains an operator recovery decision rather than a reason to run pull-request code with write credentials.

## Recovery and failure handling

A mismatched or unavailable read-back stops dependent mutations. Earlier verified controls remain reported as partial progress. Recovery begins by reading current state again and never assumes that a failed request was atomic.

The organization-administrator bypass exists only to recover from a broken policy configuration. S008 verification inspects mergeability without using that bypass, and the agent stops before final merge regardless of its technical ability to bypass.
