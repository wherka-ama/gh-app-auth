package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func NewConfigCmd() *cobra.Command {
	var showPath bool
	var showContent bool

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Show configuration file location and content",
		Long: `Display the configuration file path and optionally its content.

The configuration file location follows this priority:
  1. GH_APP_AUTH_CONFIG environment variable (if set)
  2. Default: ~/.config/gh/extensions/gh-app-auth/config.yml`,
		Example: `  # Show config file path
  gh app-auth config

  # Show config file path only
  gh app-auth config --path

  # Show config file content
  gh app-auth config --show`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return configRun(showPath, showContent)
		},
	}

	cmd.Flags().BoolVarP(&showPath, "path", "p", false, "Show only the config file path")
	cmd.Flags().BoolVarP(&showContent, "show", "s", false, "Show the config file content")

	return cmd
}

func configRun(showPath, showContent bool) error {
	configPath := getConfigPath()

	// If only path requested
	if showPath {
		fmt.Println(configPath)
		return nil
	}

	// Check if config file exists
	exists := fileExists(configPath)

	// Default: show path and status
	if !showContent {
		fmt.Printf("📁 Configuration file: %s\n", configPath)
		if envPath := os.Getenv("GH_APP_AUTH_CONFIG"); envPath != "" {
			fmt.Println("   (set via GH_APP_AUTH_CONFIG environment variable)")
		}
		if exists {
			fmt.Println("   Status: ✅ exists")
		} else {
			fmt.Println("   Status: ⚠️  not found")
			fmt.Println("\n💡 Run 'gh app-auth setup' to create a configuration.")
		}
		return nil
	}

	// Show content
	if !exists {
		return fmt.Errorf("configuration file not found: %s\nRun 'gh app-auth setup' to create one", configPath)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	fmt.Printf("# Configuration file: %s\n", configPath)
	fmt.Println("---")
	fmt.Print(string(content))

	return nil
}

// getConfigPath returns the configuration file path
func getConfigPath() string {
	// Check environment variable first
	if path := os.Getenv("GH_APP_AUTH_CONFIG"); path != "" {
		return path
	}

	// Use default path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".config", "gh", "extensions", "gh-app-auth", "config.yml")
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetConfigPath is exported for use by other packages
func GetConfigPath() string {
	return getConfigPath()
}
