# S027 hardening and evidence contract

## Owning boundaries

Independent scripted source-capture inspection preflights the encoded representation against its existing decoded 64 MiB limit before allocation. Declaration parsing uses a bounded split that can detect one excess field and reject beyond 128. Existing decoded integrity and exact capture checks remain binding; public source grammar is unchanged.

Complexity excess rejects complete operations. Known-invalid native/model/source ownership, unsafe active references, corrupt length/hash/base64, and privacy violations refuse without input mutation or partial rendered/converted publication.

## Fuzz contract

Every scheduled fuzz mutation is bounded and in memory. Static fixture reads at initialization are permitted; callback filesystem writes are prohibited. Direct all-format source construction preserves exact bytes, encoding observations, identity, unavailable timestamps, semantic summaries, diagnostics and native owners and validates before conversion.

Conversion invariants cover all distinct targets, deterministic bytes and complete reports, target reparse, immutable source/document, strict losses and fatal no payload. Structured adversarial scripted documents include genuine accepted controls and known-invalid owner/reference/capture/envelope/privacy cases, not blanket rejection.

## Evidence contract

The matrix decoder requires io.EOF after one JSON document, refusing both another value and malformed suffix. Named native conformance test covers both dialects' encode/validate/inspect/render/reparse/restore and strict/partial output safety under Windows/macOS/Linux CI. Existing governed corpus remains byte-identical. Portable grammar rows remain inapplicable to platform-specific claims. Stable/candidate proof rows remain explicitly #64/#65, with no development-hardening deferral left under #63.
