package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Movie represents a movie entity matching the PostgreSQL schema
type Movie struct {
	ID         uuid.UUID      `json:"id" db:"id"`
	Judul      string         `json:"judul" db:"judul"`
	Genre      string         `json:"genre" db:"genre"`
	TahunRilis int            `json:"tahun_rilis" db:"tahun_rilis"`
	Sutradara  string         `json:"sutradara" db:"sutradara"`
	Pemeran    pq.StringArray `json:"pemeran" db:"pemeran" swaggertype:"array,string"`

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	CreatedBy *string    `json:"created_by,omitempty" db:"created_by"`
	UpdatedBy *string    `json:"updated_by,omitempty" db:"updated_by"`
	Version   int        `json:"version" db:"version"`
}

// CreateMovieRequest represents the request data for creating a new movie
// Kolom audit diisi otomatis di backend, user hanya input data utama
// Pemeran tetap array of string agar mudah di-parse dari JSON
// ID, CreatedAt, UpdatedAt, dsb di-generate backend
// CreatedBy bisa diisi dari JWT username jika perlu
// Version diisi default 1
// DeletedAt tidak diinput user
type CreateMovieRequest struct {
	Judul      string   `json:"judul" validate:"required"`
	Genre      string   `json:"genre" validate:"required"`
	TahunRilis int      `json:"tahun_rilis" validate:"required,min=1888"`
	Sutradara  string   `json:"sutradara" validate:"required"`
	Pemeran    []string `json:"pemeran" validate:"required"`
	CreatedBy  *string  `json:"created_by,omitempty"` // opsional, bisa diisi dari JWT
}

// UpdateMovieRequest represents the request data for updating a movie
// Hanya field yang boleh diubah user
// UpdatedBy bisa diisi dari JWT username jika perlu
// Version bisa diincrement di backend
// DeletedAt tidak diinput user
type UpdateMovieRequest struct {
	Judul      *string   `json:"judul,omitempty"`
	Genre      *string   `json:"genre,omitempty"`
	TahunRilis *int      `json:"tahun_rilis,omitempty"`
	Sutradara  *string   `json:"sutradara,omitempty"`
	Pemeran    *[]string `json:"pemeran,omitempty"`
	UpdatedBy  *string   `json:"updated_by,omitempty"` // opsional, bisa diisi dari JWT
}
