# S024 verification and delivery record

## Authority and analysis

The operator designated S024 and explicitly authorized automatic push, an official pull request, all review remediation and at most one second Codex review. Final merge, releases/tags and production remain separately authorized. Base is merged S023 main `0b1922f8a50cdfda766db63012f0d9f7d259146c`; implementation branch is `codex/S024-extend-versioned-schema-model`.

Installed Spec Kit specification, clarification, checklist, plan, tasks, analysis and implementation prerequisites ran before implementation. Five design clarifications are recorded in spec.md. Blocking read-only analysis passed with 20 requirements mapped to 20 tasks, 100% coverage and no critical/high/medium findings. All custom quality criteria were assessed; reviewer-owned checklist markers remain unchanged under the explicit autopilot authorization.

## Implementation and issue acceptance

Issue #55 selects only exact current `1.1.0-dev` and frozen historical `1.0.0` pairs from unmodified parsed input, validates local artifact identities independently and denies unbundled resources. Historical semantics remain separate from scripted current rules. Both historical formats have command regressions for validate, inspect, restore, matching render and conversion, third-party producer/source/native retention, strict loss refusal and invalid identity/integrity refusal before publication. Current discovery remains current. Header identity failures retain the established structural error classification.

Issue #56 adds matching typed ASS/SSA branches, ordered physical section/record ownership, declarations/styles/events/attachments, raw field occurrences and typed lexical checks, Unicode scalar spans/tags/karaoke, derived readable Text/speaker/token provenance, explicit complexity ceilings and original/edited metadata privacy. Recursive annotations and portable complete examples cover the new definitions. Recognition is `schema_only`, with native ingest/render absent and generic validate/inspect/restore available. Native conversion remains deferred. Current typed JSON round trips preserve accepted fields.

The coordinating implementation also routes direct render/restore through shared Cue JSON finishing, reports loaded schema identity from inspection, adds recognized nil-codec registrations and stages software/schema/output identity in lockstep. Current development packaging uses the current canonical artifact and explicit development proof; immutable release schemas and historical evidence remain unchanged.

## Verification completed before final integration

- Exact historical/registry tests passed, including independently versioned producer and retained WebVTT occurrences/blocks.
- Root CLI/version/codec tests passed with both historical command families and scripted generic capability/publication tests.
- All six nested module test suites passed; nested vet and Staticcheck passed.
- All six nested vulnerability scans passed with Go 1.27.1 and zero reachable vulnerabilities. Corpus verifier imports one advisory-affected package without calling affected symbols.
- Workflow Actionlint v1.7.12 passed.
- Site full verification passed: lint, deterministic generation/check (18 documents, two immutable schemas and 12 brand assets), 78 unit tests, production build, 74 browser/accessibility tests passed with four intentional project skips, 22-route artifact proof and Cloudflare deployment dry run.
- Documentation verifier passed (24 documents, 122 local links, ten registered examples and 17 format rows); repository formatter passed.

Initial integration failures were actionable development dependencies: current identity staging preceded new schema compilation, a schema regex used an unsupported escape, the guarded karaoke `K` lexeme needed a narrow native-value exception, and exact identity selection needed its established structural error wrapper. These were corrected without weakening historical fixtures or naming coverage. Staticcheck identified two capitalized error openings, which were lowercased.

## Limits and deferred work

Model bounds govern physical items, per-item occurrences, aggregate fields/spans, physical lines, nesting, derived data and embedded attachments. Physical-record/80-byte encoded-line bounds are tighter than 16/32 MiB decoded attachment bounds; both remain explicit. The 64 MiB canonical native output ceiling is a renderer publication gate for S025/#59 because S024 installs no scripted renderer. Corpus and codec outcomes #57/#58/#59 remain the proposed S025 group; conversion, full native CLI/conformance and stable release/public hosting remain later slices.

No additional active issues arrived at the S024 assessment: the same 14 open issues comprise epic #51 and children #55 through #67. Project items #55/#56 were reconciled to In progress with Slice S024 and unused default Status; native blockers remain authoritative. Final PR/CI/review evidence will be appended after publication.

## Final integration and hosted delivery

Final Go 1.25 `go test -count=1 ./...` and root vet passed. Root Staticcheck v0.7.0 passed with compatible Go 1.26.8; Go 1.27.1 vulnerability analysis found zero reachable vulnerabilities (one imported-package advisory without affected calls). The first Staticcheck run using Go 1.27.1 encountered its unsupported export-data format; using the verified compatible toolchain corrected the tool mismatch without changing product code.

Ten fixed-work fuzz targets passed with 1,000-execution budgets: format selection, SubRip parsing/timecodes, WebVTT parsing/settings-markup/render cycles, source encoded integrity, SubRip/WebVTT conversion cycles and new scripted projection. Six Go 1.25 pure-Go trimpath builds passed for Linux/Windows/macOS on amd64/arm64. Gofmt inventory was empty and repository text verification passed. Final installed analysis prerequisites passed; read-only consistency audit found complete issue acceptance coverage and no unresolved blocking design findings.

Historical v0.0.0 SHA-256 remains `d15c7fa5227156109dd6be3d39b711aca3503794bb862169dfca96ee80adb975`. Historical v1.0.0 release and locally embedded copies both remain `1aad14567033d7e14d9beb78985e18007aefb5345095370b11b6b887df7ec541`. Valid and malformed non-dialogue event offsets/observations also resolve against native Text; aggregate tests use genuine projected occurrences.

Final complete site verification was repeated after maintained-document updates and passed with the same 78 unit/74 browser results, four intentional skips and accepted artifact/dry-run proof. Spec and historical resource UTF-8/BOM/mojibake sanity checks passed, and repository whitespace checks passed.

At initial commit preparation, official PR publication and hosted review/CI completion were pending. Subsequent current-head checks, review threads/replies and delivery readbacks on the official PR are authoritative for those later transitions. Human final merge follows satisfied hosted gates under the surrounding protocol.

## First hosted review remediation

Official PR #69 closes #55 and #56. CodeQL identified a signed color conversion; commit `29626e8` added bounded signed/unsigned parsing and boundary regressions, and hosted CodeQL passed. Automatic Codex round one reviewed `029b793` and raised five findings: unknown-field privacy, arbitrary metadata paths, truthful source capture, unknown override identity and current verifier instructions. All five received implementation/documentation fixes and meaningful regressions before second-round consumption.

Native field privacy now uses actual dialect/context profiles and closed inert extension values, including wrong-role name spoofing. Path detection covers arbitrary absolute/home/URI values at punctuation and Unicode boundaries. Capture observations resolve to integrity-verified original physical lines and source Format declarations, including Comment/repeated Events, while constructed owners omit invented captures and common timing/style edits remain independent. Ambiguous alphabetic reset/font suffixes retain their complete unknown names unless declared content identifies them; source Text is unchanged. Maintained executable instructions consistently use `1.1.0-dev -development` and separate historical stable reproduction.

Review integration passed the full Go 1.25 suite, root vet, Staticcheck with Go 1.26.8, zero-reachable-vulnerability analysis with Go 1.27.1, 1,000-execution projection fuzz, repository formatting and documentation verification. Final focused model/schema verification additionally passed captured Comment/repeated Events and Unicode-boundary regressions. Complete site verification passed 78 unit tests and 74 browser/accessibility checks with four intentional skips, plus build/artifact/dry-run gates. Current-head hosted checks, replies and resolved threads remain authoritative for subsequent review completion; no third automatic review or final merge is authorized by this delivery record.

The `a87c816` hosted race job passed after a substantial slowdown. A bounded audit found linear but unnecessary work across 1,640,625 adversarial unknown fields: duplicate identity regex scans and repeated Unicode scalar automata. The follow-up replaces redundant scans with an exact marker prefilter and an equivalent ASCII-fast Unicode letter/number predicate, preserving accepted punctuation, Unicode grammar, privacy outcomes and all ceilings. Explicit grammar/prefilter regressions and unchanged aggregate fixtures passed; focused model runtime decreased from 7.736 seconds to 0.601 seconds on the same local launcher/toolchain. Second-round consumption remains deferred until this final head passes hosted CI.

## Second hosted review remediation

All product/security checks passed on `33aead2`, including hosted race detection in 39 seconds versus 5 minutes 11 seconds before the equivalent-grammar optimization. Exactly one second review was requested for that full head and read back from GitHub. Codex completed round two with four findings: malformed record privacy, relative metadata paths, code execution Effects and nonpositive original/native dialogue intervals. Each reproduced in regression coverage before correction.

Context-correct canonical declarations and complete framing now precede native content exemptions. Malformed Dialogue rejects; retained malformed Comment/Style values and structured owners undergo conservative inspection. Relative resource separators reject outside explicit content roles. Code Effects reject case/whitespace/modifier variants. Checked original and structured native intervals require positive duration independently of editable common timing, including original records whose current capture claims or cues have been omitted.

A bounded independent audit additionally closed structural field/declaration names and semicolon-comment provenance channels. The earlier development allowance for path-looking semicolon comments is deliberately replaced with conservative inert-comment inspection because comments can capture original resource identity. Safe unknown names, safe Unicode comments and bounded malformed non-dialogue records remain retained; source assets are never sanitized or rewritten. Focused final Go 1.25 model/schema tests passed (0.964/0.064 seconds), and complete site verification again passed 78 unit/74 browser checks with four intentional skips and artifact/dry-run proof. Current-head hosted CI, replies and resolved second-round threads establish the final descendant remediation gate. No third review request is permitted or emitted; final merge remains human-owned.

The final context audit also confirmed and corrected the development model's unconditional Actor-to-Name alias: aliasing now applies only to events, matching the ratified contract. Styles require actual Name; safe extra Actor fields remain unknown inert extensions and cannot borrow Name privacy exemptions. Captured ASS/SSA event Actor names/provenance remain valid with raw names unchanged. Focused model/schema tests passed again (0.960/0.068 seconds); subsequent complete integration and current-head hosted checks provide publication evidence.
