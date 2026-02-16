package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// API endpoints and constants
const (
	apiBaseURL     = "https://groupietrackers.herokuapp.com/api"
	artistsURL     = apiBaseURL + "/artists"
	locationsURL   = apiBaseURL + "/locations"
	datesURL       = apiBaseURL + "/dates"
	relationsURL   = apiBaseURL + "/relation"
	requestTimeout = 10 * time.Second // Timeout for HTTP requests
)

// makeRequest performs an HTTP GET request with a timeout
func makeRequest(url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to create request: %v", err))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to execute request: %v", err))
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// Continue
	case http.StatusBadRequest:
		return nil, NewBadRequestError(fmt.Sprintf("bad request to API: %s", url))
	case http.StatusNotFound:
		return nil, NewInternalServerError(fmt.Sprintf("API resource not found, possible misconfiguration: %s", url))
	case http.StatusInternalServerError:
		return nil, NewInternalServerError(fmt.Sprintf("API server error: %s", url))
	default:
		return nil, NewInternalServerError(fmt.Sprintf("unexpected response status: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to read response: %v", err))
	}

	return body, nil
}

// getArtistsData fetches artist data from the API
func getArtistsData() ([]Artist, error) {
	body, err := makeRequest(artistsURL)
	if err != nil {
		return nil, err
	}

	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to decode JSON (artists): %v", err))
	}

	// Validate artists
	for _, artist := range artists {
		if artist.Name == "" {
			return nil, NewInternalServerError("invalid artist data: empty name")
		}
		if artist.Year < 0 {
			return nil, NewInternalServerError("invalid artist data: negative year")
		}
	}

	return artists, nil
}

// getLocationsData fetches location data from the API
func getLocationsData() ([]Location, error) {
	body, err := makeRequest(locationsURL)
	if err != nil {
		return nil, err
	}

	var locationResponse LocationResponse
	err = json.Unmarshal(body, &locationResponse)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to decode JSON (locations): %v", err))
	}

	return locationResponse.Index, nil
}

// getDatesData fetches concert date data from the API
func getDatesData() ([]Date, error) {
	body, err := makeRequest(datesURL)
	if err != nil {
		return nil, err
	}

	var dateResponse DateResponse
	err = json.Unmarshal(body, &dateResponse)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to decode JSON (dates): %v", err))
	}

	return dateResponse.Index, nil
}

// getRelationsData fetches relation data from the API
func getRelationsData() ([]Relation, error) {
	body, err := makeRequest(relationsURL)
	if err != nil {
		return nil, err
	}

	var relationResponse RelationResponse
	err = json.Unmarshal(body, &relationResponse)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to decode JSON (relations): %v", err))
	}

	return relationResponse.Index, nil
}

// getFullArtistData combines artist data with locations, dates, and relations
func getFullArtistData() ([]Artist, error) {
	var wg sync.WaitGroup
	var artistsErr, locationsErr, datesErr, relationsErr error
	var artists []Artist
	var locations []Location
	var dates []Date
	var relations []Relation

	wg.Add(4)
	go func() {
		defer wg.Done()
		artists, artistsErr = getArtistsData()
	}()
	go func() {
		defer wg.Done()
		locations, locationsErr = getLocationsData()
	}()
	go func() {
		defer wg.Done()
		dates, datesErr = getDatesData()
	}()
	go func() {
		defer wg.Done()
		relations, relationsErr = getRelationsData()
	}()

	wg.Wait()

	if artistsErr != nil {
		return nil, artistsErr
	}
	if locationsErr != nil {
		return nil, locationsErr
	}
	if datesErr != nil {
		return nil, datesErr
	}
	if relationsErr != nil {
		return nil, relationsErr
	}

	locationsMap := make(map[int][]string)
	for _, loc := range locations {
		locationsMap[loc.ID] = loc.Locations
	}

	datesMap := make(map[int][]string)
	for _, date := range dates {
		datesMap[date.ID] = date.Dates
	}

	relationsMap := make(map[int]map[string][]string)
	for _, rel := range relations {
		if rel.DatesLocations == nil {
			continue // Skip invalid relations
		}
		relationsMap[rel.ID] = rel.DatesLocations
	}

	for i := range artists {
		artists[i].Concerts = []Concert{}
		if dateLocations, exists := relationsMap[artists[i].ID]; exists {
			for loc, dates := range dateLocations {
				for _, date := range dates {
					// Clean date by removing asterisk
					cleanDate := strings.TrimPrefix(date, "*")
					// Format date for readability
					parsedDate, err := time.Parse("02-01-2006", cleanDate)
					formattedDate := cleanDate
					if err == nil {
						formattedDate = parsedDate.Format("02 January 2006")
					}
					// Parse location into city and country
					parts := strings.Split(loc, "-")
					city := ""
					country := ""
					if len(parts) > 0 {
						city = strings.Title(strings.TrimSpace(parts[0]))
					}
					if len(parts) > 1 {
						country = strings.Title(strings.TrimSpace(parts[1]))
						if country == "Usa" {
							country = "USA"
						}
						if country == "New_zealand" {
							country = "New Zealand"
						}
					}
					concert := Concert{
						Date:    formattedDate,
						City:    city,
						Country: country,
					}
					artists[i].Concerts = append(artists[i].Concerts, concert)
				}
			}
		}
	}

	return artists, nil
}
