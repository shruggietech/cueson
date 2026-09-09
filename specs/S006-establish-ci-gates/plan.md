# Implementation Plan: Establish CI and Cross-Platform Build Gates

**Branch**: `S006-establish-ci-gates` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/S006-establish-ci-gates/spec.md`

## Summary

Implement issue #8 as the first hosted-delivery automation slice. Add a least-privilege `CI` workflow whose stable jobs cover repository formatting, repository-authored text, vet, root and nested module tests, native Windows/macOS/Linux behavior, race detection, schema and conformance tests, pinned static analysis, pinned vulnerability analysis, and six pure-Go cross-build targets. Add an independent `CodeQL` workflow with narrowly scoped result-upload authority. Prove fail-closed behavior through a temporary marker committed only to the draft pull-request history, then remove it and require the corrected run to be green. Update the README badge, architecture, changelog, Spec Kit evidence, issue, and Project state without configuring required checks, security settings, releases, or product codecs.

## Technical Context

**Language/Version**: Go 1.25.0 minimum; CI compatibility jobs use the latest Go 1.25 patch; security and workflow-lint tooling may use the current stable Go toolchain

**Primary Dependencies**: GitHub Actions; `actions/checkout` v7.0.1 at commit `3d3c42e5aac5ba805825da76410c181273ba90b1`; `actions/setup-go` v7.0.0 at commit `b7ad1dad31e06c5925ef5d2fc7ad053ef454303e`; `github/codeql-action` v4.38.0 at commit `b96794f015dfd88f77b49b1c93e0fa7110f94c63`; Staticcheck v0.7.0; govulncheck v1.8.0; actionlint v1.7.12

**Storage**: Repository YAML and Markdown files only; hosted build outputs remain temporary runner files

**Testing**: Existing root Go tests, race tests, schema/conformance tests, publication-formatter tests, actionlint, Staticcheck, govulncheck, CodeQL, controlled failing and corrected GitHub workflow runs

**Target Platform**: Native `ubuntu-24.04`, `windows-2025`, and `macos-15`; pure-Go builds for Windows, macOS, and Linux on amd64 and arm64

**Project Type**: Go CLI repository with GitHub-hosted delivery automation

**Performance Goals**: All S006 pull-request checks finish within 20 minutes under normal hosted-runner availability

**Constraints**: Explicit least privilege; immutable action commits; explicit analysis-tool versions; stable job names; no secrets or `pull_request_target`; no background verification; no artifact upload or release; no repository-setting mutation; no checks for unimplemented codec, conversion, review-automation, or release surfaces; UTF-8 without BOM; `.gitattributes` line endings; one source line per Markdown paragraph and list item

**Scale/Scope**: Two workflows, three native test operating systems, six cross-build targets, one controlled failure proof, issue #8 only

## Constitution Check

*GATE: Passed before Phase 0 research and re-checked after Phase 1 design.*

- **I. Lossless Source Preservation**: Pass. S006 does not modify source envelopes or fixture bytes; hosted tests exercise existing integrity contracts.
- **II. Schema and Official Software Discipline**: Pass. The schema/software lockstep tests receive an explicit stable gate; no release identity changes.
- **III. Common Model Plus Native Fidelity**: Pass. Existing conformance tests run without modifying the common model or format-native data.
- **IV. No Silent Loss**: Pass. Automation fails closed and does not replace unavailable feature tests with empty successes.
- **V. Test-First Format Work**: Pass. Native platform tests execute on all three owning operating systems; no codec work enters the slice.
- **VI. Portable and Secure Operation**: Pass. Cross-builds cover six `CGO_ENABLED=0` targets, ordinary jobs are read-only, untrusted pull requests receive no secrets, and action/tool executions are pinned.
- **VII. Documentation and Delivery Authority**: Pass. S006 uses Spec Kit, the blocking analysis gate, issue #8, Project state, foreground verification, an issue-closing pull request, and no more than two Codex rounds. The operator retains final merge authority.
- **Contract Boundaries**: Pass. The change affects delivery automation and documentation only, not the public CLI or schema contracts.
- **Development and Delivery**: Pass. Repository text and line endings are gated, checks remain foreground, and push/PR authority is supplied explicitly for S006.

Post-design re-check: Pass. The deliberate refinement to the working roadmap is that S006 omits cross-format, codec, Codex-review-gate, and release tests until their owning implementations exist. Adding placeholder jobs would contradict the constitution's truthful-delivery and no-silent-loss principles. Issue #9 retains review automation, issue #10 retains repository controls, issue #11 retains release proof, and later feature slices retain codec and conversion behavior.

## Project Structure

### Documentation (this feature)

```text
specs/S006-establish-ci-gates/
├── checklists/
│   └── requirements.md
├── contracts/
│   ├── check-contract.md
│   ├── failure-proof.md
│   └── permissions-and-events.md
├── data-model.md
├── plan.md
├── quickstart.md
├── research.md
├── spec.md
└── tasks.md
```

### Source Code (repository root)

```text
.github/
└── workflows/
    ├── ci.yml
    └── codeql.yml

cmd/
internal/
scripts/
└── github-format/

README.md
CHANGELOG.md
docs/architecture.md
docs/Cueson-Project-Specification-v0.0.0.md
```

**Structure Decision**: Keep delivery automation in the standard `.github/workflows` boundary and reuse existing repository verification commands. Refine the working specification's CI list where it conflicts with staged issue ownership. Do not add production packages, a workflow abstraction layer, uploaded artifacts, or a new checker when the standalone publication formatter already owns UTF-8, BOM, mojibake, Markdown, and line-ending verification.

## Complexity Tracking

No constitutional violation requires an exception.
