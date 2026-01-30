package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tunajam/packs/internal/api"
)

// PackInfo represents pack metadata for search results (JSON output)
type PackInfo struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Stars       int      `json:"stars"`
	Tags        []string `json:"tags,omitempty"`
	Source      string   `json:"source"` // "registry" or "github"
}

func FindCmd() *cobra.Command {
	var typeFlag string
	var limitFlag int
	var jsonFlag bool

	cmd := &cobra.Command{
		Use:   "find [query]",
		Short: "Search for packs",
		Long: `Search the packs registry for skills, contexts, and prompts.

SEARCH:
  packs find                          List popular packs
  packs find "commit message"         Search by keyword
  packs find --type skill             Filter by type
  packs find --json                   Output as JSON (for agents)

TYPES:
  skill     Procedural instructions (how to do X)
  context   Domain knowledge (what is X)
  prompt    Ready-to-use prompts

OUTPUT FORMATS:
  Default:  Human-readable table
  --json:   Machine-readable JSON array

EXAMPLES:
  packs find git                      # Search for git-related packs
  packs find --type context react     # React context packs
  packs find --limit 5                # Top 5 results
  packs find --json | jq '.[0].name'  # Parse with jq`,
		RunE: func(cmd *cobra.Command, args []string) error {
			query := ""
			if len(args) > 0 {
				query = strings.Join(args, " ")
			}
			return runFind(query, typeFlag, limitFlag, jsonFlag)
		},
	}

	cmd.Flags().StringVarP(&typeFlag, "type", "t", "", "Filter by type: skill, context, prompt")
	cmd.Flags().IntVarP(&limitFlag, "limit", "l", 20, "Maximum results to return")
	cmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output as JSON")

	return cmd
}

func runFind(query string, packType string, limit int, jsonOutput bool) error {
	client := api.New()
	ctx := context.Background()

	opts := api.SearchOpts{
		Query: query,
		Type:  packType,
		Limit: int32(limit),
		Sort:  "stars", // Default sort by popularity
	}

	packs, total, err := client.Search(ctx, opts)
	if err != nil {
		// If API fails, fall back to demo data for offline/dev use
		return runFindOffline(query, packType, limit, jsonOutput)
	}

	// Convert to output format
	var results []PackInfo
	for _, p := range packs {
		results = append(results, PackInfo{
			Name:        p.Name,
			Version:     p.Version,
			Type:        p.Type,
			Description: p.Description,
			Author:      p.Author,
			Stars:       int(p.Stars),
			Tags:        p.Tags,
			Source:      "registry",
		})
	}

	// Output
	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}

	// Human-readable output
	if len(results) == 0 {
		fmt.Println("No packs found.")
		return nil
	}

	fmt.Printf("\n  Found %d packs (total: %d):\n\n", len(results), total)
	for _, p := range results {
		typeIcon := "📦"
		switch p.Type {
		case "context":
			typeIcon = "📚"
		case "prompt":
			typeIcon = "💬"
		}
		fmt.Printf("  %s %-24s %-8s  ★ %-4d  %s\n",
			typeIcon, p.Name, p.Version, p.Stars, truncate(p.Description, 40))
	}
	fmt.Printf("\n  Run: packs get <name> to install\n\n")

	return nil
}

// runFindOffline provides fallback demo data when API is unavailable
func runFindOffline(query string, packType string, limit int, jsonOutput bool) error {
	packs := getDemoPacks()

	// Filter by query
	if query != "" {
		query = strings.ToLower(query)
		var filtered []PackInfo
		for _, p := range packs {
			if strings.Contains(strings.ToLower(p.Name), query) ||
				strings.Contains(strings.ToLower(p.Description), query) {
				filtered = append(filtered, p)
			}
		}
		packs = filtered
	}

	// Filter by type
	if packType != "" {
		var filtered []PackInfo
		for _, p := range packs {
			if p.Type == packType {
				filtered = append(filtered, p)
			}
		}
		packs = filtered
	}

	// Apply limit
	if limit > 0 && len(packs) > limit {
		packs = packs[:limit]
	}

	// Output
	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(packs)
	}

	// Human-readable output
	if len(packs) == 0 {
		fmt.Println("No packs found.")
		return nil
	}

	fmt.Printf("\n  Found %d packs (offline mode):\n\n", len(packs))
	for _, p := range packs {
		typeIcon := "📦"
		switch p.Type {
		case "context":
			typeIcon = "📚"
		case "prompt":
			typeIcon = "💬"
		}
		fmt.Printf("  %s %-24s %-8s  ★ %-4d  %s\n",
			typeIcon, p.Name, p.Version, p.Stars, truncate(p.Description, 40))
	}
	fmt.Printf("\n  Run: packs get <name> to install\n\n")

	return nil
}

func getDemoPacks() []PackInfo {
	return []PackInfo{
		{Name: "commit-message", Version: "1.0.0", Type: "skill", Description: "Generate conventional commit messages", Author: "tunajam", Stars: 892, Source: "registry"},
		{Name: "pr-description", Version: "1.0.0", Type: "skill", Description: "Write PR descriptions from branch diff", Author: "tunajam", Stars: 654, Source: "registry"},
		{Name: "humanizer", Version: "1.0.0", Type: "skill", Description: "Remove AI patterns from writing", Author: "blader", Stars: 543, Source: "registry"},
		{Name: "claudeception", Version: "1.0.0", Type: "skill", Description: "Extract learnings into reusable skills", Author: "blader", Stars: 421, Source: "registry"},
		{Name: "test-driven-development", Version: "1.0.0", Type: "skill", Description: "TDD workflow for features and bugfixes", Author: "obra", Stars: 389, Source: "registry"},
		{Name: "brainstorming", Version: "1.0.0", Type: "skill", Description: "Structured ideation and design exploration", Author: "obra", Stars: 312, Source: "registry"},
		{Name: "changelog-generator", Version: "1.0.0", Type: "skill", Description: "Generate changelogs from git history", Author: "composio", Stars: 287, Source: "registry"},
		{Name: "git-worktrees", Version: "1.0.0", Type: "skill", Description: "Work with isolated git worktrees", Author: "obra", Stars: 245, Source: "registry"},
		{Name: "react-query", Version: "2.1.0", Type: "context", Description: "React Query patterns and best practices", Author: "tunajam", Stars: 1247, Source: "registry"},
		{Name: "drizzle-orm", Version: "1.0.0", Type: "context", Description: "Drizzle ORM conventions and patterns", Author: "tunajam", Stars: 876, Source: "registry"},
		{Name: "mcp-builder", Version: "1.0.0", Type: "skill", Description: "Build MCP servers for LLM integrations", Author: "composio", Stars: 534, Source: "registry"},
		{Name: "youtube-transcript", Version: "1.0.0", Type: "skill", Description: "Fetch and summarize YouTube transcripts", Author: "tapestry", Stars: 423, Source: "registry"},
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
