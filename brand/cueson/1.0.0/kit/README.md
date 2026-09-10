# Cueson brand system

Cueson is a lossless, structured interchange layer for subtitle and caption content. It performs format-specific work at the file boundary, exposes a predictable Cue JSON surface, and retains the original carrier and native details as authoritative evidence.

## Verbal identity

- Name: `Cueson`
- Domain: `cueson.io`
- Slogan: `Universal captions and subtitles`
- Description line 1: `A lossless, structured interchange layer`
- Description line 2: `for subtitle and caption content.`

Use the two description lines together and in that order. They may be set on one line only when the available width makes that composition clearer.

## Identity idea

Cueframe places three unequal cue rails inside two open brackets. The brackets represent the retained source envelope. The rails represent directly readable timing and caption content. Their coexistence communicates normalization without replacement.

The full master uses two 34-unit brackets and three 32-unit rails on a 512-unit grid. The reduced master thickens both brackets and keeps one central rail. Use the reduced master at and below 32 pixels, and never use the full master below 24 pixels.

## Color

| Role | Dark | Light |
| --- | --- | --- |
| Identity accent | `#62BEB2` | `#005D55` |
| Surface | `#080A14` | `#F7F8FC` |
| Card | `#101425` | `#FFFFFF` |
| Secondary | `#181D33` | `#ECEFFC` |
| Hover | `#222944` | `#E0E5FA` |
| Text | `#F4F6FF` | `#14172A` |
| Muted text | `#A7AEC7` | `#5B6078` |
| Line | `#596181` | `#8A8FA8` |

Cue Teal is the product identity color. ShruggieTech orange remains a scarce inherited emphasis and warning role. Success and failure always retain a label or shape and never rely on color alone.

## Typography

Space Grotesk at 500 and 700 handles display text and the outlined wordmark. Geist at 400 and 500 handles body and interface copy. Geist Mono at 400 handles timestamps, cue identifiers, hashes, formats, and provenance.

## Voice

Lead with the caption content, the source format, or the normalized Cue JSON result. Name OCR and transcription as derived interpretations. Avoid language that implies the source was discarded or that inferred text is authoritative.

## Approval status

Gate 1 approved `cueframe-r1` and the revised `cue-teal-r1` palette on 2026-09-09. Gate 2 approved retry 2 of the exact wordmarks, lockups, icons, applications, and eight public-surface compositions on 2026-09-09. Consumer import into the Cueson repository still requires its own issue and Spec Kit slice.

## Verification

Run `build/verify.py` after any generated-kit change. Delivery requires zero verification problems and zero glyph failures. Generated output belongs under ignored `dist/` paths and must never be committed.

## Parent endorsement

Set `A ShruggieTech project` in Geist Mono with positive tracking. Keep it visually subordinate and outside the Cueson logo clear space. Never combine the two marks into one lockup.

Names, wordmarks, logos, endorsement lockups, and logo path geometry remain reserved as described in the repository's [brand asset terms](../../LICENSE-BRAND.md).
