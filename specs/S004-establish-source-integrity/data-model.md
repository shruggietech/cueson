# Data Model: Source Integrity and Restoration

## Validated Source Asset

A validated source asset is an ordered `model.SourceAsset` whose envelope has passed all source-specific checks before any destination is opened.

| Field | Type | Rules |
|---|---|---|
| AssetID | string | Existing model identifier rules apply; unique within the document. |
| Role | string | Existing schema enumeration applies. |
| StoredBasename | string | Portable safe basename; collision key is NFD, default Unicode case fold, then NFD. |
| DeclaredBytes | uint64 | Must equal the canonical base64 decoded byte count. |
| DeclaredSHA256 | lowercase hex string | Must equal SHA-256 of decoded bytes. |
| EncodedSource | string | Canonical standard padded base64 only. |
| Timestamps | model.Timestamps | Optional observations with explicit creation provenance. |

The validated representation may retain offsets and declarations needed to decode again, but it does not retain a second whole decoded copy of the bundle.

## Portable Identity

Portable identity is the normalized comparison key used for stored basenames and destination collision detection.

1. Normalize the input to Unicode NFD.
2. Apply Unicode default case folding.
3. Normalize the folded value to Unicode NFD again.

The original spelling remains authoritative for output. The identity is comparison-only.

## Destination Plan

A destination plan is the complete ordered asset-to-path mapping built before output begins.

| Field | Type | Rules |
|---|---|---|
| Asset | validated source asset | Exactly one source asset. |
| Destination | absolute runtime path | Derived from the selected destination mode. |
| CollisionIdentity | portable identity | Unique within the plan, plus native path-equivalence checks. |
| ExistingState | absent or regular | Links and all non-regular entries are rejected. |
| ForceAuthorized | bool | True only when `--force` covers an existing regular destination. |

Planning states are `planned`, `staged`, `committed`, `metadata-verified`, `accepted`, and `rolled-back`. A bundle is accepted only after every asset reaches metadata-verified or an allowed metadata result.

## Staged Asset

A staged asset is a same-directory temporary regular file with mode `0600` whose bytes have been flushed, closed, reopened, and verified against the declared length and digest. Its identity is recorded so cleanup and rollback never remove a path that an external actor replaced.

## Timestamp Observation

| Field | Type | Rules |
|---|---|---|
| Kind | creation, modification, or access | Matches the schema field. |
| Instant | Unix nanoseconds | Exact stored value; range-checked before native conversion. |
| Provenance | platform creation, birth time, ctime fallback, or unavailable | Unix change time is never promoted to creation time. |

Capture obtains these observations from the same no-follow regular-file handle before the first content read.

## Metadata Mode

| Mode | Behavior |
|---|---|
| Default | Apply observable timestamps; warn and continue for unsupported values; fail for application or verification failures. |
| Strict | Roll back unless every captured requested timestamp is reproduced exactly at the platform's effective precision. |
| None | Skip timestamp application intentionally and return no timestamp claims. |

`Strict` and `None` are mutually exclusive.

## Timestamp Result

| Field | Type | Rules |
|---|---|---|
| Kind | creation, modification, or access | One result per requested observation. |
| Status | restored, unsupported, unavailable, or failed | Exactly one terminal classification. |
| EffectivePrecision | duration or platform label | Present when relevant to verification. |
| Detail | string | Deterministic explanation suitable for diagnostics. |

Byte restoration success is reported independently from metadata fidelity.

## Asset Restoration Result

| Field | Type | Rules |
|---|---|---|
| AssetID | string | Source asset identifier. |
| Destination | path | Final caller-selected path. |
| Bytes | uint64 | Recomputed from the staged or final file. |
| SHA256 | lowercase hex string | Recomputed from the staged or final file. |
| TimestampResults | ordered list | Empty in no-metadata mode; otherwise truthful per-request results. |

## Capture Result

Capture returns a `model.SourceAsset` containing only the safe basename, canonical encoded bytes, exact length, exact digest, and timestamp observations. It never stores the caller's source path.
