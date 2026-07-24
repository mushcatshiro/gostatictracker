package markup

import (
	"strings"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "Basic Markdown",
			input:    "# Hello",
			contains: "<h1>Hello</h1>",
		},
		{
			name:     "Code Block with Chroma",
			input:    "```go\npackage main\nimport \"fmt\"\n\nfunc main() {\nfmt.Println(\"Hello, World!\")\n}```",
			contains: "chroma", // Check if chroma classes are injected
		},
		{
			name:     "Mermaid Block",
			input:    "```mermaid\ngraph TD;\nA-->B;\n```",
			contains: "language-mermaid",
		},
		{
			name:     "Malicious Script Injection",
			input:    "<script>alert('XSS')</script>",
			contains: "&lt;script&gt;alert('XSS')&lt;/script&gt;", // Ensure HTML is escaped
		},
	}

	converter := NewGoldmarkConverter()
	renderCtx := RenderContext{FilePath: "test.md"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := converter.Convert([]byte(tt.input), renderCtx)
			if err != nil {
				t.Errorf("conversion failed: %v", err)
			}
			if !strings.Contains(string(output), tt.contains) {
				t.Errorf("expected output to contain %q, got %q", tt.contains, string(output))
			}
		})
	}
}
