package movie

import (
	"context"
	"errors"

	"go-flix-api/models"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// FindAll returns all movies from the database
func (r *Repository) FindAll(ctx context.Context) ([]models.Movie, error) {
	var movies []models.Movie
	query := `SELECT * FROM movies WHERE deleted_at IS NULL`
	err := r.db.SelectContext(ctx, &movies, query)
	return movies, err
}

// FindByID returns a movie by its ID
func (r *Repository) FindByID(ctx context.Context, id string) (*models.Movie, error) {
	var movie models.Movie
	query := `SELECT * FROM movies WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &movie, query, id)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

// Save inserts a new movie into the database using an explicit transaction
func (r *Repository) Save(ctx context.Context, movie models.Movie) error {
	query := `INSERT INTO movies (
		id, judul, genre, tahun_rilis, sutradara, pemeran,
		created_at, updated_at, created_by, updated_by, version
	) VALUES (
		:id, :judul, :genre, :tahun_rilis, :sutradara, :pemeran,
		:created_at, :updated_at, :created_by, :updated_by, :version
	)`

	// 1. Mulai sesi transaksi baru
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	
	// 2. Defer Rollback: Ini adalah jaring pengaman.
	// Jika ada error di tengah jalan, transaksi akan otomatis dibatalkan.
	defer tx.Rollback()

	// 3. Jalankan query di dalam transaksi (menggunakan tx, bukan r.db)
	_, err = tx.NamedExecContext(ctx, query, movie)
	if err != nil {
		// Jika ada error di sini, Rollback akan otomatis terpanggil
		return err
	}

	// 4. PENTING: Jika tidak ada error sama sekali, simpan permanen dengan COMMIT.
	// Ini adalah tombol "Konfirmasi Akhir".
	return tx.Commit()
}

// Update updates an existing movie in the database
func (r *Repository) Update(ctx context.Context, movie models.Movie) error {
	query := `UPDATE movies SET
		judul = :judul,
		genre = :genre,
		tahun_rilis = :tahun_rilis,
		sutradara = :sutradara,
		pemeran = :pemeran,
		updated_at = :updated_at,
		updated_by = :updated_by,
		version = :version
	WHERE id = :id AND deleted_at IS NULL`
	result, err := r.db.NamedExecContext(ctx, query, &movie)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

// Delete performs a soft delete by setting deleted_at
func (r *Repository) Delete(ctx context.Context, id string, deletedAt interface{}) error {
	query := `UPDATE movies SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, deletedAt, id)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("no rows deleted")
	}
	return nil
}
