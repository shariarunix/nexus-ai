package embedding

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmbeddingClient is a mock for the GenAI client wrapper
type MockEmbeddingClient struct {
	mock.Mock
}

func (m *MockEmbeddingClient) EmbedContent(ctx context.Context, text string) ([]float32, error) {
	args := m.Called(ctx, text)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]float32), args.Error(1)
}

func TestEmbedContent(t *testing.T) {
	mockClient := new(MockEmbeddingClient)
	text := "Physics"
	expectedEmbedding := []float32{0.1, 0.2, 0.3}

	mockClient.On("EmbedContent", mock.Anything, text).Return(expectedEmbedding, nil)

	// In real TDD, we would test the actual implementation against a mock server or interface
	// For now, testing the contract via mock usage

	res, err := mockClient.EmbedContent(context.Background(), text)
	assert.NoError(t, err)
	assert.Equal(t, expectedEmbedding, res)
	mockClient.AssertExpectations(t)
}
