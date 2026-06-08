package blog

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mushcatshiro/gostatictracker/models"
	"github.com/spf13/afero"
)

var SkippedFolderNames = []string{".git", "siderepo", "node_modules", ".venv"}

// there are a few possibilities
// 1. a folder with many subfolders (only folders allowed)
// 2. a folder with one index.md or index.html, with optional `siderepo`
// 3. a folder with many *.md and *.html with optional `siderepo/`
// resources i.e. images, will only exists in 2. and 3.
// 1. to continue to the next depth
// 2. will prioritize index.html > index.md
// in 3. it is possible to mix *.md and *.html
func (bm *BlogManager) Walk(root, dir string) error {
	entries, err := afero.ReadDir(bm.Fs, dir)
	if err != nil {
		return err
	}

	var nextLevelDirs []string
	indexLeaf := false
	var effectiveEntries []os.FileInfo

	for _, e := range entries {
		if e.IsDir() {
			if slices.Contains(SkippedFolderNames, e.Name()) {
				continue
			}
			nextLevelDirs = append(nextLevelDirs, filepath.Join(dir, e.Name()))
			continue
		}

		fname := e.Name()
		if fname == "_index.md" {
			continue
		}
		effectiveEntries = append(effectiveEntries, e)
		if fname == "index.md" || fname == "index.html" {
			indexLeaf = true
			continue
		}

	}
	if len(nextLevelDirs) > 0 {
		for _, dir := range nextLevelDirs {
			if err := bm.Walk(root, dir); err != nil {
				return err
			}
		}
	}
	if indexLeaf {
		return bm.AssembleIndexLeaf(root, dir, effectiveEntries)
	}
	return bm.AssembleBundle(root, dir, effectiveEntries)
}

// Url prefix referring to "content/blog/post-1/index.md" -> "blog/post-1"
func (bm *BlogManager) generateIndexLeafUrl(root, path string) string {
	relPath, _ := filepath.Rel(root, path)
	return filepath.ToSlash(filepath.Dir(relPath))
}

func (bm *BlogManager) generateBundleUrl(root, path, ext string) string {
	relPath, _ := filepath.Rel(root, path)
	noExtPath := strings.ReplaceAll(relPath, ext, "")
	return filepath.ToSlash(noExtPath)
}

func (bm *BlogManager) AssembleIndexLeaf(root, dir string, entries []os.FileInfo) error {
	var indexPath, sideRepoPath string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "index.") {
			indexPath = filepath.Join(dir, e.Name())
			continue
		}
		if e.IsDir() && e.Name() == "siderepo" {
			sideRepoPath = filepath.Join(dir, e.Name())
			continue
		}
	}
	return bm.UpdatePages(indexPath, sideRepoPath, s.generateIndexLeafUrl(root, indexPath))
}

func (bm *BlogManager) AssembleBundle(root, dir string, entries []os.FileInfo) error {
	var sideRepoPath string
	for _, e := range entries {
		if e.IsDir() && e.Name() == "siderepo" {
			sideRepoPath = filepath.Join(dir, e.Name())
			break
		}
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		if ext == ".md" || ext == ".html" {
			indexPath := filepath.Join(dir, e.Name()) // keeping naming convention the same
			key := bm.generateBundleUrl(root, indexPath, ext)
			if err := bm.UpdatePages(indexPath, sideRepoPath, key); err != nil {
				return err
			}
		}
	}
	return nil
}

func (bm *BlogManager) UpdatePages(indexPath, sideRepoPath, key string) error {
	file, err := bm.Fs.Open(indexPath)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w\n", indexPath, err)
	}
	defer file.Close()
	frontMatter, err := ExtractFrontMatter(file)
	if err != nil {
		return fmt.Errorf("failed to extract %s: %w\n", indexPath, err)
	}
	fullUrl, _ := url.JoinPath("/blog", key)
	editUrl, _ := url.JoinPath("/editor", key)
	bm.BlogEntries[key] = models.BlogEntry{
		Path:         indexPath,
		FullURL:      fullUrl,
		FrontMatter:  frontMatter,
		SideRepoPath: sideRepoPath,
		EditUrl:      editUrl,
	}
	return nil
}
