---
title: Git credential helper
description: What gh app-auth gitconfig writes into your Git configuration, how to scope it, and how to remove it.
sidebar:
  order: 4
---

Git only calls gh-app-auth if its configuration says so. `gh app-auth gitconfig` writes and removes
that configuration for you.

## Sync and clean

```bash
gh app-auth gitconfig --sync
gh app-auth gitconfig --clean
```

- `--sync` adds a helper entry for every pattern of every configured App and PAT.
- `--clean` removes every helper entry that gh-app-auth wrote, for all credentials.
- `--global` (the default) edits your user Git configuration. `--local` edits the current
  repository only.
- `--auto` configures a single App from environment variables; see
  [environment variables](../reference/environment-variables.md).

Run `--sync` again after you add or change a credential.

## What gets written

For each pattern, `--sync` sets a credential helper in the matching Git credential context, along
these lines:

```text
[credential "https://github.com/example-org"]
    helper =
    helper = !"/path/to/gh-app-auth" git-credential --pattern github.com/example-org/
    useHttpPath = true
```

It clears earlier helpers in that context first, sets `useHttpPath` for hosts with path-specific
patterns so Git sends the repository path, and puts path-specific helpers ahead of generic host
helpers. Patterns that contain spaces or glob characters are quoted. On Windows the helper is a
generated wrapper script.

Inspect the result with:

```bash
git config --show-origin --get-regexp '^credential\.'
```

## Why the pattern is passed to the helper

The helper receives `--pattern` so it can honor the Git context that selected it. That shortcut and
the URL matcher are different mechanisms; see
[how credential routing works](../understand/credential-routing.mdx).

## Other helpers

Git consults helpers in configuration order and uses the first that answers. If another helper
answers first, gh-app-auth is never asked. Check the order if Git still prompts, and see
[troubleshooting](../troubleshoot/index.md).

## Remove credentials

> **Destructive scope.** `gitconfig --clean` removes the helper entries for **all** credentials, and
> `remove --all` removes **all** configured Apps and their managed secrets. Neither can be undone
> without registering the credentials again. Outside an isolated, disposable CI job, prefer removing
> one credential at a time, or run `remove --all` interactively and read the prompt.

```bash
gh app-auth gitconfig --clean
gh app-auth remove --all
```

The first command removes helper entries and the second removes configured Apps and their stored
secrets; it asks for confirmation unless you add `--force`, which belongs only in an isolated CI job.
`remove --all` does not remove PATs. Remove those in a separate invocation, with `remove --all-pats` or
`remove --pat-name`; target flags cannot be combined.
