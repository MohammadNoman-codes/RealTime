package handlers

import (
	"database/sql"
	"fmt"
	"forum/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Comment struct {
	ID               int
	Content          string
	UserID           int
	PostID           int
	Username         string
	CreatedAt        time.Time
	CLikes           int
	CDislikes        int
	UserHasCLiked    bool
	UserHasCDisliked bool
}

type PostWithComments struct {
	PostID          int
	PostTitle       string
	Comments        []Comment
	UserID          int
	UserHasLiked    bool
	UserHasDisliked bool
	LikesCount      int
	DislikesCount   int
}

// FetchComments fetches comments for a post along with like/dislike status.
func FetchComments(postID, userID int) (PostWithComments, error) {
	db, err := sql.Open("sqlite3", "./storage/storage.db")
	if err != nil {
		return PostWithComments{}, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	var post PostWithComments

	// Fetch post details with like and dislike counts
	err = db.QueryRow(`
		SELECT p.id, p.title, 
		       (SELECT COUNT(*) FROM likes WHERE post_id = p.id) AS likes_count, 
		       (SELECT COUNT(*) FROM dislikes WHERE post_id = p.id) AS dislikes_count,
		       EXISTS (SELECT 1 FROM likes WHERE post_id = p.id AND user_id = ?) AS user_has_liked,
		       EXISTS (SELECT 1 FROM dislikes WHERE post_id = p.id AND user_id = ?) AS user_has_disliked
		FROM posts p
		WHERE p.id = ?`, userID, userID, postID).Scan(&post.PostID, &post.PostTitle, &post.LikesCount, &post.DislikesCount, &post.UserHasLiked, &post.UserHasDisliked)

	if err != nil {
		return post, fmt.Errorf("failed to fetch post details: %v", err)
	}

	// Fetch comments for the post
	rows, err := db.Query(`
		SELECT c.id, c.content, c.user_id, u.username, c.created_at,
		       (SELECT COUNT(*) FROM likes WHERE comment_id = c.id) AS comment_likes,
		       (SELECT COUNT(*) FROM dislikes WHERE comment_id = c.id) AS comment_dislikes,
		       EXISTS (SELECT 1 FROM likes WHERE comment_id = c.id AND user_id = ?) AS user_has_liked,
		       EXISTS (SELECT 1 FROM dislikes WHERE comment_id = c.id AND user_id = ?) AS user_has_disliked
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at DESC`, userID, userID, postID)

	if err != nil {
		return post, fmt.Errorf("failed to fetch comments: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.Username, &comment.CreatedAt, &comment.CLikes, &comment.CDislikes, &comment.UserHasCLiked, &comment.UserHasCDisliked)
		if err != nil {
			return post, fmt.Errorf("failed to scan comment row: %v", err)
		}
		post.Comments = append(post.Comments, comment)
	}

	if err = rows.Err(); err != nil {
		return post, fmt.Errorf("row iteration error: %v", err)
	}

	return post, nil
}

// CommentsHandler handles the display of a single post and its comments.
func CommentsHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("post_id")
	if postIDStr == "" {
		http.Error(w, "400: Bad Request - post_id missing", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "400: Bad Request - invalid post_id", http.StatusBadRequest)
		return
	}

	userID, err := models.GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Please log in to view comments", http.StatusUnauthorized)
		return
	}

	postWithComments, err := FetchComments(postID, userID)
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("./templates/comments.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, postWithComments)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	userID, err := models.GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "403: Forbidden", http.StatusForbidden)
		return
	}
	postID := r.FormValue("post_id")
	content := strings.TrimSpace(r.FormValue("content"))

	if content == "" {
		http.Error(w, "400: Bad Request - Content cannot be empty", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./storage/storage.db")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO comments (content, user_id, post_id, created_at) VALUES (?, ?, ?, ?)", content, userID, postID, time.Now())
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect back to the comments page
	http.Redirect(w, r, "/comments?post_id="+postID, http.StatusSeeOther)
}
