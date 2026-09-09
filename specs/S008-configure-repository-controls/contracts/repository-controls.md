# Repository Controls Contract

## Repository settings

| Control | Required final state |
|---|---|
| Default branch | `main` |
| Squash merge | Enabled |
| Merge commits | Disabled |
| Rebase merge | Disabled |
| Automatic merge | Disabled |
| Delete merged head branches | Enabled |
| Default workflow permission | Read |
| Workflow approval of pull requests | Disabled |
| Allowed action sources | GitHub-owned actions only |
| Full-SHA action pinning | Required |

## Security settings

| Control | Required final state |
|---|---|
| Dependency graph and vulnerability alerts | Enabled |
| Dependabot security updates | Enabled when available |
| Versioned CodeQL workflow | Retained and successful |
| GitHub default CodeQL setup | Not configured, to avoid duplicate analysis |
| Secret scanning | Enabled when available |
| Secret push protection | Enabled when available |
| Private vulnerability reporting | Enabled when available |

## Repository ruleset

The repository owns exactly one S008-created active branch ruleset named `cueson verified default branch`. It targets only the default branch, blocks deletion and non-fast-forward updates, requires pull requests, requires resolved review conversations, allows squash only, applies strict required status checks, and grants one always-available bypass category to organization administrators.

The pre-existing organization ruleset `20478126`, `default-branch PR gate`, is read-only S008 input. Its identifier and update timestamp must not change during the slice.

## Required checks

Every context selected by the ruleset must have a successful result on the current S008 head and match the stable S006 contract. The two S007 policy contexts are candidates only after the complete first-post-merge activation proof in [hosted-evidence.md](hosted-evidence.md).

## Failure contract

Every mutation is followed immediately by a separate read request. A mismatched read-back is failure even when the mutation request returned success. Dependent operations stop after failure. Already verified unrelated controls remain reported as partial progress, and rollback or recovery requires current-state inspection rather than assumptions.
