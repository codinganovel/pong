package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send username message",
	Short: "Send a pong to someone",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		message := strings.Join(args[1:], " ")
		
		if len(message) > 140 {
			fmt.Fprintf(os.Stderr, "Message too long (%d chars). Max 140 characters.\n", len(message))
			os.Exit(1)
		}

		// Get stored token
		token, err := GetStoredToken()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Validate target username exists
		err = ValidateGitHubUsername(username, token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Send to server
		payload := map[string]string{
			"to_user": username,
			"message": message,
			"token":   token,
		}

		jsonData, _ := json.Marshal(payload)
		resp, err := http.Post(serverURL+"/pong", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to send pong: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			fmt.Fprintf(os.Stderr, "Server error: %s\n", string(body))
			os.Exit(1)
		}

		fmt.Printf("✓ Pong sent to %s: %s\n", username, message)
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}