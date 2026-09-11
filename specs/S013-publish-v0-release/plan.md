# Implementation Plan: Publish and Verify v0.0.0

**Branch**: `S013-publish-v0-release` | **Date**: 2026-09-10 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/S013-publish-v0-release/spec.md`

## Summary

Publish an unsigned annotated `v0.0.0` tag at the exact accepted post-S012 revision, create the final GitHub Release from the already reviewed notes and thirteen exact files retained by the accepted default-branch release proof, independently download and compare every public byte with that evidence, and then update repository release-facing documentation through an official reviewed pull request. The transaction excludes milestone closure, signatures, attestations, public schema hosting, production changes, native codec claims, and pull-request merge.

## Technical Context

**Language/Version**: Go 1.25 for existing verification and formatting tools; Git, GitHub CLI, PowerShell, Markdown, and JSON for the release transaction and permanent record

**Primary Dependencies**: Git, GitHub Releases, GitHub Actions artifact `10178231596` from run `34543376814`, Go standard-library repository tools, and the existing six-target GoReleaser/Syft evidence

**Storage**: One immutable annotated Git tag, one public GitHub Release with thirteen assets, temporary local download directories, repository documentation, and Spec Kit artifacts

**Testing**: Exact digest comparison, checksum bijection, existing release verifier evidence, host-compatible packaged-binary execution, release/tag/API read-back, documentation verification, repository formatter, product tests, race detection, vet, vulnerability analysis, actionlint, hosted CI, CodeQL, and bounded Codex review

**Target Platform**: GitHub plus six published pure-Go targets across Windows, macOS, and Linux on amd64 and arm64; current Windows host for compatible public-binary execution

**Project Type**: Operational release publication with static repository documentation; no shipped product-code change

**Performance Goals**: Complete release read-back and thirteen-file digest audit in under five minutes after GitHub makes assets available; keep the full local gate within the existing release-proof budget

**Constraints**: Exact tag target, exact reviewed notes, exact accepted asset bytes, no rebuild or substitution, immutable versioned schema, UTF-8 without BOM, governed line endings, one source line per Markdown paragraph or list item, no milestone closure or production mutation, no PR merge

**Scale/Scope**: One tag, one release, six archives, six SBOMs, one checksum manifest, one GitHub issue and Project item, one official pull request, and release-facing repository documentation

## Constitution Check

*GATE: Passed before Phase 0 research and rechecked after Phase 1 design.*

- **Lossless source preservation**: PASS. Product source, source envelopes, and restoration behavior are unchanged. Published bytes are copied exactly from accepted evidence and compared by SHA-256 after public download.
- **Schema/software discipline**: PASS. The tag fixes the already reviewed versioned schema, and the accepted evidence binds repository, embedded, packaged, and binary versions to `0.0.0` with one schema digest.
- **Common model and no silent loss**: PASS. No model or codec behavior changes. Public notes and repository documentation retain the `envelope_only` boundary and explicitly deny native SRT/WebVTT capability.
- **Test-first format work**: PASS. No format implementation occurs. The release contract defines expected GitHub state, asset names, and digests before publication.
- **Portable and secure operation**: PASS. The accepted six-target proof is preserved byte for byte, the public set excludes internal evidence and metadata, and only the compatible Windows binary executes locally.
- **Documentation and delivery authority**: PASS. Spec Kit artifacts, a blocking analysis gate, issue #27, its single Project item, exact GitHub read-back, and the official pull request form the delivery record.
- **Protected actions**: PASS. The operator explicitly authorized the exact tag and release publication plus S013 push and pull-request creation. Milestone closure, production/schema-hosting changes, signatures, attestations, PR merge, and auto-merge remain excluded.
- **Partial-publication safety**: PASS. Preflight proves absence of an existing tag and release. Any mismatch or partial upload blocks follow-on documentation and requires operator judgment rather than replacement, movement, or silent retry.

Post-design recheck: PASS. The design adds no publication workflow or product dependency. It uses a one-time, explicitly authorized transaction, preserves accepted asset bytes, records expected state in a reviewable contract, and independently compares the public download with the accepted evidence before claiming success.

## Project Structure

### Documentation (this feature)

```text
specs/S013-publish-v0-release/
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

**Structure Decision**: Keep S013 operational. Reuse the accepted S012 artifact and verifier evidence, publish only the thirteen user-facing files, and prove public byte identity by exact digest comparison rather than adding a second release builder or modifying shipped code. Update the existing release-facing documents after verification so the repository distinguishes the public GitHub release from deferred production schema hosting and native codec work.

## Complexity Tracking

No constitution violations require justification.
