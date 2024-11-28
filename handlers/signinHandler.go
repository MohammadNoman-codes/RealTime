package handlers

import (
	"fmt"
	"net/http"

	"forum/models"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	authenticated, err := models.AuthenticateUser(email, password, w)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication error: %v", err), http.StatusInternalServerError)
		return
	}

	if !authenticated {
		// Handle invalid login (e.g., show error message)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	// Redirect to the home page or user dashboard
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
