package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show your local pong history",
	Run: func(cmd *cobra.Command, args []string) {
		history, err := LoadHistory()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading history: %v\n", err)
			os.Exit(1)
		}

		if len(history) == 0 {
			fmt.Println("No pong history found.")
			return
		}

		fmt.Printf("Your pong history (%d pongs):\n\n", len(history))

		// Show most recent first
		for i := len(history) - 1; i >= 0; i-- {
			pong := history[i]
			timeStr := pong.FetchedAt.Format("Jan 2, 3:04 PM")
			fmt.Printf("📝 %s (%s): %s\n", pong.FromUser, timeStr, pong.Message)
		}
	},
}

var clearHistoryCmd = &cobra.Command{
	Use:   "clear-history",
	Short: "Clear your local pong history",
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		historyFile := fmt.Sprintf("%s/.pong/history.json", homeDir)
		err = os.Remove(historyFile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No history file to clear.")
				return
			}
			fmt.Fprintf(os.Stderr, "Error clearing history: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✓ History cleared!")
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
	rootCmd.AddCommand(clearHistoryCmd)
}