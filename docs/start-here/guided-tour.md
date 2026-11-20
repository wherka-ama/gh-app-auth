---
title: Guided tour
description: The shortest route from the problem gh-app-auth solves to the deeper reference, in five steps.
sidebar:
  order: 5
---

This tour is a route through the documentation, not another manual. Each step names the one page
to read next.

## 1. The problem

Several organizations, several GitHub Apps, one machine, and tokens that expire. Git asks for
credentials per request, so something has to choose the right one each time. The
[overview](overview.md) states what gh-app-auth does and where it stops, and
[origins](../understand/origins.md) explains why the project exists.

## 2. Identity

A credential is either a GitHub App (identified by Client ID or App ID, plus a private key and an
installation) or a personal access token. They differ in lifetime, scope and storage. Read
[identity and credentials](../understand/identity-and-credentials.md).

## 3. Routing

Every credential carries patterns. For GitHub Apps the longest matching path prefix wins; Apps and
PATs are then compared by priority. Read [how credential routing works](../understand/credential-routing.mdx).

## 4. Setup

Register the credential, write the Git configuration, then check it:

1. [Make your first authenticated clone](first-clone.md) for a workstation.
2. [Use it in CI](../guides/ci/index.md) for an ephemeral job.
3. [Run a command with exec](first-exec.md) when you need an API call instead of a clone.

## 5. The deeper documentation

| You want to | Read |
| --- | --- |
| Know every flag | [CLI reference](../reference/cli/index.md) |
| Know every config field | [Configuration reference](../reference/configuration.md) |
| Understand the security limits | [Security model](../understand/security-model.md) |
| See the moving parts | [Architecture](../understand/architecture.mdx) |
| Fix a problem | [Troubleshooting](../troubleshoot/index.md) |
| Change the project | [Contributing](../contribute/index.md) |

The project also keeps an optional slide deck in
[`docs/presentation.md`](../presentation.md); it points back to these pages.
