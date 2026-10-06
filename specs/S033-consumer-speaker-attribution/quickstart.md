# S033 Quickstart Verification

Use current source or a1.2.0-dev development binary. Released1.1.0 is not expected to accept new fields.

1. Encode a native fixture to Cue JSON, add speaker_attributions to an existing cue and retain its native speakers, then validate and inspect.
2. Save/read the JSON and confirm exact IDs, repeat order and intervals. Restore and compare original source bytes.
3. Add media_timing with duration600000 and timeline_start30000. An assignment ending630000 within its cue is valid;630001 rejects.
4. Remove media_timing: cue-valid assignments remain valid and inspection reports unavailable media evaluation.
5. Retain a cue outside declared media but no invalid consumer interval: runtime warns, source restore remains exact.
6. Render/convert to each native target: omission reports are complete; strict creates/replaces no output.
7. Load historical1.0.0 and1.1.0 fixtures: their identities and semantics are preserved.

Run focused tests for consumer handling in model/schema/convert/codec/CLI/source, then full foreground Go/nested-module/format/docs/security/build verification. Hosted native/race/site/release-proof checks must be green for the final PR head. Record actual evidence in verification.md.
