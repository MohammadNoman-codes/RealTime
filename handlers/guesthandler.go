package handlers

import (
	"html/template"
	"net/http"
)

// GuestPageHandler handles displaying the posts for guests.
func GuestPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "405: Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/guest" {
		t, err := template.ParseFiles("templates/error.html")
		if err != nil {
			http.Error(w, "500: Internal Server Error This is the error", http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
		return
	}

	category := r.URL.Query().Get("category")

	t, err := template.ParseFiles("templates/guestpage.html")
	if err != nil {
		http.Error(w, "500: Internal Server Error this is ", http.StatusInternalServerError)
		return
	}

	// Use a dummy user ID for fetching posts without user-specific like/dislike data
	posts, err := FetchPosts(0, category)
	if err != nil {
		http.Error(w, "500: Internal Server Error maybe this one ", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Posts":            posts,
		"SelectedCategory": category,
	}

	if len(posts) == 0 {
		data["NoPosts"] = true
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
