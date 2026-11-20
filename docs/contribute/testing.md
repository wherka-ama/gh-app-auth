---
title: Testing guide
description: The test levels in gh-app-auth, how to run them, how tests are isolated, and where live end-to-end testing is documented.
sidebar:
  order: 3
---

gh-app-auth has three kinds of automated tests, and one legacy script that is not part of any gate.
This page describes how they are organized, not their current counts or pass rates; run the commands to
see those.

## Test levels

| Level | Where | Needs credentials? |
| --- | --- | --- |
| Unit tests | `*_test.go` next to the code in `cmd/` and `pkg/` | No |
| Integration tests | [`test/integration`](../../test/integration), using helpers in [`test/testutil`](../../test/testutil) | No |
| Live end-to-end tests | [`test/e2e`](../../test/e2e), behind the `e2e` build tag | **Yes**: a real GitHub App and test organizations |

Unit and integration tests use temporary directories, fake configurations, static test fixtures and a
mock GitHub server (`MockGitHubServer` in `test/testutil`) instead of real accounts. They never need real
credentials.

## Running the default tests

```bash
make test
go test -race ./...
go test -v ./pkg/auth/...
go test -v -run TestName ./pkg/...
```

`make test` runs `go test ./...`. Every file in `test/e2e` carries `//go:build e2e`, so the default
`go test ./...` and `make test` **do not compile or run the live end-to-end suite**. The CI test job
runs `go test -race` with coverage over `./...`, and `make quality` adds the lint, coverage threshold and
security scans. See [development](development.md) for the everyday commands.

## Writing unit tests

- Use table-driven tests with cases for valid input, edge cases and errors.
- Isolate the filesystem and environment with `t.TempDir()` and `t.Setenv()`. Point
  `GH_APP_AUTH_CONFIG` at a file in the temporary directory so a test cannot read your real configuration
  file, and also point `HOME` at a temporary directory so default locations, such as the secrets
  directory under `HOME`, stay inside the test. `GH_APP_AUTH_CONFIG` alone does not isolate stored
  secrets, and neither override isolates the **OS keyring**: a test that stores or reads secrets must
  use a mock secrets or keyring backend and must never touch real entries.
- Build configuration with the helpers in `test/testutil`, such as `ConfigBuilder` and
  `CreateTempConfig`. `CreateTestPrivateKey` writes a **static legacy fixture**; it does not generate a
  key. New tests that need a key should generate a throwaway key at runtime rather than extend the
  embedded fixtures.
- Mock external services through interfaces or the mock server, and never call the real GitHub API from a
  default test.
- Never log, print or commit tokens or keys, and never commit test credentials, even expired ones.
- Keep coverage at or above the thresholds enforced by `make quality`.

```go
func TestFunction(t *testing.T) {
    home := t.TempDir()
    t.Setenv("HOME", home)
    t.Setenv("GH_APP_AUTH_CONFIG", filepath.Join(home, "config.yml"))
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "example", "example", false},
        {"empty input", "", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Fatalf("Function() error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("Function() = %q, want %q", got, tt.want)
            }
        })
    }
}
```

## Integration tests

The integration tests in `test/integration` exercise the Git credential helper with the credential
protocol's input and output: no configuration, no matching App, a host-only query, `store` and
`erase`, an invalid operation, multi-stage protocols and the output format. Each test points
`GH_APP_AUTH_CONFIG` at a temporary directory and none contacts GitHub.

For the host-only case, the helper returns empty output when Git provides a host without a path. Do not
assume Git then asks again with the full path. `gitconfig --sync` writes `useHttpPath` for path-specific
patterns, which is what makes Git send the repository path in the first place. See the
[Git credential helper guide](../guides/git-credential-helper.md).

## Diagnosing a failing test

- Run one package or test with `-v` and `-run`, and add `-race` when a failure looks timing-related.
- For helper behavior, enable the [diagnostic log](../troubleshoot/diagnostics.md) with an absolute
  `GH_APP_AUTH_DEBUG_LOG` path and a **test** configuration, and read the flow steps. Do not capture Git or
  curl traces together with real credentials.
- A test that needs a built binary can use `testutil.BuildBinary` instead of relying on one you built
  earlier.

## Live end-to-end tests

The live suite authenticates against real private repositories in two dedicated test organizations. It
is documented, with its safety rules, in the [end-to-end testing tutorial](e2e-testing.md), and its
fixture provisioning is in [end-to-end test infrastructure](e2e-infrastructure.md).

## Legacy manual script

`test/manual/test-git-credential.sh` predates the current suite and is **not a gate and not
recommended**. It unsets `GH_APP_AUTH_CONFIG`, so it may fall back to your actual configuration, and it
captures credential output, which can contain a real token. There is no execution recipe for it here.
If you must read it, treat it as a live-secret risk and do not run it with real keys.

## Best practices

- Test the security-sensitive paths first: key handling, permission checks, token handling and logging.
- Keep tests deterministic and independent; no test should rely on another's state.
- Prefer small tests that fail with a precise message over large end-to-end scripts.
- Add the test with the fix: reproduce a bug, watch the test fail, then fix it.
