package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGeneratorService
type MockGeneratorService struct {
	mock.Mock
}

func (m *MockGeneratorService) GenerateQuestions(ctx context.Context, topic string, chapter, count int, language string) ([]string, error) {
	args := m.Called(ctx, topic, chapter, count, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func TestGenerateQuestions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockGeneratorService)
	handler := NewQuestionHandler(mockService) // We need to define NewQuestionHandler

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// JSON Input
	input := `{"topic": "Physics", "chapter": 1, "count": 5, "language": "en"}`
	req, _ := http.NewRequest("POST", "/questions/generate", bytes.NewBufferString(input))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// Mock Expectation
	mockService.On("GenerateQuestions", mock.Anything, "Physics", 1, 5, "en").Return([]string{"Q1", "Q2"}, nil)

	handler.Generate(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
