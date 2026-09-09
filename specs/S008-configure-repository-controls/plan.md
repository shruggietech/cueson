# Implementation Plan: Verified Repository Controls

**Branch**: `S008-configure-repository-controls` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/S008-configure-repository-controls/spec.md`

## Summary

Implement issue #10 as a staged repository-configuration slice. First publish the source-controlled specification, desired-state contracts, and durable operator documentation through the official S008 pull request. Use that pull request to read back successful current-head CI, CodeQL, and trusted pull-request-policy evidence. Then mutate only `shruggietech/cueson`: restrict workflow defaults and action sources, enable supported security facilities, make squash the sole merge method, enable automatic merged-branch deletion, and add a repository-owned default-branch ruleset containing only proven checks. Read back every mutation immediately and record unavailable or deferred controls without touching the existing organization-owned ruleset.

## Technical Context

**Language/Version**: Markdown and YAML for source-controlled contracts; GitHub REST API version `2022-11-28` through GitHub CLI 2.70.0 for live state

**Primary Dependencies**: GitHub repository, Actions, code-security, ruleset, commit-status, check-run, issue, pull-request, and Project APIs; existing CI, CodeQL, and pull-request-policy workflows

**Storage**: GitHub repository settings and rulesets; versioned evidence and operating instructions in repository Markdown

**Testing**: GitHub API read-back, current-head check and status inspection, effective-rules evaluation, actionlint through existing CI, repository publication formatting, root and nested Go tests for regression safety

**Target Platform**: Public GitHub repository `shruggietech/cueson` and its `main` default branch

**Project Type**: Repository governance and delivery-system configuration

**Performance Goals**: Every mutation receives an immediate read-back in the same foreground session; no background polling process is introduced

**Constraints**: Repository-scoped mutations only; trusted current-head evidence before protection; one operator recovery bypass category; no merge, auto-merge, tag, release, schema publication, production-domain change, or organization-ruleset mutation

**Scale/Scope**: One repository, one default branch, one repository-owned ruleset, 17 stable S006 check contexts, and 2 conditional S007 policy contexts

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Lossless Source Preservation**: Pass. S008 does not read, rewrite, or publish subtitle or Cue JSON source data.
- **II. Schema and Official Software Discipline**: Pass. No schema, version, tag, or release identity changes are in scope.
- **III. Common Model Plus Native Fidelity**: Pass. Product model and codec boundaries remain untouched.
- **IV. No Silent Loss**: Pass. Every unsupported setting and partial mutation is recorded explicitly rather than normalized into success.
- **V. Test-First Format Work**: Pass. S008 contains no format implementation; existing regression gates remain mandatory.
- **VI. Portable and Secure Operation**: Pass. The slice reduces delivery privilege, restricts action sources, and enables supported security controls without changing shipped runtime behavior.
- **VII. Documentation and Delivery Authority**: Pass. Spec Kit, GitHub issue #10, the delivery Project, the official pull request, hosted verification, and human final merge authority remain explicit.
- **Contract Boundaries**: Pass. No public CLI, schema, or Go package contract changes.
- **Development and Delivery**: Pass. Repository text remains UTF-8 without BOM, Markdown remains unwrapped, verification stays foreground, and the operator has authorized this slice's push and pull-request publication only.

## Project Structure

### Documentation (this feature)

```text
specs/S008-configure-repository-controls/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
└── tasks.md
```

### Repository touch points

```text
.github/
└── workflows/                  # Existing least-privilege CI, CodeQL, and PR-policy definitions

docs/
├── architecture.md             # Effective protection and workflow-security architecture
├── project-management.md       # Required-check activation and recovery process
└── repository-controls.md      # Desired state, mutation order, limitations, and read-back evidence

CHANGELOG.md                    # Repository-control addition and architecture decision
SECURITY.md                     # Truthful private-reporting availability

specs/S008-configure-repository-controls/
├── checklists/requirements.md
├── contracts/
├── data-model.md
├── plan.md
├── quickstart.md
├── research.md
├── spec.md
└── tasks.md
```

**Structure Decision**: Keep live GitHub configuration outside shipped product code and record its desired state and evidence in documentation plus S008 design artifacts. Direct, foreground `gh api` calls perform the one-time mutations after hosted evidence exists. A new general-purpose administration program would add maintenance surface without improving this single-repository outcome.

## Post-Design Constitution Check

Pass. The design keeps mutations repository-scoped, does not execute pull-request-controlled code with administrative credentials, defers unproven policy contexts, preserves an organization-admin recovery path, records external-state evidence, and leaves final merge and release authority with the operator. No constitutional exception or complexity waiver is required.
