# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initialized the repository, Spec Kit integration, governance documents, publication tooling, and GitHub planning foundation.
- Added the v0.0.0 architecture, Cue JSON schema, and CLI implementation baselines with their Spec Kit decision and verification artifacts.
- Added the dependency-free Go module and first buildable `cueson` executable with truthful root help, exact `cueson version` output, stable stream and exit-code handling, global diagnostic policy, and focused tests.
- Added the canonical Draft 2020-12 Cue JSON `0.0.0` schema, representative contract document, internal common model and semantic validation, embedded schema access, software/schema lockstep checks, and the `schema` command with explicit safe output replacement.
- Added canonical source-integrity validation, portable canonical-caseless collision protection, rollback-safe staged restoration, capture-before-read timestamp boundaries, native Windows/Linux/macOS timestamp adapters, and the public codec-independent `restore` command.
- Added a manifest-governed, byte-stable fixture corpus with explicit provenance and redistribution decisions, reusable path-free golden comparisons, accepted and malformed cross-package conformance cases, and bounded fuzz boundaries for Cue JSON and source-envelope validation.
- Added stable least-privilege CI and independent CodeQL workflows with pinned actions and analysis tools, native Windows/macOS/Linux tests, race and conformance gates, vulnerability scanning, and six-target pure-Go build proof.

### Changed

- Replaced unavailable dynamic CI and release badges with truthful planned and unreleased state badges until those public resources exist.
- Changed the v0.0.0 envelope capability to advertise generic exact restoration while native ingest, model-driven render, and OCR-required capabilities remain false.
- Replaced the planned CI badge with the live workflow badge and refined the draft CI roadmap so codec, conversion, review-automation, and release gates activate only when their owning implementations exist.
- Raised the minimum Go version from 1.24 to 1.25 and upgraded `golang.org/x/text` to the first compatible fixed release after the vulnerability gate found reachable issues that could not be corrected on the old floor.

### Fixed

- Made fixture-alias and closed-descriptor tests portable across case-insensitive macOS filesystems and Linux timestamp syscall behavior exposed by the first hosted native matrix.

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
- 2026-09-09: Set `restore_supported` true in S004 only after canonical integrity checks, complete destination planning, safe publication, rollback, and the public restore command pass together.
- 2026-09-09: Keep hosted native Windows, Linux, and macOS execution in downstream issue #8 while S004 provides native-selectable tests, current-host Windows proof, and CGO-disabled foreign-platform compilation to avoid a circular dependency.
- 2026-09-09: Govern root test payloads through one strict provenance and integrity manifest, keep generic comparisons domain-neutral under `internal/testutil`, and reserve codec grammar fixtures plus hosted execution for their downstream slices.
- 2026-09-09: Establish immutable-action, least-privilege hosted gates with native three-platform tests and separate six-target pure-Go build proof; defer empty codec, conversion, review-automation, release, and repository-control checks to their owning issues.
- 2026-09-09: Raise the compatibility floor to Go 1.25 rather than suppress reachable Go standard-library and Unicode-normalization vulnerabilities discovered while establishing the mandatory vulnerability gate.

## [0.0.0]

No public release has been published.

[Unreleased]: https://github.com/shruggietech/cueson/compare/v0.0.0...HEAD
[0.0.0]: https://github.com/shruggietech/cueson/releases/tag/v0.0.0
