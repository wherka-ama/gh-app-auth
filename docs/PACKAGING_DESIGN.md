# RPM and DEB Package Support - Design Document

## Overview

This document outlines the design for adding RPM (RHEL/Fedora) and DEB (Ubuntu/Debian) package creation and releasing to the gh-app-auth project.

## Research Summary

### Current State
- **Release Mechanism**: Uses `cli/gh-extension-precompile@v2` for cross-platform binary builds
- **Current Assets**: Raw binaries for Linux (amd64/arm64), macOS (amd64/arm64), Windows (amd64/arm64)
- **Missing**: Native Linux package formats (RPM/DEB)

### Selected Tool: nFPM

**Why nFPM?**
- Pure Go implementation (0 dependencies on Ruby, tar, etc.)
- Supports DEB, RPM, APK, IPK, and Arch Linux packages
- Simple YAML-based configuration
- Active maintenance by GoReleaser team
- Can be used as CLI tool or GitHub Action

**Alternative Considered**: FPM
- Rejected due to Ruby dependencies and complexity

## Solution Architecture

### 1. Makefile Targets

New targets to be added:

```
package-deps        Install nFPM packaging tool
package-deb         Build DEB package for amd64
package-rpm         Build RPM package for amd64
package-deb-arm64   Build DEB package for arm64
package-rpm-arm64   Build RPM package for arm64
packages            Build all packages (deb/rpm for amd64/arm64)
packages-local      Build packages for local architecture only
```

### 2. nFPM Configuration

**File**: `nfpm.yaml` (templated for architecture-specific builds)

Key configuration elements:
- **Name**: gh-app-auth
- **Arch**: amd64, arm64 (mapped to x86_64, aarch64 for RPM)
- **Platform**: linux
- **Version**: Extracted from git tags
- **Section**: utils
- **Priority**: optional
- **Maintainer**: AmadeusITGroup <maintainers@amadeus.com>
- **Description**: GitHub CLI extension for GitHub App authentication
- **License**: Apache-2.0
- **Homepage**: https://github.com/AmadeusITGroup/gh-app-auth
- **Dependencies**: gh (GitHub CLI)
- **Files**: Binary → /usr/bin/gh-app-auth

### 3. GitHub Actions Workflow Changes

**File**: `.github/workflows/release.yml`

New job: `build-packages`

```yaml
build-packages:
  runs-on: ubuntu-latest
  needs: release  # Run after precompile creates binaries
  strategy:
    matrix:
      include:
        - arch: amd64
          rpm_arch: x86_64
        - arch: arm64
          rpm_arch: aarch64
  steps:
    1. Checkout code
    2. Download binary artifacts from precompile job
    3. Create DEB package using nFPM action
    4. Create RPM package using nFPM action
    5. Upload packages to release
```

### 4. Architecture Mapping

| Go Arch | DEB Arch | RPM Arch |
|---------|----------|----------|
| amd64   | amd64    | x86_64   |
| arm64   | arm64    | aarch64  |

## Implementation Details

### nFPM Installation

**Local (Makefile)**:
```bash
go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest
```

**CI (GitHub Actions)**:
Install via `go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest` (official, always latest)

### Package Metadata

**DEB-specific**:
- Priority: optional
- Section: utils
- Recommends: gh

**RPM-specific**:
- Group: System Environment/Base
- Requires: gh

### File Installation

```yaml
files:
  "./dist/gh-app-auth-linux-{{ .Arch }}": "/usr/bin/gh-app-auth"
```

## Benefits

1. **Native Installation**: Users can install via `apt install` or `yum install`
2. **Dependency Management**: Automatic resolution of gh CLI dependency
3. **Enterprise Ready**: Standard package formats for enterprise Linux distributions
4. **Repository Ready**: Can be published to custom APT/YUM repositories

## Compatibility

### Target Distributions

**DEB**:
- Ubuntu 20.04+ (Focal)
- Ubuntu 22.04+ (Jammy)
- Ubuntu 24.04+ (Noble)
- Debian 11+ (Bullseye)
- Debian 12+ (Bookworm)

**RPM**:
- RHEL 8+
- RHEL 9+
- Fedora 38+
- Amazon Linux 2023
- CentOS Stream 8/9
- Rocky Linux 8/9
- AlmaLinux 8/9

## Testing Strategy

1. **Local Testing**: `make packages-local` builds for local arch
2. **CI Testing**: Build packages on PRs (dry-run mode)
3. **Package Validation**: Use `dpkg-deb -I` and `rpm -qip` to verify metadata
4. **Installation Testing**: Test in Docker containers of target distributions

## Rollout Plan

1. **Phase 1**: Implement Makefile targets
2. **Phase 2**: Add GitHub Actions workflow
3. **Phase 3**: Test with next release
4. **Phase 4**: Document installation instructions

## Future Enhancements

- Sign packages with GPG key
- Publish to package repositories (APT/YUM)
- Add APK support for Alpine Linux
