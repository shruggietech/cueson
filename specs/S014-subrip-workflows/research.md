# S014 Research

## Schema and executable identity

**Decision**: Advance the evolving canonical schema and development executable to `0.1.0`; preserve `schema/releases/v0.0.0/cueson.schema.json` and published assets unchanged.

**Rationale**: Native capability changes cannot truthfully reuse a released immutable schema identity.

## Codec composition

**Decision**: Register canonical names, aliases, and independently optional detector, decoder, and renderer functions. Exact restoration remains outside codecs.

**Rationale**: WebVTT remains schema-recognized and restorable without falsely claiming native S014 methods.

## Detection and source acquisition

**Decision**: Explicit selection wins, content grammar outranks extension, disagreement warns, and ambiguity fails. Open a regular no-follow file once, capture timestamps first, enforce 64 MiB while streaming, and return the exact parsed bytes with metadata.

**Rationale**: Extensions are hints, and reopening creates a preservation/interpretation race.

## Text decoding

**Decision**: Honor UTF BOMs, strict UTF-8, and only strong tested BOM-less UTF-16 evidence. Legacy single-byte decoding requires an explicit encoding; conflicts and malformed Unicode fail.

## SubRip parsing

**Decision**: Scan physical lines without destructive normalization and use an explicit block state machine. Preserve native sequence/timing lines, normalized timing, coordinates, payload lines, and ordered diagnostics.

**Recorded departure**: Reject reversed timing rather than swapping it, and never remove a speaker prefix from raw/plain content. The legacy converter behavior conflicts with source fidelity.

## Rendering and publication

**Decision**: Render canonical ordered SubRip with LF and final newline. Default render output is stdout. File output stages beside the destination and does not overwrite without `--force`.

## Release proof

**Decision**: Keep historical v0.0.0 release assertions while validating current non-publishing development snapshots against their development executable/schema identity without treating them as a release.
