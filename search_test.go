package main

import (
	"testing"
)

func TestPerformSearch(t *testing.T) {
	// Setup mock data
	mockArtists := []Artist{
		{
			ID:   1,
			Name: "Queen",
			Members: []string{
				"Freddie Mercury",
				"Brian May",
				"Roger Taylor",
				"John Deacon",
			},
			Year:  1970,
			Album: "14-12-1973",
			Concerts: []Concert{
				{City: "London", Country: "UK"},
				{City: "New York", Country: "USA"},
			},
		},
		{
			ID:   2,
			Name: "The Beatles",
			Members: []string{
				"John Lennon",
				"Paul McCartney",
				"George Harrison",
				"Ringo Starr",
			},
			Year:  1960,
			Album: "22-03-1963",
			Concerts: []Concert{
				{City: "Liverpool", Country: "UK"},
				{City: "Hamburg", Country: "Germany"},
			},
		},
	}

	// Lock the mutex and set cachedArtists
	cacheMutex.Lock()
	cachedArtists = mockArtists
	cacheMutex.Unlock()

	tests := []struct {
		name     string
		query    string
		expected int // Number of expected suggestions
	}{
		{"Artist Match", "Queen", 1},
		{"Artist Case Insensitive", "queen", 1},
		{"Member Match", "Freddie", 1},
		{"Location Match", "London", 2}, // Matches "London (Queen)" and "London, UK (Queen)"
		{"Creation Date Match", "1970", 1},
		{"First Album Match", "1973", 1},
		{"No Match", "Justin Bieber", 0},
		{"Multiple Matches", "John", 2}, // John Deacon (Queen) and John Lennon (Beatles)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := performSearch(tt.query)
			if len(suggestions) != tt.expected {
				t.Errorf("expected %d suggestions, got %d for query '%s'", tt.expected, len(suggestions), tt.query)
			}
		})
	}
}
