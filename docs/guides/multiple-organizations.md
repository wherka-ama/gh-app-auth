---
title: Several organizations and submodules
description: Route repositories in different organizations, and submodules, to different GitHub Apps.
sidebar:
  order: 5
---

A GitHub App installation belongs to one account. To reach repositories in several organizations,
register one App (or installation) per organization and give each its own pattern.

## One App per organization

```bash
gh app-auth setup --client-id Iv1.orga --key-file "$HOME/keys/org-a.pem" \
  --patterns 'github.com/org-a/' --name 'Org A'
gh app-auth setup --client-id Iv1.orgb --key-file "$HOME/keys/org-b.pem" \
  --patterns 'github.com/org-b/' --name 'Org B'
gh app-auth gitconfig --sync
```

Every value is a placeholder. After `--sync`, a clone of `github.com/org-a/...` uses `Org A` and a
clone of `github.com/org-b/...` uses `Org B`.

## Narrower patterns for exceptions

A longer pattern wins for App-to-App matching. Give one repository its own App without touching the
organization-wide entry:

```text
github.com/org-a/        -> Org A
github.com/org-a/secret  -> Secret App
```

Avoid two Apps with identical patterns, and avoid overlapping PAT patterns. Equal specificity makes
the outcome depend on implementation details. See
[how credential routing works](../understand/credential-routing.mdx).

## Submodules

Git asks for credentials separately for each submodule URL, so each is routed by its own URL:

```bash
git clone --recurse-submodules https://github.com/org-a/parent.git
```

If `parent` has a submodule in `org-b`, the helper selects `Org B` for that request. Nothing else
is needed, provided `--sync` ran after both Apps were registered.

## Verify each route

```bash
gh app-auth test --repo github.com/org-a/parent
gh app-auth test --repo github.com/org-b/child
```

`test` contacts GitHub with your real credentials, so run it deliberately. Use
`gh app-auth list` to review what is configured.

## Remember

Patterns choose which credential is sent. The installation's permissions and repository selection
limit what the token can do, so a route to an App that cannot see a repository still fails.
