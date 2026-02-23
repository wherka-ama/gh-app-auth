# Immutable Releases Setup

## Overview

This document describes the manual step required to enable immutable releases for the gh-app-auth repository.

## What are Immutable Releases?

Immutable releases provide an additional layer of supply chain security by:
- **Locking release assets**: Once published, release assets cannot be added, modified, or deleted
- **Protecting tags**: Tags associated with immutable releases cannot be deleted or force-pushed
- **Providing attestations**: Releases receive signed attestations for verification

## Enabling Immutable Releases

This is a one-time manual configuration step that must be done by a repository administrator.

### Step-by-Step Instructions

1. **Navigate to Repository Settings**
   - Go to: `https://github.com/AmadeusITGroup/gh-app-auth/settings`
   - (Requires repository admin access)

2. **Find the Releases Section**
   - Scroll down to the "Releases" section
   - Or use the left sidebar: **Settings** → **General** → **Releases**

3. **Enable Immutability**
   - Check the box labeled **"Make new releases immutable"**
   - Click **Save changes**

4. **Verify the Setting**
   - Create a test release to confirm immutability is enforced
   - The release should show an "Immutable" indicator

### What Happens When Enabled?

- ✅ All **new** releases will be immutable
- ✅ Release assets will be locked after publication
- ✅ Tags will be protected from deletion/modification
- ⚠️ **Existing** releases remain mutable unless republished
- ⚠️ Disabling the setting doesn't affect already-immutable releases

## Verification

### Check if Immutability is Enabled

1. Go to any release page
2. Look for the "Immutable" badge/indicator
3. Attempting to edit assets should fail

### Using the GitHub API

```bash
# Check repository settings
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/AmadeusITGroup/gh-app-auth

# Look for: "releases_are_immutable": true in the response
```

## Impact on Release Workflow

Our release workflow (`.github/workflows/release.yml`) is designed to work with immutable releases:

- **Draft releases**: Created first, then published
- **Attestations**: Generated before final publication
- **Immutable on publish**: When the draft is published, it becomes immutable

No changes to the workflow are required when enabling this setting.

## Important Notes

### For Repository Administrators

- **One-time setup**: This only needs to be enabled once
- **Non-reversible for existing releases**: Once a release is immutable, it stays immutable
- **No impact on CI/CD**: The workflow works transparently with this setting

### For Users

- **Enhanced security**: Releases cannot be tampered with after publication
- **Verifiable**: Users can verify release integrity using `gh attestation verify`
- **Trust**: Immutable releases + SLSA attestations = high confidence in artifact authenticity

## Related Documentation

- [SLSA Compliance](SLSA_COMPLIANCE.md)
- [GitHub Docs: Immutable Releases](https://docs.github.com/code-security/supply-chain-security/understanding-your-software-supply-chain/immutable-releases)
- [GitHub Changelog: GA Announcement](https://github.blog/changelog/2025-10-28-immutable-releases-are-now-generally-available/)

## Checklist for Implementation

- [ ] Repository admin navigates to Settings → Releases
- [ ] "Make new releases immutable" is enabled
- [ ] Changes are saved
- [ ] Test release created to verify
- [ ] Documentation updated (this file)

---

**Status**: ⏳ Pending manual configuration by repository administrator

**Last updated**: February 2026
