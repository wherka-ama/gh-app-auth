# CI/CD integration

This page shows how to use gh-app-auth in automated pipelines. The examples
assume the extension is already installed and a GitHub App or PAT has been
created.

## GitHub Actions

### Using the reusable action (recommended)

The repository includes a reusable action that handles installation,
configuration, and cleanup.

**Auto-detection mode** (recommended for single-org Apps):

```yaml
steps:
  - name: Setup GitHub App Auth
    uses: AmadeusITGroup/gh-app-auth/.github/actions/setup-gh-app-auth@main
    with:
      app-id: ${{ secrets.GH_APP_ID }}
      private-key: ${{ secrets.GH_APP_PRIVATE_KEY }}
      # patterns is optional - auto-detects from github.repository_owner

  - name: Clone
    run: git clone --recurse-submodules https://github.com/${{ github.repository }}
```

The action auto-detects the installation ID and uses `github.repository_owner`
to set the prefix dynamically.

**Explicit prefix** (for multi-org or specific patterns):

```yaml
steps:
  - name: Setup GitHub App Auth
    uses: AmadeusITGroup/gh-app-auth/.github/actions/setup-gh-app-auth@main
    with:
      app-id: ${{ secrets.GH_APP_ID }}
      private-key: ${{ secrets.GH_APP_PRIVATE_KEY }}
      patterns: 'github.com/myorg/'
      cleanup-on-exit: 'true'
```

`cleanup-on-exit: 'true'` (default) removes credentials when the job finishes.
Use this on self-hosted runners to avoid leaking credentials between jobs.

See [`.github/actions/README.md`](../.github/actions/README.md) for all action inputs.

### Manual setup

For more control, install and configure manually:

```yaml
steps:
  - name: Install extension
    run: gh extension install AmadeusITGroup/gh-app-auth

  - name: Configure credential
    env:
      GH_APP_PRIVATE_KEY: ${{ secrets.GH_APP_PRIVATE_KEY }}
    run: |
      gh app-auth setup \
        --app-id ${{ secrets.GH_APP_ID }} \
        --patterns "github.com/${{ github.repository_owner }}/"
      gh app-auth gitconfig --sync --global

  - name: Clone with submodules
    run: git clone --recurse-submodules https://github.com/${{ github.repository }}.git
```

The `setup` command auto-detects the installation ID. If your App is installed
on multiple organizations, specify `--patterns` for each org you need.

### Multiple organizations

Configure one credential per org. The helper picks the right one based on URL:

```yaml
- name: Configure orgs
  env:
    ORG1_KEY: ${{ secrets.ORG1_PRIVATE_KEY }}
    ORG2_KEY: ${{ secrets.ORG2_PRIVATE_KEY }}
  run: |
    GH_APP_PRIVATE_KEY="$ORG1_KEY" gh app-auth setup \
      --app-id ${{ secrets.ORG1_APP_ID }} \
      --patterns "github.com/org1/"

    GH_APP_PRIVATE_KEY="$ORG2_KEY" gh app-auth setup \
      --app-id ${{ secrets.ORG2_APP_ID }} \
      --patterns "github.com/org2/"

    gh app-auth gitconfig --sync --global
```

### Adding a PAT alongside Apps

```yaml
- name: Configure Bitbucket PAT
  run: |
    gh app-auth setup \
      --pat "${{ secrets.BB_PAT }}" \
      --patterns "bitbucket.example.com/" \
      --username "${{ secrets.BB_USERNAME }}" \
      --name "Bitbucket"
    gh app-auth gitconfig --sync --global
```

## Jenkins

```groovy
pipeline {
    agent any

    environment {
        GITHUB_APP_ID = credentials('github-app-id')
        GH_APP_PRIVATE_KEY = credentials('github-app-private-key')
    }

    stages {
        stage('Setup') {
            steps {
                sh '''
                    gh extension install AmadeusITGroup/gh-app-auth || true
                    gh app-auth setup \
                        --app-id "$GITHUB_APP_ID" \
                        --patterns "github.com/myorg/"
                    gh app-auth gitconfig --sync --global
                '''
            }
        }

        stage('Build') {
            steps {
                sh 'git clone --recurse-submodules https://github.com/myorg/repo.git'
            }
        }
    }

    post {
        always {
            sh '''
                gh app-auth gitconfig --clean --global || true
                gh app-auth remove --all --force || true
            '''
        }
    }
}
```

## Token expiry

GitHub App installation tokens expire after 1 hour. The extension caches tokens
for 55 minutes and generates a new one when the cached token is close to expiry.
This happens transparently — long-running jobs do not need special handling.

PATs do not expire through the extension (they follow whatever lifetime GitHub
or your provider sets).

## GitHub Enterprise

```bash
gh app-auth setup \
  --app-id 123456 \
  --key-file enterprise-app.pem \
  --patterns "github.example.com/corp/"
gh app-auth gitconfig --sync --global
```

## Cleanup on self-hosted runners

Non-ephemeral runners persist state between jobs. Always clean up:

```bash
gh app-auth gitconfig --clean --global
gh app-auth remove --all --force
gh app-auth remove --all-pats --force
```

Or use the composite action with `cleanup-on-exit: 'true'`.
