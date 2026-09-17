# Implementation Plan: Repair protected production deployment

**Branch**: `codex/S032-repair-production-deploy` | **Date**: 2026-09-17 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/S032-repair-protected-production-deploy/spec.md`

## Summary

Repair the protected default-branch production workflow by removing malformed package-script separators, testing the real Corepack/pnpm command boundary, binding workflow execution ref and SHA to the requested exact-main revision, limiting the Cloudflare token to three Cloudflare-facing steps, and documenting a durable least-privilege credential lifecycle. Close a new independently testable implementation issue through the official pull request while retaining #78 for the post-merge default-branch validation. After human merge, provision the protected credential and use the already authorized production continuation to prove matching GitHub/public revision evidence and complete Cloudflare/public read-back without a recovery branch.

## Technical Context

**Language/Version**: ECMAScript modules on Node.js 24+, YAML GitHub Actions, Markdown, repository Go 1.25 publication tooling

**Primary Dependencies**: Corepack/pnpm, Node built-in test runner and child process API, Cloudflare Wrangler 4.131, GitHub Actions protected environments

**Storage**: UTF-8/no-BOM YAML, JavaScript, Markdown, normalized Cloudflare JSON snapshots, GitHub environment secret metadata

**Testing**: Node built-in subprocess tests, workflow policy assertions, complete site suite, Go repository checks, hosted Actions, CodeQL, Codex/security review, post-merge Cloudflare/public verification

**Target Platform**: GitHub-hosted Linux deployment runner and Cloudflare Workers/Assets serving cueson.io; tests remain cross-platform on Windows and Linux

**Project Type**: Deployment automation plus documentation and production verification

**Performance Goals**: Keep focused command-boundary tests bounded to parser startup; preserve the current complete deployment verification duration and no-network static build behavior

**Constraints**: Manual-only production mutation; exact current main; protected environment approval; least-privilege secret with step-only exposure; no secret values in repository or logs; hidden non-interactive child processes on Windows; no automatic final merge

**Scale/Scope**: One workflow, one focused test file plus existing policy tests, three maintained documentation surfaces, one implementation issue, one retained operational issue, one official PR, and one post-merge production run

## Constitution Check

*GATE: Passed before implementation and re-checked after design.*

- **I. Lossless Source Preservation**: PASS. No Cue JSON or source-media behavior changes.
- **II. Schema and Official Software Discipline**: PASS. No schema, tag, release, or product executable changes. Production verification preserves exact immutable schema checks.
- **III. Common Model Plus Native Fidelity**: PASS. No format model or native fidelity change.
- **IV. No Silent Loss**: PASS. Deployment and state verifiers continue to fail closed on identity, inventory, and preservation mismatches.
- **V. Test-First Format Work**: PASS. No format work. The workflow defect receives a failing real-boundary test before implementation.
- **VI. Portable and Secure Operation**: PASS. Subprocess tests use direct executables with hidden Windows execution. The credential is restricted by resource, permission, protected environment, lifecycle, and step exposure.
- **VII. Documentation and Delivery Authority**: PASS. Spec Kit artifacts, blocking analysis, a separately closeable implementation issue, official CI/reviews, and retained post-merge validation preserve honest delivery state.

No constitutional exception is required.

## Project Structure

### Documentation (this feature)

```text
specs/S032-repair-protected-production-deploy/
├── checklists/
│   └── requirements.md
├── contracts/
│   └── protected-deployment.md
├── data-model.md
├── plan.md
├── quickstart.md
├── research.md
├── spec.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
.github/workflows/site-deploy.yml
CHANGELOG.md
docs/
├── architecture.md
└── release-process.md
site/
├── package.json
├── scripts/
│   ├── verify-cloudflare-state.mjs
│   └── verify-production.mjs
└── tests/
    ├── generator.test.mjs
    └── workflow-command.test.mjs
```

**Structure Decision**: Keep deployment behavior in the existing manual workflow and verifier scripts. Add one focused subprocess test file at the package-script boundary where the defect escaped. Do not expand runtime architecture or duplicate verifier parsers.

## Phase 0: Research Conclusions

1. Separate implementation merge from production proof. The implementation issue closes through the PR; #78 remains open until a merged main run proves the environment and public state.
2. Bind workflow dispatch identity to main. `GITHUB_REF` and `GITHUB_SHA` must match the requested revision before checkout so GitHub deployment metadata describes the same commit later recorded publicly.
3. Pass named arguments directly through pnpm. The standalone separator is a workflow defect, not a verifier feature.
4. Test the actual command boundary. Direct Corepack subprocesses cover package-manager forwarding and parser behavior together without credentials or mutation.
5. Use an account-owned resource-restricted API token. Entire Account Workers Scripts Legacy Edit plus `cueson.io` Zone Read, DNS Read, Zone Transform Rules Read, and Workers Routes Read cover Wrangler upload, Custom Domains, route-conflict detection, and the full state verifier.
6. Expose the token at step scope only. Checkout, setup, dependency installation, tests, and local artifact proof do not need production credentials.
7. Rotate within 90 days, validate replacements before retirement, and revoke immediately on suspected exposure or ownership change.

## Phase 1: Design Decisions

1. Add pre-checkout assertions for `refs/heads/main` and `GITHUB_SHA == REVISION`; keep both existing fresh `origin/main` checks.
2. Add `environment.url` pointing to the public deployment record so protected deployment evidence links to the exact public authority.
3. Remove the job-level token and add the environment secret only to Cloudflare preflight, Wrangler deployment, and post-deployment read-back.
4. Replace all three malformed commands with direct named arguments and update static workflow expectations.
5. Add cross-platform subprocess tests using `corepack` or `corepack.cmd`, `shell: false`, `windowsHide: true`, and sanitized Cloudflare variables. Assert each invocation passes argument parsing and reaches a deterministic downstream validation failure.
6. Update architecture, release process, and changelog to record execution identity, least-privilege credential ownership/scope/lifecycle, and the corrected package-script boundary.
7. Repair #78's Markdown body during publication, create the implementation child issue with all governed sections, and ensure both issues appear exactly once in the Delivery Project with Slice S032 and empty default Status.

## Verification Strategy

1. Red phase: add subprocess boundary and static workflow expectations, run them against the current workflow, and capture the invalid-separator failure.
2. Green phase: update the workflow and documentation, then rerun focused tests.
3. Complete local gates: full site suite; applicable Go tests/build/vet/race/staticcheck/govulncheck; github-format, UTF-8/no-BOM, mojibake, whitespace, generated-state, and diff checks.
4. Spec Kit gates: blocking analysis before source implementation and convergence after implementation.
5. Hosted gates: push authorized branch; publish formatted/read-back PR; wait for every exact-head CI, CodeQL, security, and Codex result; resolve all findings; request at most one second Codex review if round one has findings.
6. Post-merge gate: provision/read back protected credential metadata, dispatch exact current main from main, inspect workflow/environment record, verify Cloudflare before/after state and public deployment evidence, then close #78.

## Lifecycle and Authority

1. Create an S032 implementation issue as a child/dependency of #78. Move it from Specced to In progress to PR review, then Done only after merge.
2. Repair #78, retain it In progress through the PR, and move it to Release verification only after human merge and credential provisioning.
3. The user explicitly authorizes branch push, official PR publication, review handling, and necessary production deployment work. The AI remains prohibited from final merge, auto-merge, or merge queue entry.
4. Production mutation occurs only from the merged default-branch workflow. This pull request changes no live Cloudflare state.

## Complexity Tracking

No constitutional violations or exceptional complexity are introduced. The extra implementation issue is required to keep PR completion and post-merge operational proof independently truthful.
