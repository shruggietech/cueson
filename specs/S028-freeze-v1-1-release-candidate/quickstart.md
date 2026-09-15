# Quickstart: Validate the frozen candidate

Run from repository root through verified hidden Windows launchers or hosted noninteractive shells.

```text
go test -count=1 ./internal/schema ./internal/model ./internal/cli ./internal/convert ./internal/conformance
go -C scripts/release-verify test -count=1 ./...
go -C scripts/docs-verify run . -repo ../..
go run ./scripts/github-format/main.go .
goreleaser check
goreleaser release --snapshot --clean --skip=publish
go -C scripts/release-verify run . -dist ../../dist -repo ../.. -version 1.1.0 -commit <full-clean-commit> -execute-host
```

Expected current version/schema is exact 1.1.0 with byte-identical canonical/immutable copy, stable official ASS/SSA capability and unchanged historical input/source bytes. Package verifier rejects dirty/source/identity/target/hash/member/unsafe data before accepting evidence. Old-consumer compatibility proof uses the pinned verified v1.0.0 executable. Hosted same-bundle native smokes supply Linux/macOS/Windows execution; arm64 proof is structural.

Run the complete pinned Site pnpm test gate in site/ (lint/generation, unit, browser, build and artifact/dry-run). ASS/SSA authored preview/navigation is available; published schema/download inventory still contains only existing released versions. These commands never publish a tag/release or deploy production.
