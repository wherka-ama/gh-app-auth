---
title: Overview and limitations
description: What gh-app-auth does, what it deliberately does not do, and the limits to know before you rely on it.
sidebar:
  order: 1
---

gh-app-auth is a [GitHub CLI](https://cli.github.com/) extension that doubles as a Git credential
helper. It lets one machine or pipeline use several GitHub Apps, and optionally personal access
tokens, and picks the right credential for each repository URL.

This documentation is maintained alongside the code in the repository. See the
[changelog](../../CHANGELOG.md) for release history, and run `gh app-auth <command> --help` for the
help of the version you have installed.

## What it does

- Registers GitHub Apps (by Client ID or App ID) and personal access tokens (PATs) with the
  repository patterns they cover.
- Answers Git's credential requests by selecting a credential for the remote URL. For an App, it
  signs a JWT with the App's private key and exchanges it for an installation token.
- Runs commands with a token in their environment (`exec`): a short-lived installation token for a
  GitHub App, or the stored token for a PAT. It can also print an App token explicitly (`token`).
- Writes the Git configuration that makes Git call the helper (`gitconfig --sync`).

## What it does not do

- It does not run as a daemon, and it does not renew tokens in the background.
- It does not write installation tokens to disk.
- It does not limit what a token can do. GitHub App installation permissions do that. Patterns only
  choose which credential is sent.
- It does not make your key storage encrypted by itself. The OS keyring is the preferred store; a
  filesystem store is plain files.

## Limits to know first

| Limit | Detail |
| --- | --- |
| Automatic filesystem fallback for App keys | Can leave a configuration that is rejected later. Use `--use-filesystem` together with `--key-file` instead. |
| PATs without a keyring | Reading a PAT from the filesystem store is not implemented. Keep PATs on machines with an OS keyring. |
| Long-running `exec` children | The token is a snapshot taken at start. It is not renewed while the child runs. |
| `--config` and `--debug` flags | They are registered on every command but are not wired the way you might expect. Use `GH_APP_AUTH_CONFIG` and `GH_APP_AUTH_DEBUG_LOG`. |
| Windows key file permissions | The permission check on key files is skipped on Windows. File ACLs are your responsibility. |

Each limit is explained in [known limitations](../troubleshoot/known-limitations.md).

## Where to go next

1. [Install the extension](installation.md).
2. [Make your first authenticated clone](first-clone.md).
3. Try [`exec`](first-exec.md) for a one-off command with a credential token.
4. Read the [guided tour](guided-tour.md) for the shortest route through the deeper documentation.
