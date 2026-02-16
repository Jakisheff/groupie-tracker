# Groupie Tracker Geolocalization

Groupie Tracker is a web application that fetches and displays information about music artists, their concert dates, locations, and relationships using the [Groupie Tracker API](https://groupietrackers.herokuapp.com/api). It provides a user-friendly interface to browse artists, search for specific data, filter results, and view detailed concert and band information.

## Features
- **Smart Search Bar**: Real-time search with autocomplete suggestions for:
    - **Artists/Bands**: e.g., "Queen", "Scorpions"
    - **Members**: e.g., "Freddie Mercury", "Phil Collins" (shows associated band)
    - **Locations**: e.g., "London", "Japan" (shows associated band)
    - **Creation Dates**: e.g., "1975"
    - **First Album Dates**: e.g., "14-12-1973"
- **Artist List**: View a list of artists with names, images, creation dates, first albums, and members.
- **Integration with Maps**: Geolocalization of concert locations.
- **Advanced Filters**: Filter artists by creation date, first album date, number of members, and locations.
- **Artist Details**: Access detailed artist information, including concert dates, cities, and countries.
- **Responsive Design**: Styled with CSS for an intuitive user interface on all devices.
- **Error Handling**: Custom error pages for invalid requests, API failures, and missing data.
- **Concurrent Data Fetching**: Uses goroutines to efficiently fetch and cache API data.

## Search Functionality
The application features a powerful search bar with intelligent type detection. As you type, suggestions appear indicating the type of match and the associated context:

- **Artist**: Searching "Queen" suggests "Queen - artist/band".
- **Member**: Searching "Billie Joe" suggests "Billie Joe Armstrong (Green Day) - member".
- **Location**: Searching "Japan" suggests "Saitama, Japan (Green Day) - location".
- **Date**: Searching "1967" suggests "05-08-1967 (Pink Floyd) - first album".

This context ensures you know exactly which artist a member, location, or date belongs to before clicking.

## Technologies
- **Go**: Backend server, API integration, and search logic.
- **HTML/CSS/JavaScript**: Frontend templates, search interactivity, and styling.
- **Groupie Tracker API**: Source for artist, location, date, and relation data.
- **Go HTTP Server**: Built-in `net/http` package for request handling.
- **Templates**: Go `html/template` for dynamic HTML rendering.

## Installation

### Prerequisites
- Go (version 1.16 or higher)
- Git
- Web browser

### Steps
1. **Clone the repository**:
   ```bash
   git clone https://01.tomorrow-school.ai/git/dshadykh/groupie-tracker-geolocalization.git
   cd groupie-tracker-geolocalization
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the application**:
   ```bash
   go run .
   ```
   The server starts at `http://localhost:8080`.

4. **Access the application**:
   Open `http://localhost:8080` in your browser.

## Usage
- **Search**: Use the search bar in the header to find artists, members, locations, or dates. Click a suggestion to go directly to the artist's page.
- **Filters**: Use the sidebar to filter the artist list by various criteria.
- **Home Page**: Visit `http://localhost:8080/` to view the artist list.
- **Artist Details**: Navigate to `http://localhost:8080/artist/<id>` (e.g., `/artist/1`) for artist-specific details.

## Project Structure
```
groupie-tracker-geolocalization/
├── main.go                  # Entry point, caches data, sets up server & routes
├── search.go                # Search logic and API handler
├── api.go                   # API data fetching and processing
├── filters.go               # Filtering logic
├── handlers.go              # HTTP handlers and template rendering
├── types.go                 # Data structures for artists, concerts, etc.
├── error.go                 # Custom error types and rendering
├── static/                  # Static assets
│   ├── script.js            # General client-side scripts
│   ├── search.js            # Search bar interaction logic
│   └── style.css            # CSS styles
├── templates/               # HTML templates
│   ├── index.html           # Home page template
│   ├── artist_detail.html   # Artist details template
│   └── error.html           # Error page template
└── README.md                # Project documentation
```

## Authors
- Aigerim Ongalbayeva #aongalba
- Amir Zhakyshev #azhakysh
- Daniyar Shadykhanov #dshadykh