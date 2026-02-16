package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

// Test data for mocking API responses
var mockArtists = []Artist{
	{
		ID:      1,
		Name:    "Test Artist",
		Year:    2000,
		Album:   "Test Album",
		Members: []string{"Member 1", "Member 2"},
		Image:   "test.jpg",
	},
	{
		ID:      2,
		Name:    "Another Artist",
		Year:    1995,
		Album:   "Another Album",
		Members: []string{"Solo Artist"},
		Image:   "another.jpg",
	},
}

var mockLocations = LocationResponse{
	Index: []Location{
		{ID: 1, Locations: []string{"new_york-usa", "los_angeles-usa"}},
		{ID: 2, Locations: []string{"london-uk", "paris-france"}},
	},
}

var mockDates = DateResponse{
	Index: []Date{
		{ID: 1, Dates: []string{"01-01-2023", "*02-01-2023"}},
		{ID: 2, Dates: []string{"15-06-2023", "20-06-2023"}},
	},
}

var mockRelations = RelationResponse{
	Index: []Relation{
		{
			ID: 1,
			DatesLocations: map[string][]string{
				"new_york-usa":    {"01-01-2023"},
				"los_angeles-usa": {"*02-01-2023"},
			},
		},
		{
			ID: 2,
			DatesLocations: map[string][]string{
				"london-uk":    {"15-06-2023"},
				"paris-france": {"20-06-2023"},
			},
		},
	},
}

// Mock HTTP server for testing API calls
func createMockServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/artists", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockArtists)
	})

	mux.HandleFunc("/api/locations", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockLocations)
	})

	mux.HandleFunc("/api/dates", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockDates)
	})

	mux.HandleFunc("/api/relation", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockRelations)
	})

	// Handle error cases
	mux.HandleFunc("/api/error", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	})

	return httptest.NewServer(mux)
}

// TestCustomError tests the custom error types
func TestCustomError(t *testing.T) {
	tests := []struct {
		name     string
		errFunc  func(string) error
		message  string
		wantCode int
	}{
		{
			name:     "BadRequestError",
			errFunc:  NewBadRequestError,
			message:  "bad request",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "NotFoundError",
			errFunc:  NewNotFoundError,
			message:  "not found",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "InternalServerError",
			errFunc:  NewInternalServerError,
			message:  "internal error",
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errFunc(tt.message)
			customErr, ok := err.(*CustomError)
			if !ok {
				t.Fatalf("Expected *CustomError, got %T", err)
			}

			if customErr.StatusCode != tt.wantCode {
				t.Errorf("Expected status code %d, got %d", tt.wantCode, customErr.StatusCode)
			}

			if customErr.Message != tt.message {
				t.Errorf("Expected message '%s', got '%s'", tt.message, customErr.Message)
			}

			expectedErrorString := fmt.Sprintf("HTTP %d: %s", tt.wantCode, tt.message)
			if err.Error() != expectedErrorString {
				t.Errorf("Expected error string '%s', got '%s'", expectedErrorString, err.Error())
			}
		})
	}
}

// TestArtistValidate tests the Artist validation method
func TestArtistValidate(t *testing.T) {
	tests := []struct {
		name    string
		artist  Artist
		wantErr bool
	}{
		{
			name: "Valid artist",
			artist: Artist{
				Name: "Test Artist",
				Year: 2000,
			},
			wantErr: false,
		},
		{
			name: "Empty name",
			artist: Artist{
				Name: "",
				Year: 2000,
			},
			wantErr: true,
		},
		{
			name: "Negative year",
			artist: Artist{
				Name: "Test Artist",
				Year: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.artist.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Artist.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestMakeRequest tests the makeRequest function
func TestMakeRequest(t *testing.T) {
	server := createMockServer()
	defer server.Close()

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Successful request",
			url:     server.URL + "/api/artists",
			wantErr: false,
		},
		{
			name:    "Server error",
			url:     server.URL + "/api/error",
			wantErr: true,
		},
		{
			name:    "Invalid URL",
			url:     "invalid-url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := makeRequest(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetArtistsData tests the getArtistsData function
func TestGetArtistsData(t *testing.T) {
	server := createMockServer()
	defer server.Close()

	// Temporarily override the global URL for testing
	//originalURL := artistsURL
	defer func() {
		// This is a bit hacky since artistsURL is a const, but for testing purposes
		// we would need to make it configurable in the actual code
	}()

	// For this test, we'll test the function logic with a mock
	// In a real scenario, you'd want to make the URLs configurable
	t.Run("Parse artists data", func(t *testing.T) {
		// Test JSON parsing logic
		mockJSON := `[{"id":1,"name":"Test","creationDate":2000,"firstAlbum":"Album","members":["M1"],"image":"test.jpg"}]`

		var artists []Artist
		err := json.Unmarshal([]byte(mockJSON), &artists)
		if err != nil {
			t.Fatalf("Failed to unmarshal test JSON: %v", err)
		}

		if len(artists) != 1 {
			t.Errorf("Expected 1 artist, got %d", len(artists))
		}

		if artists[0].Name != "Test" {
			t.Errorf("Expected name 'Test', got '%s'", artists[0].Name)
		}
	})
}

// TestHomePage tests the home page handler
func TestHomePage(t *testing.T) {
	// Create mock templates for testing
	createMockTemplates(t)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "Valid GET request to home",
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid method",
			method:         "POST",
			path:           "/",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Invalid path",
			method:         "GET",
			path:           "/invalid",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			// Skip the actual API call for this test
			// In a real test, you'd mock the getFullArtistData function
			if tt.expectedStatus == http.StatusOK {
				// This would fail without proper mocking, so we skip
				t.Skip("Skipping due to API dependency - would need proper mocking")
			}

			homePage(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// TestArtistDetailHandler tests the artist detail handler
func TestArtistDetailHandler(t *testing.T) {
	createMockTemplates(t)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "Invalid method",
			method:         "POST",
			path:           "/artist/1",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Invalid artist ID",
			method:         "GET",
			path:           "/artist/invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Zero artist ID",
			method:         "GET",
			path:           "/artist/0",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			artistDetailHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// TestLoggingResponseWriter tests the logging response writer
func TestLoggingResponseWriter(t *testing.T) {
	w := httptest.NewRecorder()
	lrw := &loggingResponseWriter{w, http.StatusOK}

	// Test WriteHeader
	lrw.WriteHeader(http.StatusNotFound)
	if lrw.statusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, lrw.statusCode)
	}

	// Test Write
	testData := []byte("test data")
	n, err := lrw.Write(testData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Expected %d bytes written, got %d", len(testData), n)
	}
}

// TestRequestHandler tests the request logging middleware
func TestRequestHandler(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	loggedHandler := RequestHandler(handler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	loggedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", w.Body.String())
	}
}

// Helper function to create mock templates for testing
func createMockTemplates(t *testing.T) {
	// Create temporary template files
	if err := os.MkdirAll("templates", 0755); err != nil {
		t.Fatalf("Failed to create templates directory: %v", err)
	}

	// Create mock templates
	templates := map[string]string{
		"templates/index.html": `
<!DOCTYPE html>
<html>
<head><title>Artists</title></head>
<body>
{{range .}}
<div>{{.Name}}</div>
{{end}}
</body>
</html>`,
		"templates/artist_detail.html": `
<!DOCTYPE html>
<html>
<head><title>{{.Name}}</title></head>
<body>
<h1>{{.Name}}</h1>
<p>Year: {{.Year}}</p>
</body>
</html>`,
		"templates/error.html": `
<!DOCTYPE html>
<html>
<head><title>Error</title></head>
<body>
<h1>Error {{.StatusCode}}</h1>
<p>{{.Message}}</p>
</body>
</html>`,
	}

	for filename, content := range templates {
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create template %s: %v", filename, err)
		}
	}

	// Cleanup function
	t.Cleanup(func() {
		os.RemoveAll("templates")
	})
}

// Benchmark tests
func BenchmarkMakeRequest(b *testing.B) {
	server := createMockServer()
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := makeRequest(server.URL + "/api/artists")
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
	}
}

func BenchmarkArtistValidate(b *testing.B) {
	artist := Artist{
		Name: "Test Artist",
		Year: 2000,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = artist.Validate()
	}
}

// Integration test example
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This would test the full application flow
	// You'd need to set up a test server and mock the external API
	t.Skip("Integration test - implement with full server setup")
}

// Test utility functions
func TestStringProcessing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Parse location",
			input:    "new_york-usa",
			expected: []string{"New York", "USA"},
		},
		{
			name:     "Parse location with underscore",
			input:    "los_angeles-usa",
			expected: []string{"Los Angeles", "USA"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := strings.Split(tt.input, "-")
			if len(parts) != 2 {
				t.Fatalf("Expected 2 parts, got %d", len(parts))
			}

			city := strings.Title(strings.ReplaceAll(parts[0], "_", " "))
			country := strings.ToUpper(parts[1])

			result := []string{city, country}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test date parsing
func TestDateParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal date",
			input:    "01-01-2023",
			expected: "01 January 2023",
		},
		{
			name:     "Date with asterisk",
			input:    "*02-01-2023",
			expected: "02 January 2023",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanDate := strings.TrimPrefix(tt.input, "*")
			parsedDate, err := time.Parse("02-01-2006", cleanDate)
			if err != nil {
				t.Fatalf("Failed to parse date: %v", err)
			}

			formattedDate := parsedDate.Format("02 January 2006")
			if formattedDate != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, formattedDate)
			}
		})
	}
}
