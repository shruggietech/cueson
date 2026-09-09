# Quickstart: Build CLI Foundation

## Prerequisites

- Go 1.24.0 or newer compatible toolchain
- Repository checkout on `S002-build-cli-foundation`

## Verify the module

```text
go fmt ./...
go vet ./...
go test ./...
go test -race ./...
go build ./cmd/cueson
```

Run the repository publication-tooling tests separately because `scripts/github-format` is an intentionally independent module:

```text
go -C scripts/github-format test ./...
```

## Exercise successful commands

```text
go run ./cmd/cueson version
go run ./cmd/cueson --help
go run ./cmd/cueson version --help
go run ./cmd/cueson --quiet version
go run ./cmd/cueson version --silent
```

Expected version stdout is exactly `0.0.0` plus one LF. Help lists only `version`. Quiet and silent do not suppress the requested version payload.

## Exercise invocation failure

```text
go run ./cmd/cueson schema
go run ./cmd/cueson --force
go run ./cmd/cueson version unexpected
go run ./cmd/cueson -- --help
```

Each invocation writes no stdout payload, writes an error and relevant usage to stderr, and the compiled executable returns exit code 2. `go run` wraps a nonzero program status, so use the built executable when asserting the exact process exit code.

## Repository gates

```text
go run ./scripts/github-format/main.go .
git diff --check
specify check
```

Confirm that source and generated documentation are UTF-8 without BOM, contain no mojibake, and follow `.gitattributes` line endings before publication.
