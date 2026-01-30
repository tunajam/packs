package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tunajam/packs/internal/commands"
)

const banner = `                    __        
    ____  ____ _____/ /_______
   / __ \/ __ ` + "`" + `/ __/ //_/ ___/
  / /_/ / /_/ / /_/ ,< (__  ) 
 / .___/\__,_/\__/_/|_/____/  
/_/                           

  Skills for AI agents. One command.
`

var version = "0.1.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "packs",
		Short: "Skills for AI agents. One command.",
		Long:  banner + "\n  packs is a TUI-first tool for discovering, installing,\n  and sharing AI skills, prompts, and context.",
		Run: func(cmd *cobra.Command, args []string) {
			// No args = launch TUI
			commands.RunTUI()
		},
	}

	// Version flag
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(banner + "\n  v{{.Version}} · packs.sh\n\n")

	// Add commands
	rootCmd.AddCommand(commands.GetCmd())
	rootCmd.AddCommand(commands.FindCmd())
	rootCmd.AddCommand(commands.InfoCmd())
	rootCmd.AddCommand(commands.SubmitCmd())
	rootCmd.AddCommand(commands.ConfigCmd())

	// Global flags
	rootCmd.PersistentFlags().BoolP("json", "j", false, "Output as JSON")
	rootCmd.PersistentFlags().Bool("no-cache", false, "Bypass local cache")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
