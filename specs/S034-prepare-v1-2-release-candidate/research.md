# S034 research

## Exact candidate identity

**Decision**: Change only current contract identity to 1.2.0 and add a byte-identical immutable candidate copy.

**Rationale**: S033 deliberately staged a future minor identity. Current/historical local dispatch already supports exact released 1.0/1.1.

**Alternatives considered**: Editing released 1.1 violates immutable history; retaining dev does not complete the selected stable candidate; admitting dev as historical invents an unapproved contract.

## Native and old-consumer proof

**Decision**: Add explicit stable 1.2 dispatch, retain six historical paths and existing 1.0 proof, and independently authenticate/execute checksummed published 1.1 host archives.

**Rationale**: Existing native gate only dispatches stable 1.1, and existing old-consumer downloader is pinned to 1.0. Stable identity promotion alone silently omits native proof. Public 1.1 identity refusal is generic exact-contract selection rather than a diagnostic containing its version.

**Alternatives considered**: A rebuilt old binary fails public provenance; overwriting frozen S020/S028/S029 files changes historical proof; requiring a literal 1.1 diagnostic fails valid exact-version refusal.

## Public state and decision package

**Decision**: Prepare source and candidate documentation at 1.2.0 while public installation and site identity stay at 1.1.0. Reconcile stale current-status assertions, preserving historical snapshots and release records.

**Rationale**: Reviewed stable candidate metadata is not publication. Current roadmap/procedure prose contains old pending #67/#76 statements despite those native issues being closed.

**Alternatives considered**: Publishing new routes/downloads prematurely crosses the withheld production boundary; retaining contradictory current status misleads consumers.

## Independent research

Schema/version, package/native proof and documentation/governance were researched by separate existing agents before implementation. No remaining technical clarification or new issue blocks the slice.
