---
name: "✨ Feature request"
about: Suggest a new feature or enhancement for gh-app-auth
title: 'Add RPM and DEB package support for Linux distributions'
labels: enhancement
assignees: ''
---

### Feature description

Add native RPM (for RHEL/Fedora) and DEB (for Ubuntu/Debian) package creation and distribution to the gh-app-auth release process. This will allow users to install gh-app-auth using standard package managers (`apt`, `yum`, `dnf`) instead of manually downloading and installing binaries.

### Problem or use case

**Current limitations:**
- Users must manually download binaries from GitHub releases
- No automatic dependency resolution (gh CLI must be manually installed first)
- No native integration with system package managers
- Enterprise environments often prefer/require package-based installations
- No automatic updates through package manager (`apt upgrade`, `yum update`)

**Use cases:**
- Enterprise Linux deployments (RHEL, Ubuntu LTS)
- Automated provisioning scripts using package managers
- CI/CD pipelines that use clean container images
- Users who prefer standard package manager workflows

### Proposed solution

**Implementation approach:**

1. **Tool Selection**: Use [nFPM](https://github.com/goreleaser/nfpm) directly via `go install` - official GoReleaser project
   - Simpler than FPM (no Ruby dependencies)
   - Supports both DEB and RPM formats
   - Active maintenance by GoReleaser team (latest release: v2.45.0)
   - Installed via `go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest`

2. **Makefile targets** (local development):
   - `make package-deps` - Install nFPM tool
   - `make package-deb` - Build DEB package for amd64
   - `make package-rpm` - Build RPM package for amd64
   - `make packages` - Build all packages (deb/rpm for amd64/arm64)

3. **GitHub Actions workflow** (release automation):
   - New job `build-packages` in release.yml
   - Run after `cli/gh-extension-precompile` creates binaries
   - Install nFPM via `go install` (official, always up-to-date)
   - Build matrix for amd64 and arm64 architectures
   - Upload packages to release assets

4. **Package metadata**:
   - **Name**: gh-app-auth
   - **Dependencies**: gh (GitHub CLI)
   - **License**: Apache-2.0
   - **Section**: utils
   - **Arch mapping**: amd64→x86_64 (RPM), arm64→aarch64 (RPM)

5. **File layout**:
   - Binary: `/usr/bin/gh-app-auth`
   - Config: `nfpm.yaml` (templated for architecture)

### Alternative solutions

**Alternative 1: FPM (Effing Package Management)**
- Rejected: Requires Ruby and additional dependencies
- More complex than needed for simple binary packaging

**Alternative 2: Native packaging tools (dpkg-deb, rpmbuild)**
- Rejected: Complex spec/control file syntax
- Platform-specific tooling requirements

**Alternative 3: goreleaser (full solution)**
- Rejected: Would require significant workflow refactoring
- nfpm provides the packaging without replacing existing build system

### Use case examples

```bash
# Ubuntu/Debian - Download and install latest release
curl -LO https://github.com/AmadeusITGroup/gh-app-auth/releases/latest/download/gh-app-auth_linux_amd64.deb
sudo apt install ./gh-app-auth_linux_amd64.deb

# RHEL/CentOS/Rocky - Download and install latest release
curl -LO https://github.com/AmadeusITGroup/gh-app-auth/releases/latest/download/gh-app-auth_linux_amd64.rpm
sudo yum install ./gh-app-auth_linux_amd64.rpm

# Future: Repository-based installation (Phase 2)
# sudo apt install gh-app-auth
# sudo yum install gh-app-auth
```

### Impact

- **Users affected**: All Linux users, especially enterprise/IT administrators
- **Frequency**: One-time setup improvement, ongoing updates benefit
- **Priority**: High for enterprise adoption, medium for individual users

### Additional context

**Target distributions:**
- **DEB**: Ubuntu 20.04+, Debian 11+, Pop!_OS, Linux Mint
- **RPM**: RHEL 8+, Fedora 38+, CentOS Stream 8/9, Rocky Linux 8/9, AlmaLinux 8/9, Amazon Linux 2023

**Related work:**
- Similar implementation pattern used by many Go CLI tools (fzf, lazygit, etc.)
- Follows the same asset release model currently used for binaries

### Implementation considerations

**Security implications:**
- Packages will be built in GitHub Actions (trusted environment)
- Future enhancement: GPG signing for package verification

**Compatibility:**
- Fully backward compatible - binary releases continue as-is
- Package installation is additive option
- No changes to existing CLI behavior

**Performance:**
- Minimal impact on CI time (packages build in parallel with existing jobs)
- Local builds: ~5 seconds per package

**Breaking changes:**
- None - this is purely additive

**Dependencies:**
- nFPM tool (installed automatically via Makefile or GitHub Action)
- Go 1.21+ (already required for build)

**Architecture support:**
- Primary: amd64 (x86_64)
- Secondary: arm64 (aarch64) - following existing binary release pattern
