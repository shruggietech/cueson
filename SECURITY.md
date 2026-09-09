# Security policy

## Supported versions

Cueson has not published a supported release yet. Security fixes and disclosure expectations will be versioned before the first public release.

## Reporting a vulnerability

Do not open a public issue for a suspected vulnerability. Use [GitHub private vulnerability reporting](https://github.com/shruggietech/cueson/security/advisories/new) and include affected inputs, observed behavior, impact, and a minimal reproduction when safe.

Subtitle files and Cue JSON documents must be treated as untrusted input. Reports involving path traversal, unsafe restoration, parser denial of service, malformed base64, integrity bypass, or unintended disclosure of local filesystem information are especially relevant.
