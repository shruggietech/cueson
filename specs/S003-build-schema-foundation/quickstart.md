# Quickstart: Build Schema Foundation

## Prerequisites

- Go 1.24 or newer is available.
- Commands run from the repository root.
- No network access is required after module dependencies are available.

## Format and static checks

```powershell
gofmt -w ./cmd ./internal ./scripts
go vet ./...
```

Expected result: formatting makes no remaining source diff and vet reports no finding.

## Schema and unit verification

```powershell
go test ./...
go test -race ./...
```

Expected result: the schema compiles as Draft 2020-12, the representative document validates, invalid structural and semantic variants fail, naming and lockstep checks pass, and all CLI tests pass.

## Pure-Go build

```powershell
$env:CGO_ENABLED = '0'
go build ./cmd/cueson
Remove-Item -LiteralPath './cueson.exe' -ErrorAction SilentlyContinue
Remove-Item Env:CGO_ENABLED
```

Expected result: the executable builds without a native runtime dependency and the generated binary is removed after verification.

## CLI smoke checks

```powershell
go run ./cmd/cueson schema --version
go run ./cmd/cueson schema
go run ./cmd/cueson --help
```

Expected result: version output is exactly `0.0.0` plus LF, schema output is valid JSON identical to the embedded artifact, and help lists only `version` and `schema`.

For file output, select a disposable path, verify refusal without force, verify replacement with force, and remove the disposable file after comparison.

## Repository checks

```powershell
go -C scripts/github-format test ./...
git diff --check
git status --short
```

Inspect changed text for UTF-8 without BOM, expected line endings, and mojibake before publication.

## Spec Kit completion gates

Confirm the requirements checklist is complete, run cross-artifact analysis before implementation, mark every implementation task complete, then run convergence. Publication proceeds only after these gates pass.
