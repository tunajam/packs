package commands

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func GetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <pack>",
		Short: "Get a pack's content",
		Long: `Get a pack's content and output to stdout.

Examples:
  packs get commit-message              # From packs.sh registry
  packs get commit-message@1.0.0        # Specific version
  packs get gh:user/repo/pack           # Direct from GitHub
  packs get gh:user/repo/path/to/pack   # With subdirectories

Pipe to clipboard:
  packs get commit-message | pbcopy     # macOS
  packs get commit-message | xclip      # Linux`,
		Args: cobra.ExactArgs(1),
		RunE: runGet,
	}

	return cmd
}

func runGet(cmd *cobra.Command, args []string) error {
	pack := args[0]

	// Check if it's a GitHub direct reference
	if strings.HasPrefix(pack, "gh:") {
		return getFromGitHub(pack[3:]) // Strip "gh:" prefix
	}

	// Otherwise, fetch from registry
	return getFromRegistry(pack)
}

func getFromGitHub(ref string) error {
	// Parse: user/repo/path/to/pack
	parts := strings.SplitN(ref, "/", 3)
	if len(parts) < 3 {
		return fmt.Errorf("invalid GitHub reference: %s\nExpected format: gh:user/repo/pack", ref)
	}

	user := parts[0]
	repo := parts[1]
	path := parts[2]

	// Try to determine content file - check for pack.yaml first
	contentFile := "SKILL.md" // default

	// Try gh CLI first (handles auth, private repos)
	if ghInstalled() {
		content, err := ghGetContent(user, repo, path, contentFile)
		if err == nil {
			fmt.Print(content)
			return nil
		}
		// Fall through to raw fetch
	}

	// Fallback: raw.githubusercontent.com (public repos only)
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/%s/%s",
		user, repo, path, contentFile)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch from GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		// Try CONTEXT.md
		url = fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/%s/CONTEXT.md",
			user, repo, path)
		resp, err = http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch from GitHub: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 404 {
			// Try PROMPT.md
			url = fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/%s/PROMPT.md",
				user, repo, path)
			resp, err = http.Get(url)
			if err != nil {
				return fmt.Errorf("failed to fetch from GitHub: %w", err)
			}
			defer resp.Body.Close()
		}
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("pack not found: %s (HTTP %d)", ref, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Print(string(body))
	return nil
}

func getFromRegistry(pack string) error {
	// Parse version if present: pack@version
	name := pack
	version := "latest"
	if idx := strings.Index(pack, "@"); idx != -1 {
		name = pack[:idx]
		version = pack[idx+1:]
	}

	// TODO: Implement registry fetch
	// For now, return a helpful message
	return fmt.Errorf("registry fetch not yet implemented\nPack: %s, Version: %s\n\nTry using GitHub direct: packs get gh:user/repo/%s", name, version, name)
}

func ghInstalled() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

func ghGetContent(user, repo, path, file string) (string, error) {
	// Use gh api to get raw content
	cmd := exec.Command("gh", "api",
		fmt.Sprintf("/repos/%s/%s/contents/%s/%s", user, repo, path, file),
		"--jq", ".content",
		"-H", "Accept: application/vnd.github.raw+json")

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
