# Feature Specification: Complete CLI Workflows

**Feature Branch**: `codex/S017-complete-cli-workflows`

**Created**: 2026-09-11

**Status**: Implemented, pending pull-request verification

**Input**: Deliver issue #34 through Spec Kit as work slice S017, completing validation, inspection, deterministic shell completion, and the shared help, stream, diagnostic, option, and exit-code contract for the v1 command surface.

## User Scenarios & Testing

### User Story 1 - Validate without producing output (Priority: P1)

A user supplies Cue JSON, SubRip, or WebVTT input and receives a truthful success or failure result plus actionable diagnostics without creating a converted or encoded payload.

**Why this priority**: Validation is the last missing correctness workflow and gives users and automation a safe way to prove documents before restoration, rendering, conversion, or release use.

**Independent Test**: Validate representative valid, warning-bearing, malformed, corrupted, ambiguous, and unsupported inputs in every supported class and compare exact statuses and streams while proving no output artifact or stdout payload is produced.

**Acceptance Scenarios**:

1. **Given** valid Cue JSON, **When** validation runs, **Then** JSON syntax, schema structure, semantic invariants, source-envelope integrity, and software/schema compatibility all pass and the command exits successfully without a stdout payload.
2. **Given** valid native SubRip or WebVTT, **When** validation runs, **Then** bounded decoding and the selected native grammar are checked, ordered warnings use stderr, and no Cue JSON or subtitle output is written.
3. **Given** malformed, corrupted, ambiguous, or unsupported input, **When** validation runs, **Then** it fails with the established typed status and diagnostic category and leaves every unrelated path unchanged.

---

### User Story 2 - Inspect safely in human or JSON form (Priority: P2)

A user can understand an input's format, declared and available capabilities, safe source-asset facts, cue and block summaries, diagnostics, integrity state, and conversion-loss assessment state without exposing preserved bytes or local machine information.

**Why this priority**: Users need a concise observability surface to understand what Cueson recognized and can safely do before selecting a mutating workflow.

**Independent Test**: Inspect equivalent native and Cue JSON examples in human and JSON modes, compare deterministic goldens, validate the JSON shape, and scan every output for source bytes, caller paths, hostnames, usernames, and other local identifiers.

**Acceptance Scenarios**:

1. **Given** a valid supported input, **When** default inspection runs, **Then** stdout contains one concise deterministic human-readable report and stderr contains only permitted diagnostics.
2. **Given** the same input, **When** `--json` inspection runs, **Then** stdout contains only one documented lowercase `snake_case` JSON object with the same report facts.
3. **Given** a valid Cue JSON document with multiple assets or WebVTT body blocks, **When** inspection runs, **Then** ordered safe asset entries and complete aggregate cue and block counts are reported with verified integrity state.
4. **Given** an input for which no conversion target was requested, **When** inspection reports conversion-loss state, **Then** it states that loss has not been assessed and does not imply losslessness.

---

### User Story 3 - Configure shells and discover the complete CLI (Priority: P3)

A user can generate a static completion definition for an explicitly supported shell and can rely on root and command help to describe the complete implemented command surface, options, streams, and statuses.

**Why this priority**: Deterministic completion and accurate help make the v1 command surface usable interactively and safe to automate without registering placeholders or hidden behavior.

**Independent Test**: Compare completion output for every supported shell with committed snapshots, reject every other shell without output, and mechanically prove that help and completion name every implemented command and accepted option with the documented stream and exit-code rules.

**Acceptance Scenarios**:

1. **Given** a supported shell selector, **When** completion generation runs, **Then** stdout contains exactly one deterministic non-interactive static script and repeated runs are byte-identical.
2. **Given** an unknown or omitted shell selector, **When** completion generation runs, **Then** it exits as invalid invocation, writes no completion payload, and lists only supported selectors in its diagnostic usage.
3. **Given** any implemented command, **When** help is requested before or after the command name, **Then** help names every accepted command-specific and global option and accurately states payload, diagnostic, and exit-code behavior.

### Edge Cases

- Input is JSON-looking or JSON-named but fails parsing, schema validation, semantic validation, integrity validation, or version compatibility.
- Native content evidence disagrees with a file extension or an explicit input selector.
- An encoding selector is unknown, incompatible with WebVTT, or supplied for Cue JSON.
- Input is empty, whitespace-only, larger than the established source ceiling, non-regular, missing, unreadable, or canceled during acquisition or integrity checking.
- Warning-bearing native input is valid, and quiet or silent filtering must not alter its success status.
- Cue JSON contains multiple assets, zero cues where the schema permits it, many diagnostics, source-order gaps rejected by semantics, or hashes and lengths inconsistent with embedded bytes.
- Inspection receives payload text, markup, identifiers, native raw fields, or source bytes that resemble a path, terminal escape, or machine identifier.
- Standard output or standard error fails during human, JSON, completion, help, success, warning, or error output.
- A completion selector differs only by case or whitespace, appears after `--`, or is followed by an extra operand.
- Repeated completion and inspection runs occur across Windows, macOS, and Linux with different working directories and environment values.

## Requirements

### Functional Requirements

- **FR-001**: The product MUST register `validate`, `inspect`, and `completion` only with their complete S017 behavior and MUST remove the prior documentation statements that they are unavailable.
- **FR-002**: `validate` and `inspect` MUST accept Cue JSON, SubRip, and WebVTT through one bounded, regular-file, no-follow input acquisition path.
- **FR-003**: Automatic input classification MUST recognize valid Cue JSON before native formats, MUST treat JSON-looking or JSON-named invalid content as Cue JSON failure, and MUST use content-first native detection with deterministic extension-disagreement diagnostics.
- **FR-004**: Both commands MUST support explicit `auto`, `cueson`, `srt`, and `vtt` input selection plus documented aliases; explicit selection MUST NOT silently fall through to another input class.
- **FR-005**: An encoding selector MUST be available only for native text, MUST use the existing canonical names and aliases, and MUST reject incompatible WebVTT or Cue JSON combinations as invalid invocation.
- **FR-006**: Cue JSON validation MUST check UTF-8 JSON parsing, the embedded schema, semantic invariants, source-envelope integrity, and executable/schema lockstep before success.
- **FR-007**: Native validation MUST run the installed decoder and native grammar while retaining and reporting ordered recoverable diagnostics without constructing, writing, or printing Cue JSON.
- **FR-008**: `validate` MUST write no stdout payload and create no file in every success and failure case.
- **FR-009**: Successful validation MUST emit one concise stderr success diagnostic unless suppressed by quiet or silent mode; warnings and errors MUST follow the established filters.
- **FR-010**: `inspect` MUST validate the selected input to the same depth as `validate` before reporting it.
- **FR-011**: The inspection report MUST identify the input class, canonical format, declared support state, installed ingest and render availability, exact-restoration availability, validation availability, and inspection availability without claiming unimplemented capabilities.
- **FR-012**: The inspection report MUST include ordered safe source-asset summaries containing source indexes, primary status, roles, safe basenames, media types, byte lengths, and text-encoding observations when present, but MUST exclude embedded source bytes, content hashes, asset identifiers, and timestamp values.
- **FR-013**: The inspection report MUST include cue count, non-cue block count, body-item count, per-kind WebVTT block counts, media bounds and span when present, word-level-timing state, and diagnostic severity totals, plus ordered structural cue and block summaries limited to ordinals, source order, timing, line and observation counts, type, setting counts, and non-content state flags.
- **FR-014**: The inspection report MUST include ordered diagnostic summaries containing stable severity, code, source order, and resolved cue ordinal fields when present, but MUST exclude free-form diagnostic messages and cue identifiers.
- **FR-015**: The inspection report MUST state source-integrity status and checked asset count, and MUST represent conversion loss as not assessed when no target-format analysis occurred.
- **FR-016**: Default inspection MUST write one concise deterministic human-readable report to stdout and permitted diagnostics to stderr.
- **FR-017**: `inspect --json` MUST write only one deterministic UTF-8 LF JSON object to stdout using documented lowercase `snake_case` keys and a fixed report version.
- **FR-018**: Human and JSON inspection MUST express the same facts, preserve source-defined ordering where meaningful, and remain byte-identical across repeated runs on every supported platform.
- **FR-019**: Inspection output MUST NOT contain `data_base64`, raw source bytes, payload text, native raw fields, caller-supplied paths, runtime working directories, usernames, hostnames, drive or mount details, or other local machine identifiers.
- **FR-020**: `completion` MUST accept exactly one supported shell selector from `bash`, `zsh`, `fish`, and `powershell`, with selector matching documented as case-sensitive.
- **FR-021**: Completion output MUST be a deterministic UTF-8 LF static script that performs no network, filesystem discovery, prompt, subprocess invocation, or environment mutation beyond registering completion when the user deliberately evaluates it.
- **FR-022**: Completion definitions MUST include every implemented command and accepted option, MUST exclude unimplemented capabilities, and MUST keep global and command-specific candidates contextually distinct.
- **FR-023**: Unsupported, missing, or extra completion operands MUST exit 2, emit usage on stderr, and write no stdout payload; an stdout write failure after valid invocation MUST exit 1.
- **FR-024**: Root help MUST list every implemented command, every global option, the `--` delimiter, payload and diagnostic stream rules, all three exit codes, and representative examples without listing placeholders.
- **FR-025**: Each command help page MUST document every accepted command option, applicable global options, output and diagnostic streams, success and failure statuses, and at least one valid example.
- **FR-026**: Explicit help MUST write only help text to stdout and exit 0; error-associated usage MUST use stderr and retain the originating invocation status.
- **FR-027**: Global quiet, silent, no-color, delimiter, diagnostic, stdout-payload, overwrite, and typed-error behavior MUST remain consistent across every applicable command.
- **FR-028**: Missing inputs and deterministic pre-execution path or option failures MUST exit 2; parsing, schema, semantic, integrity, capability, cancellation, and runtime stream or I/O failures after valid command acceptance MUST exit 1.
- **FR-029**: Command-level golden and process tests MUST cover valid, invalid, malformed, warning-bearing, ambiguous, unsupported, and canceled inputs; exact stdout/stderr separation; JSON shape; completion snapshots; help coverage; and statuses.
- **FR-030**: Cross-platform tests, race tests, bounded fuzz seeds, pure-Go builds, repository text checks, security checks, and the existing release-proof suite MUST remain green.
- **FR-031**: Maintained CLI, architecture, format, schema, README, and changelog documentation MUST describe only implemented S017 behavior, distinguish validation from inspection and conversion, and preserve development version `0.1.0` plus experimental native format states.
- **FR-032**: S017 MUST NOT add new subtitle formats, change Cue JSON schema, declare v1 stability, publish a schema, tag or release a version, modify production services, implement downstream media inspection, or merge the pull request.
- **FR-033**: The official pull request MUST close issue #34, pass every current-head check, address every external review finding, request no more than one second Codex round, and stop for the operator's final merge ritual.

### Key Entities

- **Validated Input**: A classified Cue JSON or native subtitle input that has passed the complete validation pipeline and retains ordered safe diagnostics for reporting.
- **Inspection Report**: A versioned format-neutral summary of capabilities, source assets, document structure, diagnostics, integrity state, and unassessed conversion-loss state.
- **Asset Summary**: The safe non-payload facts for one source asset in source order.
- **Block Summary**: Aggregate counts for cues and format-native non-cue body items without their raw content.
- **Completion Definition**: One static shell-specific script derived from the implemented command and option vocabulary.

## Success Criteria

### Measurable Outcomes

- **SC-001**: One hundred percent of accepted Cue JSON, SubRip, and WebVTT validation fixtures receive the expected deterministic status and ordered diagnostics, with zero stdout bytes and zero created output files.
- **SC-002**: One hundred percent of malformed, corrupted, ambiguous, unsupported, or version-incompatible fixtures fail at the documented classification or validation stage with no fallback to an unrelated grammar.
- **SC-003**: Human and JSON inspection goldens cover every report field and supported input class, and automated scans find zero source payload bytes, content hashes, user-controlled identifiers or messages, caller paths, usernames, hostnames, or other local identifiers in output.
- **SC-004**: Every machine-readable inspection result validates against the documented report shape, uses only lowercase `snake_case` keys owned by Cueson, and is byte-identical across repeated runs.
- **SC-005**: Completion generation for all four supported shells is byte-identical across repeated runs and snapshots name 100 percent of implemented commands and accepted options with zero unimplemented entries.
- **SC-006**: Root and command help coverage proves every accepted command and option appears in the correct help surface and that every help page states streams and exit statuses accurately.
- **SC-007**: Full tests, race detection, bounded fuzz seeds, pure-Go cross-builds, repository checks, and development release verification complete with zero failures.
- **SC-008**: The final S017 pull-request head has zero failed or pending checks, zero unresolved review threads, every external finding addressed, and no more than two Codex review rounds.

## Assumptions

- The existing Cue JSON, model, source-integrity, codec registry, parser, diagnostic, and output boundaries remain authoritative and are reused rather than duplicated.
- `validate` and `inspect` use `--format` for input classification, matching `encode`; `cueson`, `json`, and `cue-json` are additional explicit selectors for Cue JSON.
- Native inspection derives the same in-memory document used by encoding but never marshals or publishes Cue JSON.
- Bash, Zsh, Fish, and PowerShell are the explicitly supported v1 completion targets; completion remains static and does not inspect the filesystem or invoke Cueson recursively.
- Inspection intentionally omits source timestamp values because they add privacy and reproducibility risk without helping the stated v1 workflow.
- Conversion loss is target-dependent, so inspection without a target records `not_assessed` rather than `lossless` or a numeric loss count.
- This slice completes issue #34 only. Hardening, end-to-end documentation, release-candidate preparation, schema publication, and v1 release remain owned by later issues.
