---
title: Credential lifecycle
description: Verify, migrate, rotate and remove gh-app-auth credentials without deleting files you still need.
sidebar:
  order: 7
---

## Inspect

```bash
gh app-auth config --path
gh app-auth list --verify-keys
gh app-auth list --format json
```

`list --verify-keys` checks that each private key can be read. `--format` accepts `table`, `json` or
`yaml`. The default configuration file is `~/.config/gh/extensions/gh-app-auth/config.yml`; set
`GH_APP_AUTH_CONFIG` to use another file.

## Migrate key storage to the keyring

```bash
gh app-auth migrate --dry-run
gh app-auth migrate
```

`migrate` with the default `--storage keyring` copies each App's private key from its file into the OS
keyring. Check first that a keyring is actually available on the machine, because migration to a system
without one warns and can end up outside the keyring. Run `--dry-run` first to preview the move. Original
key files are **kept** unless you pass `--force`; do not add `--force` until `gh app-auth list --verify-keys`
succeeds and you have another approved copy of the key if your policy needs one.

## Moving to the filesystem

`migrate --storage filesystem` does **not** copy or export a key into a managed store. It only marks the
App's key source as `filesystem` and needs an existing configured `private_key_path`; without one it
reports an error and skips the App. The reliable way to use a file is to register it explicitly:

```bash
gh app-auth setup --use-filesystem --client-id Iv1.example \
  --key-file "$HOME/keys/example-app.pem" --patterns 'github.com/example-org/' --name 'Example App'
```

Use an existing, approved key file. See [known limitations](../troubleshoot/known-limitations.md).

## Rotate a private key

1. In the App's settings on GitHub, generate a new private key.
2. Re-run `setup` with the new key, and keep the **same identifier style, installation and name**:
   if you registered with `--client-id`, use `--client-id` again, with the same `--installation-id` if you
   set one, and the same `--name`. `setup` matches an existing entry by the App identifier (the Client ID
   when present) **and** the installation ID, merges the patterns you pass with the existing ones and
   replaces that one entry. It does not update every entry that belongs to the same App.
3. Do not change `--name` during rotation. Secrets are looked up by name, so a new name stores the key
   under a new lookup and can leave the old entry's secret orphaned.
4. Run `gh app-auth list --verify-keys` and `gh app-auth test --repo <repository>`.
5. Delete the old key in GitHub, then delete old key files according to your policy.

## Refresh scope metadata

`gh app-auth scope --refresh` re-reads which repositories an installation can see. The cached
metadata lasts 24 hours and is separate from installation tokens, which are cached in memory for 55
minutes.

## Remove

```bash
gh app-auth remove --client-id Iv1.example
gh app-auth remove --pat-name 'Bitbucket PAT'
gh app-auth gitconfig --clean
```

`remove` deletes the configuration entry and the secrets gh-app-auth manages for it. It does **not**
delete a key file you supplied with `--key-file`; delete that yourself if it should not remain.
`gitconfig --clean` removes every helper entry gh-app-auth wrote, for all credentials, so run
`gitconfig --sync` afterwards if some credentials remain.

`remove` has no installation selector. If an App has more than one configured installation, read the
confirmation prompt to check what will be removed, and do not use `--force` for an ambiguous removal.
Apps and PATs are removed with separate invocations; target flags cannot be combined.
