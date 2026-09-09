# Closing-Reference Contract

## Accepted evidence

A normal pull request passes when GitHub's fully paginated `closingIssuesReferences` connection contains at least one issue. GitHub remains the authority for supported closing keywords, Markdown interpretation, issue existence, and local or cross-repository resolution.

The result identifies every resolved target by repository and issue number. A non-empty list is sufficient even when the pull request closes more than one issue.

## Rejected evidence

The policy fails when GitHub resolves no closing issue. This includes placeholders, bare issue mentions, inaccessible or nonexistent targets, examples that GitHub ignores, and non-closing prose without requiring the repository to reproduce GitHub's parser.

## Exceptions

Dependabot passes the issue-link context as an explicit automated-source exception. Any other exception requires a pull-request comment with an exact logical line `skip: issue-link - <reason>`, a non-empty reason, and the configured operator login `h8rt3rmin8r`.

The policy result names the closing targets or exception comment. It never reports a bare pass without evidence.
