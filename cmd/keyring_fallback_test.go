package cmd

import (
	"runtime"
	"strings"
	"testing"

	"github.com/AmadeusITGroup/gh-app-auth/pkg/config"
)

// TestGetKeyringInstallInstructions tests OS-specific keyring installation instructions
func TestGetKeyringInstallInstructions(t *testing.T) {
	tests := []struct {
		name            string
		goos            string
		expectedKeyword string
		checkContent    func(string) bool
	}{
		{
			name:            "Linux instructions",
			goos:            "linux",
			expectedKeyword: "gnome-keyring",
			checkContent: func(s string) bool {
				return strings.Contains(s, "gnome-keyring") ||
					strings.Contains(s, "libsecret") ||
					strings.Contains(s, "pass")
			},
		},
		{
			name:            "macOS instructions",
			goos:            "darwin",
			expectedKeyword: "Keychain",
			checkContent: func(s string) bool {
				return strings.Contains(s, "Keychain") ||
					strings.Contains(s, "macOS")
			},
		},
		{
			name:            "Windows instructions",
			goos:            "windows",
			expectedKeyword: "Credential Manager",
			checkContent: func(s string) bool {
				return strings.Contains(s, "Credential Manager") ||
					strings.Contains(s, "Windows")
			},
		},
		{
			name:            "FreeBSD instructions",
			goos:            "freebsd",
			expectedKeyword: "gnome-keyring",
			checkContent: func(s string) bool {
				return strings.Contains(s, "gnome-keyring") ||
					strings.Contains(s, "pass")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instructions := getKeyringInstallInstructions(tt.goos)

			if instructions == "" {
				t.Error("Expected non-empty instructions")
			}

			if !tt.checkContent(instructions) {
				t.Errorf("Instructions missing expected keyword %q. Got:\n%s",
					tt.expectedKeyword, instructions)
			}

			// Check for non-root installation hints
			if tt.goos == "linux" || tt.goos == "freebsd" {
				if !strings.Contains(instructions, "without root") &&
					!strings.Contains(instructions, "user-level") &&
					!strings.Contains(instructions, "home directory") {
					t.Error("Linux/FreeBSD instructions should mention non-root options")
				}
			}
		})
	}
}

// TestConfigureAppStorage_KeyringUnavailableWithEnvVar tests the scenario where
// keyring is unavailable and user tries to use GH_APP_PRIVATE_KEY env var
func TestConfigureAppStorage_KeyringUnavailableWithEnvVar(t *testing.T) {
	t.Run("keyring unavailable with env var should provide helpful error", func(t *testing.T) {
		app := &config.GitHubApp{
			AppID:          123456,
			Name:           "test-app",
			InstallationID: 789,
			Patterns:       []string{"github.com/org/*"},
			Priority:       1,
		}

		// Simulate key from env var (no file path)
		expandedKeyFile := "" // Empty because key came from env var
		useKeyring := true    // User wants keyring (default)

		// This should fail when keyring is unavailable
		// We'll need to mock keyring unavailability in the implementation
		_, err := configureAppStorageWithKeyringCheck(
			app, expandedKeyFile, useKeyring, false, // false = keyring unavailable
		)

		if err == nil {
			t.Fatal("Expected error when keyring unavailable with env var")
		}

		errMsg := err.Error()

		// Error should mention the problem
		if !strings.Contains(errMsg, "keyring") && !strings.Contains(errMsg, "unavailable") {
			t.Errorf("Error should mention keyring unavailability. Got: %s", errMsg)
		}

		// Error should mention GH_APP_PRIVATE_KEY
		if !strings.Contains(errMsg, "GH_APP_PRIVATE_KEY") {
			t.Errorf("Error should mention GH_APP_PRIVATE_KEY. Got: %s", errMsg)
		}

		// Error should mention --key-file as alternative
		if !strings.Contains(errMsg, "--key-file") && !strings.Contains(errMsg, "key-file") {
			t.Errorf("Error should mention --key-file as alternative. Got: %s", errMsg)
		}

		// Error should include OS-specific instructions
		osName := runtime.GOOS
		switch osName {
		case "linux":
			if !strings.Contains(errMsg, "gnome-keyring") && !strings.Contains(errMsg, "libsecret") {
				t.Errorf("Linux error should mention keyring installation options. Got: %s", errMsg)
			}
		case "darwin":
			if !strings.Contains(errMsg, "Keychain") {
				t.Errorf("macOS error should mention Keychain. Got: %s", errMsg)
			}
		case "windows":
			if !strings.Contains(errMsg, "Credential Manager") {
				t.Errorf("Windows error should mention Credential Manager. Got: %s", errMsg)
			}
		}
	})

	t.Run("keyring unavailable with key file should fallback to filesystem", func(t *testing.T) {
		app := &config.GitHubApp{
			AppID:          123456,
			Name:           "test-app",
			InstallationID: 789,
			Patterns:       []string{"github.com/org/*"},
			Priority:       1,
		}

		expandedKeyFile := "/path/to/key.pem" // File path provided
		useKeyring := true                    // User wants keyring (default)

		// This should gracefully fallback to filesystem when keyring unavailable
		backend, err := configureAppStorageWithKeyringCheck(
			app, expandedKeyFile, useKeyring, false, // false = keyring unavailable
		)

		if err != nil {
			t.Fatalf("Should fallback to filesystem gracefully. Got error: %v", err)
		}

		if backend != "filesystem" {
			t.Errorf("Expected filesystem backend, got %v", backend)
		}

		if app.PrivateKeyPath != expandedKeyFile {
			t.Errorf("Expected PrivateKeyPath %q, got %q", expandedKeyFile, app.PrivateKeyPath)
		}
	})

	t.Run("keyring available with env var should work", func(t *testing.T) {
		app := &config.GitHubApp{
			AppID:          123456,
			Name:           "test-app",
			InstallationID: 789,
			Patterns:       []string{"github.com/org/*"},
			Priority:       1,
		}

		expandedKeyFile := "" // Empty because key came from env var
		useKeyring := true    // User wants keyring (default)

		// This should work when keyring is available
		backend, err := configureAppStorageWithKeyringCheck(
			app, expandedKeyFile, useKeyring, true, // true = keyring available
		)

		if err != nil {
			t.Fatalf("Should work with keyring available. Got error: %v", err)
		}

		if backend != "keyring" {
			t.Errorf("Expected keyring backend, got %v", backend)
		}
	})
}

// TestFormatKeyringUnavailableError tests the error message formatting
func TestFormatKeyringUnavailableError(t *testing.T) {
	tests := []struct {
		name         string
		goos         string
		hasKeyFile   bool
		checkContent func(string) bool
	}{
		{
			name:       "Linux without key file",
			goos:       "linux",
			hasKeyFile: false,
			checkContent: func(s string) bool {
				return strings.Contains(s, "keyring") &&
					strings.Contains(s, "GH_APP_PRIVATE_KEY") &&
					strings.Contains(s, "--key-file") &&
					(strings.Contains(s, "gnome-keyring") || strings.Contains(s, "libsecret"))
			},
		},
		{
			name:       "macOS without key file",
			goos:       "darwin",
			hasKeyFile: false,
			checkContent: func(s string) bool {
				return strings.Contains(s, "keyring") &&
					strings.Contains(s, "Keychain")
			},
		},
		{
			name:       "Windows without key file",
			goos:       "windows",
			hasKeyFile: false,
			checkContent: func(s string) bool {
				return strings.Contains(s, "keyring") &&
					strings.Contains(s, "Credential Manager")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errMsg := formatKeyringUnavailableError(tt.goos, tt.hasKeyFile)

			if errMsg == "" {
				t.Error("Expected non-empty error message")
			}

			if !tt.checkContent(errMsg) {
				t.Errorf("Error message missing expected content. Got:\n%s", errMsg)
			}
		})
	}
}
