package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Server struct {
	db *sql.DB
}

type Pong struct {
	ID        int    `json:"id"`
	FromUser  string `json:"from_user"`
	ToUser    string `json:"to_user"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

func main() {
	db, err := sql.Open("sqlite", "pongs.db") // TODO: Make database path configurable
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Create tables - simplified schema
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS pongs (
			id INTEGER PRIMARY KEY,
			from_username TEXT,
			to_username TEXT,
			message TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	server := &Server{db: db}

	// Start background cleanup goroutine
	go server.startCleanupScheduler()

	http.HandleFunc("/pong", server.handleSendPong)
	http.HandleFunc("/pongs", server.handleGetPongs)
	http.HandleFunc("/clear", server.handleClearOld)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Background cleanup scheduled to run every 24 hours")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (s *Server) handleSendPong(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ToUser  string `json:"to_user"`
		Message string `json:"message"`
		Token   string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(req.Message) > 140 {
		http.Error(w, "Message too long", http.StatusBadRequest)
		return
	}

	// Validate token and get sender's username
	fromUser, err := validateGitHubToken(req.Token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Simple storage - just store sender username directly
	// Delete any existing pong from this sender to this recipient
	_, err = s.db.Exec(`
		DELETE FROM pongs WHERE from_username = ? AND to_username = ?
	`, fromUser, req.ToUser)
	if err != nil {
		log.Printf("Error deleting existing pong: %v", err)
	}

	// Insert new pong
	_, err = s.db.Exec(`
		INSERT INTO pongs (from_username, to_username, message) 
		VALUES (?, ?, ?)
	`, fromUser, req.ToUser, req.Message)
	if err != nil {
		http.Error(w, "Failed to send pong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
}

func (s *Server) handleGetPongs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// Extract token from "Bearer <token>" format
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
		return
	}

	// Validate token and get username  
	username, err := validateGitHubToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Get pongs for this user
	rows, err := s.db.Query(`
		SELECT id, from_username, to_username, message, created_at
		FROM pongs 
		WHERE to_username = ?
		ORDER BY created_at DESC
	`, username)
	if err != nil {
		http.Error(w, "Failed to get pongs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pongs []Pong
	var pongIDs []int

	for rows.Next() {
		var p Pong
		err := rows.Scan(&p.ID, &p.FromUser, &p.ToUser, &p.Message, &p.CreatedAt)
		if err != nil {
			log.Printf("Error scanning pong: %v", err)
			continue
		}
		pongs = append(pongs, p)
		pongIDs = append(pongIDs, p.ID)
	}

	// Immediately delete the pongs we just fetched (ephemeral!)
	if len(pongIDs) > 0 {
		for _, id := range pongIDs {
			s.db.Exec("DELETE FROM pongs WHERE id = ?", id)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pongs)
}

func (s *Server) handleClearOld(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Use the same cleanup logic as background task
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	result, err := s.db.Exec("DELETE FROM pongs WHERE created_at < ?", sevenDaysAgo)
	if err != nil {
		http.Error(w, "Failed to clear old pongs", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Manual cleanup: removed %d old pongs", rowsAffected)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"cleared": rowsAffected})
}

// validateGitHubToken validates a GitHub token and returns the username
func validateGitHubToken(token string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("User-Agent", "pong-server")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var user struct {
		Login string `json:"login"`
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return "", err
	}

	return user.Login, nil
}

// startCleanupScheduler runs cleanup every 24 hours in background
func (s *Server) startCleanupScheduler() {
	// Run initial cleanup on startup
	s.cleanupOldPongs()
	
	// Schedule cleanup every 24 hours
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		s.cleanupOldPongs()
	}
}

// cleanupOldPongs removes pongs older than 7 days
func (s *Server) cleanupOldPongs() {
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	result, err := s.db.Exec("DELETE FROM pongs WHERE created_at < ?", sevenDaysAgo)
	if err != nil {
		log.Printf("Cleanup error: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("Cleanup: removed %d old pongs", rowsAffected)
	}
}