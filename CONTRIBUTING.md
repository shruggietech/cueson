# Contributing to Cueson

Cueson is in pre-release bootstrap. Start with an issue so the intended outcome, dependencies, and verification can be agreed before implementation begins.

## Development workflow

1. Select or create an atomic GitHub issue with explicit acceptance criteria and verification.
2. Group compatible issues into a coherent work slice when they share one implementation and review story.
3. Use the repository's Spec Kit workflow for non-trivial product work.
4. Update tests, documentation, and `CHANGELOG.md` alongside behavior changes.
5. Run the repository verification commands in the foreground.
6. Open a pull request containing a complete closing reference such as `Closes #123`.

Repository-authored Markdown uses one source line per paragraph and per list item. Before publishing GitHub text, pass it through the repository formatter and publish from a UTF-8 file when practical.

An AI agent may prepare and verify a pull request, but the final merge remains a human decision unless the operator grants one explicit, pull-request-specific override.
