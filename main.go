package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	// Initialize cached data
	var err error
	fmt.Println("Fetching and caching artist data...")
	cachedArtists, err = getFullArtistData()
	if err != nil {
		log.Fatalf("Failed to fetch initial data: %v", err)
	}
	fmt.Println("Data cached successfully.")

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "..") {
			renderError(w, "Invalid path", http.StatusBadRequest)
			return
		}
		fs.ServeHTTP(w, r)
	})))

	// Register handlers
	mux.HandleFunc("/", homePage)
	mux.HandleFunc("/artist/", artistDetailHandler)
	mux.HandleFunc("/api/search", SearchHandler) // Register search handler

	// Wrap with logging middleware
	loggedMux := RequestHandler(mux)

	// Configure server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      loggedMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Println("Server started on http://localhost:8080")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
