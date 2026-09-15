# Implementation Plan: Fix production DNS verification

**Branch**: `codex/S022-fix-production-dns` | **Date**: 2026-09-15 | **Spec**: [spec.md](spec.md)

**Input**: S022 specification and issue #49.

## Summary

Extract only DNS verification into a small dependency-injected script module. Resolve A and AAAA concurrently with independent settled outcomes, validate family-specific IP evidence, deduplicate and sort accepted addresses, and fail with both family outcomes when neither is usable. Call the helpers for apex and www from the existing production verifier. Preserve the configured Google resolver, request timeout, and every non-DNS verification operation.

## Technical Context

**Language/Version**: JavaScript ESM on the existing Node.js 24 site baseline; compatible root Go toolchain.

**Primary Dependencies**: Built-in Node DNS, IP validation, global fetch, and node:test; no new dependencies or lockfile change.

**Storage**: Runtime-only DNS evidence; existing generated site metadata remains disposable.

**Testing**: Injected resolver/fetch tests in site/tests/production-dns.test.mjs, complete pnpm site suite, root tests, repository text/document checks, and full production verification.

**Target Platform**: Portable Node site verification and hosted Linux CI; local Windows launches use CREATE_NO_WINDOW and redirected non-interactive streams.

**Project Type**: Operational verifier within the existing static site package.

**Performance Goals**: Both family queries start independently; DoH retains its existing 20-second request timeout without added retries.

**Constraints**: Each resolver independently needs usable addresses. Preserve TLS, redirects, routes, downloads, metadata, schema bytes, the existing network-identity guard, and the manual deployment workflow.

**Scale/Scope**: Two hostnames, two resolver paths, two families per resolver. One issue and PR; no production mutation or release.

## Constitution Check

- Source fidelity, schema/software discipline, and common model: no source assets, Cue JSON, executable, or released schemas change.
- No silent loss and truthful operation: retain both failed family outcomes; invalid evidence cannot establish acceptance.
- Test-first behavior: focused regressions precede implementation; complete site and native root tests prove integration.
- Portable and secure operation: built-in IP validation, fixed existing DoH endpoint, unchanged trusted TLS/production authority, and hidden non-interactive Windows launches.
- Documentation and delivery: installed Spec Kit skills, blocking analysis, foreground verification, UTF-8 without BOM, logical Markdown paragraphs, issue closing reference, bounded reviews, explicit push/PR authority, and human merge.

Pre-research and post-design checks pass without constitutional exceptions.

## Project Structure

```text
specs/S022-fix-production-dns/
  spec.md
  plan.md
  research.md
  data-model.md
  quickstart.md
  contracts/dns-verification.md
  checklists/requirements.md
  checklists/dns.md
  tasks.md
site/scripts/production-dns.mjs
site/scripts/verify-production.mjs
site/tests/production-dns.test.mjs
site/package.json
docs/architecture.md
specs/S021-launch-cueson-site/contracts/deployment.md
CHANGELOG.md
```

**Structure Decision**: Keep DNS ownership under site/scripts with explicit dependency injection for deterministic tests. Avoid refactoring argument parsing, TLS, artifact verification, or Worker code. Existing Git/ESLint ignores cover dependencies and generated artifacts.

## Phases and Decisions

1. Specify and clarify: each resolver requires one valid family. Matching address sets or requiring both families is unnecessary; #49 and the constitution resolve routine decisions.
2. Checklist and research: generate a focused quality checklist, review it separately before implementation, and consult primary references through an independent research agent. The installed checklist prerequisites require a plan file, so setup-plan template initialization precedes the prerequisite retry; design completion follows requirements review.
3. Design and tasks: define helper contracts and test-first tasks; run read-only cross-artifact analysis before implementation.
4. Implementation: tests first, then DNS helpers, both-hostname integration, and registration in the explicit test:unit list.
5. Verification and publication: focused and complete foreground checks, public verification at the recorded S021 SHA, commit, authorized push/PR, and every review response with at most one second Codex request.

**Explicit departure from existing logic**: Replace the A-first failure path and A-only DoH query with independent A/AAAA evaluation. Require typed, valid addresses rather than any non-empty Answer list to close a false-pass hole at the same boundary.

**Alternatives rejected**: Adding production IPv4 DNS is unnecessary mutation; swallowing every error loses diagnostics; ANY does not guarantee both families; mocking the full production verifier expands scope; new DNS libraries duplicate built-in behavior.

**Verification authority**: The expected live SHA remains 46838fd5cc888b299a09b89da676db0005b00f16 although the verifier runs from S022. Generate local S022 inputs normally and pass the explicit live SHA; never require an unmerged S022 SHA from production or bypass network identity.

**Delivery evidence decision**: Local Spec Kit tasks cover verified implementation and publication/review preparation. The authorized post-commit transitions are tracked on GitHub so local checkboxes cannot falsely claim reviews that have not arrived. The full slice still ends only with a green reviewed PR returned for human merge.

**Local verification environment**: The machine's resolver returns ENOTFOUND for dns.google while production hostnames resolve. Google's public DNS response independently confirms A addresses 8.8.8.8 and 8.8.4.4; an authenticated HTTPS request to the original dns.google endpoint succeeds with a temporary process-only hostname bootstrap. Use that bootstrap for local full public verification, retain hostname and certificate validation, and change no repository/system resolver configuration. Preserve the underlying fetch cause in diagnostics so the ordinary failure remains actionable.

## Issue Acceptance Mapping

| #49 criterion | Planned evidence |
|---|---|
| IPv4-only passes both checks | US1 matrices for both helpers |
| IPv6-only passes both checks | US1 failed/empty A plus usable AAAA |
| Dual-stack deterministic | US1/US2 duplicates and completion order |
| Neither family fails usefully | US2 empty/error/invalid diagnostics |
| Complete suite and live pass | US3 site suite and full public verifier |
| Exact-main manual boundary retained | Unchanged workflow diff and deployment contract |
