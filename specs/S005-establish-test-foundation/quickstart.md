# Quickstart: S005 Verification

Run commands from the repository root on a clean checkout of `S005-establish-test-foundation`.

## Focused Helper and Conformance Verification

```powershell
go test ./internal/testutil ./internal/conformance ./internal/schema ./internal/source
```

The manifest verifier must accept every declared root-corpus artifact, reject its unit-test invalid cases, and report no orphan. Cross-package conformance must accept the governed seed, reproduce its exact source bytes and digest, reject each malformed stage, and keep all portable expectations free of injected local sentinels.

## Repository Verification

```powershell
go test -count=1 ./...
go test -count=1 -race ./...
go vet ./...
go build ./cmd/cueson
git diff --check
```

Verify changed repository-authored text as strict UTF-8 without BOM and scan it for mojibake. Verify `.gitattributes` and `.editorconfig` assign normal UTF-8/LF treatment to manifest, README, and normalized expected JSON while authoritative source and malformed input subtrees remain byte-preserved.

## Bounded Fuzz Smoke

```powershell
go test ./internal/schema '-run=^$' '-fuzz=^FuzzDecodeCueJSON$' '-fuzztime=1000x' '-parallel=1'
go test ./internal/source '-run=^$' '-fuzz=^FuzzInspectEncoded$' '-fuzztime=1000x' '-parallel=1'
go test ./internal/source '-run=^$' '-fuzz=^FuzzValidateSafeBasename$' '-fuzztime=1000x' '-parallel=1'
```

Each command must complete in the foreground without panic or retained fuzz artifact. Mutation ordering is not a deterministic claim; a fixed input's behavior and every promoted regression are deterministic.

## Checkout-Normalization Proof

For every authoritative source and malformed artifact reported by the manifest, inspect Git attributes and require `text`, `eol`, and `whitespace` to be unset. Recompute each declared byte length and SHA-256 from the checkout and compare it with `testdata/manifest.json`.

After the slice has a local commit, materialize a separate checkout and rerun conformance verification there so Git applies the committed attribute rules rather than reusing working-tree bytes:

```powershell
$proofRoot = Join-Path ([System.IO.Path]::GetTempPath()) ("cueson-s005-" + [guid]::NewGuid())
git clone --no-local . $proofRoot
git -C $proofRoot checkout S005-establish-test-foundation
go test -C $proofRoot ./internal/conformance
git -C $proofRoot check-attr text eol whitespace -- testdata/fixtures/source-envelope/basic-lf/source/captions.srt testdata/malformed/cue-json/truncated/input/document.cueson.json
```

The conformance run re-verifies the committed manifest inventory, lengths, and hashes from the materialized checkout. Remove only the printed temporary proof directory after confirming the command target.

S005 does not run hosted platform jobs, create workflows, add native SRT or WebVTT codecs, or claim subtitle grammar coverage. Those outcomes remain assigned to downstream issues.
