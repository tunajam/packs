package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "View and manage configuration",
		Long: `View and manage packs configuration.

CONFIGURATION:
  packs stores configuration in ~/.packs/config.yaml

  registry:     URL of the packs registry (default: https://api.packs.sh)
  skills_dir:   Where to install packs (auto-detected by default)
  telemetry:    Enable anonymous usage statistics (default: true)

COMMANDS:
  packs config              Show current configuration
  packs config path         Show config file path
  packs config reset        Reset to defaults

ENVIRONMENT VARIABLES:
  PACKS_REGISTRY    Override registry URL
  PACKS_SKILLS_DIR  Override skills directory
  PACKS_NO_TELEMETRY=1  Disable telemetry`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return showConfig()
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "path",
		Short: "Show config file path",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(getConfigPath())
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "reset",
		Short: "Reset to default configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return resetConfig()
		},
	})

	return cmd
}

func showConfig() error {
	home, _ := os.UserHomeDir()
	configPath := getConfigPath()
	skillsDir := detectAgentSkillsDir()

	registry := "https://api.packs.sh"
	if env := os.Getenv("PACKS_REGISTRY"); env != "" {
		registry = env + " (from env)"
	}

	telemetry := "enabled"
	if os.Getenv("PACKS_NO_TELEMETRY") == "1" {
		telemetry = "disabled (from env)"
	}

	fmt.Printf("\n  📦 packs configuration\n")
	fmt.Printf("  %s\n\n", "────────────────────────────────────────")
	fmt.Printf("  %-14s %s\n", "Config file:", configPath)
	fmt.Printf("  %-14s %s\n", "Registry:", registry)
	fmt.Printf("  %-14s %s\n", "Skills dir:", skillsDir)
	fmt.Printf("  %-14s %s\n", "Telemetry:", telemetry)
	fmt.Printf("  %-14s %s\n", "Cache:", filepath.Join(home, ".packs", "cache"))
	fmt.Println()

	return nil
}

func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".packs", "config.yaml")
}

func resetConfig() error {
	configPath := getConfigPath()
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	defaultConfig := `# packs configuration
# https://packs.sh

registry: https://api.packs.sh
telemetry: true

# Override auto-detected skills directory:
# skills_dir: ~/.packs/skills
`

	if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
		return err
	}

	fmt.Printf("✓ Reset configuration to defaults: %s\n", configPath)
	return nil
}
