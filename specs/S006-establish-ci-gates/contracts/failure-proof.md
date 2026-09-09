# Controlled Failure-Proof Contract

## Probe

The `Repository text` job fails when `.github/ci-failure-probe` exists. The marker contains a short explanation and no executable content.

## Required sequence

1. Complete and locally verify the intended S006 implementation without the marker.
2. Add the marker in a dedicated commit.
3. Push the S006 branch and open the official pull request as a draft with `Closes #8`.
4. Wait for the hosted `Repository text` check and overall `CI` workflow to fail because of the marker.
5. Record the failed run URL and verify that no unrelated job produced a false-success override.
6. Remove the marker in a dedicated correction commit and push it.
7. Wait for every CI and CodeQL check to pass.
8. Mark the pull request ready for review so external Codex and security reviews evaluate the corrected tree.

## Completion boundary

The final branch and pull-request diff do not contain the marker. GitHub run history retains both red and green evidence. No third automated Codex round may be requested, and the AI does not merge the pull request.
