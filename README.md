# Cueson

<p align="center">
  <a href="https://github.com/shruggietech/cueson/actions/workflows/ci.yml"><img alt="CI" src="https://github.com/shruggietech/cueson/actions/workflows/ci.yml/badge.svg?branch=main&amp;event=push"></a>
  <a href="https://github.com/shruggietech/cueson/milestone/1"><img alt="Release: unreleased" src="https://img.shields.io/badge/release-unreleased-9CA3AF"></a>
  <a href="LICENSE"><img alt="License" src="https://img.shields.io/badge/license-Apache--2.0-62D9B7"></a>
  <a href="docs/"><img alt="Docs" src="https://img.shields.io/badge/docs-repository-58A6FF"></a>
</p>

**A lossless, structured interchange layer for subtitle and caption content.**<br>**Initial target formats:** SubRip (`.srt`) and WebVTT (`.vtt`)<br>**Status:** Pre-release source foundation

Cueson will convert subtitle and caption formats into and out of a canonical, versioned JSON representation called Cue JSON. Its common cue model is designed for direct use by search, analysis, automation, and AI systems, while a source envelope preserves the original assets for byte-exact restoration.

The repository contains a buildable `cueson` executable and the [canonical Draft 2020-12 Cue JSON `0.0.0` schema](internal/schema/cueson.schema.json). The executable provides truthful help, `cueson version`, embedded schema retrieval, and codec-independent exact restoration from valid source envelopes. It does not yet provide a public release or native subtitle-format ingest and render support.

## Project direction

- Preserve original source bytes, names, hashes, and observable filesystem metadata.
- Expose cue timing, raw text, plain text, line structure, speakers, and token timing through a stable common model.
- Retain format-native structures and surface unsupported or lossy interpretations through diagnostics.
- Keep the Cueson JSON Schema and official executable versioned together.
- Deliver portable native binaries for Windows, macOS, and Linux.

The ratified implementation baselines cover [architecture](docs/architecture.md), the [Cue JSON schema](docs/schema.md), and the [CLI contract](docs/cli.md). The broader [Cueson Project Specification](docs/Cueson-Project-Specification-v0.0.0.md) remains a working roadmap, and the [media-format guide](docs/cueson-media-format-guide.html) explains the longer-term format landscape and product intent.

## Development

Development is specification-driven with [GitHub Spec Kit](https://github.com/github/spec-kit). The current source requires Go 1.25.0 or newer and remains pure Go with native dependencies disabled. The minimum increased when the Go 1.24 line and its compatible text dependency could no longer satisfy the repository's vulnerability gate.

Run the available command directly from source:

```text
go run ./cmd/cueson --help
go run ./cmd/cueson version
go run ./cmd/cueson schema --version
go run ./cmd/cueson schema
go run ./cmd/cueson restore --no-metadata --output restored.srt document.cueson.json
```

The restore command validates canonical base64, byte length, SHA-256, portable names, the complete destination plan, and overwrite safety before accepting output. Timestamp restoration is platform-aware; use `--strict-metadata` to require reproducible captured timestamps or `--no-metadata` to skip metadata application. Pull requests and `main` updates run stable quality, native Windows/macOS/Linux, pure-Go cross-build, vulnerability, and independent CodeQL gates. Repository protection and later codec-specific gates remain separately tracked work.

Run the product tests and build:

```text
go test ./...
go build ./cmd/cueson
```

The repository publication formatter remains an intentionally separate module with its own test command:

```text
go -C scripts/github-format test ./...
```

The non-publishing release proof builds and inspects the complete v0.0.0 candidate matrix through a separate verifier module:

```text
goreleaser release --snapshot --clean --skip=publish
go -C scripts/release-verify run . -dist ../../dist -repo ../.. -version 0.0.0 -commit <full-commit> -execute-host
```

See [release verification](docs/release-verification.md) for exact tool versions, artifact contents, checksums, SBOM expectations, and the boundary between CI evidence and an authorized public release.

See [CONTRIBUTING.md](CONTRIBUTING.md) before proposing changes. Report security concerns privately through [GitHub Security Advisories](https://github.com/shruggietech/cueson/security/advisories/new) rather than a public issue.

## License

Licensed under the [Apache License 2.0](LICENSE). Copyright 2026 [ShruggieTech](https://shruggie.tech).

<sub>A ShruggieTech project.</sub>
