# Conformance Matrix Contract

## Matrix identity

The maintained machine-readable matrix is `testdata/conformance-matrix.json`. It is repository-authored metadata, not a source payload governed by the fixture manifest.

Each row contains a stable `row_id`, one format, a concise behavior description, and evidence entries. Row IDs appear verbatim in the appropriate maintained format guide.

## Evidence classes

- `accepted_fixture`: A governed accepted fixture that exercises the behavior.
- `malformed_fixture`: A governed deterministic rejection fixture.
- `render_test`: Focused canonical renderer or parser-renderer-cycle evidence.
- `conversion_test`: Focused loss-free, lossy, strict, or fatal conversion evidence.
- `fuzz_target`: A named bounded fuzz surface.
- `platform_test`: A native operating-system or static target-build claim.
- `inapplicable`: A required class does not apply, with a non-empty reason.

## Validation

- Matrix and row versions are supported.
- Format and row identifiers use canonical portable names and are unique.
- Every fixture ID resolves exactly once in `testdata/manifest.json`.
- Every test or fuzz reference names an existing repository test function.
- Every required format-guide row exists in the matrix and vice versa.
- An evidence class cannot be omitted silently; inapplicability is explicit.
- Values contain no absolute path or local machine identifier.
