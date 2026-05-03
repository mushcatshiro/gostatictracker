package markup

import (
	"bytes"
	"fmt"

	"github.com/mushcatshiro/gostatictracker/gvfs"
	"github.com/yuin/goldmark/parser"
)

// Converter defines the interface for markdown to HTML conversion
type Converter interface {
	Convert(content []byte, page gvfs.Page) ([]byte, error)
}

// Render is the high-level function that uses the Goldmark implementation
func Render(content []byte, page gvfs.Page) ([]byte, error) {
	c := NewGoldmarkConverter()
	return c.Convert(content, page)
}

// to clean up the args
func (c *GoldmarkConverter) Convert(content []byte, page gvfs.Page) ([]byte, error) {
	var buf bytes.Buffer
	pc := parser.NewContext()
	pc.Set(pageContextKey, page)
	if err := c.engine.Convert(content, &buf); err != nil {
		return nil, fmt.Errorf("failed to convert markdown: %w", err)
	}
	return buf.Bytes(), nil
}
