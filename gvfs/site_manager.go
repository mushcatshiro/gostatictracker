package gvfs

import (
	"context"

	"github.com/spf13/afero"
)

type SiteConf struct {
	AllowedExts []string
}

type SiteManager struct {
	Fs           *afero.Afero
	StyleFs      *afero.Afero
	JsFs         *afero.Afero
	TmplFs       *afero.Afero
	Cfg          SiteConf
	Pages        map[string]Page
	PendingRepos map[string]string
}

func NewSiteManager(baseFs afero.Fs) SiteManager {
	p := make(map[string]Page)
	return SiteManager{
		Fs:    &afero.Afero{Fs: baseFs},
		Pages: p,
	}
}

func (s *SiteManager) BuildIndex(ctx context.Context, root string) error {
	// TODO add goroutines to speed up
	err := s.Walk(root, root)
	if err != nil {
		return err
	}
	// TODO process PendingRepos when goroutines are added
	return nil
}
