# Implementation Plan: v0.0.0 Documentation and Milestone Verification

**Branch**: `S010-complete-v0-docs` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/S010-complete-v0-docs/spec.md`

## Summary

Complete the required v0.0.0 repository documentation set, reconcile every current-versus-planned capability and delivery claim, and produce operator-ready milestone evidence for issues #12 and #2. Add a dependency-free standalone documentation verifier that checks the canonical document inventory plus local file and heading links without network access, integrate it into the existing Repository text gate, and leave all release, tag, immutable-schema, milestone-closure, and production actions outside S010.

## Technical Context

**Language/Version**: Repository Markdown and HTML; Go 1.25.0 for the standalone documentation verifier

**Primary Dependencies**: Go standard library only; existing repository `scripts/github-format`, CLI/schema tests, GitHub Actions, Issues, milestones, and Project metadata

**Storage**: Versioned repository files and read-only GitHub delivery metadata; no application storage

**Testing**: Table-driven standalone verifier tests, repository-level accepted and rejected link fixtures created in temporary directories, existing product and standalone-module suites, actionlint, race detection, vet, and hosted CI/CodeQL

**Target Platform**: Offline local verification on Windows, macOS, and Linux plus the existing Ubuntu Repository text job

**Project Type**: Documentation closeout plus a standalone repository-verification command outside the shipped product module

**Performance Goals**: Verify the complete audited documentation set and all local links in under five seconds on a normal checkout without network access

**Constraints**: UTF-8 without BOM, `.gitattributes` line endings, one Markdown source line per paragraph/list item, no hard-coded local machine values, no network-dependent link result, no shipped-product behavior change, no tag/release/schema-publication/milestone-close/production mutation

**Scale/Scope**: Twelve required contract documents, all maintained Markdown and HTML beneath `docs/`, additional root and governance documentation, the v0.0.0 foundation gate, eleven epic children, two remaining open issues, five stale issue evidence records, and one small standard-library verifier module

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Lossless Source Preservation**: PASS. Documentation describes source authority and restoration without modifying source-handling code or weakening path prohibitions.
- **II. Schema and Official Software Discipline**: PASS. Existing lockstep tests remain authoritative, and S010 does not create the immutable release schema before an authorized release.
- **III. Common Model Plus Native Fidelity**: PASS. Format pages distinguish the common model, native data, source envelope, and unimplemented codecs.
- **IV. No Silent Loss**: PASS. Current and planned diagnostic/loss boundaries remain explicit; no unavailable conversion is claimed.
- **V. Test-First Format Work**: PASS. No codec behavior is implemented. Verifier behavior receives table-driven tests before repository integration.
- **VI. Portable and Secure Operation**: PASS. The verifier is standard-library, offline, path-safe, and cross-platform; it performs no writes during validation.
- **VII. Documentation and Delivery Authority**: PASS. Documentation, GitHub evidence, Spec Kit artifacts, and full verification converge before publication; the operator retains merge and release authority.
- **Contract boundaries**: PASS. The work does not add a public Go API or shipped command and does not alter Cue JSON.
- **Development and delivery**: PASS. Text, changelog, review, push, and publication requirements remain enforced. Push and PR creation are covered by the operator's explicit S010 authorization; merge and release actions are not.

Post-design check: PASS. The documentation inventory, verifier contract, evidence model, and validation guide preserve the same boundaries with no justified exception.

## Project Structure

### Documentation (this feature)

```text
specs/S010-complete-v0-docs/
├── plan.md              # This file ($speckit-plan command output)
├── research.md          # Phase 0 output ($speckit-plan command)
├── data-model.md        # Phase 1 output ($speckit-plan command)
├── quickstart.md        # Phase 1 output ($speckit-plan command)
├── contracts/           # Phase 1 output ($speckit-plan command)
└── tasks.md             # Phase 2 output ($speckit-tasks command - NOT created by $speckit-plan)
```

### Source Code (repository root)
```text
README.md
CHANGELOG.md
.github/workflows/ci.yml
docs/
├── architecture.md
├── cli.md
├── formats/
│   ├── srt.md
│   └── webvtt.md
├── project-management.md
├── release-process.md
├── release-verification.md
├── repository-controls.md
├── schema.md
├── Cueson-Project-Specification-v0.0.0.md
└── cueson-media-format-guide.html
scripts/docs-verify/
├── go.mod
├── main.go
├── verify.go
└── verify_test.go
```

**Structure Decision**: Preserve existing release-specific documents as focused authorities, add the three missing canonical documents at the paths already required by the project specification, and place offline documentation validation in its own standard-library module so publication formatting, shipped product behavior, and documentation contract checks remain separate responsibilities.
