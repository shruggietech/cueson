# Quickstart: Pull-Request Policy Automation

Run commands from the repository root unless they use `go -C`.

## Prerequisites

- Go 1.25 or later
- A read-scoped `CUESON_GITHUB_TOKEN` environment value for the live probe
- No unrelated working-tree changes

## Local policy verification

```text
gofmt -w scripts/pr-policy
go -C scripts/pr-policy test -count=1 ./...
go -C scripts/pr-policy vet ./...
go run ./scripts/github-format/main.go .
go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.12
```

The tests cover the issue-link corpus, every review state, tenfold replay, stale heads, API pagination and partial errors, mutation reservation and read-back, Dependabot, and operator versus non-operator exceptions.

## Controlled read-only pull-request probe

After publishing S007, place a read-scoped token in `CUESON_GITHUB_TOKEN` without passing it on the command line or printing it, then audit without mutation:

```text
go -C scripts/pr-policy run . -repo shruggietech/cueson -pr <PR_NUMBER> -dry-run
```

The command refuses non-dry-run execution outside GitHub Actions. Compare its JSON with the pull-request closing references, Codex summary, reactions, reviews, threads, ancestry, CI, and commit statuses visible through GitHub. Clear the environment variable after the audit.

## Workflow safety inspection

Confirm `.github/workflows/pr-policy.yml`:

- runs only trusted default-branch code for write-capable events;
- never subscribes write authority to merge-ref review events;
- never checks out a pull-request head or merge ref;
- uses only `checks: read`, `contents: read`, `issues: write`, `pull-requests: read`, and `statuses: write`;
- disables checkout credential persistence and dependency caching;
- passes untrusted values as structured data or environment values;
- serializes all runs without canceling an in-flight mutation;
- exposes scheduled and operator-comment recovery without a secret.

## Whole-repository verification

```text
go vet ./...
go test -count=1 ./...
go test -race -count=1 ./...
```

Stage the intended final tree, then run:

```text
git diff --cached --check
git diff --check
```

The new trusted workflow becomes active only after merge. S007's own pull request provides read-only live API-shape evidence; fixture-backed HTTP tests prove status, reservation, second-round comment, and read-back behavior. The first post-merge pull request must prove hosted mutation identity before issue #10 makes the contexts required.
