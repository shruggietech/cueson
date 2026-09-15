# S024 research and decisions

**Date:** 2026-09-15

## Exact historical boundary

Decision: select exact current1.1.0-dev and historical1.0.0 URI/version pairs from the parsed unmodified object, then validate using the selected bundled structural and semantic authority.

Rationale: the baseline relevant model/schema/CLI/convert/source files match tagv1.0.0. Downstream typed model revalidation means semantic dispatch belongs in model validation as well as schema decoding. Historical identity/capabilities/provenance must survive every operation.

Alternatives: permissive semver identity, relabeling before validation, network schema loading and whole executable duplication were rejected. The installed compiler supports an explicit denying URL loader, while bundled references are registered locally.

## Current identity and proof

Decision:1.1.0-dev is the unreleased staged current contract and software version. #65 owns final1.1.0 promotion and immutable new release copy.

Rationale: S023 explicitly permits P05 staging; candidate/release tools currently hardcode1.0.0 and canonical-equals-released proof, so minimal current candidate assumptions must change without altering historical release resources/evidence.

Alternatives: final1.1.0 now would blur promotion ownership; retaining1.0.0 with new shape would mutate a released identity.

## Native ownership

Decision: shared typed scripted document/cue structures are exposed through separate matching ass and ssa branches. Declarations belong to format_declaration records; declaration IDs are record IDs. Constructed recognized records may omit declaration references/capture fields. Captured known records resolve the applicable section-local preceding declaration.

Rationale: avoids an independently ordered undeclared root array and source-observation fabrication. Physical sections/records own one contiguous ordering; styles/events/attachments/cues refer into it.

Alternatives: generic maps, deduplication, flattening, mandatory capture fields and source-byte-only wrapping fail native fidelity or construction.

## Common projection and privacy

Decision: model-level bounded native Text projection verifies logical hard/soft breaks, raw_text ownership, readable plain_text, drawing exclusion and derived evidence. Current scripted zero-dialogue models are valid with null media summaries.

Rationale: source grammar decoding cannot be required of common-model consumers. Producer declarations are not proof of correct projection or privacy. Known/unknown metadata and original source records are classified conservatively; unknown authoring safety ambiguity rejects, while explicit dialogue/font-family/embedded-content roles remain exempt content.

Alternatives: copying physical ASS Text into logical lines, fabricated timed units, broad path-string rejection in dialogue and self-declared safe flags fail ratified requirements.

## Capability boundary and CLI

Decision: register recognized nil-codec ass/ssa identities where needed for generic validation/inspection/restoration, with schema_only declared status and installed ingest/render false. Direct render/restore use the shared Cue JSON finishing boundary. Inspection reports loaded input schema identity rather than executing current identity.

Rationale: structural support, declared capabilities, installed native operations and stable release are separate facts. Historical command behavior and source integrity cannot bypass the shared boundary.

## Resolved research

Independent compatibility and native agents supplied file-level findings. Annotation count/digest/current-equals-released assumptions and format-specific payload fidelity tests require deliberate updates; immutable released hash assertions remain unchanged. All research questions are resolved for blocking analysis.
