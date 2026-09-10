# Research: Import the Official Cueson Brand Kit

## Decision 1: Pin the live hosted archive selected by the operator

**Decision**: Use `https://brand.shruggie.tech/cueson/downloads/cueson-brand-1.0.0.zip` as the sole S011 acquisition authority. Two independent downloads on 2026-09-10 matched at 3,889,985 bytes and SHA-256 `095e572a0db73472c9fdcf5001f68e702105a3d0ec1fcc34f830ae57c608cd07`.

**Rationale**: The operator owns both projects and explicitly designated the published ZIP as the comprehensive handoff. The live acquisition supersedes stale hashes and build-comparison gates in issue #23.

**Alternatives considered**: Rebuilding from the upstream repository or blocking on the older certified payload would contradict the operator decision and add no consumer value. Mixing old and current payloads would destroy snapshot integrity.

## Decision 2: Retain the archive and every safe payload file

**Decision**: Commit the exact ZIP under `brand/cueson/1.0.0/archive/` and its complete, unmodified extraction under `brand/cueson/1.0.0/kit/`.

**Rationale**: The user requested the whole kit in the repository. The current archive contains 265 regular files totaling 5,469,058 extracted bytes with no absolute, traversal, drive-qualified, backslash, link-like, duplicate, case-colliding, Windows-reserved, or trailing-dot paths.

**Alternatives considered**: Importing only the old handoff's 11 destinations would fail the complete-kit request. Keeping only the ZIP would make assets inconvenient to use and review. Keeping only extracted files would lose the exact delivered artifact.

## Decision 3: Add an independent import manifest

**Decision**: Generate `import-manifest.json` with acquisition metadata, the archive identity, all 265 entry paths, sizes and SHA-256 digests, the payload root, and direct repository references.

**Rationale**: The kit's internal manifest is useful but covers 261 entries and its preserved `VERIFY.md` reports an obsolete count of 246. An independent Cueson manifest accounts for the archive as actually received without modifying upstream content.

**Alternatives considered**: Trusting only upstream prose would leave four legal and manifest files outside the proof. Editing the kit's records would violate byte preservation.

## Decision 4: Add a standalone offline brand verifier

**Decision**: Add a Go 1.25 standard-library module at `scripts/brand-verify` that strictly parses the import manifest, verifies archive identity, rejects unsafe catalog entries, proves the archive-manifest-extraction bijection, rejects non-regular extracted content, and checks every declared documentation reference without writing or fetching the network.

**Rationale**: Brand integrity is a repository concern, not shipped Cueson behavior. A focused portable verifier gives deterministic local and CI evidence and can safely stream every entry rather than trusting filenames or visual similarity.

**Alternatives considered**: Extending product code would blur the runtime boundary. Shell-only verification would be less portable and harder to unit test. Online verification in every CI run would make the official host's availability part of repository correctness.

## Decision 5: Reference retained assets directly

**Decision**: Use the retained dark and light horizontal SVGs in the README, and use the retained dark horizontal SVG, favicon, and font stylesheet in the media-format guide. Do not create duplicate consumer copies.

**Rationale**: Direct references keep one authoritative byte source while delivering offline rendering. The dark-surface `color` lockup contains light ink and requires the `light` variant on GitHub's light theme. The guide can consume the kit's relative font stylesheet and WOFF2 files unchanged.

**Alternatives considered**: Copying files to a second asset tree would add drift surfaces without serving a runtime packaging need. Hotlinking the brand host would make core identity network-dependent. Wiring the kit's web manifest into the standalone guide would create broken root-relative icon URLs.

## Decision 6: Protect immutable content from the repository formatter

**Decision**: Teach `scripts/github-format` to skip `brand/`, `docs/assets/brand/`, `docs/assets/favicons/`, and `docs/assets/fonts/` in both check and fix modes, with regression tests proving protected bytes are untouched.

**Rationale**: The formatter currently classifies many kit extensions as mutable text and empirically rejects seven exact upstream files. Running repair mode could change canonical bytes despite `.gitattributes` correctly disabling Git normalization.

**Alternatives considered**: Removing those extensions from text classification would weaken authored-file validation elsewhere. Making the kit conform to Cueson prose rules would corrupt source truth.

## Decision 7: Preserve bundled terms without a separate legal-review gate

**Decision**: Retain the kit's Apache-2.0, reserved-mark, NOTICE, and SIL Open Font License materials byte for byte, summarize their boundaries in `docs/brand.md`, and update the root NOTICE for bundled content.

**Rationale**: Ownership makes an external legal approval checkpoint redundant, as the operator directed, but retaining the supplied license and attribution facts is part of faithful redistribution.

**Alternatives considered**: Dropping legal files would make the retained kit incomplete. Treating reserved marks or OFL fonts as ordinary Apache-licensed code would state a false license boundary.

## Decision 8: Keep release and production membership unchanged

**Decision**: Do not change `.goreleaser.yaml`, the four-member release archive contract, repository social-preview settings, `cueson.io`, or production metadata in S011.

**Rationale**: The release configuration explicitly includes only the executable, schema, `LICENSE`, and `NOTICE`, and the release verifier enforces that membership. The user requested a repository import, not a product release or domain activation.

**Alternatives considered**: Adding the complete kit to every binary archive would inflate distributions and change legal and SBOM surfaces without runtime need. Publishing site or repository metadata would require separate surface-specific decisions.

## Decision 9: Simplify issue #23 under the operator override

**Decision**: Replace the obsolete multi-gate issue body with the approved complete-import outcome, retain independently testable safety, integration, documentation, CI, and review acceptance criteria, and remove the `needs: decision` label.

**Rationale**: GitHub issue state is a current planning authority and must not continue to contradict the explicit operator decision or the S011 specification.

**Alternatives considered**: Leaving the stale body would make the pull request appear noncompliant with requirements the operator expressly rejected. Closing the issue before merge would claim default-branch completion prematurely.
