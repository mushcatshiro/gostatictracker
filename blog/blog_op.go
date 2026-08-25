package blog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/afero"
)

var validBlogName = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*$`)

func validateBlogName(fname string) error {
	if filepath.IsAbs(fname) {
		return errors.New("blog name must not be an absolute path")
	}

	if !validBlogName.MatchString(fname) {
		return fmt.Errorf(
			"invalid blog name %q: only lowercase letters, numbers, '-' and '_' are allowed",
			fname,
		)
	}
	return nil
}

func (bm *BlogManager) CreateNewBlog(fname string, content string) error {
	if err := validateBlogName(fname); err != nil {
		return err
	}
	blogPath := filepath.Join(bm.basePath, fname)

	if err := bm.Fs.MkdirAll(blogPath, 0755); err != nil {
		return err
	}

	indexPath := filepath.Join(blogPath, "index.md")

	file, err := bm.Fs.OpenFile(
		indexPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func (bm *BlogManager) GetSpecificBlogByPath(fname string) (string, error) {
	if err := validateBlogName(fname); err != nil {
		return "", err
	}

	filePath := filepath.Join(bm.basePath, fname, "index.md")

	content, err := afero.ReadFile(bm.Fs, filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("blog post not found: %s", fname)
		}
		return "", fmt.Errorf("failed to read blog post: %w", err)
	}

	return string(content), nil
}

func (bm *BlogManager) UpdatePage() {}

func (bm *BlogManager) DeletePage() {
	// rename with prefix "."
}

func (bm *BlogManager) CommitChange() {}

func ToSlug(title string) string {
	// 1. Convert to lowercase
	slug := strings.ToLower(title)

	// 2. Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// 3. Remove all characters except a-z, 0-9, and hyphens
	// (Included 0-9 as keeping numbers in slugs is standard practice)
	regSafeChars := regexp.MustCompile("[^a-z0-9-]+")
	slug = regSafeChars.ReplaceAllString(slug, "")

	// 4. Remove consecutive hyphens (e.g., "my---post" becomes "my-post")
	regMultipleHyphens := regexp.MustCompile("-+")
	slug = regMultipleHyphens.ReplaceAllString(slug, "-")

	// 5. Trim leading and trailing hyphens (in case special chars were at the ends)
	slug = strings.Trim(slug, "-")

	return slug
}
