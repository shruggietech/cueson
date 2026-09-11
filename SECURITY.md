# Security policy

## Supported versions

The published v0.0.0 release is an immutable envelope-only foundation. Current security hardening occurs on the v0.1.0 development line while the first stable v1.0.0 release remains pending. A release's publication does not extend its documented capability boundary or imply that the unreleased v1 contract is already supported.

Security corrections are applied to the current development line and assessed for any affected published release. Supported-version declarations will accompany stable releases and do not authorize rewriting an immutable release artifact or schema.

## Reporting a vulnerability

Do not open a public issue for a suspected vulnerability. Use the repository's enabled [GitHub private vulnerability reporting](https://github.com/shruggietech/cueson/security/advisories/new) and include affected inputs, observed behavior, impact, and a minimal reproduction when safe.

Subtitle files and Cue JSON documents are untrusted input. Relevant classes include path traversal, link following, unsafe overwrite, parser denial of service, amplification through cues or diagnostics, malformed base64, length or hash bypass, timestamp abuse, markup or completion-script injection, uncontrolled output, and disclosure of source bytes or local filesystem identity.

## Implemented boundaries

- Public file inputs use bounded regular-file acquisition and reject links, non-regular inputs, cancellation, and input above the documented limit.
- Exact restoration validates safe basenames, portable collisions, canonical base64, declared byte length, SHA-256, and the complete destination plan before publication.
- Model-driven rendering and conversion never execute source markup. Strict conversion rejects every known loss before target publication.
- Inspection and failure diagnostics omit source content, caller paths, usernames, hostnames, and other local identifiers.
- Deterministic collection limits reject excessive work rather than silently truncating accepted source content or loss reports.

These controls reduce risk but do not make subtitle content trusted. Run Cueson with only the filesystem access required for the intended input and destination, preserve original inputs for investigation, and include the exact command, version, platform, and safe reproducer in a private report.
