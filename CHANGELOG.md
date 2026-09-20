# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0](https://github.com/wherka-ama/gh-app-auth/compare/v0.0.17...v0.1.0) (2026-09-20)


### ⚠ BREAKING CHANGES

* **cli:** gh app-auth exec explicit selectors now fail closed: --repo must match the selected App's configured route and --installation-id must be a configured installation of that App. Invocations that relied on unconfigured repository hosts or installation IDs now return an error.

### Features

* **cli:** add token command ([#62](https://github.com/wherka-ama/gh-app-auth/issues/62)) ([544775e](https://github.com/wherka-ama/gh-app-auth/commit/544775ee31e7d5bfd208eb03779d26e3a34fb1e6))
* **cli:** support GitHub App Client ID ([#53](https://github.com/wherka-ama/gh-app-auth/issues/53)) ([59cb28a](https://github.com/wherka-ama/gh-app-auth/commit/59cb28aa6f14bf80fe6286a8d50c24feff85f9dd))


### Bug Fixes

* **auth:** leave clock skew margin on JWT exp and iat claims ([#58](https://github.com/wherka-ama/gh-app-auth/issues/58)) ([620f73d](https://github.com/wherka-ama/gh-app-auth/commit/620f73d8e27a81ea5736acbf5643b461da61c0f4))
* restoring NFPM_CMD which was removed by mistake ([#57](https://github.com/wherka-ama/gh-app-auth/issues/57)) ([055ff64](https://github.com/wherka-ama/gh-app-auth/commit/055ff64a8995ac7021917b77e5c454ecf1ba9636))

## [Unreleased]

### Added

- Add `gh app-auth token` for fresh, explicitly selected GitHub App installation tokens.
- Add `--client-id` selector to `gh app-auth exec`.

### Changed

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
