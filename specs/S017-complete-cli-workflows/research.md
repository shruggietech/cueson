# Research: Complete CLI Workflows

## Decision: One validated-input authority

Extract Cue JSON precedence and native selection from the conversion workflow into one read-only loader returning the validated document, input class, selection basis, native evidence, and ordered diagnostics. Validation, inspection, and conversion use this authority with operation-specific options.

**Rationale**: JSON fallthrough, schema lockstep, source integrity, format selection, decoding, and native model validation are correctness boundaries that must not diverge among commands.

**Alternatives considered**: Duplicate command-local loaders, which would drift; extension-only selection, which contradicts content-first detection; reopen the source after classification, which breaks single-acquisition fidelity.

## Decision: Assertion-only validation

Expose `validate --format auto|cueson|srt|vtt [--encoding NAME] INPUT` without output, force, pretty, stdout, strict, or speaker-detection options. Native validation disables heuristic speaker derivation and performs no model serialization.

**Rationale**: Validation proves input correctness and reports diagnostics. Encoding owns Cue JSON publication, and adding an optional validation payload would duplicate a mutating workflow and weaken stdout guarantees.

**Alternatives considered**: `validate --output`, which duplicates encode; `validate --json`, which was not requested and would create a second report contract; silent success, which is unfriendly for interactive assertion and makes quiet meaningless.

## Decision: Cue JSON precedence and staged validation

Auto mode first accepts valid Cue JSON. JSON-looking or JSON-named content that fails UTF-8 parsing, schema structure, semantics, source integrity, or lockstep remains a Cue JSON failure and never falls through to native grammar. Explicit selectors never fall through.

**Rationale**: A corrupt canonical document must not be reinterpreted as subtitle text, and users need the failing validation stage to remain truthful.

**Alternatives considered**: Native detection first, which can misclassify JSON payload text; generic invalid-input errors, which erase actionable boundaries.

## Decision: Private privacy-bounded inspection report

Keep report types under `internal/cli` with fixed ordered structs and report version `1`. Report input classification, format, schema compatibility, declared and installed capabilities, verified integrity, safe asset summaries, aggregate document facts, structural cue and block summaries, diagnostic codes and locations, and loss state.

**Rationale**: `inspect --json` is a public CLI output contract but not Cue JSON, a domain model, or a public Go API. Fixed structs make output deterministic and snake_case while preserving those boundaries.

**Alternatives considered**: Add inspection fields to Cue JSON, which pollutes source truth; expose internal model objects, which leaks raw content and ties the report to schema evolution; maps, which broaden shape and ordering risk.

## Decision: Omit content and user-controlled identifiers

Inspection excludes `data_base64`, hashes, timestamps, asset IDs, cue IDs, source identifiers, payload and token text, speaker and OCR text, native raw fields, diagnostic messages, caller paths, and machine information. It resolves diagnostic cue references to safe ordinals and proves source integrity through status and counts.

**Rationale**: Source bytes are explicitly forbidden, content hashes fingerprint them, and user-controlled identifiers or messages can contain paths or source excerpts. Structural facts and stable codes deliver the requested inspection value without echoing sensitive content.

**Alternatives considered**: Emit hashes as proof, which is unnecessary once verification status is stated; sanitize free-form messages, which risks changing meaning and duplicating provenance rules; truncate raw text, which still leaks content.

## Decision: Conversion loss remains unevaluated

Emit `loss.status = "not_evaluated"` and `loss.reason = "target_format_required"` because inspection has no target. Do not run speculative conversions or report a zero loss count.

**Rationale**: S016 defines loss as a runtime property of one source-target pair. A target-free zero would falsely imply losslessness.

**Alternatives considered**: Evaluate both targets, which expands scope and output; omit loss state, which fails the issue's request to expose relevant loss state.

## Decision: Four static completion targets

Support exactly case-sensitive `bash`, `zsh`, `fish`, and `powershell` selectors. Generate dependency-free static scripts that use shell-native tokenization and registration, never invoke Cueson recursively, execute subprocesses, access a network, inspect subtitle content, select a shell from the environment, or modify a profile.

**Rationale**: These targets cover first-class project platforms with stable programmable-completion mechanisms while preserving determinism and non-interactivity.

**Alternatives considered**: Infer `$SHELL`, which is nondeterministic; support cmd.exe, which lacks a comparable native contract; add a CLI framework solely for completion, which is disproportionate and would force parser churn.

## Decision: Ordered surface catalogue for help and completion

Create one ordered catalogue of shipped commands, option spellings, value kinds, canonical enum candidates, and descriptions. Render root and command help plus completion candidates from it, while the existing parser continues to own semantic conflicts and typed invocation errors.

**Rationale**: The current parser and help constants duplicate surface vocabulary. Adding four more hand-maintained copies would make drift likely; replacing the parser entirely is riskier than the feature.

**Alternatives considered**: Static independent help and scripts, which can disagree; a complete parser rewrite or Cobra dependency, which expands review and regression surface.

## Decision: Diagnostics retain stream failures

Give the shared diagnostic writer a retained first-write error and finalize successful command status as runtime failure if stderr output could not be completed. Invalid invocation remains exit 2 even if its usage stream also fails; already-failing runtime operations remain exit 1.

**Rationale**: Success cannot be claimed when the command's promised success or warning diagnostic was not delivered, and this closes a current shared stream-handling gap without changing error precedence.

**Alternatives considered**: Ignore stderr errors, which contradicts reliable process behavior; return a new exit code, which is forbidden; buffer every command's complete output, which is unnecessary.

## Decision: Help and process evidence are exact

Every command help includes local and global options, streams, exit codes, and examples. Golden tests cover root plus all commands and all completion scripts; catalogue invariants prove parser acceptance and absence of unimplemented vocabulary; one actual-binary process matrix proves statuses and stream separation.

**Rationale**: Exact public-text and process evidence catches drift that substring tests miss and closes the issue's explicit help and process requirements.

**Alternatives considered**: Documentation-only review, which is not executable evidence; interpreter-required completion tests only, which would make environments without every shell fail spuriously.
