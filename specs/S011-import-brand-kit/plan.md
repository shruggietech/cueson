# Implementation Plan: Import the Official Cueson Brand Kit

**Branch**: `S011-import-brand-kit` | **Date**: 2026-09-10 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/S011-import-brand-kit/spec.md`

## Summary

Pin the operator-designated ShruggieTech archive as the S011 acquisition, retain both the original ZIP and its complete safe extraction, record an independent 265-file inventory, add an offline read-only integrity verifier, adopt official assets directly from the retained kit in the README and media-format guide, document the brand boundary, and close issue #23 through a fully green and reviewed pull request.

## Technical Context

**Language/Version**: Go 1.25 for standalone verification; HTML, CSS, Markdown, JSON, and repository configuration for integration

**Primary Dependencies**: Go standard library plus `golang.org/x/text v0.39.0` for repository-equivalent Unicode normalization and case folding; Git and GitHub Actions for repository enforcement

**Storage**: Versioned ZIP, exact extracted files, and one strict JSON import manifest in the Git repository

**Testing**: Go unit and integration tests, repository text verification, documentation link verification, archive and extraction integrity verification, `git diff --check`, existing product tests, hosted CI, CodeQL, and bounded Codex review

**Target Platform**: Deterministic verification on Windows, Linux, and macOS; GitHub rendering for README; modern browsers and print for the standalone media guide

**Project Type**: Go CLI repository with separate repository-maintenance tools and static documentation

**Performance Goals**: Verify the 3,889,985-byte archive, 265 extracted files, and repository references in under 10 seconds on a typical development host without network access

**Constraints**: Preserve all retained and referenced brand bytes exactly; reject unsafe or ambiguous ZIP entries; never fetch during routine checks; do not add kit content to product releases; preserve CRLF policy for PowerShell; do not merge or perform release or production actions

**Scale/Scope**: One 1.0.0 archive, 265 payload files totaling 5,469,058 extracted bytes, six independently bound documentation asset references, one verifier module, and one GitHub issue

## Constitution Check

*GATE: Passed before Phase 0 research and rechecked after Phase 1 design.*

- **Spec-driven work**: PASS. S011 uses the installed Spec Kit phases, records the operator override, and treats analysis as blocking.
- **Source fidelity**: PASS. The archive, extracted payload, and referenced assets remain byte exact; verification rejects normalization, missing content, additional content, and unsafe names.
- **Schema and text style**: PASS. Cueson JSON is unchanged. Repository-authored JSON uses lowercase snake_case. Authored text remains UTF-8 without BOM and Markdown remains unwrapped.
- **GitHub publication integrity**: PASS. Issue and pull-request bodies will be formatted from files and read back after publication.
- **Work-slice governance**: PASS. Issue #23 remains one independently closeable outcome, is simplified according to the operator decision, and uses Slice `S011` with governed Stage transitions.
- **Review boundary**: PASS. Autopilot permits publication and at most two Codex rounds, but merge remains reserved for the operator.
- **Release and production boundary**: PASS. The current exact four-member release contract remains unchanged, and no tag, release, schema publication, domain activation, or production mutation is included.
- **Deviation from prior issue prose**: PASS. The operator explicitly superseded the dual-build reproduction, corrected-handoff, and separate legal-review gates. S011 pins the current live hosted ZIP as the sole acquisition authority, retains all safe content instead of only an old handoff subset, and preserves bundled terms without a redundant ownership checkpoint.

Post-design recheck: PASS. Direct references avoid unnecessary duplicate assets, the verifier remains a standalone non-product module, and protected paths are excluded from repair tooling without weakening checks for repository-authored text. Review remediation intentionally replaced the original standard-library-only preference with the root module's existing `golang.org/x/text v0.39.0` dependency because complete NFD normalization and Unicode case folding are required to enforce the repository's portable path-identity contract.

## Project Structure

### Documentation (this feature)

```text
specs/S011-import-brand-kit/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── import-manifest.schema.json
└── tasks.md
```

### Source Code and retained content

```text
brand/cueson/1.0.0/
├── archive/
│   └── cueson-brand-1.0.0.zip
├── kit/
│   └── <265 exact archive payload files>
└── import-manifest.json

scripts/brand-verify/
├── go.mod
├── go.sum
├── main.go
├── verify.go
└── verify_test.go

docs/
├── brand.md
└── cueson-media-format-guide.html

.github/workflows/ci.yml
.gitattributes
CHANGELOG.md
CONTRIBUTING.md
NOTICE
README.md
docs/architecture.md
docs/Cueson-Project-Specification-v0.0.0.md
scripts/docs-verify/
scripts/github-format/
```

**Structure Decision**: Keep the official acquisition isolated beneath a versioned `brand/` root, distinguish the immutable archive from its exact `kit/` extraction, store Cueson-owned acquisition evidence beside them, and keep verification in a narrowly dependent nested module outside the shipped product. The verifier uses the same pinned Unicode package version as the root module to enforce the same portable path identity without coupling the maintenance tool into shipped product code. Documentation references retained kit assets directly so the repository has one authoritative copy of each asset.

## Complexity Tracking

No constitution violations require justification.
