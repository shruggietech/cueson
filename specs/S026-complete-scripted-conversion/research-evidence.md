# S026 conversion evidence research

**Date**: 2026-09-15

**Authority reviewed**: Root AGENTS.md, constitution 0.1.0, S026 feature specification, architecture of record, and the current bodies of GitHub #60 and #61. This artifact is research only; implementation remains behind the coordinating agent's blocking analysis gate.

## Existing authority and integration boundaries

The root manifest is the sole inventory of all payloads below root testdata/fixtures, testdata/malformed, and testdata/fuzz. It permits exactly five artifact roles: source, expected_model, expected_diagnostics, expected_bytes, and malformed_input. Artifact identifiers are independent from roles. A conversion expectation containing reparsed timing/text plus its complete ordered loss report can therefore use identifier expected_conversion and role expected_model, without changing generic fixture infrastructure. Native byte goldens belong under expected/bytes and use role expected_bytes. Existing conversion report and warning digests remain unchanged regression oracles.

Existing S025 scripts and projection/diagnostic expectations remain native ingest/render/restoration evidence. S026 should reuse those sources by fixture identity in targeted cross-package tests, and append conversion-specific expectations beneath new scripted-conversion case directories. No S025 source, native projection, native diagnostic, or canonical-render expectation should be rewritten to bless conversion output.

The existing conformance matrix has three format groups (srt, vtt, scripted) and exactly seventeen scripted rows. Its verifier resolves named Test/Fuzz functions and fixture IDs, not individual subtest coverage. A new conversion evidence artifact must not add a fourth format group or change the row count. Relevant conversion_test deferrals currently all cite #60/#61; replace them with real focused tests only after those tests enforce the row. Retain platform/hardening #63 and stable/release #64/#65/#66 deferrals, and retain genuine inapplicability for acquisition/dialect-decoding rows.

The existing Report.Validate association set contains common cue positions, WebVTT block positions, and diagnostic positions. Atomic losses for scripted sections, declarations, styles, comments, attachments, and uninterpreted physical records need validated scripted source-order association as well. New tests must distinguish a valid document record reference from a fabricated position, and require JSON pointers to resolve against the immutable source document.

## Compact governed fixture plan

Use the four loss-free text/dialogue baselines as shared semantic seeds, retaining overlap, multiline order, literal ampersand, bold, italic, and underline. Keep the two existing SubRip/WebVTT conversion fixture identities and their old exact target bytes as authoritative regressions. Add one ASS baseline and one SSA baseline for the six scripted-source directions, and one SubRip/WebVTT target-default baseline for each text source. Each source fixture can own three expected target-byte artifacts and one expected_conversion JSON artifact indexed by canonical target format. This yields twelve explicit distinct-pair rows without duplicating source records for every target.

Suggested new fixture IDs are scripted-conversion/ass-baseline, scripted-conversion/ssa-baseline, scripted-conversion/subrip-baseline, and scripted-conversion/webvtt-baseline. Every artifact receives explicit synthetic/project-authored origin, reproducible recipe, MIT redistribution approval, Copyright ShruggieTech attribution, NOTICE false, byte characteristics, exact byte count, and SHA-256. Expected conversion JSON is independently authored from the agreed contract; it must not be accepted through an unexplained update-goldens mode.

Suggested loss-family sources are scripted-conversion/ass-native-losses and scripted-conversion/ssa-native-losses, with compact script/style metadata, harmless unknown fields and records, non-dialogue events, actor/effect, event margins, drawing spans, unsupported overrides, karaoke, and an inert font/graphic attachment. The source must remain accepted by the native profile; malformed/unsafe source examples already exist and should be reused for explicit fatal conversion refusals rather than being sanitized.

Text-to-scripted lossy evidence should reuse the existing conversion/srt-lossy and conversion/webvtt-lossy source fixtures to protect old loss meanings, while new expectations distinguish target-specific representable emphasis from omitted font, settings, voice, inline timing, and metadata. Add compact precision cases for exact multiples of ten milliseconds, below-half, tie, above-half, equal quantized endpoints, and checked upper-bound arithmetic, once the plan ratifies its quantization rule.

Drawing-only, empty dialogue, zero-dialogue scripts, diagnosed malformed braces, unused invalid styles, malformed comments, and malformed attachments already have paired S025 native fixtures. Conversion-specific tests should assert their documented edge outcome in each affected direction. Empty readable output must never become invented placeholder text or duration. A safely accepted native source can still have a fatal target representation edge; the root fixture class describes source acceptance, while expected_conversion records conversion rejection.

## Proposed executable evidence owners

- TestScriptedTwelveDirectionConversionConformance: table of all twelve distinct canonical source/target pairs, exact authored target-byte comparisons, independent target codec reparse, expected timing/plain-text/lines comparison, complete ordered loss comparison, repeated byte/report determinism, and immutable before/after source-model comparison.
- TestScriptedConversionCompleteLossConformance: paired ASS/SSA native-loss sources and reused text-loss sources, atomic expected entries for every field/tag/record omission or degradation, stable code membership, valid source references, deterministic source ordering, and no source excerpt or caller identity in loss context.
- TestScriptedConversionStrictFatalPublicationConformance: every governed lossy/fatal edge tested for default stdout, explicit stdout, new destination, and forced existing destination; failure returns no payload, creates no new destination, preserves existing bytes, and leaves input bytes and directory entries intact except intentionally prepared files.
- TestScriptedConversionPrecisionAndDefaultsConformance: exact/tie/neighbor quantization cases, endpoint loss cardinality, interval-collapse rejection, maximum time/overflow rejection, deterministic default style/script output, and reparsed semantic agreement. Tests must assert expected numbers and fields, not only self-consistency of implementation output.
- TestScriptedConversionSourceIntegrityPrivacyAndHistoricalConformance: native versus encoded Cue JSON parity, supported historical text Cue JSON to scripted targets, bad source hashes/lengths rejected before output, complete model/source before/after identity, and explicit native/slash/backslash/JSON-escaped caller-path privacy checks.

The new conformance file may reuse existing fixtureRoot, verifiedManifest, mustFixture, mustArtifact, readArtifact, encodeScriptedFixture, and decodeScriptedExpected helpers from the external conformance_test package. Use cli.Run for publication behavior and internal/convert.Convert on schema-decoded encode output for report/model immutability. Build a valid source envelope through the actual encode workflow; the existing internal convert test helper reuses representative original source bytes and is unsuitable as the source-truth oracle for independently authored conversion fixtures.

## Twelve-direction acceptance obligations

| Source | Targets | Baseline | Known-loss evidence | Fatal evidence |
|---|---|---|---|---|
| SubRip | WebVTT, ASS, SSA | Existing WebVTT golden plus two deterministic scripted target goldens | Existing coordinates/font/common observations; millisecond endpoints for scripted targets | Positive target duration and text/control representation rejection |
| WebVTT | SubRip, ASS, SSA | Existing SubRip golden plus two deterministic scripted target goldens | Document blocks/settings/voice/tokens, target-specific markup and precision | Target-incompatible literal/control syntax, quantized interval collapse |
| ASS | SubRip, WebVTT, SSA | Readable dialogue/timing/shared emphasis and equivalent representable dialect semantics | Every native style/record/attachment/tag/event/variant loss | Preservation-only malformed owners, unsafe edited/native content, unsupported text or target shape |
| SSA | SubRip, WebVTT, ASS | Readable dialogue/timing/shared emphasis and equivalent representable dialect semantics | Every native style/record/attachment/tag/event/variant loss | Preservation-only malformed owners, unsafe edited/native content, unsupported text or target shape |

A baseline should have the contract's expected loss list; do not assert that arbitrary script/style defaults are loss-free when observable presentation is omitted. Shared emphasis must be compared semantically after target parse, and omitted native typography/layout must remain accounted for. Every direction has at least one successful baseline and applicable lossy/fatal assertions; same-format rendering remains outside the twelve-direction graph.

## Ownership and verification boundary

The independent evidence implementer can own new internal/conformance/scripted_conversion_test.go, new testdata/fixtures/scripted-conversion payloads, testdata/manifest.json append-only additions, testdata/conformance-matrix.json conversion evidence updates, and the corresponding README explanation. The coordinating agent owns production conversion integration, CLI surface/goldens, docs/formats and full foreground verification. Coordinate fixture-contract decisions and named production tests before replacing matrix deferrals, and do not concurrently mutate the manifest or matrix from another task.

Focused verification is go test -count=1 ./internal/testutil ./internal/conformance ./internal/convert ./internal/cli, followed by root verification and standalone docs-verify checks coordinated by root. Check git attributes for new source and expected/bytes artifacts, UTF-8 without BOM for metadata, exact integrity from working-tree bytes, and no mojibake in authored delivery artifacts. Windows non-Git tooling must run through a verified CREATE_NO_WINDOW launcher with redirected noninteractive input/output.
