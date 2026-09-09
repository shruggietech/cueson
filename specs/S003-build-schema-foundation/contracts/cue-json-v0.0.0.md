# Cue JSON v0.0.0 Contract

## Artifact identity

- Dialect: `https://json-schema.org/draft/2020-12/schema`
- Canonical schema identifier: `https://cueson.io/schema/v0.0.0/cueson.schema.json`
- Cue JSON instance `$schema`: `https://cueson.io/schema/v0.0.0/cueson.schema.json`
- Cue JSON instance `schema_version`: `0.0.0`

## Validation contract

Validation parses one JSON value, applies the embedded Draft 2020-12 schema without network resolution, then applies format-neutral semantic invariants. A failure is returned as data to the caller and is never printed by the model or schema package.

Structural validation owns required properties, fixed object boundaries, primitive shapes, enum membership, schema/version constants, format-conditioned extension objects, safe single-basename syntax, lowercase SHA-256 syntax, and standard-base64 syntax.

Semantic validation owns document-local references and uniqueness, exactly one primary role, cue timing arithmetic, ordinal and source ordering, summary counts, timestamp equivalence and provenance, OCR source references, matching format data, and S003 capability truth.

Issue #6 owns decoded source length and SHA-256 comparison, Unicode canonical caseless basename collision enforcement, source timestamp capture ordering, destination planning, exact restoration, and restored metadata verification.

## Availability contract

The schema recognizes `subrip` and `webvtt` as envelope-only formats. During S003, native ingest, model-driven render, exact restore, and OCR-required flags are all false. Recognition is not a codec claim.

## Extension contract

Stable fixed objects reject additional properties. Intentionally extensible maps accept only namespaced lowercase keys in the form `vendor:name` where each segment uses lowercase letters, digits, and underscores and begins with a lowercase letter.

## Compatibility contract

The executable version, embedded schema version, schema identifier version, and instance version remain exactly equal for official artifacts. Third-party producer version values are independent.
