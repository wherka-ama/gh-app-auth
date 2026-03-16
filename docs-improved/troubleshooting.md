# Troubleshooting

## Quick diagnostics

```bash
# Check configured credentials
gh app-auth list --verify-keys

# Test authentication for a repo
gh app-auth test --repo github.com/myorg/repo

# See which credential matches a URL
gh app-auth scope --repo github.com/myorg/repo

# Show git credential helper configuration
git config --show-origin --get-regexp credential
```

## Common issues

### Git still prompts for username/password

1. Confirm your repo URL matches a configured pattern: `gh app-auth list`
2. Re-sync credential helpers: `gh app-auth gitconfig --sync --global`
3. Check that git sees the helper: `git config --show-origin --get-regexp credential`

### Wrong credential is selected

For Apps, the longest matching prefix wins. Priority only matters when choosing
between an App and a PAT. Check your config:

```bash
gh app-auth list            # see all patterns and priorities
gh app-auth scope --repo <url>  # see what would be selected
```

Fix by adjusting `--priority` or using more specific patterns.

### "No credential found"

- Ensure at least one App or PAT is configured (`gh app-auth list`).
- Patterns must match the URL prefix. Include a trailing slash for org-level
  matches (e.g., `github.com/myorg/`).

### Extension not found by git

```bash
gh extension list | grep app-auth
```

If missing, reinstall: `gh extension install AmadeusITGroup/gh-app-auth`

### Credential helper conflicts

Clear other helpers before syncing:

```bash
git config --global --unset-all credential.helper
gh app-auth gitconfig --sync --global
```

### Bitbucket rejects credentials (HTTP 401)

Bitbucket Server requires a real username. Reconfigure the PAT:

```bash
gh app-auth setup \
  --pat bbpat_xxx \
  --username your_username \
  --patterns "bitbucket.example.com/" \
  --name "Bitbucket"
gh app-auth gitconfig --sync --global
```

### CI job fails after 1 hour

GitHub App tokens auto-refresh. If the job still fails, check that the
workspace and keyring are accessible throughout the job. See [ci-cd.md](ci-cd.md).

## Debug logging

### Enable

```bash
# Log to default location
export GH_APP_AUTH_DEBUG_LOG=1

# Or specify a file
export GH_APP_AUTH_DEBUG_LOG="/tmp/gh-app-auth-debug.log"
```

Then reproduce the issue (e.g., `git clone ...`).

### Read logs

Default location: `~/.config/gh/extensions/gh-app-auth/debug.log`

```bash
# Follow live
tail -f ~/.config/gh/extensions/gh-app-auth/debug.log

# Find errors
grep FLOW_ERROR ~/.config/gh/extensions/gh-app-auth/debug.log

# Trace a specific session
grep SESSION_START ~/.config/gh/extensions/gh-app-auth/debug.log | tail -1
```

### Log format

```
[TIMESTAMP] EVENT [OPERATION_ID] key=value ...
```

| Event | Meaning |
|-------|---------|
| `SESSION_START` | New credential request |
| `FLOW_STEP` | Authentication step (match, JWT, token exchange) |
| `FLOW_SUCCESS` | Credential returned to git |
| `FLOW_ERROR` | Authentication failed |
| `SESSION_END` | Request complete |

Tokens are never logged in plain text. Token hashes (`sha256:...`) appear in
debug output for correlation.

## Reset everything

```bash
gh app-auth gitconfig --clean --global
gh app-auth remove --all --force
gh app-auth remove --all-pats --force
```

Then start fresh with `gh app-auth setup`.
