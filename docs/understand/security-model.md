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
[tokens and caching](tokens-and-caching.md).

## Logging

Diagnostic logging is off unless `GH_APP_AUTH_DEBUG_LOG` is set. When on, log lines record flow steps
and hash tokens instead of printing them, and obvious secret-looking values are redacted. Logs still
contain command-line arguments, repository URLs and file paths, so review them before sharing. See
[diagnostics](../troubleshoot/diagnostics.md).

## Operational guidance

- Prefer a GitHub App over a PAT, and the keyring over files.
- In CI use disposable runners and the explicit-key-file layout from the [CI guide](../guides/ci/index.md).
- Rotate keys and PATs on a schedule, and remove credentials you no longer use; see
  [credential lifecycle](../guides/credential-lifecycle.md).
- Do not rely on the legacy
  [encrypted-storage proposal](../ENCRYPTED_STORAGE_ARCHITECTURE.md). It is a historical design
  document, not a description of current behavior.

## Not covered

- Defense against malware or a compromised machine.
- Protection of a key file after it has been copied elsewhere.
- Any guarantee about third-party credential helpers configured next to gh-app-auth.
