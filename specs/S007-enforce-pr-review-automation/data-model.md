# Data Model: Pull-Request Policy Automation

## Pull Request Snapshot

| Field | Meaning | Validation |
|---|---|---|
| `number` | Repository-local pull-request number | Positive integer |
| `head_sha` | Revision receiving statuses | Full lowercase hexadecimal commit identifier |
| `base_ref` | Target branch governing policy eligibility | Must equal `main` before any mutation |
| `author_login` | Pull-request author | Canonical login comparison |
| `draft` | Whether native review can start | Draft review policy remains pending |
| `state` | Open or closed lifecycle | Recovery sweeps only open pull requests |
| `closing_issues` | GitHub-resolved closing targets | Fully paginated; non-empty passes ordinary link policy |
| `comments` | Conversation comments | Fully paginated and treated as untrusted data |
| `reactions` | Pull-request reactions | Fully paginated; actor, kind, and time retained |
| `reviews` | Submitted pull-request reviews | Fully paginated; actor, commit, body, and time retained |
| `review_threads` | Current inline review threads | Fully paginated with resolution and comment evidence |
| `historical_heads` | Every head revision observed in pull-request history | Fully paginated, full lowercase commit identifiers used to bind operator review requests |

## Comment

| Field | Meaning | Validation |
|---|---|---|
| `id` | Stable GitHub identifier | Positive integer |
| `author_login` | Comment author | Exact configured identity after narrow bot normalization |
| `body` | Untrusted Markdown | Parsed only for exact markers and directives |
| `created_at` | Publication time | Orders round boundaries |
| `updated_at` | Last edit time | Detects terminal activity after a request |

## Review Thread

| Field | Meaning | Validation |
|---|---|---|
| `id` | GraphQL thread identifier | Non-empty |
| `resolved` | Current GitHub thread state | Unresolved Codex findings block |
| `comments` | Ordered thread comments | Codex authorship and originating review commit retained |
| `first_codex_at` | First Codex finding time | Assigns the finding before or after round two |
| `codex_commit` | Revision reviewed by Codex | Must relate to the current head or remain stale evidence |

## Codex Summary

| Field | Meaning | Validation |
|---|---|---|
| `commit_prefix` | Reviewed commit in the latest summary row | At least seven hexadecimal characters and a unique prefix of `head_sha` |
| `status` | Latest native result | `pending`, `completed`, `failed`, or `unknown` |
| `updated_at` | Latest activity time | Must occur within the active round |

## Operator Exception

| Field | Meaning | Validation |
|---|---|---|
| `policy` | Waived policy | Exactly `issue-link` or `codex-review` |
| `reason` | Audit rationale | Non-empty after trimming |
| `actor` | Authority | Exactly the configured operator login |
| `comment_id` | PR-specific evidence | Stable positive identifier |

## Policy Result

| Field | Meaning | Validation |
|---|---|---|
| `context` | Commit-status context | One of the two governed values |
| `state` | Commit-status state | `pending`, `success`, `failure`, or `error` |
| `description` | Maintainer-facing summary | Non-empty and within GitHub limits |
| `head_sha` | Evaluated revision | Equals the fetched current head |
| `request_second` | Permitted mutation | True only for resolved first-round findings with no prior second request |
| `evidence` | Deterministic diagnostics | No secrets or local paths |

## Review Lifecycle

```text
draft -> round_one_pending -> round_one_clean
                         \-> round_one_findings -> round_one_resolved -> round_two_pending
                                                                     \-> round_two_clean
                                                                     \-> round_two_blocked
                         \-> round_one_failed
```

Dependabot and operator waivers enter explicit success states. A changed head invalidates terminal review evidence. Only `round_one_resolved` can set `request_second`; every state at or after `round_two_pending` fixes it to false.
