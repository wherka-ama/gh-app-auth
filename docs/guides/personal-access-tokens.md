---
title: Personal access tokens
description: Register a personal access token with gh-app-auth, including a username for non-GitHub hosts, and understand the limits.
sidebar:
  order: 3
---

Use a personal access token (PAT) when a host cannot issue GitHub App tokens: a Bitbucket server, or
a case where automation must act as a person. Prefer a GitHub App whenever one works, because its
tokens are short-lived and scoped by installation.

## Register a PAT

```bash
read -rs PAT_VALUE
gh app-auth setup \
  --pat "$PAT_VALUE" \
  --patterns 'bitbucket.example.com/' \
  --username example-user \
  --name 'Bitbucket PAT'
unset PAT_VALUE
```

`read -rs` keeps the token out of your shell history and off the screen. It does not hide it from
other users of the machine: `--pat` is a command-line argument, so it can appear in the process list
while `setup` runs. Do this only on a machine you trust, and revoke the token if you suspect exposure.

For github.com or GitHub Enterprise hosts, omit `--username`: the default is `x-access-token`. Hosts
that need a real account name, such as Bitbucket Data Center, need `--username`.

## Storage and its limits

- A PAT is stored in the OS keyring. Reading a PAT from the filesystem store is **not
  implemented**, so a machine without a keyring cannot use a PAT. Do not plan a headless PAT setup.
- A PAT is a static credential. How long it lives is controlled by the provider and organization
  policy, its expiry date, and revocation, not only by whoever created it. Rotate it on your own
  schedule too.

## How a PAT competes with an App

A PAT pattern matches by plain string prefix after removing `https://`. Unlike App patterns, it does
**not** respect path segment boundaries, so `github.com/acme` also matches `github.com/acme-labs/tool`.
End PAT patterns with `/`.

When an App and a PAT both match, the helper compares their priorities. `--priority` defaults to
`5` for both, an App is considered first, and a PAT must be strictly higher to win. Equal priorities
therefore keep the App in the current implementation; do not build on tie behavior. Details are in
[how credential routing works](../understand/credential-routing.mdx).

## Use it

```bash
gh app-auth gitconfig --sync
git ls-remote https://bitbucket.example.com/project/repo.git HEAD
```

`git ls-remote` makes Git ask the helper for credentials and asks the provider to authenticate, without
cloning. It succeeds only if the token has the permissions the provider requires for that repository, so
check the token's scopes with the provider if it is rejected. Reserve `gh app-auth test` for GitHub
hosts: it calls the GitHub API and says nothing reliable about other providers. Clone with plain
`git clone` afterwards.

## Remove it

```bash
gh app-auth remove --pat-name 'Bitbucket PAT'
```

See [credential lifecycle](credential-lifecycle.md).
