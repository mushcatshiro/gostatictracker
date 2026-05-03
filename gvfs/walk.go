package gvfs

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

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
func (s *SiteManager) Walk(root, dir string) error {
	entries, err := afero.ReadDir(s.Fs, dir)
	if err != nil {
		return err
	}

	var nextLevelDirs []string
	proceedNextDepth := true
	indexLeaf := false

	for _, e := range entries {
		// fmt.Printf("processing %s\n", e.Name())
		if e.IsDir() {
			if slices.Contains(SkippedFolderNames, e.Name()) {
				// fmt.Printf("skipping %s\n", e.Name())
				continue
			}
			nextLevelDirs = append(nextLevelDirs, filepath.Join(dir, e.Name()))
			// fmt.Printf("appending %s %s\n", dir, e.Name())
			continue
		}

		fname := e.Name()
		if fname == "index.md" || fname == "index.html" {
			// fmt.Print("found index leaf\n")
			indexLeaf = true
			break
		}

		ext := filepath.Ext(e.Name())
		if ext == ".md" || ext == ".html" {
			proceedNextDepth = false
			break
		}
	}
	if proceedNextDepth && !indexLeaf {
		for _, dir := range nextLevelDirs {
			if err := s.Walk(root, dir); err != nil {
				return err
			}
		}
		return nil
	}
	if indexLeaf {
		return s.AssembleIndexLeaf(root, dir, entries)
	}

	return s.AssembleBundle(root, dir, entries)
}

// Url prefix referring to "content/blog/post-1/index.md" -> "blog/post-1"
func (s *SiteManager) generateIndexLeafUrl(root, path string) string {
	relPath, _ := filepath.Rel(root, path)
	return filepath.ToSlash(filepath.Dir(relPath))
}

func (s *SiteManager) generateBundleUrl(root, path, ext string) string {
	relPath, _ := filepath.Rel(root, path)
	noExtPath := strings.ReplaceAll(relPath, ext, "")
	return filepath.ToSlash(noExtPath)
}

func (s *SiteManager) AssembleIndexLeaf(root, dir string, entries []os.FileInfo) error {
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
	return s.UpdatePages(indexPath, sideRepoPath, s.generateIndexLeafUrl(root, indexPath))
}

func (s *SiteManager) AssembleBundle(root, dir string, entries []os.FileInfo) error {
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
			key := s.generateBundleUrl(root, indexPath, ext)
			if err := s.UpdatePages(indexPath, sideRepoPath, key); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *SiteManager) UpdatePages(indexPath, sideRepoPath, key string) error {
	file, err := s.Fs.Open(indexPath)
	if err != nil {
		return err
	}
	pageMeta, err := extractFrontMatter(file)
	if err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}
	s.Pages[key] = Page{
		Path:         indexPath,
		FullURL:      fmt.Sprintf("blog/%s", key),
		Meta:         pageMeta,
		SideRepoPath: sideRepoPath,
	}
	return nil
}
