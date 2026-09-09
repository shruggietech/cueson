# Fixture Manifest Contract

## Location and authority

`testdata/manifest.json` is the only central fixture inventory. `testdata/README.md` explains contributor procedure but does not override the manifest contract. The embedded representative at `internal/schema/testdata/representative.cueson.json` remains a separately governed canonical schema example and is not an orphan in the root corpus.

## Shape

The root object contains exactly `manifest_version` and `fixtures`. Version 1 fixture records contain exactly `id`, `purpose`, `class`, `origin`, `redistribution`, `artifacts`, and `expectation`. Artifact and provenance subobjects reject unknown fields.

The complete field semantics are defined in [data-model.md](../data-model.md). Cueson-owned property names use lowercase `snake_case`.

## Inventory and path rules

Every regular payload beneath `testdata/fixtures` and `testdata/malformed` appears in exactly one artifact record. The verifier rejects missing declared files and unlisted payloads.

Manifest paths use `/`, are relative to `testdata`, and remain beneath that root after native resolution. Empty segments, `.`, `..`, backslashes, absolute or UNC forms, drive or URI prefixes, control characters, Windows-invalid punctuation, trailing dots or spaces, reserved device stems, components longer than 255 UTF-8 bytes, and portable canonical-caseless collisions are invalid. Links and non-regular payloads are invalid.

## Provenance and redistribution

Every fixture has a non-empty purpose and origin. Synthetic and derived origins include a reproducible recipe. Third-party origins include a stable source reference and retrieval date. Every committed fixture has redistribution status `approved`, a license identifier or bundled license reference, attribution, and an explicit NOTICE decision. Unknown or denied redistribution is never accepted.

## Integrity and byte contract

The verifier reads payload bytes without text decoding, checks exact length and SHA-256, and validates declared byte characteristics without rewriting. Authoritative byte areas are excluded from Git text and whitespace normalization. Manifest, README, and normalized expected JSON remain UTF-8 without BOM and LF-normalized repository text.

## Diagnostics

Failures use manifest order and labels such as `fixture <id> artifact <artifact-id>`. They do not include the discovered repository root, temporary directory, username, hostname, or absolute payload path.
