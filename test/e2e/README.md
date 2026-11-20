# End-to-end tests

This directory holds the live end-to-end (E2E) suite for gh-app-auth. It authenticates against real
private repositories in two dedicated test organizations, so it **needs real credentials** and is not
part of `make test`.

The maintained documentation is in `docs/`:

- [End-to-end testing tutorial](../../docs/contribute/e2e-testing.md): how to run the suite safely.
- [End-to-end test infrastructure](../../docs/contribute/e2e-infrastructure.md): the maintainer recipe for
  the fixture, secrets and workflows.
- [Testing guide](../../docs/contribute/testing.md): all test levels.

## Build tag and environment

Every Go file here has the `//go:build e2e` constraint, so it is compiled only with `-tags=e2e`:

```bash
go test -v -tags=e2e -run TestPreflight -timeout=2m ./test/e2e/...
```

| Variable | Required | Meaning |
| --- | --- | --- |
| `E2E_APP_ID` | Yes | GitHub App ID |
| `E2E_PRIVATE_KEY` | One of the two | Raw PEM; takes precedence |
| `E2E_PRIVATE_KEY_B64` | One of the two | Base64-encoded PEM |
| `E2E_BINARY_PATH` | No | Binary to test; built from source when unset |
| `E2E_TEST_ORG_1` | No | Defaults to `gh-app-auth-test-1` |
| `E2E_TEST_ORG_2` | No | Defaults to `gh-app-auth-test-2` |

The suite fails, rather than skips, when the credentials are missing. Never print these variables or
commit key material.

The [`scripts/`](scripts/README.md) directory holds older manual scripts. They are legacy and are not the
acceptance path.
