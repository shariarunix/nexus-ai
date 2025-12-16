package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GeneratorService interface {
	GenerateQuestions(ctx context.Context, topic string, chapter, count int, language string) ([]string, error)
}

type QuestionHandler struct {
	service GeneratorService
}

func NewQuestionHandler(service GeneratorService) *QuestionHandler {
	return &QuestionHandler{service: service}
}

type GenerateRequest struct {
	Topic    string `json:"topic" binding:"required"`
	Chapter  int    `json:"chapter" binding:"required,gt=0"`
	Count    int    `json:"count" binding:"required,gt=0"`
	Language string `json:"language" binding:"required,oneof=en bn"`
}

// Generate godoc
// @Summary      Generate Questions
// @Description  Generates exam-style questions based on a topic and chapter context.
// @Tags         questions
// @Accept       json
// @Produce      json
// @Param        request  body      GenerateRequest  true  "Generation Request"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /questions/generate [post]
func (h *QuestionHandler) Generate(c *gin.Context) {
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	questions, err := h.service.GenerateQuestions(c.Request.Context(), req.Topic, req.Chapter, req.Count, req.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"questions": questions,
		},
	})
}
