# Implementation Plan: Prepare v1.2.0 public schema and site

**Branch**: `codex/S035-publish-v1-2-public-schema-site` | **Date**: 2026-10-06 | **Spec**: [spec.md](spec.md)

**Input**: Owner-authorized S035 and preparation #88; deployment #86 remains separate.

## Summary

Promote the existing static content authority to verified v1.2.0, add its immutable schema plus release/speaker documents, reconcile current publication records, and extend independent site expectations. Retain generic generator, Worker and protected deployment architecture. Review and validate the finished artifact, push and open the official PR, handle all findings and stop for final human review/merge. After the owner merges, continue S035 with actual-main production deployment and independent live verification under their subsequent instruction; no second kickoff/approval is required.

## Technical Context

**Language/Version**: TypeScript 5.9/ECMAScript modules, Node 24+, Go 1.25 verification modules, PowerShell 7 hidden-process tooling.

**Primary Dependencies**: Existing Next 16 static export, Fumadocs, React, Playwright/axe, Wrangler; no new dependencies.

**Storage**: Maintained UTF-8/no-BOM JSON/Markdown, immutable schema snapshots; ignored deterministic generation/build outputs.

**Testing**: ESLint, generation/drift, Node unit/negative tests, Next build, Playwright/accessibility, artifact validation, Wrangler dry run, documentation/module checks and full hosted required checks.

**Target Platform**: Static Cloudflare Workers/Assets site, desktop/tablet/mobile browsers; headless local Windows and hosted Linux site proof.

**Project Type**: Documentation/static artifact preparation.

**Performance Goals**: Retain current static export behavior with no new client service or media processing.

**Constraints**: Immutable historical schema/package/tag bytes; no source-runtime behavior change; no production credentials/mutation before reviewed merge; no final merge authority.

**Scale/Scope**: 22 document entries, 24 HTML routes, four schemas, 28 total declared routes, seven primary downloads.

## Constitution Check

Pre-research and post-design gates PASS: preserve original source/runtime contracts; peer software/schema identity remains 1.2.0; released schema bytes unchanged; no silent loss or new conversion behavior; independent inventory/negative tests before authority updates; existing security/accessibility tests retained; Spec Kit and blocking analysis precede implementation; push/PR explicitly authorized; final merge remains operator-owned; subsequent owner steering authorizes S035 production deployment immediately after that merge.

## Project Structure

```text
specs/S035-publish-v1-2-public-schema-site/
  spec.md, plan.md, research.md, data-model.md, quickstart.md
  tasks.md, analysis.md, verification.md, issue-map.json
  checklists/, contracts/
site/content-map.json
site/tests/generator.test.mjs
site/tests/production-verification.test.mjs
site/tests/site.spec.ts
site/worker/index.test.ts
README.md, SECURITY.md, CONTRIBUTING.md, CHANGELOG.md
docs/architecture.md, schema.md, compatibility.md, consumer-speakers.md
docs/roadmap.md, project-management.md, release-process.md
docs/release-verification.md, releases/v1.2.0.md
docs/conversion.md, Cueson-Project-Specification-v0.0.0.md
scripts/docs-verify/verify.go, verify_test.go, candidate_v120_test.go
scripts/docs-verify/s035_publication_test.go
```

**Structure Decision**: Reuse maintained authorities and generic site generation. Independent fixed test inventories are intentionally separate from mutable content-map values. Generated files stay ignored.

## Decision log

- Split preparation #88 from live #86 rather than auto-closing an incomplete production outcome. Both remain in milestone 4; #86 natively depends on #88 and completed #85.
- Promote v1.2.0 only after reading back actual public release/assets and publication evidence. Preserve prior release navigation and immutable sources.
- Extend existing inventory expectations rather than introduce a redundant generator/configuration layer. The invalid-download test must continue using mismatched version/filename after promotion.
- Reconcile maintained current-state prose, preserve dated prior evidence and frozen specs. Non-publishing proof remains valid after a public release.
- Use isolated browser port 4174 through a temporary ignored configuration; do not reuse a potentially stale 4173 preview.
- Recorded production authorization applies to actual squash-merged main and fresh protected proof, not this PR's pre-merge commit. DNS alone does not prove address-family connectivity.
- No architecture or workflow mutation is planned; if verification proves one necessary, record its reason and a dated changelog decision before review.
- Integration exposed stale documentation-verifier rules that still required v1.1 current-state and prohibited v1.2 publication. Update those current markers and focused regression tests to the verified publication state; preserve historical parsers and evidence contracts.

## Complexity Tracking

No constitutional exceptions or new runtime layers.
