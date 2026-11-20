---
title: Jenkins
description: Requirements and an untested pipeline sketch for using gh-app-auth on disposable Jenkins agents.
sidebar:
  order: 3
---

This project has no automated test against Jenkins. The requirements below follow from how gh-app-auth
works; the pipeline is a sketch to adapt and validate in your own environment.

## Agent requirements

- A **disposable** agent: a fresh container or virtual machine per build. Do not use a shared,
  long-lived agent.
- `gh`, `git`, a POSIX `sh` and outbound HTTPS access to GitHub (or your GitHub Enterprise Server).
- The extension installed inside the build, pinned to a release, with a job-owned `HOME` and
  `GH_CONFIG_DIR`, as in the [shared POSIX recipe](index.md#shared-posix-job-recipe).
- Secrets stored as Jenkins credentials: the App's Client ID, its private key and a GitHub CLI
  token.

## Pipeline sketch

Bind the credentials with `withCredentials`, and put the shell logic in a **single-quoted** Groovy
string so Jenkins never interpolates a secret into script text:

```groovy
pipeline {
  agent { label 'disposable-gh-app-auth' }
  stages {
    stage('Clone private repository') {
      steps {
        withCredentials([
          string(credentialsId: 'gh-app-client-id', variable: 'APP_CLIENT_ID'),
          string(credentialsId: 'gh-app-private-key', variable: 'APP_PRIVATE_KEY'),
          string(credentialsId: 'gh-cli-token', variable: 'GH_TOKEN')
        ]) {
          sh '''
            export APP_PATTERN='github.com/example-org/'
            sh ./ci/gh-app-auth-job.sh
          '''
        }
      }
    }
  }
}
```

`ci/gh-app-auth-job.sh` is your checked-in copy of the shared recipe followed by your Git commands.
Credential IDs, label and pattern are placeholders.

## Cautions

- Do not use double-quoted Groovy strings around secrets.
- Do not echo environment variables or enable `set -x` while secrets are bound; Jenkins masking is
  best effort.
- Never mutate the agent's real home directory, and never run a broad credential wipe.
- Cleanup runs in the recipe's `EXIT` trap. It is best effort; disposable agents are the guarantee.
