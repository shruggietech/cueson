# Data Model: Import the Official Cueson Brand Kit

S011 adds no application data. These entities define the retained acquisition and its deterministic repository evidence.

## Import Manifest

Fields:

- `schema_version`: manifest contract version, initially `1`.
- `brand`: canonical brand identifier `cueson`.
- `kit_version`: published kit version `1.0.0`.
- `source_url`: operator-designated acquisition URL.
- `effective_url`: final URL observed after redirects.
- `acquired_at`: UTC RFC 3339 timestamp for the pinned acquisition.
- `archive`: archive repository path, exact byte size, lowercase SHA-256, entry count, and total uncompressed bytes.
- `payload_root`: repository path containing the exact extraction.
- `entries`: path-sorted complete collection of payload entry records.
- `references`: path-sorted repository documentation references to retained assets.

Validation rules:

- The document uses strict JSON with no unknown or duplicate fields.
- Paths use forward slashes, are repository relative, and remain beneath the declared roots.
- Entry paths are unique by exact spelling and Unicode case folding.
- Entry records are sorted by path and exactly match the archive catalog.
- Recorded counts and totals equal calculated values.
- Digests contain exactly 64 lowercase hexadecimal characters.

## Archive Identity

Fields:

- `path`: stable repository-relative ZIP path.
- `bytes`: exact archive length.
- `sha256`: exact archive digest.
- `entry_count`: accepted regular-file entry count.
- `uncompressed_bytes`: sum of accepted entry sizes.

Validation rules:

- The archive is a regular file and not a symbolic link.
- Calculated size and digest match the manifest.
- Every ZIP entry is a regular file with a safe, unambiguous portable path.
- Bounded entry and total sizes prevent decompression abuse.

## Payload Entry

Fields:

- `path`: archive-relative and payload-root-relative path.
- `bytes`: exact uncompressed length.
- `sha256`: digest of uncompressed bytes.

Validation rules:

- The archive stream and retained file both match the same size and digest.
- The payload tree contains exactly the manifest entry set and no additional directories with data, symlinks, or non-regular nodes.
- Name spelling and case match exactly.

## Repository Reference

Fields:

- `document`: repository-authored README or HTML document.
- `asset`: retained payload entry path.
- `reference`: exact repository-relative or document-relative text expected in the document.
- `purpose`: human-readable rendering role.

Validation rules:

- The document and asset are regular, non-symlinked files beneath the repository.
- The asset exists in the verified entry set.
- The exact reference text appears in the declared document.
- Records are unique and sorted by document, then asset, then reference.

## Import Evidence

Fields:

- `local_verification`: formatter, documentation, brand, product, workflow, and whitespace results.
- `hosted_verification`: current-head CI, CodeQL, policy, and security results.
- `review_state`: first-round result, optional second-round result, and resolved threads.
- `github_state`: issue #23, Delivery Project Stage, Slice, Status, and pull-request closing reference.

State transitions:

- Issue Stage moves from `Backlog` to `In progress` at kickoff and to `PR review` after publication.
- A second Codex round is requested only if round one reports findings after remediation.
- Issue #23 remains open until the operator merges the official pull request.
- S011 stops at green and fully reviewed; merge and release actions remain absent.
