# S014 CLI Contract

## Encode

`cueson encode [options] INPUT`

Options: `--output PATH`, `-o PATH`, `--force`, `--format auto|srt|vtt`, `--encoding NAME`, `--pretty`, `--stdout`, `--no-speaker-detection`, `--quiet`, `--silent`, `--no-color`, `--help`.

Default output is `INPUT.cueson.json`. `--stdout` and `--output -` emit one Cue JSON document to stdout. `--force` is invalid without a filesystem destination.

## Render

`cueson render [options] INPUT.cueson.json --to srt`

Options: `--to srt`, `--output PATH`, `-o PATH`, `--force`, `--strict`, `--quiet`, `--silent`, `--no-color`, `--help`.

Without `--output`, render writes canonical SubRip to stdout. `--output -` is an alias. `--force` is invalid for stdout.

## Output and exit behavior

Invalid invocation exits 2. Runtime, content, validation, capability, or publication failures exit 1 and leave no partial output. Stdout contains only requested data. Canonical render uses ordered sequences, `HH:MM:SS,mmm`, complete coordinates, raw payload, LF, one blank line between cues, and final LF.
