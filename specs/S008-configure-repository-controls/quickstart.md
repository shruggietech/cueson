# Quickstart: Verify Repository Controls

## Prerequisites

- Authenticate GitHub CLI as the configured operator with repository administration, workflow, security, and Project permissions.
- Work from a clean `S008-configure-repository-controls` checkout.
- Confirm the official S008 pull request targets `main` and resolves issue #10.
- Keep all verification foreground and non-interactive.

## 1. Verify repository-authored artifacts

```powershell
go run ./scripts/github-format/main.go .
go test -count=1 ./...
go test -race -count=1 ./...
go vet ./...
go -C scripts/github-format test -count=1 ./...
go -C scripts/pr-policy test -count=1 ./...
git diff --check
```

Expected: every command exits successfully, no repository text is repaired, and the worktree remains clean after committed changes.

## 2. Capture the live baseline

Read repository merge settings, Actions policy and workflow defaults, security-and-analysis state, vulnerability-alert state, automated-security-fix state, private-reporting state, CodeQL configuration, repository rulesets, the effective organization ruleset, and classic branch protection.

Expected: the baseline matches the `Before S008` column in `docs/repository-controls.md`. The organization ruleset is recorded but never mutated.

## 3. Publish and prove hosted checks

Push the S008 branch, create the official pull request with `Closes #10`, read the body back, and wait for current-head CI, CodeQL, PR-policy, Codex, and security-review evidence. Inspect both check runs and combined commit statuses rather than relying only on the pull-request summary.

Expected: every selected S006 context succeeds on the current head from GitHub Actions. Both S007 policy status sources are visible. If a finding-bearing first review causes the trusted workflow to publish the single second-round request, validate the Actions-bot marker and reservation; otherwise record both policy contexts as deferred from protection.

## 4. Apply controls in dependency order

1. Restrict Actions to GitHub-owned actions and require full-SHA references, then read the policy back.
2. Set default workflow permission to read and keep workflow approvals disabled, then read the defaults back.
3. Enable dependency security updates, secret scanning, push protection, and private reporting where available, reading each state back separately.
4. Make squash the only enabled merge method, keep auto-merge disabled, and enable automatic merged-head deletion, then read repository settings back.
5. Create the repository-owned default-branch ruleset with only the proven check list, then read the complete ruleset and effective rule set back.

Expected: every mutation matches the contract in `contracts/repository-controls.md`; an unavailable capability is recorded as limited, and dependent writes stop after any unexplained mismatch.

## 5. Re-prove workflows and final state

After the last source update, wait for all current-head checks and permitted reviews. Re-read repository settings, security states, Actions defaults, repository rulesets, effective branch rules, issue #10, milestone, and Project item. Compare the organization ruleset identifier and update timestamp with the baseline.

Expected: required contexts are successful on the final head, the repository-owned ruleset is active, the organization rule is unchanged, issue #10 remains open for the final merge, and its Project stage is `PR review` with Slice `S008` and no default Status value.

## 6. Publish final evidence and stop

Format the final head-specific evidence through `go run ./scripts/github-format/main.go -stdin`, publish from a file, read it back, verify Markdown structure, and stop without merging, tagging, releasing, or changing production configuration.
