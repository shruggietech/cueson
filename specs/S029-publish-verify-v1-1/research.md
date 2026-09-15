# Research: S029 release preparation

## Source and authority

S028 merged at 6a55b48c654e14f5baaf89a9add51fa4f5a4d04a; all four main workflows/three package smokes passed. Fresh kickoff found only #51/#66/#67 before the coordinator created preparation #74. GitHub returned 404 for v1.1.0 tag/release, establishing absence rather than mutation authority.

## Independent findings

Release-owner research confirmed existing verifier negative coverage already rejects source, target, inventory, checksums, schema/legal, path and clean-build drift. Use a normative contract and focused preparation policy checks rather than a new publisher, schema evaluator or duplicate integrity tests.

Documentation research identified outdated readiness wording postponing dated history until publication, duplicated Changed headings, and prospective notes unsuitable for exact public publication. Finalize metadata before tag selection, preserve prospective status and freeze a separate public notes body.

## Chosen sequencing

Preparation PR closes #74 only. Human merge produces a new source revision; fresh accepted main proof then supplies the machine decision record and exact action-specific approval gate for #66. Public hosting #67 remains blocked. This explicitly narrows the earlier S029 publication grouping to the work authorized now without claiming the publication outcome complete.

## Alternatives rejected

Publishing S028 main before metadata finalization omits dated history from the tag. Predicting a future squash SHA or approving stale assets breaks source binding. Publishing candidate notes after silently removing banners breaks exact reviewed-body identity. Treating branch push authority as release/production authority conflicts with the constitution.
