# S014 Quickstart

```text
go build -o ./dist/cueson ./cmd/cueson
./dist/cueson encode sample.srt
./dist/cueson restore sample.srt.cueson.json --output restored.srt
./dist/cueson render sample.srt.cueson.json --to srt
./dist/cueson render sample.srt.cueson.json --to srt --output rendered.srt
go test ./...
```

Expected: Cue JSON validates against embedded `0.1.0`, restore matches the source digest, render is deterministic LF SubRip, and existing destinations require `--force`.
