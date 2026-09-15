# S026 validation quickstart

Use the repository's verified hidden Windows launcher for console tools; these arguments are shown for reproducibility on headless or integrated environments. Go compatibility is 1.25.0.

```text
go test -count=1 ./internal/convert ./internal/model ./internal/codec/scripted ./internal/conformance ./internal/cli
go test -count=1 ./...
go vet ./...
go test ./internal/convert -run=^$ -fuzz=^FuzzScriptedConversionCycle$ -fuzztime=1000x -parallel=1
go -C scripts/docs-verify run . -repo ../..
```

Build the CLI and convert each new baseline source to each other target. Compare target bytes and the independently authored expected_conversion artifact, reparse through the target codec, and check timing/text/emphasis/defaults. The new conformance suite performs all twelve directions and verifies source immutability and deterministic reports.

Run strict/fatal publication cases on default stdout, explicit stdout, a new destination, and an existing destination with force. Expect zero target bytes and unchanged existing files. Run precision ties/neighbors/collapse/overflow, full native loss categories, variant alignment/colors/fields/attachments, historical text Cue JSON, malformed source, full source integrity, privacy, bounds, and cancellation cases.

Full CI parity additionally covers nested modules, Staticcheck, vulnerability scans, all meaningful governed fuzz targets, six pure-Go builds, actionlint, brand/docs integrity, and existing site tests. The official PR triggers hosted native-platform CI, CodeQL, and release proof without authorizing a release or deployment.
