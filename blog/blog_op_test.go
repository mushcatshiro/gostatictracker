package blog

import (
	"strings"
	"testing"
)

func TestExtractFrontMatter(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		title   string
	}{
		{
			name:    "Valid YAML",
			input:   "---\ntitle: Hello\n---\nContent here",
			wantErr: false,
			title:   "Hello",
		},
		{
			name:    "Missing Closing Delimiter",
			input:   "---\ntitle: Hello\nContent without end",
			wantErr: true,
		},
		{
			name:    "No Front Matter",
			input:   "Just regular content",
			wantErr: true,
			title:   "",
		},
		{
			name:    "Exceeds Max Lines",
			input:   "---\n" + strings.Repeat("key: value\n", 40) + "---\n",
			wantErr: true,
		},
		{
			name:    "Empty File",
			input:   "",
			wantErr: true,
			title:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			meta, err := ExtractFrontMatter(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractFrontMatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if meta.Title != tt.title {
				t.Errorf("Expected title %q, got %q", tt.title, meta.Title)
			}
		})
	}
}
