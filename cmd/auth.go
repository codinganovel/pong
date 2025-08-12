package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type GitHubUser struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

// GetStoredToken returns the saved GitHub token
func GetStoredToken() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	tokenFile := filepath.Join(homeDir, ".pong", "token")
	token, err := os.ReadFile(tokenFile)
	if err != nil {
		return "", fmt.Errorf("not logged in. Run 'pong login' first")
	}

	return string(token), nil
}

// ValidateGitHubUsername checks if a username exists on GitHub
func ValidateGitHubUsername(username, token string) error {
	req, err := http.NewRequest("GET", "https://api.github.com/users/"+username, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("User-Agent", "pong-cli")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return fmt.Errorf("GitHub user '%s' not found", username)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	return nil
}

const serverURL = "YOUR_PONG_SERVER_URL" // e.g., "http://localhost:8080" for local development