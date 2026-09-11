# Schema Annotation Contract

## Coverage

- Every Cueson-owned property at the root has a specific `description` and representative `examples` where the property can stand alone meaningfully.
- Every direct property of every locally reachable public `$defs` node has a specific `description`, including property nodes that delegate validation through `$ref`.
- Public object, union, and constrained-scalar definitions have a `title` when it materially improves generated documentation.
- Every enum site and each non-obvious constrained-value category has at least one valid example.

## Semantics

- `title`, `description`, and `examples` are non-normative annotations.
- Validation keywords, conditional branches, model invariants, and source-integrity checks remain authoritative.
- Descriptions identify units, null meaning, provenance, source authority, native-format boundaries, and capability meaning where applicable.
- Producer software versions remain independent from the Cue JSON schema version.

## Example verification

- Every example is validated against the exact compiled fragment on which it appears.
- Complete root examples are additionally validated through the full public schema and semantic decoder.
- Tests reject stale schema identity or version examples, unsafe basenames, absolute paths, machine identifiers, blank descriptions, missing enum examples, and type-invalid examples.
- Schema emission remains byte-identical to the canonical annotated file.
- `schema/releases/v0.0.0/cueson.schema.json` remains unchanged with SHA-256 `d15c7fa5227156109dd6be3d39b711aca3503794bb862169dfca96ee80adb975`.
