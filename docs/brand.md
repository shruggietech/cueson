# Cueson Brand Kit

**Status:** Imported official Cueson brand kit 1.0.0

**Acquired:** 2026-09-10 through Spec Kit slice `S011-import-brand-kit`

The official Cueson identity is an externally produced ShruggieTech asset set. This repository retains the complete published handoff so project documentation can use approved assets from a stable local source without recreating, optimizing, normalizing, or hotlinking them.

## Acquisition authority

The operator-designated acquisition source is the [official ShruggieTech Cueson brand-kit download](https://brand.shruggie.tech/cueson/downloads/cueson-brand-1.0.0.zip). Two independent downloads made on 2026-09-10 produced the same 3,889,985-byte ZIP with SHA-256 `095e572a0db73472c9fdcf5001f68e702105a3d0ec1fcc34f830ae57c608cd07`.

The repository retains that exact delivery as [cueson-brand-1.0.0.zip](../brand/cueson/1.0.0/archive/cueson-brand-1.0.0.zip). Its complete safe extraction is under the versioned [retained kit](../brand/cueson/1.0.0/kit/README.md), which contains 265 regular files totaling 5,469,058 extracted bytes.

The Cueson-owned [import manifest](../brand/cueson/1.0.0/import-manifest.json) is the acquisition inventory of record. It records the source, archive identity, all 265 accepted paths, exact sizes, SHA-256 digests, and repository consumer references. The kit's own [manifest](../brand/cueson/1.0.0/kit/manifest.json), provenance records, approval records, and source notes remain preserved as delivered, but they are not substituted for the independent inventory of the received archive.

## Identity use

The exact slogan is `Universal captions and subtitles`.

The canonical description is `A lossless, structured interchange layer for subtitle and caption content.`

The README selects the `color` horizontal SVG for dark surfaces and the `light` horizontal SVG for light surfaces and fallback rendering. The standalone media-format guide uses the dark-surface `color` horizontal SVG, the retained SVG favicon, and the retained font stylesheet with its local WOFF2 files. These are direct references into the retained kit, not consumer copies.

Use Cue Teal for identity, focus, and selection. Use ShruggieTech orange sparingly for inherited emphasis and warning. Success and failure must retain a label or shape rather than relying on color alone.

| Role | Dark surface | Light surface |
| --- | --- | --- |
| Identity accent | `#62BEB2` | `#005D55` |
| Base surface | `#080A14` | `#F7F8FC` |
| Card | `#101425` | `#FFFFFF` |
| Secondary | `#181D33` | `#ECEFFC` |
| Hover | `#222944` | `#E0E5FA` |
| Text | `#F4F6FF` | `#14172A` |
| Muted text | `#A7AEC7` | `#5B6078` |
| Line | `#596181` | `#8A8FA8` |

Space Grotesk at weights 500 and 700 is the display face. Geist at 400 and 500 is the body and interface face. Geist Mono at 400 is used for timestamps, cue identifiers, hashes, formats, and provenance.

Use the supplied artwork without redrawing, recoloring outside the approved variants, path normalization, non-uniform scaling, raster-derived tracing, or alteration of its alpha silhouette. Keep at least 48 units of external clear space. The full mark has a minimum size of 24 pixels, the reduced mark is required at and below 32 pixels, the horizontal lockup has a minimum width of 160 pixels, the stacked lockup has a minimum width of 112 pixels, and the wordmark has a minimum width of 120 pixels.

Set the exact endorsement `A ShruggieTech project` in Geist Mono with positive tracking. Keep it visually subordinate and outside the Cueson logo clear space. Do not combine the Cueson and ShruggieTech marks into one lockup.

The retained kit's [brand-system guide](../brand/cueson/1.0.0/kit/README.md), [brand data](../brand/cueson/1.0.0/kit/brand.json), and [PDF guide](../brand/cueson/1.0.0/kit/brand-guide.pdf) remain the detailed authorities for logo geometry, voice, applications, and approved variants.

## Licensing and attribution

The Cueson name, wordmarks, logos, endorsement lockup, and logo geometry are reserved ShruggieTech product identity assets. Their permitted reference use and restrictions are recorded in the retained [brand asset terms](../brand/cueson/1.0.0/kit/LICENSE-BRAND.md).

The kit's code, templates, and reference documentation are supplied under Apache License 2.0 as described by its retained [license](../brand/cueson/1.0.0/kit/LICENSE) and [notice](../brand/cueson/1.0.0/kit/NOTICE). Geist, Geist Mono, and Space Grotesk remain under the SIL Open Font License 1.1, with the supplied license texts retained under the kit's [font license directory](../brand/cueson/1.0.0/kit/fonts/licenses/). The repository root `NOTICE` summarizes this bundled material without changing any of those terms.

## Offline verification

Run the standalone [brand verifier](../scripts/brand-verify/) from the repository root:

```text
go -C scripts/brand-verify test -count=1 ./...
go -C scripts/brand-verify run . -repo ../..
```

The verifier performs no network access and writes no repository content. It validates the archive identity, rejects unsafe or ambiguous archive names, proves a bijection among archive entries, import-manifest entries, and extracted regular files, verifies every size and digest, and confirms that each declared repository consumer reference resolves to its recorded retained asset.

The ordinary documentation and repository text checks remain separate and must also pass. Byte-protected brand paths are excluded from formatter repair so upstream content cannot be rewritten to conform to Cueson-authored prose rules.

## Updating the kit

A later kit revision requires its own Spec Kit decision and a new versioned directory. Download the operator-designated official archive without overwriting an existing retained snapshot, record its resolved source and acquisition identity, audit archive paths before extraction, preserve every accepted file byte for byte, generate a new independent import manifest, update consumer references intentionally, and run the full repository verification suite.

Never edit an imported snapshot in place, mix payloads from separate kit versions, or repair upstream files. Differences in the official archive require an explicit new acquisition, not an unrecorded replacement.

## Release and production boundary

The complete kit is repository source material. It is not a Cueson runtime dependency and is not a member of the executable release archives. This import does not activate `cueson.io`, update production metadata, create repository social-preview settings, publish a schema, create or move a tag, or publish a GitHub Release. Any later release inclusion or production use requires a separately specified and explicitly authorized change.
