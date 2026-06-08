package blog

import (
	"context"
	"fmt"
	"testing"

	"github.com/mushcatshiro/gostatictracker/models"
	"github.com/spf13/afero"
)

func TestBuildIndex(t *testing.T) {
	fs := afero.NewMemMapFs()
	appFS := &afero.Afero{Fs: fs}

	files := map[string]string{
		"content/nested/folder/keep-going/index.md":  "---\ntitle: Nested\n---\nBody",
		"content/nested/folder/keep-going2/index.md": "---\ntitle: Nested2\n---\nBody",
		"content/blog/post-1/index.md":               "---\ntitle: Post One\n---\nBody",
		"content/blog/post-1/image.png":              "fake-image-data",
		"content/docs/page-a.md":                     "---\ntitle: Alpha\n---\nBody",
		"content/docs/page-b.md":                     "---\ntitle: Beta\n---\nBody",
		"content/docs/siderepo/config.json":          "{}",
	}

	for path, content := range files {
		appFS.WriteFile(path, []byte(content), 0644)
	}

	sm := NewBlogManager(fs)

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
	for k, v := range sm.BlogEntries {
		t.Logf("%s, %+v\n", k, v)
	}

	for _, tt := range tests {
		page, exists := sm.BlogEntries[tt.url]
		if !exists {
			t.Errorf("Expected page at URL %q not found", tt.url)
			continue
		}
		if page.FrontMatter.Title != tt.expected {
			t.Errorf("URL %q: expected title %q, got %q", tt.url, tt.expected, page.FrontMatter.Title)
		}
	}
}

func TTestBuildIndexRaceCondition(t *testing.T) {
	fs := afero.NewMemMapFs()
	appFS := &afero.Afero{Fs: fs}

	// Create 100 dummy files to ensure multiple workers are busy
	for i := range 100 {
		path := fmt.Sprintf("content/post-%d.md", i)
		appFS.WriteFile(path, []byte("---\ntitle: Race Test\n---\nBody"), 0644)
	}

	sm := NewBlogManager(fs)
	sm.BlogEntries = make(map[string]models.BlogEntry)

	// This will trigger the race detector if your map is not protected
	err := sm.BuildIndex(context.Background(), "content")
	if err != nil {
		t.Fatalf("BuildIndex failed: %v", err)
	}
}
