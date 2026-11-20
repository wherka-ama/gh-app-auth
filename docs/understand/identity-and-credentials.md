---
title: Identity and credentials
description: How GitHub Apps, installations and personal access tokens differ, and what gh-app-auth stores for each.
sidebar:
  order: 1
---

gh-app-auth moves two kinds of credential around. Knowing which is which explains most of its
behavior.

## GitHub App credentials

A GitHub App is an identity owned by a user or organization. Four pieces work together:

| Piece | What it is | Secret? |
| --- | --- | --- |
| **App identity** | The Client ID (preferred) or numeric App ID | No |
| **Private key** | A PEM key used to sign JWTs as the App | Yes |
| **Installation** | The App installed on one account, with a repository selection and permissions | No |
| **Installation token** | A short-lived token for one installation | Yes |

gh-app-auth signs a JWT with the private key, calls GitHub to exchange it for an installation
token, and hands that token to Git. Installation tokens expire after an hour. The token's reach is
set by the installation's permissions and repository selection, which gh-app-auth cannot widen
or narrow.

Ownership and installation are different things. A GitHub App is **owned** by one account (a user or
an organization), and it can be **installed** on that account and on many others. Each installation is
specific to one account, so reaching repositories in several organizations means one installation per
organization, which can be installations of the same App or of different Apps. See
[several organizations](../guides/multiple-organizations.md).

## Personal access tokens

A personal access token (PAT) is a static credential issued to a person or service account by a
provider. gh-app-auth stores it and returns it unchanged. Its lifetime is decided by provider and
organization policy, its expiry date and revocation. PATs are the only option for providers that
cannot issue GitHub App tokens, such as Bitbucket, and they carry the permissions of the account that
created them. See [personal access tokens](../guides/personal-access-tokens.md).

## What gh-app-auth stores

| Item | Stored in |
| --- | --- |
| Patterns, names, priorities, installation IDs | The configuration file |
| App private key | The OS keyring by default, or a file you provide with `--use-filesystem --key-file` |
| PAT value | The OS keyring (filesystem retrieval is not implemented) |
| Installation token | Memory only, per process |
| Installation scope metadata | The configuration file, refreshed after 24 hours |

Nothing in this table is encrypted by gh-app-auth itself. The OS keyring encrypts at rest; files are
plain. See the [security model](security-model.md).

## Identifiers in commands

`setup`, `exec`, `token` and `remove` accept `--client-id` and, for legacy entries, `--app-id`.
`remove` also accepts `--pat-name`, `--all`, `--all-pats` and `--force`. The `--installation-id` selector
exists only on `setup`, `exec` and `token`, where it disambiguates one of an App's installations.
`remove` has no installation selector: when an App has several installations, read the confirmation
prompt carefully before agreeing, and do not use `--force` for an ambiguous removal. Client IDs appear on
the App's settings page, for example `Iv1.…`.
