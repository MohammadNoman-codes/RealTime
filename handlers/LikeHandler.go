package handlers

import (
	"database/sql"
	"forum/models"
	"log"
	"net/http"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405: Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	postID := r.FormValue("post_id")
	userID, err := models.GetUserIDFromSession(r) // Get the logged-in user ID
	if err != nil {
		http.Error(w, "Please log in to perform this action", http.StatusUnauthorized)
		return
	}

	db, err := sql.Open("sqlite3", "storage/storage.db")
	if err != nil {
		log.Println("Error opening database:", err)
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Remove dislike if it exists
	_, err = tx.Exec("DELETE FROM dislikes WHERE user_id = ? AND post_id = ?", userID, postID)
	if err != nil {
		log.Println("Error deleting dislike:", err)
		http.Error(w, "Failed to remove dislike", http.StatusInternalServerError)
		return
	}

	// Insert like
	_, err = tx.Exec("INSERT INTO likes (user_id, post_id) VALUES (?, ?)", userID, postID)
	if err != nil {
		log.Println("Error inserting like:", err)
		http.Error(w, "Failed to add like", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
