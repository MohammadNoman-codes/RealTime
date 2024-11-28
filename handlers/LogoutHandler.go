package handlers

import (
	"database/sql"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session ID from the cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	sessionID := cookie.Value

	// Delete the session from the database
	db, err := sql.Open("sqlite3", "./storage/storage.db")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Delete the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0), // Expire the cookie immediately
	})

	// Redirect to the login page after logout
	http.Redirect(w, r, "/", http.StatusFound)
}
