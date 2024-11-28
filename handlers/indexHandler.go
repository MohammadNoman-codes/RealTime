package handlers

import (
	"html/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method == "GET" {
		// Serve error page if the request path is not "/"
		if r.URL.Path != "/" && r.URL.Path != "/login" {
			t, _ := template.ParseFiles("templates/error.html")
			t.Execute(w, http.StatusNotFound)
			return
		}

		// Parse the index template
		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "500: internal server error", http.StatusInternalServerError)
			return
		}

		// Execute the index template
		err = t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Return an error for invalid request methods
		http.Error(w, "400: Bad Request", http.StatusBadRequest)
		return
	}
}
