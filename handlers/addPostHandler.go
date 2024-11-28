package handlers

import (
	"database/sql"
	"fmt"
	"forum/models"
	"html/template"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func AddPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/addpost" {
		http.Error(w, "404: Page Not Found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Display the add post form
		t, err := template.ParseFiles("templates/addPost.html")
		if err != nil {
			http.Error(w, "500: Internal Server Error (Parsing Template)", http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
	case http.MethodPost:
		// Process the form submission
		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))
		category := strings.TrimSpace(r.FormValue("category"))

		if title == "" || content == "" || category == "" {
			http.Error(w, "400: Bad Request (Title, Content, or Category Missing)", http.StatusBadRequest)
			return
		}

		err := addPostToDB(r, title, content, category)
		if err != nil {
			http.Error(w, "500: Internal Server Error (Adding Post to DB)", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	default:
		http.Error(w, "405: Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func addPostToDB(r *http.Request, title, content, category string) error {
	user, err := models.GetUserIDFromSession(r)
	if err != nil {
		return fmt.Errorf("failed to retrieve user ID: %v", err)
	}
	db, err := sql.Open("sqlite3", "storage/storage.db")
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO posts (title, content, category, user_id) VALUES (?, ?, ?, ?)", title, content, category, user)
	if err != nil {
		return fmt.Errorf("failed to insert post: %v", err)
	}

	return nil
}
