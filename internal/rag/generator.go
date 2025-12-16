package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"backend/internal/domain"
)

type GenerationClient interface {
	GenerateContent(ctx context.Context, prompt string) (string, error)
}

type RetrieverInterface interface {
	Retrieve(ctx context.Context, query string, limit int, filter map[string]interface{}) ([]*domain.DocumentChunk, error)
}

type GeneratorService struct {
	client    GenerationClient
	retriever RetrieverInterface
}

func NewGeneratorService(client GenerationClient, retriever RetrieverInterface) *GeneratorService {
	return &GeneratorService{
		client:    client,
		retriever: retriever,
	}
}

func (s *GeneratorService) GenerateQuestions(ctx context.Context, topic string, chapter, count int, language string) ([]string, error) {
	// 1. Retrieve Context
	filter := map[string]interface{}{
		"chapter":  chapter,
		"language": language,
	}
	// Retrieve ample context chunks, e.g., 20
	chunks, err := s.retriever.Retrieve(ctx, topic, 20, filter)
	if err != nil {
		return nil, fmt.Errorf("retrieval failed: %w", err)
	}

	if len(chunks) == 0 {
		return nil, fmt.Errorf("no context found for topic %s in chapter %d", topic, chapter)
	}

	// 2. Build Context String
	var sb strings.Builder
	for _, c := range chunks {
		sb.WriteString(c.Content)
		sb.WriteString("\n---\n")
	}

	// 3. Construct Prompt
	prompt := fmt.Sprintf(`
You are a physics examiner. Generate %d exam-style questions on the topic "%s".
Use ONLY the following context to generate the questions.
Language: %s.

Context:
%s

Output STRICT JSON and nothing else:
{
  "questions": ["Question 1", "Question 2"]
}
`, count, topic, language, sb.String())

	// 4. Generate
	resp, err := s.client.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("generation failed: %w", err)
	}

	// 5. Parse JSON
	cleaned := cleanJSON(resp)

	var result struct {
		Questions []string `json:"questions"`
	}
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("parsing response failed: %w. Response: %s", err, resp)
	}

	if len(result.Questions) > count {
		result.Questions = result.Questions[:count]
	}

	return result.Questions, nil
}

func cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}
