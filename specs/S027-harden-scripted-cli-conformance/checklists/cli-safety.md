# CLI and source-safety requirements checklist: S027

**Purpose**: Review requirement completeness and consistency before implementation.

**Created**: 2026-09-15

**Feature**: [spec.md](../spec.md)

**Review Ownership**: Reviewer-owned requirements-quality artifact. Markers remain unchecked until reviewer evaluation; they do not record implementation completion.

## Completeness and clarity

- [ ] CHK001 Are native, Cue JSON, historical identity, alias/content selection, and installed/declared capability outcomes specified? [FR-001]
- [ ] CHK002 Are scripted inspection counts defined separately from existing text report fields, with prohibited source values explicit? [FR-002, assumptions]
- [ ] CHK003 Are help/four completion vocabulary and established owning-command streams/error classes covered without a public reclassification? [FR-003, clarifications]
- [ ] CHK004 Are accepted exact restoration, native/common edit cycles, and twelve-direction conversion regressions required? [FR-004]
- [ ] CHK005 Are malformed ownership, corrupt envelopes, unsafe references, privacy, strict, and partial-output refusals specified? [FR-005]
- [ ] CHK006 Are complete structures/diagnostics/loss limits and cancellation requirements explicit? [FR-006]

## Evidence and delivery boundaries

- [ ] CHK007 Are meaningful filesystem-free fuzz mutation boundaries and bounded accepted controls required? [FR-007]
- [ ] CHK008 Are development conformance, portable inapplicability, and future stable/release gates distinguished? [FR-008]
- [ ] CHK009 Are three native platforms, six pure-Go builds, race/static/security, site/docs, and non-publishing package proof included? [FR-009]
- [ ] CHK010 Are blocking analysis/convergence, exact-head CI, every review finding, two-round cap, and human final merge required? [FR-010]

## Notes

Autopilot requirements audit found all ten criteria specified. This command leaves reviewer-owned markers untouched. The operator's full unattended kickoff authorizes continuation after this documented audit; implementation reads these markers without changing them.
