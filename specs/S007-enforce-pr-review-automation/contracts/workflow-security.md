# Workflow Security and Event Contract

## Events

The trusted reconciler observes:

- `pull_request_target` for opened, reopened, edited, synchronize, converted-to-draft, ready-for-review, and closed transitions targeting `main`;
- `issue_comment` for created, edited, and deleted pull-request conversation comments;
- a low-frequency non-hourly schedule for reaction-only completion, thread-resolution, and missed-event recovery.

A configured-operator `/cueson reconcile` pull-request comment provides bounded manual recovery through the same default-branch `issue_comment` path. The workflow does not accept `workflow_dispatch` refs and does not subscribe a write-capable job to pull-request review or review-comment merge-ref events.

Every event resolves zero, one, or all open pull-request numbers and then fetches a complete current snapshot. Event payload content is a routing hint only. Because `issue_comment` has no branch filter, the adapter verifies the fetched base ref is `main` before any mutation.

## Trusted code boundary

The workflow explicitly checks out `main` with credential persistence disabled. It never checks out a pull-request head, merge ref, contributor branch, generated patch, artifact, or executable cache.

The workflow invokes only default-branch `scripts/pr-policy`. Pull-request bodies, comments, reactions, reviews, paths, and identifiers are untrusted structured data and are never interpolated into shell programs.

## Permissions

```yaml
permissions:
  checks: read
  contents: read
  issues: write
  pull-requests: read
  statuses: write
```

The checks scope is read-only and supplies current-head CI conclusions for the remediation bridge. No secret is consumed. `GITHUB_TOKEN` is passed only through the environment.

## Mutation contract

The command may create only latest current-head statuses for `Cueson PR policy / Issue link` and `Cueson PR policy / Codex review`, plus the one marked second-round comment when its state-machine preconditions hold.

It cannot merge, approve, dismiss, label, close, reopen, edit settings, dispatch another workflow, upload an artifact, push content, or mutate releases.

Every changed status is read back and matched by SHA, context, state, and description. Unchanged scheduled results do not create another status. Every second-round comment is preceded by a read-back reservation and followed by a complete comment refetch proving exactly one authentic marker and invocation.

## Recovery and concurrency

Scheduled recovery lists at most 100 open pull requests and evaluates each independently. One global `cueson-pr-policy` concurrency group uses `cancel-in-progress: false`, so a sweep cannot overlap a numbered event mutation. A failure in one pull request is reported and never relaxes another result.
