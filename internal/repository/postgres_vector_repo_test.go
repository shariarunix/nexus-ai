package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"backend/internal/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
)

func TestSaveChunk(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewPostgresVectorRepo(sqlxDB)

	chunk := &domain.DocumentChunk{
		ID:        uuid.New(),
		Subject:   "Physics",
		Chapter:   1,
		Content:   "Newton's First Law",
		Embedding: []float32{0.1, 0.2, 0.3},
		Language:  "en",
		Page:      10,
		CreatedAt: time.Now(),
	}

	query := regexp.QuoteMeta(`INSERT INTO embeddings (id, subject, chapter, content, embedding, language, page, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)

	mock.ExpectExec(query).
		WithArgs(chunk.ID, chunk.Subject, chunk.Chapter, chunk.Content, pgvector.NewVector(chunk.Embedding), chunk.Language, chunk.Page, chunk.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SaveChunk(context.Background(), chunk)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchSimilar(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewPostgresVectorRepo(sqlxDB)

	embedding := []float32{0.1, 0.2, 0.3}
	limit := 5

	// Expected rows
	rows := sqlmock.NewRows([]string{"id", "subject", "chapter", "content", "language", "page", "created_at"}).
		AddRow(uuid.New(), "Physics", 1, "Content 1", "en", 10, time.Now()).
		AddRow(uuid.New(), "Physics", 1, "Content 2", "en", 11, time.Now())

	query := regexp.QuoteMeta(`SELECT id, subject, chapter, content, language, page, created_at FROM embeddings ORDER BY embedding <=> $1 LIMIT $2`)

	mock.ExpectQuery(query).
		WithArgs(pgvector.NewVector(embedding), limit).
		WillReturnRows(rows)

	results, err := repo.SearchSimilar(context.Background(), embedding, limit, nil)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}
