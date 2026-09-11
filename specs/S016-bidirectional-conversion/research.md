# Research: Deliver Bidirectional Subtitle Conversion

## Decision: Runtime-only loss reports

Define loss reporting inside `internal/convert` and do not add target-specific loss fields to Cue JSON or extend persisted parser diagnostics.

**Rationale**: A loss is a property of one requested source-target conversion and policy, not an intrinsic fact about the reusable source document. Keeping the report ephemeral avoids stale target claims and unnecessary schema/version churn.

**Alternatives considered**: Add a root Cue JSON loss array, which would persist target-specific state; reuse `model.Diagnostic`, which lacks target, kind, field path, and the correct source-truth meaning.

## Decision: Atomic validated loss entries

Represent each omitted, degraded, or ambiguous semantic occurrence as one validated entry containing a stable code, warning severity, kind, source and target formats, RFC 6901 source path, message, optional source order and cue ID, and bounded sorted context attributes.

**Rationale**: Atomic entries are machine-testable, preserve precise traceability, and avoid the current renderer behavior that aggregates unrelated field sets into one prose warning.

**Alternatives considered**: Free-form strings, which are unstable; one aggregate per cue, which hides multiple losses; severity escalation in strict mode, which incorrectly conflates policy with the nature of the loss.

## Decision: Structural deterministic ordering

Accumulate the complete report in a fixed compatibility-matrix rule order: document-level fields first, then body items by numeric source order, then fixed per-item rule rank and occurrence index. Never derive order from maps, lexical JSON-pointer sorting, or platform traversal.

**Rationale**: Warning and test stability must be identical across platforms and repeated runs even when several losses occur at one location.

**Alternatives considered**: Append order scattered across conversion functions, which is fragile; global lexical sorting, which places cue indexes and field names in surprising order.

## Decision: Strict policy before rendering and publication

Build and validate the entire loss report first. If strict mode sees any loss, return a typed error retaining the complete report and naming the deterministic first entry before invoking a renderer or publisher.

**Rationale**: This proves strict no-output behavior and lets tests and internal callers inspect every reason without partial target generation.

**Alternatives considered**: Fail on the first discovered loss, which hides later losses; render then discard, which performs unnecessary work and expands failure states; treat renderer conformance diagnostics as conversion losses, which mixes distinct concerns.

## Decision: Private target-native projection

Construct a fresh private target document from validated common semantics and explicit native mapping decisions without mutating the input. Render it canonically through the installed target renderer, then return only target bytes, ordered losses, and target-render diagnostics.

**Rationale**: Existing renderers correctly require target-native data. A private projection prevents callers from mistaking the original source envelope for a newly authored target Cue JSON document.

**Alternatives considered**: Invoke a renderer against mismatched source-native data, which fails or drops fields; return the target document publicly, which creates a misleading source-envelope contract; render source bytes, which is restoration rather than conversion.

## Decision: Two input classes with explicit `--from`

`convert` accepts `--from auto|cueson|srt|vtt`, defaults to auto, and supports native-source encoding controls. Auto classification recognizes valid Cue JSON before native formats and treats JSON-looking or JSON-named invalid input as a Cue JSON validation failure rather than falling through.

**Rationale**: `--from` describes conversion direction more clearly than reusing encode's `--format`, and Cue JSON must not accidentally be interpreted as subtitle text after a schema error.

**Alternatives considered**: Extension-only detection, which contradicts content-first detection; `--format`, which is ambiguous about source versus target; separate commands for source and Cue JSON, which duplicates the workflow.

## Decision: Stdout default and existing file transaction

Successful conversion writes target bytes to stdout unless a real `--output` is selected; `--output -` is equivalent. Real files reuse the established safe publisher, require `--force` for replacement, and remain unchanged on any pre-publication failure.

**Rationale**: This matches `render`, avoids surprising inferred filenames, and composes naturally in pipelines while retaining proven file safety.

**Alternatives considered**: Default `INPUT.<target>`, which risks unexpected files and ambiguous Cue JSON name stripping; a new publisher, which would duplicate tested transaction logic.

## Decision: Target-aware markup conversion

Scan payload syntax iteratively, preserve readable text and proven-equivalent balanced bold, italic, and underline spans, escape literal target grammar triggers, and emit atomic losses for font, class, language, ruby, voice, inline timing, or malformed source constructs that cannot retain meaning.

**Rationale**: Passing raw text between grammars can turn literal input into markup or cue boundaries. Plain-text-only conversion unnecessarily discards shared presentation semantics.

**Alternatives considered**: Raw payload passthrough, which permits semantic reinterpretation; plain-text flattening, which over-reports and loses equivalent markup; a DOM/CSS dependency, which is disproportionate and conflicts with the pure-Go boundary.

## Decision: Fatal unrepresentable grammar boundaries

Treat a zero-duration SubRip cue targeting WebVTT, invalid target timing, or payload that cannot be made parser-safe while retaining readable meaning as a typed fatal projection error in both modes. Do not invent timing or silently omit a cue.

**Rationale**: Permissive loss policy authorizes disclosed omission or degradation, not fabricated timing or output that the target parser reads differently.

**Alternatives considered**: Add one millisecond, which invents timing; drop the cue, which violates no-silent-loss and user expectations; emit invalid target syntax, which fails the parser-acceptance gate.

## Decision: Shared representability assessment without schema changes

Move current SubRip same-format representability checks out of CLI-local prose aggregation into reusable atomic analysis so `render` and `convert` use the same loss vocabulary where their target constraints overlap. Keep preserved source-conformance diagnostics owned by native renderers.

**Rationale**: One authority prevents render and conversion from disagreeing about unsupported metadata, speaker, token, OCR, placement, and ambiguous payload content.

**Alternatives considered**: Leave `subRipRenderLosses` in the CLI, which duplicates policy; move every conformance diagnostic into conversion, which violates codec ownership.

## Decision: Source-envelope integrity for Cue JSON input

Expose a read-only source-envelope integrity validator that checks safe names, canonical base64, decoded lengths, hashes, and portable collisions without creating destination plans or writing files.

**Rationale**: Issue #33 requires Cue JSON conversion to reject corrupted source truth even though conversion reads structured fields rather than restoration bytes.

**Alternatives considered**: Call restore preparation indirectly, which couples validation to output planning; skip integrity because bytes are not rendered, which violates the accepted input contract.
