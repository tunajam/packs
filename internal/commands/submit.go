package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func SubmitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit <gh:user/repo/pack>",
		Short: "Submit a GitHub pack for indexing",
		Long: `Submit a GitHub pack to be indexed in packs.sh.

This makes your pack discoverable via 'packs find' and the TUI.
No authentication required - just indexes public repos.

Examples:
  packs submit gh:hsbacot/packs/commit-message
  packs submit gh:hsbacot/packs/skills/code-review`,
		Args: cobra.ExactArgs(1),
		RunE: runSubmit,
	}

	return cmd
}

func runSubmit(cmd *cobra.Command, args []string) error {
	ref := args[0]

	if !strings.HasPrefix(ref, "gh:") {
		return fmt.Errorf("submit requires a GitHub reference\nExample: packs submit gh:user/repo/pack")
	}

	// TODO: Implement submit to packs.sh
	fmt.Printf("📤 Submitting: %s\n", ref)
	fmt.Println("   (Submit endpoint coming soon)")
	fmt.Println("\n   Your pack is already accessible via:")
	fmt.Printf("   packs get %s\n", ref)

	return nil
}
