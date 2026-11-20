---
title: Origins
description: The problems gh-app-auth was created to solve, and what it does and does not promise.
sidebar:
  order: 6
---

gh-app-auth began with four recurring problems in enterprise automation that touches several GitHub
organizations.

## 1. Robot accounts versus GitHub Apps

Teams that automate Git access choose between shared "robot" user accounts, which behave like
people, and GitHub Apps, which have their own identity and installation-scoped permissions. Apps
fit governance better, but their tokens are harder to use from tools that expect a username and a
password.

## 2. GitHub Apps across organizations

An installation token belongs to one installation, and an installation belongs to one account. A
person can clone a repository with submodules from several organizations using one login. Automation
needs one credential per organization.

## 3. CI complexity

CI systems that clone across organizations need submodule credential configuration, propagation of
parent credentials and per-organization token handling. That logic tends to be re-implemented in each
pipeline.

## 4. Short-lived tokens

GitHub App installation tokens expire after an hour. Long builds can fail mid-run unless a fresh token is
obtained.

## What gh-app-auth does about them

- One tool registers GitHub Apps and personal access tokens with the repository patterns they cover.
- A Git credential helper picks the credential per request and mints a fresh installation token when
  Git asks, so **each Git request** starts with a valid token.
- Git configuration is written for you, including for submodules.
- Private keys go to the OS keyring by default.

## What it does not promise

- It does not renew a token inside a running child process: `exec` takes a snapshot.
- It does not make files encrypted, and it does not make PAT or key fallback to files reliable on
  machines without a keyring. See [known limitations](../troubleshoot/known-limitations.md).
- It does not replace the need to scope each GitHub App installation to the least privilege.

Further reading: [identity and credentials](identity-and-credentials.md) and
[routing](credential-routing.mdx).
