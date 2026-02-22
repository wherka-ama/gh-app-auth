//go:build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// Expected test repositories — must be pre-created per docs/E2E_INFRASTRUCTURE.md.
const (
	mainRepo      = "main-repo"      // In TestOrg1: private, contains submodule from TestOrg2
	submoduleRepo = "submodule-repo" // In TestOrg2: private, referenced as submodule

	// Marker files must exist in cloned repos for content-verification assertions.
	mainMarkerFile         = "data/main-marker.txt"
	submoduleMarkerFile    = "data/submodule-marker.txt"
	mainMarkerContent      = "gh-app-auth-e2e-main-marker"
	submoduleMarkerContent = "gh-app-auth-e2e-submodule-marker"
)

// TestPreflight validates that the E2E infrastructure is correctly set up.
// This is a hard pre-flight gate: all sub-tests must pass before any workflow
// tests can be meaningful. Failures here indicate infrastructure issues, not
// product bugs.
func TestPreflight(t *testing.T) {
	t.Run("binary_accessible", testPreflightBinary)
	t.Run("github_token_valid", testPreflightGitHubToken)
	t.Run("org1_main_repo_is_private", testPreflightRepoPrivate(globalConfig.TestOrg1, mainRepo))
	t.Run("org2_submodule_repo_is_private", testPreflightRepoPrivate(globalConfig.TestOrg2, submoduleRepo))
}

func testPreflightBinary(t *testing.T) {
	env, _ := isolatedAppConfig(t)
	stdout, stderr, err := RunCmd(t, env, "--version")
	if err != nil {
		t.Fatalf("binary not accessible or returned error: %v\nstderr: %s", err, stderr)
	}
	combined := stdout + stderr
	if combined == "" {
		t.Error("expected version output, got nothing")
	}
	t.Logf("binary version: %s", combined)
}

func testPreflightGitHubToken(t *testing.T) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+globalConfig.GitHubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to call GitHub API: %v", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("E2E_GITHUB_TOKEN is invalid or lacks 'repo' scope (HTTP %d): %s", resp.StatusCode, body)
	}

	var user struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatalf("failed to parse user response: %v", err)
	}
	t.Logf("GitHub token valid (user: %s)", user.Login)
}

// testPreflightRepoPrivate returns a test function that verifies a repo exists
// and is private. Fails fast with a clear, actionable error if not.
func testPreflightRepoPrivate(org, repo string) func(*testing.T) {
	return func(t *testing.T) {
		info, err := fetchRepoInfo(globalConfig.GitHubToken, org, repo)
		if err != nil {
			t.Fatalf(
				"repo %s/%s not accessible: %v\n\nAction required: create the repository per docs/E2E_INFRASTRUCTURE.md",
				org, repo, err,
			)
		}
		if !info.Private {
			t.Fatalf(
				"SECURITY GATE FAILED: %s/%s is PUBLIC\n"+
					"E2E tests authenticate against private repos to validate real auth flows.\n"+
					"Set the repository to private before running E2E tests.",
				org, repo,
			)
		}
		t.Logf("✓ %s/%s exists and is private", org, repo)
	}
}

// repoMeta holds the relevant fields from the GitHub repos API response.
type repoMeta struct {
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
}

// fetchRepoInfo calls the GitHub repos API and returns repository metadata.
func fetchRepoInfo(token, org, repo string) (*repoMeta, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", org, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("repository %s/%s not found (404)", org, repo)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API HTTP %d: %s", resp.StatusCode, body)
	}

	var meta repoMeta
	if err := json.Unmarshal(body, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &meta, nil
}
