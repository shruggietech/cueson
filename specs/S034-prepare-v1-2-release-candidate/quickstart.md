# S034 candidate adoption checks

From a built candidate, `cueson --version` prints `1.2.0` and `cueson schema` emits the matching immutable 1.2.0 schema. Use conforming Cue JSON containing `speaker_attributions` and optional `media_timing`; `validate` accepts consistent input and `inspect` exposes counts and evaluation state without IDs.

`restore` recreates source bytes unchanged. Native `render` and `convert` report omitted consumer fields; strict mode refuses before creating or replacing output. Invalid speaker intervals reject even if a retained source cue conflicts with declared media duration; source-only conflicts remain warnings.

Run `go test -count=1 ./...` and nested verification-module tests. Run repository release proof without `-development`, verifying version 1.2.0, the exact full source commit and immutable schema. Native hosts execute the accepted packaged binary and authenticate published 1.1.0 consumer bytes. Historical 1.0.0 and 1.1.0 inputs retain their loaded identity.

Local PR proof is review evidence. After an authorized merge, a fresh actual-main proof is required before the final public release decision. No tag, GitHub Release or production endpoint is changed by these checks.
