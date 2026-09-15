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

`fixtures/` contains accepted source material and its expectations. `malformed/` contains human-named permanent rejection cases with stable stages and diagnostic fragments. `fuzz/` contains byte-preserved mutation seeds that have a distinct reason not to live in the accepted or malformed corpus. Useful fuzz discoveries normally move into a minimized manifest-governed malformed or fuzz-regression case rather than remaining opaque package-local files.

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

The governed corpus now covers source-envelope integrity, Cue JSON rejection, native SubRip and WebVTT parsing and rendering, bidirectional conversion, strict loss handling, diagnostics, encodings, line endings, malformed input, and bounded fuzz seeds. `testdata/conformance-matrix.json` maps maintained format-contract rows to governed fixtures, focused tests, fuzz targets, platform evidence, or explicit inapplicability reasons.

S025 adds 86 project-authored synthetic ASS/SSA fixture records (28 accepted and 58 rejected), with paired dialect coverage and explicit MIT redistribution approval. Accepted scripts include BOM/CRLF capture, declaration spelling/order and duplicate inert additions, complete physical native retention, zero-dialogue and empty/drawing-only units, common text projection, native speakers and karaoke, inert attachments, malformed retained content, and content-role privacy exceptions. `expected/projection.json` separates ordered native style/event fields from common cue text, speaker and token observations; `expected/diagnostics.json` records severity, closed code and encounter position. These expectations were independently authored from the selected profile. They are not parser-generated golden output.

Each accepted source artifact is the byte-exact restoration oracle, with its own manifest byte count, SHA-256 and encoding/line-ending contract. Native rendering has separate canonical-byte expectations for the basic LF and BOM/CRLF pairs, deterministic repeated-render checks and complete native/common semantic reparse comparisons. Malformed retained content may ingest and restore while rendering refuses publication. Compact permanent malformed fixtures cover profile, grammar, timing, resource identity, active Effects, declaration amplification and S024 review regressions; executable tests exercise larger line/item/diagnostic/tag/output bounds without committing amplified payloads.

The shared `scripted-*` matrix rows map to the single ASS/SSA contract and require both dialects rather than duplicate row identifiers. Evidence kind `deferred` identifies applicable future assertions and names their owning issues. S025 native evidence does not mark future scripted conversion (#60/#61), full platform hardening (#63), stable freeze (#64) or reviewed release proof (#65/#66) passed. `inapplicable` remains reserved for assertions that do not apply to the row.

Large, private, or non-redistributable real-world corpora remain external. The optional maintainer corpus verifier may read such a caller-selected directory, but it must not modify source files, admit external bytes or absolute paths into the repository, or make the private corpus a hosted-CI prerequisite. Local corpus paths, usernames, hostnames, and machine identifiers must never appear in portable expectations or retained reports.

## Optional external corpus verification

Run the standalone verifier from the repository root with an explicitly selected private directory:

```text
go -C scripts/corpus-verify run . -root DIRECTORY
```

Add `-json` for a machine-readable report. One run accepts at most 10,000 regular files and applies the shared 64 MiB per-file limit. The verifier rejects a linked or non-directory root, every symbolic link or non-regular entry, unsafe relative identity, limit overflow, and cancellation. It uses bounded no-follow reads, the real CLI encode path, and native parser-renderer cycles without creating corpus, sibling, or repository files.

Output excludes the supplied root, absolute paths, source content, source bytes, and machine identity. It contains aggregate counts for SubRip, WebVTT, ASS and SSA and portable paths relative to the selected root when needed for a rejection. `accepted` counts successful validated encode operations; `render_rejected` separately counts accepted scripted preservation-only models with confirmed malformed owners, so it is a subset of `accepted`. Unexpected render-cycle errors fail verification. Recognized rejected scripted content counts as `rejected` even with a misleading extension. External bytes and results never enter `manifest.json`, the governed fixture tree, or hosted CI.
