package generation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGenerationClient struct {
	mock.Mock
}

func (m *MockGenerationClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	args := m.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}

func TestGenerateContent(t *testing.T) {
	mockClient := new(MockGenerationClient)
	prompt := "Explain Gravity"
	expectedResponse := "Gravity is a fundamental interaction..."

	mockClient.On("GenerateContent", mock.Anything, prompt).Return(expectedResponse, nil)

	res, err := mockClient.GenerateContent(context.Background(), prompt)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, res)
	mockClient.AssertExpectations(t)
}
