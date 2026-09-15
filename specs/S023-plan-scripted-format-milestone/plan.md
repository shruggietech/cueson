# Implementation Plan: Plan the scripted-format milestone

**Branch**: `codex/S023-plan-scripted-format-milestone` | **Date**: 2026-09-15 | **Spec**: [spec.md](spec.md)

**Input**: Operator-approved S023 planning/backlog scope.

## Summary

Deliver three complete planning outcomes in one PR: governed v1.1.0 roadmap/backlog, bounded ASS v4+/SSA v4 native fidelity/rendering contract, and exact historical v1 input compatibility contract. Publish one epic and sixteen children with native relationships and Project metadata; reconcile maintained current-state prose. Later codecs/schema/releases remain separate slices and protected authority gates.

## Technical Context

**Language/Version**: UTF-8 Markdown and GitHub native metadata; existing Go 1.25+ and Node/pnpm toolchains are verification dependencies only.

**Primary Dependencies**: Installed hyphenated Spec Kit skills/scripts, existing GitHub formatter/docs verifier, GitHub REST issue/milestone/sub-issue/dependency APIs and Project v2; maintained Aegisub/libass primary grammar references.

**Storage**: Maintained docs and Spec Kit artifacts; GitHub owns current planning facts. Traceability snapshots are read-back evidence, not competing planning authority.

**Testing**: Existing publication/docs tests, native planning graph/cardinality/body audit, full unchanged root tests and complete site generation/lint/unit/build/browser/artifact/dry-run suite. Hosted CI/security/release proof remains required.

**Target Platform**: Windows development with verified hidden process launches; native hosted Linux/macOS/Windows and six pure-Go build targets retain established proof.

**Project Type**: Documentation/contracts and external planning-state delivery for an existing CLI/schema product.

**Performance Goals**: No runtime performance change. Exactly seventeen new planning items, sixteen native children and eighteen direct acyclic blocker edges; partial publication is inspected and resumed without duplicates.

**Constraints**: No runtime implementation, version bump, immutable artifact mutation, arbitrary asset fetch/execution, production deployment, release publication or final merge. Preserve unknown native/source information or reject unsafe whole operations explicitly.

**Scale/Scope**: Twenty-five requirements, three stories/outcomes, sixteen future child issues, one epic and eight chronological proposed execution slices.

## Constitution Check

Pre-research and post-design gates: PASS, without exception.

| Principle | Applied design and evidence |
|---|---|
| I source preservation | Source bytes stay authoritative; accepted native unknown content retained, unsafe path/reference cases reject without sanitization. |
| II schema/software | No version/schema change now; current lockstep and immutable historical bytes remain intact. |
| III common/native | Native/editable structures and derived dialogue semantics are distinct; no opaque-only complete support claim. |
| IV no silent loss | Explicit malformed/unsupported rows, Timer restriction and complete future conversion loss accounting. |
| V test-first | Existing applicable checks and future corpus/codec tests are acceptance owners; no mirrored text-only tests added. |
| VI portable/security | Untrusted input/reference/attachment limits are specified; Windows launchers hidden and no trust expansion. |
| VII delivery authority | Native issue/project facts, blocking analysis, existing authorized push/PR, human merge, separate release/production authorization. |

## Project Structure

### Documentation (this feature)

```text
specs/S023-plan-scripted-format-milestone/
  spec.md
  plan.md
  research.md
  data-model.md
  contracts/delivery.md
  contracts/version-compatibility.md
  quickstart.md
  tasks.md
  checklists/requirements.md
  checklists/contracts.md
  backlog/*.md
  issue-map.json
  verification.md
```

### Maintained source surfaces

```text
docs/roadmap.md
docs/formats/ass-ssa.md
docs/architecture.md
docs/compatibility.md
docs/schema.md
docs/project-management.md
docs/Cueson-Project-Specification-v0.0.0.md
README.md
CHANGELOG.md
```

**Structure Decision**: Root docs remain authored authority. Existing site generation consumes root sources; no new production route or authored site copy is introduced. The ASS/SSA page is a future contract, not a claim of installed support. Backlog body snapshots trace published acceptance intent; native metadata is audited independently.

## Phase 0 research

Independent native-contract and version-compatibility researchers were dispatched as required by the installed plan skill. Root owns GitHub API research, roadmap, integration and final verification. Decisions/options are consolidated in [research.md](research.md). No new irreversible architecture is needed: preserve current public boundaries and use explicit bounded future acceptance.

## Phase 1 design

[Data model](data-model.md) defines atomic outcomes, graph and Project lifecycle. [Delivery contract](contracts/delivery.md) defines publication/read-back/recovery. [Version contract](contracts/version-compatibility.md) specifies exact identities and command expectations. The maintained ASS/SSA contract defines native shapes and source-linked acceptance rows. [Quickstart](quickstart.md) gives reproducible local/native review checks.

## Implementation sequencing

After task generation and read-only blocking analysis, root publishes/reuses native backlog state, resolves P aliases to real issue links, integrates both researchers' contract findings, reconciles maintained prose and verifies all body/native/Project facts. S023 outcomes stay open in PR review until human merge; future issues remain Backlog until hard blockers clear.

Existing formatter/docs/Go/site checks run in the foreground through a previously verified CREATE_NO_WINDOW launcher. No test logic, CI permissions, pinned dependencies, current schema or runtime code is changed. Existing hosted security, race, cross-platform and packaged proof supplies the remaining platform evidence.

## Decision log and deviations

- Adopt ASS/SSA-first as the operator-approved recommendation, rather than infer ordering from the draft's section order. Keep every other family explicitly on the later roadmap.
- Define Timer=100 bounded acceptance and safe whole-operation rejection for incompatible metadata/reference cases; never silently ignore timing semantics or sanitize raw fields.
- Keep 0.0.0 input unsupported in the next executable, matching the actual current v1 baseline; its historical binary remains available. Adding adapters later needs its own outcome.
- Preserve exact 1.0.0 input support on the future executable, while new documents have exact 1.1.0 identity and old consumers explicitly reject that unknown version.
- Reviewer-owned custom checklist markers remain unchecked; an independent assessment and clean analysis establish requirement quality under existing autopilot authorization. This intentionally avoids treating them as implementation tasks.
- Reconcile present-tense delivery prose, preserving dated historical records and actual release/deployment revision identity.
- Kickoff explicitly authorizes push and official PR creation, so the autopilot pre-push approval requirement is already satisfied. Tag/release, production and human merge authority are not inferred.

## Complexity Tracking

No constitutional violation or new runtime abstraction is introduced. Native state is represented by existing GitHub mechanisms; existing verification is reused.
