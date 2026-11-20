---
title: GitHub Enterprise Server
description: Configure gh-app-auth for GitHub Enterprise Server hosts alongside github.com.
sidebar:
  order: 6
---

For hosts other than `github.com`, gh-app-auth talks to the GitHub Enterprise Server REST API under
`https://<host>/api/v3`. For `github.com` it uses `api.github.com`.

## Register an App for your server

Create and install the GitHub App on the server, as described in
[create and install a GitHub App](create-github-app.md), then include the host in the pattern:

```bash
gh app-auth setup \
  --client-id Iv1.example \
  --key-file "$HOME/keys/ghes-app.pem" \
  --patterns 'ghes.example.com/example-org/' \
  --name 'GHES App'
gh app-auth gitconfig --sync
```

Every value is a placeholder.

## Mix with github.com

Register separate entries for each host. Patterns start with the host, so
`github.com/example-org/` and `ghes.example.com/example-org/` never collide:

```bash
gh app-auth list
```

## Using exec against a server

`gh app-auth exec` injects the token as `GH_ENTERPRISE_TOKEN` for hosts other than github.com and
sets `GH_HOST`, so `gh` talks to the right server:

```bash
gh app-auth exec --repo ghes.example.com/example-org/private-repo -- gh api repos/{owner}/{repo}
```

## Checks

- The server must be reachable over HTTPS from the machine running Git.
- `gh app-auth test --repo ghes.example.com/example-org/private-repo` exercises the full path.
- Installation scope metadata is cached for 24 hours; refresh it with `gh app-auth scope --refresh`.
