package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// GuestComment represents a comment on a post with like/dislike counts.
type GuestComment struct {
	ID        int
	Content   string
	UserID    int
	PostID    int
	Username  string
	CreatedAt time.Time
	CLikes    int
	CDislikes int
}

// GuestPostWithComments contains a post's details along with its comments.
type GuestPostWithComments struct {
	PostID    int
	PostTitle string
	Comments  []GuestComment
}

// GuestCommentsHandler displays comments for a post without allowing new comments to be added.
func GuestCommentsHandler(w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the query parameters
	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "400: Bad Request - post_id missing", http.StatusBadRequest)
		return
	}

	// Open a connection to the database
	db, err := sql.Open("sqlite3", "./storage/storage.db")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Retrieve the title of the post
	var postTitle string
	err = db.QueryRow("SELECT title FROM posts WHERE id = ?", postID).Scan(&postTitle)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "404: Not Found - Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Query to retrieve comments and their like/dislike counts
	rows, err := db.Query(`
		SELECT c.id, c.content, c.user_id, u.username, c.post_id, c.created_at,
			   (SELECT COUNT(*) FROM likes WHERE comment_id = c.id) AS comment_likes,
			   (SELECT COUNT(*) FROM dislikes WHERE comment_id = c.id) AS comment_dislikes
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?`, postID)
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Populate comments into the struct
	var comments []GuestComment
	for rows.Next() {
		var comment GuestComment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.Username, &comment.PostID, &comment.CreatedAt, &comment.CLikes, &comment.CDislikes); err != nil {
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	// Convert postID to an integer
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Struct to pass data to the template
	postWithComments := GuestPostWithComments{
		PostID:    postIDInt,
		PostTitle: postTitle,
		Comments:  comments,
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("./templates/guestcomments.html")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, postWithComments)
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
}
