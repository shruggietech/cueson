# Executable Documentation Contract

## Registered workflows

The maintained README must contain independently executable examples for build or source execution, encode, exact restoration, model-driven rendering, bidirectional conversion, validation, human and JSON inspection, and static completion generation.

## Isolation

- Examples use governed fixtures or safe temporary copies.
- Tests execute operations with temporary destinations and captured streams.
- No example creates or replaces a file beneath `testdata/`, `schema/releases/`, `brand/`, or another protected source directory.
- A strict conversion success example uses an explicitly loss-free fixture.

## Assertions

- Expected exit status is exact.
- Stdout is classified as empty, Cue JSON, native subtitle, inspection JSON, help, or completion content.
- Stderr contains only the documented success, warning, or error class.
- Exact restoration matches the governed input bytes.
- Rendered and converted output reparses as the requested format.
- Inspection output passes privacy scans.
- Every command, option, alias, selector, stream rule, and exit code in the executable surface has a maintained reference.

## Publication boundary

Documentation distinguishes the published envelope-only v0.0.0 artifacts from current v0.1.0 source behavior and the future v1.0.0 release. It does not claim that v1 binaries, immutable v1 schema, tag, public schema URL, or production website already exist.
