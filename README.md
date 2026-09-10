# Cueson

<p align="center">
  <a href="https://github.com/shruggietech/cueson/actions/workflows/ci.yml"><img alt="CI" src="https://github.com/shruggietech/cueson/actions/workflows/ci.yml/badge.svg?branch=main&amp;event=push"></a>
  <a href="https://github.com/shruggietech/cueson/milestone/1"><img alt="Release: unreleased" src="https://img.shields.io/badge/release-unreleased-9CA3AF"></a>
  <a href="LICENSE"><img alt="License" src="https://img.shields.io/badge/license-Apache--2.0-62D9B7"></a>
  <a href="docs/"><img alt="Docs" src="https://img.shields.io/badge/docs-repository-58A6FF"></a>
</p>

**A lossless, structured interchange layer for subtitle and caption content.**<br>**Initial stable-release targets:** SubRip (`.srt`) and WebVTT (`.vtt`)<br>**Status:** Unreleased v0.0.0 foundation candidate

Cueson is designed to convert subtitle and caption formats into and out of a canonical, versioned JSON representation called Cue JSON. Its common cue model is intended for direct use by search, analysis, automation, and AI systems, while a source envelope preserves original assets for byte-exact restoration. Native format conversion remains planned rather than implemented.

The repository contains a buildable `cueson` executable and the [canonical Draft 2020-12 Cue JSON `0.0.0` schema](internal/schema/cueson.schema.json). The executable provides truthful help, `cueson version`, embedded schema retrieval, and codec-independent exact restoration from valid source envelopes. It does not yet provide a public release or native subtitle-format ingest and render support.

## Planned stable direction

- Preserve original source bytes, names, hashes, and observable filesystem metadata.
- Expose cue timing, raw text, plain text, line structure, speakers, and token timing through a stable common model.
- Retain format-native structures and surface unsupported or lossy interpretations through diagnostics.
- Keep the Cueson JSON Schema and official executable versioned together.
- Deliver portable native binaries for Windows, macOS, and Linux.

The ratified implementation baselines cover [architecture](docs/architecture.md), the [Cue JSON schema](docs/schema.md), and the [CLI contract](docs/cli.md). Dedicated format pages separate current envelope-only behavior from the planned [SubRip](docs/formats/srt.md) and [WebVTT](docs/formats/webvtt.md) v1 contracts. The broader [Cueson Project Specification](docs/Cueson-Project-Specification-v0.0.0.md) remains a working roadmap, and the [media-format guide](docs/cueson-media-format-guide.html) explains the longer-term format landscape and product intent.

Release-candidate construction and inspection are documented in [release verification](docs/release-verification.md). The separately governed [release process](docs/release-process.md) describes future tag, GitHub Release, immutable schema-copy, milestone, and production-publication actions without authorizing them.

## Development

Development is specification-driven with [GitHub Spec Kit](https://github.com/github/spec-kit). The current source requires Go 1.25.0 or newer and remains pure Go with native dependencies disabled. The minimum increased when the Go 1.24 line and its compatible text dependency could no longer satisfy the repository's vulnerability gate.

No installable release or supported binary distribution exists yet. Run the implemented v0.0.0 command surface directly from source:

```text
go run ./cmd/cueson --help
go run ./cmd/cueson version
go run ./cmd/cueson schema --version
go run ./cmd/cueson schema
go run ./cmd/cueson restore --no-metadata --output restored.srt document.cueson.json
```

The restore command accepts Cue JSON with an already-populated source envelope. It validates canonical base64, byte length, SHA-256, portable names, the complete destination plan, and overwrite safety before accepting output. It does not ingest an SRT or WebVTT file, derive semantic cues, render from the model, or convert between formats. Timestamp restoration is platform-aware; use `--strict-metadata` to require reproducible captured timestamps or `--no-metadata` to skip metadata application.

The default branch is protected by pull-request, resolved-conversation, squash-only, deletion, non-fast-forward, and strict current-base rules. Seventeen GitHub Actions-owned CI and CodeQL checks are required. The two pull-request policy statuses remain visible but are not required checks while the GitHub Actions-authored second-round comment path lacks complete activation proof. Native codec-specific gates do not exist because native codecs are not implemented.

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

The non-publishing release proof builds and inspects the complete v0.0.0 candidate matrix through a separate verifier module:

```text
goreleaser release --snapshot --clean --skip=publish
go -C scripts/release-verify run . -dist ../../dist -repo ../.. -version 0.0.0 -commit <full-commit> -execute-host
```

See [release verification](docs/release-verification.md) for exact tool versions, artifact contents, checksums, SBOM expectations, and the boundary between CI evidence and an authorized public release. Candidate verification does not install or publish Cueson.

See [CONTRIBUTING.md](CONTRIBUTING.md) before proposing changes. Report security concerns privately through [GitHub Security Advisories](https://github.com/shruggietech/cueson/security/advisories/new) rather than a public issue.

## License

Licensed under the [Apache License 2.0](LICENSE). Copyright 2026 [ShruggieTech](https://shruggie.tech).

<sub>A ShruggieTech project.</sub>
