# Hosted Evidence Contract

## S006 checks

Before ruleset creation, query the official S008 pull request's current head and require successful check runs for all contexts named by `specs/S006-establish-ci-gates/contracts/check-contract.md`. Record each exact name, head SHA, conclusion, GitHub App identity and identifier, start time, completion time, and run URL.

A workflow name copied from source, a result on another head, a result on the merge ref, or an additional unstabilized summary context does not qualify.

## S007 policy statuses

Query combined commit statuses on the current S008 head. `Cueson PR policy / Issue link` and `Cueson PR policy / Codex review` must be created by the authenticated GitHub Actions identity, apply to the exact head, and expose the expected current state and target URL.

The S007 activation contract additionally requires a marked second-round request comment authored by `github-actions[bot]`, plus the matching durable reservation and post-comment status read-back. If the first Codex review is clean and the second-round path is correctly unused, the proof is incomplete and both policy statuses remain non-required.

## Ruleset selection

The required-check list is assembled only after the complete evidence query. Check runs use their exposed GitHub App integration identifier. Commit statuses use an integration identifier only when GitHub exposes and accepts one for the actual status provider. No copied, inferred, or placeholder provider identifier is permitted.

## Final-head record

After the last source push and review remediation, publish one formatted issue or pull-request comment containing the final head SHA, required contexts, excluded contexts with reasons, successful run URLs, policy-status source evidence, ruleset identifier, and read-back timestamp. Read the published body back immediately and verify its Markdown structure.
