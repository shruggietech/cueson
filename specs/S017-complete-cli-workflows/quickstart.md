# Quickstart: Complete CLI Workflows

## Validate Cue JSON and native subtitles

```text
go run ./cmd/cueson validate document.cueson.json
go run ./cmd/cueson validate captions.srt
go run ./cmd/cueson validate --format vtt captions.vtt
```

Successful validation writes no stdout payload. Its success line and any warnings use stderr; `--quiet` removes the success line and `--silent` also removes warnings.

## Inspect safely

```text
go run ./cmd/cueson inspect document.cueson.json
go run ./cmd/cueson inspect --json captions.vtt
```

The first command writes a concise structural report. The second writes exactly one report-version-1 JSON object. Neither output contains preserved source bytes, content text, user-controlled identifiers or messages, content hashes, or a caller path.

## Generate shell completion

```text
go run ./cmd/cueson completion bash
go run ./cmd/cueson completion zsh
go run ./cmd/cueson completion fish
go run ./cmd/cueson completion powershell
```

Each command writes one static script to stdout and changes no profile. Users deliberately evaluate or save the result using their shell's normal configuration mechanism.

## Verify invalid boundaries

```text
go run ./cmd/cueson validate --format cueson --encoding utf-8 document.cueson.json
go run ./cmd/cueson inspect --json malformed.json
go run ./cmd/cueson completion cmd
```

The first and third commands are invocation failures with exit code 2 and empty stdout. The malformed Cue JSON command is an accepted inspection operation that fails with exit code 1 and never falls through to native subtitle detection.

## Run implementation evidence

```text
go test ./internal/cli ./internal/conformance
go test ./...
go test -race ./...
```

The final delivery pass also runs format checks, bounded fuzz seeds, static analysis, vulnerability scanning, pure-Go target builds, documentation verification, and the development release proof recorded by the repository.
