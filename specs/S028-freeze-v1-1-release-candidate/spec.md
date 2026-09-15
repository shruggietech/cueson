# Feature Specification: Freeze the v1.1.0 release candidate

**Feature Branch**: `codex/S028-freeze-v1-1-release-candidate`

**Created**: 2026-09-15

**Status**: Specified

**Input**: Operator S028 autopilot kickoff, combining issues #64 and #65 after merged S027 revision `351036cc0c6972aa6f5953af591d4f49df718695`. Push and official PR publication are explicitly authorized; final merge, release/tag publication and production changes are excluded.

## Clarifications

### Session 2026-09-15

- Q: Does stable promotion accept the old unreleased development identity? A: No. The ratified registry supports exact 1.1.0 current and immutable 1.0.0 historical pairs; current examples may be promoted but no identity-relaxed development adapter is introduced.
- Q: Are lower-capability observations retained? A: Yes. Schema-only declarations and complete experimental native declarations remain valid observations under the exact current contract, while official current encode and private targets declare stable native capability. Loaded declarations are never rewritten to match installed capability.
- Q: Does source/candidate freeze activate public availability? A: No. Authored ASS/SSA preview documentation is included, but published download metadata and the two public schema routes stay unchanged until #67. Notes are prepared without tagging or releasing.
- Q: How is old-consumer incompatibility proved? A: Execute a digest/source-bound published v1.0.0 binary against fresh candidate output, with identity refusal, no stdout and no destination replacement; packaged current native workflows are exercised on each governed hosted amd64 platform.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Know the exact supported contract (Priority: P1)

A subtitle user or downstream consumer can read one coherent maintained contract describing the completed bounded ASS/SSA workflows, historical input handling, current output identity, and conversion losses before adopting the candidate.

**Why this priority**: Stable promotion must reflect executed behavior and explicit limitations, avoiding misleading support or forward-compatibility promises.

**Independent Test**: Audit maintained documentation, annotated schema examples and generated preview content against every approved native/conversion acceptance row, with all authored links/examples and generated-content checks passing.

**Acceptance Scenarios**:

1. **Given** accepted ASS v4+ and SSA v4 sources, **When** the user checks support, **Then** documentation distinguishes semantic dialogue, retained styles/overrides/karaoke/drawings/attachments, canonical model rendering and exact original restoration, including unsupported processing and bounded UTF-8 input.
2. **Given** released v1.0.0 documents and newly encoded v1.1.0 documents, **When** the user reads compatibility guidance, **Then** exact historical acceptance and old-consumer rejection of new output are explicit without relabeling or mutating source truth.
3. **Given** the completed twelve-direction conversion graph, **When** the user selects a target, **Then** losses, defaults, centisecond precision, strict refusal and fatal no-publication behavior are documented consistently.

### User Story 2 - Use a coherent stable candidate (Priority: P1)

A user runs the candidate and receives exact 1.1.0 version/schema identity, truthful stable ASS/SSA capabilities within the ratified profile, and compatible historical behavior.

**Why this priority**: A stable release identity is a tested contract across all commands and packages, not a version-string change.

**Independent Test**: Exercise all four native formats, schema discovery, validation, inspection, restoration, rendering and conversion under the promoted candidate, including historical positive/negative and output-safety cases.

**Acceptance Scenarios**:

1. **Given** each supported native format, **When** it is encoded, **Then** the new document uses exact current 1.1.0 identity and official producer version, validated capabilities and preserved source bytes without local identifiers.
2. **Given** tagged valid or invalid v1.0.0 documents, **When** the candidate executes every promised command, **Then** the exact historical contract, source integrity, established streams/exit classes and safe publication rules remain in force.
3. **Given** schema discovery and a restored source asset, **When** the user compares artifacts, **Then** current discovery equals the immutable candidate schema and restoration equals the original bytes, while historical schemas remain unchanged.
4. **Given** a source/model mismatch, unsafe envelope or strict-loss conversion, **When** the command is forced toward an existing destination, **Then** it rejects at its owning boundary and preserves that destination.

### User Story 3 - Review one independently proven candidate (Priority: P2)

The operator receives an official PR identifying one exact candidate revision, with independently verified packages and concise publication-ready notes, ready for final review without publishing a release.

**Why this priority**: The eventual publication must use the reviewed revision and accepted artifact set, with each protected transition independently authorized.

**Independent Test**: Produce and accept the six target archives, six target-bound SBOMs and checksum manifest, execute supported host packages, and read back exact-head CI/review/publication evidence.

**Acceptance Scenarios**:

1. **Given** the exact reviewed revision, **When** candidate proof runs, **Then** every package binds that revision, target, pure-Go build and exact schema/legal bytes, with six archives, six SBOMs and a one-to-one checksum manifest.
2. **Given** publication-ready notes, **When** their body is checked, **Then** they state highlights and compatibility limitations, remain substantially shorter than the detailed changelog, and end with the versioned full-changelog link.
3. **Given** green exact-head native/quality/security/package checks and clean or fully addressed reviews, **When** the PR is handed off, **Then** independent #64/#65 acceptance evidence is recorded and the operator retains final merge, tag/release and production authority.

### Edge Cases

- Schema identity pairs are exact; development, v0.0.0, mismatched and unknown identities must not acquire undocumented acceptance.
- Third-party producer versions do not select the input contract; existing valid schema-only/native observations must not be reinterpreted as source capabilities.
- Styles, drawings, attachments, transforms and unsupported karaoke remain retained bounded facts; stable support does not promise pixel rendering, external resource loading or legacy codepage decoding.
- Malformed preservation-only records and unrepresentable conversion semantics retain diagnostics/losses or explicit whole-operation refusal.
- Frozen historical schema/legal/source bytes must not change while current identity/examples are promoted.
- Cross-compilation is structural proof, not execution on a foreign host; supported hosted amd64 execution must be identified separately from arm64 package inspection.
- Generated preview content must not advertise an unpublished release as publicly downloadable or activate a new public schema route.
- A required contract change incompatible with established v1.0.0 behavior blocks minor promotion and requires an operator decision.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Freeze the ratified bounded ASS v4+/SSA v4 profile and all twelve distinct four-format conversion outcomes against completed development evidence, including retained constructs, edit ownership, defaults, precision, losses and refusal boundaries.
- **FR-002**: Maintained architecture, compatibility, schema, format, CLI, conversion, installation/download guidance, support tables and examples MUST agree with executed scope and distinguish current candidate status from official released/publicly hosted status.
- **FR-003**: Current annotated schema and its examples MUST describe exact current identity, valid declared capabilities, native/common ownership and bounded extension semantics; every example MUST validate under its target contract.
- **FR-004**: Generated Site preview content MUST remain reproducible from maintained authored authority, with passing route/link/support/content checks and no premature release/download/schema activation.
- **FR-005**: Detailed changelog MUST record candidate additions, fixes and architecture decisions chronologically, without advertising deferred families or claiming an unpublished release is published.
- **FR-006**: Executable, emitted/embedded current schema, canonical identity, immutable candidate schema and planned version tag MUST agree exactly on 1.1.0; every new official native encode MUST declare the current identity and truthful stable implemented capabilities.
- **FR-007**: Candidate stable promotion MUST close every remaining stable platform/render evidence row with executed proof and preserve all accepted native fixture restoration, bounded conformance, privacy and output-safety outcomes.
- **FR-008**: Exact historical v1.0.0 structural/semantic and command behavior MUST remain supported without relabeling, field dropping, source mutation or network loading. New output MUST be explicitly rejected by the released v1.0.0 executable, including output originating from existing text formats.
- **FR-009**: Released v0.0.0 and v1.0.0 schema bytes and bundled historical resources MUST remain immutable; unsupported identity/capability and corrupt/unsafe inputs MUST fail before any payload or destination publication.
- **FR-010**: All six target archives, six target-bound SBOMs and checksum manifest MUST be accepted against the exact candidate revision, with schema/legal/source-revision/target/pure-Go/local-identifier checks and truthful host execution proof.
- **FR-011**: Publication-ready release notes MUST contain highlights and exact-version compatibility limitations, be substantially shorter than the detailed changelog, pass publication formatting and end with `Full changelog: https://github.com/shruggietech/cueson/blob/v1.1.0/CHANGELOG.md`.
- **FR-012**: Installed Spec Kit specification, clarification, requirements/domain audit, planning, ordered tasks and blocking analysis MUST complete before implementation; independent convergence and foreground CI-parity verification MUST complete before publication readiness.
- **FR-013**: The official non-draft PR MUST close #64 and #65, identify the exact reviewed head, have green hosted native/quality/security/Site/package gates and handle every review finding with at most one automatically requested second Codex round.
- **FR-014**: Native issue dependencies, parent/milestone and unique delivery Project membership MUST be preserved; #64 acceptance precedes #65 candidate acceptance, Project Stage tracks actual progress and default Status remains unused.
- **FR-015**: No tag, release, public schema activation, production configuration/deployment or final merge may occur in S028. #66/#67 and epic/milestone closure remain pending their independently verified authorized outcomes.

### Key Entities

- **Frozen profile**: Accepted native/common semantics, retained source fidelity and explicit unsupported processing, linked to executed matrix evidence.
- **Versioned contract**: Exact current or historical identity plus structural/semantic authority and declared/installed capabilities.
- **Candidate bundle**: Six target archives, six SBOMs, checksum manifest, immutable schema/legal hashes and exact source revision.
- **Publication-ready notes**: Short highlights and compatibility disclosures for a later authorized official release.
- **Delivery evidence**: Per-issue acceptance trace, local/native hosted checks, exact-head review convergence and Project lifecycle.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All #64 and #65 acceptance criteria have explicit independent coverage with zero unresolved material contract or analysis findings.
- **SC-002**: Every annotated schema example and maintained/generated documentation gate passes; all seventeen required scripted evidence rows have valid evidence with zero remaining stable deferrals.
- **SC-003**: All four supported native formats emit exact 1.1.0 current identity; all promised historical input commands pass their positive/negative contract matrix and the released old executable rejects new output.
- **SC-004**: Six archives, six target-bound SBOMs and one exact six-entry checksum manifest pass candidate verification on the exact reviewed revision; all governed native host executions pass with no foreign-execution claim.
- **SC-005**: Every retained historical schema hash and accepted source restoration hash is unchanged; no source path/local identity appears in candidate documents or public evidence.
- **SC-006**: Official PR publication is read back correctly, exact-head checks are green, and no actionable review remains, using no more than two Codex rounds total.
- **SC-007**: Zero tag/release/production/final-merge transitions occur during S028; the epic and remaining publication/hosting children stay open.

## Assumptions

- S027 merge `351036cc0c6972aa6f5953af591d4f49df718695` is the clean starting authority; the fresh kickoff assessment found only epic #51 and #64-#67 open with no new arrivals.
- #64 is unblocked by closed #63; #65 is coordinated in this slice after #64, with #66/#67 excluded because publication and hosting require separate protected authority.
- Existing current development identity is unreleased; current candidate fixtures/examples may be promoted while original source fixtures, provenance, released schema resources and earlier Spec Kit snapshots remain untouched.
- Existing six-target release configuration, native CI, version-specific registry and generated Site tooling supply the established proof contract; no new signing, installer, arm64 runner or runtime is implied.
- Explicit unattended autopilot authorizes routine spec/checklist/design decisions and explicit push/PR authority removes the usual pre-push pause; unresolved constitutional/public compatibility conflicts still require human judgment.
