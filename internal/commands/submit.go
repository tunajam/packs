package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tunajam/packs/internal/api"
)

func SubmitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit <github-ref>",
		Short: "Submit a pack to the registry",
		Long: `Submit a GitHub-hosted pack to the packs.sh registry for indexing.

REQUIREMENTS:
  Your pack must be a public GitHub repository containing:
    • pack.yaml    - Metadata (name, version, type, description)
    • SKILL.md     - For skill packs
    • CONTEXT.md   - For context packs  
    • PROMPT.md    - For prompt packs

PACK.YAML FORMAT:
  name: my-skill
  version: 1.0.0
  type: skill
  description: What this skill does
  author: your-name
  license: MIT
  tags:
    - tag1
    - tag2

SUBMIT FORMATS:
  packs submit @user/repo/path        GitHub shorthand
  packs submit gh:user/repo/path      GitHub explicit

EXAMPLES:
  packs submit @myname/skills/commit-helper
  packs submit gh:anthropics/skills/docx

WHAT HAPPENS:
  1. Validates pack structure and metadata
  2. Fetches content and computes hash
  3. Indexes in packs.sh registry
  4. Pack becomes available via 'packs get'`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSubmit(args[0])
		},
	}

	return cmd
}

func runSubmit(ref string) error {
	// Normalize @ to gh:
	if strings.HasPrefix(ref, "@") {
		ref = ref[1:] // Strip @
	} else if strings.HasPrefix(ref, "gh:") {
		ref = ref[3:] // Strip gh:
	}

	// Validate format
	parts := strings.SplitN(ref, "/", 3)
	if len(parts) < 3 {
		return fmt.Errorf("invalid reference: %s\nExpected format: @user/repo/path or gh:user/repo/path", ref)
	}

	fmt.Printf("\n  📦 Submitting %s...\n\n", ref)

	// Submit to API
	client := api.New()
	ctx := context.Background()
	
	name, version, message, err := client.Submit(ctx, ref)
	if err != nil {
		return fmt.Errorf("failed to submit: %w", err)
	}

	fmt.Printf("  ✓ Submitted to registry\n")
	if message != "" {
		fmt.Printf("  ℹ %s\n", message)
	}
	fmt.Printf("\n  🎉 Pack submitted successfully!\n")
	fmt.Printf("  Available via: packs get %s", name)
	if version != "" {
		fmt.Printf("@%s", version)
	}
	fmt.Printf("\n\n")

	return nil
}
