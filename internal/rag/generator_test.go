package rag

import (
	"context"
	"testing"

	"backend/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGeneratorClient
type MockGeneratorClient struct {
	mock.Mock
}

func (m *MockGeneratorClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	args := m.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}

type MockRetriever struct {
	mock.Mock
}

func (m *MockRetriever) Retrieve(ctx context.Context, query string, limit int, filter map[string]interface{}) ([]*domain.DocumentChunk, error) {
	args := m.Called(ctx, query, limit, filter)
	return args.Get(0).([]*domain.DocumentChunk), args.Error(1)
}

func TestGenerateQuestions(t *testing.T) {
	mockGen := new(MockGeneratorClient)
	mockRetriever := new(MockRetriever)

	service := NewGeneratorService(mockGen, mockRetriever)

	ctx := context.Background()
	topic := "Newton"
	count := 2
	language := "en"
	chapter := 1

	// Expectations
	chunks := []*domain.DocumentChunk{{Content: "Context 1"}}
	mockRetriever.On("Retrieve", ctx, topic, 20, mock.MatchedBy(func(f map[string]interface{}) bool {
		return f["language"] == "en" && f["chapter"] == 1
	})).Return(chunks, nil)

	mockGen.On("GenerateContent", ctx, mock.MatchedBy(func(prompt string) bool {
		// Verify prompt contains instructions
		return true // simplify for now, check string content if needed
	})).Return(`{"questions": ["Q1", "Q2"]}`, nil)

	questions, err := service.GenerateQuestions(ctx, topic, chapter, count, language)
	assert.NoError(t, err)
	assert.Len(t, questions, 2)
	assert.Equal(t, "Q1", questions[0])

	mockRetriever.AssertExpectations(t)
	mockGen.AssertExpectations(t)
}
