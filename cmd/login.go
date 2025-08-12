package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cli/oauth"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with GitHub",
	Run: func(cmd *cobra.Command, args []string) {
		host, err := oauth.NewGitHubHost("https://github.com")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		flow := &oauth.Flow{
			Host:         host,
			ClientID:     "YOUR_GITHUB_OAUTH_CLIENT_ID",
			ClientSecret: "YOUR_GITHUB_OAUTH_CLIENT_SECRET",
			CallbackURI:  "YOUR_OAUTH_CALLBACK_URI", // e.g., "http://127.0.0.1:8080/callback"
			Scopes:       []string{"read:user"},
		}

		fmt.Println("Starting GitHub authentication...")
		accessToken, err := flow.DetectFlow()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Authentication failed: %v\n", err)
			os.Exit(1)
		}

		// Save token to ~/.pong/token
		homeDir, _ := os.UserHomeDir()
		pongDir := filepath.Join(homeDir, ".pong")
		os.MkdirAll(pongDir, 0755)
		
		tokenFile := filepath.Join(pongDir, "token")
		err = os.WriteFile(tokenFile, []byte(accessToken.Token), 0600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save token: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✓ Authentication successful!")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}