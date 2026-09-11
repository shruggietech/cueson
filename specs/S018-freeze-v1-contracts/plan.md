# Implementation Plan: Freeze and Prove the v1 Contract

**Branch**: `codex/S018-freeze-v1-contracts` | **Date**: 2026-09-11 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/S018-freeze-v1-contracts/spec.md`

## Summary

Close issues #35, #36, and #41 as one v1 contract-freeze slice. Add common bounded Cue JSON acquisition and explicit complexity ceilings, strengthen the governed corpus and seven required fuzz surfaces, introduce a machine-checked format conformance matrix and an opt-in external-corpus verifier, enrich every public canonical-schema surface with non-normative descriptions and valid examples, make maintained command examples executable, and reconcile all public documentation with the behavior already delivered by S014 through S017. Preserve development identity `0.1.0`, the immutable v0.0.0 schema, source fidelity, and every publication boundary.

## Technical Context

**Language/Version**: Go 1.25.0 for product and repository verifiers; JSON Schema Draft 2020-12; UTF-8 Markdown without hard-wrapped prose

**Primary Dependencies**: Go standard library, `github.com/santhosh-tekuri/jsonschema/v6` v6.0.3, `golang.org/x/text` v0.39.0, `golang.org/x/sys` v0.41.0; pinned CI-only Actionlint v1.7.12, Staticcheck v0.7.0, govulncheck v1.8.0, GoReleaser v2.18.1, and Syft v1.51.1

**Storage**: Canonical and immutable JSON Schema files, governed byte fixtures and `testdata/manifest.json`, a machine-readable conformance matrix, Markdown documentation, help/completion goldens, and ignored release-proof output under `dist/`

**Testing**: Go unit, integration, process, conformance, fuzz, race, schema-compilation, annotation-coverage, executable-documentation, nested-module, native-platform, static-build, security, link, text, brand-integrity, and non-publishing release-proof checks

**Target Platform**: Windows, macOS, and Linux; amd64 and arm64; official binaries remain `CGO_ENABLED=0`

**Project Type**: Single command-line product with internal domain packages and standalone repository-verifier modules

**Performance Goals**: Reject inputs above the existing 64 MiB ceiling before unbounded allocation; cap adversarial document collections, per-item occurrences, diagnostics, and loss observations at explicit deterministic limits; run every required fuzz surface at a fixed bounded work count in hosted CI; keep ordinary valid subtitle workflows within existing test latency

**Constraints**: Preserve original bytes and unknown content; never silently truncate; never expose caller paths or machine identity; do not mutate the immutable v0.0.0 schema; keep schema annotations non-normative; no public Go API; no network requirement for product or mandatory verification; no tag, release, public schema, production-domain, or merge action

**Scale/Scope**: Three independently closeable v1 issues; 12 root schema properties, 44 reachable public definitions, 135 direct public properties, 11 enum sites, two native subtitle grammars, seven required fuzz surfaces, six build targets, three native operating systems, nine shipped commands, and the maintained root documentation set

## Constitution Check

*GATE: Passed before Phase 0 research and passed again after Phase 1 design.*

- **Lossless Source Preservation**: PASS. Hardening rejects unsafe work without altering accepted source bytes; restoration evidence remains byte exact.
- **Schema and Official Software Discipline**: PASS. Development schema and executable remain `0.1.0`; emitted bytes stay identical; v0.0.0 remains immutable.
- **Common Model Plus Native Fidelity**: PASS. Matrix and documentation cover both views and do not elevate annotations or derived content to source truth.
- **No Silent Loss**: PASS. Complexity ceilings reject with one stable failure rather than truncating; conversion matrices account for every known loss.
- **Test-First Format Work**: PASS. Tasks put hostile, corpus, schema, documentation, and fuzz evidence before or alongside changes.
- **Portable and Secure Operation**: PASS. Bounded regular-file acquisition, hostile tests, three native runners, six static targets, and hidden Windows process execution remain required.
- **Documentation and Delivery Authority**: PASS. S018 uses Spec Kit, preserves issues #35, #36, and #41, and stops before merge, tagging, release, or production work.

## Project Structure

### Documentation (this feature)

```text
specs/S018-freeze-v1-contracts/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── conformance-matrix.md
│   ├── executable-documentation.md
│   └── schema-annotations.md
├── checklists/
│   └── requirements.md
└── tasks.md
```

### Source Code (repository root)

```text
cmd/cueson/                         # Minimal process adapter
internal/cli/                       # Shared bounded input use and executable documentation scenarios
internal/codec/                     # Detection, parser, renderer, settings, markup, and fuzz limits
internal/convert/                   # Loss accounting, conversion cycles, documentation equality, and bounds
internal/model/                     # Shared semantic collection ceilings and validation
internal/schema/                    # Canonical annotated schema and annotation validation
internal/source/                    # Bounded regular-file no-follow acquisition and envelope fuzzing
internal/conformance/               # Matrix, hostile cases, corpus cycles, and path-leak proof
internal/testutil/                  # Governed fixture and portable-identity helpers
scripts/docs-verify/                # Maintained-document inventory and static reference checks
scripts/corpus-verify/              # Opt-in read-only external corpus verifier
testdata/
├── conformance-matrix.json
├── manifest.json
├── fixtures/
├── malformed/
└── fuzz/
docs/                               # CLI, schema, format, conversion, compatibility, architecture, release guidance
.github/workflows/ci.yml            # Bounded fuzz and all-module quality/security execution
README.md
CONTRIBUTING.md
SECURITY.md
CHANGELOG.md
```

**Structure Decision**: Retain the existing dependency-light internal package boundaries. Put product safety limits with the packages that own the amplified collection, put cross-package evidence in `internal/conformance`, keep optional external corpus work in a standalone verifier that cannot alter source inputs, and make repository documentation mechanically traceable without exposing an additional product command.

## Implementation Phases

### Phase 1 - Safety and evidence foundation

Define shared limits, route every Cue JSON workflow through bounded no-follow acquisition, add matrix and corpus data contracts, create focused hostile regressions, strengthen the seven fuzz surfaces, and extend existing CI contexts to execute bounded fuzz work and every nested verifier module.

### Phase 2 - Machine-readable schema contract

Apply the annotation policy to every reachable public property and definition, add valid complete and fragment examples, compile and validate those examples, enforce recursive coverage, and pin the immutable v0.0.0 digest.

### Phase 3 - Executable user documentation

Replace placeholders with isolated governed-fixture workflows, add conversion and compatibility references, reconcile stale repository documents, tie CLI and format tables to code and matrix identifiers, and make every advertised workflow executable under tests.

### Phase 4 - Integrated release-readiness proof

Run focused and full tests, fixed-work fuzz sessions, race, vet, Staticcheck, vulnerability checks for every module, three-platform hosted CI, six static builds, schema and docs consistency, encoding and mojibake checks, brand integrity, and a clean non-publishing development snapshot. Converge against every S018 requirement before publication.

## Complexity Tracking

No constitutional violation or architecture exception is required. The additional standalone verifier is justified by the explicit maintainer-corpus requirement and remains outside the shipped product.
