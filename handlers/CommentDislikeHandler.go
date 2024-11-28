package handlers

import (
	"database/sql"
	"forum/models"
	"log"
	"net/http"
)

func CommentDislikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405: Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	commentID := r.FormValue("ID") // Use "ID" for comment identifier
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

	// Remove like if it exists
	_, err = tx.Exec("DELETE FROM likes WHERE user_id = ? AND comment_id = ?", userID, commentID)
	if err != nil {
		log.Println("Error deleting like:", err)
		http.Error(w, "Failed to remove like", http.StatusInternalServerError)
		return
	}

	// Insert dislike
	_, err = tx.Exec("INSERT INTO dislikes (user_id, comment_id) VALUES (?, ?)", userID, commentID)
	if err != nil {
		log.Println("Error inserting dislike:", err)
		http.Error(w, "Failed to add dislike", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/comments?post_id="+postID, http.StatusSeeOther)
}
