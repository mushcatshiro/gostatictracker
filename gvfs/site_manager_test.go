package gvfs

import (
	"context"
	"testing"

	"github.com/spf13/afero"
)

func TestBuildIndex(t *testing.T) {
	fs := afero.NewMemMapFs()
	appFS := &afero.Afero{Fs: fs}

	files := map[string]string{
		"content/nested/folder/keep-going/index.md": "---\ntitle: Nested\n---\nBody",
		"content/nested/folder/keep-going2/index.md": "---\ntitle: Nested2\n---\nBody",
		"content/blog/post-1/index.md":              "---\ntitle: Post One\n---\nBody",
		"content/blog/post-1/image.png":             "fake-image-data",
		"content/docs/page-a.md":                    "---\ntitle: Alpha\n---\nBody",
		"content/docs/page-b.md":                    "---\ntitle: Beta\n---\nBody",
		"content/docs/siderepo/config.json":         "{}",
	}

	for path, content := range files {
		appFS.WriteFile(path, []byte(content), 0644)
	}

	sm := &SiteManager{
		Fs:    appFS,
		Pages: make(map[string]Page),
	}

	err := sm.BuildIndex(context.Background(), "content")
	if err != nil {
		t.Fatalf("BuildIndex failed: %v", err)
	}

	tests := []struct {
		url      string
		expected string
	}{
		{"blog/post-1", "Post One"},
		{"docs/page-a", "Alpha"},
		{"nested/folder/keep-going", "Nested"},
		{"nested/folder/keep-going2", "Nested2"},
	}
	for k, v := range sm.Pages {
		t.Logf("%s, %+v\n", k, v)
	}

	for _, tt := range tests {
		page, exists := sm.Pages[tt.url]
		if !exists {
			t.Errorf("Expected page at URL %q not found", tt.url)
			continue
		}
		if page.Meta.Title != tt.expected {
			t.Errorf("URL %q: expected title %q, got %q", tt.url, tt.expected, page.Meta.Title)
		}
	}
}
