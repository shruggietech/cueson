# Tasks: Non-Publishing Release Proof

**Input**: Design documents from `specs/S009-prove-release-pipeline/`

**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/](contracts/), and [quickstart.md](quickstart.md)

**Tests**: S009 explicitly requires test-first artifact verification, seeded invalid-artifact cases, a real foreground snapshot, and hosted proof.

**Organization**: Tasks are grouped by independently testable user story. The release verifier is a standalone repository tool and never enters the shipped product dependency graph.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it owns a different file and does not depend on incomplete work
- **[Story]**: Maps the task to a user story in [spec.md](spec.md)
- Every task names its implementation or verification path

## Phase 1: Setup

**Purpose**: Establish the isolated verifier boundary and fixed tool identities.

- [x] T001 Create the dependency-free standalone verifier module in `scripts/release-verify/go.mod`
- [x] T002 [P] Record exact GoReleaser and Syft tool identities plus the release-proof command contract in `docs/release-verification.md`
- [x] T003 Verify `/dist/` and Go build output exclusions remain complete in `.gitignore` without broadening ignored source paths

---

## Phase 2: Foundational Verification Model

**Purpose**: Define the target matrix and deterministic evidence model that every story uses.

- [x] T004 Write failing table-driven tests for the exact six-target catalog, filenames, archive formats, binary names, and evidence ordering in `scripts/release-verify/verify_test.go`
- [x] T005 Implement target, archive, checksum, SBOM, metadata, and evidence types plus the exact six-target catalog in `scripts/release-verify/verify.go`
- [x] T006 Implement strict verifier argument parsing, foreground diagnostics, and deterministic JSON evidence output in `scripts/release-verify/main.go`

**Checkpoint**: The verifier module builds and its target/evidence tests pass without any release artifacts.

---

## Phase 3: User Story 1 - Produce the Complete Release Candidate Matrix (Priority: P1) 🎯 MVP

**Goal**: Build exactly one correctly named archive for each supported target without publishing.

**Independent Test**: Validate the GoReleaser configuration, run one clean snapshot, and observe exactly six archives with the approved names, formats, executable suffixes, and required root members.

### Tests for User Story 1

- [x] T007 [US1] Extend failing archive tests for missing, duplicate, unknown, wrong-format, snapshot-suffix, unsafe-member, wrong-binary, extra-member, and schema-byte-drift cases in `scripts/release-verify/verify_test.go`

### Implementation for User Story 1

- [x] T008 [US1] Implement strict ZIP and tar.gz readers that reject unsafe or non-regular members and extract exact member bytes in `scripts/release-verify/verify.go`
- [x] T009 [US1] Implement matrix completeness, exact archive-name, four-member, binary-name, executable-mode, and canonical-schema verification in `scripts/release-verify/verify.go`
- [x] T010 [US1] Add the snapshot-only six-target pure-Go build, archive, legal-file, canonical-schema, source-timestamp, and disabled-release configuration in `.goreleaser.yaml`
- [x] T011 [US1] Validate `.goreleaser.yaml` with GoReleaser v2.18.1 and run the first foreground snapshot into `dist/`

**Checkpoint**: User Story 1 independently produces and structurally verifies the complete six-archive matrix with no remote mutation.

---

## Phase 4: User Story 2 - Verify Artifact Integrity and Contents (Priority: P1)

**Goal**: Prove checksum bijection, schema/software/source lockstep, per-archive SBOM identity, path hygiene, and compatible CLI behavior.

**Independent Test**: Run the verifier against the real snapshot, then run seeded negative cases for each artifact boundary and confirm all invalid candidates fail with actionable diagnostics.

### Tests for User Story 2

- [x] T012 [US2] Add failing checksum tests for missing, duplicate, unknown, self-referential, malformed, uppercase, unsafe-name, and digest-mismatch entries in `scripts/release-verify/verify_test.go`
- [x] T013 [US2] Add failing metadata and build-information tests for wrong target, CGO enabled, missing trimpath, absent or mismatched VCS revision, dirty source, wrong version injection, and local-path leakage in `scripts/release-verify/verify_test.go`
- [x] T014 [US2] Add failing SBOM tests for missing, duplicate, malformed, wrong-format, wrong-target, untraceable-source, path-leaking, and publication-claim cases in `scripts/release-verify/verify_test.go`
- [x] T015 [US2] Add failing host-probe tests for nonzero exit, stderr output, extra stdout, wrong executable version, wrong schema version, and emitted-schema drift in `scripts/release-verify/verify_test.go`

### Implementation for User Story 2

- [x] T016 [US2] Implement exact checksum-manifest parsing and archive-to-digest bijection verification in `scripts/release-verify/verify.go`
- [x] T017 [US2] Implement Go build-information, target, CGO, trimpath, VCS revision, dirty-state, and injected-version verification in `scripts/release-verify/verify.go`
- [x] T018 [US2] Implement SPDX JSON identity, archive association, source-version, target, and prohibited-publication-claim verification in `scripts/release-verify/verify.go`
- [x] T019 [US2] Implement structural and caller-supplied local-identifier scans across names, members, binaries, metadata, and SBOMs in `scripts/release-verify/verify.go`
- [x] T020 [P] [US2] Implement hidden non-interactive Windows process execution in `scripts/release-verify/process_windows.go`
- [x] T021 [P] [US2] Implement ordinary foreground non-Windows process execution in `scripts/release-verify/process_other.go`
- [x] T022 [US2] Implement compatible-host `version`, `schema --version`, and schema-byte probes plus truthful foreign-target skipping in `scripts/release-verify/verify.go`
- [x] T023 [US2] Configure one target-bound SPDX JSON SBOM per archive and one six-entry SHA-256 manifest in `.goreleaser.yaml`
- [x] T024 [US2] Verify the complete real snapshot and write accepted `dist/release-evidence.json` through `scripts/release-verify/main.go`

**Checkpoint**: User Story 2 independently rejects every seeded invalid candidate and accepts the complete real candidate with exact schema and source identity.

---

## Phase 5: User Story 3 - Reproduce the Proof Locally and in Hosted Automation (Priority: P2)

**Goal**: Apply the same complete release contract through a safe pull-request workflow and a documented local sequence.

**Independent Test**: Statistically inspect workflow authority, run its commands locally, and observe the hosted `Release proof / Non-publishing snapshot` check plus retained verified evidence on the pull-request head.

### Tests for User Story 3

- [x] T025 [US3] Add repository-contract tests for exact tool pins, release disablement, mandatory snapshot arguments, safe events, read-only permissions, immutable GitHub-owned actions, and absent secret or publication surfaces in `scripts/release-verify/policy_test.go`

### Implementation for User Story 3

- [x] T026 [US3] Add the read-only pull-request and manual release-proof job with exact tool installs, snapshot verification, PR-scoped concurrency, and successful-only short retention in `.github/workflows/release-proof.yml`
- [x] T027 [US3] Add release-verifier formatting and tests to routine repository gates and cache inputs in `.github/workflows/ci.yml`
- [x] T028 [US3] Complete the local installation, snapshot, verification, evidence, and non-publication guide in `docs/release-verification.md`

**Checkpoint**: User Story 3 has one reproducible local contract and one least-privilege hosted check using identical acceptance logic.

---

## Phase 6: Polish and Cross-Cutting Verification

**Purpose**: Reconcile architecture, release claims, Spec Kit evidence, and full repository quality.

- [x] T029 [P] Update hosted-delivery and release-candidate boundaries without changing product contracts in `docs/architecture.md`
- [x] T030 [P] Record the release-proof capability and pinned-tool/non-publication decision under `[Unreleased]` in `CHANGELOG.md`
- [x] T031 [P] Add the standalone verifier command and release-proof documentation link without claiming a public release in `README.md`
- [x] T032 Run Go formatting, repository text checks, root and standalone module tests, race detection, vet, GoReleaser validation, the full snapshot proof, workflow linting, `git diff --check`, UTF-8/BOM/line-ending checks, and mojibake scans
- [x] T033 Run Spec Kit convergence against `spec.md`, `plan.md`, and `tasks.md`; append and implement any remaining work before publication
- [x] T034 Confirm issue #11 and the `cueson Delivery` item remain at Stage `In progress` with Slice `S009` and default Status unused before publication
- [x] T035 Format the official pull-request body through `scripts/github-format`, publish it with `Closes #11`, read it back, verify rendering, and move the Project item to `PR review`
- [ ] T036 Wait for all hosted checks and first-round third-party reviews, address every finding, request at most one `@codex review` second round only when required, and stop for the operator's final merge ritual after all checks and reviews are satisfied
- [x] T037 [Review] Replace host-dependent Windows absolute-path detection with portable drive and UNC syntax checks in `scripts/release-verify/verify.go` and `scripts/release-verify/verify_test.go`
- [x] T038 [Review] Bind every packaged binary's public version surface to a release-only marker and verify that marker for foreign targets in `internal/version/version.go`, `.goreleaser.yaml`, and `scripts/release-verify/verify.go`
- [x] T039 [Hosted] Scope local-identifier inspection to semantic Go build information, release metadata, and SBOM JSON so arbitrary binary bytes cannot create false path matches in `scripts/release-verify/verify.go`
- [x] T040 [Review] Replace the common-root path list with host-independent structural detection for arbitrary drive, UNC, and POSIX absolute paths in `scripts/release-verify/verify.go`

---

## Dependencies and Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Starts after the approved S009 specification and planning artifacts exist.
- **Foundational (Phase 2)**: Depends on setup and blocks all user stories.
- **User Story 1 (Phase 3)**: Depends on the verifier model and creates the real matrix used by later stories.
- **User Story 2 (Phase 4)**: Depends on User Story 1's real archives and completes release-candidate acceptance.
- **User Story 3 (Phase 5)**: Depends on the accepted local proof so hosted automation never publishes an unverified process.
- **Polish (Phase 6)**: Depends on all user stories and owns convergence, publication, and review monitoring.

### User Story Dependencies

- **User Story 1 (P1)**: Provides the independently testable six-archive MVP.
- **User Story 2 (P1)**: Builds on User Story 1 archives but remains independently testable through synthesized fixtures and the real snapshot.
- **User Story 3 (P2)**: Uses User Story 2's verifier as its only artifact acceptance authority.

### Within Each User Story

- Negative tests are written before the corresponding implementation.
- The verifier model precedes archive parsing; archive parsing precedes GoReleaser acceptance.
- Checksums, build information, SBOMs, and path hygiene precede host execution and final evidence.
- Local proof precedes hosted workflow publication.
- Hosted findings are resolved before a second review is requested.

### Parallel Opportunities

- T002 and T003 can run while the verifier module is initialized.
- T020 and T021 own mutually exclusive platform files and can run in parallel.
- T029, T030, and T031 own separate documentation files and can run in parallel after implementation stabilizes.
- Research and final read-only audits can be delegated; the coordinating agent owns integration and foreground verification.

## Parallel Example: User Story 2

```text
Task: "Implement hidden non-interactive Windows process execution in scripts/release-verify/process_windows.go"
Task: "Implement ordinary foreground non-Windows process execution in scripts/release-verify/process_other.go"
```

## Implementation Strategy

### MVP First

1. Complete setup and the verifier model.
2. Write User Story 1 negative tests.
3. implement the snapshot-only configuration and archive verifier.
4. Run one real snapshot and validate the exact six-archive matrix.
5. Continue immediately into integrity and hosted proof because issue #11 is not closeable at the archive-only checkpoint.

### Incremental Delivery

1. Exact matrix and archive structure.
2. Checksums, build metadata, schema lockstep, SBOMs, and path hygiene.
3. Host-compatible public CLI probes and deterministic evidence.
4. Read-only hosted reproduction and artifact retention.
5. Architecture, changelog, documentation, convergence, publication, and bounded review.

## Notes

- `[P]` tasks modify separate files and have no unmet dependency.
- All product and validation commands run in the foreground.
- `dist/` is disposable ignored output and never committed.
- Syft output is semantically verified but is not claimed byte-for-byte reproducible.
- S009 never creates tags, releases, signatures, attestations, or production-domain state.
