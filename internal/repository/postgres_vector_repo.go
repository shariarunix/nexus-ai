package repository

import (
	"context"
	"fmt"

	"backend/internal/domain"

	"github.com/jmoiron/sqlx"
	"github.com/pgvector/pgvector-go"
)

type PostgresVectorRepo struct {
	db *sqlx.DB
}

func NewPostgresVectorRepo(db *sqlx.DB) *PostgresVectorRepo {
	return &PostgresVectorRepo{db: db}
}

func (r *PostgresVectorRepo) SaveChunk(ctx context.Context, chunk *domain.DocumentChunk) error {
	query := `INSERT INTO embeddings (id, subject, chapter, content, embedding, language, page, created_at) 
			  VALUES (:id, :subject, :chapter, :content, :embedding, :language, :page, :created_at)`

	// map domain struct to db struct if needed, or use struct tags.
	// We need to handle the []float32 -> pgvector.Vector conversion explicitly if sqlx doesn't handle it automatically with the driver.
	// However, pgvector-go provides a type that implements Scanner/Valuer.
	// Let's create a temporary struct or modify the domain model.
	// To keep domain clean, we wrap it here.

	dbChunk := struct {
		*domain.DocumentChunk
		Embedding pgvector.Vector `db:"embedding"`
	}{
		DocumentChunk: chunk,
		Embedding:     pgvector.NewVector(chunk.Embedding),
	}

	_, err := r.db.NamedExecContext(ctx, query, dbChunk)
	return err
}

func (r *PostgresVectorRepo) SearchSimilar(ctx context.Context, embedding []float32, limit int, filter map[string]interface{}) ([]*domain.DocumentChunk, error) {
	// Basic similarity search using cosine distance (<=> operator)
	// Note: sqlx named args for SELECT is tricky with order by operators sometimes, but we can use $1

	query := `SELECT id, subject, chapter, content, language, page, created_at 
			  FROM embeddings 
			  ORDER BY embedding <=> $1 
			  LIMIT $2`

	vector := pgvector.NewVector(embedding)

	var chunks []*domain.DocumentChunk
	err := r.db.SelectContext(ctx, &chunks, query, vector, limit)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	return chunks, nil
}
