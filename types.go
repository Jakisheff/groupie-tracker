package main

// Artist represents an artist with their details and concerts
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Year         int      `json:"creationDate"`
	Album        string   `json:"firstAlbum"`
	Members      []string `json:"members"`
	Concerts     []Concert
	LocationsURL string `json:"locations"`
	DatesURL     string `json:"concertDates"`
	RelationsURL string `json:"relations"`
}

// Validate checks if the artist data is valid
func (a *Artist) Validate() error {
	if a.Name == "" {
		return NewInternalServerError("invalid artist data: empty name")
	}
	if a.Year < 0 {
		return NewInternalServerError("invalid artist data: negative year")
	}
	return nil
}

// Concert represents a concert event with location and date
type Concert struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Date    string `json:"date"`
}

// Location represents a set of locations for an artist
type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

// LocationResponse wraps the list of locations
type LocationResponse struct {
	Index []Location `json:"index"`
}

// Date represents a set of concert dates for an artist
type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// DateResponse wraps the list of dates
type DateResponse struct {
	Index []Date `json:"index"`
}

// Relation represents the relationship between dates and locations
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// RelationResponse wraps the list of relations
type RelationResponse struct {
	Index []Relation `json:"index"`
}

// MapMarker represents a marker on the map
type MapMarker struct {
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
	Title string  `json:"title"`
}
