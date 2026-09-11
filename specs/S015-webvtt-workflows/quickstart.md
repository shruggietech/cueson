# Quickstart: Native WebVTT Workflows

## Prerequisites

- Go 1.25 or newer
- GoReleaser 2.18.1 and Syft 1.51.1 for the development release proof
- A clean checkout on `codex/S015-webvtt-workflows`

## Encode, render, re-encode, and restore

```powershell
New-Item -ItemType Directory -Force ./tmp | Out-Null
go run ./cmd/cueson encode ./testdata/fixtures/webvtt/complete/source/complete.vtt --output ./tmp/complete.cueson.json
go run ./cmd/cueson render ./tmp/complete.cueson.json --to vtt --output ./tmp/rendered.vtt
go run ./cmd/cueson encode ./tmp/rendered.vtt --stdout | Out-Null
go run ./cmd/cueson restore ./tmp/complete.cueson.json --output ./tmp/restored.vtt
```

Expected results: all commands exit zero; Cue JSON validates; rendered WebVTT is accepted by native encode; restored bytes match the original source; rendered bytes represent the structured model and are not required to match the source.

## Standard-output isolation

```powershell
go run ./cmd/cueson encode ./testdata/fixtures/webvtt/minimal/source/minimal.vtt --stdout 1> ./tmp/stdout.json 2> ./tmp/diagnostics.txt
go run ./cmd/cueson render ./tmp/stdout.json --to vtt 1> ./tmp/stdout.vtt 2> ./tmp/render-diagnostics.txt
```

Expected results: standard output contains only one valid Cue JSON document or WebVTT stream. Permitted warnings appear only on standard error.

## Verification

```powershell
go test ./internal/codec/webvtt ./internal/model ./internal/schema ./internal/cli ./internal/conformance
go test -race ./...
go test ./...
goreleaser release --snapshot --clean --skip=publish
$sourceCommit = git rev-parse HEAD
go -C scripts/release-verify run . -dist ../../dist -repo ../.. -version 0.1.0 -commit $sourceCommit -development -execute-host
```

Expected results: all focused, race, full, and development release-proof checks pass without changing immutable v0.0.0 release material.
