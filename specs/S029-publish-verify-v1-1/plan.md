# Implementation Plan: Prepare reviewed v1.1.0 publication

**Branch**: `codex/S029-publish-verify-v1-1` | **Date**: 2026-09-15 | **Spec**: [spec.md](spec.md)

**Input**: S029 autopilot after merged S028 main 6a55b48c654e14f5baaf89a9add51fa4f5a4d04a; preparation #74 is a native child/blocker of #66.

## Summary

Finalize reviewed v1.1.0 history and exact publication notes before selecting a release target. Deliver an official preparation PR closing #74 and referencing #66. Specify a normative decision contract for the actual future post-merge main proof, thirteen accepted assets, same-bundle native evidence, exact notes and action-specific authority. Leave publication #66 and hosting #67 open. No shipped runtime/schema/Site/config/workflow/publisher change.

## Technical Context

**Language/Version**: Go 1.25 compatibility baseline; standard-library standalone verification modules; existing pinned Site Node/pnpm policy.

**Primary Dependencies**: Existing release/docs verifiers and formatting tool, exact local GoReleaser v2.18.1/Syft v1.51.1, existing CI/security/platform/Site checks. No dependency addition.

**Storage**: UTF-8/no-BOM Markdown contracts/docs and ignored candidate/evidence output.

**Testing**: Meaningful positive/negative release-heading/date/order/notes and decision-policy checks; root/nested tests/vet/static/vulnerability, formatting/docs/brand, six pure-Go builds and full Site suite; exact source-bound six-target package proof.

**Target Platform**: Existing Windows/macOS/Linux amd64 native hosted proof and six amd64/arm64 structural build targets.

**Project Type**: Governed metadata and operational decision preparation.

**Performance Goals**: Existing CI/candidate proof budgets; no new product performance target.

**Constraints**: No tag/release/production/final merge, no guessed squash target, no stale bundle, no source/schema/public inventory changes, foreground verification, hidden noninteractive Windows launch.

**Scale/Scope**: One preparation child/PR, one dated history, one public notes body, one normative decision contract, thirteen future public assets.

## Constitution Check

| Principle | Assessment |
|---|---|
| I Source preservation | PASS: No source/corpus/envelope behavior or byte changes. |
| II Peer products | PASS: Exact 1.1.0 software/schema/immutable package identity retained; historic bytes unchanged. |
| III Common/native fidelity | PASS: Existing bounded profile and truthful disclosures retained. |
| IV No silent loss | PASS: No conversion/model behavior change; existing loss/refusal checks retained. |
| V Test first | PASS: Metadata/notes policy checks precede finalization and prose changes. |
| VI Portable/security | PASS: Existing six targets/three governed native proofs; no publisher or new authority surface. |
| VII Delivery authority | PASS: Installed blocking analysis/convergence; explicit branch push/PR authority; separate protected release/main-merge authority. |

No exception is required.

## Project Structure

```text
specs/S029-publish-verify-v1-1/
  spec.md, plan.md, research.md, data-model.md, quickstart.md, tasks.md
  contracts/publication-decision.md
  contracts/release-notes.md
  checklists/requirements.md, checklists/publication-integrity.md
  verification.md
CHANGELOG.md
docs/roadmap.md, docs/Cueson-Project-Specification-v0.0.0.md
README.md, docs/release-process.md, docs/release-verification.md
scripts/docs-verify/verify.go, scripts/docs-verify/verify_test.go
scripts/release-verify/policy_test.go
```

**Structure Decision**: Add only normative Markdown contracts and policy/metadata checks. Existing release-verifier CLI and accepted integrity checks remain the execution authority. Populate the actual machine decision record only after preparation merge and accepted main proof.

## Recorded decisions

1. Separate #74 preparation from #66 verified public publication. Tagged metadata must be on reviewed main before target selection; a PR closing #66 now would claim unexecuted acceptance. #74 is a native child/blocker of #66; preserve existing #65/#67 relationships.
2. Preserve S029 operator code and branch. Explicit push/official PR authority overrides the skill's usual pre-push pause; final merge/tag/release/production remain excluded.
3. Move all accumulated Unreleased history into a dated 1.1.0 section with fresh Unreleased. Consolidate the duplicated Changed category; preserve each item and chronological Decisions. Earlier 1.0.0/0.0.0 release sections remain byte-identical.

Release finalization deliberately places the new S029 feature/dated governance entry into the prepared 1.1.0 section, retaining a fresh empty Unreleased section, rather than leaving tagged preparation history under Unreleased as the skill's ordinary feature default would. Existing release-section byte proof excludes the comparison-link footer, where adding the new version and advancing Unreleased is intentional.
4. Keep current candidate notes/status truthful. Author a separate exact public body under contracts/release-notes.md with compatibility/limitations and required suffix, rather than silently stripping candidate banners at publication time.
5. Bind future approval to the actual post-merge main source and freshly accepted same-bundle proof. S028 main/PR proof is historical preparation evidence; the future squash result cannot be populated before merge.
6. Normative contract lists the exact target/inventory/proof/notes/schema/legal/authority/refresh/partial-publication obligations; no machine record with speculative hashes, populated approvals or guessed source is committed.
7. Do not add a publisher or new build mode. Existing verifier already rejects source/target/checksum/schema/legal/privacy/clean-build drift; avoid duplicating its negative tests merely under another version string.
8. Test first for new date/order/category/notes boundaries and focused policy checks; complete appropriate current CI parity and fresh committed candidate proof. Full repeated fuzz is unnecessary for unchanged product code; retain automatic hosted fuzz/race/native coverage.
9. Independent release-contract and documentation owners have isolated files; coordinator owns integration, blocking analysis, convergence and final verification.
10. Installed checklist requires plan.md; create only the installed setup-plan scaffold before custom checklist and substantive planning afterward. Reviewer-owned custom markers remain unchanged; all eight requirements-quality criteria are explicit under routine autopilot.
11. Record external exact-head completion in a formatted/read-back PR comment to avoid evidence-only source changes. Closing #74 leaves #66/#67, epic #51 and milestone open.

## Ownership and ordered implementation

1. Coordinator: Spec Kit artifacts excluding the two release-owner contracts, CHANGELOG/roadmap/main working specification, issue/Project/PR lifecycle, integration and complete verification.
2. Release contract owner: contracts/publication-decision.md, contracts/release-notes.md and standalone scripts/release-verify/policy_test.go additions.
3. Documentation owner: README, release-process/release-verification/current candidate notes if needed, scripts/docs-verify checks/tests; preserve technical docs and Site unless a direct contradiction requires a justified narrow change.
4. All product edits follow T004 blocking analysis. Owner tests precede coordinator changelog/prose completion. Independent cross-owner audit and installed convergence follow integration.

## Verification and delivery

Run focused metadata/policy tests, root and six nested module tests/vet/static/vulnerability, six pure-Go builds, Go/text/encoding/actionlint/docs/brand checks and full generated Site lint/unit/build/browser/artifact/dry-run. Existing hosted CI also runs unchanged sixteen-target fixed-work fuzz, race and three-platform native tests. Commit conventionally without a co-author trailer, matching repository history. Build/verify exact clean committed packages and matching Windows native/old-consumer proof. Push and create formatted/read-back non-draft PR closing #74, referencing #66, reconcile Project Stage PR review/Slice S029/default Status clear, wait all checks/reviews, and request at most one second Codex review only for first-round findings. Hand off to the human final merge ritual.

After a separately authorized preparation merge, housekeeping must accept fresh exact resulting-main proof and populate a reviewable exact publication decision package before requesting tag/release authority. This follow-on gate is not completed by a green preparation PR.
