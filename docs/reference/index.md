---
title: Reference overview
description: Exact reference material for commands, configuration, environment variables and the composite actions.
sidebar:
  order: 1
---

| Reference | Contents |
| --- | --- |
| [CLI](cli/index.md) | Every public `gh app-auth` command, with help text captured from the current release |
| [Configuration](configuration.md) | The `config.yml` schema, validation rules and storage locations |
| [Environment variables](environment-variables.md) | Every variable the code reads, with precedence and cautions |
| [GitHub Actions](github-actions.md) | Inputs, outputs and limits of the shipped composite actions |

The CLI pages are generated from `--help` output and describe v0.2.0. Hand-written pages state the
release they describe where it matters. For the reasoning behind these details, see
[Understand](../understand/credential-routing.mdx); for symptoms, see
[troubleshooting](../troubleshoot/index.md).
