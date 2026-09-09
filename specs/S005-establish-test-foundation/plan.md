# Implementation Plan: Establish Fixture and Conformance-Test Infrastructure

**Branch**: `S005-establish-test-foundation` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/S005-establish-test-foundation/spec.md`

## Summary

Implement issue #7 as the test-foundation slice that later codecs and CI consume. Add a governed root fixture corpus and machine-verified provenance manifest, narrow repository byte-preservation rules to authoritative artifact subtrees, introduce domain-neutral comparison and path-leak helpers under `internal/testutil`, prove cross-package schema/model/source conformance through seed fixtures, and add bounded fuzz boundaries for Cue JSON decoding, canonical base64 integrity, and safe basenames. S005 deliberately adds no native subtitle codec, format grammar claim, CI workflow, release artifact, or production mutation.

## Technical Context

**Language/Version**: Go 1.24.0

**Primary Dependencies**: Go standard library; existing `golang.org/x/text` v0.31.0 for portable Unicode path identity; existing schema, model, and source packages

**Storage**: Version-controlled files beneath `testdata/`; no database, network storage, or generated persistent state

**Testing**: Go unit and integration tests, race detection, seed execution, bounded native fuzz smoke runs, manifest integrity validation, Git attribute inspection, and repository text checks

**Target Platform**: Portable Windows, macOS, and Linux test infrastructure, with native execution on the current Windows host and hosted platform execution retained by issue #8

**Project Type**: Pure-Go command-line repository with internal domain packages and repository-level test corpus

**Performance Goals**: Verify the small committed corpus in linear time with bounded memory; cap each fuzz callback input at 64 KiB as a harness resource policy without creating a product input limit; finish each smoke target after a fixed mutation count

**Constraints**: Exact fixture bytes; no implicit text normalization for authoritative byte artifacts; no absolute path or machine identifier in expectations; no symlink following; no inferred redistribution permission; deterministic Cueson-owned diagnostics; no external process, network, restoration output, or codec call inside fuzz callbacks; no CI, release, tag, or production-domain work

**Scale/Scope**: One central manifest, a small synthetic seed corpus, four deterministic malformed boundary cases, reusable helpers for future format corpora, and three initial fuzz targets

## Constitution Check

*GATE: Passed before Phase 0 research and after Phase 1 design.*

- **I. Lossless Source Preservation**: Pass. Manifest verification hashes authoritative source and expected byte artifacts without decoding, normalization, or mutation. Path checks prevent runtime source locations from entering portable expectations.
- **II. Schema and Official Software Discipline**: Pass. Conformance tests exercise the canonical embedded schema and current v0.0.0 representative without changing the schema or capability state.
- **III. Common Model Plus Native Fidelity**: Pass. Golden helpers compare the typed common model and ordered format-neutral observations without replacing source bytes or adding premature codec semantics.
- **IV. No Silent Loss**: Pass. Mismatches identify the exact logical comparison surface, malformed inputs preserve bytes and declared rejection stages, and timestamp outcomes remain explicit.
- **V. Test-First Format Work**: Pass. S005 establishes the fixture, malformed, golden, and fuzz foundation before SRT or WebVTT codec implementation begins.
- **VI. Portable and Secure Operation**: Pass. The design rejects unsafe or colliding manifest paths, links, orphaned files, unapproved redistribution, integrity drift, and machine-path leakage while keeping fuzz callbacks pure and bounded.
- **VII. Documentation and Delivery Authority**: Pass. S005 uses Spec Kit, issue #7 remains authoritative, push and PR publication are explicitly authorized, and merge, release, tag, hosted rules, and production changes remain operator-controlled.

Post-design re-check: Pass. The test utility package remains domain-neutral, the cross-package conformance layer owns domain projections, and byte-preserved subtrees are narrower than repository-authored manifest or expectation text. The schema/source 255-character versus 255-byte basename distinction is not redefined by S005; fixture paths obey the stricter byte-safe policy while product-contract reconciliation remains outside this test-infrastructure slice.

## Project Structure

### Documentation (this feature)

```text
specs/S005-establish-test-foundation/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── fixture-manifest.md
│   ├── test-helpers.md
│   └── fuzz-boundaries.md
└── tasks.md
```

### Source Code (repository root)

```text
.editorconfig
.gitattributes
CHANGELOG.md
docs/
└── architecture.md
internal/
├── conformance/
│   └── conformance_test.go
├── schema/
│   └── fuzz_test.go
├── source/
│   └── fuzz_test.go
└── testutil/
    ├── fixtures.go
    ├── fixtures_test.go
    ├── golden.go
    ├── golden_test.go
    ├── paths.go
    └── paths_test.go
testdata/
├── README.md
├── manifest.json
├── fixtures/
│   └── source-envelope/
│       └── basic-lf/
│           ├── source/
│           │   └── captions.srt
│           └── expected/
│               ├── diagnostics.json
│               └── model.json
└── malformed/
    ├── cue-json/
    │   ├── structural-extra-field/input/document.cueson.json
    │   └── truncated/input/document.cueson.json
    └── source-envelope/
        ├── hash-mismatch/input/document.cueson.json
        └── semantic-count/input/document.cueson.json
```

**Structure Decision**: `internal/testutil` owns only portable corpus loading, integrity, generic comparisons, and forbidden-token checks so it cannot create domain import cycles. `internal/conformance` maps Cueson model, diagnostic, restore, and timestamp values into those helpers. Package-local fuzz tests access the smallest unexported validation boundaries without exposing a public API. Root `testdata` owns the governed corpus; the existing embedded schema representative remains an official contract example outside the corpus and is explicitly documented rather than silently treated as an orphan.

## Complexity Tracking

No constitution violation requires justification.
