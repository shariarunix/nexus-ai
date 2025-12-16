package ingestion

// Chunker splits text into smaller chunks
type Chunker struct {
	MaxChunkSize int
	Overlap      int
}

func NewChunker(maxChunkSize, overlap int) *Chunker {
	return &Chunker{
		MaxChunkSize: maxChunkSize,
		Overlap:      overlap,
	}
}

func (c *Chunker) Chunk(text string) []string {
	if len(text) <= c.MaxChunkSize {
		return []string{text}
	}

	var chunks []string
	runes := []rune(text)

	for i := 0; i < len(runes); {
		end := i + c.MaxChunkSize
		if end > len(runes) {
			end = len(runes)
		}

		chunks = append(chunks, string(runes[i:end]))

		if end == len(runes) {
			break
		}

		// Advance by chunkSize - overlap
		// Ensure we move forward at least 1 character to avoid infinite loop
		step := c.MaxChunkSize - c.Overlap
		if step < 1 {
			step = 1
		}
		i += step
	}

	return chunks
}
