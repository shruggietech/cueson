# Cueson source and handoff notes

- Gate 1 approval binds `cueframe-r1` geometry SHA-256 `779a7c79ec00a6b62b83468101f6ec6606804833d8c1fbd0b77b5d4baa795a49` and the revised `cue-teal-r1` palette SHA-256 `15bb46c278449a9b7dbd65238bcd5566abaaf4dcf60c4cfc56bb7a7ecc7c7be2`. The earlier Cue Iris approval is retained as superseded history in the S026 approval ledger.
- Geometry provenance is `glyphkit`. `build/mk_paths.py` is the producing source, and its full and reduced arrays must remain exactly equal to `brand.json`.
- The empty approval-ledger source binding is intentional because Cueson is a newly constructed identity with no imported authoritative mark. Gate 1 drift is instead detected through `derivative_config_sha256`.
- Generated `measured` and `color` blocks are absent from source `brand.json` and restored in ignored output by the generator.
- Gate 2 retry 2 was approved on 2026-09-09. The exact derivative approval manifest SHA-256 is `e9f3aef341ed8910426abd6d7876dac1f39f69afe370efc9c1cddb7af41f9dc5`; approved surfaces are the showcase card, brand landing page, guideline topics, downloads, registry endpoints, public metadata, structured data, and social preview.
- Full and Reduced standalone framing uses declared presentation padding of 87 and 97 units respectively. This framing is distinct from the 48-unit external clear-space rule and preserves the Gate 1 composition.
- A later Cueson repository issue and Spec Kit slice must consume `consumer-handoff.json`, verify every source and artifact hash, copy only declared destinations, and preserve license notices. S026 does not modify `A:\Code\cueson` or activate `cueson.io`.
