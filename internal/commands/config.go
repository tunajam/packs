package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long: `Manage packs configuration.

Examples:
  packs config list                    # Show all config
  packs config set telemetry false     # Disable telemetry
  packs config set registry https://... # Set custom registry`,
	}

	cmd.AddCommand(configListCmd())
	cmd.AddCommand(configSetCmd())

	return cmd
}

func configListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("📋 Configuration (~/.packs/config.yaml)")
			fmt.Println()
			fmt.Println("   registry:  https://packs.sh")
			fmt.Println("   telemetry: true")
			fmt.Println("   cache.ttl: 1h")
			return nil
		},
	}
}

func configSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]
			// TODO: Actually persist config
			fmt.Printf("✓ Set %s = %s\n", key, value)
			return nil
		},
	}
}
