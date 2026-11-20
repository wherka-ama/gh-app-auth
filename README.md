# gh-app-auth

A [GitHub CLI](https://cli.github.com/) extension and Git credential helper that routes each Git remote
to the right GitHub App or personal access token (PAT).

A GitHub App installation token belongs to one installation and expires after an hour. If you or
your pipelines touch repositories in several organizations, including submodules, you end up
juggling credentials. gh-app-auth registers each credential with the repository patterns it covers,
and Git asks it for credentials per request. This README describes v0.2.0.

## What it does

- Registers GitHub Apps (Client ID or App ID, private key) and PATs together with repository patterns.
- Acts as a Git credential helper: it picks the credential for the remote URL and, for an App, mints an
  installation token.
- Runs commands with a token in their environment: `gh app-auth exec`. For a GitHub App the token is a
  short-lived installation token; for a PAT it is the stored token.
- Keeps installation tokens in memory only (55-minute cache per process) and private keys in the OS
  keyring by default.
- Supports github.com, GitHub Enterprise Server, and PATs for hosts such as Bitbucket.

It does not run a daemon, does not renew tokens in the background, and does not limit what a token can do:
the GitHub App installation's permissions do that.

## Install

```bash
gh extension install AmadeusITGroup/gh-app-auth
```

Pin a release with `--pin v0.2.0`, and upgrade with `gh extension upgrade app-auth`.

## A safe starting path

This registers one GitHub App that is already installed on your organization with **Contents: read**
access, and clones a private repository. It needs the GitHub CLI and Git. All values are placeholders.
The `chmod` command is for Linux and macOS; on Windows, restrict the key file with ACLs, because gh-app-auth
skips its permission check there.

```bash
chmod 600 "$HOME/keys/example-app.pem"
gh app-auth setup \
  --client-id Iv1.example \
  --key-file "$HOME/keys/example-app.pem" \
  --patterns 'github.com/example-org/' \
  --name 'Example App'
gh app-auth gitconfig --sync
gh app-auth list --verify-keys
git clone https://github.com/example-org/private-repo.git
```

The key is stored in the OS keyring by default. On a headless machine without a keyring, add
`--use-filesystem` together with `--key-file`; do not rely on the automatic filesystem fallback. The
[first clone tutorial](docs/start-here/first-clone.md) explains each step and what to expect.

To call the GitHub API with the same credential, prefer `exec` over printing a token. Name the repository
explicitly, because `exec` otherwise uses the repository of your current directory:

```bash
gh app-auth exec --repo github.com/example-org/private-repo -- gh api repos/example-org/private-repo --jq .full_name
```

## Commands

| Command | Purpose |
| --- | --- |
| `setup` | Register a GitHub App or PAT with repository patterns |
| `gitconfig` | Add or remove the Git credential helper entries |
| `list`, `config`, `scope` | Inspect credentials, the configuration file and installation scope |
| `test` | Check that a credential works for a repository |
| `exec`, `token` | Run a command with a token, or print one explicitly |
| `migrate`, `remove` | Move key storage, or remove credentials |

Every command is documented in the [CLI reference](docs/reference/cli/index.md).

## Documentation

The documentation sources are in [`docs/`](docs/README.md). A documentation site is built from them
with Astro Starlight, but it is **not published yet**: it goes live after the documentation workflow deploys an
approved `main`. Until then, read the files on GitHub.

| If you want to | Read |
| --- | --- |
| Get started | [Overview](docs/start-here/overview.md), [installation](docs/start-here/installation.md), [first clone](docs/start-here/first-clone.md) |
| Do a task | [Guides](docs/guides/index.md), including [CI](docs/guides/ci/index.md) |
| Look something up | [CLI reference](docs/reference/cli/index.md), [configuration](docs/reference/configuration.md), [environment variables](docs/reference/environment-variables.md) |
| Understand it | [Credential routing](docs/understand/credential-routing.mdx), [security model](docs/understand/security-model.md), [architecture](docs/understand/architecture.mdx) |
| Fix a problem | [Troubleshooting](docs/troubleshoot/index.md), [known limitations](docs/troubleshoot/known-limitations.md) |
| Contribute | [Contributing](CONTRIBUTING.md), [development](docs/contribute/development.md), [release runbook](docs/contribute/releasing.md) |

## Known limitations

- The automatic filesystem fallback for App keys can leave a configuration that is rejected later.
  Use `--use-filesystem --key-file`.
- A PAT cannot be read from the filesystem store, so PATs need an OS keyring.
- `exec` gives its child a token snapshot that is not renewed during the run.
- The `--config` and `--debug` flags are not wired as their help suggests; use `GH_APP_AUTH_CONFIG`
  and `GH_APP_AUTH_DEBUG_LOG` (with an absolute file path).

Details and alternatives are in [known limitations](docs/troubleshoot/known-limitations.md).

## Contributing and security

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and the pull request process, and the
[documentation guide](docs/contribute/documentation.md) for working on the docs. Report vulnerabilities
privately as described in the [security policy](.github/SECURITY.md).

## License

MIT License. See [LICENSE](LICENSE) for details.
