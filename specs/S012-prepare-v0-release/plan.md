# Implementation Plan: Prepare the v0.0.0 Release

**Branch**: `S012-prepare-v0-release` | **Date**: 2026-09-10 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/S012-prepare-v0-release/spec.md`

## Summary

Freeze the v0.0.0 schema into its versioned release path, stage a dated changelog and concise release notes, make the versioned schema the artifact packaging authority, extend the standalone verifier and evidence contract to prove immutable-schema identity, and run the same non-publishing proof on pull requests and pushes to `main`. The pull request proves the reviewed head; after squash merge, the push workflow binds the resulting default-branch commit as the proposed publication target. S012 stops before merge, tag, release, milestone, schema-hosting, or production actions.

## Technical Context

**Language/Version**: Go 1.25 for standalone verification; YAML, Markdown, JSON Schema Draft 2020-12, and repository configuration

**Primary Dependencies**: Go standard library, GoReleaser v2.18.1, Syft v1.51.1, GitHub Actions, and the existing repository modules

**Storage**: Versioned schema and release notes in Git; transient candidate artifacts and deterministic evidence under ignored `dist/`; GitHub Actions evidence retained for three days

**Testing**: Go unit and integration tests, byte-identity tests, policy tests, repository formatter, documentation verifier, product tests, race detection, vet, vulnerability analysis, actionlint, GoReleaser configuration validation, full six-target snapshot verification, hosted CI, CodeQL, and bounded Codex review

**Target Platform**: Six pure-Go targets across Windows, macOS, and Linux on amd64 and arm64; GitHub Actions on Ubuntu for complete structural proof; compatible native host execution when requested

**Project Type**: Go CLI repository with separate standard-library release-verification tooling and static release documentation

**Performance Goals**: Preserve the existing release-proof timeout of 20 minutes and verify repository-only schema and notes prerequisites in under one second before artifact inspection

**Constraints**: Exact schema bytes and version lockstep; exact six-target and four-member archive contracts; no product dependency added; no publication-capable configuration; no stale commit recorded before squash merge; UTF-8 without BOM; LF except governed CRLF scripts; no merge, tag, release, milestone closure, schema hosting, or production mutation

**Scale/Scope**: One release version, two repository schema copies, six archives, six binaries, six SBOMs, one six-entry checksum manifest, one deterministic evidence file, one release-note document, one changelog transition, one workflow, and one GitHub issue

## Constitution Check

*GATE: Passed before Phase 0 research and rechecked after Phase 1 design.*

- **Lossless source preservation**: PASS. Product source envelopes and restoration behavior are unchanged; release verification continues to reject local paths and machine identifiers in semantic metadata.
- **Schema/software discipline**: PASS. The release schema, canonical embedded schema, packaged schema, binary marker, and requested version must match exactly before evidence exists.
- **Common model and no silent loss**: PASS. S012 does not alter the common model or codec behavior, and release notes explicitly retain the envelope-only limitation.
- **Test-first work**: PASS. Schema-drift, notes-policy, evidence-shape, workflow-authority, and artifact-contract tests precede or accompany implementation.
- **Portable and secure operation**: PASS. The six CGO-disabled targets and bounded archive inspection remain intact; the workflow retains read-only authority and no secrets.
- **Documentation and delivery authority**: PASS. Spec Kit artifacts, a blocking analysis gate, one issue, one Project item, and the official pull request remain the delivery record.
- **Protected actions**: PASS. Push and pull-request creation are specifically authorized. Merge, tag, GitHub Release publication, milestone closure, schema hosting, and production mutation remain excluded.
- **Exact candidate binding**: PASS. S012 deliberately does not write the pull-request head into durable release prose. The push-to-main proof binds the eventual squash-merge revision, avoiding a stale or nonexistent final commit.

Post-design recheck: PASS. The evidence contract adds immutable-schema identity without widening the product API, and `release.disable: true`, no tag trigger, read-only permissions, no checkout credentials, and semantic SBOM verification preserve the established non-publishing boundary.

## Project Structure

### Documentation (this feature)

```text
specs/S012-prepare-v0-release/
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

### Source Code and release records

```text
schema/releases/v0.0.0/
└── cueson.schema.json

docs/releases/
└── v0.0.0.md

scripts/release-verify/
├── main.go
├── verify.go
├── verify_test.go
└── policy_test.go

.github/workflows/release-proof.yml
.goreleaser.yaml
CHANGELOG.md
README.md
docs/release-process.md
docs/release-verification.md
scripts/docs-verify/
```

**Structure Decision**: Keep the canonical authoring schema under `internal/schema`, admit an exact versioned release copy under `schema/releases`, and package from the versioned path so repository, archive, and evidence identities converge on the release record. Keep release notes under `docs/releases` so offline documentation verification audits them. Extend the existing standalone verifier rather than adding release behavior to the shipped executable. Trigger the same read-only proof on `main` pushes to establish post-squash evidence without introducing a publishing workflow.

## Complexity Tracking

No constitution violations require justification.
