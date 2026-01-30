package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func FindCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find <query>",
		Short: "Search for packs",
		Long: `Search the packs.sh registry for packs.

Examples:
  packs find commit           # Search for "commit"
  packs find "code review"    # Search phrase
  packs find --type skill git # Filter by type`,
		Args: cobra.MinimumNArgs(1),
		RunE: runFind,
	}

	cmd.Flags().StringP("type", "t", "", "Filter by type (skill, context, prompt)")

	return cmd
}

func runFind(cmd *cobra.Command, args []string) error {
	query := args[0]
	packType, _ := cmd.Flags().GetString("type")

	// TODO: Implement search
	fmt.Printf("🔍 Searching for: %s\n", query)
	if packType != "" {
		fmt.Printf("   Type filter: %s\n", packType)
	}
	fmt.Println("\n   (Registry search coming soon)")
	fmt.Println("\n   For now, try GitHub direct:")
	fmt.Printf("   packs get gh:tunajam/packs/%s\n", query)

	return nil
}
