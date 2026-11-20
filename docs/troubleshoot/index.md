---
title: Troubleshooting
description: Find your symptom, check the cause, and apply the fix for the most common gh-app-auth problems.
sidebar:
  order: 1
---

Start with the symptom. Each entry names the check that confirms the cause before the fix. For
logging, see [diagnostics](diagnostics.md); for behaviors that are limits rather than bugs, see
[known limitations](known-limitations.md).

## First checks

```bash
gh app-auth list --verify-keys
gh app-auth config --path
git config --show-origin --get-regexp '^credential\.'
```

These show which credentials exist and whether their keys are readable, which configuration file is
in use, and which Git credential helpers are configured.

## Symptoms

| Symptom | Likely cause | Check | Fix |
| --- | --- | --- | --- |
| Git still prompts for a username and password | The helper is not configured, or another helper answers first | `git config --show-origin --get-all credential.helper` | Run `gh app-auth gitconfig --sync`; reorder or remove the other helper |
| The wrong credential is used | Overlapping patterns, or a PAT priority higher than the App | `gh app-auth list` | Make patterns unambiguous; see [routing](../understand/credential-routing.mdx) |
| `setup` reports `Both --key-file and GH_APP_PRIVATE_KEY` | Both were provided | Environment and flags | Provide one. In CI use a file and a differently named variable |
| A configuration is rejected after `setup` with no keyring | Automatic filesystem fallback recorded no key path | `gh app-auth config --show` | Re-run `setup --use-filesystem --key-file <path>` |
| `list --verify-keys` cannot read a key | Missing keyring entry or key file | The `private_key_source` and path | Re-run `setup`, or `migrate` |
| Key file rejected | Group or world readable on POSIX | `ls -l` on the file | `chmod 600` the file |
| `test` fails at the installation step | The App is not installed on that owner, or the repository is outside its selection | The App's installation settings on GitHub | Install the App or add the repository |
| `test` fails with a clock-skew message | The machine clock is off | `date -u` against a trusted source | Fix the clock; see the [ADR](../adr/0001-jwt-clock-skew-margin.md) |
| A long `exec` command fails after about an hour | The installation token expired | Runtime of the child | Split the work, or run `exec` again for each phase |
| A PAT does not work on a headless machine | Filesystem PAT retrieval is not implemented | Whether an OS keyring exists | Use a machine with a keyring |
| `--json` is not recognized | The flag does not exist | `gh app-auth list --help` | Use `list --format json` |
| `scope --repo` is not recognized | `scope` has only `--refresh` | `gh app-auth scope --help` | Use `scope --refresh` |
| `--config` or `--debug` have no effect | The flags are not wired in v0.2.0 | n/a | Use `GH_APP_AUTH_CONFIG` and `GH_APP_AUTH_DEBUG_LOG` |

## Non-GitHub hosts

`test` is written for GitHub hosts and does not guarantee anything about other providers. For a
Bitbucket-style host, confirm the PAT works with the provider directly, make sure `--username` was set
when the provider needs a real account name, and end the PAT pattern with `/`.

## Git helper checks

- The helper path stored in Git's configuration is the extension binary. After you move or reinstall the
  extension, run `gh app-auth gitconfig --sync` again.
- Git uses the first helper in configuration order that answers. Path-specific entries written by
  `gitconfig --sync` come before generic host helpers.
- There is no cache-removal command for installation tokens, because they live only in process
  memory and disappear when the process exits.

## Recovery

> **Destructive scope.** `gitconfig --clean` removes the helper entries for all credentials, and
> `remove --all` removes all configured Apps and their managed secrets. Prefer removing one credential
> at a time. Outside an isolated, disposable CI job, run `remove --all` interactively without `--force`
> and read the prompt.

```bash
gh app-auth gitconfig --clean
gh app-auth remove --all
```

Then register the credentials again with `setup` and run `gitconfig --sync`. `remove --all` does not
remove PATs; run `remove --all-pats` or `remove --pat-name` separately, because target flags cannot be
combined. See [credential lifecycle](../guides/credential-lifecycle.md).

## Still stuck

Capture a [diagnostic log](diagnostics.md), review it for sensitive values and open an issue on
[GitHub](https://github.com/AmadeusITGroup/gh-app-auth/issues). For vulnerabilities, follow the
[security policy](../../.github/SECURITY.md) instead.
