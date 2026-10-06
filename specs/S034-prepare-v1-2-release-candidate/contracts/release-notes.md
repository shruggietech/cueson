# Cueson v1.2.0

- Store optional application-assigned speaker IDs and cue-contained speaker intervals independently from native subtitle labels. Names, UUIDs and bounded Unicode values retain their exact spelling; applications determine identifier scope.
- Store optional declared media duration and timeline alignment. Speaker intervals must fit their cues and supplied media boundaries; unavailable duration stays unavailable, and retained source-cue conflicts are reported without rewriting source data.
- Inspect consumer annotations through counts and validation state. Native rendering and all twelve conversions report omitted consumer data; strict mode refuses loss before output publication.

Exact 1.0.0 and 1.1.0 inputs remain supported through immutable local contracts. New output uses exact 1.2.0; older consumers reject it. Restoration reproduces original bytes. This release adds storage and validation for downstream diarization results; it includes no diarization or audio probing engine. Six pure-Go packages cover Windows, macOS and Linux amd64/arm64 with matching schema and legal bytes, checksums and target-bound SBOMs; packages are unsigned and unattested. Public schema and site hosting remain separately authorized.

Full changelog: https://github.com/shruggietech/cueson/blob/v1.2.0/CHANGELOG.md
