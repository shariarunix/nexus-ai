package ingestion

import (
	"context"
	"fmt"
	"io"
	"time"

	"backend/internal/domain"

	"github.com/google/uuid"
)

// Interfaces for dependencies
type Parser interface {
	Parse(r io.ReaderAt, size int64) (string, error)
}

type Embedder interface {
	EmbedContent(ctx context.Context, text string) ([]float32, error)
}

// IngestionService coordinates the document ingestion process
type IngestionService struct {
	parser   Parser
	chunker  *Chunker
	embedder Embedder
	repo     domain.VectorRepository
}

func NewIngestionService(parser Parser, chunker *Chunker, embedder Embedder, repo domain.VectorRepository) *IngestionService {
	return &IngestionService{
		parser:   parser,
		chunker:  chunker,
		embedder: embedder,
		repo:     repo,
	}
}

func (s *IngestionService) Ingest(ctx context.Context, reader io.ReaderAt, size int64, subject string, chapter int, language string) error {
	// 1. Parse PDF
	text, err := s.parser.Parse(reader, size)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	// 2. Chunk text
	chunks := s.chunker.Chunk(text)

	// 3. Process each chunk
	for i, content := range chunks {
		// Embed
		embedding, err := s.embedder.EmbedContent(ctx, content)
		if err != nil {
			return fmt.Errorf("embedding failed for chunk %d: %w", i, err)
		}

		// Save to DB
		chunk := &domain.DocumentChunk{
			ID:        uuid.New(),
			Subject:   subject,
			Chapter:   chapter,
			Content:   content,
			Embedding: embedding,
			Language:  language,
			Page:      0, // Parser simplified, real one might map pages
			CreatedAt: time.Now(),
		}

		if err := s.repo.SaveChunk(ctx, chunk); err != nil {
			return fmt.Errorf("saving chunk %d failed: %w", i, err)
		}
	}

	return nil
}
