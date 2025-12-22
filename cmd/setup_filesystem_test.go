package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/AmadeusITGroup/gh-app-auth/pkg/config"
	"github.com/AmadeusITGroup/gh-app-auth/pkg/secrets"
)

// TestConfigureAppStorage_FilesystemWithEnvVar tests the scenario where
// --use-filesystem is specified but the key comes from GH_APP_PRIVATE_KEY env var
func TestConfigureAppStorage_FilesystemWithEnvVar(t *testing.T) {
	t.Run("filesystem storage requires key file path", func(t *testing.T) {
		app := &config.GitHubApp{
			AppID:          123456,
			Name:           "test-app",
			InstallationID: 789,
			Patterns:       []string{"github.com/org/*"},
			Priority:       1,
		}

		// Simulate the case where key comes from env var (no file path)
		privateKeyContent := "-----BEGIN RSA PRIVATE KEY-----\ntest\n-----END RSA PRIVATE KEY-----"
		expandedKeyFile := "" // Empty because key came from env var
		useKeyring := false   // User specified --use-filesystem

		_, err := configureAppStorage(app, privateKeyContent, expandedKeyFile, useKeyring)
		if err == nil {
			t.Fatal("Expected error when using --use-filesystem with GH_APP_PRIVATE_KEY env var")
		}

		expectedErrMsg := "filesystem storage requires --key-file"
		if err.Error() != expectedErrMsg {
			t.Errorf("Expected error message %q, got %q", expectedErrMsg, err.Error())
		}
	})

	t.Run("filesystem storage works with key file", func(t *testing.T) {
		// Skip on Windows - file permissions are not enforced the same way
		if runtime.GOOS == "windows" {
			t.Skip("Skipping on Windows: file permissions not supported")
		}

		tempDir := t.TempDir()
		keyPath := filepath.Join(tempDir, "test.pem")
		keyContent := "-----BEGIN RSA PRIVATE KEY-----\ntest\n-----END RSA PRIVATE KEY-----"
		if err := os.WriteFile(keyPath, []byte(keyContent), 0600); err != nil {
			t.Fatalf("Failed to create test key file: %v", err)
		}

		app := &config.GitHubApp{
			AppID:          123456,
			Name:           "test-app",
			InstallationID: 789,
			Patterns:       []string{"github.com/org/*"},
			Priority:       1,
		}

		expandedKeyFile := keyPath
		useKeyring := false // User specified --use-filesystem

		backend, err := configureAppStorage(app, keyContent, expandedKeyFile, useKeyring)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if backend != secrets.StorageBackendFilesystem {
			t.Errorf("Expected backend %v, got %v", secrets.StorageBackendFilesystem, backend)
		}

		if app.PrivateKeyPath != keyPath {
			t.Errorf("Expected PrivateKeyPath %q, got %q", keyPath, app.PrivateKeyPath)
		}

		if app.PrivateKeySource != config.PrivateKeySourceFilesystem {
			t.Errorf("Expected PrivateKeySource %v, got %v", config.PrivateKeySourceFilesystem, app.PrivateKeySource)
		}
	})
}

// TestGetPrivateKey_FilesystemCompatibility tests that getPrivateKey behavior
// is compatible with filesystem storage requirements
func TestGetPrivateKey_FilesystemCompatibility(t *testing.T) {
	t.Run("env var provides content but no path for filesystem", func(t *testing.T) {
		// Set env var
		t.Setenv("GH_APP_PRIVATE_KEY", "test-key-content")

		// No key file specified
		content, path, err := getPrivateKey("")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if content != "test-key-content" {
			t.Errorf("Content = %q, want %q", content, "test-key-content")
		}

		// This is the key point: path is empty when using env var
		if path != "" {
			t.Errorf("Path should be empty when using env var, got %q", path)
		}

		// This empty path will cause configureAppStorage to fail with --use-filesystem
		// which is the correct behavior
	})

	t.Run("key file provides both content and path for filesystem", func(t *testing.T) {
		// Skip on Windows - file permissions are not enforced the same way
		if runtime.GOOS == "windows" {
			t.Skip("Skipping on Windows: file permissions not supported")
		}

		tempDir := t.TempDir()
		t.Setenv("GH_APP_PRIVATE_KEY", "") // Clear env var

		keyPath := filepath.Join(tempDir, "test.pem")
		keyContent := "-----BEGIN RSA PRIVATE KEY-----\ntest\n-----END RSA PRIVATE KEY-----"
		if err := os.WriteFile(keyPath, []byte(keyContent), 0600); err != nil {
			t.Fatalf("Failed to create test key file: %v", err)
		}

		content, path, err := getPrivateKey(keyPath)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if content != keyContent {
			t.Error("Content mismatch")
		}

		// This is the key point: path is populated when using --key-file
		if path != keyPath {
			t.Errorf("Path = %q, want %q", path, keyPath)
		}

		// This populated path allows configureAppStorage to work with --use-filesystem
	})
}
