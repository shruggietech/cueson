# Implementation Plan: Prepare the v1.0.0 Release Candidate

**Branch**: `codex/S019-prepare-v1-release-candidate` | **Date**: 2026-09-11 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/S019-prepare-v1-release-candidate/spec.md`

## Summary

Promote the frozen v1-bound executable and Cue JSON contract from development identity `0.1.0` and `experimental` format support to coordinated candidate identity `1.0.0` and `stable` SubRip/WebVTT support. Admit a byte-identical immutable v1 schema, finalize the detailed changelog and concise candidate release notes, switch the existing snapshot-only six-target pipeline into non-development candidate proof, strengthen exact schema and legal-file verification, add native packaged-binary smoke proof on Linux, Windows, and macOS, reconcile maintained documentation and governed fixtures, and prove the complete artifact set from a clean exact commit. S019 publishes only its branch and official pull request; merge, tag, GitHub Release, milestone closure, schema hosting, and production changes remain excluded.

## Technical Context

**Language/Version**: Go 1.25 for the product and standalone verification tools; YAML, Markdown, and JSON Schema Draft 2020-12 for delivery contracts

**Primary Dependencies**: Go standard library, `golang.org/x/text`, GoReleaser v2.18.1, Syft v1.51.1, GitHub Actions, and the existing repository modules

**Storage**: Canonical and immutable versioned schemas plus release records in Git; governed Cue JSON fixtures in `testdata/`; transient candidate artifacts and deterministic evidence under ignored `dist/`; short-lived GitHub Actions evidence

**Testing**: Go unit, integration, conformance, fuzz-seed, schema-lockstep, fixture-integrity, CLI, documentation, policy, and release-verifier tests; formatter, race detection, vet, vulnerability analysis, actionlint, GoReleaser validation, six-target snapshot verification, Linux/Windows/macOS packaged-binary execution, hosted CI, CodeQL, security review, and bounded Codex review

**Target Platform**: Six pure-Go targets across Windows, macOS, and Linux on amd64 and arm64; hosted amd64 packaged-binary smoke checks on each operating system; structural and build-information proof for every target

**Project Type**: Go CLI repository with embedded JSON Schema, maintained static documentation, governed fixtures, and separate standard-library verification tools

**Performance Goals**: Keep each hosted candidate-proof job within 20 minutes, preserve existing bounded input limits, and reject repository identity or policy failures before expensive artifact inspection

**Constraints**: Exact `1.0.0` lockstep; byte-identical schema and legal-file copies; stable support only for the frozen SubRip/WebVTT v1 contracts; exact six-target, four-member archive contract; clean VCS metadata; no product dependency expansion; UTF-8 without BOM; `.gitattributes` line endings; no publication-capable configuration; no predicted post-squash commit; no merge, tag, release, milestone closure, public schema hosting, or production mutation

**Scale/Scope**: One coordinated version and maturity transition, one immutable schema copy, six archives, six binaries, six software bills of materials, one six-entry checksum manifest, one deterministic evidence document per proof runner, one release-note document, one changelog release section, governed fixture migration, maintained current-state documentation, one workflow, and one existing GitHub issue

## Constitution Check

*GATE: Passed before Phase 0 research and rechecked after Phase 1 design.*

- **Lossless source preservation**: PASS. The release transition does not rewrite source envelopes or weaken restoration. Updated fixtures retain exact governed payload identities and path-leak checks.
- **Schema and official software discipline**: PASS. Runtime, canonical schema, embedded schema, immutable release copy, packaged schema, evidence, and intended tag must all agree on `1.0.0`; the released v0.0.0 copy remains unchanged.
- **Common model plus native fidelity**: PASS. Stable status recognizes the already frozen SubRip and WebVTT model/native contracts without altering source authority or moving derived observations into source truth.
- **No silent loss**: PASS. Conversion loss accounting, strict refusal, diagnostics, and bounded rejection remain unchanged and must pass the existing conformance suite.
- **Test-first format work**: PASS. Lockstep, stable-support, fixture, exact-identity, legal-file, packaging, workflow, notes, and evidence assertions are updated before or alongside the coordinated implementation.
- **Portable and secure operation**: PASS. The six CGO-disabled targets, hosted native smoke checks, archive safety rules, clean build metadata, bounded verification, and hidden Windows child-process behavior remain required.
- **Documentation and delivery authority**: PASS. The S019 Spec Kit artifacts, blocking analysis gate, issue #37, v1 milestone, Delivery Project item, official pull request, and review record remain the authority.
- **Protected actions**: PASS. Branch push and official pull-request creation are specifically authorized. Merge, tag, GitHub Release and asset publication, milestone closure, public schema hosting, and production mutation remain excluded.
- **Intentional prior-behavior departure**: PASS. Replacing `0.1.0`/`experimental` development behavior and the `-development` verifier path with `1.0.0`/`stable` candidate behavior is the explicit outcome of issue #37 and the transition required by `docs/architecture.md`; historical v0.0.0 and completed Spec Kit artifacts are not rewritten.

Post-design recheck: PASS. The design uses existing internal ownership, keeps publication disabled, strengthens identity and archive verification at the standalone verifier boundary, preserves immutable historical evidence, and adds no public Go API or production dependency.

## Project Structure

### Documentation (this feature)

```text
specs/S019-prepare-v1-release-candidate/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── release-evidence.schema.json
├── checklists/
│   └── requirements.md
└── tasks.md
```

### Source Code, contracts, and release records

```text
internal/version/
internal/model/
internal/schema/
internal/cli/
internal/convert/
testdata/
schema/releases/v1.0.0/
docs/releases/v1.0.0.md
scripts/release-verify/
scripts/docs-verify/
.github/workflows/release-proof.yml
.goreleaser.yaml
CHANGELOG.md
README.md
CONTRIBUTING.md
SECURITY.md
docs/
```

**Structure Decision**: Keep version ownership under `internal/version`, schema identity and embedding under `internal/schema`, capability validation under `internal/model`, and produced stable capability values in the existing CLI and conversion owners. Finalize the canonical schema first, copy its exact bytes once into `schema/releases/v1.0.0`, and package from that immutable path. Extend the existing standalone verifier and documentation verifier rather than adding release behavior to the shipped executable. Build the candidate once in the hosted proof job, retain its exact bundle, and fan that bundle out to native Linux, Windows, and macOS smoke jobs instead of rebuilding per platform. Preserve the snapshot-only GoReleaser invocation and read-only workflow while running the verifier in release-candidate mode.

## Complexity Tracking

No constitution violations require justification.
