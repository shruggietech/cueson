# S026 conversion interface contract

## Supported graph

The canonical source/target values are subrip, webvtt, ass, and ssa. Every ordered distinct pair is supported for the documented baseline, giving twelve directions. CLI aliases srt/subrip and vtt/webvtt remain. Same-format conversion directs the user to render. Unsupported formats retain explicit refusal.

## Text and timing

Outbound scripted conversion preserves readable dialogue and representable b/i/u emphasis. Drawing-only, unreadable dialogue, and zero-dialogue text output are fatal; mixed drawings are atomically omitted. Literal text that changes target grammar semantics is fatal.

Scripted targets use nearest-centisecond nonnegative rounding with ties upward and checked quotient/remainder arithmetic. Each changed endpoint reports conversion_scripted_centisecond_quantized. Collapse or overflow is fatal. Actual LF and NBSP become native controls; literal braces or literal backslash-N/n/h constructs are fatal.

## Defaults and variants

Text-source scripted targets use matching ScriptType, WrapStyle 0, and Default Arial 20 style with documented neutral presentation, bottom-center alignment, margins 10, and encoding 1. No video dimensions or source provenance are inferred.

Variant conversion retains representable native owners and explicitly maps alignment. Layer/Marked roles and non-equivalent color/alpha and dialect-only fields are accounted atomically. Unknown controls are preserved only when their target meaning is proven; otherwise report or reject. No pixel-equivalence claim is made.

## Reports and refusal

Loss reports remain runtime-only, use the existing bounded shape, retain all old codes/meanings, and append the thirteen scripted codes described in data-model.md. Native references use positional pointers, validated physical order, fixed messages, and controlled context without raw field values.

Strict refuses any complete known-loss report before invoking a renderer. Fatal failures in either mode return no payload; stdout has no partial payload and existing destinations remain unchanged. Warnings/errors follow established diagnostic filters on stderr.

## Constructed target boundary

Internal ValidateScriptedTarget(source,target) validates the original normally, requires exact SourceEnvelope equality, and validates constructed target owners without making a target-dialect claim about original bytes. Retained captures remain source-checked; capture-free targets do not decode unrelated source bytes as scripted captures. RenderTarget(ctx,source,target) verifies complete source integrity and uses bounded serialization and reparsing. Public Cue JSON Validate and native Render remain unchanged.

## Delivery

Official PR must close #60 and #61, pass current-head checks, and satisfy every review finding. At most one manual second @codex review is allowed. Final merge is human-owned; no tag, release, or production operation is part of S026.
