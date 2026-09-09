<!--
Sync Impact Report
- Version: template -> 0.1.0
- Added principles: lossless preservation; schema/software discipline; common model/native fidelity; no silent loss; test-first format work; portable and truthful behavior; documentation and delivery authority
- Added sections: Contract Boundaries; Development and Delivery
- Templates reviewed: .specify/templates/spec-template.md, plan-template.md, tasks-template.md
- Follow-up: Ratify schema details through the planned v0.0.0 schema slice
-->

# Cueson Constitution

## Core Principles

### I. Lossless Source Preservation

Every supported encode operation MUST preserve the original source asset bytes and sufficient integrity metadata to restore and verify those bytes exactly. Parsing, normalization, speaker inference, markup processing, rendering, conversion, and future OCR MUST NOT rewrite the source envelope. Filesystem paths and local machine identifiers MUST NOT enter Cue JSON; stored file names are safe basenames only.

### II. Schema and Official Software Discipline

The Cueson JSON Schema and official Cueson executable are peer products. Every official release MUST embed and emit the schema version matching the executable version, and released schema artifacts are immutable. During pre-1.0 development, breaking contract changes MAY occur in a documented minor release. Third-party producer versions are independent from the Cue JSON schema version they target.

### III. Common Model Plus Native Fidelity

Cue JSON MUST expose normalized timing and semantic subtitle content without requiring consumers to decode original bytes or parse a source grammar. The common model MUST NOT replace format-native information needed for fidelity. Source bytes remain authoritative for restoration; normalized fields and OCR observations are derived representations with explicit provenance.

### IV. No Silent Loss

Unknown, malformed, unsupported, or non-representable source information MUST be preserved when practical and surfaced through deterministic diagnostics. Cross-format conversion MUST account for every known loss, and strict mode MUST prevent lossy output. Unsupported platform capabilities and missing codecs MUST be reported truthfully.

### V. Test-First Format Work

Parser, renderer, restoration, conversion, schema, and platform behavior changes require focused tests and representative fixtures before or alongside implementation. Accepted source fixtures MUST prove byte-exact restoration. Malformed input MUST never panic or disappear silently. Platform-specific metadata claims require native platform tests.

### VI. Portable and Secure Operation

Official releases target Windows, macOS, and Linux as first-class platforms and remain pure Go with `CGO_ENABLED=0` unless a recorded decision justifies otherwise. Subtitle and Cue JSON inputs are untrusted. Restoration must prevent traversal and unsafe overwrite, validate lengths and hashes, bound allocations where practical, and treat stored timestamps as untrusted input.

### VII. Documentation and Delivery Authority

Documentation, schema, CLI behavior, and tests MUST agree. Non-trivial product work uses Spec Kit and an analysis gate. GitHub Issues, milestones, pull requests, and the `cueson Delivery` Project are the planning authority. Human operators retain final merge, push, tag, release, and production-domain authority unless they grant a specific one-time authorization.

## Contract Boundaries

- The public v1 contracts are the CLI and Cue JSON Schema. Go packages remain under `internal/` until a separate public API is approved.
- Cueson-owned JSON property names use lowercase `snake_case`.
- A multi-asset source envelope is foundational so paired and binary formats do not require a new root model.
- Exact restoration and model-driven rendering are separate operations and MUST be described separately.
- OCR results are repeatable derived observations, never source truth.
- The current project specification is a working pre-release draft. Product invariants in this constitution control when draft prose conflicts with an explicit operator decision or a ratified project rule.

## Development and Delivery

- Repository-authored text uses UTF-8 without BOM. Line endings follow `.gitattributes`, including CRLF for PowerShell, command, and batch scripts.
- Markdown prose uses one logical source line per paragraph and one logical source line per list item.
- Work slices SHOULD combine compatible atomic issues when one implementation and verification story remains clear.
- Architecture-affecting decisions update `[Unreleased]` in `CHANGELOG.md`.
- Verification runs in the foreground and must complete before success is reported.
- Eligible pull requests receive at most two automated Codex review rounds. Dependabot is excluded unless explicitly requested.
- An AI agent MUST stop before push, merge, tag, release publication, or production `cueson.io` changes unless the operator explicitly authorizes that specific action.

## Governance

This constitution governs project implementation and delivery. Amendments require a documented rationale, an impact review against active specifications and templates, an entry in `CHANGELOG.md`, and operator approval. Version changes follow semantic intent: major for incompatible governance changes, minor for new or materially expanded principles, and patch for clarifications that do not change obligations.

Pull requests and reviews MUST verify applicable constitutional principles. Any intentional exception must be explicit, narrowly scoped, justified in the relevant issue or Spec Kit artifact, and approved by the operator when it crosses a protected authority boundary.

**Version**: 0.1.0 | **Ratified**: 2026-09-09 | **Last Amended**: 2026-09-09
