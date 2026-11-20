---
title: Continuous integration
description: Principles and a shared POSIX recipe for using gh-app-auth in disposable CI jobs, with provider-specific guides.
sidebar:
  order: 1
---

CI is where credential handling goes wrong most often: long-lived runners, shared home directories
and secrets interpolated into scripts. These guides keep one rule at their center: **run gh-app-auth
in a disposable, isolated job, and clean up explicitly.**

| Provider | Guide |
| --- | --- |
| GitHub Actions | [GitHub Actions](github-actions.md) (complete, copyable workflow) |
| Jenkins | [Jenkins](jenkins.md) (agent requirements and pipeline sketch) |
| GitLab CI | [GitLab CI](gitlab-ci.md) (runner requirements and job sketch) |

This project has no automated test against Jenkins or GitLab. Treat those two guides as grounded
requirements and untested sketches, and validate them in your own environment.

## Principles

- **Disposable and isolated.** Use an ephemeral runner, container or agent, with a job-owned `HOME`
  and `GH_CONFIG_DIR`. Never edit a shared home directory, and never run broad credential cleanup on a
  shared machine.
- **Secrets through the environment.** Bind the private key and the Client ID to environment
  variables with your CI system's secret mechanism. Do not interpolate secret values into script
  text, such as a Groovy string or a shell line assembled by the CI system.
- **Explicit key file.** Write the key to a file you own with `umask 077`, then use
  `setup --use-filesystem --key-file`. This is the reliable headless route. Automatic filesystem
  fallback is not: see [known limitations](../../troubleshoot/known-limitations.md).
- **Supply the key under a different name.** `GH_APP_PRIVATE_KEY` takes precedence over the key file inside
  `setup`, and `setup --use-filesystem` rejects it together with `--key-file`. Name the CI variable
  `APP_PRIVATE_KEY`.
- **Always clean up, and know what cleanup covers.** `gitconfig --clean` removes every
  gh-app-auth helper entry, and `remove` deletes the managed secrets for the App you name. Neither
  deletes the key file you supplied, so delete it yourself.
- **Pin the extension.** Install a specific release, for example
  `gh extension install AmadeusITGroup/gh-app-auth --pin v0.2.0`, so a pipeline does not change
  behavior on its own. Installing needs outbound access to GitHub and a GitHub CLI token.
- **Cleanup is best effort.** A cleanup step does not run if the machine disappears. An ephemeral
  machine is the real guarantee.
- **Tokens expire.** An installation token lasts up to an hour and the helper caches it in memory for 55
  minutes per process. Git starts a fresh helper for each request. `exec` takes a snapshot for its
  child and does not renew it.

## Shared POSIX job recipe

The same shell logic works under any CI system that gives the job a disposable POSIX shell. It expects
`APP_CLIENT_ID`, `APP_PATTERN` and `GH_TOKEN` in the environment, plus **exactly one** of
`APP_PRIVATE_KEY` (the PEM text) or `APP_PRIVATE_KEY_FILE` (the path of a file that holds the PEM, which is
how GitLab File variables arrive).

```sh
set -eu
umask 077
if [ -n "${APP_PRIVATE_KEY:-}" ] && [ -n "${APP_PRIVATE_KEY_FILE:-}" ]; then
  echo "set only one of APP_PRIVATE_KEY and APP_PRIVATE_KEY_FILE" >&2
  exit 1
fi
if [ -z "${APP_PRIVATE_KEY:-}" ] && [ -z "${APP_PRIVATE_KEY_FILE:-}" ]; then
  echo "set APP_PRIVATE_KEY or APP_PRIVATE_KEY_FILE" >&2
  exit 1
fi

job_dir="$(mktemp -d)"
export HOME="$job_dir/home"
export GH_CONFIG_DIR="$HOME/.config/gh"
mkdir -p "$GH_CONFIG_DIR"
key_file="$job_dir/app-key"

cleanup() {
  status=$?
  gh app-auth gitconfig --clean --global || status=1
  gh app-auth remove --client-id "$APP_CLIENT_ID" --force || status=1
  rm -rf -- "$job_dir" || status=1
  exit "$status"
}
trap cleanup EXIT

if [ -n "${APP_PRIVATE_KEY_FILE:-}" ]; then
  cat -- "$APP_PRIVATE_KEY_FILE" > "$key_file"
else
  printf '%s\n' "$APP_PRIVATE_KEY" > "$key_file"
  unset APP_PRIVATE_KEY
fi

gh extension install AmadeusITGroup/gh-app-auth --pin v0.2.0
gh app-auth setup \
  --use-filesystem \
  --client-id "$APP_CLIENT_ID" \
  --key-file "$key_file" \
  --patterns "$APP_PATTERN" \
  --name 'CI App'
gh app-auth gitconfig --sync --global
```

Run your Git commands after this block in the same shell. Because `HOME` is job-owned, the cleanup also
removes the files `setup` wrote.

## Not covered

Shared or long-lived self-hosted machines need stronger isolation than a temporary home directory, and
that setup is out of scope here. Windows agents are not covered.
