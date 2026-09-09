# Initial Fuzz Boundary Contract

## Cue JSON decoding

The Cue JSON target accepts arbitrary bytes up to the harness guard, calls the complete canonical decode boundary, and returns normally on rejection. An accepted typed document is re-encoded and must pass the same complete decode boundary. Seeds cover representative input, empty input, malformed JSON, invalid UTF-8, multiple values, paired and unpaired Unicode surrogate escapes, structural violations, and semantic violations.

## Canonical base64 integrity

The source-integrity target accepts arbitrary strings up to the harness guard and calls the canonical encoded-byte inspector. On success, an independent strict decode must reproduce the byte count and SHA-256, and a repeated inspection must return the same result. Seeds cover empty, one- to three-byte encodings, binary bytes, malformed padding, whitespace, alternate alphabets, illegal characters, missing padding, and non-canonical trailing bits.

## Safe basename validation

The basename target accepts arbitrary strings up to the harness guard and calls portable source-basename validation twice to prove deterministic acceptance. Seeds cover ordinary ASCII, canonical Unicode equivalents, boundary lengths, multibyte names, traversal, absolute and drive forms, UNC and URI-like forms, separators, alternate data streams, controls, trailing characters, reserved device stems, superscript aliases, and invalid UTF-8.

The target does not assert that the schema's character-count limit and source validation's UTF-8 byte-count limit are equivalent; S005 does not resolve that existing product-policy distinction.

## Execution and promotion

Ordinary tests run all seed inputs. Each mutation target runs in the foreground for a fixed count with one worker. Callbacks use no restoration, filesystem mutation, native timestamps, clocks, network, or external processes. Panics are not recovered.

Useful minimized failures are promoted to human-named root fixtures with manifest provenance, integrity values, and deterministic expected outcomes. Package-local automatic fuzz corpus files are not a second source of truth.
