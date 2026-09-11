# Data Model: Complete CLI Workflows

## Validated Input

An internal read-only result shared by validation, inspection, and conversion.

| Field | Meaning | Rules |
|---|---|---|
| `document` | Validated common and native Cue model | Passed schema and integrity validation for Cue JSON, or native decoding plus model and generated-envelope integrity validation |
| `input_kind` | `cue_json` or `native_subtitle` | Derived once; never guessed after a Cue JSON candidate fails |
| `selection_basis` | `cue_json`, `explicit`, `content`, or `extension` | Describes why the class or native format won |
| `content_format` | Optional canonical native format | Present only when content detection produced evidence |
| `extension_format` | Optional canonical native format | Present only when the safe input basename extension is registered |
| `diagnostics` | Ordered safe model diagnostics | Source order is retained; presentation applies command filters |

## Inspection Report Version 1

The report is a fixed CLI-owned JSON object. It is not Cue JSON and does not enter the public Go package surface.

| Field | Type | Meaning |
|---|---|---|
| `inspect_report_version` | string | Exactly `1` |
| `input` | Input Summary | Classification and non-path evidence |
| `format` | string | Canonical `subrip` or `webvtt` |
| `schema` | Schema Summary | Embedded identifier, version, and compatibility result |
| `capabilities` | Capability Summary | Declared document support and installed runtime operations |
| `integrity` | Integrity Summary | Verified status plus checked asset and byte totals |
| `assets` | ordered Asset Summary array | Safe source-envelope facts in stored order |
| `document` | Document Summary | Aggregate cue, body, timing, word-timing, and diagnostic facts |
| `cues` | ordered Cue Summary array | Structural cue facts ordered by ordinal |
| `blocks` | Block Summary | Aggregate native block counts plus ordered safe block entries |
| `diagnostics` | ordered Diagnostic Summary array | Stable codes and structural locations without messages |
| `loss` | Loss State | Target-free conversion-loss assessment state |

## Input Summary

| Field | Type | Rules |
|---|---|---|
| `kind` | string | `cue_json` or `native_subtitle` |
| `selection_basis` | string | `cue_json`, `explicit`, `content`, or `extension` |
| `content_format` | nullable string | Canonical native format or null |
| `extension_format` | nullable string | Canonical native format or null |

The report never includes the caller path or working directory.

## Schema Summary

| Field | Type | Rules |
|---|---|---|
| `id` | string | Embedded canonical schema identifier |
| `version` | string | Embedded schema version |
| `compatible` | boolean | True only after executable/schema lockstep passes |

## Capability Summary

`declared` reproduces the validated format-support booleans and status from the document. `installed` contains fixed `ingest`, `render`, `restore`, `validate`, and `inspect` booleans derived from runtime authorities. Schema recognition alone never sets native ingest or render true.

## Integrity Summary

| Field | Type | Rules |
|---|---|---|
| `status` | string | Exactly `verified` after complete source-envelope validation |
| `asset_count` | integer | Number of checked assets |
| `total_bytes` | integer | Sum of validated decoded lengths with overflow rejection |

No content hash is emitted.

## Asset Summary

| Field | Type | Rules |
|---|---|---|
| `index` | integer | Zero-based stored order |
| `primary` | boolean | True for the primary asset |
| `role` | string | Validated source role |
| `file_name` | string | Validated portable safe basename |
| `media_type` | nullable string | Validated declared media type |
| `bytes` | integer | Validated decoded byte count |
| `encoding` | nullable Encoding Summary | Canonical encoding, BOM, line-ending, and confidence observation without text |

Asset ID, timestamps, hash, and base64 data are prohibited.

## Document and Cue Summaries

Document Summary contains `cue_count`, `non_cue_block_count`, `body_item_count`, nullable media bounds and span, `has_word_level_timing`, and diagnostic totals by severity.

Each Cue Summary contains `ordinal`, `source_order`, `start_milliseconds`, `end_milliseconds`, `duration_milliseconds`, `payload_line_count`, `speaker_count`, `token_count`, `ocr_observation_count`, `has_source_identifier`, `has_placement`, `native_setting_count`, `native_setting_occurrence_count`, and `invalid_native_setting_count`. Cue IDs and every content-bearing value are prohibited.

## Block Summary

Aggregate fields are `total`, `note_count`, `style_count`, `region_count`, and `unknown_count`. Ordered entries contain `source_order`, `type`, `line_count`, `has_parsed_region`, `native_setting_count`, `native_setting_occurrence_count`, and `invalid_native_setting_count`. Raw block content and settings values are prohibited.

## Diagnostic Summary

| Field | Type | Rules |
|---|---|---|
| `severity` | string | Validated severity |
| `code` | string | Stable diagnostic code |
| `source_order` | nullable integer | Original structural location |
| `cue_ordinal` | nullable integer | Resolved from cue ID internally, never the identifier itself |

Entries retain validated source order. Free-form message text is prohibited.

## Loss State

| Field | Value |
|---|---|
| `status` | `not_evaluated` |
| `reason` | `target_format_required` |

No loss count is present because zero would imply an evaluation occurred.

## Completion Surface

The ordered surface catalogue contains commands and their option definitions. Each option records canonical spellings, value kind, canonical enum candidates, repeatability, description, and applicable command. Completion scripts are deterministic projections of this catalogue; semantic conflict validation remains in the parser.
