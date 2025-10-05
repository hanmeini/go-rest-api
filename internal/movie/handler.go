package movie

import (
	"encoding/json"
	"go-flix-api/models"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// @Summary Get all movies
// @Description Get a list of all movies
// @Tags movies
// @Produce json
// @Success 200 {array} models.Movie
// @Router /movies [get]
func (h *Handler) GetAllMovies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	movies, err := h.service.GetAllMovies(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// @Summary Get movie by ID
// @Description Get a movie by its ID
// @Tags movies
// @Produce json
// @Param id path string true "Movie ID"
// @Success 200 {object} models.Movie
// @Failure 404 {object} map[string]string
// @Router /movies/{id} [get]
func (h *Handler) GetMovieByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	movie, err := h.service.GetMovieByID(ctx, id)
	if err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

// @Summary Create a new movie
// @Description Create a new movie
// @Tags movies
// @Accept json
// @Produce json
// @Param movie body models.CreateMovieRequest true "Movie to create"
// @Success 201 {object} models.Movie
// @Failure 400 {object} map[string]string
// @Router /movies [post]
func (h *Handler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req models.CreateMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	// Ambil username dari header (hasil middleware JWT)
	username := r.Header.Get("X-Username")
	if username != "" {
		req.CreatedBy = &username
	}
	movie, err := h.service.CreateMovie(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

// @Summary Update a movie
// @Description Update a movie by its ID
// @Tags movies
// @Accept json
// @Produce json
// @Param id path string true "Movie ID"
// @Param movie body models.UpdateMovieRequest true "Movie fields to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /movies/{id} [put]
func (h *Handler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	var req models.UpdateMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	username := r.Header.Get("X-Username")
	if err := h.service.UpdateMovie(ctx, id, req, username); err != nil {
		if err.Error() == "no rows updated" {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "movie updated"})
}

// @Summary Delete a movie
// @Description Soft delete a movie by its ID
// @Tags movies
// @Produce json
// @Param id path string true "Movie ID"
// @Success 204 {object} nil
// @Failure 404 {object} map[string]string
// @Router /movies/{id} [delete]
func (h *Handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	if err := h.service.DeleteMovie(ctx, id); err != nil {
		if err.Error() == "no rows deleted" || err.Error() == "movie already deleted" {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
