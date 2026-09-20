# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1](https://github.com/wherka-ama/gh-app-auth/compare/v0.1.0...v0.1.1) (2026-09-15)


### Bug Fixes

* **ci:** bump golangci-lint to v2.13.2 in lint.yml too ([d4431e7](https://github.com/wherka-ama/gh-app-auth/commit/d4431e7a95c7f96a0bf1cb54ff9938efb491656e))
* **ci:** exclude release-please-generated CHANGELOG.md from markdownlint ([2d1d1dd](https://github.com/wherka-ama/gh-app-auth/commit/2d1d1ddc866b6bb283fa5f1066cbf2b3465b0d52))
* **ci:** pass --repo to gh commands in publish job ([1aeea65](https://github.com/wherka-ama/gh-app-auth/commit/1aeea65ebca1025f3eaaabb1ff1cd3f6f28f7eb6))
* **ci:** pass --repo to gh commands in the publish job ([d80f4b9](https://github.com/wherka-ama/gh-app-auth/commit/d80f4b9afa08aadbb2c05b5fe1e24255cfbcf37f))

## [0.1.0](https://github.com/wherka-ama/gh-app-auth/compare/v0.0.16...v0.1.0) (2026-09-15)


### ⚠ BREAKING CHANGES

* **release:** the maintainer release runbook changes — releases are dispatched via `gh workflow run release.yml`, not created with `gh release create --prerelease`.
* **cli:** gh app-auth exec explicit selectors now fail closed: --repo must match the selected App's configured route and --installation-id must be a configured installation of that App. Invocations that relied on unconfigured repository hosts or installation IDs now return an error.

### Features

* add RPM/DEB package builds for amd64, arm64 ([#37](https://github.com/wherka-ama/gh-app-auth/issues/37)) ([8aab024](https://github.com/wherka-ama/gh-app-auth/commit/8aab0245ee6928dad96066dfe34ed487fdb9b411))
* auto-detect installation ID and add config command ([6734926](https://github.com/wherka-ama/gh-app-auth/commit/673492603b760fbc49d6b14f199129b93266d8f7))
* auto-detect mode for app installation id ([ca0b4ce](https://github.com/wherka-ama/gh-app-auth/commit/ca0b4ce131c05cc0faab093123c81cfe14f1b9c6))
* Automatic setup of GitHub App during git clone ([fd8ba89](https://github.com/wherka-ama/gh-app-auth/commit/fd8ba89d4f25d9bb638470d2a77219c708b63e2f))
* Automatic setup of GitHub App during git clone ([8f640ac](https://github.com/wherka-ama/gh-app-auth/commit/8f640ac472c890d23ab3bef70e3e262d57ad8850))
* **ci:** release-please automation + release environment hook ([dc055a2](https://github.com/wherka-ama/gh-app-auth/commit/dc055a269de941d9b3cd44865498854ffed1322e))
* **cli:** add token command ([#62](https://github.com/wherka-ama/gh-app-auth/issues/62)) ([544775e](https://github.com/wherka-ama/gh-app-auth/commit/544775ee31e7d5bfd208eb03779d26e3a34fb1e6))
* **cli:** support GitHub App Client ID ([#53](https://github.com/wherka-ama/gh-app-auth/issues/53)) ([59cb28a](https://github.com/wherka-ama/gh-app-auth/commit/59cb28aa6f14bf80fe6286a8d50c24feff85f9dd))
* enable useHttpPath for path-specific patterns ([3656d03](https://github.com/wherka-ama/gh-app-auth/commit/3656d03d59cc359c6c1fe1cc02460c0aa5dd97c8))
* GitHub CLI extension for GitHub App and PAT authentication ([6df981b](https://github.com/wherka-ama/gh-app-auth/commit/6df981b046cbfabe872061841385f92c37efd7a5))
* run commands with app authentication ([#51](https://github.com/wherka-ama/gh-app-auth/issues/51)) ([d959ae1](https://github.com/wherka-ama/gh-app-auth/commit/d959ae10591f6e2f00bca30b32c6a9b8f81e6afe))


### Bug Fixes

* **auth:** leave clock skew margin on JWT exp and iat claims ([#58](https://github.com/wherka-ama/gh-app-auth/issues/58)) ([620f73d](https://github.com/wherka-ama/gh-app-auth/commit/620f73d8e27a81ea5736acbf5643b461da61c0f4))
* **ci:** default dispatch ref to the triggering branch, not main ([ef6c239](https://github.com/wherka-ama/gh-app-auth/commit/ef6c23940694293f7af42f65f2ea0f779fd6f85a))
* **ci:** disable cgo for e2e test compilation ([6f76080](https://github.com/wherka-ama/gh-app-auth/commit/6f760808887b77409b6ecbbd3fea54a049a5f0af))
* **ci:** grant callee permissions on reusable workflow calls ([1bbadb9](https://github.com/wherka-ama/gh-app-auth/commit/1bbadb95343cb88732eee090970709a53c5d8ed1))
* **ci:** handle fork PRs and edge cases in workflows ([#30](https://github.com/wherka-ama/gh-app-auth/issues/30)) ([a3b8475](https://github.com/wherka-ama/gh-app-auth/commit/a3b84757eb650486c5508dea6c3c2ac2961f8146))
* **ci:** make auto-bump pre-major-aware on 0.x versions ([f40cdc5](https://github.com/wherka-ama/gh-app-auth/commit/f40cdc5a3538d3f40c68c456ec130ac80adeacec))
* **config:** split multi-pattern entries with unique installation IDs per org ([#35](https://github.com/wherka-ama/gh-app-auth/issues/35)) ([601632a](https://github.com/wherka-ama/gh-app-auth/commit/601632a86f197294ad294981dd30f3ff0e4b35f7)), closes [#34](https://github.com/wherka-ama/gh-app-auth/issues/34)
* create draft release before building assets ([8fb03f5](https://github.com/wherka-ama/gh-app-auth/commit/8fb03f5d324dc2f8268ef9ad1e79cb58858a1255))
* handle existing releases in release workflow ([76890af](https://github.com/wherka-ama/gh-app-auth/commit/76890af085b4982ac3fcbc1cd6e10c13ba1cc412))
* improving the NFPM flow + small tweaks on the next-version script ([43615b9](https://github.com/wherka-ama/gh-app-auth/commit/43615b903c22dd4fa0b4dcfbf4d4d637961ab2a9))
* reorder credential helpers for correct precedence ([16af226](https://github.com/wherka-ama/gh-app-auth/commit/16af226b0dddc38c6b4bc063047733b5dac1240e)), closes [#20](https://github.com/wherka-ama/gh-app-auth/issues/20)
* resolve git credential helper precedence for path-specific patterns ([5aa0f3f](https://github.com/wherka-ama/gh-app-auth/commit/5aa0f3f3d0881bdb92817590c1e536f0ef4583b8))
* restoring NFPM_CMD which was removed by mistake ([#57](https://github.com/wherka-ama/gh-app-auth/issues/57)) ([055ff64](https://github.com/wherka-ama/gh-app-auth/commit/055ff64a8995ac7021917b77e5c454ecf1ba9636))
* **security:** address CodeQL security alerts in diagnostic logging ([#28](https://github.com/wherka-ama/gh-app-auth/issues/28)) ([570e7b0](https://github.com/wherka-ama/gh-app-auth/commit/570e7b01ad041de9809dbf759e40d04a557be200))
* **security:** implement multi-layered sensitive data redaction ([#29](https://github.com/wherka-ama/gh-app-auth/issues/29)) ([9982156](https://github.com/wherka-ama/gh-app-auth/commit/99821564bc543f836df8b4743daf20b908c1a5c1))
* When both filesystem+env var are set, no app could be configured anymore ([60e7336](https://github.com/wherka-ama/gh-app-auth/commit/60e7336ebf0b4317497887e9de548709f2b5283a))
* When both filesystem+env var are set, no app could be configured anymore ([9a35515](https://github.com/wherka-ama/gh-app-auth/commit/9a355150253643bb4446c87ca8f6c0e0b2f98367))
* Windows test compatibility and security workflow ([f9af0be](https://github.com/wherka-ama/gh-app-auth/commit/f9af0becc8e5b75d595ad2d8c265d3ec68e1e54a))


### Continuous Integration

* **release:** draft-first gated release pipeline ([10522b1](https://github.com/wherka-ama/gh-app-auth/commit/10522b12f69d17c3f3b60a2f36210b3caa3cb661))

## [Unreleased]

### Added

- Add `gh app-auth token` for fresh, explicitly selected GitHub App installation tokens.
- Add `--client-id` selector to `gh app-auth exec`.
- Add a cross-platform E2E pipeline (`e2e.yml`) that installs and exercises the
  real release assets on Linux (deb/rpm, amd64+arm64), macOS (arm64+intel), and
  Windows (amd64+arm64), including cross-org submodule clones.
- Add `scripts/next-version.sh` and SLSA build-provenance attestation
  (`attest-release.yml`) to the release pipeline.

### Changed

- **Release process**: releases are now triggered by `workflow_dispatch` (or a
  `v*` tag push) instead of prerelease creation. The pipeline stages assets on a
  **draft** release, gates publication on the E2E suite and attestation, then
  publishes the draft as latest. This is compatible with GitHub's immutable
  releases, which lock assets at publish — the previous prerelease-promotion
  model would fail with `HTTP 422` once enabled. See `docs/RELEASE_PROCESS.md`.
  **Breaking change** for the maintainer release runbook: create releases via
  `gh workflow run release.yml`, not `gh release create --prerelease`.

- `gh app-auth exec` explicit selectors now fail closed: `--repo` must match the
  selected App's configured route and `--installation-id` must be a configured
  installation of that App. Previously a single App-ID match bypassed repository
  validation and `--installation-id` could override the configured installation.
  **Breaking change** for invocations that relied on those behaviours.

### Fixed

- `gh app-auth setup` rejects malformed App patterns instead of silently
  creating no configuration entry.

### Security

- Enforce configured repository routes and exact installation IDs for explicit App selectors.
- Bound authenticated redirects: refuse cross-host redirects and stop after 10 hops.

[Unreleased]: https://github.com/AmadeusITGroup/gh-app-auth/compare/v1.0.0...HEAD
