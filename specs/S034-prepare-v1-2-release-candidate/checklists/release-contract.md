# Release contract requirements checklist

**Purpose**: Reviewer-owned clarity check for candidate identity, compatibility, proof and publication authority.

**Feature**: [spec.md](../spec.md)

## Contract clarity

- [x] CHK001 Are candidate identity, latest public version and production hosting state distinguished? [FR-001, FR-010, FR-014]
- [x] CHK002 Are historical support and former-development refusal stated without forward-compatibility claims? [FR-003, FR-006]
- [x] CHK003 Do packaged proof requirements include consumer fields and all native loss paths? [FR-004]
- [x] CHK004 Are exact bundle/source/schema bindings and same-bundle host requirements quantified? [FR-002, FR-007, FR-008]
- [x] CHK005 Is the published old-consumer provenance independently authenticated? [FR-006]
- [x] CHK006 Are stale evidence, partial upload and changed-state recovery defined? [FR-011]
- [x] CHK007 Are planning/project facts owned by native metadata? [FR-012]
- [x] CHK008 Are review limits, executed evidence and action-specific authority explicit? [FR-013, FR-014]

## Ownership

Unchecked markers remain reviewer-owned requirements-quality criteria. Autopilot's explicit instruction to continue authorizes implementation after the independent analysis assesses their substance; these markers do not represent implementation or test completion.

## Independent review notes

Reviewed on 2026-10-06 by the independent model/documentation reviewer. All eight criteria pass as requirements-clarity checks. These marks do not certify implementation completion, hosted native execution, public availability or authorization for a later action. The separate implementation review found four proof gaps; subsequent source review confirmed their remediation through historical strict destination refusals, minimal selector probes, independent declared-media rejection and validate/inspect output privacy guards. Actual public-consumer execution and all three hosted same-bundle native proofs remain separately reported verification gates.

- CHK001: [spec.md](../spec.md) separates exact 1.2.0 candidate identity from latest published v1.1.0; FR-014 excludes production route, download and deployment changes. [publication-decision.md](../contracts/publication-decision.md) assigns preparation, publication and hosting to #84, #85 and #86.
- CHK002: FR-003 and the edge cases require exact historical 1.0.0/1.1.0 support, immutable released contracts and former-development refusal. Acceptance scenarios require older consumers to reject new output rather than claim forward compatibility.
- CHK003: FR-004, [candidate-contract.md](../contracts/candidate-contract.md) and the plan require consumer storage, inert identifiers, cue/media timing, explicit unavailable state, exact restoration, deterministic omission and strict refusal. The publication decision quantifies four native formats and twelve conversion directions, retaining S033's consumer rules.
- CHK004: FR-002/007/008 and the candidate/publication contracts specify six archives, six target-bound SBOMs, one checksum manifest, exact source/schema binding and three amd64 hosts executing the same accepted bundle without rebuilding. Arm64 proof is explicitly structural.
- CHK005: FR-006 and [published-v1.1.0-consumer.json](../contracts/published-v1.1.0-consumer.json) bind public archive names, lengths, digests, source revision, released schema and legal bytes. The publication decision requires genuine positive controls and refusal/probe counts for separately identified published 1.0.0 and 1.1.0 binaries; source rebuilds cannot substitute.
- CHK006: The publication decision's refresh and partial-failure requirements invalidate expired or changed evidence, prohibit silent artifact substitution and require retained partial state plus operator recovery judgment when read-back or publication differs.
- CHK007: FR-012 and [issue-map.json](../issue-map.json) identify native GitHub issues, dependencies, milestone and Project fields as current authority. The issue map explicitly identifies itself as a dated traceability snapshot.
- CHK008: FR-013/014, SC-005 and the plan require actual executed checks and review evidence, at most two automated Codex rounds, and separate authority for final specific-PR merge, tag creation/push, release/assets publication and production hosting. [tasks.md](../tasks.md) distinguishes local completion from final-head hosted/review runtime gates.
