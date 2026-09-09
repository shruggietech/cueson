# Governed test fixtures

This directory contains the small, deterministic, redistributable corpus used by Cueson's fixture and conformance tests. `manifest.json` is the machine-readable authority for fixture identity, provenance, redistribution approval, byte characteristics, integrity, and expected outcomes. This README explains the workflow but does not override the manifest contract.

## Directory layout

```text
testdata/
├── README.md
├── manifest.json
├── fixtures/
│   └── <boundary>/
│       └── <case-id>/
│           ├── source/
│           └── expected/
│               └── bytes/
├── malformed/
│   └── <boundary>/
│       └── <case-id>/
│           ├── input/
│           └── expected/
└── fuzz/
    └── <boundary>/
        └── <case-id>/
            └── input/
```

`fixtures/` contains accepted source material and its expectations. `malformed/` contains human-named permanent rejection cases with stable stages and diagnostic fragments. `fuzz/` reserves byte-preserved input paths for a future retained seed that has a distinct reason not to live in the accepted or malformed corpus. Useful fuzz discoveries normally move into a minimized manifest-governed malformed or fuzz-regression case rather than remaining opaque package-local files.

Every regular payload beneath `fixtures/`, `malformed/`, or `fuzz/` must be declared exactly once by `manifest.json`. The manifest and README are repository metadata and are not self-listed. The canonical embedded example at `internal/schema/testdata/representative.cueson.json` is governed by the schema package and is not part of this root inventory.

## Manifest and provenance

The manifest is ordered, versioned JSON. Each fixture record has a stable portable identifier, purpose, classification, origin, explicit redistribution decision, ordered artifacts, and an acceptance or rejection expectation. Each artifact has a logical identifier and role, a slash-relative path below `testdata`, a media type, exact byte length, lowercase SHA-256, and declared encoding, byte-order-mark, line-ending, and final-newline characteristics.

Origin must be stated as project-authored, synthetic, derived, or third-party. Synthetic and derived material includes a reproducible recipe. Third-party material includes a stable source reference and retrieval date. Runtime source paths, checkout paths, user names, host names, and temporary-directory identifiers are not provenance and must never enter the manifest or generated expectations.

Every committed fixture requires redistribution status `approved`, its SPDX license identifier or repository-relative bundled license reference, explicit attribution, and a decision about whether `NOTICE` must change. Project authorship does not imply redistribution approval. Material with unknown or denied redistribution rights must remain outside the repository. Third-party material must be reviewed for license compatibility, and required attribution must be added to `NOTICE` before the material is committed.

## Byte protection

Repository-authored metadata, including this README, `manifest.json`, and normalized JSON expectations, uses UTF-8 without a byte-order mark and LF line endings. Authoritative files beneath `fixtures/**/source/`, `fixtures/**/expected/bytes/`, `malformed/**/input/`, and `fuzz/**/input/` bypass Git text, end-of-line, and whitespace normalization. Known subtitle extensions remain protected as a second line of defense. Matching EditorConfig rules prevent compliant editors from rewriting those authoritative byte paths while ordinary PowerShell, command, and batch scripts elsewhere retain CRLF.

The manifest records byte characteristics for review and verification; verification reads bytes without decoding, normalizing, or repairing them. Do not open and save an authoritative artifact with a text editor that may change its encoding, line endings, byte-order mark, trailing whitespace, or final newline.

## Contributor workflow

1. Create a human-readable, stable case identifier and place source or malformed bytes in the matching authoritative subtree.
2. Record a concise purpose, truthful origin, reproducible recipe or source reference, redistribution approval, license, attribution, and `NOTICE` decision.
3. Add normalized model or diagnostic expectations outside `expected/bytes/`; place any byte-exact rendered expectation inside `expected/bytes/`.
4. Compute each artifact's exact byte length and SHA-256 from the working-tree bytes, then add one manifest entry for every payload.
5. Run manifest verification and conformance tests before reviewing the diff. Never use an update-goldens mode to bless unexplained drift.
6. Inspect Git attributes for every authoritative artifact and confirm that `text`, `eol`, and `whitespace` are unset.
7. If a fuzz run finds a useful failure, minimize it where practical, promote it under a descriptive case identifier, document its provenance and expected result, and rerun the deterministic regression twice.

## Verification

Run focused fixture and conformance verification from the repository root:

```powershell
go test ./internal/testutil ./internal/conformance ./internal/schema ./internal/source
```

Inspect normalization behavior for representative metadata and byte paths:

```powershell
git check-attr text eol whitespace -- testdata/README.md testdata/manifest.json testdata/fixtures/source-envelope/basic-lf/source/captions.srt testdata/fixtures/source-envelope/basic-lf/expected/model.json testdata/malformed/cue-json/truncated/input/document.cueson.json
```

Metadata and normalized expectations report text treatment with LF. Authoritative source and malformed input report `text: unset`, `eol: unset`, and `whitespace: unset`. The manifest verifier independently recomputes every declared length and SHA-256 and rejects missing, extra, linked, non-regular, unsafe, colliding, unapproved, or drifted artifacts.

Run the complete deterministic repository checks before publication:

```powershell
go test -count=1 ./...
go test -count=1 -race ./...
go vet ./...
git diff --check
```

Bounded fuzz smoke commands and checkout-normalization proof are maintained in `specs/S005-establish-test-foundation/quickstart.md`. Fuzz callbacks do not restore files and must run in the foreground without retained artifacts. Promoted fuzz regressions become manifest-governed cases.

## Scope boundary

This corpus foundation does not implement or claim native SRT or WebVTT parsing, rendering, conversion, dialect, or grammar coverage. Format-specific fixture coverage begins with the later codec slices.

S005 does not add hosted CI workflows or prove cross-platform execution. Issue #8 owns that automation. Large or non-redistributable real-world corpora remain external, and local corpus paths must never appear in portable expectations.
