# Quickstart: Verify v0.0.0 Documentation and Milestone Readiness

## 1. Validate repository text and documentation links

```powershell
go run ./scripts/github-format/main.go .
go -C scripts/github-format test -count=1 ./...
go -C scripts/docs-verify test -count=1 ./...
go -C scripts/docs-verify run . -repo ../..
git diff --check
```

Expected: all required documents exist, repository text is valid UTF-8 with governed line endings, and every local documentation file and heading link resolves without network access.

## 2. Verify documented product behavior

```powershell
go run ./cmd/cueson --help
go run ./cmd/cueson version
go run ./cmd/cueson schema --version
go test -count=1 ./...
go test -race -count=1 ./...
go vet ./...
```

Expected: help lists only implemented commands, software and schema report `0.0.0`, and existing product tests remain green without S010 changing shipped behavior.

## 3. Verify standalone repository tooling and workflows

```powershell
go -C scripts/pr-policy test -count=1 ./...
go -C scripts/release-verify test -count=1 ./...
go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.12
```

Expected: delivery policy, release verification, and workflow syntax remain green with the documentation verifier included in the Repository text job.

## 4. Audit documentation claims

Review the canonical set in [documentation-contract.md](contracts/documentation-contract.md). Confirm that SRT and WebVTT are described as `envelope_only`, current CLI examples execute, release verification remains non-publishing, protected release actions remain separate, and the changelog does not imply that v0.0.0 has been released.

## 5. Audit GitHub delivery evidence

Read issues #1 through #12, epic #2's child summary, milestone `v0.0.0`, and every issue's `cueson Delivery` Project fields. During review, issues #12 and #2 remain open with Stage `PR review`, Slice `S010`, default Status empty, and complete closing references on the official pull request.

## 6. Validate hosted behavior

After the authorized branch push, wait for every current-head CI, CodeQL, pull-request policy, Codex, and security result. Address every finding, request at most one second Codex round if round one has findings, and stop for the operator after all checks and threads are satisfied.

Expected: the pull request is clean and mergeable, and S010 has not merged, tagged, released, copied an immutable schema, closed the milestone, or changed production state.
