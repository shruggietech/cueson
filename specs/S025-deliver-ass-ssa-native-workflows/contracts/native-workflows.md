# S025 native workflow contract

## Profile and source acquisition

Use the ratified [ASS/SSA grammar](../../../docs/formats/ass-ssa.md). Accepted dialects are ASS v4+ and SSA v4 with UTF-8 optional BOM and LF/CRLF. Reject preamble content before the first section, missing/mixed/unsupported dialects, NUL, bare CR and unsupported encoding. Source bytes and timestamp acquisition precede parsing and remain exact.

## Interfaces and ownership

The shared internal scripted package exposes detection facts, a bounded parser result and model-driven renderer without importing its parent registry. Registry adapters construct official current documents and native capability observations using the existing source/model/schema boundaries. The coordinator wires existing encode/render/shared-input commands; downstream conversion and full discovery enhancements stay unavailable.

Model helpers expose only existing authoritative native field/declaration/timing/scalar/attachment facts. No helper grants context-free privacy exemption. Native Name/Actor provenance remains present when heuristic speaker detection is disabled.

## Diagnostics and fatal behavior

Use bounded stable identifiers for unsupported/malformed dialect, declaration, style/event, override/karaoke/attachment and unsafe metadata/active content. Fatal acquisition/encoding/dialect/dialogue/privacy/complexity failures reject before payload publication. Safely retained malformed non-dialogue/declaration/attachment content has ordered diagnostics and is preservation-only: every native rendering mode refuses it. Strict rendering additionally refuses safely interpretable diagnosed ambiguity.

## Rendering and equivalence

Recognized records serialize from validated owning lexical/typed fields and common timing. Captured raw fields do not override edits. Constructed recognized owners omit unobserved captures and require complete fields. Declaration insertion/restoration must preserve every captured owner interpretation. Scalar canonicalization never modifies retained source observations; it rejects incompatible commas/newlines or centisecond precision instead of sanitizing/truncating.

Output is deterministic UTF-8 without BOM, LF with final LF, at most 1 MiB per physical line and 64 MiB total. Reparse comparison remaps position-based identities and compares dialect, occurrence ordering, declarations/unknown values, styles/events/assets and common timing/text/speakers/tokens. Canonical whitespace, numeric lexemes, introduced declarations and capture-only source observations may differ. Original source restore remains byte-exact independently.

## Verification and deferred evidence

Every accepted corpus source proves schema/model/source validity and exact restore; each renderable source and model edit/construction has independently expected semantic and canonical-byte evidence. Malformed/unsafe/limit/strict/fatal outcomes prove no partial destination or stdout and no source modification. Preserve old-format/historical regression behavior. Conversion assertions remain deferred to #60/#61, release-wide CLI/conformance to #62/#63 and stable/publication to #64-#67.
