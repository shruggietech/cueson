# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Reconciled repository release status and verification records after publishing the immutable v0.0.0 tag and the exact thirteen-asset GitHub Release from accepted default-branch evidence.
- Advanced the evolving executable and canonical Cue JSON contract to development version 0.1.0 while preserving every published v0.0.0 schema and release artifact unchanged.

### Added

- Added bounded native SubRip detection and text decoding, complete documented cue parsing, exact source-envelope preservation, deterministic model-driven SubRip rendering, and public `encode` and `render` commands.
- Added bounded native WebVTT detection and UTF-8 decoding, source-ordered cue and block parsing, setting and markup fidelity, deterministic diagnostics, exact source-envelope preservation, model-driven rendering, governed fixtures, and CLI workflows.
- Added a capability-based internal codec registry for experimental native SubRip and WebVTT support.

### Fixed

- Made shared schema, encode, and render file publication transactional across short writes, failed writes, new destinations, and forced replacements so a failed operation does not leave partial or lost output.

### Decisions

- 2026-09-11: Advance development schema and software identity together to 0.1.0 because native capability constraints change canonical schema bytes; keep the released v0.0.0 identity immutable.
- 2026-09-11: Reject reversed SubRip timing instead of silently swapping endpoints, and retain speaker prefixes in raw and plain payload views while exposing speaker identity only as a derived observation.
- 2026-09-11: Keep WebVTT source bytes authoritative while representing cues and non-cue blocks in one contiguous source order, retain ordered raw setting occurrences beside effective values, and derive plain text, speakers, and inline timing without replacing native payload syntax.

## [0.0.0] - 2026-09-10

### Added

- Initialized the repository, Spec Kit integration, governance documents, publication tooling, and GitHub planning foundation.
- Added the v0.0.0 architecture, Cue JSON schema, and CLI implementation baselines with their Spec Kit decision and verification artifacts.
- Added the dependency-free Go module and first buildable `cueson` executable with truthful root help, exact `cueson version` output, stable stream and exit-code handling, global diagnostic policy, and focused tests.
- Added the canonical Draft 2020-12 Cue JSON `0.0.0` schema, representative contract document, internal common model and semantic validation, embedded schema access, software/schema lockstep checks, and the `schema` command with explicit safe output replacement.
- Added canonical source-integrity validation, portable canonical-caseless collision protection, rollback-safe staged restoration, capture-before-read timestamp boundaries, native Windows/Linux/macOS timestamp adapters, and the public codec-independent `restore` command.
- Added a manifest-governed, byte-stable fixture corpus with explicit provenance and redistribution decisions, reusable path-free golden comparisons, accepted and malformed cross-package conformance cases, and bounded fuzz boundaries for Cue JSON and source-envelope validation.
- Added stable least-privilege CI and independent CodeQL workflows with pinned actions and analysis tools, native Windows/macOS/Linux tests, race and conformance gates, vulnerability scanning, and six-target pure-Go build proof.
- Added a standalone pull-request policy engine and trusted-default-branch workflow for GitHub-resolved issue links, current-head Codex review state, operator-only exceptions, scheduled recovery, and one idempotent second-round request.
- Added a versioned repository-control contract covering secure Actions defaults, dependency and secret protections, private vulnerability reporting, squash-only delivery, automatic merged-branch cleanup, and evidence-backed default-branch rules.
- Added a snapshot-only six-target release proof with canonical schema and legal-file packaging, archive checksums, target-bound SPDX SBOMs, strict repository-owned artifact verification, and read-only hosted evidence retention.
- Added a dependency-free offline documentation verifier with a canonical maintained-document inventory, exact-case repository-confined path resolution, local-anchor validation, deterministic diagnostics, and CI coverage.
- Added dedicated SubRip, WebVTT, and release-process documentation that separates the current envelope-only foundation from planned native format support and protected publication actions.
- Added the complete official Cueson brand-kit 1.0.0 archive and byte-exact extraction, an independent 265-file import manifest, an offline integrity verifier, and repository brand guidance.
- Added the versioned v0.0.0 schema candidate, concise release notes, schema-bound release evidence, and non-publishing proof for the eventual squash-merge commit on `main`.

### Changed

- Replaced unavailable dynamic CI and release badges with truthful planned and unreleased state badges until those public resources exist.
- Changed the v0.0.0 envelope capability to advertise generic exact restoration while native ingest, model-driven render, and OCR-required capabilities remain false.
- Replaced the planned CI badge with the live workflow badge and refined the draft CI roadmap so codec, conversion, review-automation, and release gates activate only when their owning implementations exist.
- Raised the minimum Go version from 1.24 to 1.25 and upgraded `golang.org/x/text` to the first compatible fixed release after the vulnerability gate found reachable issues that could not be corrected on the old floor.
- Added the pull-request policy module to repository formatting and test gates while keeping it separate from the shipped Cueson module.
- Reconciled maintained documentation and foundation delivery records with merged evidence through S009 while keeping the v0.0.0 candidate under `[Unreleased]`.
- Adopted official local Cueson identity, palette, typography, favicon, slogan, and subordinate ShruggieTech endorsement across the README and media-format guide without adding brand assets to product release artifacts.

### Fixed

- Made fixture-alias and closed-descriptor tests portable across case-insensitive macOS filesystems and Linux timestamp syscall behavior exposed by the first hosted native matrix.
- Allowed a resolved second-round Codex finding to complete on a proven descendant remediation head with green required checks, removing the contradiction between stale-head rejection and the no-third-review rule.

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
- 2026-09-09: Reconcile native Codex round one from trusted default-branch code, use current-head GitHub evidence as the idempotency authority, and permit exactly one marked second-round request after resolved findings and green remediation-head CI.
- 2026-09-09: Layer a Cueson-owned default-branch ruleset over the unchanged organization baseline, require only successful current-head checks with verified providers, and retain one organization-administrator recovery path without granting normal agent merge authority.
- 2026-09-09: Pin GoReleaser and Syft as build-only tools, keep the checked-in release configuration unable to publish, and verify stable SBOM meaning and source binding without claiming byte-identical upstream SBOM output.
- 2026-09-09: Keep the v0.0.0 release candidate under `[Unreleased]` until a separately authorized tag and GitHub Release exist; a completed foundation milestone is not publication.
- 2026-09-09: Keep executable non-publishing candidate verification separate from the protected release process for tags, GitHub Releases, immutable schema copies, and production publication.
- 2026-09-10: Treat the operator-designated ShruggieTech 1.0.0 ZIP as the sole S011 acquisition authority, retain its archive and every safe payload file byte for byte, and preserve bundled terms without a redundant separate legal-review gate.
- 2026-09-10: Treat pull-request release proof as review evidence and bind the proposed publication target only after the S012 squash-merge commit passes its own non-publishing proof on `main`.

[Unreleased]: https://github.com/shruggietech/cueson/compare/v0.0.0...HEAD
[0.0.0]: https://github.com/shruggietech/cueson/releases/tag/v0.0.0
