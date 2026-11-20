---
title: Known limitations
description: Behaviors of gh-app-auth v0.2.0 that are limits of the current implementation, with the safe alternative for each.
sidebar:
  order: 3
---

These are limits of v0.2.0, not misconfigurations. Each lists what to do instead. Fixes will be marked
as unreleased until they ship.

## Automatic filesystem fallback for App keys

When no OS keyring is usable, `setup` may fall back to the filesystem and record
`private_key_source: filesystem` **without** recording a key path. Configuration validation then rejects
the entry. The managed fallback copy is raw bytes with mode `0400` under
`~/.config/gh/extensions/gh-app-auth/secrets/`, not encrypted.

**Do instead:** pass `--use-filesystem` together with `--key-file <path>`. The configuration keeps your
path, and nothing depends on the fallback.

## PATs need a keyring

Reading a personal access token from the filesystem store returns "not implemented". A machine without
an OS keyring cannot use a PAT.

**Do instead:** keep PATs on machines with a keyring, or use a GitHub App.

## GH_APP_PRIVATE_KEY and files

`GH_APP_PRIVATE_KEY` takes precedence over `--key-file`, and `setup --use-filesystem` rejects the
two together. Supplying a key through the environment does not make storage memory-only.

**Do instead:** in CI, supply the secret under another name and write it to a file you own.

## Long-running exec children

`exec` gives its child a token snapshot. Installation tokens expire after about an hour and nothing
renews them during the run.

**Do instead:** run `exec` per phase, or call it again before a long step.

## --config and --debug flags

The flags are registered on every command but are not wired as their help suggests.

**Do instead:** use `GH_APP_AUTH_CONFIG` for the configuration path and `GH_APP_AUTH_DEBUG_LOG` with an
absolute file path for logging.

## GH_APP_AUTH_CONFIG and secrets

The variable moves the configuration file but not the secrets directory, which stays under `HOME`.

**Do instead:** isolate `HOME` when you need isolated secrets, as the CI guides do.

## Key file permissions

POSIX checks reject group or world **read** access, and are skipped on Windows. They do not detect
symlinks, group write or later changes, and a key read from a stored configuration path is not
re-checked.

**Do instead:** use `chmod 600` on POSIX and set ACLs on Windows yourself.

## Composite action cleanup

`setup-gh-app-auth` registers its cleanup as a shell `trap` inside the setup step, not as a post-job hook.
The trap runs when that step's shell exits, so it can delete the credentials just configured before later
steps run. `cleanup-gh-app-auth` removes all Apps but not PATs.

**Do instead:** use the explicit steps in the [GitHub Actions guide](../guides/ci/github-actions.md).
Details are in the [GitHub Actions reference](../reference/github-actions.md).

## Equal priorities

When two credentials tie on priority, the winner depends on implementation details. Do not rely on
which one wins.

**Do instead:** make patterns unambiguous. See [routing](../understand/credential-routing.mdx).
