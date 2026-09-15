# Requirements Quality Checklist: Scripted milestone contracts

**Purpose**: Reviewer assessment of requirement clarity, coverage and consistency.

**Created**: 2026-09-15

**Feature**: [spec.md](../spec.md)

**Ownership**: Reviewer-owned. A checked marker means requirements-quality approval, not completed implementation. Autopilot may record assessment separately without modifying these markers.

## Completeness and Clarity

- [ ] CHK001 Are sixteen atomic outcomes and three S023 ownership boundaries explicit? [Completeness, Spec FR-002]
- [ ] CHK002 Are native relationships distinct from coordination and Project metadata? [Clarity, Spec FR-004 through FR-006]
- [ ] CHK003 Are dialect, field declarations, overrides, karaoke and drawings covered by source-linked acceptance rows? [Coverage, Spec FR-012 through FR-016]
- [ ] CHK004 Are rendering edit ownership and failure conditions explicit? [Clarity, Spec FR-014]

## Fidelity, Security and Compatibility

- [ ] CHK005 Are source bytes, native/raw fields and common derived observations distinguished consistently? [Consistency, Spec FR-013]
- [ ] CHK006 Are reference/attachment/path cases reconciled with source preservation and privacy without silent sanitization? [Conflict, Spec FR-015]
- [ ] CHK007 Are exact historical identities, commands, current output and old-consumer behavior defined? [Completeness, Spec FR-017 through FR-019]
- [ ] CHK008 Are unknown/mismatched identity and provisional capability states unambiguous? [Clarity, Spec FR-017, FR-020]

## Acceptance and Delivery

- [ ] CHK009 Are twelve conversion directions, timing precision, complete losses and strict/fatal rejection specified? [Coverage, Spec FR-021]
- [ ] CHK010 Are partial publication recovery and new-issue reconciliation requirements measurable? [Measurability, Spec FR-007, FR-008]
- [ ] CHK011 Are future XML, SAMI, broadcast and bitmap/OCR outcomes retained without hidden v1.1 scope? [Consistency, Spec FR-010]
- [ ] CHK012 Are protected release/production transitions and final human merge distinct from S023 completion? [Clarity, Spec FR-024, FR-025]

## Notes

The blocking analysis and separate assessment record evaluate these requirements. Implement reads these markers and does not turn them into an implementation task list.
