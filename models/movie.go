package models

import (
	"time"

	"github.com/google/uuid"
)

// Movie represents a movie entity
type Movie struct {
	ID         string    `json:"id"`
	Judul      string    `json:"judul"`
	Genre      string    `json:"genre"`
	TahunRilis int       `json:"tahun_rilis"`
	Sutradara  string    `json:"sutradara"`
	Pemeran    []string  `json:"pemeran"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateMovieRequest represents the request data for creating a new movie
type CreateMovieRequest struct {
	Judul      string   `json:"judul" validate:"required"`
	Genre      string   `json:"genre" validate:"required"`
	TahunRilis int      `json:"tahun_rilis" validate:"required,min=1888"`
	Sutradara  string   `json:"sutradara" validate:"required"`
	Pemeran    []string `json:"pemeran" validate:"required"`
}

// UpdateMovieRequest represents the request data for updating a movie
type UpdateMovieRequest struct {
	Judul      *string   `json:"judul,omitempty"`
	Genre      *string   `json:"genre,omitempty"`
	TahunRilis *int      `json:"tahun_rilis,omitempty"`
	Sutradara  *string   `json:"sutradara,omitempty"`
	Pemeran    *[]string `json:"pemeran,omitempty"`
}

// NewMovie creates a new Movie instance from CreateMovieRequest
func NewMovie(req CreateMovieRequest) *Movie {
	now := time.Now()
	return &Movie{
		ID:         uuid.New().String(),
		Judul:      req.Judul,
		Genre:      req.Genre,
		TahunRilis: req.TahunRilis,
		Sutradara:  req.Sutradara,
		Pemeran:    req.Pemeran,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
