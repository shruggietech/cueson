# Research: Verified Repository Controls

## Decision: Add a repository-owned ruleset without editing the organization baseline

**Rationale:** The effective repository state already includes organization ruleset `20478126`, `default-branch PR gate`. Editing that ruleset would change an organization-owned object and could affect repositories beyond Cueson. GitHub supports repository-owned branch rulesets with default-branch targeting, pull-request requirements, required status checks, deletion protection, non-fast-forward protection, and bypass actors. S008 will create one active repository-owned ruleset named `cueson verified default branch` and verify that the organization ruleset identifier and `updated_at` value remain unchanged.

**Alternatives considered:** Extending the organization rule was rejected because the operator authorized Cueson work, not organization-wide policy mutation. Classic branch protection was rejected because rulesets are already the effective policy mechanism and the classic endpoint currently reports no protection. Replacing the organization rule was rejected because rules combine cumulatively and the existing baseline remains useful.

## Decision: Keep one organization-administrator recovery bypass

**Rationale:** The configured operator is an organization administrator and the current effective policy already recognizes that recovery authority. A repository-owned ruleset can retain exactly one `OrganizationAdmin` bypass entry. This avoids new bypasses for automations, deploy keys, teams, or contributor roles while preventing an unproven or renamed check from irrecoverably locking the default branch.

**Alternatives considered:** No bypass was rejected because a bad required-check configuration could strand the repository. A general repository-role bypass was rejected because it can silently widen as role assignments change. A user-specific bypass was rejected because the repository ruleset API's durable organization-administrator category matches existing operator authority and recovery practice.

## Decision: Require all 17 stable S006 workflow contexts after current-head proof

**Rationale:** The S006 check contract names 16 CI jobs and `CodeQL / Analyze Go`. PR #19 and the post-merge `main` run proved those names, and the official S008 pull request will prove them again on its current head before the ruleset is created. Each required check will bind to the GitHub Actions integration identifier when the check-run API exposes that common source. Strict required-check policy will require the pull-request head to include the latest default-branch state.

**Alternatives considered:** Requiring only an aggregate workflow result was rejected because the workflows expose independently selectable jobs and no stable aggregate context. Copying names from workflow YAML without hosted read-back was rejected by issue #10 and the S008 specification. Requiring the additional code-scanning summary named only `CodeQL` was rejected because it is not part of the ratified S006 stable-check contract.

## Decision: Make both S007 policy contexts conditional on complete hosted proof

**Rationale:** S007 explicitly requires the first eligible post-merge pull request to demonstrate both hosted status provenance and GitHub Actions-authored second-round request behavior before either policy context becomes required. S008 will inspect `Cueson PR policy / Issue link` and `Cueson PR policy / Codex review` on its head. If the first review is clean, no legitimate second-round comment should exist, so both contexts remain non-required and the missing proof is documented. If findings are remediated and the workflow emits the single allowed marked request, both contexts may be added only after source, body, reservation, and status read-back succeeds.

**Alternatives considered:** Requiring the issue-link context alone after its successful status was rejected because the recorded S007 activation contract conditions either context on the complete proof. Triggering a synthetic third-party finding or spoofed comment was rejected because it would not prove the trusted identity path. Holding all other repository controls open indefinitely was rejected because CI and CodeQL already have independent complete evidence.

## Decision: Reduce workflow and action-supply-chain defaults

**Rationale:** Repository workflow defaults currently grant write permission even though every versioned workflow declares explicit permissions. S008 will change the default to read and keep pull-request approvals disabled. The repository currently allows every public action and does not require SHA pinning, while every existing external action reference is already a full commit SHA and owned by GitHub. S008 will allow GitHub-owned actions only and require full-SHA pinning, then rerun all workflows.

**Alternatives considered:** Leaving write as the default was rejected because it provides unnecessary fallback privilege. Allowing all Marketplace actions was rejected because no current workflow needs that supply-chain breadth. An explicit per-action allowlist was rejected as unnecessary while all current actions are GitHub-owned and full-SHA pinning supplies the stronger immutable-reference requirement.

## Decision: Enable supported security facilities and preserve versioned CodeQL

**Rationale:** Dependency alerts are already enabled, but automated security fixes are disabled. Secret scanning, push protection, and private vulnerability reporting are disabled even though the repository is public and the security policy directs reporters to the private channel. The versioned CodeQL workflow is healthy, so GitHub default setup will remain unconfigured to avoid competing analysis. S008 will enable automated security fixes, secret scanning, push protection, and private reporting, then read back each state. Plan- or policy-limited subfeatures will be recorded rather than forced.

**Alternatives considered:** Enabling default CodeQL setup was rejected because it would duplicate the reviewed versioned workflow and destabilize the required context. Claiming the existing `SECURITY.md` reporting link was sufficient was rejected because the live private-reporting feature is disabled. Enabling experimental non-provider pattern or validity-check options without entitlement evidence was rejected as outside the issue's required baseline.

## Decision: Use direct foreground API operations plus durable evidence

**Rationale:** GitHub settings are external state, and the issue requires immediate API read-back rather than a reusable deployment service. Direct `gh api` operations keep credentials in the authenticated client, avoid checked-in secrets, and allow each dependent mutation to stop on failure. The source-controlled `docs/repository-controls.md` records the desired state, evidence model, and final observed configuration; issue and pull-request comments record final head-specific hosted evidence that cannot be self-referentially embedded in the commit it identifies.

**Alternatives considered:** A new administrative Go program was rejected as disproportionate for one repository and would duplicate GitHub CLI authentication and REST behavior. Background polling was rejected by the repository's Windows process and verification rules. Recording only successful API exit codes was rejected because a successful mutation response is not authoritative read-back.

## Decision: Bind status read-back to the requested endpoint instead of a redundant response field

**Rationale:** The first hosted S008 reconciliation proved that GitHub's commit-status create and list representations omit a `sha` property even though both requests are already scoped to `/commits/{sha}/statuses`. S007's fixture incorrectly supplied that redundant field, so the first accepted status mutation was reported as a read-back mismatch. The adapter will treat the validated request path as the commit binding and will compare the returned status identifier, context, state, description, target URL, and exact GitHub Actions creator. Existing matching statuses are no-ops only when that creator is trusted.

**Alternatives considered:** Accepting any status with matching text was rejected because an unexpected writer could suppress a trusted replacement. Inferring the SHA from unrelated response text was rejected because the request path is the authoritative binding. Retrying the failed mutation without fixing the adapter was rejected because every first mutation would fail once and violate immediate read-back success.

## Sources

- [GitHub REST API endpoints for repository rulesets](https://docs.github.com/en/rest/repos/rules)
- [GitHub rules available in rulesets](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/available-rules-for-rulesets)
- [GitHub REST API endpoints for Actions permissions](https://docs.github.com/en/rest/actions/permissions)
- [GitHub REST API endpoints for repositories and security settings](https://docs.github.com/en/rest/repos/repos)
- [GitHub private vulnerability reporting configuration](https://docs.github.com/en/code-security/how-tos/report-and-fix-vulnerabilities/configure-vulnerability-reporting/configure-for-a-repository)
