# Cueson

<p align="center">
  <a href="https://github.com/shruggietech/cueson/issues/8"><img alt="CI: planned" src="https://img.shields.io/badge/CI-planned-9CA3AF"></a>
  <a href="https://github.com/shruggietech/cueson/milestone/1"><img alt="Release: unreleased" src="https://img.shields.io/badge/release-unreleased-9CA3AF"></a>
  <a href="LICENSE"><img alt="License" src="https://img.shields.io/badge/license-Apache--2.0-62D9B7"></a>
  <a href="docs/"><img alt="Docs" src="https://img.shields.io/badge/docs-repository-58A6FF"></a>
</p>

**A lossless, structured interchange layer for subtitle and caption content.**<br>**Initial target formats:** SubRip (`.srt`) and WebVTT (`.vtt`)<br>**Status:** Pre-release executable foundation

Cueson will convert subtitle and caption formats into and out of a canonical, versioned JSON representation called Cue JSON. Its common cue model is designed for direct use by search, analysis, automation, and AI systems, while a source envelope preserves the original assets for byte-exact restoration.

The repository contains the first buildable `cueson` command-line foundation. It currently provides truthful help and `cueson version`; it does not yet provide a public release, an embedded schema, source restoration, or native subtitle-format support.

## Project direction

- Preserve original source bytes, names, hashes, and observable filesystem metadata.
- Expose cue timing, raw text, plain text, line structure, speakers, and token timing through a stable common model.
- Retain format-native structures and surface unsupported or lossy interpretations through diagnostics.
- Keep the Cueson JSON Schema and official executable versioned together.
- Deliver portable native binaries for Windows, macOS, and Linux.

The ratified implementation baselines cover [architecture](docs/architecture.md), the [Cue JSON schema](docs/schema.md), and the [CLI contract](docs/cli.md). The broader [Cueson Project Specification](docs/Cueson-Project-Specification-v0.0.0.md) remains a working roadmap, and the [media-format guide](docs/cueson-media-format-guide.html) explains the longer-term format landscape and product intent.

## Development

Development is specification-driven with [GitHub Spec Kit](https://github.com/github/spec-kit). The current source requires Go 1.24.0 or newer and uses no third-party runtime dependencies.

Run the available command directly from source:

```text
go run ./cmd/cueson --help
go run ./cmd/cueson version
```

Run the product tests and build:

```text
go test ./...
go build ./cmd/cueson
```

The repository publication formatter remains an intentionally separate module with its own test command:

```text
go -C scripts/github-format test ./...
```

See [CONTRIBUTING.md](CONTRIBUTING.md) before proposing changes. Report security concerns privately through [GitHub Security Advisories](https://github.com/shruggietech/cueson/security/advisories/new) rather than a public issue.

## License

Licensed under the [Apache License 2.0](LICENSE). Copyright 2026 [ShruggieTech](https://shruggie.tech).

<sub>A ShruggieTech project.</sub>
