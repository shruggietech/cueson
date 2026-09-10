# Data Model: v0.0.0 Documentation and Milestone Verification

S010 persists no application data. These entities define the repository documentation and delivery evidence inspected by the verifier and operator.

## Documentation Contract

Fields:

- `path`: canonical repository-relative path.
- `required`: whether absence is always a verification failure.
- `authority`: the subject for which this document is current source of truth.
- `current_claims`: implemented or externally verified statements that require evidence.
- `future_claims`: planned statements that require an explicit future or gated label.
- `links`: local and external references discovered in the document.

Validation rules:

- Every required path is unique and exists as a regular file beneath the repository.
- Current claims agree with implementation or read-back evidence.
- Future claims do not read as current capability.
- Repository-authored text follows the established encoding and Markdown source rules.

## Documentation Link

Fields:

- `source_document`: document containing the reference.
- `source_line`: one-based source line for diagnostics.
- `raw_target`: destination exactly as authored.
- `kind`: Markdown link, Markdown image, reference definition, HTML link, or HTML image.
- `classification`: local file, local fragment, repository-root target, or external target.
- `resolved_path`: normalized repository path for a local target.
- `resolved_fragment`: decoded heading identity for a fragment target.

Validation rules:

- Local resolution never escapes the repository root.
- A local path exists as a regular file or directory.
- A fragment names a generated heading anchor in the target Markdown document.
- External targets have a recognized scheme or network-path form and are not fetched.
- Code spans and fenced code do not create link records.

## Heading Anchor

Fields:

- `document`: Markdown document containing the heading.
- `line`: one-based heading line.
- `heading_text`: visible heading content after Markdown formatting is removed.
- `slug`: GitHub-compatible lowercase anchor identity.
- `occurrence`: zero-based duplicate count used for `-1`, `-2`, and later suffixes.

Validation rules:

- Headings inside fenced code are ignored.
- Duplicate base slugs receive stable suffixes in document order.
- URL-encoded link fragments are decoded before comparison.

## Foundation Gate

Fields:

- `criterion`: one v0.0.0 completion statement.
- `state`: satisfied, merge-pending, release-gated, or unsatisfied.
- `evidence`: repository file, test, workflow run, issue, Project field, or release boundary.

State transitions:

- `unsatisfied` to `satisfied` only after evidence exists.
- `merge-pending` to `satisfied` only after the official S010 change reaches `main`.
- `release-gated` remains incomplete in S010.

## Milestone Evidence Record

Fields:

- `epic`: issue #2.
- `children`: issues #1 through #12.
- `milestone`: v0.0.0.
- `project`: `cueson Delivery`.
- `version_evidence`: software/schema lockstep results.
- `quality_evidence`: local and hosted checks.
- `release_proof`: S009 non-publishing result.
- `protected_actions`: tag, release, immutable schema copy, milestone closure, and production publication.

Validation rules:

- Every child retains independent closure evidence.
- Issue #12 and epic #2 remain open until S010 merges.
- All in-scope issues appear exactly once in the Project.
- Stage and Slice carry workflow state; default Status is empty.
- Protected actions remain absent unless separately authorized.
