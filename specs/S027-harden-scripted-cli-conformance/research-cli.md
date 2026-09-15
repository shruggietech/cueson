# S027 CLI integration research

**Assessed**: 2026-09-15 against merged S026 `771240800aa9ebec35ba3fab19e054cc10a81e74`.

## Existing authority and implemented behavior

The constitution requires immutable source envelopes, truthful installed capabilities, privacy-safe derived output, bounded hostile input, and unchanged historical CLI compatibility. `docs/compatibility.md` distinguishes the public CLI streams, classes and options from implementation inspection-report internals and diagnostic prose. S027 retains `1.1.0-dev` and experimental ASS/SSA; stable promotion belongs to #64/#65.

`internal/cli/input.go` already provides the shared Cue JSON precedence, scripted content-first detector, explicit selector aliases, local historical schema selection, semantic validation, source integrity and native codec lookup. `normalizeInputFormat` already recognizes `ass` and `ssa`. Native and Cue JSON validation/inspection already succeed through this path, including misleading `.srt`, `.vtt` and `.json` names, as exercised by `scripted_native_test.go`, `scripted_input_test.go` and `versioned_input_test.go`. Reimplementing a separate scripted classification or validation stack would create unnecessary drift.

Inspection already separates a document's declared capability state from installed decoder/renderer availability. Earlier schema-only scripted observations remain schema-only while installed ingest/render report true. Exact restoration remains independent. Scripted reports currently show cue/common counts but contain no native section, record, style, event, override, karaoke or attachment aggregates. Existing `blocks` and document `non_cue_block_count` describe retained WebVTT blocks; they must not be repurposed as a competing scripted record inventory.

## Concrete gaps

- `internal/cli/surface.go` still describes validate/inspect as Cue JSON, SubRip or WebVTT only and limits their advertised `--format` values to `auto`, `cueson`, `srt` and `vtt`. All four static completion definitions derive from these values and therefore omit `ass`/`ssa` for these two commands.
- `finalizeValidateOptions` and `finalizeInspectOptions` advertise the same obsolete list on invalid-format errors, despite accepting scripted selectors.
- `docs/cli.md` repeats that obsolete validate/inspect input vocabulary. Its maintained documentation test will require the approved new values once the catalogue changes.
- Exact help and completion goldens will need regeneration for the deliberately changed vocabulary and descriptions. Other command definitions and goldens should remain unchanged unless their root summary is derived from the changed summaries.
- Scripted inspection currently omits the native structural counts needed by #62. Existing tests check success and privacy but do not check scripted native structural inventory.

## Recommended additive inspection contract

Add a pointer-valued root `scripted` object with `omitempty`, present only for `ass`/`ssa` documents and containing only the following integer aggregates. Keep inspection report version `1`, every existing field's meaning/order, and every existing SubRip/WebVTT human and JSON report byte sequence unchanged. Append the optional object after existing root fields to preserve their relative ordering. The human report adds one deterministic counts-only scripted section only when this object is present.

| Field | Exact definition |
| --- | --- |
| `section_count` | Length of the matching native `sections` collection. |
| `record_count` | Length of the matching native `records` collection, including blank, comment, attachment and preservation-only records. |
| `format_declaration_count` | Number of native records whose kind is `format_declaration`. |
| `style_count` | Length of the native `styles` collection, including invalid preserved styles. |
| `invalid_style_count` | Number of native styles with `valid` false. |
| `event_count` | Length of the native `events` collection, including comment and invalid preserved events. |
| `dialogue_event_count` | Number of native events whose validated `event_type` is `dialogue`. |
| `comment_event_count` | Number of native events whose validated `event_type` is `comment`. |
| `invalid_event_count` | Number of native events with `valid` false. |
| `attachment_count` | Length of the native `attachments` collection; attachment names, payloads, lengths and hashes are excluded. |
| `override_tag_count` | Sum of the lengths of every native event's `tags` collection. |
| `karaoke_span_count` | Sum of the lengths of every native event's `karaoke` collection. |
| `unsupported_karaoke_span_count` | Number of those karaoke entries with `supported` false. |
| `unknown_record_count` | Number of native records whose kind is `unknown`. |
| `malformed_record_count` | Number of native records whose kind is `malformed`. |

These counts deliberately distinguish physical records from logical native owners and common projected dialogue cues. Drawing-only and invalid events are counted as native events regardless of projected cue count. Known collection/aggregate limits are already enforced before report construction; projection does not truncate and never renders native values. A nil matching native branch must not silently produce a truthful-looking empty scripted inventory; the existing validated-input invariant should be explicit in the helper/error path or covered by report-construction rejection.

Do not expose dialect strings, section/style/event/attachment names, actor text, field names or values, tag names/parameters, drawing contents, karaoke text, raw timestamps, captured lines, IDs or other user-controlled native data. The existing report-level `format` identifies the branch; capabilities continue reporting declared and installed facts independently.

## Catalogue and diagnostic maintenance

Keep `orderedCommandSurface` as the advertised vocabulary authority. Extend validate/inspect `--format` option values to `auto`, `cueson`, `srt`, `vtt`, `ass`, `ssa`; add experimental and UTF-8 bounded-profile notes. Derive invalid-selector diagnostics from each command's catalogue option values with a small helper that retrieves the option and formats its values, rather than introducing another hardcoded choice list. Help and completions already consume the catalogue and require no separate scripted generator logic.

Add a focused vocabulary test that checks the canonical catalogue input values against `normalizeInputFormat` and asserts the expected per-command values. Existing `assertCompleteVocabulary` only checks command and option spellings; it does not prove that command-specific value lists are present. Add value checks for each command's generated shell context or rely on exact generated goldens plus an explicit per-command catalogue assertion, rather than accepting a global substring `ass` already present for encode/render/convert.

## Preserve existing encoding error classes explicitly

S026 convert preflights explicit `ass`/`ssa` with a supported non-UTF-8 encoding as invocation failure `2`. Encode, validate and inspect currently defer that known profile rejection to the scripted codec and return runtime failure `1`; `TestScriptedSelectedEncodingRequiresMatchingProfile` explicitly pins encode's existing result. The selected S027 approach retains these owning-command classes instead of silently changing them under a minor compatibility slice. All commands still agree that the bounded scripted profile accepts UTF-8/UTF-8 BOM only. Tests should state this deliberate existing distinction and retain unknown encoding, missing-input, Cue JSON encoding prohibition, automatic classification and WebVTT preflight behavior.

## Focused implementation and verification ownership

- CLI implementation owns `internal/cli/surface.go`, `validate.go`, `inspect.go`, a small catalogue-derived selector diagnostic helper, and a new focused scripted inspection/compatibility test file. No new codec, model branch, source mutation or schema change is necessary for #62.
- Golden ownership includes the changed validate/inspect help files and all four completion goldens under `internal/cli/testdata`; generation must use the executable catalogue and preserve UTF-8 without BOM and LF output.
- Documentation integration updates the validate/inspect sections of `docs/cli.md` and explains the optional counts object and existing encoding-class distinction. Historical released identities and old schema copies remain unchanged.
- Test accepted native and Cue JSON in both dialects; schema-only versus experimental declarations; misleading extensions, aliases, automatic/explicit selection; drawing-only events, valid/invalid styles/events, declarations, preservation-only records, attachments, overrides and supported/unsupported karaoke; report non-mutation; exact repeat determinism and old text report omission/byte equivalence; both human/JSON privacy sentinels; compact one-LF JSON and lowercase snake-case keys; short stdout writes and writer errors; quiet/silent filtering; malformed, integrity and unknown identity failures.
- Run all existing CLI tests, especially the complete frozen historical both-format command matrix, command option/stream/error tests, documentation scenarios, exact help/completion goldens, static completion safety and native shell syntax checks. Root integration owns full corpus, conformance/fuzz/platform/security/release-proof verification and external review.

## Exclusions

This research does not authorize implementation before the blocking analysis gate. It does not promote ASS/SSA to stable, change software/schema identity, add a public Go API, emit native source data in inspection, revise old report semantics, publish releases/schema/site, or perform merge. #63 retains its independent closeable hardening/evidence acceptance after #62 integration within the shared S027 slice.
