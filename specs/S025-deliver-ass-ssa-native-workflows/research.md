# S025 research

## Shared native parser and model authority

Decision: one internal scripted parser handles ASS and SSA profiles and reuses narrow wrappers around authoritative model grammar/scalar/timestamp/projection helpers. Final model/schema/integrity validation remains mandatory.

Rationale: the S024 privacy and capture review regressions make copied native validation dangerous. Parser framing creates ordered observed owners; model validation proves the whole retained view.

Alternatives considered: independent ASS/SSA parsers and duplicated validation were rejected for drift; source-only replay was rejected because common semantics and edits are required.

## Content selection

Decision: add explicit scripted-candidate rejection evidence to the registry selection boundary. Missing/mixed/unsupported dialect signatures never fall through to filename or a mismatched explicit selector.

Rationale: existing positive/negative evidence cannot distinguish unsupported scripted source from unknown content.

Alternatives considered: guessing from extensions and choosing one mixed signature violate content-first selection.

## Experimental capabilities and compatibility

Decision: retain valid schema_only input declarations and additionally allow experimental ingest/render/restore true with OCR false. New native output uses experimental. Executable/schema stays 1.1.0-dev and released 1.0.0 semantics stay frozen.

Rationale: loaded document capability observations need not be rewritten when executable codecs become available. Stable support requires downstream full conformance/release evidence.

Alternatives considered: changing all schema-only inputs or advertising stable from parser tests alone were rejected.

## Model-driven renderer

Decision: shared serialization validates framing, complete declaration transitions, scalar ownership, common timing and source integrity before publication. Recognized owners regenerate native syntax; only established inert raw records/encoded data replay. New Format insertion preserves subsequent captured-owner declaration requirements.

Rationale: raw capture is observation, not edit authority; canonical UTF-8/LF output must reparse with retained semantics and checked 64 MiB output.

Alternatives considered: silent comma/newline sanitation, raw-line precedence, precision truncation and unconditional replay are rejected.

## Corpus and matrix

Decision: authored paired dialect fixtures reuse manifest v1 roles/provenance. Extend the existing conformance matrix with shared scripted row IDs and explicit deferred operation evidence. Portable native/common/diagnostic expectations are audited independently; conversion/stable assertions remain pending.

Rationale: one shared inventory avoids an ungoverned testing dialect, and semantic cycles alone cannot catch paired parser/renderer errors.

Alternatives considered: generated unreviewed golden blessing and marking future conversion inapplicable were rejected.

## Independent research ownership

Parser research: s025_ingest. Renderer research: s025_render. Corpus/matrix research: s025_corpus. Coordinator owns schema/current capability, registry/CLI integration, optional corpus verifier, docs, integration and final verification. Each implementation boundary has exclusive files; helper API coordination precedes consumers.
