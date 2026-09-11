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
  <a href="https://github.com/shruggietech/cueson/releases/tag/v0.0.0"><img alt="Release" src="https://img.shields.io/github/v/release/shruggietech/cueson?display_name=tag&amp;sort=semver"></a>
  <a href="LICENSE"><img alt="License" src="https://img.shields.io/badge/license-Apache--2.0-62D9B7"></a>
  <a href="docs/"><img alt="Docs" src="https://img.shields.io/badge/docs-repository-58A6FF"></a>
</p>

**A lossless, structured interchange layer for subtitle and caption content.**<br>**Initial stable-release targets:** SubRip (`.srt`) and WebVTT (`.vtt`)<br>**Status:** v0.1.0 development (experimental native SubRip and WebVTT)

Cueson encodes SubRip and WebVTT subtitle files into a canonical, versioned JSON representation called Cue JSON, renders the structured model back to deterministic native syntax, and converts between those formats with explicit loss reporting. Its common cue model is directly usable by search, analysis, automation, and AI systems, while a source envelope preserves original assets for byte-exact restoration.

The repository contains the development `cueson` executable and canonical Draft 2020-12 Cue JSON `0.1.0` schema. The separate [official v0.0.0 release](https://github.com/shruggietech/cueson/releases/tag/v0.0.0) and its byte-identical [immutable schema](schema/releases/v0.0.0/cueson.schema.json) remain the published foundation. Current source adds bounded native SubRip and WebVTT detection, decoding, semantic ingest, deterministic rendering, and codec-independent exact restoration.

## Capability direction

- Preserve original source bytes, names, hashes, and observable filesystem metadata.
- Expose cue timing, raw text, plain text, line structure, speakers, and token timing through a stable common model.
- Retain format-native structures and surface unsupported or lossy interpretations through diagnostics.
- Keep the Cueson JSON Schema and official executable versioned together.
- Deliver portable native binaries for Windows, macOS, and Linux.

The ratified implementation baselines cover [architecture](docs/architecture.md), the [Cue JSON schema](docs/schema.md), and the [CLI contract](docs/cli.md). Dedicated format pages define the experimental native [SubRip](docs/formats/srt.md) and [WebVTT](docs/formats/webvtt.md) contracts on the path to v1. The broader [Cueson Project Specification](docs/Cueson-Project-Specification-v0.0.0.md) remains a working roadmap, and the [media-format guide](docs/cueson-media-format-guide.html) explains the longer-term format landscape and product intent.

The official Cueson identity is retained in the repository as the complete [brand kit](brand/cueson/1.0.0/kit/README.md). See the [Cueson brand guide](docs/brand.md) for provenance, integrity verification, asset selection, licensing boundaries, and the [official ShruggieTech download](https://brand.shruggie.tech/cueson/downloads/cueson-brand-1.0.0.zip).

The dated v0.0.0 history and concise [release notes](docs/releases/v0.0.0.md) are bound to immutable tag [`v0.0.0`](https://github.com/shruggietech/cueson/tree/v0.0.0) at `b294a6952c8bd041d852c502f5d7206c0b58edd6`. [Release verification](docs/release-verification.md) records the accepted six-target proof and independent public-download audit. The governed [release process](docs/release-process.md) keeps milestone closure, signatures, attestations, public schema hosting, and production actions separately authorized.

## Development

Development is specification-driven with [GitHub Spec Kit](https://github.com/github/spec-kit). The current source requires Go 1.25.0 or newer and remains pure Go with native dependencies disabled. The minimum increased when the Go 1.24 line and its compatible text dependency could no longer satisfy the repository's vulnerability gate.

Download the archive and matching SBOM for a supported target from the [v0.0.0 GitHub Release](https://github.com/shruggietech/cueson/releases/tag/v0.0.0), then verify the archive with `cueson_0.0.0_checksums.txt`. That published foundation provides the version, schema, and exact-restore commands:

```text
cueson version
cueson schema --version
cueson schema
cueson restore --no-metadata --output restored.srt document.cueson.json
```

Current development source adds experimental SubRip and WebVTT encode, render, conversion, validation, privacy-bounded inspection, and static shell-completion workflows:

```text
go run ./cmd/cueson --help
go run ./cmd/cueson encode --pretty captions.srt
go run ./cmd/cueson render captions.srt.cueson.json --to srt --output rendered.srt
go run ./cmd/cueson encode --pretty captions.vtt
go run ./cmd/cueson render captions.vtt.cueson.json --to vtt --output rendered.vtt
go run ./cmd/cueson convert captions.srt --to vtt --output captions.vtt
go run ./cmd/cueson convert captions.vtt.cueson.json --to srt --strict --output captions.srt
go run ./cmd/cueson validate captions.srt
go run ./cmd/cueson inspect --json captions.vtt.cueson.json
go run ./cmd/cueson completion powershell
go run ./cmd/cueson restore --no-metadata --output restored.srt document.cueson.json
```

`encode` writes Cue JSON to `INPUT.cueson.json` by default and retains the exact source bytes. SubRip supports explicit encoding selection for ambiguous legacy files; WebVTT accepts UTF-8 only, with an optional UTF-8 BOM. `render --to srt` and `render --to vtt` serialize the structured model as canonical LF syntax and remain intentionally distinct from exact `restore`. `convert` accepts Cue JSON or native subtitle input, writes the requested native format to stdout by default, reports every known omission or degradation on stderr, and rejects any known loss before publication when `--strict` is selected. `validate` checks canonical or native input without producing a payload, while `inspect` reports structural facts without source bytes, content text, or local identifiers. `completion` emits deterministic static definitions for Bash, Zsh, Fish, and PowerShell without modifying a profile.

The default branch is protected by pull-request, resolved-conversation, squash-only, deletion, non-fast-forward, and strict current-base rules. Seventeen GitHub Actions-owned CI and CodeQL checks are required. The two pull-request policy statuses remain visible but are not required checks while the GitHub Actions-authored second-round comment path lacks complete activation proof. Native codec behavior is exercised by the schema-and-conformance, native-test, race-detection, static-analysis, and vulnerability gates.

Run the product tests and build:

```text
go test ./...
go build ./cmd/cueson
```

Repository publication formatting and offline documentation-link verification remain intentionally separate modules:

```text
go -C scripts/github-format test ./...
go -C scripts/docs-verify test ./...
go -C scripts/docs-verify run . -repo ../..
```

The non-publishing release proof remains the acceptance path for candidate builds and inspects the complete six-target matrix through a separate verifier module:

```text
goreleaser release --snapshot --clean --skip=publish
go -C scripts/release-verify run . -dist ../../dist -repo ../.. -version 0.1.0 -commit <full-commit> -development -execute-host
```

See [release verification](docs/release-verification.md) for exact tool versions, artifact contents, checksums, SBOM expectations, the published v0.0.0 evidence, and the boundary between candidate proof and authorized publication. Candidate verification does not publish Cueson by itself.

See [CONTRIBUTING.md](CONTRIBUTING.md) before proposing changes. Report security concerns privately through [GitHub Security Advisories](https://github.com/shruggietech/cueson/security/advisories/new) rather than a public issue.

## License

Licensed under the [Apache License 2.0](LICENSE). Copyright 2026 [ShruggieTech](https://shruggie.tech).

<sub>A ShruggieTech project.</sub>
