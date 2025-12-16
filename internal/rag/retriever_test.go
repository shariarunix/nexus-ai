package rag

import (
	"context"
	"testing"

	"backend/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks define in test or reused if shared.
// Ideally shared mocks in a `mocks` package.
// For now, inline or duplicate is fine for speed.

type MockEmbedder struct {
	mock.Mock
}

func (m *MockEmbedder) EmbedContent(ctx context.Context, text string) ([]float32, error) {
	args := m.Called(ctx, text)
	return args.Get(0).([]float32), args.Error(1)
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.DocumentChunk), args.Error(1)
}

func TestRetrieve(t *testing.T) {
	mockRepo := new(MockRepo)
	mockEmbedder := new(MockEmbedder)
	retriever := NewRetriever(mockEmbedder, mockRepo)

	ctx := context.Background()
	query := "Newton laws"
	embedding := []float32{0.1, 0.2}
	mockEmbedder.On("EmbedContent", ctx, query).Return(embedding, nil)

	chunks := []*domain.DocumentChunk{
		{Content: "Law 1"},
		{Content: "Law 2"},
	}
	// Filter check
	mockRepo.On("SearchSimilar", ctx, embedding, 2, mock.MatchedBy(func(f map[string]interface{}) bool {
		return f["language"] == "en"
	})).Return(chunks, nil)

	results, err := retriever.Retrieve(ctx, query, 2, map[string]interface{}{"language": "en"})
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "Law 1", results[0].Content)

	mockEmbedder.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
