# Cueson

<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="brand/cueson/1.0.0/kit/logos/svg/cueson-horizontal-color.svg">
    <source media="(prefers-color-scheme: light)" srcset="brand/cueson/1.0.0/kit/logos/svg/cueson-horizontal-light.svg">
    <img alt="Cueson" src="brand/cueson/1.0.0/kit/logos/svg/cueson-horizontal-light.svg" width="520">
  </picture>
</p>

<p align="center"><strong>Universal captions and subtitles</strong></p>

<p align="center">
  <a href="https://github.com/shruggietech/cueson/actions/workflows/ci.yml"><img alt="CI" src="https://github.com/shruggietech/cueson/actions/workflows/ci.yml/badge.svg?branch=main&amp;event=push"></a>
  <a href="https://github.com/shruggietech/cueson/releases/tag/v1.0.0"><img alt="Release" src="https://img.shields.io/github/v/release/shruggietech/cueson?display_name=tag&amp;sort=semver"></a>
  <a href="LICENSE"><img alt="License" src="https://img.shields.io/badge/license-Apache--2.0-62D9B7"></a>
  <a href="docs/"><img alt="Docs" src="https://img.shields.io/badge/docs-repository-58A6FF"></a>
</p>

**A lossless, structured interchange layer for subtitle and caption content.**<br>**v1 contract formats:** SubRip (`.srt`) and WebVTT (`.vtt`)<br>**Status:** v1.0.0 released and independently verified

Cueson encodes SubRip and WebVTT subtitle files into a canonical, versioned JSON representation called Cue JSON, renders the structured model back to deterministic native syntax, and converts between those formats with explicit loss reporting. Its common cue model is directly usable by search, analysis, automation, and AI systems, while a source envelope preserves original assets for byte-exact restoration.

The [official v1.0.0 release](https://github.com/shruggietech/cueson/releases/tag/v1.0.0) contains the stable `cueson` executable, its canonical Draft 2020-12 Cue JSON schema, and a byte-identical [immutable v1 schema](schema/releases/v1.0.0/cueson.schema.json). It provides bounded native SubRip and WebVTT detection, decoding, semantic ingest, deterministic rendering, cross-format conversion, validation, inspection, shell completion, and codec-independent exact restoration. The earlier [v0.0.0 envelope-only release](https://github.com/shruggietech/cueson/releases/tag/v0.0.0) and its [immutable schema](schema/releases/v0.0.0/cueson.schema.json) remain available as historical foundations. The public product and documentation home is [cueson.io](https://cueson.io).

## Capability direction

- Preserve original source bytes, names, hashes, and observable filesystem metadata.
- Expose cue timing, raw text, plain text, line structure, speakers, and token timing through a stable common model.
- Retain format-native structures and surface unsupported or lossy interpretations through diagnostics.
- Keep the Cueson JSON Schema and official executable versioned together.
- Deliver portable native binaries for Windows, macOS, and Linux.

The ratified implementation baselines cover [architecture](docs/architecture.md), the [Cue JSON schema](docs/schema.md), and the [CLI contract](docs/cli.md). Dedicated format pages define the stable v1 [SubRip](docs/formats/srt.md) and [WebVTT](docs/formats/webvtt.md) contracts. The broader [Cueson Project Specification](docs/Cueson-Project-Specification-v0.0.0.md) remains a working roadmap, and the [media-format guide](docs/cueson-media-format-guide.html) explains the longer-term format landscape and product intent.

The official Cueson identity is retained in the repository as the complete [brand kit](brand/cueson/1.0.0/kit/README.md). See the [Cueson brand guide](docs/brand.md) for provenance, integrity verification, asset selection, licensing boundaries, and the [official ShruggieTech download](https://brand.shruggie.tech/cueson/downloads/cueson-brand-1.0.0.zip).

The v1.0.0 history, immutable schema, and concise [release notes](docs/releases/v1.0.0.md) are bound to immutable annotated tag [`v1.0.0`](https://github.com/shruggietech/cueson/tree/v1.0.0) at `2cad4c816340404289b4d1d87179a4071713bb46`. The [GitHub Release](https://github.com/shruggietech/cueson/releases/tag/v1.0.0) publishes the exact thirteen files accepted by default-branch proof run [34621429626](https://github.com/shruggietech/cueson/actions/runs/34621429626). [Release verification](docs/release-verification.md) records the independent public-byte check, and the governed [release process](docs/release-process.md) keeps milestone closure, signatures, attestations, public schema hosting, and production actions separately authorized.

## Installation

Download the archive for your platform from the [v1.0.0 GitHub Release](https://github.com/shruggietech/cueson/releases/tag/v1.0.0) and verify it with `cueson_1.0.0_checksums.txt`. Six pure-Go archives cover Windows, macOS, and Linux on amd64 and arm64, and each archive has a matching SPDX JSON SBOM.

Developers building from source require Go 1.25.0 or newer. Clone the repository, then use `go run ./cmd/cueson`, `go build ./cmd/cueson`, or `go install ./cmd/cueson`. The published release remains pure Go with `CGO_ENABLED=0`. Versioned public schemas resolve at `https://cueson.io/schema/vVERSION/cueson.schema.json`; the schema embedded in each executable and archive remains available without network access.

## Executable quick start

Run these examples from the repository root. Create an empty `quickstart` directory first. Every input below is a committed, redistributable fixture, and every created file stays beneath that scratch directory. Existing destinations are refused unless `--force` is explicitly supplied.

<!-- docs-verify:example source-run -->

```text
go run ./cmd/cueson version
```

Expected result: exit status 0, stdout is exactly `1.0.0` plus LF, and stderr is empty.

<!-- docs-verify:example encode -->

```text
go run ./cmd/cueson encode --pretty --output quickstart/document.cueson.json testdata/fixtures/conversion/srt-loss-free/source/input.srt
```

Expected result: exit status 0, stdout and stderr are empty, and `quickstart/document.cueson.json` is valid Cue JSON containing the exact source bytes in its source envelope.

<!-- docs-verify:example restore -->

```text
go run ./cmd/cueson restore --no-metadata --output quickstart/restored.srt quickstart/document.cueson.json
```

Expected result: exit status 0 with empty stdout and stderr. `quickstart/restored.srt` is byte-for-byte identical to `testdata/fixtures/conversion/srt-loss-free/source/input.srt`; `restore` does not render the structured cue model.

<!-- docs-verify:example render -->

```text
go run ./cmd/cueson render --to srt --output quickstart/rendered.srt quickstart/document.cueson.json
```

Expected result: exit status 0 with empty stdout and stderr. `quickstart/rendered.srt` is canonical LF SubRip produced from the structured model and is not claimed to preserve the source's original bytes.

<!-- docs-verify:example convert-srt-vtt -->

```text
go run ./cmd/cueson convert --strict --no-speaker-detection --to vtt --output quickstart/converted.vtt testdata/fixtures/conversion/srt-loss-free/source/input.srt
```

<!-- docs-verify:example convert-vtt-srt -->

```text
go run ./cmd/cueson convert --strict --to srt --output quickstart/converted.srt testdata/fixtures/conversion/webvtt-loss-free/source/input.vtt
```

Both governed inputs are explicitly loss-free for the requested direction. Each command exits 0 with empty stdout and stderr and writes parser-valid native output. For other inputs, normal conversion reports every known loss on stderr; `--strict` rejects the complete conversion before publishing any output when one or more losses exist. See the [conversion contract](docs/conversion.md).

<!-- docs-verify:example validate -->

```text
go run ./cmd/cueson validate testdata/fixtures/webvtt/minimal/source/minimal.vtt
```

Expected result: exit status 0, empty stdout, and one successful-validation diagnostic on stderr. `--quiet` suppresses that success diagnostic; `--silent` also suppresses warnings.

<!-- docs-verify:example inspect-human -->

```text
go run ./cmd/cueson inspect testdata/fixtures/webvtt/minimal/source/minimal.vtt
```

<!-- docs-verify:example inspect-json -->

```text
go run ./cmd/cueson inspect --json testdata/fixtures/webvtt/minimal/source/minimal.vtt
```

Each inspection command exits 0 with its report on stdout and no stderr for this fixture. Reports contain bounded structural facts, not preserved source bytes, content text, caller paths, usernames, hostnames, or other local identifiers.

<!-- docs-verify:example completion -->

```text
go run ./cmd/cueson completion bash > quickstart/cueson.bash
```

Expected result: exit status 0 and one deterministic UTF-8 LF Bash completion definition on stdout, redirected here to the scratch file. Cueson does not modify a shell profile.

The complete command, option, alias, stream, exit-status, overwrite, and strict-mode behavior is in the [CLI contract](docs/cli.md). The [compatibility contract](docs/compatibility.md) defines the stable v1 boundary, platform targets, immutable release surfaces, and separately deferred production hosting.

## Development

Development is specification-driven with [GitHub Spec Kit](https://github.com/github/spec-kit). Run the product tests and build with `go test ./...` and `go build ./cmd/cueson`. Repository publication formatting, executable-documentation checks, and offline link verification run through `scripts/github-format` and `scripts/docs-verify` as described in [CONTRIBUTING.md](CONTRIBUTING.md).

The default branch is protected by pull-request, resolved-conversation, squash-only, deletion, non-fast-forward, and strict current-base rules. Required CI and CodeQL checks cover schema and conformance, native platforms, race detection, static analysis, vulnerability analysis, and pure-Go cross-builds. The non-publishing release proof and its exact commands are documented in [release verification](docs/release-verification.md); successful candidate verification does not publish Cueson.

Report security concerns privately through [GitHub Security Advisories](https://github.com/shruggietech/cueson/security/advisories/new) rather than a public issue.

## License

Licensed under the [Apache License 2.0](LICENSE). Copyright 2026 [ShruggieTech](https://shruggie.tech).

<sub>A ShruggieTech project.</sub>
