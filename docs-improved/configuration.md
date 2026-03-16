# Configuration reference

## Config file location

1. `$GH_APP_AUTH_CONFIG` environment variable (if set)
2. `~/.config/gh/extensions/gh-app-auth/config.yml` (default)

```bash
gh app-auth config          # show path and status
gh app-auth config --path   # print path only
gh app-auth config --show   # print file contents
```

## File format

```yaml
version: "1"
github_apps:
  - name: Org Automation
    app_id: 123456
    installation_id: 987654        # optional, auto-detected during setup
    private_key_source: keyring    # keyring | filesystem
    private_key_path: ""           # set when source=filesystem
    patterns:
      - github.com/myorg/
    priority: 5

pats:
  - name: Personal
    token_source: keyring
    patterns:
      - github.com/personal-org/
    priority: 15
    username: ""                   # optional, defaults to x-access-token
```

At least one `github_apps` or `pats` entry is required.

## GitHub App fields

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `name` | string | yes | Display label |
| `app_id` | int | yes | GitHub App ID |
| `installation_id` | int | no | Auto-detected if omitted |
| `private_key_source` | string | yes | `keyring` or `filesystem` |
| `private_key_path` | string | no | Only when `source=filesystem` |
| `patterns` | list | yes | Prefixes: `host/`, `host/org/`, or `host/org/repo` |
| `priority` | int | no | Deprecated for App-vs-App (longest prefix wins). Only used when choosing between an App and a PAT. |

## PAT fields

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `name` | string | yes | Display label and keyring key |
| `token_source` | string | yes | `keyring` or `filesystem` |
| `patterns` | list | yes | Prefixes: `host/`, `host/org/`, or `host/org/repo` |
| `priority` | int | yes | Used when an App and a PAT both match the same URL |
| `username` | string | no | Defaults to `x-access-token`; set for Bitbucket |

### Username by provider

| Provider | Username |
|----------|----------|
| GitHub.com / GitHub Enterprise | leave blank (`x-access-token`) |
| Bitbucket Server / Data Center | your Bitbucket username |
| Other HTTPS Git hosts | whatever the provider expects |

## Prefix matching

Despite the field name `patterns`, values are plain string prefixes — there is
no glob or wildcard support. The `/*` suffix is stripped for backward
compatibility with older configs.

When Git requests credentials for a URL:

1. The protocol is stripped (`https://github.com/org/repo` → `github.com/org/repo`).
2. Each configured prefix is compared with `strings.HasPrefix`.
3. The credential with the longest matching prefix wins.
4. If an App and a PAT both match, the one with the higher `priority` wins.
   Priority between two Apps is ignored (longest prefix always wins).

Examples:

| URL | Matching prefixes | Winner |
|-----|-------------------|--------|
| `github.com/org/repo` | App `github.com/org/`, PAT `github.com/` | App (longer prefix) |
| `github.com/org/repo` | App `github.com/org/` (pri 5), PAT `github.com/org/` (pri 15) | PAT (higher priority) |
| `bitbucket.example.com/scm/team/repo` | PAT `bitbucket.example.com/` | PAT (only match) |
| `github.com/org/special-repo` | App A `github.com/org/special-repo`, App B `github.com/org/` | App A (longer prefix) |

### Git credential context mapping

`gitconfig --sync` converts prefixes to git credential contexts:

| Prefix | Git context |
|--------|-------------|
| `github.com/org/` | `https://github.com/org` |
| `github.enterprise.com/` | `https://github.enterprise.com` |
| `bitbucket.example.com/` | `https://bitbucket.example.com` |

## Secret storage

Private keys and PATs are stored separately from the config file:

- **Keyring** (default) — macOS Keychain, Windows Credential Manager, Linux
  Secret Service (GNOME Keyring / KWallet)
- **Filesystem** (fallback) — `~/.config/gh/extensions/gh-app-auth/secrets/`,
  used only when the keyring is unavailable

Removing a credential with `gh app-auth remove` deletes its secret from both
locations.

## Multi-organization examples

### Two GitHub orgs

```yaml
version: "1"
github_apps:
  - name: Frontend
    app_id: 111111
    private_key_source: keyring
    patterns:
      - github.com/frontend-org/

  - name: Backend
    app_id: 222222
    private_key_source: keyring
    patterns:
      - github.com/backend-org/
```

### GitHub Enterprise + GitHub.com

```yaml
version: "1"
github_apps:
  - name: Enterprise
    app_id: 333333
    private_key_source: keyring
    patterns:
      - github.enterprise.com/

  - name: Cloud
    app_id: 444444
    private_key_source: keyring
    patterns:
      - github.com/myorg/
```

### GitHub App + Bitbucket PAT

```yaml
version: "1"
github_apps:
  - name: GitHub App
    app_id: 555555
    private_key_source: keyring
    patterns:
      - github.com/myorg/

pats:
  - name: Bitbucket
    token_source: keyring
    patterns:
      - bitbucket.example.com/
    priority: 40
    username: jsmith
```

## Manual git configuration

If you prefer to set credential helpers by hand instead of using
`gitconfig --sync`:

```bash
git config --global credential.'https://github.com/myorg'.helper \
  '!gh app-auth git-credential --pattern "github.com/myorg/"'
```

Use `gh app-auth scope --repo <url>` to check which credential would be
selected for a given URL.
