# Quickstart: Verify the v1.0.0 Release Candidate

## 1. Verify stable release identity and repository records

```powershell
go test -count=1 ./internal/version ./internal/model ./internal/schema ./internal/cli ./internal/convert
go run ./scripts/github-format/main.go .
go -C scripts/docs-verify test -count=1 ./...
go -C scripts/docs-verify run . -repo ../..
git diff --check
```

Expected: executable and schema identity is exactly `1.0.0`; SubRip and WebVTT documents use stable capability declarations; canonical and immutable v1 schema bytes match; governed documentation and release records are UTF-8 without BOM, use repository line endings, and contain no stale current-source development identity.

## 2. Verify schema admission, fixtures, and release policy

```powershell
go test -count=1 ./internal/conformance ./internal/testutil
go -C scripts/corpus-verify test -count=1 ./...
go -C scripts/corpus-verify run . -root ../..
go -C scripts/release-verify test -count=1 ./...
go -C scripts/release-verify test -race -count=1 ./...
go -C scripts/release-verify vet ./...
go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.12
```

Expected: exact schema identity includes `$id`, root instance `$schema`, and `schema_version`; the immutable v1 copy is required and byte-identical; packaged license and notice bytes match the repository; governed fixture hashes match; the workflow uses candidate mode and preserves its read-only, credential-free, non-publishing boundary.

## 3. Build and verify the exact candidate

Install the exact GoReleaser and Syft versions documented in `docs/release-verification.md`, then run from a clean committed checkout:

```powershell
goreleaser check
goreleaser release --snapshot --clean --skip=publish
$sourceRevision = git rev-parse HEAD
go -C scripts/release-verify run . -dist ../../dist -repo ../.. -version 1.0.0 -commit $sourceRevision -execute-host
```

Expected: six archives, six software bills of materials, six checksums, and four members per archive pass; every binary and schema reports `1.0.0`; `dist/release-evidence.json` records the same full commit, intended tag, immutable schema and legal-file digests, complete target digests, native host execution, and `published: false` without development mode.

## 4. Verify all product and repository behavior

```powershell
go test -count=1 ./...
go test -race -count=1 ./...
go vet ./...
go -C scripts/github-format test -count=1 ./...
go -C scripts/pr-policy test -count=1 ./...
go -C scripts/brand-verify test -count=1 ./...
go -C scripts/brand-verify run . -repo ../..
```

Expected: the full CLI, schema, native codecs, restoration, conversion, inspection, completion, conformance, documentation, policy, and retained brand behavior remain green.

## 5. Audit GitHub delivery evidence

Inspect issue #37, its single `cueson Delivery` item, milestone v1.0.0, parent epic #29, dependency on publication issue #38, and the official pull request. During implementation Stage is `In progress`; after pull-request publication Stage is `PR review`; Slice is `S019`; default Status is empty; the pull request contains `Closes #37`.

## 6. Validate hosted pull-request behavior

Wait for every current-head CI, CodeQL, release-proof, native packaged-smoke, pull-request-policy, Codex, security, and review result. Address every finding and request exactly one second Codex review only if round one reports findings or the repository-documented condition requires it.

Expected: the pull request is green, fully reviewed, conflict-free, and ready for the operator's final review and merge ritual. No merge, tag, GitHub Release, milestone closure, schema publication, or production mutation has occurred.

## 7. Post-merge candidate binding

After the operator separately authorizes and completes the squash merge, verify the Release proof run for the exact new `main` commit and use its evidence as the proposed v1.0.0 publication package for issue #38.

Expected: default-branch evidence identifies the squash-merge commit itself. Tag creation, GitHub Release and asset publication, milestone closure, schema hosting, and production work still require explicit authorization and are not part of S019 autopilot.
