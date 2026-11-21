---
title: Install gh-app-auth
description: Prerequisites and installation options for the gh-app-auth extension, release packages and source builds.
sidebar:
  order: 2
---

## Prerequisites

- The [GitHub CLI](https://cli.github.com/) (`gh`), version 2.45 or newer, and Git 2.30 or newer.
  These are the minimum versions the project has documented.
- A GitHub App or a personal access token to register. See
  [create and install a GitHub App](../guides/create-github-app.md).
- Go is **not** required to use the extension. It is only needed to build from source or to
  contribute. Source builders need the Go version specified in [go.mod](../../go.mod).

## Install the extension

```bash
gh extension install AmadeusITGroup/gh-app-auth
```

`gh` installs the precompiled binary for your platform from the latest GitHub release.

To stay on a specific release, pin it. This example pins `v0.2.0`; choose and test the release you want:

```bash
gh extension install AmadeusITGroup/gh-app-auth --pin v0.2.0
```

Upgrade an unpinned install with `gh extension upgrade app-auth`. Check what is installed with
`gh extension list`; the output of `gh app-auth --version` is not a reliable indicator of the
release.

## Release packages

Each release also attaches Linux DEB and RPM packages and cross-platform binaries. DEB packages
install directly with `dpkg -i`. RPM packages require importing the release signing key first. The
[release runbook](../contribute/releasing.md) lists every asset and how packages are signed.

## Build from source

```bash
git clone https://github.com/AmadeusITGroup/gh-app-auth.git
cd gh-app-auth
go build -o gh-app-auth .
gh extension install .
```

See [development](../contribute/development.md) for the full contributor setup.

## Shell completion

`gh app-auth completion` generates completion scripts for bash, fish, PowerShell and zsh. See the
[`completion` reference](../reference/cli/completion.md).

## Next step

[Make your first authenticated clone](first-clone.md).
