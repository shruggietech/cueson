# Research: Establish Source Integrity and Restoration

## Decision: Keep hosted native matrix execution in issue #8

S004 supplies build-tagged Windows, Linux, and macOS timestamp implementations and tests, runs Windows behavior natively on the current host, and cross-compiles Linux and macOS sources with `CGO_ENABLED=0`. Issue #8 remains responsible for executing the complete hosted native matrix.

**Rationale**: Issue #8 is explicitly blocked by source and test foundations. Requiring its workflows as a prerequisite for S004 would create a dependency cycle. Cross-compilation is useful build evidence but is never described as native behavior proof.

**Alternatives considered**: Adding CI to S004 was rejected because issue #8 owns that independent outcome. Omitting non-Windows platform tests was rejected because the constitution requires native-selectable evidence before capability claims.

## Decision: Pin the newest Go 1.24-compatible native API module

Use `golang.org/x/sys` v0.41.0. It is the newest release whose module declares Go 1.24 compatibility; v0.42.0 requires Go 1.25. Move the already-used `golang.org/x/text` v0.31.0 from indirect to direct because portable collision keys call it directly.

**Rationale**: New platform calls belong in the maintained Go system-call modules instead of the frozen `syscall` package, while the repository's Go 1.24 contract must remain intact.

**Alternatives considered**: Frozen `syscall`, CGO wrappers, and a toolchain upgrade were rejected as less maintainable or outside S004.

## Decision: Return a typed document from the schema validation boundary

Add a typed decode operation that rejects non-UTF-8 input, parses exactly one JSON value, applies the embedded Draft 2020-12 schema and model semantics, and returns the validated document. Keep the existing validate operation as a delegating compatibility surface.

**Rationale**: The restore command must not parse the same untrusted document under a weaker path or duplicate schema/model logic in the source package.

**Alternatives considered**: Re-unmarshalling after validation was rejected because it duplicates parsing and can drift. Moving JSON validation into `internal/source` was rejected because it violates the ratified responsibility map.

## Decision: Prove canonical encoded bytes without whole-bundle decoded retention

For each asset, reject CR/LF, decode with strict standard base64, require re-encoding to equal the original spelling, and stream through a byte counter and SHA-256 digest. Check safe integer conversions and compare exact declared length and lowercase digest. Complete this preflight for every asset before destination inspection. Decode again into bounded-buffer staging files after the entire bundle passes.

**Rationale**: Strict base64 still ignores CR/LF. Re-encoding closes canonical-spelling gaps. A streaming preflight avoids holding every decoded asset simultaneously while retaining the required all-assets-before-output boundary.

**Alternatives considered**: Trusting schema `contentEncoding` was rejected because the validator does not prove decoded integrity. A single decode directly to destination staging was rejected because a later corrupt asset would open output before full validation. An arbitrary fixed size limit was rejected because the product contract defines none.

## Decision: Use the normative portable canonical-caseless key

Compute basename and destination collision identities with NFD normalization, default Unicode case folding, then NFD normalization again. Defensively revalidate every stored basename in `internal/source` even after schema validation. Compare planned absolute paths by this portable key, native cleaned equality, and `os.SameFile` when multiple existing paths can identify one object.

**Rationale**: This is the ratified cross-platform collision contract and prevents case- or normalization-insensitive filesystems from collapsing distinct assets.

**Alternatives considered**: Lowercasing, simple Unicode folding, NFC, and platform-only comparison were rejected because they do not implement the approved contract.

## Decision: Preflight and stage the whole transaction before publication

Resolve every destination, validate its parent or output directory, reject links and non-regular entries, enforce force policy, and snapshot existing-file identity before creating a stage. Then create exclusive mode-0600 stage files in each destination directory, stream decoded bytes while hashing and counting, synchronize and close them, reopen them, and verify length and SHA-256 independently. All stages complete before the first final destination is published.

**Rationale**: Full planning prevents intra-bundle replacement. Same-directory stages make final publication same-filesystem. Independent readback catches incomplete or corrupted staging before acceptance.

**Alternatives considered**: Direct final writes, per-asset plan/write interleaving, and truncation under force were rejected because they can destroy destinations or leave partial bundles after deterministic validation failures.

## Decision: Publish new files without replacement and preserve forced targets with hard-link rollback material

For an absent destination, create the final name as a hard link to the verified stage and then unlink the stage, so a raced-in destination causes failure rather than replacement. For a forced existing regular file, create and verify an exclusive same-directory hard-link backup while the original remains in place, revalidate the target identity, then replace it with the verified stage in one rename operation. Track committed identities and roll back controlled failures in reverse order. Refuse force when safe rollback material cannot be created. Preserve a recovery backup and report a compound conflict if external mutation makes rollback unsafe.

**Rationale**: Hard links preserve exact bytes and filesystem metadata without a missing-destination window. Identity-checked compensating rollback is the strongest portable multi-file guarantee available without claiming cross-filesystem atomicity or external transaction isolation.

**Alternatives considered**: Renaming originals aside was rejected because it creates a race window already identified during S003 review. Copy backups were rejected because they can lose metadata. Platform-specific multi-file transactions do not exist. Crash consistency and concurrent external mutation isolation remain explicitly unclaimed.

## Decision: Treat cleanup failure after acceptance as a warning

Once every destination has passed final integrity and metadata policy, failure to remove a stage name or rollback backup produces a warning with the exact recovery path and returns success. Before acceptance, rollback or cleanup failure is part of the runtime error and retains recovery material.

**Rationale**: The requested destinations are complete after acceptance, so reporting total failure would encourage unnecessary retry and possible overwrite. This follows the S003 post-commit cleanup decision.

**Alternatives considered**: Returning failure after accepted output was rejected as semantically false. Silently ignoring retained recovery files was rejected because it hides operator action.

## Decision: Capture timestamps from the same no-follow handle before content reads

Open a regular source without following links, capture native timestamp fields from that handle, then read and hash bytes through the same handle. Windows maps nonzero creation FILETIME to `windows_creation_time`. Linux uses `statx` birth time when its mask confirms availability and never maps ctime to creation. macOS uses birth time only when it differs from the documented ctime fallback; an equal value is labeled `ctime_fallback` rather than overclaimed.

**Rationale**: A single handle closes path-swap races and preserves the capture-before-read invariant on filesystems that update access time.

**Alternatives considered**: Path stat followed by a separate open was rejected because the object can change. `os.FileInfo.ModTime` alone cannot expose all required timestamps or provenance.

## Decision: Use four timestamp outcomes and exact readback

Each requested field is `restored`, `unsupported`, `unavailable`, or `failed`. A null source field is unavailable. A genuine platform limitation or unrepresentable precision is unsupported. Native API, identity, range, or readback errors are failed. Only exact native readback equality is restored; the schema has no precision metadata from which to justify a tolerance. No-metadata mode returns an explicit skipped mode and no fabricated timestamp results.

**Rationale**: Exact comparisons prevent silent timestamp loss. The four states separate missing source observations from destination limitations and operational errors.

**Alternatives considered**: Tolerance-based equality was rejected because guessed filesystem precision would overclaim. A fifth `skipped` timestamp state was rejected because skip is an operation mode, not an observation outcome.

## Decision: Restore timestamps through native handles after final byte verification

Windows uses no-follow handles with `GetFileTime` and `SetFileTime`, rejects reparse/directory/device handles, enforces 100-nanosecond representability, and verifies raw values after closing the setter handle. Linux and macOS use descriptor-based access/modification setters and native descriptor stat readback. Linux creation setting is always unsupported. macOS creation is captured but reported unsupported in this pure-Go baseline rather than using a pathname mutation with a symlink race. `ctime_fallback` is never applied as creation.

**Rationale**: Handle-based mutation protects the identity already verified. Exact Windows FILETIME conversion prevents overflow and truncation. Truthful unsupported status is safer than a path-based macOS creation-time overclaim.

**Alternatives considered**: Path-based timestamp setters were rejected because destination replacement can redirect them. CGO/Foundation was rejected by the pure-Go constraint. Treating Linux or Darwin ctime as birth time was rejected by the constitution.

## Decision: Roll back metadata failures according to policy

A native apply/readback `failed` result always fails and rolls back the transaction. `unsupported` warns and permits exact-byte success in default mode but fails and rolls back in strict mode. `unavailable` is informational and does not fail. No-metadata skips the entire timestamp phase. Metadata runs after final names and byte verification, and rollback restores forced originals or removes new destinations only after identity checks.

**Rationale**: Platform limitations are expected and user-controllable; operational failures are not safe to accept. Strict mode means every non-null requested timestamp must be exactly restored.

**Alternatives considered**: Treating operational failures as warnings was rejected because it hides broken supported behavior. Failing default mode for known unsupported creation timestamps would make generic restoration unusable on Linux.

## Decision: Advertise restoration only after the whole slice passes

Change both recognized v0.0.0 capability objects and model invariants to require `restore_supported: true` only in the completed S004 implementation. Keep status `envelope_only`, ingest false, render false, and OCR-required false.

**Rationale**: Generic restoration from an already-valid envelope is now a real public command, while codec capabilities remain absent.

**Alternatives considered**: Leaving restoration false after shipping was rejected as inaccurate. Changing ingest or render was rejected as unrelated and unimplemented.

## Primary references

- Go base64 canonicality: https://pkg.go.dev/encoding/base64
- Go file and hard-link operations: https://pkg.go.dev/os
- Go Unicode case folding and normalization: https://pkg.go.dev/golang.org/x/text/cases and https://pkg.go.dev/golang.org/x/text/unicode/norm
- Go platform syscall modules: https://pkg.go.dev/golang.org/x/sys/windows and https://pkg.go.dev/golang.org/x/sys/unix
- Windows file times and handles: https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-getfiletime, https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-setfiletime, and https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-createfilew
- Linux `statx` and `utimensat`: https://man7.org/linux/man-pages/man2/statx.2.html and https://man7.org/linux/man-pages/man2/utimensat.2.html
- Apple `stat(2)` birth-time fallback: https://developer.apple.com/library/archive/documentation/System/Conceptual/ManPages_iPhoneOS/man2/stat.2.html
