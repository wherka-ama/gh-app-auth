# gh-app-auth documentation

This directory holds the documentation for `gh-app-auth` v0.2.0, a GitHub CLI extension and Git
credential helper for GitHub Apps and personal access tokens. The files here are the single source
for the documentation site, which is built with Astro Starlight from the `website/` directory.

The site is not published yet. It goes live after the documentation workflow deploys an approved
`main`. Until then, read these files directly on GitHub.

## Start here

- [Overview and limitations](start-here/overview.md)
- [Install gh-app-auth](start-here/installation.md)
- [Make your first authenticated clone](start-here/first-clone.md)
- [Run your first command with exec](start-here/first-exec.md)
- [Guided tour](start-here/guided-tour.md)

## Guides

- [Guides overview](guides/index.md)
- [Create and install a GitHub App](guides/create-github-app.md)
- [Personal access tokens](guides/personal-access-tokens.md)
- [Git credential helper](guides/git-credential-helper.md)
- [Several organizations and submodules](guides/multiple-organizations.md)
- [GitHub Enterprise Server](guides/github-enterprise-server.md)
- [Credential lifecycle](guides/credential-lifecycle.md)
- [Continuous integration](guides/ci/index.md): [GitHub Actions](guides/ci/github-actions.md),
  [Jenkins](guides/ci/jenkins.md), [GitLab CI](guides/ci/gitlab-ci.md)

## Reference

- [Reference overview](reference/index.md)
- [CLI overview](reference/cli/index.md), with a generated page per command
- [Configuration](reference/configuration.md)
- [Environment variables](reference/environment-variables.md)
- [GitHub Actions](reference/github-actions.md)

## Understand

- [Identity and credentials](understand/identity-and-credentials.md)
- [How credential routing works](understand/credential-routing.mdx)
- [Tokens and caching](understand/tokens-and-caching.md)
- [Security model](understand/security-model.md)
- [Architecture](understand/architecture.mdx)
- [Origins](understand/origins.md)

## Troubleshoot

- [Troubleshooting](troubleshoot/index.md)
- [Diagnostics and logging](troubleshoot/diagnostics.md)
- [Known limitations](troubleshoot/known-limitations.md)

## Contribute

- [Contributing overview](contribute/index.md)
- [Development](contribute/development.md)
- [Testing guide](contribute/testing.md)
- [End-to-end testing tutorial](contribute/e2e-testing.md)
- [End-to-end test infrastructure](contribute/e2e-infrastructure.md)
- [Release runbook](contribute/releasing.md)
- [Documentation guide](contribute/documentation.md)
- [Architecture decisions](contribute/architecture-decisions.md)

## Older document names

These files remain as short pointers so that existing links keep working. Each names its current home.

| Older document | Current home |
| --- | --- |
| [architecture.md](architecture.md) | [Architecture](understand/architecture.mdx) |
| [ci-cd-guide.md](ci-cd-guide.md) | [Continuous integration](guides/ci/index.md) |
| [configuration.md](configuration.md) | [Configuration reference](reference/configuration.md) |
| [DIAGNOSTIC_LOGGING.md](DIAGNOSTIC_LOGGING.md) | [Diagnostics and logging](troubleshoot/diagnostics.md) |
| [E2E_INFRASTRUCTURE.md](E2E_INFRASTRUCTURE.md) | [End-to-end test infrastructure](contribute/e2e-infrastructure.md) |
| [E2E_TESTING_TUTORIAL.md](E2E_TESTING_TUTORIAL.md) | [End-to-end testing tutorial](contribute/e2e-testing.md) |
| [ENCRYPTED_STORAGE_ARCHITECTURE.md](ENCRYPTED_STORAGE_ARCHITECTURE.md) | Historical proposal; current model in [Security model](understand/security-model.md) |
| [GITCONFIG_COMMAND.md](GITCONFIG_COMMAND.md) | [Git credential helper](guides/git-credential-helper.md) |
| [installation.md](installation.md) | [Install gh-app-auth](start-here/installation.md) |
| [JSON_MULTI_ORG_CONFIG.md](JSON_MULTI_ORG_CONFIG.md) | [GitHub Actions reference](reference/github-actions.md) |
| [origin_of_the_project.md](origin_of_the_project.md) | [Origins](understand/origins.md) |
| [PATTERN_ROUTING.md](PATTERN_ROUTING.md) | [How credential routing works](understand/credential-routing.mdx) |
| [presentation.md](presentation.md) | Optional slide deck; see the [guided tour](start-here/guided-tour.md) |
| [RELEASE_PROCESS.md](RELEASE_PROCESS.md) | [Release runbook](contribute/releasing.md) |
| [security.md](security.md) | [Security model](understand/security-model.md) |
| [TESTING.md](TESTING.md) | [Testing guide](contribute/testing.md) |
| [TOKEN_CACHING.md](TOKEN_CACHING.md) | [Tokens and caching](understand/tokens-and-caching.md) |
| [troubleshooting.md](troubleshooting.md) | [Troubleshooting](troubleshoot/index.md) |
