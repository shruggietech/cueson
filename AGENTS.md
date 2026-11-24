# Cueson agent operating contract

This file is binding for AI agents working in this repository. The current project specification is a working pre-release draft. Explicit operator decisions and the ratified constitution control bootstrap sequencing when draft implementation details conflict.

## Start with current authority

Before implementation, inspect `.specify/memory/constitution.md`, the relevant Spec Kit slice, the architecture of record, and the active GitHub issues. Challenge contradictions instead of duplicating them. If an implementation departs from prior behavior or draft prose, state the departure and rationale explicitly.

## Spec-driven work

Non-trivial product work uses the installed GitHub Spec Kit workflow. When the `shruggie-speckit` skill is available, full-slice kickoff or autopilot work MUST use it to orchestrate the installed Spec Kit commands. Do not invent substitute commands when an integration is missing.

The analysis gate is blocking. Verification runs in the foreground. Independent v1 work SHOULD use sub-agents when file ownership and verification boundaries are clear; the coordinating agent owns integration and final verification.

## Source fidelity

Never drop unknown source content, normalize raw source fields, alter preserved source bytes, or label lossy conversion as lossless without an approved specification change. Derived speaker, token, text, or OCR fields never replace source truth. Cue JSON must never contain an original filesystem path or local machine identifier.

## Schema and text style

Cueson-owned JSON keys use lowercase `snake_case`. Repository-authored text uses UTF-8 without BOM, and line endings follow `.gitattributes`. Intentional non-UTF-8 or nonstandard-line-ending fixtures belong under `testdata/` with documented provenance.

Agents MUST NOT hard-wrap Markdown prose. Use one source line per paragraph and one source line per list item regardless of length. Fenced code, tables, generated files, and native format structures follow their own layout.

## GitHub publication integrity

Before publishing an issue, pull request, review, or release-note body, author it as UTF-8 Markdown with intentional blank lines and pass it through:

```text
go run ./scripts/github-format/main.go -stdin
```

Publish from a file whenever practical. Immediately read the body back from GitHub and verify headings, blank lines, lists, task checkboxes, fences, tables, and plausible line count. A successful API response is not proof of correct rendering.

## Work slices and issues

Every actionable issue owns one independently closeable and testable outcome with these sections:

```text
Outcome
Context
Scope
Acceptance criteria
Dependencies
Verification
```

Before recommending a work slice, inspect active issues, milestones, dependencies, and prior operator decisions. Prefer the largest coherent group that shares one implementation narrative, review surface, and verification story. One issue is not automatically one slice. Preserve each issue's acceptance criteria and explain relevant exclusions.

Use native assignees, milestones, parent/sub-issues, and dependencies for facts they own. Use governed labels for type, priority, effort, area, and gates. Do not duplicate native metadata in Project fields. Every in-scope repository issue appears exactly once in the `cueson Delivery` Project.

The Project `Stage` field uses exactly:

```text
Backlog
Ready
Specced
In progress
PR review
Release verification
Done
```

The Project `Slice` text field is the only default additional planning field. Leave the default `Status` field unused.

## Pull requests and Codex review

Normal pull requests to the default branch contain at least one complete closing reference such as `Closes #123`. Eligible non-Dependabot pull requests receive an automatic first Codex review. An eyeballs reaction is acknowledgement only; a thumbs-up with no findings is a clean review.

If round one finds issues, inspect and address every finding, resolve threads only after the concern is actually handled, rerun verification, update the head, and then request exactly one second review with `@codex review`. Do not request a third review automatically.

Before publishing that second-round request, include exactly one backticked full commit ID for the requested head and verify its repository-policy binding. Preserve the request receipt after publication: editing it changes its evidence timestamp and can invalidate round-two attribution. A malformed or stale receipt never authorizes another review or a self-issued waiver; finish the remediation and present the exact remaining policy decision to the operator.

An AI agent MUST NOT perform the final pull-request merge, enable auto-merge, or enter a merge queue unless the human operator explicitly authorizes that specific pull request. A merge authorization is single-use and expires when material state changes require new judgment.

## Push, release, and production boundaries

Do not push, create or move a tag, publish a release, publish a schema to the production domain, or mutate production `cueson.io` configuration without explicit authorization for that action. General instructions to finish, build, or use autopilot do not grant those authorities.

For S035 only, the operator explicitly authorized push and official PR, then directed that deployment and verification must not be left hanging. After the operator's final specific PR merge, continue S035 immediately through exact-main protected production deployment and independent live verification without requesting another kickoff or deployment approval. Keep #86 and milestone v1.2.0 open until their full acceptance criteria pass. This recorded authorization does not grant final PR merge authority or extend to unrelated releases or production changes.

## Post-merge housekeeping

Begin housekeeping only after the operator confirms a successful merge, then verify the GitHub merge state. Fetch and prune, fast-forward a clean non-divergent default branch, inspect slice branches and worktrees, and remove only state proven stale and clean. Never use broad destructive cleanup, forced worktree removal, `git reset --hard`, or routine force deletion.

Verify automatic remote head-branch deletion, post-merge CI, linked issues, parent epics, milestones, dependencies, and Project stages. Release work remains separately authorized. Report everything removed, retained, reconciled, and left pending.

## Changelog and release notes

`CHANGELOG.md` is the detailed release history and follows Keep a Changelog with a `Decisions` section for architecture changes. GitHub release notes contain highlights only, remain substantially shorter than the changelog, and end with:

```text
Full changelog: https://github.com/shruggietech/cueson/blob/vX.Y.Z/CHANGELOG.md
```

## Windows process behavior

Foreground, flashing, or focus-stealing console windows are prohibited. Direct `git` and `gh` commands are allowed. Project-owned Windows child-process launchers must use `CREATE_NO_WINDOW` or an equivalent hidden-process guarantee for console applications and disable interactive prompts.

## Workflowy

This project opts into the `workflowy` skill and its private agent workspace.
Load the skill when available. Use the operator's project link or find this
repo or folder under **AI Agent Workspace > PROJECTS**. Reuse its branch or
insert it alphabetically. Search **STANDARDS** by title or topic tag. Read the
applicable rules, defaults, and scopes before governed work.

Use normal bullets for candid discussion, drafts, decisions, and evidence.
Never store work in Workflowy comments or comment drafts. Give `info` exactly
three direct children: `tags`, `links`, and `details`. Put purpose, current state,
next actions, blockers, and supporting evidence under `details`.
Use lowercase letters and underscores for reusable field names. Use inline
code for exact technical literals. Link to the exact canonical finding or evidence.
Git and the established tracker govern official progress where used.
Otherwise use the project artifacts. The operator decides what is published.

In `notes`, keep major sessions and substantive updates as a flat, oldest-first
list. Use only the ISO 8601 date or timestamp as the title. Retain the finest
supported precision and offset. Put the description, scope tags, and prose
under `title`, `tags`, and `details`. Use label-only field bullets with each
value or link nested in its own bullet. Require the Workflowy bullet format
for fields and values. Only the contents of `details` may use other formats.
Keep tags out of names and other fields. Use one hashtag per `tags` value.
Put any project-level tags in `info > tags`. Link related updates. The operator
owns timeline construction. Agents record dates without creating timeline
metadata, views, saved searches, or exports.

For Spec Kit, use one direct `spec-kit` bullet per formal slice in code order.
Use only its code and short title. Nest the verified date under `date` and
the documented epic under `epic` when applicable. Use the skill for shared
reference records, migration, and precise field rules.

Desktop MCP requires a running, signed-in Workflowy app with Local MCP
enabled and a connected client. A cloud or other nonlocal agent uses a
separately configured, operator-authorized connector when available. Confirm
the destination and read back writes. If no connector is available, continue
under the available instructions and provide an appropriately private,
ready-to-paste nested outline with its exact destination. Mark it unsynced.
