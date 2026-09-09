# Cueson

<p align="center">
  <a href="https://github.com/shruggietech/cueson/actions/workflows/ci.yml"><img alt="CI" src="https://github.com/shruggietech/cueson/actions/workflows/ci.yml/badge.svg"></a>
  <a href="https://github.com/shruggietech/cueson/releases/latest"><img alt="Release" src="https://img.shields.io/github/v/release/shruggietech/cueson?color=58A6FF"></a>
  <a href="LICENSE"><img alt="License" src="https://img.shields.io/badge/license-Apache--2.0-62D9B7"></a>
  <a href="docs/"><img alt="Docs" src="https://img.shields.io/badge/docs-repository-58A6FF"></a>
</p>

**A lossless, structured interchange layer for subtitle and caption content.**<br>**Initial stable formats:** SubRip (`.srt`) and WebVTT (`.vtt`)<br>**Status:** Pre-release repository bootstrap

Cueson will convert subtitle and caption formats into and out of a canonical, versioned JSON representation called Cue JSON. Its common cue model is designed for direct use by search, analysis, automation, and AI systems, while a source envelope preserves the original assets for byte-exact restoration.

The repository is currently establishing its specification, governance, and delivery system. It does not yet ship a working `cueson` executable or claim complete format support.

## Project direction

- Preserve original source bytes, names, hashes, and observable filesystem metadata.
- Expose cue timing, raw text, plain text, line structure, speakers, and token timing through a stable common model.
- Retain format-native structures and surface unsupported or lossy interpretations through diagnostics.
- Keep the Cueson JSON Schema and official executable versioned together.
- Deliver portable native binaries for Windows, macOS, and Linux.

The current architecture and repository draft is in [Cueson Project Specification](docs/Cueson-Project-Specification-v0.0.0.md). The [media-format guide](docs/cueson-media-format-guide.html) explains the longer-term format landscape and product intent.

## Development

Development is specification-driven with [GitHub Spec Kit](https://github.com/github/spec-kit). Product implementation begins only after the initial repository and GitHub planning foundation is in place.

See [CONTRIBUTING.md](CONTRIBUTING.md) before proposing changes. Report security concerns privately through [GitHub Security Advisories](https://github.com/shruggietech/cueson/security/advisories/new) rather than a public issue.

## License

Licensed under the [Apache License 2.0](LICENSE). Copyright 2026 [ShruggieTech](https://shruggie.tech).

<sub>A ShruggieTech project.</sub>
