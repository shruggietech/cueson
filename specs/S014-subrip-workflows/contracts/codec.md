# Native Codec Contract

- Canonical identities are `subrip` and `webvtt`; aliases are `srt` and `vtt`.
- Detector, decoder, and renderer capabilities are independently optional and registry behavior is deterministic.
- Exact restoration is not a codec method.
- Explicit format wins; otherwise content evidence wins over extension and ambiguity fails.
- Decoders consume exact bounded bytes and never reopen the source path.
- Successful decoding returns a valid document with those exact source bytes.
- Tolerated deviations produce stable diagnostics; fatal failures return no publishable document.
- Rendering consumes a valid structured model and emits deterministic platform-independent bytes.
