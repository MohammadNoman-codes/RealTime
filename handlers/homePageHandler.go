package handlers

import (
	"database/sql"
	"fmt"
	"forum/models"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	ID              int
	Title           string
	Content         string
	Category        string
	LikesCount      int
	DislikesCount   int
	UserHasLiked    bool
	UserHasDisliked bool
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

//This is what we have to do when we have to work for a project
//This is not what I have signed up for
// Come on we need to finish this
// I am done with coding 
// I dont want to do this 


func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "405: Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/home" {
		t, err := template.ParseFiles("templates/error.html")
		if err != nil {
			http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
		return
	}

	category := r.URL.Query().Get("category")

	t, err := template.ParseFiles("templates/homePage.html")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	userID, err := models.GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Please log in to view posts", http.StatusUnauthorized)
		return
	}

	posts, err := FetchPosts(userID, category)
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Posts":            posts,
		"SelectedCategory": category,
	}

	if len(posts) == 0 {
		data["NoPosts"] = true
	}

	//This is when we excetute the things
	//I am almost done with this
	//I am almost done with this
	// Please we need to finish this
	// I am done with campus
	// Please I need to go home
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
