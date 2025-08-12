package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type HistoryEntry struct {
	FromUser  string    `json:"from_user"`
	Message   string    `json:"message"`
	FetchedAt time.Time `json:"fetched_at"`
}

// SaveToHistory saves pongs to local history file
func SaveToHistory(pongs []struct {
	FromUser  string `json:"from_user"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}) error {
	if len(pongs) == 0 {
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	pongDir := filepath.Join(homeDir, ".pong")
	os.MkdirAll(pongDir, 0755)
	
	historyFile := filepath.Join(pongDir, "history.json")

	// Read existing history
	var history []HistoryEntry
	if data, err := os.ReadFile(historyFile); err == nil {
		json.Unmarshal(data, &history)
	}

	// Add new pongs to history
	fetchTime := time.Now()
	for _, pong := range pongs {
		history = append(history, HistoryEntry{
			FromUser:  pong.FromUser,
			Message:   pong.Message,
			FetchedAt: fetchTime,
		})
	}

	// Save updated history
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(historyFile, data, 0600)
}

// LoadHistory loads pongs from local history file
func LoadHistory() ([]HistoryEntry, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	historyFile := filepath.Join(homeDir, ".pong", "history.json")
	data, err := os.ReadFile(historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []HistoryEntry{}, nil // Return empty history if file doesn't exist
		}
		return nil, err
	}

	var history []HistoryEntry
	err = json.Unmarshal(data, &history)
	return history, err
}