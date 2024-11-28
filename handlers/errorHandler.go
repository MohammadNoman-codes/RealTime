package handlers

import (
	"html/template"
	"net/http"
)

type CustomError struct {
	Code      int
	Message   string
	ErrorType string
}

func HandleErrorPage(w http.ResponseWriter, r *http.Request) {
	// The error html page does not show the error message,
	err := &CustomError{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error",
	}
	renderErrorPage(w, err)
}

func renderErrorPage(w http.ResponseWriter, er *CustomError) {
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, er)
	if err != nil {
		http.Error(w, er.Message, er.Code)
		return
	}
}
