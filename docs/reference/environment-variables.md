---
title: Environment variables
description: Every environment variable gh-app-auth reads, with precedence rules, effects and cautions.
sidebar:
  order: 3
---

This page lists the variables read by the gh-app-auth code.

| Variable | Read by | Effect |
| --- | --- | --- |
| `GH_APP_AUTH_CONFIG` | All commands and the helper | Path to the configuration file |
| `GH_APP_AUTH_DEBUG_LOG` | All commands and the helper | Enables diagnostic logging to the file at this path |
| `GH_APP_PRIVATE_KEY` | `setup` | PEM content of an App private key |
| `GH_APP_PRIVATE_KEY_PATH` | Auto mode | Path to the private key file |
| `GH_APP_ID` | Auto mode | Numeric App ID |
| `GH_APP_CLIENT_ID` | Auto mode | App Client ID |
| `GH_REPO` | `test` | Overrides the current repository |
| `HOME` | Default locations | Base of the default configuration and secrets directories |

## GH_APP_AUTH_CONFIG

Changes where the configuration file is read from. The default is
`~/.config/gh/extensions/gh-app-auth/config.yml`. A `.yaml` or `.yml` extension is read as YAML and a
`.json` extension as JSON; any other extension is tried as YAML, then JSON.

It does **not** move the secrets directory. Secrets stored by gh-app-auth stay under a fixed
directory beneath `HOME`: `~/.config/gh/extensions/gh-app-auth/secrets/`.

## GH_APP_AUTH_DEBUG_LOG

Set it to an explicit absolute path. Any non-empty value enables logging, and the value is used
as the file name, so `GH_APP_AUTH_DEBUG_LOG=1` writes a file named `1` in the current directory.

```bash
export GH_APP_AUTH_DEBUG_LOG="$HOME/gh-app-auth-debug.log"
```

The file is appended to with mode `0600`, and missing parent directories are created with mode
`0700`. If the directory cannot be created the log goes to a file in the system temporary directory,
and if the file cannot be opened logging is silently disabled. See
[diagnostics](../troubleshoot/diagnostics.md) for what is logged.

## GH_APP_PRIVATE_KEY

Supplies the PEM key to `setup` without a file. It takes precedence over `--key-file`. Using it does
**not** make storage memory-only: the key is still persisted by gh-app-auth. `setup --use-filesystem`
rejects it together with `--key-file` as conflicting options. In CI, supply the secret under another
name and write it to a file instead; see the [CI guide](../guides/ci/index.md).

## Auto mode: GH_APP_PRIVATE_KEY_PATH, GH_APP_ID, GH_APP_CLIENT_ID

If `GH_APP_PRIVATE_KEY_PATH` is set together with `GH_APP_ID` or `GH_APP_CLIENT_ID`, the credential
helper can configure an App on demand when no configured credential matches the requested
repository. It records an entry named `Auto setup` using the requested repository URL as its
pattern, with priority `5`. `gh app-auth gitconfig --sync --auto` prepares Git for this mode. Auto
mode writes configuration and key storage by the same rules as `setup`, so the
[known limitations](../troubleshoot/known-limitations.md) apply.

## Variables that gh-app-auth sets or removes

`gh app-auth exec` removes `GH_TOKEN`, `GITHUB_TOKEN`, `GH_ENTERPRISE_TOKEN`,
`GITHUB_ENTERPRISE_TOKEN`, `GH_HOST` and `GH_REPO` from the child environment, then sets the token
variable for the target host (`GH_TOKEN` for github.com, `GH_ENTERPRISE_TOKEN` otherwise), `GH_HOST`
and, when a repository is known, `GH_REPO`. See [`exec`](cli/exec.md).
