# Implementation Plan: Freeze the v1.1.0 release candidate

**Branch**: `codex/S028-freeze-v1-1-release-candidate` | **Date**: 2026-09-15 | **Spec**: [spec.md](spec.md)

**Input**: S028 operator autopilot combining #64 and #65 after exact merged S027 main `351036cc0c6972aa6f5953af591d4f49df718695`.

## Summary

Freeze maintained/native/schema/CLI/generated-preview contracts and then promote the exact current software/schema contract to 1.1.0 stable bounded ASS/SSA. Preserve exact historical 1.0.0 input and immutable historical bytes. Reuse existing stable release verifier and six-target non-publishing configuration, strengthen matching packaged native workflow execution and verify published old-consumer refusal on fresh output. Close remaining stable conformance deferrals with genuine evidence. Publish one official reviewed PR closing both atomic outcomes without protected release/production transitions.

## Technical Context

**Language/Version**: Root Go 1.25.0 compatibility; six nested standard-library modules; Site Node 24/corepack/pnpm under the pinned repository policy.

**Primary Dependencies**: Existing JSON Schema registry/compiler, x/text, Cobra-independent CLI catalogue; Next/Fumadocs and repository-authored Site generator; pinned GoReleaser 2.18.1/Syft 1.51.1; existing static/security tool versions.

**Storage**: Tracked UTF-8 schema/docs/examples/evidence contracts; ignored dist candidate artifacts and temporary native smoke directories.

**Testing**: Root/nested Go tests/vet/staticcheck/vulnerability, 16 fixed-work fuzz targets, schema/hash/CLI/conformance tests, release-policy and packaged workflow tests; Site unit/browser/generation/build/artifact/dry-run.

**Target Platform**: Windows/macOS/Linux amd64 native hosted proof; six pure-Go amd64/arm64 build/package structural targets.

**Project Type**: Offline CLI plus peer JSON Schema, non-publishing release tooling and documentation Site preview.

**Performance Goals**: Preserve existing bounded parser/model/render/capture/diagnostic limits and deterministic no-loss/no-publication boundaries; no new runtime or performance target.

**Constraints**: No original paths/local identifiers in Cue JSON; no external native resource execution; exact historical bytes; no tag/release/production/final merge; hidden Windows noninteractive subprocesses.

**Scale/Scope**: Issues #64/#65, four native formats, twelve distinct conversions, 17 required scripted evidence rows, six package targets, two retained public schema routes.

## Constitution Check

| Principle | Pre-design and post-design assessment |
|---|---|
| I Source truth | PASS: Original source fixture bytes/provenance and restore hashes unchanged; models retain native owners. |
| II Peer products | PASS: exact current 1.1.0 canonical/embedded/emitted/immutable/package equality plus immutable historical digests. |
| III Common/native model | PASS: bounded existing ownership/profile is frozen without semantic or external execution expansion. |
| IV No silent loss | PASS: twelve-edge complete loss/default/precision/strict/fatal proof retained. |
| V Test-first | PASS: focused stable identity/capability/refusal/matrix/package tests precede promotion; existing corpus/race/fuzz retained. |
| VI Portable/security | PASS: six pure-Go targets, three hosted native package smokes, hidden Windows process guarantee, offline exact registry and bounded integrity. |
| VII Delivery authority | PASS: installed Spec Kit blocking analysis and convergence, explicit push/official PR authority; human final merge; #66/#67 exclusions. |

No constitutional exception is required. An unavoidable established v1 public break is a blocking escalation, not an identity-relaxed workaround.

## Project Structure

### Documentation (this feature)

```text
specs/S028-freeze-v1-1-release-candidate/
  spec.md
  plan.md
  research.md
  data-model.md
  contracts/candidate.md
  contracts/release-evidence-contract.schema.json
  quickstart.md
  tasks.md
  checklists/requirements.md
  checklists/candidate-integrity.md
  verification.md
```

### Source Code (repository root)

```text
internal/version, internal/schema, internal/model, internal/cli, internal/convert
internal/conformance
schema/releases/v1.1.0
scripts/release-verify
scripts/docs-verify
.github/workflows/release-proof.yml
.goreleaser.yaml
docs, README.md, CHANGELOG.md
site/content-map.json, site/scripts/generate.mjs, site/tests
testdata/conformance-matrix.json
```

**Structure Decision**: Reuse exact current/historical registry and existing stable verifier rather than inventing another identity or proof mode. New workflow proof stays in the standalone release verifier with platform-hidden process control; no production package dependency imports CLI.

## Phase 0 research decisions

- Current identity promotion changes only current runtime/examples/native report identity; immutable schema/releases v0/v1, bundled historical resources and historical CLI goldens remain byte-identical. Earlier slice snapshots remain chronological observations.
- Accept schema_only and complete experimental declarations as current input observations, adding complete stable as another precise tuple; official encode/private native targets use stable. This avoids treating installed capability as loaded source truth while still advertising the executed stable profile.
- Replace the two scripted-stable-gate deferrals with stable identity/render/platform evidence and tighten scripted matrix tests to forbid any scripted deferral. Ordinary corpus/platform tests remain reused.
- Stable candidate packaging switches existing schema source to the immutable v1.1.0 file and uses stable verifier mode without -development, preserving read-only, same-bundle hosted smokes and nonpublication.
- Packaged smoke must prove all native commands, loss/refusal, exact restoration/provenance and historical inputs. Old consumer proof uses verified public v1.0.0 bytes bound to known release revision/checksum, never unverified input-selected executables.
- Published v1.0.0 inspect intentionally suppresses detailed schema diagnostics. Immediately run that binary's validate identity gate against the same unchanged fresh payload before each owning old-command refusal; require the exact historical schema identity diagnostic, runtime failure, no stdout, unchanged payload, and forced/absent destination safety. Record 32 identity probes separately from 32 owning refusal paths rather than attributing suppressed details to inspect.
- Add ASS/SSA authored preview route/navigation and four-format conversion description while retaining latest published v1.0.0 download/schema metadata and exact two public schema routes. Candidate notes remain repository publication-ready.
- Existing Windows launcher is reused with CREATE_NO_WINDOW and redirected noninteractive I/O. Direct git/gh remain allowed.
- Installed checklist prerequisites require plan.md, so setup-plan scaffolding precedes the domain checklist while substantive plan design follows. The first clarify prerequisite failed on an incorrect local feature.json key and was corrected to the installed feature_directory key before rerun.
- Custom checklist markers stay reviewer-owned; all ten requirements are audited present. Explicit unattended autopilot authorizes routine continuation; implementation never changes those markers. No extension hooks are installed.

## Ordered implementation ownership

1. Coordinator owns Spec Kit artifacts, matrix/stable cross-package conformance, chronological changelog/roadmap integration, GitHub lifecycle and full final verification.
2. Independent candidate owner owns current runtime/schema/model/CLI/convert promotion and focused tests, immutable current schema, release verifier/config/workflow and package compatibility proof.
3. Independent documentation owner owns maintained root docs other than changelog/roadmap/main plan, docs verifier, publication-ready release notes, and generated Site map/navigation/tests. Coordinate exact final schema annotation prose with candidate owner.
4. Coordinator integrates only after blocking analysis and owner-focused tests, then independent convergence audits cross-owner output before full foreground verification.

Independent file ownership is explicit; all integrated verification remains coordinator-owned.

## Verification strategy

Run focused identity/capability, historical/native CLI and matrix tests after test-first promotion, standalone release tests/policy, and docs/Site generator regressions. Then run full root and six nested CI-parity quality/fuzz/format/docs/brand gates, six pure-Go targets, and complete Site suite. Commit locally before exact candidate packaging because clean VCS source identity is required. Local candidate packaging and native execution run on Windows through a hidden noninteractive launcher; three hosted native executions and exact-head security/CI/package proof must pass after publication. No foreign native execution is claimed.

Publish bodies from UTF-8 formatted files and immediately verify exact body/rendering/closing references. First Codex review is automatic; address all findings and use exactly one second @codex review only if round one finds concerns. Record terminal external evidence as a PR completion comment to avoid an evidence-only new head. Default Status remains unused; #64 acceptance precedes #65; both stay open until human merge.
