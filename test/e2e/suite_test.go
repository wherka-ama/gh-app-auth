//go:build e2e

package e2e

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"testing"
)

// E2EConfig holds all configuration needed for E2E tests.
// All fields are populated from environment variables at startup.
type E2EConfig struct {
	// BinaryPath is the path to the gh-app-auth binary under test.
	// Set via E2E_BINARY_PATH. If empty, binary is built from source.
	BinaryPath string

	// AppID is the GitHub App ID (string form for passing to CLI flags).
	// Set via E2E_APP_ID.
	AppID string

	// PrivateKeyPEM is the decoded PEM-encoded private key content.
	// Set via E2E_PRIVATE_KEY_B64 (base64-encoded for safe CI storage).
	PrivateKeyPEM string

	// TestOrg1 is the first test organization name.
	// Set via E2E_TEST_ORG_1 (default: gh-app-auth-test-1).
	TestOrg1 string

	// TestOrg2 is the second test organization name.
	// Set via E2E_TEST_ORG_2 (default: gh-app-auth-test-2).
	TestOrg2 string

	// GitHubToken is a GitHub token for API calls (pre-flight checks, repo verification).
	// Must have 'repo' scope to read private repo metadata.
	// Set via E2E_GITHUB_TOKEN.
	GitHubToken string
}

// globalConfig is loaded once in TestMain and shared across all tests.
var globalConfig *E2EConfig

// TestMain is the entry point for the E2E test suite.
// It loads configuration and fails fast if required env vars are missing.
func TestMain(m *testing.M) {
	cfg, err := loadE2EConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "E2E configuration error: %v\n\n", err)
		fmt.Fprintf(os.Stderr, "Required environment variables:\n")
		fmt.Fprintf(os.Stderr, "  E2E_APP_ID            GitHub App ID\n")
		fmt.Fprintf(os.Stderr, "  E2E_PRIVATE_KEY_B64   Base64-encoded private key PEM\n")
		fmt.Fprintf(os.Stderr, "  E2E_GITHUB_TOKEN      GitHub token (repo scope) for API calls\n")
		fmt.Fprintf(os.Stderr, "\nOptional:\n")
		fmt.Fprintf(os.Stderr, "  E2E_BINARY_PATH       Path to binary (built from source if unset)\n")
		fmt.Fprintf(os.Stderr, "  E2E_TEST_ORG_1        Test org 1 (default: gh-app-auth-test-1)\n")
		fmt.Fprintf(os.Stderr, "  E2E_TEST_ORG_2        Test org 2 (default: gh-app-auth-test-2)\n")
		fmt.Fprintf(os.Stderr, "\nSee docs/E2E_INFRASTRUCTURE.md for infrastructure setup.\n")
		os.Exit(1)
	}

	globalConfig = cfg
	os.Exit(m.Run())
}

// loadE2EConfig loads the E2E test configuration from environment variables.
func loadE2EConfig() (*E2EConfig, error) {
	cfg := &E2EConfig{}

	cfg.AppID = os.Getenv("E2E_APP_ID")
	if cfg.AppID == "" {
		return nil, fmt.Errorf("E2E_APP_ID is required")
	}

	privateKeyB64 := os.Getenv("E2E_PRIVATE_KEY_B64")
	if privateKeyB64 == "" {
		return nil, fmt.Errorf("E2E_PRIVATE_KEY_B64 is required")
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(privateKeyB64))
	if err != nil {
		return nil, fmt.Errorf("failed to base64-decode E2E_PRIVATE_KEY_B64: %w", err)
	}
	cfg.PrivateKeyPEM = string(decoded)

	cfg.GitHubToken = os.Getenv("E2E_GITHUB_TOKEN")
	if cfg.GitHubToken == "" {
		return nil, fmt.Errorf("E2E_GITHUB_TOKEN is required")
	}

	cfg.TestOrg1 = os.Getenv("E2E_TEST_ORG_1")
	if cfg.TestOrg1 == "" {
		cfg.TestOrg1 = "gh-app-auth-test-1"
	}

	cfg.TestOrg2 = os.Getenv("E2E_TEST_ORG_2")
	if cfg.TestOrg2 == "" {
		cfg.TestOrg2 = "gh-app-auth-test-2"
	}

	cfg.BinaryPath = os.Getenv("E2E_BINARY_PATH")

	return cfg, nil
}
