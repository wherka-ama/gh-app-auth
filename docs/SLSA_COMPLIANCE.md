# SLSA Compliance

This document describes the Supply Chain Levels for Software Artifacts (SLSA) compliance status of gh-app-auth releases.

## Current SLSA Level

**SLSA Build Level 3** ✅

## What is SLSA?

[SLSA](https://slsa.dev) (Supply Chain Levels for Software Artifacts) is a security framework that provides a checklist of standards and controls to prevent tampering, improve integrity, and secure packages and infrastructure in software supply chains.

## SLSA Build Track Levels

### Build L0: No guarantees
- No requirements—represents lack of SLSA
- Development or test builds

### Build L1: Provenance exists
- Package has provenance showing how it was built
- Makes debugging and patching easier
- Prevents release mistakes

### Build L2: Hosted build platform
- Builds run on dedicated infrastructure
- Provenance is digitally signed
- Prevents tampering after the build

### Build L3: Hardened builds ✅ (Current Level)
- Builds run on hardened platform with strong controls
- Runs cannot influence each other (isolation)
- Signing secrets inaccessible to build steps
- Prevents tampering by insider threats or compromised credentials

## Our Implementation

### Build Platform
- **Platform**: GitHub Actions
- **Runners**: GitHub-hosted (ephemeral, isolated)
- **Workflow**: Reusable workflow for attestation separation

### Provenance Generation
- **Tool**: GitHub Artifact Attestations (`actions/attest-build-provenance`)
- **Format**: Sigstore bundle
- **Signing**: OIDC-based via Sigstore

### Controls
1. **Ephemeral builds**: Each build runs on a fresh GitHub-hosted runner
2. **Isolated signing**: Attestation happens in separate reusable workflow
3. **No secret access**: Build steps cannot access signing credentials
4. **Reusable workflow**: Centralized, consistent attestation process

### Reusable Workflow
```yaml
uses: AmadeusITGroup/gh-app-auth/.github/workflows/attest-release.yml@main
```

This workflow:
- Runs on isolated ephemeral infrastructure
- Has no access to build secrets
- Generates signed provenance for all artifacts

## Verifying Releases

### Prerequisites
- [GitHub CLI](https://cli.github.com) installed
- Authenticated to GitHub (`gh auth login`)

### Verify a Release Artifact

Download a release artifact and verify its attestation:

```bash
# Download the artifact
gh release download v0.0.15 --pattern 'gh-app-auth_*'

# Verify the attestation
gh attestation verify gh-app-auth_linux_amd64 --owner AmadeusITGroup
```

### Expected Output
```
✓ Verification succeeded!

sha256:abc123... was attested by:
REPO             PREDICATE_TYPE                  WORKFLOW
AmadeusITGroup/gh-app-auth  https://slsa.dev/provenance/v1  .github/workflows/attest-release.yml@refs/tags/v0.0.15
```

### Verify with Specific Workflow

To verify the artifact was signed by a specific reusable workflow:

```bash
gh artifact verify gh-app-auth_linux_amd64 \
  --signer-workflow AmadeusITGroup/gh-app-auth/.github/workflows/attest-release.yml
```

## Release Integrity

### Immutable Releases

Our releases are configured as **immutable**:
- Release assets cannot be modified after publication
- Tags cannot be deleted or moved
- Provides additional protection against supply chain attacks

### What This Protects Against

| Threat | Protection |
|--------|------------|
| Artifact tampering | ✅ Immutable assets + signed provenance |
| Tag manipulation | ✅ Protected tags for immutable releases |
| Insider threats | ✅ Isolated signing, no secret access |
| Compromised credentials | ✅ Requires build platform exploit |
| Build injection | ✅ Reproducible workflow, ephemeral runners |

## Workflow Details

### Release Process

1. **Tag pushed** (`v*` pattern)
2. **Build job** runs tests and creates cross-platform binaries
3. **Attest job** calls reusable workflow to generate attestations
4. **Publish job** publishes the release with attestations attached

### Reusable Workflow Architecture

```
┌─────────────┐     ┌─────────────────┐     ┌──────────────┐
│   Build     │────▶│  Attest (L3)    │────▶│   Publish    │
│  Job        │     │  Reusable WF    │     │   Job        │
└─────────────┘     └─────────────────┘     └──────────────┘
      │                       │                      │
      ▼                       ▼                      ▼
Creates binaries       Signs provenance         Publishes
No signing access      No build secrets         Verifies
```

## Technical Details

### Provenance Contents

Each attestation includes:
- Build source (repository, commit SHA)
- Build platform (GitHub Actions)
- Build parameters
- Subject (artifact name and digest)
- Timestamp
- Builder identity

### Sigstore Integration

- **Certificate Authority**: GitHub/Sigstore
- **Signature Format**: Sigstore bundle (DSSE envelope)
- **Verification**: GitHub CLI or Sigstore tools

### Permissions Required

```yaml
permissions:
  contents: write          # Create releases
  id-token: write          # OIDC token for signing
  attestations: write      # Create attestations
  artifact-metadata: write # Artifact metadata
```

## Audit and Compliance

### For Security Auditors

1. **Verify workflow**: Check `.github/workflows/release.yml` uses reusable workflow
2. **Check permissions**: Verify id-token: write is present
3. **Review attestations**: All releases have associated attestations
4. **Test verification**: Run `gh attestation verify` on any artifact

### Compliance Reports

To generate a compliance report for a specific release:

```bash
RELEASE_VERSION="v0.0.15"
ARTIFACTS=("gh-app-auth_linux_amd64" "gh-app-auth_darwin_amd64" "gh-app-auth_windows_amd64.exe")

for artifact in "${ARTIFACTS[@]}"; do
    echo "Verifying $artifact..."
    gh attestation verify "$artifact" --owner AmadeusITGroup
done
```

## References

- [SLSA Specification v1.0](https://slsa.dev/spec/v1.0/levels)
- [GitHub Artifact Attestations](https://github.blog/enterprise-software/devsecops/enhance-build-security-and-reach-slsa-level-3-with-github-artifact-attestations/)
- [Sigstore](https://www.sigstore.dev/)
- [Immutable Releases](https://github.blog/changelog/2025-10-28-immutable-releases-are-now-generally-available/)
- [attest-build-provenance Action](https://github.com/actions/attest-build-provenance)

## Questions?

For questions about our SLSA compliance or supply chain security, please:
- Open an issue: https://github.com/AmadeusITGroup/gh-app-auth/issues
- See our security policy: [SECURITY.md](./SECURITY.md)

---

*Last updated: February 2026*
