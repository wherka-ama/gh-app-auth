---
title: Contributing overview
description: How to contribute code, tests, documentation and releases to gh-app-auth, with links to each maintainer guide.
sidebar:
  order: 1
---

Thank you for helping with gh-app-auth. Start with the guide that matches what you are changing.

| You are changing | Read |
| --- | --- |
| Go code | [Development](development.md), then the [testing guide](testing.md) |
| Behavior that needs live credentials | [End-to-end testing tutorial](e2e-testing.md) and [test infrastructure](e2e-infrastructure.md) |
| Documentation | [Documentation guide](documentation.md) |
| A release | [Release runbook](releasing.md) |
| A design decision | [Architecture decisions](architecture-decisions.md) |

The repository-level entry point is [CONTRIBUTING.md](../../CONTRIBUTING.md), which covers the pull
request process and commit message format. Maintainers are listed in
[MAINTAINERS.md](../../MAINTAINERS.md). Report vulnerabilities privately; see the
[security policy](../../.github/SECURITY.md).

## Ground rules

- Never log, print or commit tokens, private keys or test credentials.
- Add or update tests with the change, and keep tests table-driven where it helps.
- Update the documentation with the behavior change. Mark behavior that is not yet released as
  unreleased.
- Ask whether a breaking change or a new feature needs an ADR before you finish it.
