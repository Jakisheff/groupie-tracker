package main

import (
	"fmt"
	"net/http"
)

// CustomError represents a custom error with an HTTP status code and message
type CustomError struct {
	StatusCode int
	Message    string
}

// Error implements the error interface
func (e *CustomError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

// NewBadRequestError creates a 400 Bad Request error
func NewBadRequestError(message string) error {
	return &CustomError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

// NewNotFoundError creates a 404 Not Found error
func NewNotFoundError(message string) error {
	return &CustomError{
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}

// NewInternalServerError creates a 500 Internal Server Error
func NewInternalServerError(message string) error {
	return &CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
	}
}

// renderError renders the error page using the template
func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := struct {
		StatusCode int
		Message    string
	}{
		StatusCode: statusCode,
		Message:    message,
	}
	err := tmpl.ExecuteTemplate(w, "error.html", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("HTTP %d: %s", http.StatusInternalServerError, "Failed to render error page"), http.StatusInternalServerError)
	}
}
