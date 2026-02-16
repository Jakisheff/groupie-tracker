package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var (
	tmpl *template.Template
)

// init initializes HTML templates
func init() {
	var err error
	funcMap := template.FuncMap{
		"contains": func(slice []int, item int) bool {
			for _, v := range slice {
				if v == item {
					return true
				}
			}
			return false
		},
		"containsString": func(slice []string, item string) bool {
			for _, v := range slice {
				if v == item {
					return true
				}
			}
			return false
		},
	}
	tmpl, err = template.New("").Funcs(funcMap).ParseFiles(
		"templates/index.html",
		"templates/artist_detail.html",
		"templates/error.html",
	)
	if err != nil {
		log.Fatal("Failed to parse templates:", err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		renderError(w, fmt.Sprintf("%s method not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		renderError(w, fmt.Sprintf("Page not found: %s", r.URL.Path), http.StatusNotFound)
		return
	}

	artists, err := getFullArtistData()
	if err != nil {
		if customErr, ok := err.(*CustomError); ok {
			renderError(w, customErr.Message, customErr.StatusCode)
			return
		}
		renderError(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	if len(artists) == 0 {
		log.Printf("No artist data available")
		renderError(w, "No artist data available", http.StatusInternalServerError)
		return
	}

	// Parse query parameters
	query := r.URL.Query()
	filter := ArtistFilter{
		CreationDateMin: parseInt(query.Get("creationMin")),
		CreationDateMax: parseInt(query.Get("creationMax")),
		FirstAlbumMin:   parseInt(query.Get("firstAlbumMin")),
		FirstAlbumMax:   parseInt(query.Get("firstAlbumMax")),
		Members:         parseInts(query["members"]),
		Locations:       query["locations"],
	}

	// Get filter options from all artists
	filterOptions := GetFilterOptions(artists)

	// Filter artists
	filteredArtists := FilterArtists(artists, filter)

	pageData := PageData{
		Artists:       filteredArtists,
		FilterOptions: filterOptions,
		CurrentFilter: filter,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.ExecuteTemplate(w, "index.html", pageData)
	if err != nil {
		log.Printf("Failed to render template: %v", err)
		renderError(w, "Failed to render page", http.StatusInternalServerError)
	}
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func parseInts(s []string) []int {
	var ints []int
	for _, v := range s {
		ints = append(ints, parseInt(v))
	}
	return ints
}

func artistDetailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		renderError(w, fmt.Sprintf("%s method not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/artist/"):]
	var artistID int
	_, err := fmt.Sscanf(idStr, "%d", &artistID)
	if err != nil || artistID <= 0 {
		renderError(w, "Invalid artist ID format", http.StatusBadRequest)
		return
	}

	artists, err := getFullArtistData()
	if err != nil {
		if customErr, ok := err.(*CustomError); ok {
			log.Printf("Failed to fetch data: %v", err)
			renderError(w, customErr.Message, customErr.StatusCode)
			return
		}
		log.Printf("Failed to fetch data: %v", err)
		renderError(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	if len(artists) == 0 {
		log.Printf("No artist data available")
		renderError(w, "No artist data available", http.StatusInternalServerError)
		return
	}

	var selectedArtist *Artist
	for i := range artists {
		if artists[i].ID == artistID {
			selectedArtist = &artists[i]
			break
		}
	}

	if selectedArtist == nil {
		renderError(w, fmt.Sprintf("Artist with ID %d not found", artistID), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.ExecuteTemplate(w, "artist_detail.html", selectedArtist)
	if err != nil {
		log.Printf("Failed to render template: %v", err)
		renderError(w, "Failed to render page", http.StatusInternalServerError)
	}
}
