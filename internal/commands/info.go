package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func InfoCmd() *cobra.Command {
	var jsonFlag bool

	cmd := &cobra.Command{
		Use:   "info <pack>",
		Short: "Show detailed pack information",
		Long: `Display detailed information about a pack.

USAGE:
  packs info commit-message           Show pack details
  packs info commit-message@1.0.0     Specific version
  packs info --json commit-message    Output as JSON

INFORMATION SHOWN:
  • Name, version, type
  • Description and author
  • Stars and download count
  • Tags and license
  • Available versions
  • Source (registry or GitHub ref)

EXAMPLES:
  packs info humanizer                # View humanizer details
  packs info --json react-query       # JSON output for scripts`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInfo(args[0], jsonFlag)
		},
	}

	cmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output as JSON")

	return cmd
}

func runInfo(pack string, jsonOutput bool) error {
	// Parse version
	name := pack
	version := "latest"
	if idx := strings.Index(pack, "@"); idx != -1 {
		name = pack[:idx]
		version = pack[idx+1:]
	}

	// TODO: Fetch from packs.sh API
	// For now, return demo data
	info := getPackInfo(name, version)
	if info == nil {
		return fmt.Errorf("pack not found: %s", pack)
	}

	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(info)
	}

	// Human-readable output
	typeIcon := "📦"
	switch info.Type {
	case "context":
		typeIcon = "📚"
	case "prompt":
		typeIcon = "💬"
	}

	fmt.Printf("\n  %s %s\n", typeIcon, info.Name)
	fmt.Printf("  %s\n\n", strings.Repeat("─", 50))
	fmt.Printf("  %-14s %s\n", "Version:", info.Version)
	fmt.Printf("  %-14s %s\n", "Type:", info.Type)
	fmt.Printf("  %-14s %s\n", "Author:", info.Author)
	fmt.Printf("  %-14s ★ %d\n", "Stars:", info.Stars)
	fmt.Printf("  %-14s %s\n", "License:", info.License)
	fmt.Printf("\n  %s\n", info.Description)
	if len(info.Tags) > 0 {
		fmt.Printf("\n  Tags: %s\n", strings.Join(info.Tags, ", "))
	}
	if len(info.Versions) > 0 {
		fmt.Printf("\n  Versions: %s\n", strings.Join(info.Versions, ", "))
	}
	fmt.Printf("\n  Install: packs get %s\n\n", info.Name)

	return nil
}

type PackDetail struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Stars       int      `json:"stars"`
	Downloads   int      `json:"downloads"`
	License     string   `json:"license"`
	Tags        []string `json:"tags"`
	Versions    []string `json:"versions"`
	GithubRef   string   `json:"github_ref,omitempty"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

func getPackInfo(name, version string) *PackDetail {
	// Demo data
	packs := map[string]*PackDetail{
		"commit-message": {
			Name:        "commit-message",
			Version:     "1.0.0",
			Type:        "skill",
			Description: "Generate conventional commit messages from staged changes. Analyzes git diff and produces well-formatted commits following the Conventional Commits specification.",
			Author:      "tunajam",
			Stars:       892,
			Downloads:   4521,
			License:     "MIT",
			Tags:        []string{"git", "commits", "conventional-commits"},
			Versions:    []string{"1.0.0"},
			CreatedAt:   "2026-01-15",
			UpdatedAt:   "2026-01-28",
		},
		"humanizer": {
			Name:        "humanizer",
			Version:     "1.0.0",
			Type:        "skill",
			Description: "Remove signs of AI-generated writing from text. Based on Wikipedia's comprehensive guide to AI writing patterns, detecting 24 common issues.",
			Author:      "blader",
			Stars:       543,
			Downloads:   2187,
			License:     "MIT",
			Tags:        []string{"writing", "editing", "ai-detection"},
			Versions:    []string{"1.0.0"},
			GithubRef:   "blader/humanizer",
			CreatedAt:   "2026-01-20",
			UpdatedAt:   "2026-01-29",
		},
		"react-query": {
			Name:        "react-query",
			Version:     "2.1.0",
			Type:        "context",
			Description: "React Query (TanStack Query) patterns, best practices, and common pitfalls. Comprehensive reference for data fetching, caching, and state management.",
			Author:      "tunajam",
			Stars:       1247,
			Downloads:   6892,
			License:     "MIT",
			Tags:        []string{"react", "tanstack", "data-fetching", "caching"},
			Versions:    []string{"2.1.0", "2.0.0", "1.0.0"},
			CreatedAt:   "2026-01-10",
			UpdatedAt:   "2026-01-29",
		},
	}

	if p, ok := packs[name]; ok {
		return p
	}
	return nil
}
