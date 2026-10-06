# S033 Consumer Annotations Contract

The full accepted field/timing/Unicode rules are [Issue #82](https://github.com/shruggietech/cueson/issues/82). This contract ratifies its speaker_attributions and media_timing property names and meanings.

Identifiers are consumer-managed strings. Cueson assigns no global/universal identity meaning and no document-local limitation. IDs are data values; schema-valid text does not authorize execution, markup, URL fetching, path resolution or secondary decoding. Native labels stay independent.

## Schema and validation

New fields are optional only under exact1.2.0-dev. Closed attribution objects require speaker_id and paired optional integer endpoints. Closed media_timing requires integer nonnegative duration and optional signed timeline start. Executable semantics check references/timing bounds, Unicode frozen sets and arithmetic. Current examples use exact development identity; released resources remain immutable.

## CLI observations

Existing command vocabulary is retained. validate reports cue conflicts as runtime warnings. inspect JSON version1 optionally adds consumer_annotations with counts and media check state when annotations/declaration exist, never ID values. With no media declaration, timed assignments are checked against cues and media check is unavailable; untimed-only assignments are not evaluated. restore warns for source-cue conflicts while restoring exact bytes; invalid assignments reject.

## Native operation outcomes

No current native format can preserve arbitrary consumer identity or media declarations independently of native labels. Matching render and all twelve conversions report every omission before output. Strict refuses known losses; invalid annotations reject in every mode. Bound warnings/losses at existing ceilings rather than silently truncate. No consumer metadata is promoted to captured source/native provenance.

## Operator boundaries

Explicit push/PR authority covers this slice. Final merge, tag, release publication and production remain human-owned. At most one second Codex review.

