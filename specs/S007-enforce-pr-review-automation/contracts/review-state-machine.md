# Codex Review State-Machine Contract

## Evidence authority

- Codex REST comment/review identity: login `chatgpt-codex-connector[bot]`, numeric user ID `199175422`, type `Bot`
- Codex REST reaction identity: the same exact login and numeric ID can be returned as type `User`; only this transport-specific tuple is normalized to the configured bot identity
- Codex GraphQL identity: exact login `chatgpt-codex-connector` after removing only the known REST `[bot]` suffix
- Summary marker: `<!-- codex-pull-request-review-summary -->`
- Automated second-round marker: `<!-- cueson-codex-review-request:v1 round=2 pr=<number> head=<full-sha> -->`, accepted only from the exact GitHub Actions bot identity
- Operator identity: `h8rt3rmin8r`
- Current revision: the pull request's full head SHA

## State rules

| Condition | Status | Mutation |
|---|---|---|
| Draft eligible pull request | Pending | None |
| Dependabot author | Success, excluded | None |
| Operator `codex-review` waiver | Success, waived | None |
| No current-head terminal Codex evidence | Pending | None; native integration owns round one |
| Eyes without current-head terminal evidence | Pending | None |
| Current-head failed or unknown result | Failure | None |
| Completed current-head result, current-round thumbs-up, no Codex finding | Success, clean | None |
| First-round Codex finding remains unresolved | Failure | None |
| First-round findings all resolved on ancestor A, current descendant B has green required CI, no second request or reservation | Pending | Reserve on B and post the single marked request, then verify both |
| Second reservation exists but request result is ambiguous | Failure | None; operator recovery required |
| Second request exists without later current-head terminal evidence | Pending | None |
| Second-round Codex finding remains unresolved | Failure | None |
| Second-round completion has no unresolved Codex finding | Success | None |
| Second-round failure, stale result, or later head change | Failure | None; operator recovery required |

## Round attribution

The summary is edited in place and is only latest terminal evidence. It is not a historical ledger. Reviews and threads retain their reviewed commit; the second-round marker and status-history reservation establish the round-two boundary. A configured-operator request must cite exactly one backticked commit prefix that uniquely resolves within the collected pull-request head history, and the later of its creation or edit time establishes the boundary; an unbound request consumes the allowance but cannot be rebound to a later head. `isResolved` controls finding resolution, while outdated findings remain blocking until resolved.

The first-round remediation bridge requires GitHub comparison evidence that the finding review commit is an ancestor of and older than the current head, plus successful current-head results for every configured S006 CI gate. This is the only transition that may use finding evidence from a prior head. Divergent or force-pushed history, missing gates, or top-level findings without resolvable threads remain blocking.

A top-level Codex review that signals findings but has no resolvable thread cannot be automatically cleared. A current-head summary and reaction never override an unresolved Codex thread.

## Round ceiling

Any authenticated S007 marker, authenticated round-two reservation status, or configured-operator comment containing a standalone `@codex review` invocation consumes round two. A marker-shaped comment from another actor fails closed without becoming authoritative. Automation never emits another invocation after any authoritative evidence exists.

Before posting, the reconciler refetches the head, all comments, and status history; then it writes and reads back a pending reservation in the Codex policy context. The marked request comment atomically publishes the invocation and durable visible marker. After posting, the reconciler refetches comments and requires exactly one matching marker and invocation before replacing the status with a verified pending result. When a configured-operator request is observed without a reservation, the first reconciliation publishes and reads back that same authenticated reservation before later terminal evidence can complete the round; deletion of the comment therefore leaves a consumed, fail-closed state.

One global non-canceling workflow concurrency group serializes every ordinary event and scheduled sweep. If any mutation response is ambiguous, the reservation remains consumed and the policy fails closed without retry.

## Failure behavior

API errors, incomplete pagination, GraphQL partial errors, malformed summaries, identity mismatches, missing current-head attribution, duplicate markers, contradictory outcomes, and protocol evidence beyond round two fail closed. A failure does not authorize an automatic retry.
