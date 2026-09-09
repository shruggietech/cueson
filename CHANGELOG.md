# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initialized the repository, Spec Kit integration, governance documents, publication tooling, and GitHub planning foundation.
- Added the v0.0.0 architecture, Cue JSON schema, and CLI implementation baselines with their Spec Kit decision and verification artifacts.
- Added the dependency-free Go module and first buildable `cueson` executable with truthful root help, exact `cueson version` output, stable stream and exit-code handling, global diagnostic policy, and focused tests.
- Added the canonical Draft 2020-12 Cue JSON `0.0.0` schema, representative contract document, internal common model and semantic validation, embedded schema access, software/schema lockstep checks, and the `schema` command with explicit safe output replacement.

### Changed

- Replaced unavailable dynamic CI and release badges with truthful planned and unreleased state badges until those public resources exist.

### Decisions

- Treat the project specification as a working pre-release draft until its contracts are ratified through the constitution and implementation slices.
- Keep repository publication utilities separate from the shipped Cueson product.
- 2026-09-09: Use `subrip` and `webvtt` as canonical schema keys and declare the completed v0.0.0 milestone `envelope_only`, with public generic restoration owned by the source-foundation slice and native ingest/render deferred.
- 2026-09-09: Require an `ocr_observations` array on every cue and keep every observation derived, independently provenanced, and subordinate to source truth.
- 2026-09-09: Show only implemented CLI commands in help, classify invocation and pre-execution failures as exit code 2, and classify missing runtime capability as exit code 1.
- 2026-09-09: Define one portable source-asset basename and collision contract, enforce single-name safety in S003, and defer normalization or case-fold collision enforcement before restoration to issue #6.
- 2026-09-09: Treat stored basenames as default and `--output-dir` names while allowing single-asset `--output` to select a separately validated runtime path.
- 2026-09-09: Keep bootstrap issue #1 open until the badge correction reaches and is verified on the default branch.
- 2026-09-09: Keep `restore_supported` false throughout S003 because schema recognition and a preservation envelope do not constitute the public restoration capability owned by issue #6.

## [0.0.0]

No public release has been published.

[Unreleased]: https://github.com/shruggietech/cueson/compare/v0.0.0...HEAD
[0.0.0]: https://github.com/shruggietech/cueson/releases/tag/v0.0.0
