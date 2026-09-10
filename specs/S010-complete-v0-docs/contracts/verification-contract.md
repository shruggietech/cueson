# Documentation Verification Contract

## Invocation

From repository root:

```text
go -C scripts/docs-verify run . -repo ../..
```

The command performs a read-only offline check and accepts no network or publication credentials.

## Audited scope

- Root: `README.md`, `CHANGELOG.md`, `CONTRIBUTING.md`, `SECURITY.md`, and `AGENTS.md`.
- Governance: `.specify/memory/constitution.md`.
- Documentation: every Markdown and HTML file recursively beneath `docs/`.
- Required-path enforcement: the twelve canonical documents in `documentation-contract.md`.

## Accepted references

- Inline Markdown links and images.
- Markdown reference definitions.
- Markdown reference uses and autolinks.
- HTML `href` and `src` attributes.
- Relative local paths resolved from the source document.
- Leading-slash repository-root paths.
- Fragment-only references and file-plus-fragment references.
- URL-encoded local paths and fragments.
- Explicit HTML `id` and `name` anchors.
- Existing local directories when no fragment is present.
- External absolute URLs, email links, and protocol-relative network references without fetching them.

## Rejected references

- Missing required documents.
- Local targets that do not exist.
- Local paths that escape the repository root.
- Local paths whose authored case differs from the repository entry or whose targets are symbolic links.
- Fragments that do not identify a heading anchor in the target Markdown document.
- Fragments attached to non-Markdown targets.
- Empty, malformed, or unsupported local target forms.

## Diagnostics and exit status

- Success exits `0` and reports the number of checked documents and local links.
- Contract violations exit `1`, print deterministic repository-relative, line-specific diagnostics to standard error, and do not modify files.
- Invalid invocation or unreadable repository state exits `2` with a concise diagnostic.
- Diagnostics are sorted so identical repository state produces identical output.

## Tests

Tests cover required-document inventory, relative and root paths, directory links, fragments, duplicate and setext headings, URL decoding, reference definitions and uses, autolinks, HTML badge links and anchors, code/comment exclusion, external-link non-fetch behavior, traversal, path casing, symlinks, controls, backslashes, malformed escapes, missing files, missing anchors, non-Markdown fragments, deterministic ordering, and repository acceptance.
