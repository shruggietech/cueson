# S024 input and capability contract

## Supported exact identities

- Historical:https://cueson.io/schema/v1.0.0/cueson.schema.json with schema_version1.0.0.
- Current development:https://cueson.io/schema/v1.1.0-dev/cueson.schema.json with schema_version1.1.0-dev.

Only those complete pairs are accepted. Missing/wrong types, URI/version mismatch, approximate URI, v0.0.0 and unknown versions reject before publication. Original parsed input is validated without modification. All schemas and references are local; no input-controlled retrieval occurs.

## Commands

Historical SubRip/WebVTT validate, inspect, restore, matching render and existing target conversion retain the ratified1.0.0 contract and all input identity/provenance/native/source facts. Current software lockstep checks its embedded current schema, never a loaded historical version. Inspect reports the loaded input identity.

New native SubRip/WebVTT encoding emits exact1.1.0-dev and official1.1.0-dev producer identity. Schema/version commands describe current output only. Unsupported identities exit1 with stderr and no payload; invocation errors retain exit2. Source integrity, cancellation, strict losses and destination/force safety remain required.

## Scripted development behavior

Current ass/ssa branches declare schema_only, ingest_supported:false, render_supported:false, restore_supported:true and ocr_required_for_semantic_output:false. Generic structural/semantic validation and privacy-bounded inspection are available. Exact restoration uses verified original assets. Recognized nil codecs produce typed unavailable-capability errors for native ingest/render; conversion to/from scripted native formats remains unavailable until downstream implementation.

Schema shapes/semantics realize the ratified [native contract](../../../docs/formats/ass-ssa.md), including optional construction capture fields, logical lines and unknown-metadata whole-operation rejection. Current model recognition is not proof of stable codec support.

## Deferred authority

No native codec, new immutable release schema, tag/release, production deployment or merge is authorized by completion of these two issues. Specific kickoff grants push/official PR and review remediation; human final merge remains pending.
