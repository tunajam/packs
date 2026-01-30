package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const (
	authBaseURL = "https://packs-api.fly.dev"
)

func LoginCmd() *cobra.Command {
	var tokenFlag string
	
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with GitHub",
		Long: `Authenticate with GitHub to publish packs.

This will open your browser to authenticate with GitHub.
After authenticating, copy the token and paste it here.

EXAMPLES:
  packs login                  Interactive login
  packs login --token "..."    Login with existing token`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if tokenFlag != "" {
				return saveToken(tokenFlag)
			}
			return interactiveLogin()
		},
	}

	cmd.Flags().StringVar(&tokenFlag, "token", "", "Provide token directly")
	
	return cmd
}

func LogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Log out and remove stored credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			tokenPath := getTokenPath()
			if err := os.Remove(tokenPath); err != nil {
				if os.IsNotExist(err) {
					fmt.Println("Not logged in")
					return nil
				}
				return err
			}
			fmt.Println("✓ Logged out")
			return nil
		},
	}
}

func WhoamiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show current authenticated user",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := getStoredToken()
			if err != nil || token == "" {
				fmt.Println("Not logged in. Run: packs login")
				return nil
			}

			// Call /auth/me to get user info
			req, _ := http.NewRequest("GET", authBaseURL+"/auth/me", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("failed to check auth: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Println("Token invalid or expired. Run: packs login")
				return nil
			}

			var user struct {
				GitHubLogin string `json:"github_login"`
				GitHubID    int64  `json:"github_id"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
				return err
			}

			fmt.Printf("Logged in as @%s (ID: %d)\n", user.GitHubLogin, user.GitHubID)
			return nil
		},
	}
}

func interactiveLogin() error {
	// First, try to use gh CLI if available
	if token, user := tryGhCLI(); token != "" {
		fmt.Printf("✓ Found gh CLI, authenticated as @%s\n", user)
		return saveTokenQuiet(token)
	}

	// Fall back to browser OAuth
	fmt.Println("Opening browser for GitHub authentication...")
	fmt.Println()

	// Open browser
	loginURL := authBaseURL + "/auth/login"
	if err := openBrowser(loginURL); err != nil {
		fmt.Printf("Could not open browser. Please visit:\n%s\n\n", loginURL)
	}

	// Wait for user to paste token
	fmt.Println("After authenticating, paste the token here and press Enter:")
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	token, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	token = strings.TrimSpace(token)

	if token == "" {
		return fmt.Errorf("no token provided")
	}

	return saveToken(token)
}

// tryGhCLI attempts to get a token from the gh CLI
func tryGhCLI() (token string, username string) {
	// Check if gh is installed
	ghPath, err := exec.LookPath("gh")
	if err != nil || ghPath == "" {
		return "", ""
	}

	// Check if gh is authenticated
	cmd := exec.Command("gh", "auth", "status")
	if err := cmd.Run(); err != nil {
		return "", ""
	}

	// Get the token
	cmd = exec.Command("gh", "auth", "token")
	out, err := cmd.Output()
	if err != nil {
		return "", ""
	}
	token = strings.TrimSpace(string(out))
	if token == "" {
		return "", ""
	}

	// Verify token works and get username
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return "", ""
	}
	defer resp.Body.Close()

	var user struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", ""
	}

	return token, user.Login
}

// saveTokenQuiet saves token without the verification step (already verified)
func saveTokenQuiet(token string) error {
	tokenPath := getTokenPath()

	if err := os.MkdirAll(filepath.Dir(tokenPath), 0700); err != nil {
		return err
	}

	if err := os.WriteFile(tokenPath, []byte(token), 0600); err != nil {
		return err
	}

	fmt.Printf("  Token saved to %s\n", tokenPath)
	return nil
}

func saveToken(token string) error {
	tokenPath := getTokenPath()

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(tokenPath), 0700); err != nil {
		return err
	}

	// Write token with restricted permissions
	if err := os.WriteFile(tokenPath, []byte(token), 0600); err != nil {
		return err
	}

	// Verify token works
	req, _ := http.NewRequest("GET", authBaseURL+"/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to verify token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.Remove(tokenPath) // Clean up invalid token
		return fmt.Errorf("invalid token")
	}

	var user struct {
		GitHubLogin string `json:"github_login"`
	}
	json.NewDecoder(resp.Body).Decode(&user)

	fmt.Printf("\n✓ Logged in as @%s\n", user.GitHubLogin)
	fmt.Printf("  Token saved to %s\n", tokenPath)
	return nil
}

func getTokenPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".packs", "token")
}

func getStoredToken() (string, error) {
	data, err := os.ReadFile(getTokenPath())
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	return cmd.Start()
}

// GetAuthToken returns the stored auth token, or empty string if not logged in
func GetAuthToken() string {
	token, _ := getStoredToken()
	return token
}
