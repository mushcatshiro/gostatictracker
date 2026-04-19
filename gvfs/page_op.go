package gvfs

import "path/filepath"

func (s *SiteManager) CreatePage(path string, content string) error {
	err := s.Fs.MkdirAll(filepath.Dir(path), 0644)
	if err != nil {
		return err
	}
	file, err := s.Fs.Create(path)
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

func (s *SiteManager) ReadPage() {}

func (s *SiteManager) UpdatePage() {}

func (s *SiteManager) DeletePage() {
	// rename with prefix "."
}

func (s *SiteManager) CommitChange() {}
