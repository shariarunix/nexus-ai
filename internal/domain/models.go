package domain

import (
	"time"

	"github.com/google/uuid"
)

// DocumentChunk represents a text chunk with its embedding
type DocumentChunk struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Subject   string    `json:"subject" db:"subject"`
	Chapter   int       `json:"chapter" db:"chapter"`
	Content   string    `json:"content" db:"content"`
	Embedding []float32 `json:"embedding" db:"embedding"` // pgvector
	Language  string    `json:"language" db:"language"`   // 'bn' or 'en'
	Page      int       `json:"page" db:"page"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Question represents a generated question
type Question struct {
	Text string `json:"text"`
}
