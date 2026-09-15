# S026 data model

## Immutable source document

The original validated Document owns schema/producer identity, native/common source observations, and complete SourceEnvelope assets. Conversion treats it as immutable and verifies all assets before analysis. Existing historical text documents retain their identity.

## Matrix analysis and private target

matrixAnalysis retains Translations and Losses and may carry a constructed Target. Text targets use one ordered translation per source cue. Scripted target builders construct or deep-copy native owners, rederive their common projections, and return complete losses. The private scripted target uses current development identity and experimental capability without being published as Cue JSON; its source envelope remains exactly equal to the original.

## Loss entries

Existing Report/Loss structures remain the only runtime model. Fourteen appended codes account for scripted metadata, style fields, event fields, records, sections, attachments, overrides omitted/degraded, override comments, drawings, changed centisecond endpoints, variant fields omitted/degraded, and common scripted source identifier annotations. Source orders resolve against common cues, diagnostics, WebVTT blocks, and scripted physical headers/records. Messages and contexts remain bounded and controlled.

## Target native owners

Scripted sections/records/styles/events/attachments retain the existing schema shapes. Text-source targets create deterministic canonical owners without raw source capture references. Variants preserve unchanged captured content at original physical positions, remove changed captures, map dialect declarations and alignment, and account for lost observed fields or presentation. Final native/common projections must agree.

## Validation transitions

Original input validation and integrity -> bounded compatibility/target construction -> complete canonical report -> strict policy -> source-aware target validation/render -> bounded target reparse and semantic comparison -> existing atomic CLI publication. Any failure publishes zero target payload; strict loss analysis returns its complete report.
