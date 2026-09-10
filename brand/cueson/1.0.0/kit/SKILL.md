---
name: cueson-brand
description: Apply the Cueson brand system to product interfaces, documentation, repositories, and generated assets.
license: Proprietary. Internal ShruggieTech use.
compatibility: Requires a filesystem for bundled fonts and assets. The optional verification scripts require Python 3.
metadata:
  version: 1.0.0
  canon: 1.2.1
  parent: ShruggieTech
allowed-tools:
  - Read
  - Glob
  - Grep
---

# Cueson brand

Read `README.md`, `brand.json`, and `enforcement/AGENTS.md` before producing Cueson material.

Cueson provides universal captions and subtitles through a lossless, structured interchange layer. Lead with the caption content, source format, retained evidence, or normalized Cue JSON result. State clearly that OCR and transcription are derived interpretations, and never imply that normalized text replaces the authoritative source.

Use Cueframe without recomposition. Two open brackets are the retained source envelope, and three unequal rails are directly readable structured cue content. At and below 32 pixels use the reduced master with one rail. Never close the brackets, equalize the rails, introduce a CC badge, add media-control imagery, or resize the glyph and wordmark independently.

Cue Teal is the identity accent: `#62BEB2` on Archive Night and `#005D55` on light surfaces. ShruggieTech orange remains scarce and carries emphasis and warning. Use the bundled Space Grotesk, Geist, and Geist Mono families and the generated semantic tokens and components.

Write the description as `A lossless, structured interchange layer` followed by `for subtitle and caption content.` in graphical compositions. Use one line only when the available width makes it clearer.

Gate 2 retry 2 is approved for the eight public surfaces declared in `brand.json`. Consumer import remains a separate governed workflow. Run `build/verify.py` after any kit change; a nonzero problem count blocks delivery.
