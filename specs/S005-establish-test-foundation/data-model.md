# Data Model: Fixture and Conformance-Test Infrastructure

## Fixture Manifest

The fixture manifest is an ordered versioned collection. Version 1 has exactly `manifest_version` and `fixtures`; unknown fields fail validation so governance changes remain explicit.

| Field | Type | Rules |
|---|---|---|
| ManifestVersion | integer | Exactly 1 for S005. |
| Fixtures | ordered fixture records | Non-empty; fixture identifiers are unique under portable canonical caseless matching. |

## Fixture Record

| Field | Type | Rules |
|---|---|---|
| ID | string | Stable slash-separated lowercase identifier; unique and portable. |
| Purpose | string | Non-empty human explanation of the behavior evidenced. |
| Class | enum | `accepted`, `malformed`, or `fuzz_regression`. |
| Origin | provenance origin | Explicit kind, source, and reproducible recipe where applicable. |
| Redistribution | redistribution decision | `approved` is required for committed material. |
| Artifacts | ordered artifact records | Non-empty; logical IDs and portable paths are unique. |
| Expectation | case expectation | Declares acceptance or one stable rejection stage and diagnostic fragment. |

## Provenance Origin

| Field | Type | Rules |
|---|---|---|
| Kind | enum | `project_authored`, `synthetic`, `derived`, or `third_party`. |
| Source | string | Stable origin description or reference, never a local path. |
| Recipe | string | Required for synthetic or derived material; omitted only when not applicable. |
| RetrievedOn | date | Required for third-party material. |

## Redistribution Decision

| Field | Type | Rules |
|---|---|---|
| Status | enum | Must be `approved`; `unknown` and `denied` are rejected. |
| License | string | SPDX identifier or repository-relative bundled license reference. |
| Attribution | string | Explicit copyright or attribution statement. |
| NoticeRequired | boolean | Drives the repository NOTICE review obligation. |

## Artifact Record

| Field | Type | Rules |
|---|---|---|
| ID | string | Stable logical identifier unique within the fixture. |
| Role | enum | `source`, `expected_model`, `expected_diagnostics`, `expected_bytes`, or `malformed_input`. |
| Path | string | Slash-relative beneath testdata; no absolute, drive, URI, backslash, traversal, reserved, colliding, or overlong component. |
| MediaType | string | Non-empty declared content type. |
| SizeBytes | integer | Non-negative exact file length. |
| SHA256 | string | Exactly 64 lowercase hexadecimal digits. |
| ByteContract | byte contract | Explicit text/binary characteristics; never inferred during validation. |

## Byte Contract

| Field | Type | Rules |
|---|---|---|
| Encoding | enum | `utf-8`, `ascii`, `binary`, or `invalid_utf8`; expanding the set requires a manifest-contract revision with deterministic verification semantics. |
| BOM | enum | `present`, `absent`, or `not_applicable`. |
| LineEndings | enum | `lf`, `crlf`, `cr`, `mixed`, `none`, or `not_applicable`. |
| FinalNewline | enum | `present`, `absent`, or `not_applicable`; explicit for text and not applicable for binary. |

## Case Expectation

| Field | Type | Rules |
|---|---|---|
| Result | enum | `accepted` or `rejected`. |
| Stage | enum | For rejection: `parse`, `structure`, `semantics`, or `integrity`; omitted for acceptance. |
| DiagnosticContains | string | Non-empty stable Cueson-owned fragment for rejection; omitted for acceptance. |

## Golden Comparison Input

A golden comparison is identified by a fixture ID and logical surface. It receives explicit expected and observed values rather than filesystem paths.

| Surface | Comparison rule |
|---|---|
| Model or generic JSON | Semantic object equality preserving null versus empty values and array order. |
| Diagnostics | Exact ordered equality of stable fields. |
| Rendered or source bytes | Exact byte equality with first differing offset and lengths. |
| Integrity | Exact byte length and lowercase SHA-256. |
| Timestamps | Exact kind and status; restored observations also require declared effective precision. |

## State Transitions

A fixture progresses from `discovered` to `manifest_validated`, `integrity_verified`, `expectation_verified`, and `accepted`. Any provenance, path, inventory, byte-contract, digest, expected-result, or path-leak failure terminates verification before acceptance. Verification is read-only and never repairs or rewrites a fixture.
