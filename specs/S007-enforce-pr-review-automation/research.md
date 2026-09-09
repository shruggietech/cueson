# Research: Pull-Request Policy Automation

## Decision: Reconcile native round one instead of duplicating it

The configured ChatGPT Codex Connector is the first-round trigger authority. Its pull-request summary states that opening a reviewable pull request or marking a draft ready starts review, and prior Cueson pull requests contain automatic Codex reviews before any maintainer `@codex review` comment. S007 records a pending current-head review status and waits for native evidence rather than posting another first request.

**Rationale:** This satisfies the automatic-first-review rule while preventing duplicate requests and preserving the one remaining automated request for post-remediation round two.

**Alternatives considered:** Always comment on open or ready was rejected because it duplicates the native integration. Treating native activity as outside policy was rejected because round attribution would become unauditable.

## Decision: Run write-capable reconciliation only from trusted default-branch events

The workflow uses `pull_request_target`, `issue_comment`, and `schedule`. These load workflow definitions from the default branch. It explicitly checks out the default branch and never checks out a pull-request head or merge ref. GitHub's review, review-comment, and arbitrary-ref manual workflow events were rejected because they cannot preserve the same trusted-code guarantee.

**Rationale:** The reconciler must comment once and publish per-head statuses, but no pull-request-controlled workflow may receive the write-capable token.

**Alternatives considered:** A write-capable `pull_request` workflow was rejected because the proposed change can alter its mutation path. A secret-backed personal token was rejected because `GITHUB_TOKEN` supplies the narrow repository permissions. An embedded script was rejected in favor of one tested Go decision and transport boundary.

**Primary references:** [Events that trigger workflows](https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows), [GITHUB_TOKEN security and recursion](https://docs.github.com/en/actions/concepts/security/github_token), and [commit statuses](https://docs.github.com/en/rest/commits/statuses).

## Decision: Publish commit statuses for the current pull-request head

The reconciler publishes `Cueson PR policy / Issue link` and `Cueson PR policy / Codex review` against the current head SHA. Each status has a concise description and workflow URL. The command reads the latest status for the same SHA and context after every write and rejects a mismatch.

**Rationale:** Comment-triggered workflows run against the default-branch SHA, so their ordinary job check cannot represent the pull-request head. Commit statuses bind policy to the revision issue #10 may later require.

**Alternatives considered:** A combined context was rejected because linkage and external review have independent recovery paths. Check runs were rejected because their additional surface does not improve S007 acceptance.

## Decision: Make visible GitHub evidence the idempotency store

The marked round-two comment, operator exceptions, Codex summary, review commit identifiers, reactions, reviews, and review threads form the state snapshot. Every collection is fully paginated. One global non-canceling workflow concurrency group serializes all runs. The second-round comment is both the request and durable reservation; before posting the command refetches comments, and afterward it requires exactly one marker.

**Rationale:** Workflow artifacts and caches expire and race. A globally serialized writer plus an atomic visible marker remains auditable. If a post is accepted but its response is lost, the next snapshot sees the marker and cannot post again.

**Alternatives considered:** External storage and labels add authority, retention, and permission problems. Per-event payload state is incomplete. Per-pull-request concurrency was rejected because a scheduled all-open sweep could overlap a numbered event group.

## Decision: Parse only narrow Codex evidence and normalize the known bot identity

REST represents the app as `chatgpt-codex-connector[bot]`, while GraphQL can expose `chatgpt-codex-connector`. The REST reactions endpoint empirically reports that exact login and immutable numeric ID with type `User`, even though comment and review payloads report type `Bot`; the adapter narrowly normalizes only that exact reaction tuple back to the configured bot identity. The policy canonicalizes only the exact optional `[bot]` suffix and rejects every other actor. A summary is recognized only by `<!-- codex-pull-request-review-summary -->`; its latest row must name a recognized `Running`, `In progress`, `Completed`, or `Failed` state and current-head commit. Reactions alone cannot prove commit attribution. Review threads use their Codex review commit, and `isResolved`, not `isOutdated`, determines resolution.

**Rationale:** The Codex summary is edited in place and retains only latest activity, so it is terminal evidence rather than a round ledger. The marked round-two request is the round boundary.

**Alternatives considered:** Trusting any bot, any thumbs-up, or free-form text was rejected as spoofable or stale. Treating absence of findings as success was rejected because acknowledgement and infrastructure failure are not completed reviews.

## Decision: Use GitHub's parsed closing references

The GraphQL pull-request `closingIssuesReferences` connection is the closing-reference authority. The policy accepts a non-empty fully paginated connection and does not implement a competing Markdown parser.

**Rationale:** GitHub owns the supported keyword grammar, ignored Markdown regions, issue existence, and cross-repository resolution. A repository parser would inevitably drift from merge behavior.

**Alternatives considered:** Regular-expression validation was rejected because it can accept fenced examples, fail grouped references, and disagree with GitHub.

## Decision: Restrict exceptions to the configured human operator

An exception is an exact line `skip: issue-link - <reason>` or `skip: codex-review - <reason>` authored by `h8rt3rmin8r`. Association fields, PR authorship, labels, and body text never grant authority. Dependabot is separately excluded from both the link gate and Codex review by canonical actor identity.

**Rationale:** The binding contract gives review waiver authority to the human operator, not every collaborator. The same exact authority is safer for the documented link exception.

**Alternatives considered:** Association-based trust permits collaborator self-waiver. Labels hide rationale and require more permissions. Body directives are author-controlled.

## Decision: Treat a second request as the permanent automated ceiling

Any exact S007 round-two marker attributed to the GitHub Actions bot, or configured-operator `@codex review` comment, consumes round two. A marker-shaped comment from any other actor fails closed without establishing a round boundary. Automation emits its request only after first-round Codex findings exist and all resolvable threads are resolved. After a second request, current-head unresolved findings, failure, stale completion, or later head changes remain blocking until explicit operator action. No later state emits a request.

**Rationale:** Conservative counting is required to prove duplicate, retried, edited, and concurrent events cannot create a third review.

**Alternatives considered:** Counting only workflow-authored comments was rejected because a manual second request shares the same cap. Retrying failures was rejected because it can exceed the bound.

## Decision: Recover missed reaction and thread-resolution events by schedule

GitHub Actions has no dependable first-class trigger for PR-body reactions or review-thread resolution. A low-frequency schedule evaluates up to 100 open pull requests, and a configured-operator `/cueson reconcile` comment provides bounded manual recovery through the trusted `issue_comment` path. Event inputs are wakeups only; every run fetches current state.

**Rationale:** Without recovery, a clean reaction-only review or final resolved thread can leave a stale status indefinitely.

**Alternatives considered:** Polling inside a long-running workflow wastes runner time. Treating pending as success hides missing evidence.

## Decision: Audit S007 live without granting unmerged code write authority

The new default-branch workflow cannot govern the PR that introduces it. After publication, the local S007 checkout runs the command in dry-run mode with a read-scoped token and compares the result with the live S007 pull request. In-memory HTTP tests prove exact status and comment mutations. The first post-merge pull request must prove the hosted status source and whether Codex honors an Actions-bot request before issue #10 makes the contexts required.

**Rationale:** A live read-only run satisfies the controlled inspection available before merge. Granting the branch a maintainer token would contradict the central security requirement. Hosted activation proof is therefore an explicit downstream repository-control prerequisite, not a reason to weaken S007.

**Alternatives considered:** Calling a branch-only workflow was rejected because privileged manual and comment workflows must exist on the default branch. Running the local branch with a write token was rejected because it executes untrusted-to-main code with mutation authority.

## Decision: Bridge first-round findings across a descendant remediation head

Round-one findings retain their reviewed commit A. Automation may request round two on current head B only when GitHub proves A is an ancestor of B, B is newer, every A finding thread is resolved, no unresolvable top-level finding remains, and the configured S006 CI gates pass on B. A force push, divergent comparison, missing check, or ambiguous result fails closed.

**Rationale:** Normal remediation necessarily adds a fix commit after review. Treating all earlier evidence as stale makes round two unreachable, while accepting unrelated stale evidence without ancestry and green-current-head checks is unsafe.

**Alternatives considered:** Requiring the current head to equal the reviewed finding head was rejected because code fixes change the head. Accepting any later head was rejected because force-pushed or unrelated changes could inherit resolved findings.
