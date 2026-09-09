# Permissions and Event Contract

## CI

- Runs for every pull request targeting `main` and every push to `main`.
- May expose `workflow_dispatch` only for a normal green rerun; the first controlled failure proof uses pull-request history.
- Declares `contents: read` and no write permission.
- Does not consume secrets.
- Does not use `pull_request_target`.
- Disables checkout credential persistence before any repository-controlled command runs.
- Cancels a superseded run only when both runs belong to the same pull request and workflow.
- Does not upload build artifacts or mutate repository state.

## CodeQL

- Runs for every pull request targeting `main`, every push to `main`, and one low-frequency schedule.
- Declares `contents: read` and `security-events: write`.
- Uses one Go analysis job with the stable name `Analyze Go` under workflow `CodeQL`.
- Does not grant pull-request, issues, actions, packages, attestations, deployments, or contents write authority.
- Does not use `pull_request_target` or repository secrets.
- Disables checkout credential persistence before analysis runs.

## Fork behavior

Ordinary CI remains read-only for fork-originated pull requests. CodeQL uses GitHub's supported pull-request event and may rely on GitHub's built-in code-scanning handling for eligible public-repository pull requests; it does not compensate for an unavailable reporting token by broadening another workflow.
