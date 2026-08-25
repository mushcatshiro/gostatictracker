package blog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestCreateNewBlog(t *testing.T) {
	fs := afero.NewMemMapFs()

	basePath := "/blogs"
	if err := fs.MkdirAll(basePath, 0755); err != nil {
		t.Fatal(err)
	}

	bm, err := NewBlogManager(fs, basePath)
	if err != nil {
		t.Fatalf("NewBlogManager() error = %v", err)
	}

	content := "# My First Blog\n\nHello, world!"
	fname := "my-first-blog"

	if err := bm.CreateNewBlog(fname, content); err != nil {
		t.Fatalf("CreateNewBlog() error = %v", err)
	}

	indexPath := filepath.Join(basePath, fname, "index.md")

	// Verify the file exists.
	info, err := fs.Stat(indexPath)
	if err != nil {
		t.Fatalf("index.md does not exist: %v", err)
	}

	if info.IsDir() {
		t.Fatal("index.md is a directory")
	}

	// Verify the permissions.
	if got, want := info.Mode().Perm(), os.FileMode(0644); got != want {
		t.Errorf("file permissions = %o, want %o", got, want)
	}

	// Verify the contents.
	got, err := afero.ReadFile(fs, indexPath)
	if err != nil {
		t.Fatalf("failed to read index.md: %v", err)
	}

	if string(got) != content {
		t.Errorf("file contents = %q, want %q", string(got), content)
	}
}

func TestValidateBlogName(t *testing.T) {
	tests := []struct {
		name    string
		fname   string
		wantErr bool
	}{
		{
			name:  "valid name",
			fname: "my-first-blog",
		},
		{
			name:  "valid name with underscore",
			fname: "my_first_blog",
		},
		{
			name:    "empty",
			fname:   "",
			wantErr: true,
		},
		{
			name:    "dot",
			fname:   ".",
			wantErr: true,
		},
		{
			name:    "dot dot",
			fname:   "..",
			wantErr: true,
		},
		{
			name:    "parent traversal",
			fname:   "../secret",
			wantErr: true,
		},
		{
			name:    "nested path",
			fname:   "foo/bar",
			wantErr: true,
		},
		{
			name:    "backslash traversal",
			fname:   `..\secret`,
			wantErr: true,
		},
		{
			name:    "absolute path",
			fname:   "/tmp/evil",
			wantErr: true,
		},
		{
			name:    "uppercase",
			fname:   "My-Blog",
			wantErr: true,
		},
		{
			name:    "spaces",
			fname:   "my blog",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBlogName(tt.fname)

			if tt.wantErr && err == nil {
				t.Errorf("validateBlogName(%q) expected error", tt.fname)
			}

			if !tt.wantErr && err != nil {
				t.Errorf("validateBlogName(%q) unexpected error: %v", tt.fname, err)
			}
		})
	}
}

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
