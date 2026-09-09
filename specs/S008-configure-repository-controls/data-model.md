# Data Model: Verified Repository Controls

## Repository Control

A repository-scoped setting whose desired and observed states can be compared exactly.

| Field | Meaning | Validation |
|---|---|---|
| `name` | Stable control identifier | Unique within S008 |
| `scope` | Repository API surface owning the setting | Must identify `shruggietech/cueson` |
| `before` | Authoritative pre-mutation value | Captured before dependent writes |
| `desired` | Approved final value | Concrete, never unknown |
| `mutation` | Attempted operation and result | Omit secrets and token material |
| `after` | Immediate authoritative read-back | Must equal desired or carry a limitation |
| `status` | `verified`, `limited`, `deferred`, or `failed` | `verified` requires matching read-back |

## Hosted Check Evidence

Evidence that a check or commit status is safe to select as required.

| Field | Meaning | Validation |
|---|---|---|
| `pull_request` | Official S008 pull request number | Must target `main` |
| `head_sha` | Exact revision under review | Must equal the current PR head when recorded |
| `context` | Exact check or status name | Case-sensitive and non-empty |
| `kind` | `check_run` or `commit_status` | Must match the GitHub evidence surface |
| `state` | Terminal observed result | Must be successful before selection |
| `provider` | GitHub App, Actions bot, or status creator | Required whenever exposed |
| `completed_at` | Terminal evidence time | Must follow creation of the head revision |
| `required` | Whether the context enters the repository rule | True only after every evidence gate passes |
| `exclusion_reason` | Reason a candidate remains non-required | Required when `required` is false |

## Repository Ruleset

The single repository-owned rules object added by S008.

| Field | Desired value |
|---|---|
| Name | `cueson verified default branch` |
| Source | Repository `shruggietech/cueson` |
| Target | Branch |
| Ref condition | Default branch only |
| Enforcement | Active |
| Bypass | One `OrganizationAdmin` entry in `always` mode |
| Ref rules | Deletion and non-fast-forward protection |
| Pull request | Required, zero approval minimum, all conversations resolved, squash only |
| Status checks | Strict and limited to proven current-head contexts |

## Security Capability

One requested GitHub security feature and its truthful availability state.

| Field | Meaning |
|---|---|
| `name` | Dependency alerts, automated security fixes, versioned CodeQL, secret scanning, push protection, or private reporting |
| `available` | Whether current plan and policy permit the feature |
| `before` | Initial enabled or configured state |
| `desired` | Intended enabled or intentionally versioned state |
| `after` | Read-back state after any mutation |
| `limitation` | Exact API status and explanation when unavailable |

## Control Evidence Record

The chronological mapping of baseline, hosted proof, mutation, and verification events. Source-controlled documentation records stable configuration and limitations. Head-specific hosted evidence is published on issue #10 or the S008 pull request after the final head exists.

### State transitions

```text
planned -> baseline_recorded -> evidence_ready -> applied -> verified
                                      |             |
                                      |             +-> failed
                                      +-> deferred

planned -> baseline_recorded -> limited
```

Dependent mutations stop at `failed`. A `deferred` required-check candidate does not block unrelated controls. A `limited` security capability satisfies truthful reporting only when the authoritative limitation is preserved.
