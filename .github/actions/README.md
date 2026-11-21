# GitHub Actions for gh-app-auth

This directory holds two composite actions for setting up and removing GitHub App authentication in a
workflow: `setup-gh-app-auth` and `cleanup-gh-app-auth`. They are conveniences with real limitations, and
the documentation recommends explicit, isolated workflow steps instead.

> **Read the limits before you use them.**
>
> - `setup-gh-app-auth` registers its "cleanup hook" as a shell `trap` inside the setup step. It is not a
>   post-job hook. The trap runs as soon as that step's shell exits, so it can **delete the credentials the
>   step just configured before later steps run**. If you use the composite action and manage cleanup
>   yourself, set `cleanup-on-exit: 'false'` and add the cleanup action with `if: always()`.
> - Inputs, including secrets, are interpolated into shell scripts. A secret that contains quotes or shell
>   syntax can break the step or run unintended commands, so "use a trusted secret" is not a guarantee of
>   safe interpolation. The JSON form (`apps-config`) places private keys directly in that text.
> - It does not use `--use-filesystem --key-file`. Where no OS keyring is available it depends on an
>   automatic filesystem fallback that is known to be unreliable.
> - `cleanup-gh-app-auth` removes all configured Apps and the Git helper entries, not PATs.
> - Do not use either action on shared or long-lived machines.

## Recommended instead

The [GitHub Actions guide](../../docs/guides/ci/github-actions.md) shows a complete workflow that avoids
the composite action: it isolates `HOME`, writes the key to a job-owned file, passes secrets through `env`
and always cleans up. The shared rules for Jenkins and GitLab CI are in the
[CI overview](../../docs/guides/ci/index.md).

## Available Actions

The exact inputs, outputs, defaults and limits of both actions are in the
[GitHub Actions reference](../../docs/reference/github-actions.md). It is checked against each
`action.yml`, so it stays accurate when they change. Reference an action by a release tag or commit, not
by `main`.

## Multi-Organization Examples

For several organizations, prefer one explicit credential per organization as described in the
[multiple organizations guide](../../docs/guides/multiple-organizations.md). The `apps-config` input of
the setup action is documented in the reference, together with the interpolation caution above.

## Best Practices

- Use ephemeral runners, secrets from the repository's secret store, and the least permissions the App
  needs.
- Pin the extension version and the action ref.
- Always clean up explicitly, and know that cleanup is best effort.

## Security Considerations

See the [security model](../../docs/understand/security-model.md) and
[known limitations](../../docs/troubleshoot/known-limitations.md).

## Troubleshooting

See [troubleshooting](../../docs/troubleshoot/index.md).

## Related Documentation

- [Main README](../../README.md)
- [Continuous integration guides](../../docs/guides/ci/index.md)
- [GitHub Actions reference](../../docs/reference/github-actions.md)
- [Architecture](../../docs/understand/architecture.mdx)
