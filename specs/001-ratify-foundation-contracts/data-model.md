# Data Model: Ratify Foundation Contracts

## Contract Baseline

Represents the reviewed implementation authority for a topic.

**Fields**:

- `topic`: architecture, schema, or CLI.
- `version`: `0.0.0` for this slice.
- `status`: ratified for implementation, while still subject to explicit future Spec Kit amendments before release.
- `authority`: constitution first, then the ratified topic document, then the working project specification for non-conflicting detail.
- `decisions`: normative choices owned by the topic.
- `deferred_work`: behavior intentionally assigned to later issues.
- `verification`: observable checks that prove internal consistency.

**Relationships**:

- One contract baseline traces to one or more constitutional principles.
- One implementation issue may depend on multiple contract baselines.
- The working project specification links to all three baselines.

## Format Capability Declaration

Represents official release capability for one schema-recognized format.

**Fields**:

- `format`: canonical lowercase identifier such as `subrip` or `webvtt`.
- `status`: maturity value; `envelope_only` for the completed v0.0.0 milestone.
- `ingest_supported`: whether native source ingest exists.
- `render_supported`: whether model-driven output exists.
- `restore_supported`: whether valid preserved assets can be restored exactly through the public CLI.
- `ocr_required_for_semantic_output`: whether semantic text depends on OCR.

**Validation rules**:

- Capability describes the official Cueson release associated with the contract.
- Schema recognition never implies native ingest or render support.
- Exact restoration depends on a valid source envelope and public restore command, not on a codec.

## Command Availability Declaration

Represents one CLI command's relationship to the current executable.

**Fields**:

- `command`: canonical command name.
- `availability`: implemented or absent.
- `listed_in_help`: true only when implemented.
- `payload_stream`: stdout when the command emits a payload.
- `diagnostic_stream`: stderr.
- `failure_class`: invocation/precondition or runtime capability.
- `exit_code`: 2 for invocation/precondition failure, 1 for runtime failure, 0 for success.

**State transition**:

- Absent commands become implemented only in their owning issue and then enter help and command listings.

## Restoration Output Plan

Represents the validated one-to-one mapping from source assets to runtime destinations before restoration performs file I/O.

**Fields**:

- `asset_id`: document-local source asset identity.
- `stored_basename`: the schema-validated original basename.
- `portable_collision_key`: basename transformed through Unicode canonical caseless matching (NFD normalization, default Unicode case folding, then NFD normalization again).
- `destination`: the complete runtime output path.
- `destination_origin`: stored default, stored basename beneath `--output-dir`, or explicit single-asset `--output`.

**Validation rules**:

- Stored basenames and portable collision keys are unique within the source bundle.
- Default output is valid only for one asset and uses its stored basename in the current directory.
- `--output-dir` maps every asset's stored basename beneath one caller-selected directory.
- `--output` is valid only for one asset, may rename it, and is validated as a caller-selected literal runtime destination rather than as stored document data.
- `--output` and `--output-dir` are mutually exclusive.
- The complete plan rejects duplicate destinations and overwrite violations before any output is opened; `--force` does not permit intra-bundle replacement.

## Delivery Record

Represents the evidence-backed state of an issue within the active Project.

**Fields**:

- `issue_number`: GitHub issue identity.
- `slice`: coherent work-slice identifier.
- `stage`: governed Project Stage value.
- `state`: open or closed.
- `evidence`: commit, checks, and readback supporting the state.

**Validation rules**:

- Closed issues use Stage Done.
- Merge-dependent criteria remain unchecked while the change is absent from the default branch.
- Local-only work does not claim a public branch or pull request.
- Issue #3 remains open until its change reaches the default branch through the approved merge process.
