# Quickstart: Validate Foundation Contract Ratification

## Prerequisites

- Work from branch `001-ratify-foundation-contracts`.
- Use the repository's configured Go, Git, GitHub CLI, PowerShell 7, and Spec Kit installations.
- Do not push, open a pull request, tag, release, or change the production domain during this validation.

## Local validation

Run from the repository root:

```powershell
go -C scripts/github-format test ./...
go run ./scripts/github-format/main.go .
git diff --check
specify check
```

Expected outcomes:

- Formatter tests pass.
- The whole-repository text check reports `github-format: OK`.
- Git reports no whitespace errors.
- Spec Kit reports the Codex integration ready.

## Contract validation

Confirm the ratified documents provide exactly one answer for:

- v0.0.0 schema URI and lifecycle;
- `subrip` and `webvtt` keys;
- completed-milestone `format_support` semantics;
- required OCR observation array cardinality;
- pre-1.0 minor and patch compatibility;
- command visibility and unavailable-command behavior;
- restoration independence from codecs;
- portable basename rejection of Windows devices and alternate-data-stream syntax;
- generic exact restoration in the v0.0.0 roadmap rather than the later codec series.

Trace issues #4 through #7 to the corresponding architecture, schema, or CLI section. No acceptance criterion may depend only on chat history.

## README validation

Confirm the CI and release badges are static state badges linked to issue #8 and milestone #1. Confirm license and documentation links resolve within the repository. No image target may depend on a missing workflow or release.

## GitHub readback

Read issues and Project state with GitHub CLI after reconciliation:

```powershell
gh issue view 1 --repo shruggietech/cueson --json state,body,url
gh issue view 3 --repo shruggietech/cueson --json state,body,url
gh issue view 6 --repo shruggietech/cueson --json state,body,url
gh project item-list 3 --owner shruggietech --format json
gh pr list --repo shruggietech/cueson --state open
```

Expected outcomes:

- Issue #1 remains open and In progress with only its merge-dependent README criterion pending.
- Issue #3 remains open, belongs to slice `001-ratify-foundation-contracts`, and Stage is In progress.
- Issue #6 assigns the generic exact-restoration CLI boundary to its source-foundation scope.
- No public pull request exists for this slice.

## Protected halt

After all tasks and checks pass, create a local conventional commit and stop. The anticipated command after separate authorization is:

```powershell
git push -u origin 001-ratify-foundation-contracts
```
