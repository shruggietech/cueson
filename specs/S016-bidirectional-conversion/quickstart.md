# Quickstart: Verify Bidirectional Subtitle Conversion

## Prerequisites

- Go 1.25 from the repository toolchain declaration
- A clean checkout on `codex/S016-bidirectional-conversion`
- No production credentials or network services

## Build and help

```powershell
New-Item -ItemType Directory -Force .tmp | Out-Null
go build -o .tmp/cueson.exe ./cmd/cueson
./.tmp/cueson.exe --help
./.tmp/cueson.exe convert --help
```

Expected: root help lists `convert`, command help documents the exact options and stream rules in [contracts/cli.md](contracts/cli.md), and no unimplemented command appears.

## Convert native sources

```powershell
./.tmp/cueson.exe convert testdata/fixtures/conversion/srt-loss-free/source/input.srt --to vtt --output .tmp/minimal.vtt
./.tmp/cueson.exe convert testdata/fixtures/conversion/webvtt-loss-free/source/input.vtt --to srt --output .tmp/minimal.srt
./.tmp/cueson.exe encode .tmp/minimal.vtt --output .tmp/minimal-vtt.cueson.json
./.tmp/cueson.exe encode .tmp/minimal.srt --output .tmp/minimal-srt.cueson.json
```

Expected: both conversions and target re-encodes succeed, converted files use canonical LF syntax, and target Cue JSON validates.

## Convert Cue JSON

```powershell
./.tmp/cueson.exe encode testdata/fixtures/conversion/srt-loss-free/source/input.srt --output .tmp/source.cueson.json
./.tmp/cueson.exe convert .tmp/source.cueson.json --from cueson --to vtt --output .tmp/from-json.vtt
```

Expected: conversion uses the structured model, does not restore source bytes, and produces the same representable target semantics as direct source conversion.

## Exercise loss and strict policy

```powershell
./.tmp/cueson.exe convert testdata/fixtures/conversion/webvtt-lossy/source/input.vtt --to srt --output .tmp/lossy.srt
./.tmp/cueson.exe convert testdata/fixtures/conversion/webvtt-lossy/source/input.vtt --to srt --strict --output .tmp/strict.srt
```

Expected: normal mode writes parser-valid SubRip and emits the ordered losses defined by [contracts/compatibility.md](contracts/compatibility.md); strict mode exits 1 and leaves `.tmp/strict.srt` absent.

## Run focused and full verification

```powershell
go test ./internal/convert ./internal/cli ./internal/conformance
go test ./...
go test -race ./...
go vet ./...
```

Expected: conversion matrices, loss-report goldens, parser cycles, transactions, path-leak checks, existing workflows, and race evidence all pass.
