# gh-app-auth

A [GitHub CLI](https://cli.github.com/) extension that provides Git credential
authentication using GitHub Apps and Personal Access Tokens (PATs).

It implements the [Git credential helper](https://git-scm.com/docs/gitcredentials)
protocol, so `git clone`, `git fetch`, and `git push` work without additional
tooling once configured.

## Installation

```bash
gh extension install AmadeusITGroup/gh-app-auth
```

Requires GitHub CLI v2.45+.

## Quick start

```bash
# 1. Configure a GitHub App
export GH_APP_PRIVATE_KEY="$(cat ~/keys/my-app.pem)"
gh app-auth setup \
  --app-id 123456 \
  --patterns "github.com/myorg/"
unset GH_APP_PRIVATE_KEY

# 2. Set up git credential helpers
gh app-auth gitconfig --sync --global

# 3. Use git as usual
git clone https://github.com/myorg/private-repo.git
```

To use a Personal Access Token instead of an App:

```bash
gh app-auth setup \
  --pat ghp_your_token \
  --patterns "github.com/myorg/" \
  --name "My PAT"
gh app-auth gitconfig --sync --global
```

## How it works

The extension registers itself as a Git credential helper. When Git needs
credentials for a URL, the helper compares the URL against configured prefixes
(e.g. `github.com/myorg/`) and picks the credential with the longest match.
For GitHub Apps this means generating a JWT and exchanging it for a short-lived
installation token. For PATs the stored token is returned directly.

By default GitHub grants 60-minute validity to installation tokens and that's what we use here. However, every interaction with gh-app-auth leads to creation of a new token. They are not reused between consequtive calls, unless it is driven by the some internal git flows. 

## Commands

| Command | Description |
|---------|-------------|
| `setup` | Add a GitHub App or PAT |
| `list` | Show configured credentials |
| `remove` | Delete a credential |
| `test` | Verify authentication works |
| `scope` | Show which repos an App can access |
| `config` | Print config file path or contents |
| `gitconfig` | Manage git credential helper entries |
| `migrate` | Move keys to encrypted storage |
| `git-credential` | Git credential helper (called by Git) |
| `help` | Show help for any command |
| `debug` | Show information about the scope of the configured credentials i.e. app installations and repositories (useful for troubleshooting) |

Run `gh app-auth <command> --help` for flags and usage.

## Documentation

- [Installation and setup](docs-improved/installation.md)
- [Configuration reference](docs-improved/configuration.md)
- [CI/CD integration](docs-improved/ci-cd.md)
- [Security](docs-improved/security.md)
- [Troubleshooting](docs-improved/troubleshooting.md)

For contributors:

- [Architecture](docs-improved/architecture.md)
- [Development guide](docs-improved/development.md)
- [Contributing](CONTRIBUTING.md)