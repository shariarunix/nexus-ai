package handlers

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockIngestionService
type MockIngestionService struct {
	mock.Mock
}

func (m *MockIngestionService) Ingest(ctx context.Context, reader io.ReaderAt, size int64, subject string, chapter int, language string) error {
	args := m.Called(ctx, reader, size, subject, chapter, language)
	return args.Error(0)
}

// IngestionService interface is defined in handler package if we want to decoupling
// But for test we mocked it.
// We need to make sure NewDocumentHandler accepts the mock.
// The test defined MockIngestionService.

func TestUploadDocument(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockIngestionService)
	handler := NewDocumentHandler(mockService) // We need to define NewDocumentHandler and interface

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create multipart request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.pdf")
	part.Write([]byte("fake pdf content"))

	writer.WriteField("chapter", "1")
	writer.WriteField("subject", "Physics")
	writer.WriteField("language", "en")
	writer.Close()

	req, _ := http.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req

	// Mock Expectation
	mockService.On("Ingest", mock.Anything, mock.Anything, mock.Anything, "Physics", 1, "en").Return(nil)

	handler.Upload(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
