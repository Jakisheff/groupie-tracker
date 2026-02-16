package main

import (
	"testing"
)

func TestFilterArtists(t *testing.T) {
	artists := []Artist{
		{
			ID:      1,
			Name:    "Queen",
			Year:    1970,
			Album:   "13-07-1973",
			Members: []string{"Freddie", "Brian", "Roger", "John"},
			Concerts: []Concert{
				{City: "London", Country: "UK"},
				{City: "New York", Country: "USA"},
			},
		},
		{
			ID:      2,
			Name:    "Soja",
			Year:    1997,
			Album:   "05-01-2002",
			Members: []string{"Jacob", "Bob", "Patrick", "Ryan", "Ken", "Hellman", "Rafael", "Trevor"},
			Concerts: []Concert{
				{City: "San Diego", Country: "USA"},
				{City: "Arlington", Country: "USA"},
			},
		},
		{
			ID:      3,
			Name:    "Pink Floyd",
			Year:    1965,
			Album:   "04-08-1967",
			Members: []string{"David", "Roger", "Nick", "Richard"},
			Concerts: []Concert{
				{City: "London", Country: "UK"},
				{City: "Rome", Country: "Italy"},
			},
		},
	}

	tests := []struct {
		name     string
		filter   ArtistFilter
		expected int
	}{
		{
			name:     "No filters",
			filter:   ArtistFilter{},
			expected: 3,
		},
		{
			name: "Filter by Creation Date (1960-1969)",
			filter: ArtistFilter{
				CreationDateMin: 1960,
				CreationDateMax: 1969,
			},
			expected: 1, // Pink Floyd
		},
		{
			name: "Filter by Members (4 members)",
			filter: ArtistFilter{
				Members: []int{4},
			},
			expected: 2, // Queen, Pink Floyd
		},
		{
			name: "Filter by Location (USA)",
			filter: ArtistFilter{
				Locations: []string{"USA"},
			},
			expected: 2, // Queen, Soja
		},
		{
			name: "Filter by First Album (Since 2000)",
			filter: ArtistFilter{
				FirstAlbumMin: 2000,
			},
			expected: 1, // Soja
		},
		{
			name: "Combined Filter (UK + 4 members)",
			filter: ArtistFilter{
				Locations: []string{"UK"},
				Members:   []int{4},
			},
			expected: 2, // Queen, Pink Floyd
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterArtists(artists, tt.filter)
			if len(result) != tt.expected {
				t.Errorf("Expected %d artists, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestGetFilterOptions(t *testing.T) {
	artists := []Artist{
		{Year: 1970, Album: "13-07-1973", Members: []string{"A", "B"}, Concerts: []Concert{{City: "A", Country: "B"}}},
		{Year: 1990, Album: "01-01-1995", Members: []string{"A"}, Concerts: []Concert{{City: "C", Country: "D"}}},
	}

	options := GetFilterOptions(artists)

	if options.CreationDateMin != 1970 || options.CreationDateMax != 1990 {
		t.Errorf("Incorrect Creation Date Range: %d-%d", options.CreationDateMin, options.CreationDateMax)
	}
	if options.FirstAlbumMin != 1973 || options.FirstAlbumMax != 1995 {
		t.Errorf("Incorrect First Album Range: %d-%d", options.FirstAlbumMin, options.FirstAlbumMax)
	}
	if len(options.Members) != 2 {
		t.Errorf("Incorrect number of member options: %d", len(options.Members))
	}
	if len(options.Locations) != 2 {
		t.Errorf("Incorrect number of location options: %d", len(options.Locations))
	}
}
