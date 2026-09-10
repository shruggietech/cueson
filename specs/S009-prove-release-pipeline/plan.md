# Implementation Plan: Non-Publishing Release Proof

**Branch**: `S009-prove-release-pipeline` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/S009-prove-release-pipeline/spec.md`

## Summary

Establish a complete v0.0.0 snapshot release proof using an exact GoReleaser v2 configuration, one target-bound SPDX software bill of materials generated from each packaged binary by an exact Syft version and named for its corresponding archive, a standard-library-only verification module, and a read-only pull-request workflow. The proof builds the six approved pure-Go targets, packages the canonical schema and legal files, validates checksums and archive structure, executes host-compatible version surfaces, rejects local path leakage, retains review artifacts temporarily, and has no tag, signing, attestation, release, release-asset, or production-domain publication path.

## Technical Context

**Language/Version**: Go 1.25.0 for the product and repository-owned verifier; YAML for GoReleaser and GitHub Actions configuration

**Primary Dependencies**: GoReleaser v2.18.1 and Syft v1.51.1 as exact build-only tools; Go standard library for `scripts/release-verify`; immutable GitHub-owned workflow actions only

**Storage**: Transient `dist/` artifacts locally and short-lived GitHub Actions artifacts in hosted verification; no application persistence

**Testing**: Table-driven Go tests for archive, checksum, metadata, SBOM, path-leak, and host-execution decisions; GoReleaser configuration validation; full snapshot generation; foreground artifact inspection; existing root and standalone-module gates

**Target Platform**: Windows, macOS, and Linux on amd64 and arm64 with `CGO_ENABLED=0`; snapshot orchestration and full-matrix structural verification run on Windows locally and Ubuntu in hosted automation

**Project Type**: Command-line product plus repository-owned release automation and standalone validation tooling

**Performance Goals**: The hosted snapshot proof completes within 20 minutes and emits one independently visible terminal check result

**Constraints**: No publication credential, tag mutation, GitHub Release mutation, signing, attestation publication, or `cueson.io` mutation; repository Actions allow only GitHub-owned actions; Windows child processes remain hidden and non-interactive; release output contains no local paths or machine identifiers; `dist/` remains ignored

**Scale/Scope**: Six archives, six SPDX JSON documents, one SHA-256 manifest, GoReleaser metadata, one release-evidence summary, and issue #11 only

## Constitution Check

*GATE: Passed before research and passed again after design.*

- **I. Lossless Source Preservation**: The packaged canonical schema is copied byte-for-byte and verified against the repository source. No Cue JSON source envelope is transformed.
- **II. Schema and Official Software Discipline**: The release version is injected through the existing `internal/version` build variable, and verification compares executable, embedded schema, packaged schema, and requested release identities.
- **III. Common Model Plus Native Fidelity**: S009 does not change Cue JSON model or codec behavior.
- **IV. No Silent Loss**: Missing archives, archive members, checksums, SBOMs, metadata, or validation evidence fail the snapshot proof visibly.
- **V. Test-First Format Work**: S009 changes no parser or renderer. Release verification receives focused positive and negative tests before hosted use.
- **VI. Portable and Secure Operation**: All six first-class targets build with `CGO_ENABLED=0`; archive traversal and local-path leakage are rejected; Windows verification subprocesses use `CREATE_NO_WINDOW`.
- **VII. Documentation and Delivery Authority**: Spec Kit artifacts, issue #11, Project Stage and Slice, changelog decisions, an issue-closing pull request, bounded Codex review, and human-only merge authority remain in force.
- **Authority boundary**: The operator authorized the S009 branch push and pull-request publication, but not merge, tag, release, signature, attestation, or production publication.

No constitutional exception is required.

## Project Structure

### Documentation (this feature)

```text
specs/S009-prove-release-pipeline/
├── checklists/
│   └── requirements.md
├── contracts/
│   ├── artifact-contract.md
│   └── workflow-authority.md
├── data-model.md
├── plan.md
├── quickstart.md
├── research.md
├── spec.md
└── tasks.md
```

### Source Code (repository root)

```text
.github/workflows/
└── release-proof.yml

scripts/release-verify/
├── go.mod
├── main.go
├── process_other.go
├── process_windows.go
├── verify.go
└── verify_test.go

.goreleaser.yaml
CHANGELOG.md
README.md
docs/architecture.md
docs/release-verification.md
```

**Structure Decision**: Keep GoReleaser at its conventional repository-root path, isolate the release verifier as a dependency-free standalone Go module so it does not enter the shipped product graph, place hosted behavior in one narrow workflow, and record public architecture and development guidance without creating a second product specification.

## Phase 0: Research

1. Confirm current GoReleaser v2 snapshot, build-matrix, archive, checksum, metadata, and SBOM behavior.
2. Confirm the current stable GoReleaser and Syft releases and select exact versions.
3. Reconcile third-party release tooling with the repository's GitHub-owned-actions-only policy.
4. Define reproducibility, path-hygiene, and non-publication controls that can be proven before a real tag exists.
5. Define host-compatible executable checks and full-matrix structural checks without pretending cross-target binaries can run locally.

## Phase 1: Design and Contracts

1. Define release-candidate, target, archive, checksum, SBOM, metadata, and evidence entities in [data-model.md](data-model.md).
2. Freeze archive names, formats, members, digest coverage, SBOM names, and version assertions in [contracts/artifact-contract.md](contracts/artifact-contract.md).
3. Freeze workflow events, permissions, tool acquisition, artifact retention, failure behavior, and prohibited publication surfaces in [contracts/workflow-authority.md](contracts/workflow-authority.md).
4. Define the foreground local and hosted validation sequence in [quickstart.md](quickstart.md).
5. Re-evaluate the constitution check against the design before task generation.

## Implementation Strategy

Build the verifier's negative cases and core archive model first, then add the GoReleaser configuration, generate a real snapshot, and refine validation against observed GoReleaser and Syft output. Only after local proof passes should the read-only hosted workflow and documentation be added. The analysis gate must report no critical or high-severity artifact inconsistency before implementation starts. Convergence runs after implementation and may only append remaining tasks.

## Complexity Tracking

No constitutional violations or complexity exceptions are present.
