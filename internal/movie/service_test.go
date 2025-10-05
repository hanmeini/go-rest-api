package movie

import (
	"context"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"go-flix-api/models"
)

// removed fake in favor of sqlmock-backed repository tests

func TestCreateMovie(t *testing.T) {
	// Use real repository methods backed by mocked DB to avoid compile issues
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRepository(sqlxDB)
	svc := NewService(repo)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO movies")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	req := models.CreateMovieRequest{
		Judul: "Title", Genre: "Action", TahunRilis: 2020, Sutradara: "Dir",
		Pemeran: []string{"A", "B"},
	}
	m, err := svc.CreateMovie(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateMovie error: %v", err)
	}
	if m.Judul != "Title" {
		t.Fatalf("unexpected title: %s", m.Judul)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUpdateMovie(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRepository(sqlxDB)
	svc := NewService(repo)

	// FindByID
	rows := sqlmock.NewRows([]string{"id", "judul", "genre", "tahun_rilis", "sutradara", "pemeran", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "version"}).
		AddRow("11111111-1111-1111-1111-111111111111", "Old", "G", 2000, "S", "{A,B}", time.Now(), time.Now(), nil, nil, nil, 1)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM movies WHERE id = $1 AND deleted_at IS NULL")).
		WithArgs("11111111-1111-1111-1111-111111111111").
		WillReturnRows(rows)

	// Update (Update uses direct DB without transaction in repository)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE movies SET")).
		WillReturnResult(sqlmock.NewResult(0, 1))

	newTitle := "New"
	req := models.UpdateMovieRequest{Judul: &newTitle}
	if err := svc.UpdateMovie(context.Background(), "11111111-1111-1111-1111-111111111111", req, "tester"); err != nil {
		t.Fatalf("UpdateMovie error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
