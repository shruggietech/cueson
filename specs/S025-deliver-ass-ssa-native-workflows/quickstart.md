# S025 validation guide

From the repository root, use Go 1.25 or newer and the installed repository verification toolchains. Windows launchers must preserve the established hidden-process guarantee.

1. Run `go test -count=1 ./internal/codec/scripted ./internal/model ./internal/schema ./internal/conformance ./internal/cli`.
2. Encode authored ASS/SSA corpus sources with `cueson encode <source> --output <model>`, then validate/inspect the model and restore it to a separate destination. Compare original/restored bytes.
3. Render the model to a separate native file and re-ingest it. Compare required native/common projection and canonical bytes rather than source identity.
4. Apply consistent owning Text/native-field/common timing edits and repeat rendering. Retained source captures remain original; restoration still reproduces the original asset.
5. Exercise malformed retained content, unsafe metadata, sub-centisecond edits and strict ambiguity. No output payload or destination replacement is accepted.
6. Run full root/nested tests/vet/static/vulnerability, fixed-work fuzz, six pure-Go builds, repository text/docs/brand and complete generated site/browser/build/artifact/dry-run gates.
7. After publication, wait for current-head CI and all first/second review findings to converge. The human owns final merge.

See [native contract](contracts/native-workflows.md) for detailed canonicalization and deferred downstream acceptance.
