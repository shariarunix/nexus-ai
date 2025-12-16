package domain

import (
	"context"
)

// VectorRepository defines the interface for vector operations
type VectorRepository interface {
	SaveChunk(ctx context.Context, chunk *DocumentChunk) error
	SearchSimilar(ctx context.Context, embedding []float32, limit int, filter map[string]interface{}) ([]*DocumentChunk, error)
}
