---
title: Make your first authenticated clone
description: Register a GitHub App with gh-app-auth, wire it into Git, and clone a private repository.
sidebar:
  order: 3
---

This walkthrough registers one GitHub App, configures Git to use it, and clones a private
repository. It describes a workstation with an OS keyring. For a headless machine or a pipeline, see
[the CI guide](../guides/ci/index.md).

## Before you start

You need:

- The [GitHub CLI](https://cli.github.com/) and Git installed, and the extension from
  [installation](installation.md).
- A GitHub App that is already installed on the organization that owns your repository, with access to
  that repository and the **Contents: read** permission. See
  [create and install a GitHub App](../guides/create-github-app.md) if you still need one.
- The App's **Client ID**, shown on its settings page (for example `Iv1.…`), and a private key
  downloaded as a `.pem` file.

Creating a GitHub App, installing it on an organization, or generating its key can require
organization owner or App manager rights. If you cannot do these yourself, ask an administrator to
install the App, and obtain the Client ID and a private key through your organization's approved
access or secret-sharing procedure. Never send a private key over chat or email. Do not assume every
user in an organization is allowed to create Apps.

This walkthrough uses a single GitHub App and no personal access tokens. When both an App and a PAT
match a repository, the helper compares them by priority, which changes which credential is chosen.
Keep the first run simple.

## Placeholders

Every value below is a placeholder. Replace it with your own.

| Placeholder | Meaning |
| --- | --- |
| `Iv1.example` | Client ID of your GitHub App |
| `$HOME/keys/example-app.pem` | Path to the private key file you downloaded |
| `github.com/example-org/` | Host and organization that the App covers |
| `Example App` | A friendly name shown in `list` output |
| `example-org/private-repo` | A private repository the App can read |

## 1. Install the extension

```bash
gh extension install AmadeusITGroup/gh-app-auth
```

## 2. Restrict the key file

The key file must not be readable by other users. On Linux and macOS:

```bash
chmod 600 "$HOME/keys/example-app.pem"
```

`chmod` does not apply on Windows, where file permissions are set through ACLs. The permission check
that gh-app-auth performs on key files is skipped on Windows, so protecting the file is your
responsibility there.

## 3. Register the App

```bash
gh app-auth setup \
  --client-id Iv1.example \
  --key-file "$HOME/keys/example-app.pem" \
  --patterns 'github.com/example-org/' \
  --name 'Example App'
```

`--client-id` is preferred. `--app-id` with the numeric App ID also works. By default the key is
stored in the OS keyring. Passing the key through the `GH_APP_PRIVATE_KEY` environment variable
instead of `--key-file` changes how you supply it, not where it is stored: it is still persisted by
gh-app-auth.

On a headless machine without an OS keyring, add `--use-filesystem` together with `--key-file`. That is
the reliable route there. Do not count on `setup` falling back to the filesystem by itself: that
automatic fallback can leave a configuration that is rejected later, and the key file it writes is
plain, unencrypted text. See [known limitations](../troubleshoot/known-limitations.md). A personal
access token has the same limitation, so keep PATs on machines with a keyring.

The installation ID is detected automatically. The full flag list is in the
[`setup` reference](../reference/cli/setup.md).

## 4. Tell Git to use the helper

```bash
gh app-auth gitconfig --sync
```

This writes credential helper entries to your global Git configuration for the credentials you
registered. `--local` limits the change to the current repository. See the
[Git credential helper guide](../guides/git-credential-helper.md).

## 5. Check the configuration

```bash
gh app-auth list --verify-keys
```

Expected outcome: `Example App` appears in the list with its pattern, and its private key is
reported as accessible. If the key cannot be read, fix that here before going further.

## 6. Test the credential

```bash
gh app-auth test --repo github.com/example-org/private-repo
```

This checks that a pattern matches the repository, that a JWT can be signed with the key, that an
installation token can be retrieved, and that the GitHub API accepts it. Expected outcome: the command
reports success for each check; any failure names the step that failed. The command contacts GitHub
with the credentials you configured. Nothing on this site runs it for you, and the examples here
never use real credentials.

## 7. Clone

```bash
git clone https://github.com/example-org/private-repo.git
```

Git asks its configured credential helpers for credentials, the helper selects `Example App` for this
remote, and the clone should complete. If Git still prompts, another credential helper earlier in your
Git configuration may be answering first; check `git config --show-origin --get-all credential.helper`.

## What changed on your machine

- A configuration file at `~/.config/gh/extensions/gh-app-auth/config.yml`. Set `GH_APP_AUTH_CONFIG`
  to use another path.
- The private key in the OS keyring. If you used `--use-filesystem --key-file` instead, the
  configuration records the path of your key file and gh-app-auth reads that file; it does not
  encrypt it.
- Credential helper entries in your Git configuration.
- No installation token on disk. Tokens are generated on demand and cached in memory only.

## Next steps

- Learn [how the helper picks a credential](../understand/credential-routing.mdx) before you add a
  second App or a PAT.
- Run a command with a credential token using [`exec`](first-exec.md).
- Set it up in a pipeline with the [CI guide](../guides/ci/index.md).
