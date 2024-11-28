package handlers

import (
	"database/sql"
	"forum/models"
	"net/http"
)

func UndislikeHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "405: Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    postID := r.FormValue("post_id")
    userID, err := models.GetUserIDFromSession(r)  // Get the logged-in user ID
    if err != nil {
        http.Error(w, "Please log in to perform this action", http.StatusUnauthorized)
        return
    }

    db, err := sql.Open("sqlite3", "storage/storage.db")
    if err != nil {
        http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    _, err = db.Exec("DELETE FROM dislikes WHERE user_id = ? AND post_id = ?", userID, postID)
    if err != nil {
        http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/home", http.StatusSeeOther)
}