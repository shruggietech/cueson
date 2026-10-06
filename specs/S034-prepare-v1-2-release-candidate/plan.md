# Implementation Plan: Prepare the v1.2.0 release candidate

**Branch**: `codex/S034-prepare-v1-2-release-candidate` | **Date**: 2026-10-06 | **Issue**: [#84](https://github.com/shruggietech/cueson/issues/84) | **Spec**: [spec.md](spec.md)

## Summary

Promote the reviewed S033 model/schema/CLI contract to exact 1.2.0 as a non-publishing stable candidate. Preserve published 0/1/1.1 bytes and public v1.1 availability. Extend package proof instead of merely relabeling development artifacts. Candidate preparation, public release and production hosting are independently testable milestone outcomes.

## Technical Context

**Language/Version**: Existing pure Go 1.25+ executable and nested verification modules, JSON Schema 2020-12, Markdown, installed PowerShell Spec Kit scripts and existing Node/pnpm site tooling.

**Primary Dependencies**: Existing local schema registry, consumer model validators, four native codecs, twelve-edge converter, release verifier, pinned GoReleaser v2.18.1/Syft v1.51.1 and published prior release bytes.

**Storage**: Immutable candidate schema and versioned release evidence/decision contracts; no new persisted consumer fields.

**Testing**: Focused exact-identity/immutable-copy/former-development refusal, version-aware old-consumer verification, consumer native proof and prepared documentation tests. Foreground product/nested tests, formatting, vet, static/security analysis, fuzz/build/site gates plus hosted race/native/CodeQL/package proof.

**Target Platform**: Windows/macOS/Linux amd64 and arm64 packages; same-bundle native execution on three amd64 hosts only.

**Constraints**: Source bytes, released artifacts and historical snapshots unchanged; no production download/route admission; all Windows tool launchers hidden and noninteractive; no public tags/releases/assets/deployment.

**Scale/Scope**: Fourteen requirements, three stories, one candidate issue/PR, later publication and hosting issues. Six archives plus six SBOMs plus one manifest. Six historical documents. Existing 1.0 proof plus authenticated 1.1 refusal.

## Constitution Check

Pre-research and post-design PASS. No governance amendment: explicit release separation, immutable history, source preservation, deterministic omission/strict refusal, test-first focused changes, pure-Go bounded validation, full installed Spec Kit sequence and blocking analysis. Specific merge governance remains applicable to the eventual reviewed PR. Public release authority is withheld.

## Phase 0: Research

Three independent agents reviewed exact schema promotion, package/old-consumer proof and public/status documentation. Consolidated decisions are in [research.md](research.md). Current stable proof dispatch is hardcoded to 1.1.0; its omission on naive promotion is a release correctness defect to fix in this slice. Published 1.1 uses a generic exact-contract diagnostic, so an authenticated version-aware matcher is required.

## Phase 1: Design

[data-model.md](data-model.md), [contracts/candidate-contract.md](contracts/candidate-contract.md) and [quickstart.md](quickstart.md) define the matching bytes, native evidence and observable adoption checks. Final publication contract freezes policy and intended inventory, but never invents an actual post-merge revision or future accepted digests.

## Structure and ownership

- Schema agent: internal/schema/**, internal/version/**, internal/model/consumer.go identity, schema/releases/v1.2.0/**, current fixture/conformance/CLI expectation-only identity updates. Preserve historical dirs and assertions.
- Proof agent: scripts/release-verify/**, .goreleaser.yaml, .github/workflows/release-proof.yml, new S034 evidence/old-consumer contracts. Preserve S020/S028/S029 files.
- Documentation agent: maintained docs, README/CONTRIBUTING, dated CHANGELOG and docs/releases/v1.2.0.md; no site published inventory or historical release/spec edits.
- Root: S034 Spec Kit/project metadata, publication decision contract, scripts/docs-verify/**, integration, site check, full final verification, GitHub delivery and reviews. No agent commits or pushes.

## Decisions and alternatives

- Promote exact stable candidate identity instead of retaining development output or modifying published 1.1.0. The owner selected a release candidate; immutable-copy equality becomes mandatory.
- Keep former development identity unsupported rather than inventing a migration contract.
- Retain immutable earlier proof schemas and add separate 1.2 evidence/old-consumer records instead of overwriting prior publication truth.
- Verify actual published 1.1 archives with target/revision/schema/legal/version/digest binding and positive historical acceptance before refusal; a source rebuild cannot substitute.
- Preserve existing 1.0 proof while proving 1.1 for the new candidate. Evidence records each consumer independently.
- Same-bundle current annotation tests include each source format, counts-only inspection, timing invalidity, warning-only source conflicts, exact restore, loss reports and strict occupied/absent-output refusal.
- Reconcile stale current-state status separately from dated narrative. Keep primary public installation and site metadata at verified 1.1.
- Custom quality markers remain reviewer-owned. Explicit unattended scope authorizes proceeding after clean analysis without a new checklist permission prompt.
- Prior explicit push/official-PR authority and the current instruction to halt only at public release allow repository delivery through review; tags/releases/production remain outside scope. Do not silently bypass final specific-PR merge governance.

## Delivery and verification

Run specify → clarify → checklist → plan → tasks → read-only analyze → implement → foreground required checks → conventional commit → authorized branch/official PR and hosted gates. Resolve every review finding and request at most one second Codex round only after a finding-driven remediation. Present concrete reviewed candidate and publication boundary; do not claim PR-head proof equals post-squash main proof.
