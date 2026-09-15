# Conversion Requirements Checklist: Complete scripted conversion

**Purpose**: Reviewer-owned requirements quality gate for the complete conversion matrix and refusal safety.
**Created**: 2026-09-15
**Feature**: [spec.md](../spec.md)

**Note**: Generated through the installed speckit-checklist workflow.
**Review Ownership**: Mark an item checked only after the reviewer determines that requirements quality is satisfied.
**Marker Semantics**: A checked item means requirements approval, not implementation completion.

## Requirement Completeness

- [x] CHK001 Are all ten new directions and the two established directions explicitly enumerated? [Completeness, Spec FR-001/002/003]
- [x] CHK002 Are atomic loss identity, ordering, and source-reference requirements defined? [Completeness, Spec FR-004]
- [x] CHK003 Are drawing-only and empty-script outcomes explicitly defined? [Clarity, Spec FR-005 and Edge Cases]
- [x] CHK004 Are deterministic defaults distinguished from observed source provenance? [Clarity, Spec FR-006]

## Requirement Clarity and Consistency

- [x] CHK005 Are quantization, ties, overflow, and collapsed interval outcomes explicit? [Clarity, Spec FR-007 and Edge Cases]
- [x] CHK006 Are variant retention and every known omitted native category accounted for? [Completeness, Spec FR-008]
- [x] CHK007 Are strict losses and fatal refusal requirements consistent across file and stdout publication? [Consistency, Spec FR-009/010]
- [x] CHK008 Are literal control sequences and target semantic reparsing boundaries defined? [Coverage, Spec FR-011 and Edge Cases]

## Safety, Compatibility, and Scope

- [x] CHK009 Are source preservation, privacy, cancellation, and allocation/report bounds explicit? [Non-Functional, Spec FR-012/013/014]
- [x] CHK010 Are historical inputs, immutable schemas, and existing loss meanings protected? [Consistency, Spec FR-003/015/019]
- [x] CHK011 Are necessary conversion CLI changes distinguished from later release-wide scope? [Scope, Spec FR-016 and Scope and authority]
- [x] CHK012 Are twelve-direction evidence and protected delivery boundaries measurable and complete? [Measurability, Spec FR-017/018/020 and SC-001 through SC-005]

## Notes

New items intentionally remain unchecked until the coordinating reviewer assesses the completed research and design. speckit-implement reads markers and does not change them. Scope: conversion fidelity and safety; depth: standard; audience: coordinating reviewer before implementation. The plan template was initialized before checklist prerequisite resolution because the installed checker requires plan.md; design authoring follows research.

Coordinating requirements review (2026-09-15): 12/12 approved after resolved research, clarification, complete design, and blocking analysis. No unresolved requirements gap or constitution conflict. This assessment is separate from speckit-implement and precedes implementation.
