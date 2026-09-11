# Feature Specification: Freeze and Prove the v1 Contract

**Feature Branch**: `codex/S018-freeze-v1-contracts`

**Created**: 2026-09-11

**Status**: Draft

**Input**: Deliver issues #35, #36, and #41 through Spec Kit as work slice S018 by hardening the complete v1 behavior, enriching the canonical Cue JSON Schema for machine consumers, and making the maintained user documentation executable and release-ready without preparing or publishing v1.0.0.

## User Scenarios & Testing

### User Story 1 - Trust the complete v1 workflow under hostile and portable conditions (Priority: P1)

A user can run every documented SubRip, WebVTT, Cue JSON, restoration, rendering, conversion, validation, and inspection workflow knowing malformed or hostile input fails safely and supported inputs behave consistently across the release platforms.

**Why this priority**: A stable v1 cannot be released until the already implemented workflows have whole-system evidence across grammars, source fidelity, conversion loss, platform behavior, architecture targets, and security boundaries.

**Independent Test**: Run the governed corpus, bounded fuzz targets, hostile-input regressions, restoration and parser-renderer cycles, native platform tests, supported architecture builds, race checks, security scans, and path-leak checks, then confirm every documented contract row has traceable passing evidence.

**Acceptance Scenarios**:

1. **Given** every documented SubRip grammar row and WebVTT structure row, **When** the conformance matrix is evaluated, **Then** each applicable row maps to accepted, malformed, rendering, and conversion evidence with no unaccounted behavior.
2. **Given** hostile size, encoding, base64, timestamp, path, markup, or malformed-grammar input, **When** any applicable workflow processes it, **Then** processing remains bounded, never panics or traverses paths, never executes input, and never emits uncontrolled or privacy-sensitive output.
3. **Given** a supported source fixture, **When** encode, exact restore, model-driven render, conversion, validation, and inspection are exercised as applicable, **Then** source bytes remain authoritative, rendering follows the structured model, known losses are complete and deterministic, and strict conversion refuses every known loss.
4. **Given** Windows, macOS, and Linux release environments and all supported architecture targets, **When** the v1 verification matrix runs, **Then** native smoke tests and pure static builds succeed without claiming execution evidence for a foreign target.

---

### User Story 2 - Discover the Cue JSON contract from the schema (Priority: P2)

A downstream integrator can understand Cue JSON fields, units, examples, provenance, null behavior, capability semantics, and source-truth boundaries directly from the canonical schema without reverse-engineering the executable or depending on prose alone.

**Why this priority**: The schema is a public v1 contract and is consumed by editors, generators, validation interfaces, and documentation systems that require machine-readable context as well as normative constraints.

**Independent Test**: Recursively evaluate every consumer-facing schema property and definition for required annotations, validate every annotation example against its applicable schema fragment, emit the embedded schema through the CLI, compare bytes with the canonical schema, and prove the immutable v0.0.0 schema hash is unchanged.

**Acceptance Scenarios**:

1. **Given** any Cueson-owned root property or consumer-facing property reachable through a reusable definition, **When** a schema-aware consumer examines it, **Then** it has a non-empty specific description and any useful reusable definition has a meaningful title.
2. **Given** a semantic area, enumeration, or non-obvious constrained value, **When** its annotations are inspected, **Then** representative portable examples explain the relevant units, provenance, nullability, source authority, format boundary, or capability meaning.
3. **Given** complete and fragment-level examples, **When** automated validation runs, **Then** all examples satisfy their applicable constraints and contain no local machine identity, stale version, unsafe path, or invalid portable value.
4. **Given** the canonical development schema, **When** the executable embeds and emits it, **Then** the annotated bytes remain identical while the released v0.0.0 schema remains unchanged.

---

### User Story 3 - Install and use Cueson from accurate v1 documentation (Priority: P3)

A user can install Cueson, select a supported workflow, understand compatibility and loss boundaries, execute every README quick-start command, and interpret options, streams, statuses, diagnostics, overwrite rules, and strict behavior from one consistent documentation set.

**Why this priority**: The implementation is not a usable v1 until users can successfully operate it and integrators can distinguish stable public contracts from internal packages, exact restoration from rendering, and supported behavior from deferred production work.

**Independent Test**: Execute every maintained command example against governed fixtures, compare generated help and schema surfaces with their references, trace every format and compatibility claim to verification evidence, validate every link and document, and scan all authored text for encoding corruption and contradictory pre-release claims.

**Acceptance Scenarios**:

1. **Given** the README installation and quick-start guidance, **When** a user follows the encode, restore, render, convert, validate, and inspect examples, **Then** every command succeeds with the documented files, streams, statuses, and output semantics.
2. **Given** any implemented command or option, **When** the user consults help or the CLI reference, **Then** its arguments, aliases, streams, exit codes, overwrite behavior, strict behavior, and completion availability agree exactly.
3. **Given** the schema and format references, **When** an integrator evaluates the v1 contract, **Then** common and native data, support states, named encodings, grammar variants, diagnostics, conversion losses, compatibility rules, and immutable release-copy policy are complete and consistent.
4. **Given** installation, security, release-transition, or production-hosting guidance, **When** it is reviewed, **Then** it accurately separates source builds, released artifacts, future v1 publication, and post-v1 production-domain work.

### Edge Cases

- A corpus row applies to parsing but not rendering or conversion, and the traceability record must state the exclusion rather than invent evidence.
- A fuzz-discovered failure already has an equivalent governed fixture or exposes an intentional bound rather than a product defect.
- A schema property is internal to schema composition but still reachable by consumers through a public root.
- An annotation example is locally valid as a scalar but invalid within its parent object because of cross-field invariants.
- A description could be read as changing a validation constraint or making derived data authoritative.
- A documentation command would overwrite a fixture, emit binary or structured payload to the wrong stream, or depend on a machine-specific working directory.
- Native platform behavior differs legitimately because the platform cannot restore a timestamp class; the documentation and evidence must state the supported boundary.
- A valid hostile input is expensive but below an established limit, or rejected only after the command has accepted the operation.
- External corpus material cannot be redistributed and must be supported through a local maintainer hook without entering the repository inventory.
- A generated document, diagnostic, report, fixture, schema example, build record, or documentation transcript contains a caller path, username, hostname, drive, mount, or other local identifier.

## Requirements

### Functional Requirements

- **FR-001**: S018 MUST preserve all accepted behavior and acceptance criteria delivered by v1 issues #30 through #34 while closing the remaining verification and documentation gaps owned by issues #35, #36, and #41.
- **FR-002**: Every documented SubRip grammar row and WebVTT structure row MUST map to governed accepted, malformed, renderer, and conversion evidence where applicable, with explicit reasons for inapplicable cells.
- **FR-003**: Governed fuzz coverage MUST include input detection, SubRip blocks and timing, WebVTT blocks and settings, markup scanning, source-envelope validation, and parser-renderer cycles.
- **FR-004**: Every reproducible fuzz defect discovered in scope MUST become a deterministic governed regression fixture or focused regression test before the slice closes.
- **FR-005**: Hostile size, base64, timestamp, path, encoding, markup, and malformed-grammar inputs MUST fail safely without panic, path traversal, input execution, catastrophic resource growth, uncontrolled output, or source-content disclosure.
- **FR-006**: Bounded-input, allocation, diagnostic-count, output-size, and parser-progress limits MUST be explicit, deterministic, and verified wherever untrusted input can otherwise amplify work or output.
- **FR-007**: Exact restoration from every accepted source fixture MUST remain byte-identical and MUST remain distinct from model-driven rendering and conversion.
- **FR-008**: Parser-renderer and encode-restore cycles MUST prove the documented common-model, native-fidelity, source-order, diagnostic, and loss-accounting contracts without silently normalizing unknown content.
- **FR-009**: Generated Cue JSON and every privacy-bounded report or diagnostic surface MUST contain no original filesystem path, local machine identifier, prohibited source bytes, or unapproved raw content.
- **FR-010**: Windows, macOS, and Linux native smoke tests MUST exercise the platform behavior each host can truthfully claim, and the supported amd64 and arm64 target matrix MUST build with static pure-Go configuration.
- **FR-011**: Dependency, static analysis, race, vulnerability, CodeQL, integrity, repository-control, documentation, and non-publishing release-proof checks MUST remain accurately configured and green.
- **FR-012**: The repository MUST provide a documented maintainer hook for optional non-redistributable external corpus verification without admitting untracked corpus bytes to governed fixtures or CI.
- **FR-013**: The canonical development schema MUST provide a non-empty specific `description` for every Cueson-owned root property and every consumer-facing property reachable through reusable definitions.
- **FR-014**: Reusable public schema definitions MUST have meaningful titles when a title improves generated documentation, while internal composition helpers MAY omit titles when their enclosing public definition supplies the consumer meaning.
- **FR-015**: Schema descriptions MUST explain applicable units, null semantics, provenance, source-authoritative versus derived status, native-format boundaries, capability states, and compatibility meaning without changing normative constraints.
- **FR-016**: Schema `examples` MUST cover every top-level semantic area, every enumeration, and non-obvious constrained values including identifiers, versions, digests, base64 payloads, timestamps, timings, diagnostics, and native-format data.
- **FR-017**: Complete-object and fragment-level schema examples MUST validate against their applicable schema and MUST be portable, current, non-identifying, and internally consistent.
- **FR-018**: Deterministic automated coverage MUST reject missing, empty, invalid, stale, unsafe, machine-identifying, or semantically contradictory annotations whenever consumer-facing schema fields are added or changed.
- **FR-019**: Schema annotations MUST remain non-normative, MUST NOT weaken validation constraints or change existing accepted and rejected fixture outcomes, and MUST be retained byte-for-byte in the embedded and emitted development schema.
- **FR-020**: The immutable v0.0.0 schema copy MUST remain byte-identical to its published baseline.
- **FR-021**: The README MUST provide tested installation guidance and executable quick-start examples for encode, exact restore, render, conversion, validation, and inspection using supported repository fixtures or self-contained safe inputs.
- **FR-022**: Maintained CLI documentation MUST cover every command, option, alias, delimiter, stream rule, exit status, overwrite behavior, strict mode, quiet and silent behavior, and completion target implemented by the executable.
- **FR-023**: Maintained schema documentation MUST identify the v1 contract, common-model invariants, format-native extensions, support states, source envelope, annotations, compatibility policy, and immutable release-copy rule.
- **FR-024**: Maintained SubRip documentation MUST define the complete accepted grammar, tolerated variants, named encodings and aliases, diagnostics, native fidelity, rendering behavior, and conversion limitations.
- **FR-025**: Maintained WebVTT documentation MUST define every supported signature and header form, block, cue, timing form, setting, markup construct, ordering behavior, rolling-caption behavior, diagnostic, native fidelity rule, rendering behavior, and conversion limitation.
- **FR-026**: Maintained conversion documentation MUST provide a complete bidirectional loss and fatal-incompatibility matrix consistent with runtime loss codes and strict-mode behavior.
- **FR-027**: Maintained architecture, security, installation, compatibility, changelog, and release-process guidance MUST agree with the implemented v1 development behavior and MUST distinguish the future v1 release transaction from post-v1 website or production-domain work.
- **FR-028**: Every repository-authored command example MUST be executable by deterministic verification and MUST assert the promised status, streams, and output or artifact behavior without modifying governed fixtures.
- **FR-029**: Help/reference, schema/reference, format-matrix, link, Markdown, UTF-8, line-ending, mojibake, local-identifier, and executable-example consistency MUST be checked automatically.
- **FR-030**: The v1 public compatibility promise MUST cover the CLI and Cue JSON Schema while leaving packages under `internal/` outside that promise.
- **FR-031**: S018 MUST update `[Unreleased]` for material contract, verification, schema-annotation, and documentation decisions without declaring or dating v1.0.0.
- **FR-032**: S018 MUST NOT create an immutable v1.0.0 schema copy, change the development identity from `0.1.0`, prepare final v1 release notes, create or move a tag, publish a release or schema, mutate production configuration, close the v1 epic or milestone, or merge its own pull request.
- **FR-033**: The official pull request MUST close issues #35, #36, and #41, pass every current-head check, address every external review finding, request no more than one second Codex round, and stop for the operator's final merge ritual.

### Key Entities

- **Conformance Matrix**: The traceable mapping from each documented format and workflow contract row to accepted, malformed, rendering, conversion, fuzz, platform, or explicit inapplicability evidence.
- **Governed Corpus Entry**: A repository-approved input with portable identity, provenance, immutable bytes, expected behavior, and integrity metadata.
- **Schema Annotation**: Non-normative machine-readable context attached to a public schema definition or property, including descriptions, titles, and examples.
- **Executable Documentation Scenario**: A maintained user command whose inputs, status, streams, and expected result are verified automatically.
- **Compatibility Contract**: The declared v1 stability boundary for the CLI and Cue JSON Schema, including support states and exclusions.

## Success Criteria

### Measurable Outcomes

- **SC-001**: One hundred percent of documented SubRip grammar rows and WebVTT structure rows have traceable governed evidence or an explicit verified inapplicability reason.
- **SC-002**: All seven required fuzz surfaces run bounded seed sessions without panic or invariant violation, and every reproducible in-scope failure is retained as deterministic regression evidence.
- **SC-003**: Automated hostile-input and privacy scans report zero traversal, execution, catastrophic growth, uncontrolled output, source-byte disclosure, original paths, or local machine identifiers across applicable public workflows.
- **SC-004**: Native Windows, macOS, and Linux smoke jobs pass, all six supported target binaries build statically, and the full quality, security, race, integrity, and non-publishing release-proof gates are green.
- **SC-005**: One hundred percent of consumer-facing schema properties meet annotation policy, every enumeration and top-level semantic area has valid examples, every checked example passes its applicable constraints, and the published v0.0.0 schema digest is unchanged.
- **SC-006**: Every maintained README workflow example runs successfully and every implemented command and option has exactly one accurate documented help/reference path with no stale unavailable or release-ready claim.
- **SC-007**: Automated traceability and consistency checks find zero unexplained differences among schema, CLI help, format guides, compatibility matrices, executable behavior, and maintained examples.
- **SC-008**: A complete foreground verification run and the official pull request's hosted checks both pass with all external review findings resolved and no more than two Codex review rounds.

## Assumptions

- Issues #30 through #34 define the implemented v1 behavior that S018 hardens and documents; S018 may fix defects exposed by verification but does not add another subtitle format or public API.
- The current canonical development version remains `0.1.0` until the separately scoped v1.0.0 release-candidate slice.
- The repository's governed fixtures are sufficient for mandatory CI, while non-redistributable corpora remain optional maintainer inputs outside version control.
- Schema annotations improve discovery but remain non-normative; validation keywords, model semantics, and the constitution continue to define acceptance.
- Any behavior change required to satisfy a pre-existing v1 contract is recorded as a correction, with tests and documentation updated together.
- Publication authority granted for S018 covers its branch and pull request only. It does not authorize merging, tagging, releasing, schema hosting, or production changes.
