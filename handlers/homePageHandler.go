package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/models"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	Category        string `json:"category"`
	LikesCount      int    `json:"likes_count"`
	DislikesCount   int    `json:"dislikes_count"`
	UserHasLiked    bool   `json:"user_has_liked"`
	UserHasDisliked bool   `json:"user_has_disliked"`
}

type Response struct {
	Message string `json:"message"`
	Posts   []Post `json:"posts,omitempty"`
}

// FetchPosts fetches posts based on the selected category.
func FetchPosts(userID int, category string) ([]Post, error) {
	db, err := sql.Open("sqlite3", "storage/storage.db")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	var rows *sql.Rows
	if category == "" {
		// Fetch all posts
		rows, err = db.Query(`
            SELECT posts.id, posts.title, posts.content, 
                   (SELECT COUNT(*) FROM likes WHERE post_id = posts.id) AS likes_count, 
                   (SELECT COUNT(*) FROM dislikes WHERE post_id = posts.id) AS dislikes_count,
                   EXISTS (SELECT 1 FROM likes WHERE post_id = posts.id AND user_id = ?) AS user_has_liked,
                   EXISTS (SELECT 1 FROM dislikes WHERE post_id = posts.id AND user_id = ?) AS user_has_disliked
            FROM posts`, userID, userID)
	} else if category == "liked" {
		// Fetch posts that the user has liked
		rows, err = db.Query(`
            SELECT posts.id, posts.title, posts.content, 
                   (SELECT COUNT(*) FROM likes WHERE post_id = posts.id) AS likes_count, 
                   (SELECT COUNT(*) FROM dislikes WHERE post_id = posts.id) AS dislikes_count,
                   EXISTS (SELECT 1 FROM likes WHERE post_id = posts.id AND user_id = ?) AS user_has_liked,
                   EXISTS (SELECT 1 FROM dislikes WHERE post_id = posts.id AND user_id = ?) AS user_has_disliked
            FROM posts 
            JOIN likes ON posts.id = likes.post_id 
            WHERE likes.user_id = ?`, userID, userID, userID)
	} else {
		// Fetch posts by category
		rows, err = db.Query(`
            SELECT posts.id, posts.title, posts.content, 
                   (SELECT COUNT(*) FROM likes WHERE post_id = posts.id) AS likes_count, 
                   (SELECT COUNT(*) FROM dislikes WHERE post_id = posts.id) AS dislikes_count,
                   EXISTS (SELECT 1 FROM likes WHERE post_id = posts.id AND user_id = ?) AS user_has_liked,
                   EXISTS (SELECT 1 FROM dislikes WHERE post_id = posts.id AND user_id = ?) AS user_has_disliked
            FROM posts WHERE category = ?`, userID, userID, category)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.LikesCount, &post.DislikesCount, &post.UserHasLiked, &post.UserHasDisliked)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return posts, nil
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405: Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON request body into a map
	var request struct {
		Category string `json:"category"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON format: %v", err), http.StatusBadRequest)
		return
	}

	// Retrieve the userID from the session
	userID, err := models.GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Please log in to view posts", http.StatusUnauthorized)
		return
	}

	// Fetch posts based on the category
	posts, err := FetchPosts(userID, "")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare the response
	var response Response
	if len(posts) > 0 {
		response.Message = "Posts fetched successfully"
		response.Posts = posts
	} else {
		response.Message = "No posts available"
	}

	// Set response header and encode the response to JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
