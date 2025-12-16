package rag

import (
	"context"
	"fmt"

	"backend/internal/domain"
	"backend/internal/ingestion"
)

// Retriever retrieves relevant content
type Retriever struct {
	embedder ingestion.Embedder
	repo     domain.VectorRepository
}

func NewRetriever(embedder ingestion.Embedder, repo domain.VectorRepository) *Retriever {
	return &Retriever{
		embedder: embedder,
		repo:     repo,
	}
}

func (r *Retriever) Retrieve(ctx context.Context, query string, limit int, filter map[string]interface{}) ([]*domain.DocumentChunk, error) {
	embedding, err := r.embedder.EmbedContent(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("embedding query failed: %w", err)
	}

	chunks, err := r.repo.SearchSimilar(ctx, embedding, limit, filter)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	return chunks, nil
}
