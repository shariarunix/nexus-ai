package ingestion

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock PDF creation is hard without external tools, so we tests against a small sample file if available, or try to mock the reader if we design the interface right.
// For TDD, let's define the interface and test basic errors first.

func TestParsePDF_EmptyFile(t *testing.T) {
	parser := NewPDFParser()
	data := []byte("not a pdf")
	reader := bytes.NewReader(data)

	// Should fail because it's not a valid PDF structure
	_, err := parser.Parse(reader, int64(len(data)))
	assert.Error(t, err)
}

// In a real scenario, we would add a test file resource and test parsing it.
// For now, testing logic that relies on the library behavior might be integration testing.
