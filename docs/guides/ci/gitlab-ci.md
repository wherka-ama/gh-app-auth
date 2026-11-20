---
title: GitLab CI
description: Requirements and an untested job sketch for using gh-app-auth on disposable GitLab runners.
sidebar:
  order: 4
---

This project has no automated test against GitLab CI. The requirements below follow from how
gh-app-auth works; the job is a sketch to adapt and validate in your own environment. No container
image is named because this project does not publish or verify one.

## Runner requirements

- A **disposable** executor, such as a Docker or Kubernetes executor that starts a clean container
  per job. Do not use a shared shell executor.
- The job environment provides `gh`, `git` and a POSIX `sh`, and can reach GitHub over HTTPS.
- The extension is installed during the job, pinned to a release, with a job-owned `HOME` and
  `GH_CONFIG_DIR`, as in the [shared POSIX recipe](index.md#shared-posix-job-recipe).
- `APP_CLIENT_ID` and `GH_TOKEN` are **masked, protected** variables of type *Variable*.
- The private key is a **protected File variable** named `APP_PRIVATE_KEY_FILE`. A multi-line PEM cannot be
  masked, so do not store it as a masked Variable. For a File variable the value your job sees is the
  **path** of a file that GitLab writes for the job. Restrict it to protected branches and tags, and see
  GitLab's documentation on [SSH keys in jobs](https://docs.gitlab.com/ci/jobs/ssh_keys/) and
  [CI/CD variables](https://docs.gitlab.com/ci/variables/).
- Because the key cannot be masked, it is **not hidden from job logs**. Never `cat` it to the log, and
  never enable shell tracing or CI debug tracing in a job that can read it.

## Job sketch

```yaml
clone-private-repository:
  stage: build
  tags:
    - disposable-gh-app-auth
  variables:
    APP_PATTERN: 'github.com/example-org/'
  script:
    - sh ./ci/gh-app-auth-job.sh
```

`ci/gh-app-auth-job.sh` is your checked-in copy of the shared recipe followed by your Git commands.
The recipe copies the contents of the File variable's path, with output redirected into a file
owned by the job (`cat -- "$APP_PRIVATE_KEY_FILE" > "$key_file"`), so nothing is printed. The job exports
no raw key text and interpolates no key into a script line. Tag and pattern are placeholders.

## Cautions

- Do not paste secret values into `script:` lines or `variables:` blocks; reference them as
  environment variables inside the script file.
- Do not print variables or enable `set -x` while secrets are available.
- Protect the variables and limit which branches and tags can read them; mask the ones that can be masked.
- Cleanup runs in the recipe's `EXIT` trap and is best effort; a disposable executor is the
  guarantee.
