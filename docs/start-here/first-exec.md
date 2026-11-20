---
title: Run your first command with exec
description: Use gh app-auth exec to run one command with a credential token, and understand what the child process can see.
sidebar:
  order: 4
---

`gh app-auth exec` runs a command with a token from a credential you already configured: a short-lived
installation token for a GitHub App, or the stored token for a PAT. It is the
preferred way to call the GitHub API or `gh` from a script, because the token goes into the child
process environment instead of being printed.

This page assumes you finished [your first authenticated clone](first-clone.md).

## Run a command

From inside a clone of a repository covered by a configured App:

```bash
gh app-auth exec -- gh api repos/{owner}/{repo} --jq .full_name
```

Outside a repository, name the target:

```bash
gh app-auth exec --repo github.com/example-org/private-repo -- gh pr list
```

Expected outcome: the command runs as the GitHub App installation and prints its normal output. If
no credential matches the repository, `exec` fails before running the command.

## What the child process receives

- The token appears only in the child's environment. gh-app-auth does not print it.
- Before the child starts, `exec` removes the four GitHub token variables it knows about
  (`GH_TOKEN`, `GITHUB_TOKEN`, `GH_ENTERPRISE_TOKEN`, `GITHUB_ENTERPRISE_TOKEN`), then injects the
  token as `GH_TOKEN` for github.com or `GH_ENTERPRISE_TOKEN` for other hosts, together with
  `GH_HOST` and, when a repository is known, `GH_REPO`.
- Anything the child starts inherits that environment. Treat the child as trusted for the lifetime
  of the token.

## Lifetime and selection rules

- An installation token is a snapshot. It is created when `exec` starts and is not renewed while
  the child runs, so a long-running child fails once the token expires (GitHub expires installation
  tokens after one hour). A PAT is static and lasts until it expires or is revoked.
- Selecting an App explicitly with `--app-id`, `--client-id` or `--installation-id` skips PATs. The
  selector must match a configured entry, and `--repo`, when given, must match that App's route.
- For a repository-independent API call, select the App and installation explicitly:

```bash
gh app-auth exec --client-id Iv1.example --installation-id 12345678 -- gh api /installation/repositories
```

Every value above is a placeholder. For the complete flag list see the
[`exec` reference](../reference/cli/exec.md).

## When you need the raw token

`gh app-auth token` prints a fresh installation token to stdout. Standard output is a secret
boundary: shell tracing, terminal capture and pipelines can all expose it. Prefer `exec` unless a
tool genuinely needs the token as a string. `token` needs an explicitly configured App plus a
matching `--repo`, and it bypasses the in-memory cache to mint a fresh token. See the
[security model](../understand/security-model.md).
