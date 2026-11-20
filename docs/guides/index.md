---
title: Guides overview
description: Task-oriented guides for creating credentials, configuring Git, working across organizations and running in CI.
sidebar:
  order: 1
---

Each guide covers one task from start to finish. Look up exact flags in the
[CLI reference](../reference/cli/index.md) and the reasoning in [Understand](../understand/identity-and-credentials.md).

| Guide | Use it to |
| --- | --- |
| [Create and install a GitHub App](create-github-app.md) | Prepare the App, permissions, installation and private key that gh-app-auth needs |
| [Personal access tokens](personal-access-tokens.md) | Register a PAT, including a username for Bitbucket-style hosts |
| [Git credential helper](git-credential-helper.md) | Write, inspect and remove the Git configuration the helper relies on |
| [Several organizations and submodules](multiple-organizations.md) | Route different organizations to different Apps |
| [GitHub Enterprise Server](github-enterprise-server.md) | Use hosts other than github.com |
| [Credential lifecycle](credential-lifecycle.md) | Verify, migrate, rotate and remove credentials |
| [Continuous integration](ci/index.md) | Run safely in GitHub Actions, Jenkins and GitLab CI |
