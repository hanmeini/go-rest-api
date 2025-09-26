package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-flix-api/models"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// In-memory database
var moviesDB = make(map[string]models.Movie)

// Global validator instance
var validate *validator.Validate

// Initialize validator
func init() {
	validate = validator.New()

	// Add dummy data for testing/demo
	seedDummyData()
}

// seedDummyData adds sample movies to the database
func seedDummyData() {
	dummyMovies := []models.Movie{
		{
			ID:         "550e8400-e29b-41d4-a716-446655440001",
			Judul:      "The Dark Knight",
			Genre:      "Action",
			TahunRilis: 2008,
			Sutradara:  "Christopher Nolan",
			Pemeran:    []string{"Christian Bale", "Heath Ledger", "Aaron Eckhart"},
			CreatedAt:  time.Now().Add(-24 * time.Hour),
			UpdatedAt:  time.Now().Add(-24 * time.Hour),
		},
		{
			ID:         "550e8400-e29b-41d4-a716-446655440002",
			Judul:      "Inception",
			Genre:      "Sci-Fi",
			TahunRilis: 2010,
			Sutradara:  "Christopher Nolan",
			Pemeran:    []string{"Leonardo DiCaprio", "Marion Cotillard", "Tom Hardy"},
			CreatedAt:  time.Now().Add(-48 * time.Hour),
			UpdatedAt:  time.Now().Add(-48 * time.Hour),
		},
		{
			ID:         "550e8400-e29b-41d4-a716-446655440003",
			Judul:      "Parasite",
			Genre:      "Thriller",
			TahunRilis: 2019,
			Sutradara:  "Bong Joon-ho",
			Pemeran:    []string{"Song Kang-ho", "Lee Sun-kyun", "Cho Yeo-jeong"},
			CreatedAt:  time.Now().Add(-72 * time.Hour),
			UpdatedAt:  time.Now().Add(-72 * time.Hour),
		},
	}

	// Add dummy movies to database
	for _, movie := range dummyMovies {
		moviesDB[movie.ID] = movie
	}
}

// CreateMovieHandler handles POST /api/movies
func CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateMovieRequest

	// Read request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	// Create new movie
	movie := models.NewMovie(req)

	// Save to database
	moviesDB[movie.ID] = *movie

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Return created movie
	json.NewEncoder(w).Encode(movie)
}

// GetMoviesHandler handles GET /api/movies
func GetMoviesHandler(w http.ResponseWriter, r *http.Request) {
	// Convert map to slice
	var movies []models.Movie
	for _, movie := range moviesDB {
		movies = append(movies, movie)
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Return all movies
	json.NewEncoder(w).Encode(movies)
}

// GetMovieByIDHandler handles GET /api/movies/{id}
func GetMovieByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if movie exists
	movie, exists := moviesDB[id]
	if !exists {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Return movie
	json.NewEncoder(w).Encode(movie)
}

// UpdateMovieHandler handles PUT /api/movies/{id}
func UpdateMovieHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if movie exists
	movie, exists := moviesDB[id]
	if !exists {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	var req models.UpdateMovieRequest

	// Read request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Update fields only if they are provided in request
	if req.Judul != nil {
		movie.Judul = *req.Judul
	}
	if req.Genre != nil {
		movie.Genre = *req.Genre
	}
	if req.TahunRilis != nil {
		movie.TahunRilis = *req.TahunRilis
	}
	if req.Sutradara != nil {
		movie.Sutradara = *req.Sutradara
	}
	if req.Pemeran != nil {
		movie.Pemeran = *req.Pemeran
	}

	// Update timestamp
	movie.UpdatedAt = time.Now()

	// Save updated movie
	moviesDB[id] = movie

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Return updated movie
	json.NewEncoder(w).Encode(movie)
}

// DeleteMovieHandler handles DELETE /api/movies/{id}
func DeleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if movie exists
	_, exists := moviesDB[id]
	if !exists {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	// Delete movie from database
	delete(moviesDB, id)

	// Return success response
	w.WriteHeader(http.StatusNoContent)
}

// HealthCheckHandler handles GET /health
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "movie-api",
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Return health status
	json.NewEncoder(w).Encode(response)
}
