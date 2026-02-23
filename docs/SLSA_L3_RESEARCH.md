# SLSA L3 Immutable Releases Research

## Research Date: February 23, 2026

## Objective
Research and document best practices for implementing SLSA L3 compliant releases with artifact attestations and immutable releases for gh-app-auth.

---

## References

### Primary Resources

1. **SLSA Specification v1.0 - Security Levels**
   - URL: https://slsa.dev/spec/v1.0/levels
   - Key topics: Build track levels (L0-L3), requirements for each level

2. **GitHub Blog: Enhance build security and reach SLSA Level 3 with GitHub Artifact Attestations**
   - URL: https://github.blog/enterprise-software/devsecops/enhance-build-security-and-reach-slsa-level-3-with-github-artifact-attestations/
   - Key topics: Build provenance, reusable workflows, ephemeral machines

3. **GitHub Changelog: Code-to-cloud traceability and SLSA Build Level 3 security**
   - URL: https://github.blog/changelog/2026-01-20-strengthen-your-supply-chain-with-code-to-cloud-traceability-and-slsa-build-level-3-security/
   - Key topics: Artifact metadata APIs, linked artifacts view

4. **SLSA: Get Started Guide**
   - URL: https://slsa.dev/how-to/get-started
   - Key topics: SLSA 1-3 requirements, GitHub Actions guidance

5. **GitHub Changelog: Releases now support immutability in public preview**
   - URL: https://github.blog/changelog/2025-08-26-releases-now-support-immutability-in-public-preview/
   - Key topics: Immutable assets, tag protection, release attestations

6. **GitHub Changelog: Immutable releases are now generally available**
   - URL: https://github.blog/changelog/2025-10-28-immutable-releases-are-now-generally-available/
   - Key topics: GA status, enabling immutability, verification

7. **Release Please (Google)**
   - URL: https://github.com/googleapis/release-please
   - Key topics: Automated releases, conventional commits, CHANGELOG generation

8. **GitHub Action: attest-build-provenance**
   - URL: https://github.com/actions/attest-build-provenance
   - Key topics: Build provenance attestation, workflow integration

---

## SLSA Build Track Levels Summary

### Build L0: No guarantees
- **Use case**: Development or test builds
- **Requirements**: None
- **Security**: n/a

### Build L1: Provenance exists
- **Use case**: Quick benefits without changing build workflows
- **Requirements**:
  - Consistent build process
  - Provenance exists (build platform, process, inputs)
  - Provenance distributed to consumers
- **Benefits**:
  - Easier debugging, patching, rebuilding
  - Prevents release mistakes
  - Software inventory creation

### Build L2: Hosted build platform
- **Use case**: Moderate security benefits while waiting for L3 hardening
- **Requirements** (all of L1, plus):
  - Build runs on dedicated infrastructure (not workstation)
  - Provenance tied to infrastructure via digital signature
  - Downstream verification validates provenance authenticity
- **Benefits**:
  - Prevents tampering after build
  - Deters adversaries facing legal/financial risk
  - Reduces attack surface
  - Allows early migration to supported platforms

### Build L3: Hardened builds
- **Use case**: Most software releases
- **Requirements** (all of L2, plus):
  - Build platform implements strong controls:
    - Prevent runs from influencing one another (even within same project)
    - Prevent secret material used for signing from being accessible to user-defined build steps
- **Benefits**:
  - Prevents tampering by insider threats, compromised credentials, other tenants
  - Greatly reduces impact of compromised upload credentials (requires difficult exploit)
  - Strong confidence package was built from official source and process
- **Key mechanism**: Provenance signed by key only accessible to build platform

---

## GitHub Artifact Attestations

### Overview
GitHub Artifact Attestations streamline provenance establishment for builds. By enabling provenance generation directly within GitHub Actions workflows, each artifact includes a verifiable record of its build history.

### SLSA Level Achievements
- **SLSA Level 1**: Generating build provenance
- **SLSA Level 2**: Using GitHub Artifact Attestations on GitHub-hosted runners (default)
- **SLSA Level 3**: Using a reusable workflow for provenance generation

### Key Features
1. **Build provenance made simple**: No need to handle cryptographic key material
2. **Secure signing with ephemeral machines**: Each build in clean, isolated environment
3. **Reusable workflows**: Central enforcement of build security across all projects

### Required Permissions
```yaml
permissions:
  id-token: write        # Mint OIDC token for Sigstore signing certificate
  attestations: write      # Persist attestation
  artifact-metadata: write # Generate artifact metadata storage records
```

### Usage Example
```yaml
- uses: actions/attest-build-provenance@v3
  with:
    subject-path: '<PATH TO ARTIFACT>'
```

### Inputs
- `subject-path`: Path to artifact (glob pattern or list allowed, max 1024)
- `subject-digest`: SHA256 digest (format: "sha256:hex_digest")
- `subject-name`: Subject name for attestation
- `subject-checksums`: Path to checksums file
- `push-to-registry`: Push attestation to image registry (default: false)
- `create-storage-record`: Create storage record (default: true)
- `show-summary`: Attach attestations to workflow summary (default: true)

### Outputs
- `attestation-id`: ID of created attestation
- `attestation-url`: URL to view attestation
- `bundle-path`: Path to attestation bundle (Sigstore bundle format)

### Verification
```bash
# Verify artifact was signed using specific reusable workflow
gh artifact verify <file-path> --signer-workflow <owner>/<repo>/.github/workflows/sign-artifact.yml
```

---

## Immutable Releases

### Overview
Immutable releases add a new layer of supply chain security by protecting assets and tags from tampering after publication.

### Features
1. **Immutable assets**: Once published, assets cannot be added, modified, or deleted
2. **Tag protection**: Tags for immutable releases cannot be deleted or moved
3. **Release attestations**: Signed attestations for easy verification of authenticity and integrity

### Benefits
- Protection from supply chain attacks
- Verifiable authenticity even outside GitHub
- Sigstore bundle format for compatibility

### Enabling Immutability
- Enable at repository or organization level in settings
- All new releases become immutable once enabled
- Existing releases remain mutable unless republished
- Disabling doesn't affect releases created while enabled

### Verification
Attestations use Sigstore bundle format, allowing verification via:
- GitHub CLI: `gh attestation verify`
- Sigstore-compatible tooling for CI/CD policy enforcement

---

## Release Please

### Overview
Release Please automates releases based on conventional commits. It maintains Release PRs that are kept up-to-date as work is merged.

### Key Concepts

#### Release PR Lifecycle
1. **autorelease: pending**: Initial state before merge
2. **autorelease: tagged**: PR merged, release tagged
3. **autorelease: snapshot**: Special state for snapshot versions
4. **autorelease: published**: GitHub release published (recommended convention)

#### Conventional Commits
- `fix:` → SemVer patch
- `feat:` → SemVer minor
- `feat!:`, `fix!:`, `refactor!:` → SemVer major (breaking change)

#### What Release Please Does on Merge
1. Updates CHANGELOG.md
2. Updates language-specific files (e.g., package.json, version.txt)
3. Tags commit with version number
4. Creates GitHub Release based on tag

### Supported Strategies
- `go`: Go repositories
- `simple`: version.txt + CHANGELOG.md
- `terraform-module`: README.md version + CHANGELOG.md
- And many more (node, python, rust, java, etc.)

### Deployment Options
1. **GitHub Action** (recommended): googleapis/release-please-action
2. **CLI**: Running release-please CLI

### GitHub Action Setup
```yaml
- uses: googleapis/release-please-action@v4
  with:
    release-type: go
```

---

## Integration Strategy for gh-app-auth

### Current State
- Release workflow: `.github/workflows/release.yml`
- Triggered on tags: `v*`
- Uses **Makefile targets** (`make release packages`) for cross-platform builds
- Builds binaries for multiple platforms using Go cross-compilation
- Creates DEB and RPM packages via nfpm
- Release artifacts uploaded to GitHub release

### Target State
1. **Automated versioning**: Release Please for CHANGELOG and version management
2. **SLSA L3 compliance**: Reusable workflow + attest-build-provenance
3. **Immutable releases**: Enable immutability at repository level
4. **Attestation verification**: Document how consumers verify releases

### Required Workflow Changes
1. Add `release-please` workflow for automated release PRs
2. Update release workflow with attestations and immutability
3. Create reusable attestation workflow for SLSA L3
4. Enable immutable releases in repository settings
5. Document verification process for users

### Permission Requirements
```yaml
permissions:
  contents: write          # Create releases
  id-token: write          # OIDC token for signing
  attestations: write      # Create attestations
  artifact-metadata: write # Artifact metadata
  pull-requests: write     # For release-please
```

---

## Research Conclusions

### Key Findings
1. GitHub Artifact Attestations + GitHub-hosted runners = SLSA L2 by default
2. Using reusable workflow for attestation = SLSA L3
3. Immutable releases provide additional tamper protection
4. Release Please automates versioning and CHANGELOG maintenance
5. Sigstore integration eliminates key management burden

### Recommended Approach
- Phase 1: Implement Release Please for automated releases
- Phase 2: Add artifact attestations (SLSA L2)
- Phase 3: Create reusable workflow for SLSA L3
- Phase 4: Enable immutable releases and update documentation

---

*Document generated from research conducted on 2026-02-23*
