# Installation and setup

## Prerequisites

- [GitHub CLI](https://cli.github.com/) v2.45+ (`gh auth status` to verify)
- Git 2.30+
- A GitHub App **or** a Personal Access Token

## Install the extension

```bash
gh extension install AmadeusITGroup/gh-app-auth
```

Upgrade later with:

```bash
gh extension upgrade app-auth
```

## Configure credentials

### GitHub App

You need the App ID and a private key file. Both are available from the App's
settings page in your GitHub organization.

Store the key via environment variable (recommended):

```bash
export GH_APP_PRIVATE_KEY="$(cat ~/keys/my-app.pem)"
gh app-auth setup \
  --app-id 123456 \
  --patterns "github.com/myorg/" \
  --name "Org Automation"
unset GH_APP_PRIVATE_KEY
```

Or reference a key file on disk:

```bash
gh app-auth setup \
  --app-id 123456 \
  --key-file ~/keys/my-app.pem \
  --patterns "github.com/myorg/"
```

### Personal Access Token

```bash
gh app-auth setup \
  --pat ghp_your_token \
  --patterns "github.com/myorg/" \
  --name "My PAT" \
  --priority 15
```

For providers that require a real username (e.g., Bitbucket Server):

```bash
gh app-auth setup \
  --pat bbpat_your_token \
  --patterns "bitbucket.example.com/" \
  --username your_username \
  --name "Bitbucket"
```

### Multiple credentials

Run `setup` once per credential. The longest matching prefix wins.
See [configuration.md](configuration.md) for details.

## Set up git credential helpers

```bash
gh app-auth gitconfig --sync --global
```

This reads your configured patterns and writes the corresponding
`credential.helper` entries into `~/.gitconfig`. Run it again after adding or
removing credentials.

Options:

- `--local` — scope to the current repository instead of global
- `--clean` — remove all gh-app-auth credential helper entries
- `--auto` — CI/CD mode using `GH_APP_ID` and `GH_APP_PRIVATE_KEY_PATH` env vars

## Verify

```bash
# List credentials and check key accessibility
gh app-auth list --verify-keys

# Test authentication against a specific repo
gh app-auth test --repo github.com/myorg/private-repo
```

## Next steps

- [Configuration reference](configuration.md) — config file format, pattern
  matching, multi-org examples
- [CI/CD integration](ci-cd.md) — GitHub Actions and Jenkins examples
- [Troubleshooting](troubleshooting.md) — common issues and debug logging
