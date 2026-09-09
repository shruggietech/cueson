# Implementation Plan: Pull-Request Policy Automation

**Branch**: `S007-enforce-pr-review-automation` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/S007-enforce-pr-review-automation/spec.md`

## Summary

Add a trusted-default-branch GitHub Actions reconciler and a separately testable Go policy engine. The engine consumes GitHub's resolved closing references, interprets native Codex evidence for the current pull-request head, enforces exact operator exceptions and Dependabot exclusion, publishes stable commit-status contexts, and allows one idempotent marked second-round request after first-round findings are resolved. The workflow never executes pull-request-controlled content with write authority.

## Technical Context

**Language/Version**: Go 1.25 for a standalone repository-automation module; GitHub Actions workflow YAML

**Primary Dependencies**: Go standard library; GitHub REST and GraphQL APIs; pinned `actions/checkout` and `actions/setup-go`

**Storage**: No durable local storage; current GitHub pull-request evidence and stable marked comments are the state authority

**Testing**: Go table tests, in-memory HTTP mutation/read-back tests, actionlint, existing repository CI, and a read-only live pull-request audit

**Target Platform**: GitHub-hosted Ubuntu runner and local maintainer environments supported by Go 1.25

**Project Type**: Repository automation with a standalone Go command

**Performance Goals**: Evaluate one pull request in a bounded paginated request set; reconcile at most 100 open pull requests in one recovery sweep

**Constraints**: Fail closed; no pull-request code execution with write authority; no secrets; exact two-round ceiling; deterministic current-head attribution; all mutations read back; UTF-8 GitHub publication integrity

**Scale/Scope**: One public repository, ordinary pull requests to `main`, one configured Codex App identity, one configured operator login, two stable status contexts, and issue #9's verification corpus

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Lossless Source Preservation**: Not applicable to subtitle assets; automation never reads or rewrites product data.
- **II. Schema and Official Software Discipline**: No schema or release-version change.
- **III. Common Model Plus Native Fidelity**: No product model change.
- **IV. No Silent Loss**: Pass. Unknown evidence, API errors, incomplete pagination, and ambiguous attribution fail closed.
- **V. Test-First Format Work**: No format work; automation behavior is test-first with representative event and API fixtures.
- **VI. Portable and Secure Operation**: Pass. The command is pure Go, treats GitHub content as untrusted, and the write-capable workflow runs only default-branch code.
- **VII. Documentation and Delivery Authority**: Pass. S007 uses Spec Kit, blocks on analysis, follows issue #9, updates architecture and changelog, and stops at human merge.
- **Development and Delivery**: Pass. Text follows `.gitattributes`, Markdown uses logical lines, verification runs in the foreground, and no automated path can exceed two Codex rounds.

The pre-research gate passes without exception.

## Project Structure

### Documentation (this feature)

```text
specs/S007-enforce-pr-review-automation/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── closing-reference.md
│   ├── review-state-machine.md
│   └── workflow-security.md
└── tasks.md
```

### Source Code (repository root)

```text
.github/workflows/
└── pr-policy.yml

scripts/pr-policy/
├── go.mod
├── main.go
├── github.go
├── policy.go
├── github_test.go
└── policy_test.go

docs/
├── architecture.md
└── project-management.md
```

**Structure Decision**: Keep automation isolated from the shipped product module under `scripts/pr-policy`. Put deterministic decisions in `policy.go`, GitHub transport and pagination in `github.go`, and retain `main.go` as a thin command adapter. One globally serialized trusted-default-branch workflow invokes the command without checking out pull-request revisions.

## Post-Design Constitution Check

The design preserves the initial gate. The workflow contract excludes merge-ref review events from its write-capable trigger set, never checks out untrusted refs, uses a global non-canceling concurrency group, fails closed, and treats the single marked request as an atomic durable round-two reservation. No hidden state or extra review loop exists.

## Complexity Tracking

No constitutional violations require justification.
