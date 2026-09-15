# Scripted-format milestone version compatibility contract

**Status:** S023 future implementation contract. S023 changes planning and documentation only; it does not change runtime validation, executable versions, schemas, codecs, or released artifacts.

**Target:** A prospective v1.1.0 executable and current output schema, conditional on the minor-version gates below. Until a later implementation and release gate passes, the shipped and current-source executable/schema remain v1.0.0.

## Authority and observed baseline

The constitution, maintained [compatibility policy](../../../docs/compatibility.md), [schema contract](../../../docs/schema.md), [CLI contract](../../../docs/cli.md), and [architecture](../../../docs/architecture.md) control the public interface. The immutable released schemas and tagged model implementations define historical validation; the older working project specification does not override them.

At the S022 merged baseline, [internal/schema/schema.go](../../../internal/schema/schema.go) registers and compiles exactly one embedded schema. `schema.Decode` validates UTF-8 and Unicode escapes, parses exactly one JSON value, validates the current structural contract, decodes the typed model, and invokes current `model.Document.Validate`. The current canonical schema fixes the instance `$schema` and `schema_version` to the v1.0.0 pair, and the current model requires stable SubRip/WebVTT capabilities. It rejects a v0.0.0 instance before it can be used as a current document. Merely retaining historical schema files in the repository does not provide historical executable input support.

The released v0.0.0 schema fixes the capability state to `envelope_only`, with `ingest_supported: false`, `render_supported: false`, `restore_supported: true`, and `ocr_required_for_semantic_output: false`. Its model includes representative common/native structures but its executable does not offer native semantic ingest, model-driven rendering, conversion, validation, inspection, or completion. Historical reference implementations are [`v0.0.0:internal/schema/schema.go`](https://github.com/shruggietech/cueson/blob/v0.0.0/internal/schema/schema.go) and [`v0.0.0:internal/model/model.go`](https://github.com/shruggietech/cueson/blob/v0.0.0/internal/model/model.go). These are a distinct immutable baseline, not the stable v1 native contract with an earlier number.

The proposed milestone supports historical v1.0.0 input only. Both the current v1.0.0 executable and proposed v1.1.0 executable reject v0.0.0 input identity. The immutable v0.0.0 executable remains available for its own valid envelopes and exact restoration. Adding v0.0.0 adapters or new envelope inspection/validation support is a separate optional future outcome, outside the v1 compatibility promise and this backlog. A representative common/native structure does not grant semantic capabilities to an envelope-only document.

## Exact contract registry

The future executable recognizes only the following complete identity pairs. A URI is compared as an exact string; redirect targets, case folding, trailing slashes, unversioned aliases, semver ranges, producer versions, and approximate URI matching are not substitutes.

| Input contract | Instance `$schema` | Instance `schema_version` | Local structural authority | Semantic authority |
|---|---|---|---|---|
| Historical stable native | `https://cueson.io/schema/v1.0.0/cueson.schema.json` | `1.0.0` | Byte-identical released v1.0.0 schema | Version-specific stable SubRip/WebVTT semantics derived from tagged v1.0.0 |
| Prospective current native | `https://cueson.io/schema/v1.1.0/cueson.schema.json` | `1.1.0` | Future reviewed canonical schema | Future current model semantics for existing and scripted formats |

The current row is a target identity, not an artifact introduced by S023. P14 owns the final promotion and immutable-copy proof; development fixtures and candidate identity staging before promotion must be documented by P05 and may not overwrite an existing release identity with a new shape.

All structural authorities are bundled locally with the future executable. Historical loading must work offline and may not retrieve a schema, import, reference, attachment, or source asset from an input-controlled URI. The selected registry entry supplies the schema resource and all dependencies. Artifact identity validation checks the registered `$id`, instance `$schema` constant, `schema_version` constant, and expected Draft 2020-12 dialect for each entry independently.

Released schema bytes are immutable, including whitespace and annotations. Their expected SHA-256 digests are:

| Released artifact | SHA-256 |
|---|---|
| v0.0.0 | `d15c7fa5227156109dd6be3d39b711aca3503794bb862169dfca96ee80adb975` |
| v1.0.0 | `1aad14567033d7e14d9beb78985e18007aefb5345095370b11b6b887df7ec541` |

## Validation and command dispatch

The implementation must retain the existing shared validated-input boundary instead of duplicating version selection separately in each command. Restore's direct Cue JSON path must use the same exact contract selection and corresponding integrity requirements, even where its bounded file acquisition differs from other commands.

1. Apply the existing input size, UTF-8, Unicode escape, single-value JSON, and cancellation checks before accepting a document. Extract both identity fields from the parsed object without normalizing them or changing the input.
2. Select exactly one registry entry using the complete identity pair. Missing, wrong-typed, mismatched, or unsupported identities fail explicitly. JSON-looking or `.json`-named invalid inputs retain Cue JSON recognition precedence and must not fall through to a native codec.
3. Validate the original parsed instance against the selected immutable structural schema. Do not rewrite its identity to make it pass the current schema, remove unknown properties, or replace the schema with an identity-relaxed approximation.
4. Decode using a representation that retains every accepted field and invoke semantics belonging to the selected contract. Historical v1.0.0 validation retains its released stable format, source-order, payload/native consistency, timing, summary, diagnostic, reference, and collection rules. Current scripted-format semantics apply only to current inputs. Applying current validation indiscriminately to all versions is prohibited. Rejected v0.0.0 inputs must never be relabeled or passed to current stable capability validation as a substitute for support.
5. Apply source integrity and portable destination safety before any output publication. Structural or semantic success alone does not prove canonical base64, declared length, hashes, source references, or a safe restoration plan. Later implementation may strengthen bounded processing and safety where consistent with existing documented constraints; any change rejecting formerly valid v1.0.0 workflows must be evaluated by the minor-version gate rather than hidden as historical validation.
6. Enforce the selected contract's semantic capability boundary, then invoke the allowed operation. Private adapters may supply runtime conveniences only after historical validation, retain the original identity/provenance and source envelope, and must not manufacture missing native facts or label derived defaults as captured source information.

P04 chooses the minimal internal organization (for example, isolated version-specific validators and a lossless internal adapter). Duplicating entire released executables is not required. Every accepted historical field must remain available to the operation; changes to private Go types are permitted because Go packages remain internal. A typed projection that drops accepted native content is not an acceptable adapter.

## Executable and input command matrix

`Accept` means complete structural, semantic, and source-integrity validation followed by the command's normal capability, strictness, and output rules. It does not promise success for malformed content, an invalid target, a failed filesystem precondition, or a document with an applicable fatal rendering/conversion concern.

| Executable | Cue JSON identity | `validate` | `inspect` | `restore` | Matching native `render` | `convert` to a supported native target |
|---|---|---|---|---|---|---|
| Released v1.0.0 | v0.0.0 | Reject identity | Reject identity | Reject identity | Reject identity | Reject identity |
| Released v1.0.0 | v1.0.0 | Accept | Accept | Accept | Accept for SubRip/WebVTT | Accept for SubRip/WebVTT |
| Released v1.0.0 | v1.1.0 | Reject identity | Reject identity | Reject identity | Reject identity | Reject identity |
| Future v1.1.0 | v0.0.0 | Reject unsupported identity | Reject unsupported identity | Reject unsupported identity | Reject unsupported identity | Reject unsupported identity |
| Future v1.1.0 | v1.0.0 | Accept historical stable contract | Accept historical stable contract | Accept | Accept for matching SubRip/WebVTT | Accept for existing targets; scripted targets follow P09/P10 |
| Future v1.1.0 | v1.1.0 | Accept current contract | Accept current contract | Accept | Accept for matching supported format | Accept under the approved four-format conversion matrix |
| Future v1.1.0 | Unknown or mismatched pair | Reject | Reject | Reject | Reject | Reject |

`encode` continues to accept native subtitle input rather than treating Cue JSON as a migration command. Existing SubRip/WebVTT selectors, aliases, encoding behavior, destinations, speaker-detection controls, and output safety remain compatible. The future executable adds the ratified ASS/SSA native selectors and emits the current output identity for all newly encoded native input. The tagged v0.0.0 executable remains unchanged; the table does not add commands to that historical executable.

An invalid or unsupported identity after an otherwise valid invocation exits 1 using stderr for the error and publishes no payload. Invalid command syntax or pre-execution requirements retain exit 2. Success, warnings, quiet/silent filtering, stdout payload separation, cancellation, overwrite preflight, and `--force` boundaries retain the v1 CLI contract. No new public migration command or target-schema option is introduced by S023 or promised by this matrix.

## Input identity, provenance, and output identity

Validation, inspection, restoration, rendering, and conversion must not silently rewrite the loaded instance's `$schema`, `schema_version`, `producer`, `format_support`, source envelope, raw source fields, or unknown content allowed by its exact schema. Inspection reports the supported historical identity and capabilities without substituting the executing binary's current schema. Unsupported v0.0.0 envelopes must not become successful inspection reports or gain stable native support by relabeling. An independent third-party `producer.version` remains valid when its document conforms to the selected contract; it does not select a schema or need to equal official Cueson versions.

Restore publishes original source asset bytes; native render and convert publish native subtitle bytes. None of these operations emits a rewritten Cue JSON instance. Target conversion projections are private derived representations and remain outside Cue JSON. If a later feature proposes Cue JSON migration or regeneration, it requires a separate specification describing identity, producer attribution, source preservation, and transformation losses.

Native encoding in the future release creates a new document with current v1.1.0 identity and official producer version 1.1.0. The source envelope contains exact captured assets and safe basenames only. Official new output is not stamped with a historical identity to trick an older consumer. This target follows the constitution's software/schema lockstep and is distinct from reading historical inputs.

The future `schema`, `schema --output PATH`, and `schema --version` commands continue to describe only the exact current embedded output contract. They are not historical-schema selection commands. `version`, current canonical schema identity, and the release tag must agree at publication. Startup/build lockstep checks compare the executable with its current embedded output schema; accepting a historical input is a separate capability and never triggers a false software-versus-input-version equality requirement.

## Minor-version release gates

v1.1.0 is conditional on compatible additions to the existing CLI and preservation of the released v1.0.0 contract. The following evidence is mandatory before P13/P14 close:

- The historical v1.0.0 schemas, fixture documents, positive/negative semantics, source restoration, existing rendering/conversion loss behavior, CLI selectors/options/aliases, streams, exit classes, and overwrite/strict behavior remain valid under the new executable for the commands promised in the matrix.
- A reviewed schema/model change inventory identifies every change to existing SubRip/WebVTT fields, required properties, enums, constraints, native/common ownership, diagnostics, and annotations. The release must not silently narrow accepted historical documents, reinterpret their existing fields, weaken integrity, or remove an existing CLI workflow. Additions for ASS/SSA must be isolated to their approved new branches.
- Version-specific validation genuinely uses the released contracts; success obtained by relabeling v1.0.0 inputs as 1.1.0 or comparing schemas with their identity constraints removed is insufficient. Historical schema hash proof and current emission proof are independent checks.
- Release notes and maintained compatibility documentation disclose that new official output targets v1.1.0 and exact-version consumers must explicitly add support for that identity. Test evidence confirms the released v1.0.0 executable rejects v1.1.0 documents, including new documents originating from SubRip/WebVTT. There is no promise that an older consumer accepts a newer schema, and adding historical input support to the new executable does not provide forward compatibility to old consumers.
- Any change incompatible with an established public guarantee blocks the minor release and requires an explicit major-version decision before implementation continues. The version/command matrix cannot be used to excuse a breaking change. If a supported downstream workflow requires historical output generation, chart that separately and resolve its public contract before claiming the minor-version gate passed; do not invent an undocumented schema-relabeling workaround.

The identity change accompanying a new lockstep release is an explicit new exact contract and does not mutate the historical contract. It requires downstream version negotiation; the word "additive" is not evidence of automatic consumer acceptance. S023 records this limitation now, and P13/P14 must verify and document the actual release behavior rather than declare compatibility from a version bump.

## Concrete future acceptance and verification ownership

| Case | Required observation | Owning future outcome |
|---|---|---|
| V01: Tagged valid v1.0.0 SubRip and WebVTT fixtures, including nontrivial native blocks, settings, tokens, and diagnostics | Every promised command validates by the v1.0.0 contract; restoration is byte-exact; rendering/conversion retain established semantics and loss behavior; the input identity and producer are unchanged | P04 integration with P11 |
| V02: Tagged valid v0.0.0 representative and envelope fixtures | Current v1.0.0 and future v1.1.0 commands reject unsupported identity before publication even when the common cue shape looks usable; the tagged v0.0.0 binary and artifacts remain unchanged | P04 integration with P11 |
| V03: Historical invalid v1.0.0 fixtures and deliberately cross-version capability values | Historical structural/semantic failures remain failures; a v1.0.0 document with envelope-only capabilities rejects; identity-swapping a v0.0.0 envelope does not create a valid stable document | P04 |
| V04: URI/version mismatch, unknown future version, unversioned alias, malformed field types, missing identities, JSON-looking content with a native extension | Deterministic Cue JSON failure, no fallback decoder, no stdout payload or filesystem replacement, and no schema network request | P04 and P11 |
| V05: Valid historical document from a third-party producer with a different software version | Accept without changing producer fields or using the producer version for contract selection | P04 |
| V06: Historical corrupt base64, wrong declared size/hash, unresolved references, or unsafe restore destinations | Full integrity/safety rejection before publication, including with `--force`; existing files remain unchanged | P04 and P11 |
| V07: Current ASS/SSA document and cross-format native-data mismatch | Validate only by the current structural and semantic branches; scripted data under a historical identity and wrong-format native extensions reject | P05 integration with P04/P11 |
| V08: Native encoding of each supported format by the release candidate | New output has exact current identity, current official producer, truthful capabilities, byte-exact source assets, and no original filesystem path or local identifier | P05 model rules, P11 CLI, P14 candidate proof |
| V09: Current schema discovery with historical support installed | `schema` and `schema --output` are byte-identical to the current canonical schema; `schema --version` and `version` agree with current candidate identity; historical registry resources do not change command emission | P04 and P14 |
| V10: Immutable resources and disconnected execution | Bundled v1.0.0 resources match their recorded hash and tagged bytes, both historical release copies remain unchanged, and supported historical/current validation works without network access or input-controlled schema retrieval | P04 and P14 |
| V11: Released v1.0.0 executable reads candidate v1.1.0 output, including newly encoded existing text formats | Reject exact identity; documentation explains the limitation and does not claim forward compatibility | P14, documented at P13 |
| V12: Existing valid/invalid v1 CLI invocations across native platforms | Existing aliases, quiet/silent, streams, codes, target matching, strict loss reports, output preflight, cancellation, and no-publication guarantees remain compatible; new scripted selectors are consistently available | P11 and P14 native release proof |

P04 owns test-first schema registry and historical semantic work. P05 owns the prospective current scripted schema/model and the reviewed branch/constraint change inventory. P11 owns end-to-end command matrix, classification, capability, stream, provenance, and publication-safety regressions. P14 owns exact candidate version promotion, immutable historical byte proof, current emitted schema equality, packaged native-platform proof, and released-old-executable rejection evidence. P13 owns compatibility documentation and the frozen public-contract decision. These are later implementation requirements, not checks already passed by S023.
