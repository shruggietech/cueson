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

S008 stopped all administrative mutations at this failure, amended its Spec Kit requirements, and added test-first remediation before continuing. Final values, ruleset identity, admitted check evidence, limitations, and verification timestamps follow only after the corrected adapter and authoritative read-backs succeed.

## Recovery and failure handling

A mismatched or unavailable read-back stops dependent mutations. Earlier verified controls remain reported as partial progress. Recovery begins by reading current state again and never assumes that a failed request was atomic.

The organization-administrator bypass exists only to recover from a broken policy configuration. S008 verification inspects mergeability without using that bypass, and the agent stops before final merge regardless of its technical ability to bypass.
