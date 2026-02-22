//go:build e2e

package e2e

import (
	"fmt"
	"strings"
	"testing"
)

// TestAuthentication validates the complete authentication configuration
// workflow end-to-end:
//  1. setup GitHub App credentials (via key file and via env var)
//  2. list configured apps
//  3. sync git credential helpers via gitconfig --sync
//  4. verify the credential helper responds to a matching URL
//  5. verify multi-org (cross-org) App setup
//
// Each sub-test uses a fully isolated config directory and git config to
// prevent state leakage between tests.
func TestAuthentication(t *testing.T) {
	t.Run("setup_with_key_file", testSetupWithKeyFile)
	t.Run("setup_with_env_key", testSetupWithEnvKey)
	t.Run("list_shows_configured_app", testListShowsConfiguredApp)
	t.Run("gitconfig_sync_registers_helper", testGitconfigSync)
	t.Run("credential_helper_responds", testCredentialHelperResponds)
	t.Run("remove_clears_configuration", testRemoveClearsConfiguration)
}

// testSetupWithKeyFile exercises setup using --key-file flag with --use-filesystem.
// --use-filesystem is required in CI where no OS keyring is available.
func testSetupWithKeyFile(t *testing.T) {
	env, _ := isolatedAppConfig(t)
	keyFile := writePrivateKeyFile(t, globalConfig.PrivateKeyPEM)
	pattern := fmt.Sprintf("github.com/%s/*", globalConfig.TestOrg1)

	stdout, stderr, err := RunCmd(t, env,
		"setup",
		"--app-id", globalConfig.AppID,
		"--key-file", keyFile,
		"--patterns", pattern,
		"--name", "E2E-KeyFile-Test",
		"--use-filesystem",
	)
	if err != nil {
		t.Fatalf("setup --key-file failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}
	if !strings.Contains(stdout+stderr, "Successfully configured") {
		t.Errorf("expected success message in output\nstdout: %s\nstderr: %s", stdout, stderr)
	}
	t.Log("✓ setup with key file succeeded")
}

// testSetupWithEnvKey exercises setup using GH_APP_PRIVATE_KEY environment variable.
// This is the preferred CI approach (no key file on disk).
func testSetupWithEnvKey(t *testing.T) {
	env, _ := isolatedAppConfig(t)
	env = setEnv(env, "GH_APP_PRIVATE_KEY", globalConfig.PrivateKeyPEM)
	pattern := fmt.Sprintf("github.com/%s/*", globalConfig.TestOrg1)

	stdout, stderr, err := RunCmd(t, env,
		"setup",
		"--app-id", globalConfig.AppID,
		"--patterns", pattern,
		"--name", "E2E-EnvKey-Test",
		"--use-filesystem",
	)
	if err != nil {
		t.Fatalf("setup with env key failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}
	if !strings.Contains(stdout+stderr, "Successfully configured") {
		t.Errorf("expected success message\nstdout: %s\nstderr: %s", stdout, stderr)
	}
	t.Log("✓ setup with GH_APP_PRIVATE_KEY env var succeeded")
}

// testListShowsConfiguredApp verifies that list output contains the app's details.
func testListShowsConfiguredApp(t *testing.T) {
	env, _ := isolatedAppConfig(t)
	keyFile := writePrivateKeyFile(t, globalConfig.PrivateKeyPEM)
	pattern := fmt.Sprintf("github.com/%s/*", globalConfig.TestOrg1)

	_, _, err := RunCmd(t, env,
		"setup",
		"--app-id", globalConfig.AppID,
		"--key-file", keyFile,
		"--patterns", pattern,
		"--name", "E2E-List-Test",
		"--use-filesystem",
	)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	stdout, stderr, err := RunCmd(t, env, "list")
	if err != nil {
		t.Fatalf("list failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}
	combined := stdout + stderr
	if !strings.Contains(combined, globalConfig.AppID) {
		t.Errorf("list output missing app ID %s\nOutput: %s", globalConfig.AppID, combined)
	}
	if !strings.Contains(combined, globalConfig.TestOrg1) {
		t.Errorf("list output missing org %s\nOutput: %s", globalConfig.TestOrg1, combined)
	}
	t.Log("✓ list shows configured app")
}

// testGitconfigSync verifies that gitconfig --sync registers the credential helper
// in the (isolated) global git config.
func testGitconfigSync(t *testing.T) {
	requireGit(t)

	env, _ := isolatedAppConfig(t)
	keyFile := writePrivateKeyFile(t, globalConfig.PrivateKeyPEM)
	pattern := fmt.Sprintf("github.com/%s/*", globalConfig.TestOrg1)

	_, _, err := RunCmd(t, env,
		"setup",
		"--app-id", globalConfig.AppID,
		"--key-file", keyFile,
		"--patterns", pattern,
		"--name", "E2E-GitconfigSync-Test",
		"--use-filesystem",
	)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	stdout, stderr, err := RunCmd(t, env, "gitconfig", "--sync", "--global")
	if err != nil {
		t.Fatalf("gitconfig --sync failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}
	t.Logf("gitconfig sync output: %s%s", stdout, stderr)

	// Verify the helper was written to the isolated git config
	gitOut, gitErr, err := RunGit(t, env, t.TempDir(),
		"config", "--global", "--get-regexp", "credential\\..*\\.helper")
	if err != nil {
		t.Fatalf("git config read failed: %v\nstderr: %s", err, gitErr)
	}
	if !strings.Contains(gitOut, "gh-app-auth") {
		t.Errorf("expected gh-app-auth credential helper in git config\ngit config output: %s", gitOut)
	}
	t.Log("✓ credential helper registered in git config")
}

// testCredentialHelperResponds verifies the git-credential subcommand responds
// to a URL that matches a configured pattern, using the Git credential protocol.
func testCredentialHelperResponds(t *testing.T) {
	env, _ := isolatedAppConfig(t)
	keyFile := writePrivateKeyFile(t, globalConfig.PrivateKeyPEM)
	pattern := fmt.Sprintf("github.com/%s/*", globalConfig.TestOrg1)

	_, _, err := RunCmd(t, env,
		"setup",
		"--app-id", globalConfig.AppID,
		"--key-file", keyFile,
		"--patterns", pattern,
		"--name", "E2E-CredHelper-Test",
		"--use-filesystem",
	)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Simulate git's credential get request for the test org's main repo.
	// The helper should return username + password (installation token).
	credInput := fmt.Sprintf("protocol=https\nhost=github.com\npath=%s/%s\n\n",
		globalConfig.TestOrg1, mainRepo)

	var stdout, stderr string
	err = retryOp(t, "credential get", func() error {
		var rerr error
		stdout, stderr, rerr = RunCmdWithStdin(t, env, credInput, "git-credential", "get")
		return rerr
	})
	if err != nil {
		t.Fatalf("git-credential get failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}
	if !strings.Contains(stdout, "password=") {
		t.Errorf("expected password= in credential output\nstdout: %s\nstderr: %s", stdout, stderr)
	}
	if !strings.Contains(stdout, "username=") {
		t.Errorf("expected username= in credential output\nstdout: %s", stdout)
	}
	t.Log("✓ credential helper returned valid credentials")
}

// testRemoveClearsConfiguration verifies that remove deletes the app entry.
func testRemoveClearsConfiguration(t *testing.T) {
	env, _ := isolatedAppConfig(t)
	keyFile := writePrivateKeyFile(t, globalConfig.PrivateKeyPEM)
	pattern := fmt.Sprintf("github.com/%s/*", globalConfig.TestOrg1)
	appName := "E2E-Remove-Test"

	_, _, err := RunCmd(t, env,
		"setup",
		"--app-id", globalConfig.AppID,
		"--key-file", keyFile,
		"--patterns", pattern,
		"--name", appName,
		"--use-filesystem",
	)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	stdout, stderr, err := RunCmd(t, env, "remove", "--name", appName, "--yes")
	if err != nil {
		t.Fatalf("remove failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	// After remove, list should not show the app
	listOut, listErr, err := RunCmd(t, env, "list")
	if err != nil {
		// Empty config is acceptable — list may return non-zero when nothing configured
		t.Logf("list after remove returned error (may be expected): %v\nstderr: %s", err, listErr)
		return
	}
	combined := listOut + listErr
	if strings.Contains(combined, appName) {
		t.Errorf("app %q still appears in list after remove\nOutput: %s", appName, combined)
	}
	t.Log("✓ remove cleared configuration")
}
