package movie

import (
	"context"
	"errors"
	"time"

	"go-flix-api/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ensureNotEmpty ensures that a string slice is never empty (to avoid null values in database)
func ensureNotEmpty(slice []string) []string {
	if len(slice) == 0 {
		return []string{}
	}
	return slice
}

// CreateMovie creates a new movie and saves it to the database
func (s *Service) CreateMovie(ctx context.Context, req models.CreateMovieRequest) (*models.Movie, error) {
	id := uuid.New()
	now := time.Now()
	movie := models.Movie{
		ID:         id,
		Judul:      req.Judul,
		Genre:      req.Genre,
		TahunRilis: req.TahunRilis,
		Sutradara:  req.Sutradara,
		Pemeran:    pq.StringArray(ensureNotEmpty(req.Pemeran)),
		CreatedAt:  now,
		UpdatedAt:  now,
		DeletedAt:  nil,
		CreatedBy:  req.CreatedBy,
		UpdatedBy:  req.CreatedBy,
		Version:    1,
	}
	if err := s.repo.Save(ctx, movie); err != nil {
		return nil, err
	}
	return &movie, nil
}

// GetAllMovies returns all movies
func (s *Service) GetAllMovies(ctx context.Context) ([]models.Movie, error) {
	return s.repo.FindAll(ctx)
}

// GetMovieByID returns a movie by its ID
func (s *Service) GetMovieByID(ctx context.Context, id string) (*models.Movie, error) {
	return s.repo.FindByID(ctx, id)
}

// UpdateMovie updates an existing movie
func (s *Service) UpdateMovie(ctx context.Context, id string, req models.UpdateMovieRequest, username string) error {
	movie, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
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
	movie.UpdatedAt = time.Now()
	movie.UpdatedBy = &username
	movie.Version++
	return s.repo.Update(ctx, *movie)
}

// DeleteMovie performs a soft delete
func (s *Service) DeleteMovie(ctx context.Context, id string) error {
	movie, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if movie.DeletedAt != nil {
		return errors.New("movie already deleted")
	}
	deletedAt := time.Now()
	return s.repo.Delete(ctx, id, deletedAt)
}
