# Inspection Report Contract: Version 1

## JSON root

The root keys appear in this exact order:

```text
inspect_report_version
input
format
schema
capabilities
integrity
assets
document
cues
blocks
diagnostics
loss
```

All Cueson-owned keys use lowercase `snake_case`. Optional scalar facts are represented as JSON `null`; arrays are always present, including when empty. No implementation-specific fields may be added without updating this contract and its golden and shape tests.

## Ordering

- Assets retain source-envelope order and expose only their zero-based report index.
- Cues sort by validated ordinal.
- Native blocks sort by validated `source_order`, then type as a deterministic tie-breaker even though valid WebVTT source order is unique.
- Diagnostics retain validated document order and resolve a cue reference to `cue_ordinal` when possible.
- Object property order follows the fixed report structs and is stable across repeated runs.

## Privacy exclusions

The following values are prohibited anywhere in JSON or human output:

- `data_base64` and decoded source bytes;
- SHA or other content hashes;
- source-envelope asset IDs and cue IDs;
- source identifiers, payload text, token text, speaker names, OCR text, native raw syntax, and setting values;
- timestamp values;
- free-form diagnostic messages;
- caller input paths, working directories, usernames, hostnames, drive or mount details, and local machine identifiers.

Safe stored basenames remain allowed because the source contract already defines them as portable document data.

## Human projection

Human output is rendered exclusively from the completed version-1 report. It groups input and format, schema and capability state, integrity and asset facts, aggregate document facts, cue summaries, block summaries, diagnostic codes and locations, and target-free loss state. It never reads the source document separately, ensuring that JSON and human modes expose identical facts.

## Loss state

Without a conversion target, the report contains exactly:

```json
{
  "status": "not_evaluated",
  "reason": "target_format_required"
}
```

This is not a zero-loss assertion. Target-specific loss evaluation requires a future contract change.
