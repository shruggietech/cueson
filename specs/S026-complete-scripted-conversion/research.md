# S026 consolidated research

**Date**: 2026-09-15

Phase 0 dispatched independent outbound, target, and evidence research through the installed speckit-plan workflow. Findings are recorded in research-outbound.md, research-targets.md, and research-evidence.md. All unknowns are resolved; no operator question is required.

## Decision: native owner traversal and balanced emphasis

Preserve readable dialogue and target-supported bold/italic/underline using native style and ordered override/reset state. Traverse all native owners for atomic omissions instead of relying only on common plain text. Reject malformed preserved owners and target literal reinterpretation. Alternative: flatten only cue plain_text, rejected because it silently loses source presentation and records.

## Decision: drawing and empty output

Drawing-only or entirely unreadable dialogue and zero-dialogue scripts are fatal for text targets. Mixed drawing/readable dialogue may omit each drawing span with its own loss. Valid empty scripts may remain valid across scripted variants. Alternative: drop cues or invent text, rejected because neither is required and both weaken fidelity.

## Decision: target construction and timing

Text-source scripted targets use deterministic ScriptType/WrapStyle and one Default style (Arial 20, white primary, red secondary, black outline/shadow, neutral emphasis/scales/angle/spacing, bottom-center, margins 10, encoding 1), without inferred video dimensions. Nearest-centisecond rounding uses quotient/remainder with ties upward, one loss per changed endpoint, checked overflow, and fatal collapse. Literal braces and literal recognized control pairs are fatal; actual LF/NBSP become explicit native controls. Alternative: unchecked truncation or invented interval enlargement, rejected because timing loss would become silent or policy-dependent.

## Decision: variant semantics

Deep-copy native owners and retain representable script metadata, styles, events, dialogue, comments, safe unknown content, and attachments at their original positions. Map alignment explicitly. Replace Layer/Marked with deterministic target defaults and account for their different roles. Account per omitted or degraded style/event field and unsupported dialect control, including neutral observed values. Color/alpha roles follow the researched conservative policy; no pixel-equivalence claim is made. Alternative: rename columns or flatten through text, rejected because dialect fields are not semantic identities.

## Decision: private constructed-target boundary

Normal Cue JSON Validate and native Render remain unchanged. Add internal ValidateScriptedTarget(source,target) and RenderTarget(ctx,source,target) that validate the original document and full integrity, require exact envelope equality, validate target native ownership/privacy/dialect/projections, and skip only the target's original-source dialect claim. Existing captured content remains independently tied to original source positions. Capture-free targets do not parse unrelated original UTF-16 or legacy text bytes as scripted captures. The private target is never emitted as Cue JSON. Alternative: fabricate or rewrite source bytes or expose a JSON bypass, rejected by the constitution.

## Decision: one report and complete evidence

Append fourteen closed loss codes after established ranks, preserving all old meanings and pair goldens. Existing common losses apply only when their precise source features are omitted. Custom common source identifiers on scripted documents require a precise generic omission code; native SubRip sequence framing remains exempt and existing WebVTT identifier accounting remains unchanged. Nonzero ignored SSA per-color high bytes require field degradation when effective AlphaLevel mapping replaces them; zero ignored high bytes are ordinary representation normalization. Add scripted section/record orders to source-reference validation. Use existing manifest roles and four shared source baselines with twelve direction-specific expectations, plus lossy/fatal/precision/privacy cases. Strict analysis completes before the renderer and both strict/fatal paths retain atomic publication behavior. Alternative: truncate losses or extend unrelated fixture schemas, rejected as unnecessary and unsafe.

## Delivery authority

Implementation refinement: edited scripted Cue JSON can legitimately carry non-centisecond common timing while retaining exact source-native timestamp captures. The same checked quantization and endpoint loss policy therefore applies to variant targets as well as text-source targets. This preserves S025 edit ownership and avoids an unreported renderer-only precision refusal.

The user expressly authorizes push and official PR publication for S026, overriding the skill's default pre-push halt. Final merge, tags, releases, and production changes remain protected. First Codex review is automatic; exactly one second request is permitted if findings require it. No extension registry exists, so installed command hooks are skipped.
