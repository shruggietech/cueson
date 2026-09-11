# Implementation Plan: Publish and Verify v1.0.0

**Branch**: `codex/S020-publish-v1-release` | **Date**: 2026-09-11 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/S020-publish-v1-release/spec.md`

## Summary

Freeze the exact post-S019 v1.0.0 candidate and its thirteen-file publication contract, complete non-mutating release preflight, and halt until the operator specifically authorizes the immutable tag and public GitHub Release. After authorization, publish and independently verify every public byte, update release-facing repository documentation, complete the official reviewed pull request, and reconcile the v1 delivery surfaces only within separately granted authority.

## Technical Context

**Language/Version**: Go 1.25 for existing release verification and GitHub-body formatting tools; Git, GitHub CLI, PowerShell, Markdown, and JSON for the release transaction and permanent record

**Primary Dependencies**: Git, GitHub Releases, accepted GitHub Actions artifact `10273380044` from release-proof run `34621429626`, Go standard-library repository tools, and existing six-target GoReleaser and Syft evidence

**Storage**: One prospective immutable annotated Git tag, one prospective public GitHub Release with thirteen assets, clean temporary verification directories, repository documentation, GitHub planning metadata, and Spec Kit artifacts

**Testing**: Exact digest comparison, checksum bijection, candidate evidence validation, archive and SBOM inspection, host-compatible packaged-binary execution, release/tag/API read-back, documentation verification, repository formatter, product tests, race detection, vet, vulnerability analysis, actionlint, hosted CI, CodeQL, release proof, and bounded Codex review

**Target Platform**: GitHub plus six published pure-Go targets across Windows, macOS, and Linux on amd64 and arm64; current Windows host for compatible public-binary execution

**Project Type**: Operational release publication with static repository documentation; no shipped product-code change is planned

**Performance Goals**: Complete release read-back and thirteen-file digest audit within five minutes after GitHub makes assets available; keep the full local verification gate within the established release-proof budget

**Constraints**: Exact tag target, exact approved notes, exact accepted asset bytes, no rebuild or substitution, immutable versioned schema, UTF-8 without BOM, governed line endings, one source line per Markdown paragraph or list item, no protected mutation without exact authorization, no pull-request merge

**Scale/Scope**: One tag, one release, six archives, six SBOMs, one checksum manifest, one GitHub issue and Project item, one parent epic and milestone, one official pull request, and release-facing repository documentation

## Constitution Check

*GATE: Passed before Phase 0 research and rechecked after Phase 1 design.*

- **Lossless source preservation**: PASS. Product source and restoration behavior are unchanged. Candidate and future public bytes are compared exactly by SHA-256.
- **Schema/software discipline**: PASS. The candidate binds repository, embedded, emitted, packaged, and binary versions to `1.0.0` with immutable schema digest `1aad14567033d7e14d9beb78985e18007aefb5345095370b11b6b887df7ec541`.
- **Common model and no silent loss**: PASS. No model or codec behavior changes. Release notes preserve the stable native-format and conversion-loss boundaries.
- **Test-first format work**: PASS. No format implementation occurs. The contract fixes expected GitHub state, asset names, sizes, and digests before publication.
- **Portable and secure operation**: PASS. One accepted six-target bundle is preserved byte for byte, public assets exclude internal evidence and metadata, and only the compatible Windows binary executes locally.
- **Documentation and delivery authority**: PASS. Spec Kit artifacts, a blocking analysis gate, issue #38, its single Project item, exact GitHub read-back, and an official pull request form the delivery record.
- **Protected actions**: PASS. Current authority covers kickoff, local specification, and read-only preflight only. Tag creation or push, release publication, branch push, pull-request publication, milestone or epic closure, production/schema-hosting changes, signatures, attestations, PR merge, and auto-merge remain blocked until specifically authorized.
- **Partial-publication safety**: PASS. Preflight proves absence of an existing tag and release. Any future mismatch or partial upload blocks follow-on documentation and requires operator judgment rather than movement, replacement, or silent retry.

Post-design recheck: PASS. The design adds no product dependency or alternate build path. It preserves the accepted S019 artifact, makes every protected transition explicit, records expected public state in a reviewable contract, and independently compares public downloads with accepted evidence before claiming success.

## Project Structure

### Documentation (this feature)

```text
specs/S020-publish-v1-release/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── release-publication-contract.json
├── checklists/
│   └── requirements.md
└── tasks.md
```

### Release-facing repository records

```text
README.md
CHANGELOG.md
docs/
├── architecture.md
├── cli.md
├── project-management.md
├── release-process.md
├── release-verification.md
└── schema.md
```

**Structure Decision**: Keep S020 operational. Reuse the accepted S019 artifact and verifier evidence, publish only the thirteen user-facing files after exact authorization, and prove public byte identity by digest comparison rather than adding another builder or changing shipped code. Update existing release-facing documents only after verified publication so the repository distinguishes the public release from deferred production schema hosting.

## Complexity Tracking

No constitution violations require justification.
