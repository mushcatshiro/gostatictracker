package gvfs

import (
	"context"
	"runtime"
	"sync"

	"github.com/spf13/afero"
	"golang.org/x/sync/errgroup"
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
	mu           sync.RWMutex
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

func (s *SiteManager) BuildIndexFast(ctx context.Context, root string) error {
	jobs := make(chan string, 100) // Channel of paths to process
	g, ctx := errgroup.WithContext(ctx)

	// 1. Workers
	for i := 0; i < runtime.NumCPU(); i++ {
		g.Go(func() error {
			for path := range jobs {
				// Each worker handles the I/O and parsing
				if err := s.processSinglePage(root, path); err != nil {
					return err
				}
			}
			return nil
		})
	}

	// 2. Producer (Your original Walk logic, slightly modified)
	g.Go(func() error {
		defer close(jobs)
		return s.ConcurrentWalk(root, root, jobs)
	})

	return g.Wait()
}

func (s *SiteManager) processSinglePage(root, path string) error {
	return nil
}

func (s *SiteManager) ConcurrentWalk(root, dir string, jobs chan string) error {
	return nil
}
