# Quickstart: Verify the Official Cueson Brand Kit Import

## 1. Verify repository and publication text

```powershell
go run ./scripts/github-format/main.go .
go -C scripts/github-format test -count=1 ./...
git diff --check
```

Expected: repository-authored text passes, PowerShell remains CRLF-governed, and byte-protected brand trees are excluded from check and repair behavior.

## 2. Verify the retained acquisition offline

```powershell
go -C scripts/brand-verify test -count=1 ./...
go -C scripts/brand-verify run . -repo ../..
```

Expected: the retained ZIP matches 3,889,985 bytes and SHA-256 `095e572a0db73472c9fdcf5001f68e702105a3d0ec1fcc34f830ae57c608cd07`; all 265 safe entries match the import manifest and exact extraction; every declared README and media-guide reference resolves to a retained asset.

## 3. Verify documentation and local rendering

```powershell
go -C scripts/docs-verify test -count=1 ./...
go -C scripts/docs-verify run . -repo ../..
```

Open `README.md` through GitHub-compatible light and dark rendering and open `docs/cueson-media-format-guide.html` offline at desktop, narrow, and print widths.

Expected: light and dark README variants remain readable, the guide loads official local fonts, favicon, identity and palette without a network request, the truthful v0.0.0 disclaimer remains present, and all maintained documentation links resolve.

## 4. Verify existing product and repository behavior

```powershell
go test -count=1 ./...
go test -race -count=1 ./...
go vet ./...
go -C scripts/pr-policy test -count=1 ./...
go -C scripts/release-verify test -count=1 ./...
go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.12
```

Expected: shipped Cueson behavior remains unchanged, release verification still requires exactly four archive members, and workflow syntax remains valid.

## 5. Audit GitHub delivery evidence

Read issue #23, its single `cueson Delivery` item, the v1.0.0 milestone, and the official S011 pull request. During implementation Stage is `In progress`; after publication Stage is `PR review`; Slice is `S011`; default Status is empty; the pull request contains `Closes #23`.

## 6. Validate hosted behavior

Wait for every current-head CI, CodeQL, pull-request policy, Codex, and security result. Address every finding and request exactly one second Codex review only if round one reports findings.

Expected: the pull request is green, fully reviewed, conflict-free, and ready for the operator's final review and merge ritual. S011 has not merged, tagged, released, published a schema, closed a milestone, or changed production state.
