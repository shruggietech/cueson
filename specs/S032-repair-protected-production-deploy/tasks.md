# Tasks: Repair protected production deployment

**Input**: Design documents from `specs/S032-repair-protected-production-deploy/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Real package-script command-boundary coverage and workflow policy coverage are required before implementation.

**Organization**: Tasks separate the independently mergeable implementation outcome from the post-merge production proof retained by #78.

## Phase 1: Setup and governed issue repair

**Purpose**: Freeze the baseline, repair issue authority, and establish independently truthful lifecycle outcomes.

- [x] T001 Record the exact-main baseline, production health, open-issue inventory, recovery-state status, protected authority boundary, and initial branch state in `specs/S032-repair-protected-production-deploy/verification.md`.
- [x] T002 Repair issue #78 as formatted/read-back Markdown with every governed section and its post-merge acceptance criteria intact.
- [x] T003 Create the S032 implementation issue with Outcome, Context, Scope, Acceptance criteria, Dependencies, and Verification; link it natively to #78 when supported; add it exactly once to the Delivery Project with Slice S032, Stage Specced, and empty default Status.

---

## Phase 2: Requirements and blocking analysis

**Purpose**: Resolve ambiguity and coverage gaps before source implementation.

- [x] T004 Validate `checklists/requirements.md`, the least-privilege permission research, and the implementation/production issue split against the constitution, architecture, S030 contract, and S031 recovery evidence.
- [x] T005 Run the Spec Kit blocking analysis across `spec.md`, `plan.md`, `tasks.md`, checklists, contract, issue acceptance, and repository authority; resolve every critical, high, or material coverage/consistency finding.
- [x] T006 Read back both issues and Delivery Project fields, retain #78 for production validation, and move the implementation issue to In progress after the analysis gate passes.

**Checkpoint**: Requirements and lifecycle authority are complete; source implementation may begin.

---

## Phase 3: User Story 1 - Run the protected deployment command path (Priority: P1)

**Goal**: Repair package-script forwarding and bind protected workflow identity to exact current main.

**Independent Test**: Real Corepack/pnpm subprocess tests reach verifier validation beyond argument parsing, while static workflow tests reject literal separators, non-main dispatch identity, or removed exact-main checks.

- [x] T007 [US1] Add focused cross-platform subprocess expectations in `site/tests/workflow-command.test.mjs` and wire them into `site/package.json` using direct hidden non-interactive execution.
- [x] T008 [US1] Update workflow policy expectations in `site/tests/generator.test.mjs` for direct named arguments, main ref/SHA binding, step-scoped token exposure, environment URL, and retained dual exact-main checks.
- [x] T009 [US1] Run focused tests against the existing implementation, prove the literal separator and missing identity/secret-scope policies fail for the intended reasons, and record red evidence in `specs/S032-repair-protected-production-deploy/verification.md`.
- [x] T010 [US1] Update `.github/workflows/site-deploy.yml` to remove literal separators, require main execution ref and exact workflow SHA before checkout, retain both fresh exact-main gates, and set the deployment environment URL.
- [x] T011 [US1] Rerun command-boundary and static workflow tests and record green US1 evidence in `specs/S032-repair-protected-production-deploy/verification.md`.

**Checkpoint**: The reviewable workflow command and revision path is exact, fail closed, and independently tested.

---

## Phase 4: User Story 2 - Operate a durable least-privilege credential (Priority: P1)

**Goal**: Restrict secret exposure and publish a complete durable credential contract.

**Independent Test**: Workflow policy proves only three Cloudflare-facing steps receive the token, and documentation defines the exact owner, resources, permissions, rotation, expiry, revocation, and replacement checks.

- [x] T012 [US2] Remove `CLOUDFLARE_API_TOKEN` from workflow job scope and add it only to Cloudflare preflight, deployment, and post-deployment read-back steps in `.github/workflows/site-deploy.yml`.
- [x] T013 [US2] Document account ownership, exact account/zone resources, Workers Scripts Write plus Zone/DNS/Transform Rules Read permissions, 90-day rotation, expiry handling, replacement validation, emergency revocation, and non-secret audit evidence in `docs/release-process.md`.
- [x] T014 [US2] Update `docs/architecture.md` and `CHANGELOG.md` with the exact execution-identity, step-scoped credential, and corrected forwarding decisions.
- [x] T015 [US2] Run focused workflow/documentation checks and review the full verifier call surface against the permission contract; record green US2 evidence in `specs/S032-repair-protected-production-deploy/verification.md`.

**Checkpoint**: The durable credential is least privilege, operationally documented, and unavailable to unrelated workflow steps.

---

## Phase 5: Cross-cutting verification and convergence

**Purpose**: Prove repository-wide quality and converge the official review head.

- [x] T016 Run the complete site suite (`lint`, generation/check, unit, build, browser, artifact, and Wrangler dry run) and record outcomes in `specs/S032-repair-protected-production-deploy/verification.md`.
- [x] T017 Run applicable root and standalone Go tests/build/vet/race/staticcheck/govulncheck, github-format checks, UTF-8/no-BOM/mojibake/whitespace checks, `git diff --check`, and generated-state review.
- [x] T018 Run Spec Kit convergence across requirements, issue acceptance, workflow/test/documentation changes, credential permissions, exact-main gates, and production boundary; resolve every material finding and update `tasks.md` plus `verification.md`.
- [x] T019 Commit conventionally, verify the exact committed tree, and confirm the implementation issue is In progress/S032 while #78 remains open/In progress/S032 with both default Status values empty.

---

## Phase 6: Official pull request and review convergence

**Purpose**: Publish the explicitly authorized implementation and reach a terminal reviewed exact head.

- [ ] T020 Push the authorized branch and publish a github-format/read-back official PR with `Closes` for the implementation issue and `Refs #78`; move both issues to PR review and clear default Status.
- [ ] T021 Wait for all exact-head CI, CodeQL, Site, security, and first-round Codex results; inspect every review, comment, thread, check, and reaction and address every actionable finding with focused plus complete verification.
- [ ] T022 If round one contains findings, request exactly one second Codex review with `@codex review`, then handle every second-round result without requesting a third automatic round; resolve threads only after concerns are satisfied.
- [ ] T023 Publish formatted/read-back completion evidence on the PR, confirm all required checks are green at the final head and no actionable review remains, then hand off for the human final review and merge ritual.

## Phase 7: Post-merge production validation

**Purpose**: Complete the operational outcome after the human merge.

- [ ] T024 After operator-confirmed merge, perform standard safe housekeeping, verify exact current main, provision/read back the protected credential metadata, and move #78 to Release verification.
- [ ] T025 Dispatch the merged workflow from `main` with exact current main under the standing production authorization; inspect its protected environment deployment record and every workflow step.
- [ ] T026 Verify complete Cloudflare before/after state, public DNS/TLS/redirects, 25 routes, seven downloads, metadata, manifest, and three immutable schema bytes; prove GitHub and public deployment records identify the same exact main revision.
- [ ] T027 Remove only temporary state proven stale and clean, publish formatted/read-back completion evidence, close #78, reconcile Project state, and report final production status.

## Dependencies and Execution Order

- T001-T006 establish the blocking governance and analysis gate.
- T007-T011 implement US1 test first. T012-T015 complete US2 against the same workflow after US1 is green.
- T016-T019 require all reviewable implementation work complete. T020-T023 require a committed converged head.
- T024-T027 require the human to merge the official PR. They remain part of S032 but outside the pre-merge handoff.

## Implementation Strategy

1. Repair issue authority and pass blocking analysis.
2. Make the real workflow boundary fail in focused tests.
3. Repair argument forwarding and exact execution identity.
4. Restrict and document the durable credential.
5. Run complete verification and Spec Kit convergence.
6. Publish and converge the official PR.
7. Stop for the human merge ritual, then complete the already authorized production validation.
