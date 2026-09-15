# Tasks: Fix production DNS verification

**Input**: S022 spec, plan, research, runtime model, and DNS contract.

**Tests**: Required test-first regressions under the operator's autopilot kickoff and constitution.

## Phase 1: Setup

- [x] T001 Confirm clean feature branch, existing dependencies, Git/ESLint ignore coverage, and headless Windows tooling against `site/package.json`, `.gitignore`, and `site/eslint.config.mjs`.

## Phase 2: Foundational

- [x] T002 Complete separate requirements review and blocking Spec Kit cross-artifact analysis for `specs/S022-fix-production-dns/spec.md`, `plan.md`, `tasks.md`, and `checklists/`.

## Phase 3: User Story 1 - Either address family (P1)

**Independent Test**: Each injected resolver passes IPv4-only, IPv6-only, dual-stack, sibling failure, duplicate, and completion-order scenarios.

- [x] T003 [US1] Write acceptance matrix and synchronous/asynchronous sibling-failure regressions in `site/tests/production-dns.test.mjs` and confirm the initial missing-helper failure.
- [x] T004 [US1] Implement independent system and DoH address-family helpers, built-in IP validation, sorted unique inventories, and the existing DoH endpoint/header/timeout in `site/scripts/production-dns.mjs` (FR-001, FR-002, FR-003, FR-005).

## Phase 4: User Story 2 - Truthful address failures (P1)

**Independent Test**: Each resolver rejects missing/invalid evidence and reports hostname, resolver, and both family outcomes; a usable sibling still succeeds.

- [x] T005 [US2] Add no-answer, alias-only, wrong-family, malformed payload/status/Answer, HTTP/transport/JSON failure, and deterministic diagnostic regressions to `site/tests/production-dns.test.mjs` before completing failure behavior.
- [x] T006 [US2] Implement fixed-family failure aggregation retaining operational codes/messages in `site/scripts/production-dns.mjs` and pass the complete focused suite (FR-003, FR-004, FR-005, FR-006).

## Phase 5: User Story 3 - Full acceptance and delivery (P1)

**Independent Test**: The full site suite and public verifier pass the recorded S021 revision; protected probes and deployment workflow remain intact.

- [x] T007 [US3] Integrate both DNS helpers for apex and www in `site/scripts/verify-production.mjs` and register the DNS regression file in `site/package.json` without altering the network guard or non-DNS probes (FR-001, FR-006, FR-007).
- [x] T008 [P] [US3] Document the dual-family rule in `docs/architecture.md` and `specs/S021-launch-cueson-site/contracts/deployment.md`; add the scoped fix and dated script decision to `CHANGELOG.md` (FR-009).
- [x] T009 [US3] Run the focused suite, complete site suite including Wrangler dry run, root tests, text/document checks, and full public verifier from `specs/S022-fix-production-dns/quickstart.md`; inspect the unchanged `.github/workflows/site-deploy.yml` (FR-007, FR-008).

## Phase 6: Publication and Review Preparation

- [x] T010 Prepare the official PR body and exact verification evidence from `specs/S022-fix-production-dns/quickstart.md`, include `Closes #49`, and pass the Markdown through the publication formatter before commit/push (FR-010).
- [x] T011 Record the authorized push/PR, Project read-back, hosted-check, and bounded review continuation in `specs/S022-fix-production-dns/quickstart.md` and `tasks.md` without falsely marking external convergence complete before publication (FR-010, SC-004).

## Dependencies and Execution Order

T001 -> T002 -> T003 -> T004 -> T005 -> T006 -> T007 -> T009 -> T010 -> T011. T008 uses separate documentation files and can run alongside T007 once the helper contract is stable. Both resolver paths share one small helper module and test file, so US1/US2 edits remain sequential. US3 depends on both resolver stories for full acceptance.

## Parallel Examples

For US1 and US2, run controlled A/AAAA promises concurrently inside tests but edit the shared helper/test files sequentially. For US3, documentation T008 can proceed independently of verifier integration T007. Independent primary-source research already ran alongside specification authoring.

## Implementation Strategy

Prove healthy-family behavior first, then failure truthfulness, then complete production acceptance. US1 is the minimal local increment; #49 remains open until all three stories and reviewed delivery are complete. No production deployment or final merge is part of agent execution.

## Authorized Delivery Continuation

After the verified local implementation and preparation tasks are complete, commit and push S022, publish the formatted official PR, immediately verify its GitHub body, and set the existing Project item to Slice S022 and Stage PR review. Wait for hosted checks and external reviews, handle and reply to every finding, resolve only handled threads, and use at most one second Codex request if necessary. Track actual external convergence on the PR, then return its green reviewed head for human final review and merge. Local task completion does not claim these subsequent GitHub transitions have already occurred.
