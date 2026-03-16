# Security

## Private key storage

The extension stores private keys and PATs in the OS keyring by default:

| Platform | Backend |
|----------|---------|
| macOS | Keychain |
| Windows | Credential Manager |
| Linux | Secret Service API (GNOME Keyring, KWallet) |

If the keyring is unavailable the extension falls back to filesystem storage
under `~/.config/gh/extensions/gh-app-auth/secrets/`.

Use `gh app-auth list` to see which backend each credential uses.

### File permissions

Key files on disk must have mode `600` or `400`. The extension rejects
world-readable or group-readable files.

```bash
chmod 600 ~/keys/my-app.pem
```

## Token handling

Installation tokens (the credentials returned to Git) are held in memory only:

- Cached for 55 minutes per process, regenerated on expiry.
- Never written to disk.
- Best-effort zeroing on cleanup. Go strings are immutable, so the
  original may remain until garbage collection.
- Lost when the process exits.

JWT tokens (used to request installation tokens) are generated on demand and not
cached.

## What is stored where

| Data | Location | Encrypted |
|------|----------|-----------|
| App IDs, patterns, priorities | Config file | No (not sensitive) |
| Private keys | OS keyring | Yes |
| PATs | OS keyring | Yes |
| Installation tokens | Memory only | N/A |
| JWT tokens | Not stored | N/A |

## Input validation

- **Secret file paths**: `filepath.Base()` is applied to app names when
  building filesystem paths, preventing `../` traversal in secret storage.
- **Key file permissions**: must be `600` or `400`; world-readable files are
  rejected.
- **App ID**: must be a positive integer.

## Logging

Debug logs (when enabled) include flow steps and token hashes (`sha256:...`).
Full tokens, private keys, and credentials are never logged.

## CI/CD considerations

- Store private keys in your CI system's secrets management (GitHub Secrets,
  Jenkins Credentials, etc.).
- On self-hosted runners, clean up after each job:

  ```bash
  gh app-auth gitconfig --clean --global
  gh app-auth remove --all --force
  ```

  Or use the composite action with `cleanup-on-exit: 'true'`.

## Reporting vulnerabilities

Report security issues through
[GitHub Security Advisories](https://github.com/AmadeusITGroup/gh-app-auth/security).
Do not use public issues for security reports.

See [SECURITY.md](../.github/SECURITY.md) for contact details.
