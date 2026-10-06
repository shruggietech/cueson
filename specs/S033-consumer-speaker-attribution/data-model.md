# S033 Data Model

## Consumer attribution

cues[].speaker_attributions is optional, max1024 per cue, ordered and repeated entries permitted. Each closed object requires speaker_id and optionally paired start_milliseconds/end_milliseconds. Missing/empty arrays mean no recorded consumer assignment; null is structurally invalid. Untimed means participation somewhere in the cue. Positive half-open timed intervals use checked nonnegative int64 milliseconds inside cue and supplied media bounds. IDs are opaque1..256 decoded Unicode scalars; exact frozen exclusions and whitespace set are those in issue #82, preserved without normalization.

## Media timing

Optional closed root media_timing requires nonnegative duration_milliseconds and optional signed timeline_start_milliseconds. Omitted start evaluates as0 but remains omitted on serialization. End=start+duration uses checked int64 arithmetic. Null/partial/unknown fields reject. Duration0 is known empty media, never unknown. The consumer owns resource/timeline association; this declaration is not measured source truth.

## Runtime assessment

Pure helper produces timed/untimed attribution counts, media declaration presence, checked/unavailable/not_evaluated state and cue conflict count. Timed validity is enforced by Document.Validate. Nonfatal cue conflict warnings are computed, bounded and content-free. Do not persist them into source-native Diagnostics or alter Stats.

## Native export

Every attribution occurrence and declaration omitted from native output produces a field-specific bounded deterministic runtime observation. Caller IDs do not enter diagnostics/loss attributes. Source-aware private targets clear consumer additions after original validation.

## Version dispatch

Current1.2.0-dev accepts additions; historical1.1.0/1.0.0 reject them even in direct typed revalidation. Historical1.1.0 scripted support and1.0.0 limitations remain unchanged. Immutable local registry selects exact identity pairs and never fetches input-controlled URIs.
