---
title: GitHub Actions
description: A complete GitHub Actions workflow that uses gh-app-auth on a GitHub-hosted runner with explicit setup and cleanup.
sidebar:
  order: 2
---

This workflow is for GitHub-hosted `ubuntu-latest` runners only. Each job starts on a fresh virtual
machine that is discarded afterwards. It is not a recipe for shared or self-hosted machines.

The job uses its own home directory under `RUNNER_TEMP`, so the extension, its configuration and the
Git configuration it writes stay inside the job. The private key is written to a file owned by the job
and passed with `--use-filesystem --key-file`, which is the reliable path on a machine without an OS
keyring.

```yaml
name: private-clone

on: workflow_dispatch

permissions:
  contents: read

jobs:
  clone:
    runs-on: ubuntu-latest
    env:
      APP_CLIENT_ID: ${{ vars.APP_CLIENT_ID }}
    steps:
      - name: Prepare isolated home
        run: |
          job_home="$(mktemp -d "$RUNNER_TEMP/gh-app-auth-home.XXXXXX")"
          echo "HOME=$job_home" >> "$GITHUB_ENV"
          echo "GH_CONFIG_DIR=$job_home/.config/gh" >> "$GITHUB_ENV"

      - name: Configure gh-app-auth
        env:
          APP_PRIVATE_KEY: ${{ secrets.APP_PRIVATE_KEY }}
          GH_TOKEN: ${{ github.token }}
        run: |
          umask 077
          app_key_file="$(mktemp "$RUNNER_TEMP/app-key.XXXXXX")"
          echo "APP_KEY_FILE=$app_key_file" >> "$GITHUB_ENV"
          printf '%s\n' "$APP_PRIVATE_KEY" > "$app_key_file"
          gh extension install AmadeusITGroup/gh-app-auth --pin v0.2.0
          gh app-auth setup \
            --use-filesystem \
            --client-id "$APP_CLIENT_ID" \
            --key-file "$app_key_file" \
            --patterns 'github.com/example-org/' \
            --name 'CI App'
          gh app-auth gitconfig --sync --global

      - name: Clone a private repository
        run: git clone https://github.com/example-org/private-repo.git

      - name: Remove credentials
        if: always()
        run: |
          status=0
          gh app-auth gitconfig --clean --global || status=1
          gh app-auth remove --client-id "$APP_CLIENT_ID" --force || status=1
          if [ -n "${APP_KEY_FILE:-}" ]; then rm -f -- "$APP_KEY_FILE" || status=1; fi
          exit "$status"
```

The Client ID lives in a repository variable and the private key in a secret. Both reach the shell
through `env`, never through an expression inside the script. The key variable is deliberately not
named `GH_APP_PRIVATE_KEY`: that variable takes precedence over the key file inside `setup`, and
`setup --use-filesystem` rejects it together with `--key-file`.

The final step runs every cleanup command even when an earlier one fails, then reports the combined
status. It deletes the key file explicitly because `remove` deletes the credentials gh-app-auth manages,
not the file you supplied.

## What to know

- **Cleanup is best effort.** An `if: always()` step runs after failures and cancellations, but not if
  the runner itself is lost. On GitHub-hosted runners the discarded virtual machine is the real
  guarantee.
- **Cleanup is not selective.** `gitconfig --clean` removes every helper entry that gh-app-auth
  wrote, for all registered credentials, not only one App. `remove` deletes only the App you name.
  That is why this sample uses a fresh job-owned home.
- **Self-hosted runners.** A shared or non-ephemeral machine needs stronger isolation than a
  temporary home directory. That setup is not covered here.
- **Automatic key fallback is unreliable.** Without `--use-filesystem` and `--key-file`, `setup` may
  fall back to filesystem storage without recording the key's path, and the resulting configuration
  is then rejected. Stored fallback key files are plain, unencrypted files.
- **Pinning.** The sample pins the extension to `v0.2.0`. Update the pin deliberately, and pin
  third-party actions according to your own policy.
- **Do not broaden cleanup.** Avoid cleanup that wipes unrelated `credential.helper` settings or deletes
  files you did not create.
- **Tokens expire.** An installation token lasts up to an hour and the helper caches it in memory for 55
  minutes, per process. Git starts a fresh helper process for each request. `gh app-auth exec`
  takes a snapshot of the token for its child process and does not renew it while the child runs.
- **Patterns select, they do not restrict.** The App's installation permissions limit what a token
  can do. See [credential routing](../../understand/credential-routing.mdx).

## The repository's composite action

The repository ships a composite action, `setup-gh-app-auth`. Its "Register cleanup hook" step sets a
shell `trap` inside the setup step, so it is **not** a post-job hook and can remove the credentials it
just configured before later steps run. Prefer the explicit steps above; if you do use the action, set
`cleanup-on-exit: 'false'` and manage cleanup yourself.
Its inputs, outputs and limitations are in the
[GitHub Actions reference](../../reference/github-actions.md).
