package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func SetupDatabase() error {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./storage/storage.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Read setup.sql file
	setupSQL, err := os.ReadFile("./storage/setup.sql")
	if err != nil {
		return err
	}

	// Execute SQL statements
	_, err = db.Exec(string(setupSQL))
	if err != nil {
		return err
	}

	return nil
}

func CreateUser(email, username, password string) error {
	db, err := sql.Open("sqlite3", "./storage/storage.db")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users (email, username, password) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(email, username, password)
	if err != nil {
		return err
	}

	return nil
}

func Dublicate(db *sql.DB, userID int) error {
	// Query to find an active session for the user
	var sessionID string
	var expiryTime time.Time
	err := db.QueryRow("SELECT session_id, expires_at FROM sessions WHERE user_id = ? AND expires_at > ?", userID, time.Now()).Scan(&sessionID, &expiryTime)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing session: %v", err)
	}

	// If an active session exists, delete it
	if sessionID != "" && expiryTime.After(time.Now()) {
		_, err := db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
		if err != nil {
			return fmt.Errorf("failed to delete existing session: %v", err)
		}
	}

	return nil
}

// AuthenticateUser checks the user's credentials and creates a session if valid.
func AuthenticateUser(email, password string, w http.ResponseWriter) (bool, error) {
	db, err := sql.Open("sqlite3", "./storage/storage.db")
	if err != nil {
		return false, err
	}
	defer db.Close()

	var userID int
	var storedPassword string
	err = db.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&userID, &storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // User not found
		}
		return false, err
	}

	// Compare passwords (in a real application, use bcrypt or similar for secure password storage)
	if storedPassword != password {
		return false, nil // Passwords do not match
	}

	// Call the Dublicate function to handle any existing active sessions
	err = Dublicate(db, userID)
	if err != nil {
		return false, fmt.Errorf("failed to handle duplicate session: %v", err)
	}

	// Generate a session ID
	sessionID := generateSessionID()

	// Set the session ID as a cookie
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(24 * time.Hour), // Cookie expires in 24 hours
	}
	http.SetCookie(w, cookie)

	// Store the session in the database
	_, err = db.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)", sessionID, userID, time.Now().Add(24*time.Hour))
	if err != nil {
		return false, err
	}

	return true, nil // Authentication and session creation successful
}
func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func GetUserIDFromSession(r *http.Request) (int, error) {
	// Retrieve the session ID from the cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			return 0, fmt.Errorf("session not found, please log in")
		}
		return 0, fmt.Errorf("failed to retrieve session: %v", err)
	}
	sessionID := cookie.Value

	// Open the database
	db, err := sql.Open("sqlite3", "./storage/storage.db")
	if err != nil {
		return 0, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Query the database to get the user ID from the session ID
	var userID int
	err = db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", sessionID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("session not found in the database")
		}
		return 0, fmt.Errorf("failed to retrieve session information: %v", err)
	}

	return userID, nil
}

func HasUserLikedPost(userID, postID int) (bool, error) {
	db, err := sql.Open("sqlite3", "./storage/storage.db")
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
