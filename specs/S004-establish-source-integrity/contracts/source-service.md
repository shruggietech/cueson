# Internal Source Service Contract

## Validated Decode Boundary

`internal/schema` exposes a typed decode operation that accepts UTF-8 JSON bytes, performs the canonical JSON Schema check, decodes `model.Document`, runs semantic validation, and returns the document only when every layer succeeds. The existing validation-only entry point delegates to this boundary.

Invalid UTF-8, malformed JSON, schema violations, unknown fields, and semantic violations are validation failures. Source-specific code never operates on an unvalidated document.

## Integrity Preparation

Preparation accepts a validated document and checks every asset in document order before filesystem output:

- canonical standard padded base64;
- decoded byte count;
- lowercase SHA-256;
- safe stored basename;
- portable basename uniqueness.

Preparation returns either the complete ordered validated set or one deterministic error. It does not open destinations and does not retain decoded copies of the complete bundle.

## Destination Planning

Planning accepts validated assets and restore options. It resolves exactly one destination mode, requires existing parent directories, rejects multi-asset single-file modes, and examines every destination with no-follow semantics.

The planner rejects links, directories, devices, non-regular entries, unapproved existing files, portable collisions, native path-equivalent collisions, and an asset targeting another bundle destination. No staging begins until the complete plan succeeds.

## Restoration Transaction

Restoration performs the following transaction for the complete bundle:

1. Decode each asset again into a same-directory exclusive staging file.
2. Flush, close, reopen, and verify each staged byte count and digest.
3. For absent destinations, publish with an exclusive same-filesystem link so a raced-in target is not overwritten.
4. For forced regular destinations, create an exclusive hard-link rollback copy, revalidate identity, and atomically rename the stage over the approved file.
5. Reopen and verify every final file.
6. Apply and read back metadata according to the selected mode.
7. Accept the bundle, then remove rollback copies and non-authoritative links.

Any controlled failure before acceptance rolls back prior commits in reverse order. Rollback and cleanup use recorded file identities so they do not erase external replacements. A cleanup failure after full acceptance is a warning; a rollback failure before acceptance is a runtime failure and retains recoverable backup material when necessary.

The contract does not claim cross-process transaction isolation or crash atomicity.

## Source Capture

Capture opens one source path without following a link, proves the handle is a regular file, records observable timestamps and provenance from that handle, then reads bytes. It returns one source asset with a safe basename, canonical base64, exact length and SHA-256, and no source path.

## Timestamp Adapter

The platform adapter receives an already-open regular-file handle. It captures timestamps before reads and applies supported timestamps only after final publication and byte verification. Every requested observation produces exactly one result: `restored`, `unsupported`, `unavailable`, or `failed`.

Windows attempts creation, modification, and access through file-handle APIs. Linux and macOS attempt modification and access through descriptor APIs. Linux birth-time setting and the S004 macOS creation-time baseline are explicitly unsupported. Readback comparisons are exact after any required, predeclared precision conversion.

## Error Classes

- Invocation and destination-mode preconditions map to command status 2.
- Document parsing, schema, semantic, integrity, filesystem, metadata, cancellation, and transaction failures map to status 1.
- Successful exact-byte restoration maps to status 0, with default-mode unsupported metadata emitted as warnings.
