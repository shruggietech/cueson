# Consumer speaker attribution

**Status:** Implemented in exact `1.2.0` stable release candidate source. Published v1.1.0 packages and immutable public schemas do not include these fields.

Cueson stores speaker assignments made by downstream applications. A `speaker_id` is an opaque string whose meaning, lifetime and scope are assigned by the consumer. Cueson never treats it as a global or universal identity. An application may share IDs across its own collection or use temporary IDs for one file. Names and UUIDs are both ordinary values. Cueson requires no registry, UUID format, automatic ID generation or diarization engine.

Each cue optionally contains an ordered `speaker_attributions` array. This is independent of source-derived `speakers`, WebVTT voice labels, ASS/SSA names and native source text. Adding consumer data never promotes native labels into identities or changes the preserved source envelope. Missing and empty arrays mean no assignments; `null` is invalid. Duplicate IDs, repeated entries, caller order and overlapping assignments remain intact. The array permits at most 1,024 entries per cue.

```json
{
  "media_timing": {
    "duration_milliseconds": 120000
  },
  "speaker_attributions": [
    {"speaker_id": "application-speaker-7"},
    {
      "speaker_id": "temporary-speaker-B",
      "start_milliseconds": 1200,
      "end_milliseconds": 1800
    }
  ]
}
```

This example shows the field shapes only. `media_timing` belongs at the document root; `speaker_attributions` belongs inside a cue. It is not a complete Cue JSON document. An untimed assignment states participation somewhere within that cue and does not assert continuous speech throughout the cue.

## Identifier limits

IDs contain 1 to 256 decoded Unicode scalar values. Valid supplementary characters, combining characters, internal spaces and joiners are allowed. Case and Unicode normalization are preserved exactly. Comparisons use the exact decoded string; Cueson never trims, normalizes, coerces or truncates it.

The following are rejected anywhere: C0 controls U+0000–U+001F, C1 controls U+007F–U+009F, U+2028/U+2029 line separators, bidi controls U+061C/U+200E/U+200F/U+202A–U+202E/U+2066–U+2069, and U+FEFF. Unpaired surrogate escapes and invalid UTF-8 reject before ordinary JSON decoding can replace them. Leading or trailing Unicode White_Space is rejected using the frozen set U+0009–U+000D, U+0020, U+0085, U+00A0, U+1680, U+2000–U+200A, U+2028, U+2029, U+202F, U+205F and U+3000.

IDs are inert data. Markup-looking, URI-looking, path-looking, Base64-looking, `__proto__` and `constructor` values have no executable or special identity meaning. Cueson never executes, resolves, fetches or decodes them as another payload. Consumers must escape values for their own display or query context and use them as values rather than code, filesystem operations or object property authority. A string schema cannot replace correct handling by a downstream consumer. No original filesystem path or local machine identifier may be sourced into Cue JSON.

## Timings and declared media boundaries

An attribution either supplies both `start_milliseconds` and `end_milliseconds` or neither. Present values are exact nonnegative signed-64-bit integers, not fractional numbers or numeric strings. A timed attribution is a positive half-open interval `[start, end)` within its cue. Values use the document's absolute millisecond timeline, so cue start/end boundaries are inclusive limits on the supplied interval endpoints. No implicit rounding, clamping, splitting or rebasing occurs. Editing a cue can therefore invalidate a previously valid assignment.

Optional root `media_timing` declares a known media duration using required nonnegative `duration_milliseconds` and optional signed `timeline_start_milliseconds`. Omitted timeline start evaluates as zero and remains omitted when serialized. Zero duration means known zero duration. The media interval is `[timeline_start, timeline_start + duration)`; overflow rejects. When this declaration exists, every timed attribution must also fit inside it. A signed offset can describe a source placed on a wider timeline; attribution endpoints remain nonnegative. Untimed assignments have no interval to check.

Cueson validates consistency with the consumer's declaration. It does not independently measure or authenticate the audio. Subtitle coverage (`document.media_end_milliseconds` or `media_span_milliseconds`) is not source duration: subtitles can omit opening silence, trailing content or entire portions of an audio track. Ingest does not infer a duration from cues, encoded byte count, file timestamps or unknown subtitle metadata. No media probe, decoder or network lookup is added.

Without declared duration, cue checks still run and the media boundary check is `unavailable`. With duration and no timed assignments it is `not_evaluated`; with timed assignments it is `checked`. A checked state says only that the supplied interval constraints passed. Original cue intervals outside the declaration produce one runtime `consumer_cue_media_conflict` warning with a conflict count. They remain unchanged and exact-restorable. Runtime warnings never rewrite archived `diagnostics` or `stats`.

## Inspection and native exports

Inspection report version `1` gains optional `consumer_annotations` containing attribution/timed/untimed counts, media declaration presence, media-check state and conflicting-cue count. Human inspection reports the same summary. IDs are never printed by either inspection form or consumer diagnostics. Documents with neither assignments nor a media declaration retain the prior inspection layout. `validate` explicitly reports media-check availability for annotated documents.

All four matching native renderers and all twelve cross-format conversion directions omit these consumer fields. Native rendering reports `consumer_speaker_attribution_omitted` once per attribution occurrence and `consumer_media_timing_omitted` once per declaration. Conversion uses the corresponding stable loss codes in the [conversion vocabulary](conversion.md#stable-loss-vocabulary), with deterministic occurrence pointers. Repeated entries each count. A bounded diagnostic or loss limit causes refusal rather than incomplete accounting. Strict mode refuses before publishing any output. IDs do not appear in generated native bytes or loss messages, and native labels remain unchanged.

## Version boundary

Current candidate source emits exact schema and executable identity `1.2.0`, with a byte-identical [immutable candidate schema](../schema/releases/v1.2.0/cueson.schema.json). It validates supported historical `1.1.0` and `1.0.0` documents locally without network retrieval and preserves their identity and source observations. Historical schemas reject the new fields; adding fields requires selecting the current contract explicitly. Older executables have no promised forward compatibility. The former `1.2.0-dev` identity is unsupported. Released schema bytes and production routes/downloads remain their independently published versions. Candidate preparation does not publish the 1.2.0 schema URI or release packages; publication and production deployment require separate authorization.

The complete implementation and acceptance mapping is in [S033](../specs/S033-consumer-speaker-attribution/spec.md), tracked by [issue #82](https://github.com/shruggietech/cueson/issues/82).
