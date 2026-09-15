# S027 validation quickstart

Run from repository root with the existing Go toolchain. On Windows use the verified hidden noninteractive command launcher or integrated terminal, as required by AGENTS.md.

```text
go test -count=1 ./internal/cli ./internal/model ./internal/convert ./internal/testutil ./internal/conformance
go test -count=1 ./...
go vet ./...
go test ./internal/convert -run=^$ -fuzz=^FuzzScriptedConversionCycle$ -fuzztime=1000x -parallel=1
go run ./scripts/github-format/main.go .
go -C scripts/docs-verify run . -repo ../..
```

Inspect an accepted ASS and SSA native file and its encoded Cue JSON with human and JSON output. Expect catalogue-approved tokens, truthful experimental capabilities, safe scripted counts, and no source/native text or private identity. Run historical command and text report snapshots unchanged.

Run all fixed-work fuzz targets named in .github/workflows/ci.yml, including new structured scripted ownership cases. Expect bounded memory-only mutation callbacks, target-reparse/strict/fatal/source invariants, and no panic. Native workflow conformance proves exact restoration and output safety; hosted CI supplies all three platform claims and six pure-Go builds.

Full final verification includes root/nested format/tests/vet/static/security, site unit/browser/build/artifact/drift and non-publishing hosted release proof. Evidence is recorded in verification.md and the current-head PR completion record. Publication of a release, tag or production site is outside this guide.
