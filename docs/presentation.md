---
title: "Solving GitHub App Authentication in CI/CD"
subtitle: "Multi-organization Git operations with gh-app-auth"
author: "gh-app-auth Project"
---

# GitHub App Authentication in CI/CD

**One credential helper for Git across several organizations**

This optional deck summarizes gh-app-auth. The canonical, maintained documentation is
[the guides](README.md); every claim here is explained there.

---

# The problem {#problem}

- A GitHub App installation token belongs to one installation, and expires after an hour.
- Reaching repositories in several organizations, including submodules, needs one credential per
  organization.
- Robot user accounts avoid this but weaken governance; Apps fit governance but are awkward for Git.

Read more: [origins](understand/origins.md).

---

# The approach {#approach}

- Register each GitHub App, or personal access token, with the repository patterns it covers.
- Git asks a credential helper per request; the helper picks the credential and mints a token.
- Private keys stay in the OS keyring by default.

Read more: [identity and credentials](understand/identity-and-credentials.md).

---

# Routing {#routing}

- App patterns match whole path segments; the longest match selects the App.
- Apps and personal access tokens are then compared by priority.
- Patterns choose a credential. The installation's permissions decide what it can do.

Read more: [credential routing](understand/credential-routing.mdx).

---

# Setup {#setup}

```bash
gh extension install AmadeusITGroup/gh-app-auth
gh app-auth setup --client-id Iv1.example --key-file "$HOME/keys/app.pem" \
  --patterns 'github.com/example-org/' --name 'Example App'
gh app-auth gitconfig --sync
git clone https://github.com/example-org/private-repo.git
```

All values are placeholders. Read more: [first clone](start-here/first-clone.md).

---

# Security posture {#security}

- Installation tokens are memory-only, cached per process for 55 minutes.
- No daemon and no background renewal; `exec` hands its child a token snapshot.
- Filesystem key storage is plain files; fallbacks have known limits.

Read more: [security model](understand/security-model.md) and
[known limitations](troubleshoot/known-limitations.md).

---

# In CI {#ci}

- Use a disposable, isolated job with its own home directory.
- Pass secrets through the environment, write the key to a job-owned file, and always clean up.

Read more: [CI guides](guides/ci/index.md).

---

# Where to go next {#next}

- [Documentation hub](README.md)
- [CLI reference](reference/cli/index.md)
- [Troubleshooting](troubleshoot/index.md)
- [Contributing](contribute/index.md)
