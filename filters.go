package main

import (
	"strings"
	"time"
)

// ArtistFilter holds the filter criteria from the request
type ArtistFilter struct {
	CreationDateMin int
	CreationDateMax int
	FirstAlbumMin   int
	FirstAlbumMax   int
	Members         []int
	Locations       []string
}

// FilterOptions holds the available options and ranges for the frontend
type FilterOptions struct {
	CreationDateMin int
	CreationDateMax int
	FirstAlbumMin   int
	FirstAlbumMax   int
	Members         []int
	Locations       []string
}

// PageData holds the data for the index page template
type PageData struct {
	Artists       []Artist
	FilterOptions FilterOptions
	CurrentFilter ArtistFilter
}

// FilterArtists applies the filters to the list of artists
func FilterArtists(artists []Artist, filter ArtistFilter) []Artist {
	var filtered []Artist

	for _, artist := range artists {
		// Filter by Creation Date
		if filter.CreationDateMin > 0 && artist.Year < filter.CreationDateMin {
			continue
		}
		if filter.CreationDateMax > 0 && artist.Year > filter.CreationDateMax {
			continue
		}

		// Filter by First Album Date
		firstAlbumYear := parseYear(artist.Album)
		if filter.FirstAlbumMin > 0 && firstAlbumYear < filter.FirstAlbumMin {
			continue
		}
		if filter.FirstAlbumMax > 0 && firstAlbumYear > filter.FirstAlbumMax {
			continue
		}

		// Filter by Number of Members
		if len(filter.Members) > 0 {
			found := false
			for _, m := range filter.Members {
				if len(artist.Members) == m {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filter by Locations
		if len(filter.Locations) > 0 {
			found := false
			for _, locFilter := range filter.Locations {
				locFilter = strings.ToLower(strings.TrimSpace(locFilter))
				for _, concert := range artist.Concerts {
					// Check city and country
					city := strings.ToLower(concert.City)
					country := strings.ToLower(concert.Country)
					// Full location string "city-country" or "city, country"
					fullLoc := city + ", " + country
					fullLocHyphen := city + "-" + country

					// User hint: "Seattle, Washington, USA is part of Washington, USA"
					// We'll check if the filter term is contained in the location string
					if strings.Contains(fullLoc, locFilter) || strings.Contains(fullLocHyphen, locFilter) ||
						strings.Contains(city, locFilter) || strings.Contains(country, locFilter) {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				continue
			}
		}

		filtered = append(filtered, artist)
	}

	return filtered
}

// GetFilterOptions generates the available filter options from the artist data
func GetFilterOptions(artists []Artist) FilterOptions {
	options := FilterOptions{
		CreationDateMin: 3000, // Initialize with high value
		CreationDateMax: 0,
		FirstAlbumMin:   3000,
		FirstAlbumMax:   0,
		Members:         []int{},
		Locations:       []string{},
	}

	membersMap := make(map[int]bool)
	locationsMap := make(map[string]bool)

	for _, artist := range artists {
		// Creation Date
		if artist.Year < options.CreationDateMin {
			options.CreationDateMin = artist.Year
		}
		if artist.Year > options.CreationDateMax {
			options.CreationDateMax = artist.Year
		}

		// First Album Date
		year := parseYear(artist.Album)
		if year < options.FirstAlbumMin {
			options.FirstAlbumMin = year
		}
		if year > options.FirstAlbumMax {
			options.FirstAlbumMax = year
		}

		// Members
		membersMap[len(artist.Members)] = true

		// Locations
		for _, concert := range artist.Concerts {
			loc := concert.City + ", " + concert.Country
			locationsMap[loc] = true
		}
	}

	// Convert maps to slices
	for m := range membersMap {
		options.Members = append(options.Members, m)
	}
	for l := range locationsMap {
		options.Locations = append(options.Locations, l)
	}

	return options
}

// parseYear extracts the year from a date string (dd-mm-yyyy)
func parseYear(dateStr string) int {
	t, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		return 0
	}
	return t.Year()
}
