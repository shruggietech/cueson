# Publication Requirements Checklist: S035

**Purpose**: Reviewer-owned requirements quality before implementation.

**Created**: 2026-10-06

**Feature**: [spec.md](../spec.md)

**Marker semantics**: A checked item means the requirement is clear and complete, not that implementation has finished.

- [x] CHK001 Are public-release identity and seven download expectations explicit and measurable? [FR-001]
- [x] CHK002 Are all four immutable schema identities and historical preservation requirements stated? [FR-002, FR-003]
- [x] CHK003 Are release/speaker navigation and prior routes covered? [FR-004, FR-005]
- [x] CHK004 Does prose distinguish public software, prepared source and pending live deployment? [FR-006]
- [x] CHK005 Are consumer identifier scope and duration limitations required without new product restrictions? [FR-007]
- [x] CHK006 Are relevant negative tests and existing credential/revision boundaries retained? [FR-008, FR-009]
- [x] CHK007 Does the later handoff require fresh main and separate live address-family evidence? [FR-010]
- [x] CHK008 Are issue closure, milestone and push/review/merge/production boundaries unambiguous? [FR-011, FR-012]

## Independent reviewer notes

Reviewed by the independent site research agent after S035 scope alignment on 2026-10-06. All eight requirements-quality checks pass; these markers do not claim implementation, CI, merge or production completion.

CHK001: FR-001, SC-002 and the data model explicitly bind the published v1.2.0 release and seven official platform/checksum URLs.

CHK002: FR-002/FR-003 and SC-001 require all four immutable identities. The data model pins the new 191170-byte v1.2.0 digest; existing content-map and independently fixed site tests retain the three historical sizes, digests and paths. T008/T009 preserve and verify those identities.

CHK003: FR-004/FR-005, SC-002/SC-003 and T004/T005/T008/T012 cover both new public documents, navigation, four prior/current release pages, historical links and the exact 28-route inventory.

CHK004: FR-006, SC-004 and T006/T007/T013 distinguish already published software from prepared source and actual live activation, preserving dated records.

CHK005: FR-007 and US1 acceptance distinguish consumer-governed IDs, existing string limits, independent native observations, optional declared duration and the absence of inferred audio-duration claims. No new product-usage restriction is introduced.

CHK006: FR-008/FR-009 and T004/T008/T010/T012 require independent inventories and negative tests while retaining manual deployment, exact-main freshness, step-scoped credentials and infrastructure preservation.

CHK007: FR-010, T016-T018 and the production handoff require actual merged main, fresh proof, explicit header probes and independent IPv4/IPv6 HTTPS evidence. The S022 DNS helper accepts one usable family and is not sufficient evidence of HTTPS connectivity over both; T018 expressly requires the additional probes.

CHK008: FR-011/FR-012, SC-005, T014-T018 and the recorded S035-only authority agree: push/PR and post-owner-merge production work are authorized, final specific PR merge remains human-owned, and #86/milestone stay open until live acceptance passes. No second deployment kickoff is required.
