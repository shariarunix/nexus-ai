package ingestion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunkText(t *testing.T) {
	chunker := NewChunker(10, 2) // MaxChunkSize 10, Overlap 2
	text := "Hello world this is a test"
	// "Hello worl" (10)
	// "rld this i" (overlap 2 "rl", then take 10?)
	// Let's implement recursive or simple sliding window.
	// For RAG, usually recursive character text splitter.
	// Let's test basic splitting.

	// If implementation is simple:
	// "Hello worl"
	// "ld this is"
	// "is a test"

	chunks := chunker.Chunk(text)
	assert.NotEmpty(t, chunks)
	for _, c := range chunks {
		assert.LessOrEqual(t, len(c), 10)
	}
}
