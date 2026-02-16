package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// SearchSuggestion represents a single search result suggestion
type SearchSuggestion struct {
	Text string `json:"text"`
	Type string `json:"type"` // "artist/band", "member", "location", "creation date", "first album"
	ID   int    `json:"id"`   // Artist ID to redirect to
}

// Global variable to hold cached artist data
var (
	cachedArtists []Artist
	cacheMutex    sync.RWMutex
)

// SearchHandler handles the /api/search endpoint
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		json.NewEncoder(w).Encode([]SearchSuggestion{})
		return
	}

	suggestions := performSearch(query)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

// performSearch searches through the cached data and returns suggestions
func performSearch(query string) []SearchSuggestion {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	var suggestions []SearchSuggestion
	queryLower := strings.ToLower(query)

	checkMap := make(map[string]bool) // To avoid duplicate suggestions for the same text+type

	for _, artist := range cachedArtists {
		// 1. Artist/Band Name
		if strings.Contains(strings.ToLower(artist.Name), queryLower) {
			addSuggestion(&suggestions, checkMap, artist.Name, "artist/band", artist.ID)
		}

		// 2. Members
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), queryLower) {
				// Format: Member Name (Artist Name)
				text := fmt.Sprintf("%s (%s)", member, artist.Name)
				addSuggestion(&suggestions, checkMap, text, "member", artist.ID)
			}
		}

		// 3. Locations
		for _, concert := range artist.Concerts {
			// Check City
			if strings.Contains(strings.ToLower(concert.City), queryLower) {
				text := fmt.Sprintf("%s (%s)", concert.City, artist.Name)
				addSuggestion(&suggestions, checkMap, text, "location", artist.ID)
			}
			// Check Country
			if strings.Contains(strings.ToLower(concert.Country), queryLower) {
				text := fmt.Sprintf("%s (%s)", concert.Country, artist.Name)
				addSuggestion(&suggestions, checkMap, text, "location", artist.ID)
			}
			// Check Full Location (City, Country)
			fullLoc := fmt.Sprintf("%s, %s", concert.City, concert.Country)
			if strings.Contains(strings.ToLower(fullLoc), queryLower) {
				text := fmt.Sprintf("%s (%s)", fullLoc, artist.Name)
				addSuggestion(&suggestions, checkMap, text, "location", artist.ID)
			}
		}

		// 4. Creation Date
		creationDateStr := strconv.Itoa(artist.Year)
		if strings.Contains(creationDateStr, queryLower) {
			text := fmt.Sprintf("%s (%s)", creationDateStr, artist.Name)
			addSuggestion(&suggestions, checkMap, text, "creation date", artist.ID)
		}

		// 5. First Album Date
		if strings.Contains(artist.Album, queryLower) {
			text := fmt.Sprintf("%s (%s)", artist.Album, artist.Name)
			addSuggestion(&suggestions, checkMap, text, "first album", artist.ID)
		}
	}

	return suggestions
}

func addSuggestion(suggestions *[]SearchSuggestion, checkMap map[string]bool, text, typeStr string, id int) {
	key := fmt.Sprintf("%s-%s-%d", text, typeStr, id)
	if !checkMap[key] {
		*suggestions = append(*suggestions, SearchSuggestion{
			Text: text,
			Type: typeStr,
			ID:   id,
		})
		checkMap[key] = true
	}
}
