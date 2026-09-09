# Feature Specification: Ratify Foundation Contracts

**Feature Branch**: `001-ratify-foundation-contracts`

**Created**: 2026-09-09

**Status**: Complete (round-one findings resolved and verified)

**Input**: User description: "Using Spec Kit and the shruggie-speckit autopilot protocol, make the README badges truthful, complete bootstrap bookkeeping, and ratify the v0.0.0 architecture, schema, CLI, and format boundaries without shipped product implementation; halt before push and public pull request."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Implement from one ratified baseline (Priority: P1)

As a Cueson maintainer, I can consult concise architecture, schema, and CLI contracts that resolve the provisional choices in the project draft so subsequent implementation issues begin from one consistent baseline.

**Why this priority**: The Go module, schema, source-integrity, and test-foundation issues all depend on these decisions. Leaving alternatives unresolved would move architectural guessing into implementation.

**Independent Test**: Review the ratified documents against the constitution and issue #3, then trace each acceptance criterion for issues #4 through #7 to an authoritative contract statement without relying on chat history.

**Acceptance Scenarios**:

1. **Given** the working project specification contains provisional alternatives, **When** a maintainer reviews the ratified schema contract, **Then** schema identity, release mutability, format keys, support semantics, source authority, and OCR observation cardinality each have one unambiguous rule.
2. **Given** the v0.0.0 executable will initially have incomplete capabilities, **When** a maintainer reviews the CLI contract, **Then** implemented command visibility, unknown command handling, missing-codec handling, output streams, overwrite behavior, and exit codes are explicit.
3. **Given** implementation issues #4 through #7, **When** their scopes are compared with the architecture contract, **Then** every package boundary and deferred capability has a documented owner and no shipped product code is required by this slice.

---

### User Story 2 - See truthful repository status (Priority: P2)

As a repository visitor, I can distinguish current bootstrap state from planned CI and release capabilities without interpreting broken or misleading live badges.

**Why this priority**: The README is the first public status surface and must not imply that workflows or releases exist before their tracked implementation slices are complete.

**Independent Test**: Render the README from the branch and confirm each badge is immediately resolvable or explicitly static, while the surrounding status statement continues to describe the repository as pre-release.

**Acceptance Scenarios**:

1. **Given** no CI workflow or public release exists, **When** a visitor opens the README, **Then** CI and release badges communicate planned or unreleased state without requesting missing dynamic resources.
2. **Given** the repository contains its license and documentation, **When** a visitor follows those badges, **Then** each link resolves to the intended repository content.

---

### User Story 3 - Trust delivery status (Priority: P3)

As a project operator, I can see that the bootstrap is closure-ready and the contract-ratification slice is actively tracked without claiming that an unmerged badge fix, unpushed branch, or unopened pull request is publicly available.

**Why this priority**: Project state must reflect evidence, but remote tracking must not outrun protected push and pull-request authorization boundaries.

**Independent Test**: Read back issues #1 and #3 and their Project entries, confirming bootstrap closure evidence, the one intentionally pending README criterion, the active slice identity, and no public pull request for this branch.

**Acceptance Scenarios**:

1. **Given** the bootstrap commit is present but the badge correction is not yet merged, **When** bootstrap bookkeeping is reconciled, **Then** issue #1 records completed criteria and evidence but remains open and In progress until post-merge verification can complete the README criterion.
2. **Given** this Spec Kit slice is underway, **When** issue #3 is read back, **Then** it identifies the contract-ratification slice as active but remains open until its change is merged.
3. **Given** autopilot reaches its protected boundary, **When** local verification and commit finish, **Then** no push or public pull request has occurred.

### Edge Cases

- A dynamic badge target does not exist yet: use an explicit static state badge rather than waiting for propagation or implying success.
- The future public schema domain does not resolve: retain its canonical versioned identifier while documenting repository, embedded, and release-artifact resolution before domain activation.
- A schema-recognized format has no native codec: validation and restoration capabilities remain distinct from ingest and render capabilities.
- A user invokes a command that is not shipped in the current executable: return an invocation error and do not advertise the command in help.
- A shipped command receives a schema-recognized format for which the required codec is absent: return a runtime capability error that distinguishes this case from an unknown format.
- A third-party producer emits a conforming document: its producer version may differ from the schema version it targets.
- A local task completes but its branch is not public: Project and issue wording must not claim a pull request or remote branch exists.
- A basename is harmless on the producing platform but maps to a Windows device or alternate data stream: reject it through the portable schema contract before restoration constructs any output path.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The slice MUST establish `docs/architecture.md`, `docs/schema.md`, and `docs/cli.md` as concise implementation baselines derived from the constitution and the working project specification.
- **FR-002**: The architecture baseline MUST assign boundaries for CLI wiring, the common model, schema handling, source preservation and restoration, codec handling, conversion, and future OCR providers while keeping all Go packages internal until separately approved.
- **FR-003**: The architecture baseline MUST keep exact source restoration independent from codec-driven rendering and MUST identify source bytes as authoritative.
- **FR-004**: The schema baseline MUST identify the in-development contract as `0.0.0` and its canonical immutable-release identifier as `https://cueson.io/schema/v0.0.0/cueson.schema.json`; the identifier is valid before the public domain resolves.
- **FR-005**: The schema baseline MUST state that the unreleased `0.0.0` contract may be refined, that every released schema is immutable, that documented pre-1.0 breaking changes require a minor version, and that patch releases remain non-breaking.
- **FR-006**: The schema baseline MUST use `subrip` and `webvtt` as the initial format keys and MUST define the completed v0.0.0 milestone capability state as `envelope_only`, with native ingest and render unsupported and exact restoration supported for valid source envelopes.
- **FR-007**: The schema baseline MUST define `format_support` as the official Cueson release capability declaration for a document's format, not as a claim about arbitrary third-party producer capabilities.
- **FR-008**: The schema baseline MUST require each cue to contain an `ocr_observations` array, allowing zero or more independently provenanced derived observations without replacing native cue text or source assets.
- **FR-009**: The schema baseline MUST preserve a multi-asset source envelope, portable safe basenames, exact byte length, SHA-256 identity, original bytes, and truthful timestamp provenance without recording an original filesystem path or machine identifier; basename validation MUST reject Windows reserved device names and NTFS alternate-data-stream syntax on every platform before restoration constructs an output path.
- **FR-010**: The schema baseline MUST preserve lowercase `snake_case` for Cueson-owned keys while retaining standards-defined JSON Schema keyword spelling.
- **FR-011**: The CLI baseline MUST expose only commands implemented by the current executable in command listings and help; at the first executable foundation this means help and `version`, with `schema` added by the schema slice and generic `restore` added by the source-foundation slice.
- **FR-012**: The CLI baseline MUST assign exit code `2` to invalid invocation or environment preconditions, including an unregistered command, and exit code `1` to runtime capability failures, including a shipped command that lacks a required codec.
- **FR-013**: The CLI baseline MUST distinguish unknown formats, schema-recognized formats lacking the required codec, and strict-mode loss rejection in diagnostics.
- **FR-014**: The CLI baseline MUST reserve stdout for payloads, send diagnostics to stderr, honor quiet, silent, color, literal-path, option-termination, and explicit-overwrite rules, and emit no decorative text in structured output.
- **FR-015**: The README MUST replace missing CI and release dynamic resources with truthful static state badges while retaining working license and repository-documentation links.
- **FR-016**: The working project specification MUST identify the ratified architecture, schema, and CLI documents as the implementation baseline and MUST remove or label any directly conflicting provisional wording touched by this slice.
- **FR-017**: The changelog MUST record the ratification decisions and the temporary static badge policy in the Unreleased section.
- **FR-018**: Bootstrap issue #1 MUST record completed criteria and verification evidence during this slice, while the README criterion, issue closure, and Project stage Done remain explicitly deferred until the badge correction is merged and verified on the default branch.
- **FR-019**: Issue #3 MUST be associated with slice `001-ratify-foundation-contracts`, moved to In progress during local work, and remain open until the eventual pull request is merged.
- **FR-020**: The slice MUST complete Spec Kit specification, clarification, planning, task generation, analysis, implementation, and local verification, then commit locally and halt before push or public pull-request creation.
- **FR-021**: The slice MUST NOT add shipped Go code, a root Go module, schema implementation, CI workflows, repository rulesets, tags, releases, or production-domain changes.
- **FR-022**: The v0.0.0 roadmap MUST include the multi-asset source envelope, integrity validation, and public generic exact restoration required by the ratified `restore_supported: true` capability while leaving codec-integrated round trips to the 0.x series.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All five contract decisions named by issue #3 (schema identity, format keys and support semantics, OCR cardinality, pre-1.0 compatibility, and CLI unavailable-command behavior) have exactly one normative answer in the ratified baseline.
- **SC-002**: Every acceptance criterion in issues #4 through #7 can cite at least one ratified architecture, schema, or CLI section without relying on chat history.
- **SC-003**: All dynamic badge requests remaining in the README refer to resources that currently exist; planned CI and release states require zero missing-resource requests.
- **SC-004**: The Spec Kit requirements checklist is fully passing and the analysis gate reports zero unresolved critical or high-severity inconsistencies.
- **SC-005**: Repository formatting, tests, encoding checks, and documentation link checks complete successfully with zero failures.
- **SC-006**: GitHub readback shows issue #1 open and In progress with only its merge-dependent README criterion pending, issue #3 open and In progress, and zero public pull requests created for this slice before the autopilot halt.

## Assumptions

- The pushed bootstrap commit `1188b5e` provides evidence for every issue #1 criterion except the README badge correction, which becomes complete only after this slice reaches the default branch.
- The operator's request authorizes routine issue and Project reconciliation but does not authorize pushing this feature branch or opening a public pull request.
- Static CI and release badges are temporary and will be replaced by live badges only after the corresponding workflow or release surface exists.
- The public schema domain remains a canonical identifier, not a currently required network endpoint.
- The completed v0.0.0 milestone includes the source-envelope restoration foundation even though native SRT and WebVTT codecs arrive later.
