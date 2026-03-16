# Development guide

## Prerequisites

- Go 1.21+
- GitHub CLI (`gh`)
- Git

## Build

```bash
git clone https://github.com/AmadeusITGroup/gh-app-auth.git
cd gh-app-auth
go build -o gh-app-auth .
```

Install as a local extension for manual testing:

```bash
gh extension install .
```

## Test

```bash
make test                      # all tests
go test ./pkg/auth/...         # single package
go test -v -run TestFunction . # single test
go test -race ./...            # with race detector
go test -coverprofile=c.out ./... && go tool cover -func=c.out  # coverage
```

## Lint and CI

```bash
make fmt           # format code
make lint          # golangci-lint
make quality       # lint + vet
make security-scan # gosec + govulncheck
make markdownlint  # markdown files
make ci            # full pipeline (build, test, lint, security, docs)
```

Run `make ci` before submitting a PR. The same checks run in GitHub Actions.

## Project structure

See [architecture.md](architecture.md) for package layout and dependency rules.

## Code conventions

### Imports

Three groups separated by blank lines: stdlib, external, internal.

```go
import (
    "fmt"
    "os"

    "github.com/spf13/cobra"

    "github.com/AmadeusITGroup/gh-app-auth/pkg/config"
)
```

### Error handling

Wrap errors with context:

```go
if err != nil {
    return fmt.Errorf("failed to load config from %s: %w", path, err)
}
```

### Console output

Use `fmt.Println` for static strings, `fmt.Printf` only when formatting
variables:

```go
fmt.Println("Operation complete")           // static
fmt.Printf("Processed %d items\n", count)   // dynamic
```

### Commands

Follow the existing Cobra pattern:

```go
func NewExampleCmd() *cobra.Command {
    var flag string
    cmd := &cobra.Command{
        Use:   "example",
        Short: "One-line description",
        RunE: func(cmd *cobra.Command, args []string) error {
            return exampleRun(flag)
        },
    }
    cmd.Flags().StringVar(&flag, "flag", "", "Description")
    return cmd
}
```

Register in `cmd/root.go`.

### Tests

Use table-driven tests and `t.TempDir()` for file isolation:

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid", "hello", "HELLO", false},
        {"empty", "", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Security

- Never log tokens or private keys in plain text.
- Use `secrets.HashToken()` when logging token references.
- Validate file permissions (600/400) before reading key files.
- Use `t.TempDir()` and `t.Setenv()` in tests to avoid side effects.

## Commit messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`,
`ci`, `chore`

Scopes: `auth`, `config`, `cli`, `cache`, `security`, `docs`, `ci`, `deps`

Examples:

```
feat(auth): add JWT token caching
fix(config): handle missing config file gracefully
test(auth): add integration tests for token refresh
```

## Pull request checklist

- [ ] `make ci` passes
- [ ] New code has tests
- [ ] No secrets in code or logs
- [ ] Documentation updated if user-facing behavior changed
- [ ] Commit messages follow conventional commits

## Adding a new command

1. Create `cmd/newcommand.go` with `NewXxxCmd()` function.
2. Create `cmd/newcommand_test.go`.
3. Register in `cmd/root.go`.
4. Run `make ci`.

## Useful links

- [GitHub CLI extension docs](https://docs.github.com/en/github-cli/github-cli/creating-github-cli-extensions)
- [go-gh library](https://pkg.go.dev/github.com/cli/go-gh/v2)
- [GitHub App authentication](https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/about-authentication-with-a-github-app)
- [Git credential helpers](https://git-scm.com/docs/gitcredentials)
