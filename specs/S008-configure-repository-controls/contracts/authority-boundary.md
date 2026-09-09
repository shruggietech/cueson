# Authority Boundary Contract

## Authorized actions

- Create and use branch `S008-configure-repository-controls`.
- Push S008 commits and publish the official pull request closing issue #10.
- Mutate settings and repository-owned rules for `shruggietech/cueson` that fall within issue #10 and the approved S008 specification.
- Enable supported security capabilities for `shruggietech/cueson`.
- Request at most one second Codex review after first-round findings are resolved.

## Prohibited actions

- Editing or replacing organization-owned ruleset `20478126` or any other organization-wide policy.
- Mutating settings, rules, or security features for another repository.
- Merging or auto-merging the S008 pull request.
- Creating or moving a tag, publishing a release or schema, or changing production `cueson.io` configuration.
- Triggering a third Codex review.
- Expanding bypass authority to a bot, deploy key, contributor role, or unrelated team.

## Recovery authority

The new repository-owned ruleset retains one `OrganizationAdmin` bypass in `always` mode. This is a recovery boundary, not normal delivery flow. Ordinary S008 completion still stops for the operator's final review and merge ritual.
