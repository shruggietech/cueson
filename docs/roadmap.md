# Cueson development roadmap

**Status:** S023-S025 merged; S026 complete scripted conversion in progress

**Assessed:** 2026-09-15 against merged S025 main revision `6173698ff112386907600e5572b9a365e4383453`, all active issues, native dependencies, and milestone/Project state.

## Current delivery state

v0.0.0 and v1.0.0 are released and independently verified, their milestones/epics are closed, and stable SubRip/WebVTT native workflows are implemented. S021 launched the public documentation site and immutable schemas; S022 corrected independent IPv4/IPv6 verification. The production site remains deployed from reviewed S021 revision `46838fd5cc888b299a09b89da676db0005b00f16`; merging later repository changes does not deploy production.

At the fresh S023 assessment there were no open issues or new active arrivals. S023 publishes native [milestone v1.1.0](https://github.com/shruggietech/cueson/milestone/3), coordination [epic #51](https://github.com/shruggietech/cueson/issues/51), and sixteen atomic children. Native issue relationships, milestone membership and cueson Delivery Project remain planning authority. This prose and the [S023 issue-map snapshot](../specs/S023-plan-scripted-format-milestone/issue-map.json) provide traceability rather than duplicate custom Project metadata.

S023 merged through PR #68 and planning issues #52, #53 and #54 are closed. S024 merged through PR #69 at `86ea7ecf8cd9166f1603a528f21e4e1c214f9fa0`, closing #55 and #56. S025 merged through PR #70 at `6173698ff112386907600e5572b9a365e4383453`, closing #57, #58 and #59 with experimental corpus, ingest, rendering and restoration evidence. The refreshed 2026-09-15 S026 kickoff assessment found no new issues: epic #51 and atomic #60-#67 remain active. S026 implements the coherent ten-new-direction conversion outcome (#60/#61), verifies the complete twelve-direction graph, and preserves later full CLI/conformance and stable/release/publication gates. Native GitHub state controls subsequent transitions.

## Next release scope and compatibility gate

The next major delivery milestone targets **v1.1.0: stable documented ASS/SSA native workflows, preserved historical v1.0.0 input handling, verified official release assets and immutable public schema/site hosting**. It is a compatible minor-release target, not a major-version change. Any unavoidable established CLI/schema break blocks this target and requires an explicit major-version decision before proceeding.

Selecting ASS/SSA first is an explicit deviation from treating draft section order as execution priority: the [main working specification](Cueson-Project-Specification-v0.0.0.md) lists future families without ranking them. Two related text-native dialects share sections/styles/events, retained overrides, semantic dialogue, model rendering and one conversion verification surface. This reuses the source envelope, registry, loss reports, corpus, annotations and pure-Go proof without prematurely adding an XML profile stack or binary/OCR runtime. Every other candidate remains below; the ordering is a cohesion judgment rather than a promised effort estimate.

The [future ASS/SSA contract](formats/ass-ssa.md) pins a bounded ASS v4+/SSA v4 UTF-8 profile, accepted/rejected semantics, retained native structures and edit ownership. It does not promise a pixel/video renderer, legacy codepage decoding or arbitrary external-resource loading. Source fidelity and privacy remain binding: retain accepted native/source information or reject an unsafe whole operation, never silently sanitize raw fields or mislabel lossy conversion.

The [future exact-version compatibility contract](../specs/S023-plan-scripted-format-milestone/contracts/version-compatibility.md) requires historical 1.0.0 schema/semantics on the new executable, current exact output identity and version-aware capabilities. Current v1.0.0 and planned v1.1.0 reject v0.0.0 identity; the historical v0 binary remains unchanged. Old exact-version consumers reject new 1.1.0 documents and must explicitly add that contract. A constant bump or identity-relaxing schema comparison does not prove compatibility.

## Chronological execution slices

Each slice runs the installed Spec Kit/autopilot workflow end-to-end with one integrated implementation/review/verification story. S026 is active; later codes/groupings are provisional until their kickoff and must be reassessed after each merge and any new issue. Preserve atomic acceptance criteria if analysis requires a narrower reviewable implementation session.

| Order | Proposed slice | Atomic children | Complete outcome and shared verification |
|---|---|---|---|
| 1 | S023-plan-scripted-format-milestone | [#52](https://github.com/shruggietech/cueson/issues/52), [#53](https://github.com/shruggietech/cueson/issues/53), [#54](https://github.com/shruggietech/cueson/issues/54) | Governed roadmap/backlog, native contracts, exact historical compatibility matrix and clean analysis. |
| 2 | S024-extend-versioned-schema-model | [#55](https://github.com/shruggietech/cueson/issues/55), [#56](https://github.com/shruggietech/cueson/issues/56) | Local exact historical schema/semantic validation plus annotated scripted current schema/model, with old/new and immutable-byte proof. |
| 3 | S025-deliver-ass-ssa-native-workflows | [#57](https://github.com/shruggietech/cueson/issues/57), [#58](https://github.com/shruggietech/cueson/issues/58), [#59](https://github.com/shruggietech/cueson/issues/59) | Provenance corpus, both ingesters/renderers, native/common edit cycles and byte-exact restoration. |
| 4 | S026-complete-scripted-conversion | [#60](https://github.com/shruggietech/cueson/issues/60), [#61](https://github.com/shruggietech/cueson/issues/61) | Twelve-direction four-format matrix, complete deterministic losses, precision and strict/fatal publication safety. |
| 5 | S027-harden-scripted-cli-conformance | [#62](https://github.com/shruggietech/cueson/issues/62), [#63](https://github.com/shruggietech/cueson/issues/63) | Shared CLI catalogue/input integration, executed conformance/fuzz/limit/privacy evidence and native-platform gates. |
| 6 | S028-freeze-v1-1-release-candidate | [#64](https://github.com/shruggietech/cueson/issues/64), [#65](https://github.com/shruggietech/cueson/issues/65) | Accurate frozen docs/stable profile, exact 1.1.0 identity/immutable copy and packaged non-publishing release proof. |
| 7 | S029-publish-verify-v1-1 | [#66](https://github.com/shruggietech/cueson/issues/66) | Explicitly authorized official tag/release and independently verified public artifacts. |
| 8 | S030-publish-v1-1-public-schema-site | [#67](https://github.com/shruggietech/cueson/issues/67) | Reviewed artifact and explicitly authorized exact-main deployment, verified docs/downloads/three immutable schemas and milestone closure. |

S025 is the largest candidate code slice. Keep corpus, ingest and rendering together so issues close on executed evidence. If S023/native implementation analysis shows the variants or retention surfaces cannot be reviewed coherently in one session, split into independently complete native outcomes and allocate fresh slice codes; do not close unfinished children or preserve these provisional numbers at the expense of verification.

## Atomic hard dependencies

The native blocked-by graph is acyclic and transitively reduced. Same-slice coordination is not a reciprocal hard blocker. Completed S022 issue #49 remains the first prerequisite; the sixteen children are native sub-issues of epic #51.

| Atomic issue | Direct native blocker | Outcome |
|---|---|---|
| [#52](https://github.com/shruggietech/cueson/issues/52) | [#49](https://github.com/shruggietech/cueson/issues/49) | Ratify the v1.1.0 milestone and reconcile delivery planning |
| [#53](https://github.com/shruggietech/cueson/issues/53) | [#52](https://github.com/shruggietech/cueson/issues/52) | Ratify ASS and SSA native preservation and semantic contracts |
| [#54](https://github.com/shruggietech/cueson/issues/54) | [#52](https://github.com/shruggietech/cueson/issues/52) | Ratify minor-release schema and CLI backward compatibility |
| [#55](https://github.com/shruggietech/cueson/issues/55) | [#54](https://github.com/shruggietech/cueson/issues/54) | Implement exact historical-schema input compatibility |
| [#56](https://github.com/shruggietech/cueson/issues/56) | [#53](https://github.com/shruggietech/cueson/issues/53), [#54](https://github.com/shruggietech/cueson/issues/54) | Add ASS and SSA schema and model structures |
| [#57](https://github.com/shruggietech/cueson/issues/57) | [#56](https://github.com/shruggietech/cueson/issues/56) | Establish the ASS and SSA fidelity corpus |
| [#58](https://github.com/shruggietech/cueson/issues/58) | [#55](https://github.com/shruggietech/cueson/issues/55), [#57](https://github.com/shruggietech/cueson/issues/57) | Implement ASS and SSA detection and native ingest |
| [#59](https://github.com/shruggietech/cueson/issues/59) | [#58](https://github.com/shruggietech/cueson/issues/58) | Implement ASS and SSA model-driven rendering |
| [#60](https://github.com/shruggietech/cueson/issues/60) | [#59](https://github.com/shruggietech/cueson/issues/59) | Convert ASS and SSA to SubRip and WebVTT with complete loss accounting |
| [#61](https://github.com/shruggietech/cueson/issues/61) | [#60](https://github.com/shruggietech/cueson/issues/60) | Convert existing text formats to scripted targets and between ASS/SSA |
| [#62](https://github.com/shruggietech/cueson/issues/62) | [#61](https://github.com/shruggietech/cueson/issues/61) | Integrate scripted formats across validation, inspection, help, and completion |
| [#63](https://github.com/shruggietech/cueson/issues/63) | [#62](https://github.com/shruggietech/cueson/issues/62) | Harden scripted-format conformance and fuzz boundaries |
| [#64](https://github.com/shruggietech/cueson/issues/64) | [#63](https://github.com/shruggietech/cueson/issues/63) | Freeze v1.1 contracts and publish-ready documentation |
| [#65](https://github.com/shruggietech/cueson/issues/65) | [#64](https://github.com/shruggietech/cueson/issues/64) | Prepare the exact v1.1.0 release candidate and proof |
| [#66](https://github.com/shruggietech/cueson/issues/66) | [#65](https://github.com/shruggietech/cueson/issues/65) | Publish and independently verify v1.1.0 |
| [#67](https://github.com/shruggietech/cueson/issues/67) | [#66](https://github.com/shruggietech/cueson/issues/66) | Publish the v1.1.0 schema and documentation on cueson.io |

Each child has the governed six-section body and independent acceptance/verification, one milestone, governed labels and exactly one Project item. Stage reflects actual progress, owning Slice is child planning text, and default Status stays unused. The spanning epic has no single Slice. S023-S025 children are Done; S026 children move through implementation and PR review, and downstream children stay Backlog until their open native prerequisites clear. Use Specced, Ready, In progress, Release verification or Done only with the corresponding evidence. Native assignees, parents/dependencies and milestone facts are not copied into custom fields.

## Subsequent roadmap candidates

| Later milestone | Main-plan families | Required foundation and explicit deferral |
|---|---|---|
| XML timed-text | TTML plus a pinned IMSC text profile first; SMPTE-TT and EBU-TT as explicit follow-on profiles | Namespace/unknown-node fidelity, style/layout inheritance, clock/frame/tick timing, profile conformance, safe XML processing, native schema and conversion losses. One generic parser cannot claim every XML family. |
| Additional legacy text | SAMI | HTML-like preservation, class/language/style/timing semantics, safe processing and common projection. |
| Broadcast binary | EBU-STL | Binary records, broadcast character sets, timecode/rate conventions, metadata, bounded parsing/corpus and exact restoration. Base64 is transport, not semantic support. |
| Bitmap semantic | PGS/SUP followed by paired VobSub IDX/SUB | Timed images/events/palettes, independently addressed assets/images, paired restoration, OCR provider provenance/quality and explicit derived conversion losses. Envelope-only experiments remain incomplete. |

OCR remains permanently in scope for complete bitmap semantics, as required by the main plan; it is not removed because this milestone is text-native. Later IMSC work must pin its version/profile, rather than claim all text/image profiles from a filename. Media demuxing, speech recognition, a video runtime, public Go library API, new brand design, signing/attestation and expanded installers are separate scope decisions rather than hidden blockers for the established release contract.

## Review, release and closure gates

Every slice has blocking analysis, foreground verification, responses to every actionable review/security finding and at most two automated Codex review rounds. Push/PR publication uses the relevant kickoff authority; the human owns each specific final merge. Tag/release and production schema/site deployment require their own explicit authority even after code is green. No S023 change grants future release/deployment authority.

Candidate readiness proves source/current schema/executable identity, historical input behavior, accepted native/common corpus cycles, exact restore, complete loss reports, native platforms, six package targets, six target-bound SPDX SBOMs and checksums. Official publication then independently verifies exact tag/release/public assets. Public hosting finally verifies the authorized reviewed deployment, DNS/trusted TLS/redirects/routes/downloads/revision and all three immutable schema hashes without a mutable latest alias.

After every human-confirmed merge, verify GitHub merge and exact-main CI, prune and synchronize clean state, reconcile closing children/parent/milestone/dependencies and Project fields, then reassess new issues before proposing the next slice. Close the v1.1.0 epic/milestone only after all sixteen scoped outcomes have independently completed and the public release/hosting verification has passed.
