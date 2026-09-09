# Research: Establish CI and Cross-Platform Build Gates

## Immutable action references

**Decision**: Pin every `uses:` action to a verified full commit SHA and retain the corresponding release tag in a nearby comment. Use `actions/checkout` v7.0.1 at `3d3c42e5aac5ba805825da76410c181273ba90b1`, `actions/setup-go` v7.0.0 at `b7ad1dad31e06c5925ef5d2fc7ad053ef454303e`, and `github/codeql-action` v4.38.0 at `b96794f015dfd88f77b49b1c93e0fa7110f94c63`.

**Rationale**: GitHub identifies a full commit SHA as the only immutable action reference. The selected commits were resolved from official release tags through GitHub's API and their commits have verified signatures. Version comments preserve auditability and update intent.

**Alternatives considered**: Major tags were rejected because they can move. Patch tags alone were rejected because they are not immutable. Vendoring the actions was rejected as unnecessary ownership and review surface.

## Workflow permission boundary

**Decision**: Declare `contents: read` at the CI workflow level. Declare `contents: read` and `security-events: write` only in the CodeQL workflow. Disable checkout credential persistence in every job. Do not use secrets, write-capable ordinary jobs, or `pull_request_target`.

**Rationale**: The repository currently defaults workflow tokens to write, so omission would violate least privilege. GitHub recommends read-only defaults with narrowly scoped elevation. CodeQL advanced setup requires result-publication authority, while normal checkout and analysis do not.

**Alternatives considered**: Waiting for issue #10 to change repository defaults was rejected because S006 workflows must be safe under current settings. Default CodeQL setup was rejected because issue #8 requires a versioned workflow and stable check contract. A shared write token was rejected as excessive.

## Runner and Go matrix

**Decision**: Use explicit GA runner labels `ubuntu-24.04`, `windows-2025`, and `macos-15`. Run the root suite natively on all three with the latest Go 1.25 patch. Run race detection on Linux amd64. Cross-build the executable on Linux for `windows`, `darwin`, and `linux`, each on `amd64` and `arm64`, with `CGO_ENABLED=0`.

**Rationale**: Explicit runner generations reduce surprise from moving `-latest` labels. Three native runs execute platform-selected timestamp tests where they belong. One reliable race job provides race coverage without conflating it with native behavior. Six temporary builds prove the supported portability matrix without publishing artifacts.

**Alternatives considered**: Moving `-latest` labels were rejected for primary compatibility evidence. Native arm64 runners were rejected because S006 needs build proof, not six native execution environments. Cross-compilation alone was rejected because the constitution requires native platform tests.

## Static and vulnerability tooling

**Decision**: Run Staticcheck v0.7.0 and govulncheck v1.8.0 with explicit module versions. Run actionlint v1.7.12 for workflow syntax and expression validation. Use the current stable Go toolchain for actionlint and vulnerability analysis, while Go 1.25 compatibility remains proven by the test and build jobs. Raise the project floor from Go 1.24 to Go 1.25 and upgrade `golang.org/x/text` from v0.31.0 to v0.39.0 because the first S006 scan found reachable findings whose fixes require those minimums.

**Rationale**: Staticcheck v0.7.0 is the first released line whose module declares Go 1.25. Govulncheck v1.8.0 is the current explicitly versioned release that loads Go 1.25 modules when built by the stable analysis runtime; v1.1.4 builds its package loader with the host's older Go launcher and fails before analysis in this environment. The initial S006 scan under Go 1.24.2 reported six reachable findings across standard-library call paths and the `golang.org/x/text` invalid-input path. The standard-library fixes were unavailable on the Go 1.24 line, while the text fix first available to this module is v0.39.0 and requires Go 1.25. After remediation, govulncheck v1.8.0 reports zero symbol-level vulnerabilities affecting the repository; it reports imported but unreachable GO-2026-5024 in `golang.org/x/sys/windows`, which does not violate the reachable-finding gate. A current stable runtime keeps vulnerability evaluation current while the separate compatibility jobs prove the minimum supported line.

**Alternatives considered**: `@latest` was rejected as non-reproducible. The official govulncheck action was rejected because its current implementation installs `govulncheck@latest` internally. A new tools module was rejected because explicit `go run module@version` commands avoid polluting product dependencies while retaining checksum verification. Suppressing reachable findings or retaining Go 1.24 with an affected text dependency was rejected as incompatible with untrusted-input handling and a fail-closed vulnerability gate.

## Stable check structure

**Decision**: Give every job a fixed human-readable name, with matrix dimensions rendered from a closed include list. Separate formatting, repository text, vet, schema/conformance, static analysis, vulnerability analysis, native tests, race detection, pure-Go builds, and CodeQL.

**Rationale**: Issue #10 can later bind rules to proven names rather than infer step state inside one opaque job. Closed include lists prevent new unsupported matrix combinations from appearing accidentally.

**Alternatives considered**: One monolithic job was rejected because it hides which contract failed and provides only one coarse required-check name. Dynamic path filters were rejected because documentation-only changes must still prove the stable check set exists.

## Controlled failing-run proof

**Decision**: Make the repository-text job reject a committed `.github/ci-failure-probe` marker. Publish the first draft pull-request revision with that marker, wait for the job and overall CI workflow to fail, then remove the marker, push the correction, and require all checks to pass before marking the pull request ready for review.

**Rationale**: This proves the real hosted workflow fails closed without keeping intentionally broken code, using a hidden background poll, or requiring a workflow-dispatch definition already present on the default branch. Draft status prevents automated review from evaluating the deliberately red revision.

**Alternatives considered**: A syntax-invalid workflow was rejected because it may create no check at all. A unit test failure was rejected because it would mix delivery proof with product behavior. `workflow_dispatch` was rejected for first-run proof because a newly added workflow cannot be dispatched reliably until it exists on the default branch.

## README badge transition

**Decision**: Change the static planned badge to the `CI` workflow badge in the same pull request and target the stable `.github/workflows/ci.yml` path. Verify the URL shape before merge and verify the live default-branch rendering during post-merge housekeeping.

**Rationale**: The workflow and its badge target become part of `main` atomically at merge. Keeping the static badge beyond that point would understate repository capability.

**Alternatives considered**: Updating the badge in a separate pre-merge change was rejected because it would recreate the previously observed missing-workflow badge. Deferring indefinitely to documentation cleanup was rejected because S006 owns the CI status transition.
