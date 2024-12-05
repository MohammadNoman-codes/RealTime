package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"forum/models"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type sResponse struct {
	Message string `json:"message"`
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON request body into the LoginRequest struct
	var loginReq LoginRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON format: %v", err), http.StatusBadRequest)
		return
	}

	// Authenticate the user with the provided credentials
	authenticated, err := models.AuthenticateUser(loginReq.Email, loginReq.Password, w)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication error: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare the response
	var response sResponse
	if authenticated {
		response.Message = "Login successful"
		w.WriteHeader(http.StatusOK)
	} else {
		response.Message = "Invalid credentials"
		w.WriteHeader(http.StatusUnauthorized)
	}

	// Encode the response as JSON and send it
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
