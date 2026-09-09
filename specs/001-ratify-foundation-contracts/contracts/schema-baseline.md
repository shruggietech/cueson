# Schema Baseline Contract

## Identity and compatibility

- Dialect: JSON Schema Draft 2020-12, identified in the schema artifact's `$schema` keyword.
- Development version: `0.0.0`.
- Schema artifact `$id` and Cue JSON instance `$schema`: `https://cueson.io/schema/v0.0.0/cueson.schema.json`.
- The canonical URI is valid as an identifier before public-domain activation.
- Unreleased v0.0.0 may be refined; released schemas are immutable.
- Before v1, breaking changes require a documented minor version and patches remain non-breaking.
- Official software and schema versions remain equal; third-party producer versions are independent.

## Initial formats and capabilities

For the completed v0.0.0 milestone, `subrip` and `webvtt` use:

```json
{
  "status": "envelope_only",
  "ingest_supported": false,
  "render_supported": false,
  "restore_supported": true,
  "ocr_required_for_semantic_output": false
}
```

The declaration describes official Cueson capability for the document's format. Schema recognition does not imply an available codec. Exact restoration is a public generic source-envelope capability owned by issue #6.

## Common and source model rules

- Cues expose normalized timing, `raw_text`, `plain_text`, logical lines, speakers, tokens, OCR observations, placement, and format-native data.
- Every cue contains an `ocr_observations` array, which may be empty.
- Each non-empty OCR observation is independently identified and provenanced, including its engine identity and a resolvable source reference; if `derived` is retained, it is always true.
- The multi-asset source envelope preserves ordered assets, portable safe basenames, exact byte lengths, SHA-256 identities, original bytes, and truthful timestamp metadata.
- Portable safe basenames reject path syntax, control characters, Windows-invalid punctuation and trailing characters, NTFS alternate-data-stream syntax, and case-insensitive Windows reserved device names even when followed by extensions.
- Basenames within one source bundle are unique under Unicode canonical caseless matching: NFD normalization, default Unicode case folding, then NFD normalization again.
- Original filesystem paths and machine identifiers are prohibited.
- Cueson-owned keys use lowercase `snake_case`; JSON Schema keywords keep standards-defined spelling.
