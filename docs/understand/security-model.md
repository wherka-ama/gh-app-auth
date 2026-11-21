---
title: Security model
description: What gh-app-auth protects, where secrets live, who can still read them, and the limits of the design.
sidebar:
  order: 4
---

gh-app-auth handles private keys, personal access tokens and installation tokens. This page states
what it does about them and, just as important, what it does not. It makes no certification or
compliance claims.

To report a vulnerability, follow the project's
[security policy](../../.github/SECURITY.md).

## Where each secret lives

| Secret | Where | Protection |
| --- | --- | --- |
| App private key | OS keyring by default | Encrypted at rest by the operating system's keyring |
| App private key, filesystem source | A file you provide with `--key-file`, or a managed copy | Plain file; relies on file permissions |
| Automatic filesystem fallback | `~/.config/gh/extensions/gh-app-auth/secrets/<name>.<type>` | Raw bytes with mode `0400`; **not encrypted** |
| PAT | OS keyring | Encrypted at rest by the keyring. Filesystem retrieval is not implemented |
| Installation token | Process memory | Never written to disk; 55-minute cache per process |
| JWT | Process memory | Short-lived, not cached |

The configuration file holds identifiers, names, patterns and priorities, not secret values.

## Key file permission checks

On POSIX systems, `setup` and JWT loading reject a key file whose permission bits include group or
world **read** (`0044`). Use `600` or `400`. The check does not cover symlinks, group write, or
changes made after the check, and it is not repeated every time a configured key is read from the path
stored in the configuration. On Windows the check is skipped, and you are responsible for file ACLs.

## Boundaries and exposure

- **Same user and root.** Anything running as your user, or as root, can read your keyring entries,
  your key files, the process memory of gh-app-auth, and the environment of its children. gh-app-auth
  does not defend against that.
- **Process memory.** Tokens are plain strings in process memory. Go strings are immutable, so
  zeroization is best effort and not guaranteed.
- **Child environment.** `exec` places a token in the child's environment, and every descendant of the
  child inherits it. Treat the child as trusted for the token's life.
- **Standard output.** `token` writes a token to standard output. Terminals, shell tracing, CI logs and
  pipelines can capture it. Prefer `exec`. This is a deliberate choice recorded in
  [ADR 0002](../adr/0002-explicit-github-app-token-output.md).
- **Process arguments.** `--pat` is a command-line argument, so a PAT can appear in the process list while
  `setup` runs.
- **Config path.** `GH_APP_AUTH_CONFIG` moves the configuration file, not the secrets directory.

## Tokens and scope

Patterns select a credential; they do not restrict it. The limits on what a token can do come from the
GitHub App installation's permissions and repository selection. Grant the minimum. See
[tokens and caching](#tokens-and-caching).

## Tokens and caching

Several clocks run at once in gh-app-auth. They are separate and should not be confused.

| Item | Lifetime | Where it lives |
| --- | --- | --- |
| App JWT | `exp` is 9 minutes ahead and `iat` is 60 seconds behind the current time | Created on demand, not cached |
| Installation token | Up to 1 hour, set by GitHub | Memory only |
| Installation token cache entry | 55 minutes, per process | Memory only |
| Installation scope metadata | 24 hours | The configuration file |
| PAT | Provider policy, expiry or revocation | OS keyring |

### JWT

To act as the App, gh-app-auth signs a JWT with the App's private key. The `iat` claim is set 60
seconds in the past to tolerate clock skew, and `exp` is 9 minutes in the future. A JWT is used to
request an installation token and is not stored. See the
[ADR on the clock-skew margin](../adr/0001-jwt-clock-skew-margin.md).

### Installation token

GitHub issues an installation token for one installation, valid for up to an hour. The helper
caches it in memory for 55 minutes, which leaves a margin before expiry.

- The cache exists only inside one process. Git starts a **fresh helper process for each request**,
  so each request normally mints a new token. The cache helps processes that stay alive and make
  several requests, such as `gh app-auth test`.
- A background worker removes expired entries. It does **not** renew them. There is no daemon and no disk
  cache.
- Tokens are never written to disk.

### exec and token

`exec` takes a snapshot of a token and gives it to its child. Nothing renews it, so a child that runs
longer than the token's life starts failing. `token` bypasses the cache and mints a fresh token each
time; its output goes to standard output, which is a secret boundary. See the
[exposure boundaries](#boundaries-and-exposure).

### Scope metadata

`gh app-auth scope` records which repositories an installation can see and caches that for 24
hours. It is metadata used for matching decisions; it does not extend or shorten any token's
life, and `scope --refresh` re-reads it.

### What caching does not change

- A cached token is still bound to its installation's permissions.
- The local cache has no independent revocation and receives no push notification. GitHub authorizes
  each request, and it may reject a token before the cache's 55 minutes are up. Suspending or uninstalling
  the installation is how access is blocked; see GitHub's
  [guide to reviewing and modifying installed GitHub Apps](https://docs.github.com/en/apps/using-github-apps/reviewing-and-modifying-installed-github-apps)
  and the [REST reference for revoking an installation access token](https://docs.github.com/en/rest/apps/installations#revoke-an-installation-access-token).
- Restarting a process starts with an empty cache.

## Logging

Diagnostic logging is off unless `GH_APP_AUTH_DEBUG_LOG` is set. When on, log lines record flow steps
and hash tokens instead of printing them, and obvious secret-looking values are redacted. Logs still
contain command-line arguments, repository URLs and file paths, so review them before sharing. See
[diagnostics](../troubleshoot/index.md#diagnostics).

## Operational guidance

- Prefer a GitHub App over a PAT, and the keyring over files.
- In CI use disposable runners and the explicit-key-file layout from the [CI guide](../guides/ci/index.md).
- Rotate keys and PATs on a schedule, and remove credentials you no longer use; see
  [credential lifecycle](../guides/credentials.md#maintain-credentials).

## Not covered

- Defense against malware or a compromised machine.
- Protection of a key file after it has been copied elsewhere.
- Any guarantee about third-party credential helpers configured next to gh-app-auth.
