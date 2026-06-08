package blog

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"

	"github.com/mushcatshiro/gostatictracker/models"
	"github.com/spf13/afero"
)

func (bm *BlogManager) CreateNewBlog(path string, content string) error {
	err := bm.Fs.MkdirAll(filepath.Dir(path), 0644)
	if err != nil {
		return err
	}
	file, err := bm.Fs.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func (bm *BlogManager) GetSpecificBlogByPath(path string, contentOnly bool) (models.BlogEntry, error) {
	var be models.BlogEntry
	be, exists := bm.BlogEntries[path]
	if !exists {
		return be, ErrPageDoesNotExists
	}

	if !contentOnly {
		content, err := afero.ReadFile(bm.Fs, be.Path)
		if err != nil {
			return be, fmt.Errorf("fail to read %s with error %v", path, err)
		}
		be.Content = content
		return be, nil
	}

	file, err := bm.Fs.Open(be.Path)
	if err != nil {
		return be, fmt.Errorf("failed to read %s", path)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for i := 0; i < be.FrontMatter.ContentSidx; i++ {
		_, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return be, fmt.Errorf("unexpected EOF for %s", path)
			}
			return be, fmt.Errorf("fail to read %s with error %v", path, err)
		}
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return be, fmt.Errorf("fail to read content with error %v", err)
	}
	be.Content = content
	return be, nil
}

func (bm *BlogManager) UpdatePage() {}

func (bm *BlogManager) DeletePage() {
	// rename with prefix "."
}

func (bm *BlogManager) CommitChange() {}
