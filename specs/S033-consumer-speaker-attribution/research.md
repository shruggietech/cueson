# S033 Research

## Exact version and historical boundaries

Decision: stage1.2.0-dev, add byte-exact historical1.1.0 from released authority, retain1.0.0. Rationale: old schemas reject additions and typed models are repeatedly revalidated by restore/render/convert. Alternatives: mutating1.1.0 violates immutability; final1.2.0 overclaims release; Decode-only gates fail programmatic historical validation.

Released1.1.0 schema SHA-256:223b61cbcf6337167039268576b2c739564585fff10e6cc6076e47a03526a0f7. Preserve published schema routes/download metadata; stage a non-publishing development snapshot separately.

## Native export and loss accounting

Decision: neutral model omission diagnostics plus converter-specific deterministic losses, one per attribution and one per media declaration. Rationale: scripted ASS/SSA variant conversion bypasses shared common cue loss handling; a conversion-wide pass reaches all twelve edges. Alternatives: appending only common cue losses misses variant edges; embedding arbitrary consumer IDs in native labels changes provenance and cannot preserve independent identity meaning.

Private targets clear omitted consumer fields after original validation to avoid duplicate native loss reporting and stale media checks.

## Runtime evaluation and source fidelity

Decision: pure summary and diagnostics helpers, no stored Diagnostics/statistics mutation. Rationale: source-cue conflicts must not block byte-exact restore, while invalid consumer intervals must fail before any publication. Alternatives: fatal source-cue conflict prevents restoration; clipping/recomputing creates undocumented loss.

## Identifier and timeline semantics

Decision: explicit frozen forbidden-code-point and boundary-whitespace sets from issue #82, Unicode scalar counting, checked signed media endpoint arithmetic. Rationale: codepoint validation is portable without destructive normalization. Alternatives: ASCII-only restricts approved consumers; broad format-character ban damages legitimate text; string blacklists fail data-role boundaries.

## Inspection and process artifacts

Decision: optional counts-only inspection projection, established absent-field output unchanged; development snapshot policy and docs-verifier current-source markers updated explicitly. Rationale:1.1.0 current-source equality cannot remain true after a staged new contract. No stable publication/production artifacts change.
