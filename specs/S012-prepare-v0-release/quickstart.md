# Quickstart: Verify the v0.0.0 Release Preparation

## 1. Verify repository release records

```powershell
go run ./scripts/github-format/main.go .
go -C scripts/github-format test -count=1 ./...
go -C scripts/docs-verify test -count=1 ./...
go -C scripts/docs-verify run . -repo ../..
git diff --check
```

Expected: repository text is UTF-8 without BOM, governed line endings and Markdown layout pass, maintained links resolve, `CHANGELOG.md` contains a fresh `[Unreleased]` section plus dated v0.0.0 history, and `docs/releases/v0.0.0.md` ends with the exact tagged changelog link.

## 2. Verify schema admission and release policy

```powershell
go -C scripts/release-verify test -count=1 ./...
go -C scripts/release-verify test -race -count=1 ./...
go -C scripts/release-verify vet ./...
go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.12
```

Expected: canonical and versioned schemas match byte for byte, release notes and evidence contracts pass, the workflow includes pull-request, manual, and `main` push proof while retaining read-only non-publishing authority, and the GoReleaser archive contract still has exactly four members.

## 3. Build and verify the exact candidate

Install the exact GoReleaser and Syft versions from `docs/release-verification.md`, then work from a clean commit:

```powershell
goreleaser check
goreleaser release --snapshot --clean --skip=publish
$sourceRevision = git rev-parse HEAD
go -C scripts/release-verify run . -dist ../../dist -repo ../.. -version 0.0.0 -commit $sourceRevision -execute-host
```

Expected: six archives, six SBOMs, six checksums, and four members per archive pass; every binary and schema reports `0.0.0`; `dist/release-evidence.json` records the same full commit, the versioned schema SHA-256, complete target digests, and `published: false`.

## 4. Verify existing product and repository behavior

```powershell
go test -count=1 ./...
go test -race -count=1 ./...
go vet ./...
go -C scripts/pr-policy test -count=1 ./...
go -C scripts/brand-verify test -count=1 ./...
go -C scripts/brand-verify run . -repo ../..
```

Expected: shipped CLI, schema, restoration, fixtures, policy, and retained brand behavior remain unchanged and green.

## 5. Audit GitHub delivery evidence

Read the S012 issue, its single `cueson Delivery` item, milestone v0.0.0, and the official pull request. During implementation Stage is `In progress`; after publication Stage is `PR review`; Slice is `S012`; default Status is empty; the pull request contains a complete closing reference.

## 6. Validate hosted pull-request behavior

Wait for every current-head CI, CodeQL, release-proof, pull-request policy, Codex, and security result. Address every finding and request exactly one second Codex review only if round one reports findings.

Expected: the pull request is green, fully reviewed, conflict-free, and ready for the operator's final review and merge ritual. No merge, tag, GitHub Release, milestone closure, schema publication, or production mutation has occurred.

## 7. Post-merge candidate binding

After the operator separately authorizes and completes the squash merge, verify the `Release proof` run for the exact new `main` commit and use its evidence as the proposed v0.0.0 publication package.

Expected: the default-branch evidence identifies the squash-merge commit itself. Tag creation, release publication, asset publication, milestone closure, and production work still require their own authorization and are not part of S012 autopilot.
