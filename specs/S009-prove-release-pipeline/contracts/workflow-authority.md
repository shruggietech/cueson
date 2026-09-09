# Non-Publishing Workflow Authority

## Entry points

The hosted snapshot proof runs for pull requests targeting `main` and explicit workflow dispatch. It does not run on tag, release, or privileged target-branch events.

## Permissions and trust

- Workflow permissions are exactly `contents: read`.
- Checkout credentials are not persisted.
- No repository, environment, organization, or personal secret is referenced.
- No `GITHUB_TOKEN` is passed to GoReleaser, Syft, the verifier, or another command.
- Repository-controlled pull-request code runs only under the ordinary unprivileged `pull_request` event.
- Every referenced workflow action is GitHub-owned and pinned to a full commit identifier.
- Tool commands are installed from exact Go module versions and never from a floating selector.

## Required sequence

1. Check out complete source history without persisting credentials.
2. Install the declared Go compatibility line.
3. Install exact GoReleaser and Syft commands into a temporary tool directory.
4. Validate the GoReleaser configuration.
5. Run `release --snapshot --clean --skip=publish` without a publication token.
6. Run the repository-owned release verifier against the complete `dist/` tree and current revision.
7. Upload the verified `dist/` tree as a short-lived workflow artifact for review.

The upload step runs only after verification succeeds. A failed proof remains a failed check and does not retain a misleading accepted artifact set.

## Explicitly prohibited operations

S009 has no path that:

- creates, updates, or pushes a Git tag;
- creates, edits, publishes, or deletes a GitHub Release;
- uploads an asset to a GitHub Release;
- writes repository contents, issues, pull requests, checks, packages, security events, attestations, or identity tokens;
- imports a signing key or invokes artifact signing;
- publishes provenance or an attestation;
- modifies DNS, hosting, or content under `cueson.io`;
- enables auto-merge or merges the pull request.

## Future release boundary

A real release workflow requires its own specification, exact release identity, operator authorization, protected environment or equivalent approval boundary, release-note contract, signing and attestation decision, tag ownership, and post-publication verification. S009's workflow must not be extended into that path by merely adding a token or removing `--snapshot`.
