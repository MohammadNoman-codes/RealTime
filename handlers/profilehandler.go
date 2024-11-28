package handlers

import (
	"database/sql"
	"forum/models"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type ProfilePost struct {
	ID      int
	Title   string
	Content string
}

type ProfileData struct {
	Username string
	Posts    []ProfilePost
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve user ID from session
	userID, err := models.GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Please log in to view your profile", http.StatusUnauthorized)
		return
	}

	// Connect to the database
	db, err := sql.Open("sqlite3", "./storage/storage.db")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query to get user posts
	rows, err := db.Query(`
		SELECT id, title, content 
		FROM posts 
		WHERE user_id = ?`, userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Retrieve posts
	var posts []ProfilePost
	for rows.Next() {
		var post ProfilePost
		if err := rows.Scan(&post.ID, &post.Title, &post.Content); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	// Query to get username
	var username string
	err = db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create profile data
	profileData := ProfileData{
		Username: username,
		Posts:    posts,
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, profileData)
}
