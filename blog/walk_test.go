package blog

import (
	"testing"

	"github.com/spf13/afero"
)

func TestWalkSkipsFolders(t *testing.T) {
	fs := afero.NewMemMapFs()
	appFS := &afero.Afero{Fs: fs}

	// Create a file inside a skipped directory
	appFS.MkdirAll(".git", 0755)
	appFS.WriteFile(".git/config", []byte("should skip"), 0644)
	appFS.WriteFile("index.md", []byte("---\ntitle: Home\n---"), 0644)

	sm := NewBlogManager(appFS)

	err := sm.Walk(".", ".")
	if err != nil {
		t.Fatal(err)
	}

	// Assert: Map should have the home page but NOT the .git config
	if len(sm.BlogEntries) != 1 {
		t.Errorf("Expected 1 page, got %d", len(sm.BlogEntries))
	}
}

func TestWalk(t *testing.T) {
	// to include assets
}
