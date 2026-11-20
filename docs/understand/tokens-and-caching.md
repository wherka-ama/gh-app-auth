---
title: Tokens and caching
description: The lifetimes of JWTs, installation tokens, the in-memory cache and scope metadata, and what each does and does not do.
sidebar:
  order: 3
---

Several clocks run at once in gh-app-auth. They are separate and should not be confused.

| Item | Lifetime | Where it lives |
| --- | --- | --- |
| App JWT | `exp` is 9 minutes ahead and `iat` is 60 seconds behind the current time | Created on demand, not cached |
| Installation token | Up to 1 hour, set by GitHub | Memory only |
| Installation token cache entry | 55 minutes, per process | Memory only |
| Installation scope metadata | 24 hours | The configuration file |
| PAT | Provider policy, expiry or revocation | OS keyring |

## JWT

To act as the App, gh-app-auth signs a JWT with the App's private key. The `iat` claim is set 60
seconds in the past to tolerate clock skew, and `exp` is 9 minutes in the future. A JWT is used to
request an installation token and is not stored. See the
[ADR on the clock-skew margin](../adr/0001-jwt-clock-skew-margin.md).

## Installation token

GitHub issues an installation token for one installation, valid for up to an hour. The helper
caches it in memory for 55 minutes, which leaves a margin before expiry.

- The cache exists only inside one process. Git starts a **fresh helper process for each request**,
  so each request normally mints a new token. The cache helps processes that stay alive and make
  several requests, such as `gh app-auth test`.
- A background worker removes expired entries. It does **not** renew them. There is no daemon and no disk
  cache.
- Tokens are never written to disk.

## exec and token

`exec` takes a snapshot of a token and gives it to its child. Nothing renews it, so a child that runs
longer than the token's life starts failing. `token` bypasses the cache and mints a fresh token each
time; its output goes to standard output, which is a secret boundary. See the
[security model](security-model.md).

## Scope metadata

`gh app-auth scope` records which repositories an installation can see and caches that for 24
hours. It is metadata used for matching decisions; it does not extend or shorten any token's
life, and `scope --refresh` re-reads it.

## What caching does not change

- A cached token is still bound to its installation's permissions.
- The local cache has no independent revocation and receives no push notification. GitHub authorizes
  each request, and it may reject a token before the cache's 55 minutes are up. Suspending or uninstalling
  the installation is how access is blocked; see GitHub's
  [guide to reviewing and modifying installed GitHub Apps](https://docs.github.com/en/apps/using-github-apps/reviewing-and-modifying-installed-github-apps)
  and the [REST reference for revoking an installation access token](https://docs.github.com/en/rest/apps/installations#revoke-an-installation-access-token).
- Restarting a process starts with an empty cache.
