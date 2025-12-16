package ingestion

import (
	"fmt"
	"io"

	"github.com/ledongthuc/pdf"
)

type PDFParser struct{}

func NewPDFParser() *PDFParser {
	return &PDFParser{}
}

func (p *PDFParser) Parse(r io.ReaderAt, size int64) (string, error) {
	// ledongthuc/pdf NewReader expects io.ReaderAt and size
	reader, err := pdf.NewReader(r, size)
	if err != nil {
		return "", fmt.Errorf("failed to create pdf reader: %w", err)
	}

	var content string
	for pageIndex := 1; pageIndex <= reader.NumPage(); pageIndex++ {
		p := reader.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		text, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}
		content += text + "\n"
	}

	return content, nil
}
