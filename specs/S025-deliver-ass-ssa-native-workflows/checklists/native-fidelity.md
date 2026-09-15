# Native fidelity requirements checklist

**Purpose**: Reviewer requirements-quality assessment for S025.

**Created**: 2026-09-15

**Feature**: [spec.md](../spec.md)

Markers belong to the reviewer: checked means requirements-quality approval, never implementation completion. Autopilot evaluates the requirements without changing these generated markers.

## Completeness and clarity

- [ ] CHK001 Are corpus inventory and independent expected semantics explicitly required? [Completeness, Spec FR-001/FR-002]
- [ ] CHK002 Are all dialect selection and rejection/fallback cases defined? [Coverage, Spec FR-003]
- [ ] CHK003 Are exact source capture and canonical output distinguished? [Clarity, Spec FR-004/FR-010/FR-013]
- [ ] CHK004 Are native ordering, unknown fields and attachment retention specified? [Completeness, Spec FR-005]
- [ ] CHK005 Are drawing exclusion, actor provenance and eligible karaoke criteria explicit? [Clarity, Spec FR-006]

## Safety and edit ownership

- [ ] CHK006 Is privacy classification grammar-dependent for original and edited views? [Consistency, Spec FR-008]
- [ ] CHK007 Are deterministic diagnostics and no-partial-publication obligations complete? [Coverage, Spec FR-009/FR-015]
- [ ] CHK008 Are captured observations distinct from structured edit owners? [Clarity, Spec FR-011/FR-012]
- [ ] CHK009 Are constructed content and declaration transitions covered? [Coverage, Spec FR-012; Clarifications]
- [ ] CHK010 Are precision, line, total output and aggregate ceilings quantified by the ratified profile? [Measurability, Spec FR-013/FR-016]

## Boundaries

- [ ] CHK011 Are experimental versus schema-only/stable capability declarations consistent? [Consistency, Spec FR-007/FR-017; Clarifications]
- [ ] CHK012 Are historical compatibility, later conversion and protected publication boundaries clear? [Coverage, Spec Scope/FR-018/FR-020]

## Notes

Autopilot assessment: all twelve criteria are clearly specified, with no unresolved requirement-quality gap. Reviewer markers remain unchanged.
