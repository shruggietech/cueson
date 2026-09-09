# Data Model: Build Schema Foundation

## Cue JSON document

The root entity is fixed and contains the canonical schema URI, schema version, canonical format, official format support, producer, source envelope, normalized metadata, document summary, ordered cues, matching root format data, diagnostics, and statistics.

Validation requires every fixed root area and rejects undeclared root properties. The root format selects exactly one matching format-specific object at both root and cue level.

## Format support

Fields are `status`, `ingest_supported`, `render_supported`, `restore_supported`, and `ocr_required_for_semantic_output`.

For S003, both `subrip` and `webvtt` require `status: envelope_only`; every capability boolean is false. Issue #6 is the only slice authorized to change restoration support to true.

## Producer

Fields are `name` and `version`. Producer software identity is independent from Cue JSON `schema_version` and may differ from `0.0.0` for third-party producers.

## Source envelope

Fields are `primary_asset_id` and ordered `assets`. At least one asset exists. Asset identifiers are unique, exactly one asset has role `primary`, and `primary_asset_id` resolves to that asset.

## Source asset

Fields are `id`, `role`, `file_name`, nullable `media_type`, `size`, `hashes`, `timestamps`, nullable `encoding`, and `data_base64`.

`role` is `primary` or `companion`. `file_name` is a portable safe basename with no path or URI surface. `size.bytes` is a non-negative integer and `size.text` is nullable display data. `hashes.sha256` is exactly 64 lowercase hexadecimal characters. `data_base64` uses standard padded or unpadded-empty base64 syntax. S003 establishes these representations; issue #6 owns decoded-byte comparison and restoration.

## Timestamp observations

The timestamp set contains nullable `created`, `modified`, and `accessed` timestamp values plus `created_source`.

A non-null timestamp value contains RFC 3339 `iso` and signed `unix_ns` values for the same instant. Creation provenance is `birthtime`, `windows_creation_time`, `ctime_fallback`, or `unavailable`. `unavailable` requires `created` to be null; every other provenance requires a non-null creation value.

## Encoding observation

Fields are nullable `bom`, `line_endings`, nullable `detected_encoding`, and nullable `confidence`. Line-ending values are `none`, `lf`, `crlf`, `cr`, `mixed`, or `unknown`; confidence ranges from zero through one.

## Metadata and document summary

Metadata contains nullable `title`, `language`, `kind`, and `description` values. The document summary contains `cue_count`, nullable media start/end/span values, and `has_word_level_timing`.

Counts equal the cue collection. When cues exist, media bounds and span agree with the minimum start and maximum end; an empty collection uses null bounds and zero count.

## Common cue

Fields are `id`, `ordinal`, `source_order`, nullable `source_identifier`, `timing`, `payload`, `speakers`, `tokens`, `ocr_observations`, nullable `placement`, and format-specific data.

Cue identifiers, ordinals, and source-order values are unique. Ordinals are contiguous from zero in array order. Source order is strictly increasing. Timing contains non-negative `start_milliseconds`, `end_milliseconds`, and `duration_milliseconds`; end is not before start and duration equals end minus start.

## Payload

Fields are `raw_text`, `plain_text`, and ordered `lines`. All three are required; they may be empty when an empty native payload is truthful.

## Speaker and token observations

A speaker contains `name` and `origin`, where origin is `native` or `heuristic`. A token contains `text`, normalized start/end milliseconds, and nullable native timing text. Token timing remains inside its containing cue.

## OCR observation

Required fields are `id`, `derived`, `engine`, `source_asset_id`, `text`, and `lines`. `derived` is always true. Optional or nullable fields cover engine version, model, language, confidence, alternatives, source event/image identity, source regions, and namespaced processing options. Observation identifiers are unique per cue and source assets resolve against the source envelope.

## Placement

Placement is either null or a conservative object with nullable normalized `line`, `position`, `size`, and `align` observations. Native placement remains in format-specific data.

## Format-specific data

Exactly one `subrip` or `webvtt` object exists and matches root `format`. SubRip cue data preserves raw sequence and timing lines plus nullable coordinates. WebVTT cue data preserves nullable raw identifier, raw timing/settings/payload, and parsed string settings. Root SubRip data contains dialect and nullable decoding notes. Root WebVTT data preserves signature, nullable description, ordered metadata lines, and ordered non-cue blocks with unique source order.

## Diagnostic and statistics

A diagnostic contains severity (`info`, `warning`, or `error`), code, message, nullable source order, and nullable cue identity.

Statistics contain cue, diagnostic, warning, and error counts, word-timing presence, and nullable media span. All counts and summary values agree with the document collections they summarize.
