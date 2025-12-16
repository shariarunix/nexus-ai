package embedding

import (
	"context"
	"errors"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiClient handles interaction with Google Gemini API
type GeminiClient struct {
	client *genai.Client
	model  *genai.EmbeddingModel
}

// NewGeminiClient creates a new client instance
func NewGeminiClient(ctx context.Context, apiKey string) (*GeminiClient, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := client.EmbeddingModel("text-embedding-004")
	return &GeminiClient{
		client: client,
		model:  model,
	}, nil
}

// EmbedContent generates embeddings for the given text
func (c *GeminiClient) EmbedContent(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, errors.New("text cannot be empty")
	}

	res, err := c.model.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		return nil, err
	}

	if res.Embedding == nil || len(res.Embedding.Values) == 0 {
		return nil, errors.New("no embedding returned")
	}

	return res.Embedding.Values, nil
}
