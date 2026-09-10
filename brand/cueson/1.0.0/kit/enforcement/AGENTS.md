# Agent Contract: Cueson

**Read this before writing any UI. It takes a minute and it is binding.**

You are working inside a brand with a fixed vocabulary. If you need a value
that is not in this document, **stop and ask**. Do not invent one, and do not
reach for a stock Tailwind palette class because it is faster.

## The stop condition

Inventing a colour, a spacing value, a radius, a font, or a component prop is
the failure this contract exists to prevent. When the vocabulary below does not
cover what you need, say so and wait.

## Colour: use the slot, never the value

Write `bg-primary`, `text-muted-foreground`, `border-border`. Never write a
hex, an `rgb()`, or `bg-slate-900`.

| Slot | Dark | Light |
| --- | --- | --- |
| `background` | `#080A14` | `#F8F8F6` |
| `foreground` | `#FFFFFF` | `#0A0A0A` |
| `card` | `#101425` | `#FFFFFF` |
| `primary` | `#62BEB2` | `#005D55` |
| `muted-foreground` | `#9A9A9A` | `#6B6B6B` |
| `destructive` | `#E9505F` | `#C0293A` |
| `border` / `input` | `#262626` | `#E5E5E5` |

### Three colour mistakes that get made constantly

1. **White text on the accent.** `#FFFFFF` on `#62BEB2` measures 2.2:1 and
   fails. The legal foreground is `#000000` at 9.53:1. Use
   `text-primary-foreground` and it is handled.
2. **The bright accent as text on a light surface.** `#62BEB2` measures 2.07:1
   on `#F8F8F6`. The light block already substitutes `#005D55`. Never override it.
3. **`#C24000` as text.** It measures 4.03:1 on the dark base. It is a fill.
   Its legal foreground is `#FFFFFF` at 5.21:1.

## Spacing and radius

Spacing scale, in px: 4/8/12/16/24/32/48/64/96/120. Nothing between them.

Radii: `rounded-sm` 6 (chips), `rounded-md` 8 (buttons, inputs, popovers),
`rounded-xl` 12 (cards, dialogs), `rounded-2xl` 16, `rounded-full` (badges).
Never `rounded-none`, never an arbitrary `rounded-[...]`.

Layout: content 1200px, narrow 720px. Gutters 24px then 48px then 80px. Section rhythm 120px then 160px then 200px.

## Type

Space Grotesk for display at 500/700. Geist for body at 400/500. Geist Mono for labels,
code, and metadata at 400.

Asking for an undeclared weight makes the renderer synthesise or substitute a face, which prints badly and forces outlined glyphs into PDFs. In mono, carry emphasis with colour.

## Affiliation

This is a ShruggieTech-owned child brand. The only approved ownership endorsement is `A ShruggieTech project`. Keep it outside the logo clear space.



## Density

Two settings ship, and both are correct in the right place. Default for
marketing and reading surfaces; compact for dense tabular data. Do not invent
a third.

## Icons

lucide, inline SVG, `currentColor`, 1.5 to 2px stroke on a 24 grid. Do not
install another icon library. If lucide lacks a domain symbol, it goes in
`icons/` drawn to the same spec.

## Accessibility, non-negotiable

- Visible 2px focus ring at 2px offset on every interactive element
- Status never carried by colour alone; pair it with a label or a shape
- Respect `prefers-reduced-motion`
- WCAG AA at rendered size

## Copy

Cueson copy is precise, calm, format-literate and evidence-conscious. Put the caption content and the source format and retained evidence first. Use familiar nouns and verbs. Keep
sentences short.

Headlines name something a reader can act on. Prefer literal product
language such as "Universal captions and subtitles" to slogans, mood, or abstract benefit claims.

Do not reach for: claims that normalization replaces the source; describing OCR output as authoritative; AI transformation hype; generic accessibility symbolism.

Cue Teal carries product identity, focus, and selection. Archive Night and its raised surfaces keep structured content primary. ShruggieTech orange remains scarce for emphasis and warning, while success and failure retain explicit labels and shapes.


Never build a sentence out of `X, not Y`, or `X over Y`, or
`rather than merely Z`. It is the clearest tell of machine-written copy.
Avoid em-dashes; use parentheses, commas, or hyphens. No testimonials, no
feature grids standing in for an explanation, no manufactured urgency.

## Before you call it done

```bash
npx eslint --config enforcement/eslint.brand.mjs .
npx stylelint --config enforcement/stylelint.config.json "**/*.css"
python3 build/verify.py
```

A build that fails any of these is not finished, whatever it looks like.
