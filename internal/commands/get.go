package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

func GetCmd() *cobra.Command {
	var outputFlag string
	var installFlag bool
	var forceFlag bool

	cmd := &cobra.Command{
		Use:   "get <pack>",
		Short: "Install a pack or output its content",
		Long: `Install a pack to your agent's skills directory, or output to stdout.

SOURCES:
  packs get commit-message              Registry (packs.sh)
  packs get commit-message@1.0.0        Specific version
  packs get @user/repo/pack             GitHub shorthand
  packs get gh:user/repo/pack           GitHub explicit

INSTALLATION:
  By default, packs installs to your detected agent's skills directory:
    • Claude Code:  ~/.claude/skills/<pack>/
    • Clawdbot:     ./skills/<pack>/
    • Codex:        ~/.codex/skills/<pack>/
    • Generic:      ~/.packs/skills/<pack>/

  Use --output to specify a custom path, or pipe to handle manually:
    packs get commit-message | pbcopy    # Copy to clipboard
    packs get commit-message > SKILL.md  # Save to file

FLAGS:
  -o, --output <path>   Install to specific directory
  -i, --install         Force install (skip stdout, always write to disk)  
  -f, --force           Overwrite existing pack

EXAMPLES:
  packs get commit-message                    # Install from registry
  packs get @anthropics/skills/docx           # Install from GitHub
  packs get commit-message -o ./my-skills/    # Custom install path
  packs get commit-message | cat              # Output to stdout (pipe detected)`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(args[0], outputFlag, installFlag, forceFlag)
		},
	}

	cmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Install to specific directory")
	cmd.Flags().BoolVarP(&installFlag, "install", "i", false, "Force install to disk")
	cmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Overwrite existing pack")

	return cmd
}

func runGet(pack string, outputDir string, install bool, force bool) error {
	// Normalize @ to gh: 
	if strings.HasPrefix(pack, "@") {
		pack = "gh:" + pack[1:]
	}

	var content string
	var packName string
	var err error

	// Fetch content based on source
	if strings.HasPrefix(pack, "gh:") {
		ref := pack[3:] // Strip "gh:" prefix
		content, packName, err = getFromGitHub(ref)
	} else {
		content, packName, err = getFromRegistry(pack)
	}

	if err != nil {
		return err
	}

	// Determine output mode
	isPiped := !isTerminal()
	
	if isPiped && outputDir == "" && !install {
		// Piped output - just print content
		fmt.Print(content)
		return nil
	}

	// Install mode
	installPath := outputDir
	if installPath == "" {
		installPath = detectAgentSkillsDir()
	}

	// Create pack directory
	packDir := filepath.Join(installPath, packName)
	
	// Check if exists
	if _, err := os.Stat(packDir); err == nil && !force {
		return fmt.Errorf("pack already exists: %s\nUse --force to overwrite", packDir)
	}

	// Create directory
	if err := os.MkdirAll(packDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write SKILL.md
	skillPath := filepath.Join(packDir, "SKILL.md")
	if err := os.WriteFile(skillPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write skill: %w", err)
	}

	fmt.Printf("✓ Installed %s to %s\n", packName, packDir)
	return nil
}

func getFromGitHub(ref string) (content string, name string, err error) {
	// Parse: user/repo or user/repo/path/to/pack
	parts := strings.SplitN(ref, "/", 3)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid GitHub reference: %s\nExpected format: @user/repo or @user/repo/path", ref)
	}

	user := parts[0]
	repo := parts[1]
	path := ""
	if len(parts) > 2 {
		path = parts[2]
	}
	
	// Extract pack name from path or repo name
	if path != "" {
		name = filepath.Base(path)
	} else {
		name = repo
	}

	// Try content files in order: SKILL.md, CONTEXT.md, PROMPT.md
	contentFiles := []string{"SKILL.md", "CONTEXT.md", "PROMPT.md"}
	
	// Try gh CLI first (handles auth, private repos)
	if ghInstalled() {
		for _, file := range contentFiles {
			content, err = ghGetContent(user, repo, path, file)
			if err == nil {
				return content, name, nil
			}
		}
	}

	// Fallback: raw.githubusercontent.com (public repos only)
	for _, file := range contentFiles {
		var url string
		if path != "" {
			url = fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/%s/%s",
				user, repo, path, file)
		} else {
			url = fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/%s",
				user, repo, file)
		}

		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			return string(body), name, nil
		}
	}

	return "", "", fmt.Errorf("pack not found: %s\nTried: SKILL.md, CONTEXT.md, PROMPT.md", ref)
}

func getFromRegistry(pack string) (content string, name string, err error) {
	// Parse version if present: pack@version
	name = pack
	version := "latest"
	if idx := strings.Index(pack, "@"); idx != -1 {
		name = pack[:idx]
		version = pack[idx+1:]
	}

	// TODO: Connect to packs.sh API
	// For now, try GitHub fallback via packs-registry
	registryRef := fmt.Sprintf("tunajam/packs-registry/packs/%s", name)
	content, _, err = getFromGitHub(registryRef)
	if err != nil {
		return "", "", fmt.Errorf("pack not found in registry: %s@%s\n\nTry GitHub direct: packs get @user/repo/%s", name, version, name)
	}
	
	return content, name, nil
}

func detectAgentSkillsDir() string {
	home, _ := os.UserHomeDir()
	
	// Check for Claude Code
	claudeDir := filepath.Join(home, ".claude", "skills")
	if dirExists(filepath.Join(home, ".claude")) {
		return claudeDir
	}

	// Check for Clawdbot (workspace skills/)
	if dirExists("skills") || fileExists("AGENTS.md") || fileExists("SOUL.md") {
		return "skills"
	}

	// Check for Codex
	codexDir := filepath.Join(home, ".codex", "skills")
	if dirExists(filepath.Join(home, ".codex")) {
		return codexDir
	}

	// Check for Cursor
	cursorDir := filepath.Join(home, ".cursor", "skills")
	if dirExists(filepath.Join(home, ".cursor")) {
		return cursorDir
	}

	// Default: ~/.packs/skills
	return filepath.Join(home, ".packs", "skills")
}

func isTerminal() bool {
	fi, _ := os.Stdout.Stat()
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func ghInstalled() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

func ghGetContent(user, repo, path, file string) (string, error) {
	var apiPath string
	if path != "" {
		apiPath = fmt.Sprintf("/repos/%s/%s/contents/%s/%s", user, repo, path, file)
	} else {
		apiPath = fmt.Sprintf("/repos/%s/%s/contents/%s", user, repo, file)
	}
	cmd := exec.Command("gh", "api", apiPath,
		"-H", "Accept: application/vnd.github.raw+json")

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// GetRuntimeInfo returns OS/arch for telemetry
func GetRuntimeInfo() (string, string) {
	return runtime.GOOS, runtime.GOARCH
}
