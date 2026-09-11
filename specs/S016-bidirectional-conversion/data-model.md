# Data Model: Deliver Bidirectional Subtitle Conversion

## Conversion request

| Field | Type | Rules |
|---|---|---|
| Input document | Validated Cue JSON document | Must identify SubRip or WebVTT and pass source-envelope integrity validation |
| Source format | Canonical format | Derived from the document after explicit or automatic input classification |
| Target format | Canonical format | Required, installed, and different from source format |
| Strict | Boolean | Blocks before rendering when the complete report contains any loss |

The request is immutable during conversion. Native-source loading occurs before this entity is created and uses the same bounded acquisition and decoder authority as encode.

## Loss

| Field | Type | Rules |
|---|---|---|
| Code | Stable token | Known compatibility-matrix code; independent of message wording |
| Severity | Enum | `warning` for every recoverable representational loss |
| Kind | Enum | `omitted`, `degraded`, or `ambiguous` |
| Message | String | Nonempty, bounded, path-free, and free of source payload excerpts |
| Source format | Canonical format | `subrip` or `webvtt` |
| Target format | Canonical format | The other supported format |
| Path | JSON Pointer | Nonempty RFC 6901 pointer locating the source semantic |
| Source order | Optional integer | Resolves to the associated cue or WebVTT block when present |
| Cue ID | Optional string | Resolves to the associated cue when present |
| Context | Ordered attributes | Unique names, fixed key order, bounded values, and no source excerpts or machine identifiers |

Every entry represents one atomic semantic occurrence. Loss identity is its code plus source path and fixed occurrence context; duplicate identities are invalid.

## Loss report

| Field | Type | Rules |
|---|---|---|
| Losses | Ordered loss list | Complete before strict evaluation, renderer invocation, or publication |

The canonical order is fixed document-rule rank, then numeric body source order, then per-item rule rank, then occurrence index. Validation rejects unknown codes, invalid pairs, unresolved associations, duplicate identities, or noncanonical order.

## Conversion result

| Field | Type | Rules |
|---|---|---|
| Bytes | Byte sequence | Canonical target subtitle bytes, UTF-8 without BOM, LF line endings, one final LF |
| Loss report | Loss report | Complete and valid; empty for loss-free conversion |
| Diagnostics | Ordered diagnostics | Target-renderer conformance observations, distinct from representational loss |

A strict-loss error retains the complete report but returns no target bytes. Fatal validation, projection, renderer, and context errors are typed failures rather than loss entries.

## Private target document

The converter builds a transient target-native document containing cloned common cues, canonical target-native fields, target capability declarations, recomputed cue and document summaries, and no source-format native data. It may retain the input source envelope only to satisfy internal invariants, is never serialized as Cue JSON, and is discarded after canonical rendering.

## State transitions

```text
input bytes
  -> bounded classification and decode
  -> validated immutable source document
  -> complete target projection and loss report
  -> strict decision
     -> loss present: strict-loss error, no renderer, no publication
     -> no blocking policy: canonical target rendering
  -> stdout write or transactional file publication
```
