---
title: Configuration reference
description: The gh-app-auth configuration file, its fields, validation rules and storage locations.
sidebar:
  order: 2
---

gh-app-auth stores which credentials exist and which repositories they cover in one configuration
file. You normally change it with `setup`, `remove` and `migrate`, but it is safe to read and, with
care, to edit.

## Location

| What | Where |
| --- | --- |
| Configuration file | `~/.config/gh/extensions/gh-app-auth/config.yml`, or the file named by `GH_APP_AUTH_CONFIG` |
| Stored secrets | `~/.config/gh/extensions/gh-app-auth/secrets/`, always under `HOME` |
| Diagnostic log | The file named by `GH_APP_AUTH_DEBUG_LOG`, when set |

`gh app-auth config` prints the path (`--path`) or the content (`--show`). A `.yaml` or `.yml`
extension is parsed as YAML and `.json` as JSON. Other extensions are tried as YAML, then JSON.
`GH_APP_AUTH_CONFIG` moves only the configuration file; it does not relocate stored secrets.

## Examples

A GitHub App whose key is in the OS keyring:

```yaml
version: "1.0"
github_apps:
  - name: Example App
    client_id: Iv1.example
    installation_id: 12345678
    private_key_source: keyring
    patterns:
      - github.com/example-org/
```

A GitHub App whose key is a file you manage, the reliable headless layout:

```yaml
version: "1.0"
github_apps:
  - name: CI App
    client_id: Iv1.example
    private_key_source: filesystem
    private_key_path: /home/runner/work/_temp/app-key
    patterns:
      - github.com/example-org/
```

A personal access token for a host that needs a username:

```yaml
version: "1.0"
pats:
  - name: Bitbucket PAT
    private_key_source: keyring
    username: example-user
    patterns:
      - bitbucket.example.com/
    priority: 5
```

An App and a PAT together:

```yaml
version: "1.0"
github_apps:
  - name: Org App
    app_id: 123456
    private_key_source: keyring
    patterns:
      - github.com/example-org/
    priority: 5
pats:
  - name: Legacy PAT
    private_key_source: keyring
    patterns:
      - github.com/legacy-org/
    priority: 10
```

The token itself is never in this file. A PAT's value is in the keyring, and an App's key is in the
keyring or at `private_key_path`.

## Top-level fields

| Field | Meaning |
| --- | --- |
| `version` | A non-empty string. New files are created with `"1.0"`. Validation only checks that it is not empty; it does not reject other strings |
| `github_apps` | List of GitHub App entries |
| `pats` | List of personal access token entries. Optional |

At least one App or PAT is required.

## GitHub App fields

| Field | Required | Meaning |
| --- | --- | --- |
| `name` | Yes | Friendly name shown in `list` |
| `client_id` | One of | The App's Client ID. Preferred |
| `app_id` | One of | The numeric App ID. Legacy; used when no Client ID is set |
| `installation_id` | No | Installation ID. `0` or absent means it is detected at runtime |
| `private_key_source` | See below | Where the key lives: `keyring` or `filesystem` |
| `private_key_path` | With `filesystem` | Path to the key file. `~` is expanded |
| `patterns` | Yes | At least one non-empty pattern |
| `priority` | No | Compared with PATs only; see [routing](../understand/credential-routing.mdx) |
| `scope` | No | Cached installation scope, written by gh-app-auth |

Validation rules:

- Without `private_key_source`, a `private_key_path` is required and the source is treated as
  `filesystem`.
- `filesystem` requires `private_key_path`. The automatic filesystem fallback in `setup` can record
  `filesystem` without a path, which validation then rejects. See
  [known limitations](../troubleshoot/known-limitations.md).
- `inline` is rejected: inline keys must be migrated to the keyring or the filesystem.
- `priority` does not affect which App the URL matcher picks. It is used only when an App and a PAT
  both match, and it is not deprecated for that purpose.

The `scope` block holds `repository_selection`, `account_login`, `account_type`, an optional
`repositories` list and the cache timestamps `last_fetched`, `last_updated` and `cache_expiry`. It is
refreshed with `gh app-auth scope --refresh` and cached for 24 hours.

## Personal access token fields

| Field | Required | Meaning |
| --- | --- | --- |
| `name` | Yes | Friendly name |
| `patterns` | Yes | At least one non-empty pattern, matched by plain string prefix |
| `private_key_source` | No | `keyring` or `filesystem`. The field keeps this name for PATs; `token_source` is not a valid field |
| `priority` | No | Compared with Apps when both match. `setup` defaults it to `5` |
| `username` | No | Username for HTTP basic authentication. Defaults to `x-access-token` for GitHub |

PAT retrieval from the filesystem store is not implemented in v0.2.0, so use `keyring` for PATs.

## What is and is not validated

The loader parses the file with ordinary YAML or JSON decoding, so it does **not** reject unknown
fields at runtime. A misspelled field is silently ignored. The examples on this page are checked more
strictly: the documentation build rejects unknown fields in them, which catches drift in the docs but
is not a runtime guarantee. Every `yaml` block on this page is a complete configuration with exactly one
document; partial snippets belong in `text` blocks. For example, this PAT entry is rejected by that check because
`token_source` is not a field:

```yaml invalid
version: "1.0"
pats:
  - name: Wrong field name
    token_source: keyring
    patterns:
      - github.com/example-org/
```

## Editing safely

- Keep a copy before you edit. Run `gh app-auth list --verify-keys` afterwards.
- Prefer `setup`, `remove` and `migrate` over manual edits so secrets and configuration stay in step.
- Re-run `gh app-auth gitconfig --sync` after you change patterns, so Git's helper entries match.
