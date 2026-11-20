---
title: Create and install a GitHub App
description: What to ask for, or create yourself, so gh-app-auth has a GitHub App with least-privilege access.
sidebar:
  order: 2
---

gh-app-auth needs four things from GitHub: the App's **Client ID** (or numeric App ID), a **private
key**, an **installation** on the owner of the repositories you need, and the **permissions** the
installation was granted.

## Who can do what

Creating an App, installing it on an organization and generating its private key can each require
an organization owner or an App manager. You may not be allowed to create Apps at all. If you are
not, ask an administrator to do the steps below and to hand you the Client ID and the key through an
approved secret-sharing or access procedure. Organization creation is not required: an existing
organization or a personal account can own the App.

## Choose the permissions

Grant the least privilege your workflow needs:

| Workflow | Repository permission |
| --- | --- |
| Clone and fetch | **Contents: Read-only** |
| Push | **Contents: Read and write** |
| API calls through `gh app-auth exec` | Only the permissions those calls need |

GitHub always adds **Metadata: Read-only**. Patterns in gh-app-auth never reduce what a token can do;
the installation's permissions and repository selection do. See
[how routing works](../understand/credential-routing.mdx).

## Create the App

1. In GitHub, open the account or organization settings, then **Developer settings**, then **GitHub Apps**,
   and choose **New GitHub App**.
2. Give it a name and a homepage URL. A webhook is not needed for gh-app-auth, so you can leave
   webhooks inactive.
3. Set the repository permissions from the table above.
4. Decide who can install it. An App limited to your own account or organization is the safest
   choice.
5. Create the App, and note the **Client ID** on its settings page. Client IDs identify the App
   publicly; they are not secrets.

## Install it

Install the App on the organization or account that owns the repositories, and choose **Only select
repositories** where you can. A fully scoped installation is easier to reason about than **All
repositories**. gh-app-auth detects the installation ID automatically; you can also pass
`--installation-id` to `setup`.

## Generate and store the private key

Generate a private key from the App's settings. GitHub downloads a `.pem` file once; it cannot be
shown again. Store it according to your organization's secret-handling rules. Then:

- On a workstation with an OS keyring, register the key with `setup --key-file`. See
  [your first authenticated clone](../start-here/first-clone.md).
- On an ephemeral CI runner, follow the [CI guide](ci/index.md), which writes the key to a job-owned file.

Protect the key file: `chmod 600` on Linux and macOS. gh-app-auth rejects a key file that is readable
by group or others on POSIX systems, and skips that check on Windows, where ACLs are your
responsibility.

## Check the result

```bash
gh app-auth setup --client-id Iv1.example --key-file "$HOME/keys/example-app.pem" \
  --patterns 'github.com/example-org/' --name 'Example App'
gh app-auth list --verify-keys
gh app-auth test --repo github.com/example-org/private-repo
```

All values are placeholders. If `test` fails, see [troubleshooting](../troubleshoot/index.md).
