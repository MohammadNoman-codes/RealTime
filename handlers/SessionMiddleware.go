package handlers

import (
	"database/sql"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SessionMiddleware checks for session validity before allowing access to protected routes
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			// No session cookie, redirect to login
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		sessionID := cookie.Value

		db, err := sql.Open("sqlite3", "./storage/storage.db")
		if err != nil {
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		var expiresAt time.Time
		err = db.QueryRow("SELECT expires_at FROM sessions WHERE session_id = ?", sessionID).Scan(&expiresAt)
		if err != nil {
			if err == sql.ErrNoRows {
				// Session not found, redirect to login
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}

		if time.Now().After(expiresAt) {
			// Session has expired, log out the user
			LogoutHandler(w, r)
			return
		}

		// Session is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
