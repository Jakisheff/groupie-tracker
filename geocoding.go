package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Coordinates represents latitude and longitude
type Coordinates struct {
	Lat float64 `json:"lat,string"`
	Lng float64 `json:"lon,string"`
}

// GeocodingCache stores cached coordinates
var (
	geocodingCache = make(map[string]Coordinates)
	geoCacheMutex  sync.RWMutex
)

// GeocodeLocation fetches coordinates for a given location (city, country)
func GeocodeLocation(city, country string) (Coordinates, error) {
	locationKey := fmt.Sprintf("%s, %s", city, country)

	// Check cache first
	geoCacheMutex.RLock()
	if coords, found := geocodingCache[locationKey]; found {
		geoCacheMutex.RUnlock()
		return coords, nil
	}
	geoCacheMutex.RUnlock()

	// Rate limiting: 1 request per second (simple implementation)
	time.Sleep(1100 * time.Millisecond)

	// Construct URL
	baseURL := "https://nominatim.openstreetmap.org/search"
	u, err := url.Parse(baseURL)
	if err != nil {
		return Coordinates{}, err
	}
	q := u.Query()
	q.Set("q", locationKey)
	q.Set("format", "json")
	q.Set("limit", "1")
	u.RawQuery = q.Encode()

	// Create request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return Coordinates{}, err
	}
	req.Header.Set("User-Agent", "GroupieTracker/1.0 (dshadykh@student.21-school.ru)") // Replace with valid email or identifier

	// Execute request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Coordinates{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Coordinates{}, fmt.Errorf("geocoding API returned status: %d", resp.StatusCode)
	}

	// Parse response
	var results []Coordinates
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return Coordinates{}, err
	}

	if len(results) == 0 {
		return Coordinates{}, fmt.Errorf("no coordinates found for: %s", locationKey)
	}

	// Cache result
	geoCacheMutex.Lock()
	geocodingCache[locationKey] = results[0]
	geoCacheMutex.Unlock()

	return results[0], nil
}
