package handlers

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IngestionService interface {
	Ingest(ctx context.Context, reader io.ReaderAt, size int64, subject string, chapter int, language string) error
}

type DocumentHandler struct {
	service IngestionService
}

func NewDocumentHandler(service IngestionService) *DocumentHandler {
	return &DocumentHandler{service: service}
}

// Upload godoc
// @Summary      Upload a PDF document
// @Description  Uploads a PDF textbook chapter for ingestion.
// @Tags         documents
// @Accept       multipart/form-data
// @Produce      json
// @Param        file      formData  file    true  "PDF File"
// @Param        chapter   formData  int     true  "Chapter Number"
// @Param        subject   formData  string  true  "Subject Name"
// @Param        language  formData  string  true  "Language (en/bn)"
// @Security     BearerAuth
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /documents/upload [post]
func (h *DocumentHandler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	chapterStr := c.PostForm("chapter")
	subject := c.PostForm("subject")
	language := c.PostForm("language")

	if subject == "" || language == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subject and language are required"})
		return
	}

	chapter, err := strconv.Atoi(chapterStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chapter must be an integer"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer file.Close()

	if err := h.service.Ingest(c.Request.Context(), file, fileHeader.Size, subject, chapter, language); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ingestion failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document uploaded successfully"})
}
