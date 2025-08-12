package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pong",
	Short: "Leave a note. Check your notes. That's it.",
	Long:  `Pong is the digital equivalent of leaving a sticky note on someone's desk. Simple, ephemeral, and pressure-free.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get stored token
		token, err := GetStoredToken()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Fetch pongs from server
		req, err := http.NewRequest("GET", serverURL+"/pongs", nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create request: %v\n", err)
			os.Exit(1)
		}

		req.Header.Set("Authorization", "Bearer "+token)
		
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to fetch pongs: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fmt.Fprintf(os.Stderr, "Server error: %s\n", string(body))
			os.Exit(1)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read response: %v\n", err)
			os.Exit(1)
		}

		var pongs []struct {
			FromUser  string `json:"from_user"`
			Message   string `json:"message"`
			CreatedAt string `json:"created_at"`
		}

		err = json.Unmarshal(body, &pongs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse response: %v\n", err)
			os.Exit(1)
		}

		if len(pongs) == 0 {
			fmt.Println("No pongs waiting for you!")
			return
		}

		// Save to local history before displaying
		err = SaveToHistory(pongs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to save history: %v\n", err)
		}

		fmt.Printf("You have %d pong%s:\n\n", len(pongs), func() string {
			if len(pongs) == 1 { return "" }
			return "s"
		}())

		for _, pong := range pongs {
			fmt.Printf("📝 %s: %s\n", pong.FromUser, pong.Message)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}