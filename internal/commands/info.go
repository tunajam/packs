package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func InfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <pack>",
		Short: "Show pack metadata",
		Long: `Show detailed metadata for a pack.

Examples:
  packs info commit-message
  packs info commit-message --versions`,
		Args: cobra.ExactArgs(1),
		RunE: runInfo,
	}

	cmd.Flags().Bool("versions", false, "List all available versions")

	return cmd
}

func runInfo(cmd *cobra.Command, args []string) error {
	pack := args[0]

	// TODO: Implement info fetch
	fmt.Printf("📦 %s\n", pack)
	fmt.Println("   (Pack info coming soon)")

	return nil
}
