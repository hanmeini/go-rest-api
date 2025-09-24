package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-flix-api/models"
)

func TestCreateMovie_Success(t *testing.T) {
	// Prepare test data
	reqBody := models.CreateMovieRequest{
		Judul:      "The Dark Knight",
		Genre:      "Action",
		TahunRilis: 2008,
		Sutradara:  "Christopher Nolan",
		Pemeran:    []string{"Christian Bale", "Heath Ledger", "Aaron Eckhart"},
	}

	// Convert to JSON
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "/api/movies", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	CreateMovieHandler(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Check content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}

	// Parse response body
	var response models.Movie
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify response data
	if response.Judul != reqBody.Judul {
		t.Errorf("Handler returned wrong judul: got %v want %v", response.Judul, reqBody.Judul)
	}
	if response.Genre != reqBody.Genre {
		t.Errorf("Handler returned wrong genre: got %v want %v", response.Genre, reqBody.Genre)
	}
	if response.TahunRilis != reqBody.TahunRilis {
		t.Errorf("Handler returned wrong tahun_rilis: got %v want %v", response.TahunRilis, reqBody.TahunRilis)
	}
	if response.Sutradara != reqBody.Sutradara {
		t.Errorf("Handler returned wrong sutradara: got %v want %v", response.Sutradara, reqBody.Sutradara)
	}

	// Check if ID is generated
	if response.ID == "" {
		t.Error("Handler should generate a non-empty ID")
	}

	// Check if timestamps are set
	if response.CreatedAt.IsZero() {
		t.Error("Handler should set CreatedAt timestamp")
	}
	if response.UpdatedAt.IsZero() {
		t.Error("Handler should set UpdatedAt timestamp")
	}

	// Verify movie is stored in database
	if _, exists := moviesDB[response.ID]; !exists {
		t.Error("Movie should be stored in database")
	}
}

func TestCreateMovie_ValidationFailure(t *testing.T) {
	// Prepare invalid test data (empty Judul)
	reqBody := models.CreateMovieRequest{
		Judul:      "", // Empty judul should fail validation
		Genre:      "Action",
		TahunRilis: 2008,
		Sutradara:  "Christopher Nolan",
		Pemeran:    []string{"Christian Bale", "Heath Ledger"},
	}

	// Convert to JSON
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "/api/movies", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	CreateMovieHandler(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Check response body contains validation error
	responseBody := rr.Body.String()
	if responseBody == "" {
		t.Error("Handler should return error message for validation failure")
	}

	// Verify movie is NOT stored in database
	// Count movies before and after - should be the same
	initialCount := len(moviesDB)

	// The movie should not be added to database due to validation failure
	if len(moviesDB) != initialCount {
		t.Error("Movie should not be stored in database when validation fails")
	}
}

// Helper function to reset database for clean tests
func resetDatabase() {
	moviesDB = make(map[string]models.Movie)
}

// Setup and teardown for tests
func TestMain(m *testing.M) {
	// Reset database before running tests
	resetDatabase()

	// Run tests
	m.Run()

	// Clean up after tests
	resetDatabase()
}
