package handlers

import (
	"fmt"
	"forum/models"
	"net/http"
	"strings"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	if email == "" || username == "" || password == "" {
		http.Error(w, "can not have empty fields please enter the data into it", http.StatusBadRequest)
		return
	}

	err := models.CreateUser(email, username, password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	// Redirect to index.html after successful registration
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
