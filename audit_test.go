package main

import (
	"strings"
	"testing"
)

func TestAuditCases(t *testing.T) {
	// Setup mock data for audit scenarios
	mockArtists := []Artist{
		{ID: 1, Name: "Green Day", Members: []string{"Billie Joe Armstrong", "Mike Dirnt", "Tre Cool"}, Year: 1987, Album: "13-04-1990", Concerts: []Concert{{City: "Osaka", Country: "Japan"}, {City: "Saitama", Country: "Japan"}, {City: "Nagoya", Country: "Japan"}}},
		{ID: 2, Name: "Scorpions", Members: []string{"Rudolf Schenker"}, Year: 1965, Album: "09-02-1972", Concerts: []Concert{{City: "Hanover", Country: "Germany"}}},
		{ID: 3, Name: "The Jimi Hendrix Experience", Members: []string{"Jimi Hendrix"}, Year: 1966, Album: "12-05-1967", Concerts: []Concert{{City: "London", Country: "UK"}}},
		{ID: 4, Name: "Genesis", Members: []string{"Phil Collins"}, Year: 1967, Album: "07-03-1969", Concerts: []Concert{{City: "London", Country: "UK"}}},
		{ID: 5, Name: "Pink Floyd", Members: []string{"David Gilmour"}, Year: 1965, Album: "05-08-1967", Concerts: []Concert{{City: "London", Country: "UK"}}},
		{ID: 6, Name: "Queen", Members: []string{"Freddie Mercury"}, Year: 1970, Album: "13-07-1973", Concerts: []Concert{{City: "London", Country: "UK"}, {City: "Dunedin", Country: "New Zealand"}}}, // Dunedin is in Otago, usually not associated with Queensland widely, but let's assume valid data for "queen" test as simple string match
		{ID: 7, Name: "AC/DC", Members: []string{"Angus Young"}, Year: 1973, Album: "17-02-1975", Concerts: []Concert{{City: "Sydney", Country: "Australia"}}},
		// "Queensland" implies a location match
		{ID: 8, Name: "Some Band", Concerts: []Concert{{City: "Brisbane", Country: "Queensland"}}},
	}

	cacheMutex.Lock()
	cachedArtists = mockArtists
	cacheMutex.Unlock()

	tests := []struct {
		desc       string
		query      string
		expectText []string // Substrings that MUST appear in the results (in `text` field)
	}{
		{
			desc:       "Billie Joe -> Billie Joe Armstrong (Green Day)",
			query:      "Billie Joe",
			expectText: []string{"Billie Joe Armstrong (Green Day)"},
		},
		{
			desc:       "Japan -> saitama-japan, osaka-japan, nagoya-japan",
			query:      "Japan",
			expectText: []string{"Saitama, Japan (Green Day)", "Osaka, Japan (Green Day)", "Nagoya, Japan (Green Day)"},
		},
		{
			desc:       "Scorpions -> Scorpions",
			query:      "Scorpions",
			expectText: []string{"Scorpions"},
		},
		{
			desc:       "Jimi Hendrix -> The Jimi Hendrix Experience",
			query:      "Jimi Hendrix",
			expectText: []string{"The Jimi Hendrix Experience"}, // Artist match
		},
		{
			desc:       "Phil Collins -> Phil Collins (Genesis)",
			query:      "Phil Collins",
			expectText: []string{"Phil Collins (Genesis)"},
		},
		{
			desc:       "london-uk -> Pink Floyd, Queen...",
			query:      "london", // "London, UK"
			expectText: []string{"London, UK (Pink Floyd)", "London, UK (Queen)", "London, UK (Genesis)", "London, UK (The Jimi Hendrix Experience)"},
		},
		{
			desc:       "queen -> Queen, Queensland",
			query:      "queen",
			expectText: []string{"Queen", "Queensland (Some Band)"},
		},
		{
			desc:       "05-08-1967 -> Pink Floyd",
			query:      "05-08-1967",
			expectText: []string{"05-08-1967 (Pink Floyd)"},
		},
		{
			desc:       "1973 -> AC/DC",
			query:      "1973",
			expectText: []string{"1973 (AC/DC)"},
		},
		{
			desc:       "1965 -> Scorpions, Pink Floyd",
			query:      "1965",
			expectText: []string{"1965 (Scorpions)", "1965 (Pink Floyd)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			suggestions := performSearch(tt.query)

			for _, expect := range tt.expectText {
				found := false
				for _, s := range suggestions {
					if strings.Contains(s.Text, expect) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Query '%s': Expected result containing '%s', but not found. Got: %v", tt.query, expect, suggestions)
				}
			}
		})
	}
}
