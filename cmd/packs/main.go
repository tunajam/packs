package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/tunajam/packs/internal/commands"
)

const banner = `                    __        
    ____  ____ _____/ /_______
   / __ \/ __ ` + "`" + `/ __/ //_/ ___/
  / /_/ / /_/ / /_/ ,< (__  ) 
 / .___/\__,_/\__/_/|_/____/  
/_/                           

  skills.sh compatible. Enterprise ready.
`

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	titleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true)
	cmdStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	descStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	accentStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
)

func styledHelp(cmd *cobra.Command, args []string) {
	fmt.Println()
	fmt.Println(accentStyle.Render(banner))
	
	fmt.Println(titleStyle.Render("  USAGE"))
	fmt.Printf("    %s\n\n", dimStyle.Render("packs [command]"))
	
	fmt.Println(titleStyle.Render("  COMMANDS"))
	fmt.Printf("    %s  %s\n", cmdStyle.Render("packs             "), descStyle.Render("Browse packs in TUI"))
	fmt.Printf("    %s  %s\n", cmdStyle.Render("packs find <query>"), descStyle.Render("Search for packs"))
	fmt.Printf("    %s  %s\n", cmdStyle.Render("packs get <name>  "), descStyle.Render("Install a pack"))
	fmt.Printf("    %s  %s\n", cmdStyle.Render("packs info <name> "), descStyle.Render("Show pack details"))
	fmt.Printf("    %s  %s\n", cmdStyle.Render("packs submit <ref>"), descStyle.Render("Submit a pack to registry"))
	fmt.Printf("    %s  %s\n", cmdStyle.Render("packs config      "), descStyle.Render("Show or set configuration"))
	fmt.Println()
	
	fmt.Println(titleStyle.Render("  GITHUB FETCH"))
	fmt.Printf("    %s\n", cmdStyle.Render("packs get gh:user/repo/path/to/pack"))
	fmt.Printf("    %s\n\n", dimStyle.Render("Fetch directly from GitHub (works with private repos)"))
	
	fmt.Println(titleStyle.Render("  AUTH"))
	fmt.Printf("    %s  %s\n", cmdStyle.Render("packs login "), descStyle.Render("Authenticate with GitHub"))
	fmt.Printf("    %s  %s\n", cmdStyle.Render("packs logout"), descStyle.Render("Log out"))
	fmt.Printf("    %s  %s\n", cmdStyle.Render("packs whoami"), descStyle.Render("Show current user"))
	fmt.Println()
	
	fmt.Println(titleStyle.Render("  OPTIONS"))
	fmt.Printf("    %s  %s\n", cmdStyle.Render("-j, --json    "), descStyle.Render("Output as JSON"))
	fmt.Printf("    %s  %s\n", cmdStyle.Render("    --no-cache"), descStyle.Render("Bypass local cache"))
	fmt.Printf("    %s  %s\n", cmdStyle.Render("-h, --help    "), descStyle.Render("Show this help"))
	fmt.Printf("    %s  %s\n", cmdStyle.Render("-v, --version "), descStyle.Render("Show version"))
	fmt.Println()
	
	fmt.Printf("  %s %s\n\n", dimStyle.Render("Docs:"), accentStyle.Render("https://packs.sh"))
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "packs",
		Short: "skills.sh compatible. Enterprise ready.",
		Long:  banner + "\n  packs is a TUI-first tool for discovering, installing,\n  and sharing AI skills, prompts, and context.",
		Run: func(cmd *cobra.Command, args []string) {
			// No args = launch TUI
			commands.RunTUI()
		},
	}
	
	// Custom help
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		styledHelp(cmd, args)
	})

	// Version flag
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(banner + "\n  v{{.Version}} · packs.sh\n\n")

	// Add commands
	rootCmd.AddCommand(commands.GetCmd())
	rootCmd.AddCommand(commands.FindCmd())
	rootCmd.AddCommand(commands.InfoCmd())
	rootCmd.AddCommand(commands.SubmitCmd())
	rootCmd.AddCommand(commands.ConfigCmd())
	rootCmd.AddCommand(commands.LoginCmd())
	rootCmd.AddCommand(commands.LogoutCmd())
	rootCmd.AddCommand(commands.WhoamiCmd())

	// Global flags
	rootCmd.PersistentFlags().BoolP("json", "j", false, "Output as JSON")
	rootCmd.PersistentFlags().Bool("no-cache", false, "Bypass local cache")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
