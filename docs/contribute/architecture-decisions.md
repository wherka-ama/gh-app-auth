---
title: Architecture decisions
description: How decisions are recorded as ADRs, and the index of the accepted decisions with links to the originals.
sidebar:
  order: 8
---

Significant technical decisions are recorded as numbered, immutable Architecture Decision Records
(ADRs) in [`docs/adr/`](../adr/README.md). That README holds the template, naming rules and the
status lifecycle. An accepted ADR is never rewritten; a new decision supersedes it with a new ADR.

Write an ADR for breaking changes, new features that introduce a concept or trust boundary,
security decisions, technology choices, significant refactoring, and rejected alternatives worth
recording. Bug fixes, test additions, documentation-only changes and formatting do not need one.

## Accepted decisions

| ADR | Decision |
| --- | --- |
| [0001: JWT clock-skew margin](../adr/0001-jwt-clock-skew-margin.md) | Leave a margin on the JWT timestamp claims |
| [0002: Explicit GitHub App token output](../adr/0002-explicit-github-app-token-output.md) | Printing a token is an explicit, separate command |
| [0003: Draft-first gated release pipeline](../adr/0003-draft-first-gated-release-pipeline.md) | Stage releases as drafts and publish only after gates pass |

These pages link to the original documents on GitHub so the immutable text is not duplicated here.
The documentation site itself did not need an ADR.
