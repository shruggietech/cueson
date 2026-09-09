# Implementation Plan: Ratify Foundation Contracts

**Branch**: `001-ratify-foundation-contracts` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/001-ratify-foundation-contracts/spec.md`

## Summary

Convert the provisional v0.0.0 architecture, schema, and CLI choices into three concise ratified documents, correct README badges that target unavailable dynamic resources, and reconcile the completed bootstrap and active contract-ratification issues. The work is documentation and delivery metadata only: it adds no shipped Go code, root module, executable schema, CI workflow, repository rule, release, or production-domain mutation.

## Technical Context

**Language/Version**: GitHub Flavored Markdown, UTF-8 without BOM; JSON only for existing Spec Kit state

**Primary Dependencies**: GitHub Spec Kit 1.0.4, the existing standard-library `scripts/github-format` utility, Git, and authenticated GitHub CLI access

**Storage**: Versioned repository documents and GitHub issue/Project metadata

**Testing**: Existing Go unit tests for `scripts/github-format`, whole-repository publication-format checks, Git whitespace checks, strict UTF-8 checks, documentation link/contract assertions, and GitHub API readback

**Target Platform**: GitHub repository rendering plus local Windows development; the ratified contracts cover future Windows, macOS, and Linux product delivery

**Project Type**: Documentation and repository-governance slice for a future command-line product

**Performance Goals**: Every README badge target in this slice renders without depending on a missing workflow or release; all local validation completes in one foreground run

**Constraints**: No shipped product implementation; no push, public pull request, tag, release, ruleset, CI workflow, or production-domain change; one logical Markdown source line per paragraph and list item; source files remain UTF-8 without BOM with line endings governed by `.gitattributes`

**Scale/Scope**: Three ratified contract documents, the working project specification, README, changelog, one Spec Kit feature directory, and GitHub issues #1, #3, and #6

## Constitution Check

*GATE: Passed before Phase 0 research and re-checked after Phase 1 design.*

- **I. Lossless Source Preservation**: PASS. The schema contract keeps source bytes authoritative, prohibits source paths, and separates restoration from rendering.
- **II. Schema and Official Software Discipline**: PASS. The plan fixes the versioned schema identifier, lockstep rule, pre-1.0 compatibility, and release immutability without publishing a schema.
- **III. Common Model Plus Native Fidelity**: PASS. The schema contract retains common cue fields alongside format-native data and independently provenanced OCR observations.
- **IV. No Silent Loss**: PASS. The CLI contract distinguishes unknown formats, missing codecs, and strict loss rejection.
- **V. Test-First Format Work**: PASS. No parser or renderer changes occur; contract assertions and documentation checks precede and verify edits.
- **VI. Portable and Secure Operation**: PASS. No runtime code is introduced; unsafe path and untrusted-source boundaries remain normative.
- **VII. Documentation and Delivery Authority**: PASS. Spec Kit, issue contracts, Project state, analysis, verification, local commit, and protected push/PR boundaries govern the slice.

Post-design re-check: PASS. Research and contract artifacts introduce no constitutional exception or additional complexity.

## Project Structure

### Documentation (this feature)

```text
specs/001-ratify-foundation-contracts/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── architecture-baseline.md
│   ├── cli-baseline.md
│   └── schema-baseline.md
├── checklists/
│   └── requirements.md
└── tasks.md
```

### Repository Documents

```text
README.md
CHANGELOG.md
docs/
├── architecture.md
├── cli.md
├── schema.md
├── Cueson-Project-Specification-v0.0.0.md
└── project-management.md
```

**Structure Decision**: Keep the Spec Kit artifacts as the traceable design record and publish concise topic documents under `docs/` as the implementation baseline. The large working specification remains product context and roadmap, but directly links to the ratified documents when implementation details are needed.

## Implementation Phases

### Phase 0: Resolve decisions

1. Reconcile schema identity, lifecycle, format-support state, OCR cardinality, and source-authority rules against the constitution.
2. Reconcile command visibility, failure classification, streams, overwrite behavior, and future package ownership.
3. Select badge and GitHub bookkeeping behavior that is truthful before CI and releases exist.
4. Record the chosen alternatives and rejected options in `research.md`.

### Phase 1: Design contracts

1. Model the contract baseline and delivery-state entities in `data-model.md`.
2. Define architecture, schema, and CLI interface contracts under `contracts/`.
3. Define end-to-end verification and protected-authority checks in `quickstart.md`.
4. Re-evaluate every constitutional gate.

### Phase 2: Implement documentation and tracking

1. Author the three ratified documents under `docs/` from the approved contracts.
2. Link the working specification to those documents and resolve touched provisional contradictions.
3. Replace unavailable dynamic README badges with truthful static state badges.
4. Record decisions in the changelog.
5. Reconcile bootstrap and active-slice GitHub metadata with readback.

### Phase 3: Analyze, verify, and halt

1. Run the Spec Kit artifact analysis gate and resolve every blocking finding.
2. Execute local format, test, encoding, link, and contract checks in the foreground.
3. Commit the verified slice locally.
4. Halt before push and before public pull-request creation.

## Complexity Tracking

No constitution violations or complexity exceptions are required.
