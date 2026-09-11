# CLI contract additions

## Encode

`cueson encode INPUT` recognizes WebVTT from a boundary-valid signature or explicit `--format vtt|webvtt`. It preserves the current default output, `--output`, `--force`, `--pretty`, `--stdout`, quiet, silent, no-color, diagnostics, and exit-code rules.

WebVTT accepts UTF-8 only. An explicit WebVTT format combined with a non-UTF-8 `--encoding` is an invalid invocation. An automatic selection that later proves to be WebVTT with an incompatible override is a runtime decode failure. Neither case publishes output.

## Render

`cueson render INPUT.cueson.json --to vtt|webvtt` writes canonical model-driven WebVTT to stdout by default or a caller-selected destination. It preserves `--output`, `--force`, `--strict`, quiet, silent, no-color, diagnostics, and exit-code rules.

## Output safety

All schema, encode, and render filesystem writes use the shared transactional publisher. A failed write, flush, verification, rename, or replacement leaves a new destination absent and an existing destination unchanged or restored.

## Exclusions

S015 exposes no `convert`, `validate`, `inspect`, or completion command and does not make WebVTT stable. Exact byte restoration remains `cueson restore`.
