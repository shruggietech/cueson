# Quickstart: Verify the v1 Contract Freeze

## Prerequisites

- Go 1.25.0 or newer.
- A clean checkout on `codex/S018-freeze-v1-contracts`.
- PowerShell 7 for the Windows examples.
- Exact release tools only for the final non-publishing snapshot.

## 1. Verify specification and repository text

```powershell
go run ./scripts/github-format/main.go .
go -C scripts/docs-verify test -count=1 ./...
go -C scripts/docs-verify run . -repo ../..
go -C scripts/brand-verify test -count=1 ./...
go -C scripts/brand-verify run . -repo ../..
git diff --check
```

Expected result: every command exits 0, all maintained links resolve, no encoding or mojibake violation appears, and retained brand evidence is unchanged.

## 2. Verify schema annotations and immutable history

```powershell
go test -count=1 ./internal/schema
go run ./cmd/cueson schema > $env:TEMP/cueson-s018.schema.json
```

Compare the emitted file byte-for-byte with `internal/schema/cueson.schema.json`. Confirm `schema/releases/v0.0.0/cueson.schema.json` still has SHA-256 `d15c7fa5227156109dd6be3d39b711aca3503794bb862169dfca96ee80adb975`.

## 3. Verify conformance and executable documentation

```powershell
go test -count=1 ./internal/conformance ./internal/cli ./internal/codec/... ./internal/convert ./internal/source
```

Expected result: the matrix resolves every format row, hostile inputs stop at their declared bounds, accepted sources restore byte-for-byte, examples run only in temporary directories, and output privacy checks pass.

## 4. Run bounded fuzz work

Run the fixed-work fuzz commands documented by the S018 task implementation and CI workflow for all required surfaces with one worker. Each target must complete without panic, invariant failure, uncontrolled output, or a new uncommitted regression.

## 5. Run the optional external corpus verifier

```powershell
go -C scripts/corpus-verify run . -root C:\path\to\private\corpus
```

Expected result: the verifier reads supported regular files without changing them, rejects links and oversized inputs, reports portable relative identities or aggregate counts only, and creates no repository artifact. This step is optional and is not a hosted-CI prerequisite.

## 6. Run the complete foreground gate

```powershell
go test -count=1 ./...
go test -race -count=1 ./...
go vet ./...
go run honnef.co/go/tools/cmd/staticcheck@v0.7.0 ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...
```

Repeat tests, vet, Staticcheck, and vulnerability checks for every module under `scripts/`. Build all six supported `CGO_ENABLED=0` targets. Then run the exact non-publishing snapshot procedure in `docs/release-verification.md` against the full current commit.

## 7. Confirm publication boundaries

The successful S018 proof must retain development version `0.1.0`, produce no immutable v1 schema, tag, GitHub Release, public schema, production-domain change, milestone closure, or pull-request merge.
