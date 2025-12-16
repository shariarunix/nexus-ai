package generation

import (
	"context"
	"errors"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiClient(ctx context.Context, apiKey string) (*GeminiClient, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("gemini-2.0-flash-exp")
	return &GeminiClient{
		client: client,
		model:  model,
	}, nil
}

func (c *GeminiClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	if prompt == "" {
		return "", errors.New("prompt cannot be empty")
	}

	res, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("no response candidate")
	}

	// Simple extraction of first part text
	if txt, ok := res.Candidates[0].Content.Parts[0].(genai.Text); ok {
		return string(txt), nil
	}

	return "", errors.New("unexpected response format")
}
