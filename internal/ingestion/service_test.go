package ingestion

import (
	"context"
	"io"
	"testing"

	"backend/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockParser struct {
	mock.Mock
}

func (m *MockParser) Parse(r io.ReaderAt, size int64) (string, error) {
	args := m.Called(r, size)
	return args.String(0), args.Error(1)
}

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) SaveChunk(ctx context.Context, chunk *domain.DocumentChunk) error {
	args := m.Called(ctx, chunk)
	return args.Error(0)
}

func (m *MockRepo) SearchSimilar(ctx context.Context, embedding []float32, limit int, filter map[string]interface{}) ([]*domain.DocumentChunk, error) {
	args := m.Called(ctx, embedding, limit, filter)
	return args.Get(0).([]*domain.DocumentChunk), args.Error(1)
}

type MockEmbedder struct {
	mock.Mock
}

func (m *MockEmbedder) EmbedContent(ctx context.Context, text string) ([]float32, error) {
	args := m.Called(ctx, text)
	return args.Get(0).([]float32), args.Error(1)
}

func TestIngestDocument(t *testing.T) {
	mockParser := new(MockParser)
	mockRepo := new(MockRepo)
	mockEmbedder := new(MockEmbedder)
	chunker := NewChunker(100, 10)

	service := NewIngestionService(mockParser, chunker, mockEmbedder, mockRepo)

	ctx := context.Background()
	content := "Physics Content"
	mockParser.On("Parse", mock.Anything, int64(10)).Return(content, nil)
	mockEmbedder.On("EmbedContent", ctx, content).Return([]float32{0.1, 0.2}, nil)
	// Note: Chunker splits "Physics Content" (15 chars) into 1 chunk if max is 100.

	mockRepo.On("SaveChunk", ctx, mock.MatchedBy(func(c *domain.DocumentChunk) bool {
		return c.Content == "Physics Content" && len(c.Embedding) == 2
	})).Return(nil)

	err := service.Ingest(ctx, nil, 10, "Physics", 1, "en")
	assert.NoError(t, err)

	mockParser.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockEmbedder.AssertExpectations(t)
}
