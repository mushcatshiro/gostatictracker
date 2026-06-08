package blog

import (
	"context"
	"errors"
	"runtime"
	"sync"

	"github.com/mushcatshiro/gostatictracker/models"
	"github.com/spf13/afero"
	"golang.org/x/sync/errgroup"
)

var ErrPageDoesNotExists = errors.New("page does not exists")

type BlogConf struct {
	AllowedExts []string
}

type BlogManager struct {
	Fs           *afero.Afero
	Cfg          BlogConf
	BlogEntries  map[string]models.BlogEntry
	PendingRepos map[string]string
	Mu           sync.RWMutex
}

func NewBlogManager(baseFs afero.Fs) BlogManager {
	p := make(map[string]models.BlogEntry)
	return BlogManager{
		Fs:          &afero.Afero{Fs: baseFs},
		BlogEntries: p,
	}
}

func (bm *BlogManager) BuildIndex(ctx context.Context, root string) error {
	// TODO add goroutines to speed up
	err := bm.Walk(root, root)
	if err != nil {
		return err
	}
	// TODO process PendingRepos when goroutines are added
	// pgs, err := json.MarshalIndent(s.Pages, "", "  ")
	// fmt.Println(string(pgs))
	return nil
}

func (bm *BlogManager) BuildIndexFast(ctx context.Context, root string) error {
	jobs := make(chan string, 100) // Channel of paths to process
	g, ctx := errgroup.WithContext(ctx)

	// 1. Workers
	for i := 0; i < runtime.NumCPU(); i++ {
		g.Go(func() error {
			for path := range jobs {
				// Each worker handles the I/O and parsing
				if err := bm.processSinglePage(root, path); err != nil {
					return err
				}
			}
			return nil
		})
	}

	// 2. Producer (Your original Walk logic, slightly modified)
	g.Go(func() error {
		defer close(jobs)
		return bm.ConcurrentWalk(root, root, jobs)
	})

	return g.Wait()
}

func (bm *BlogManager) processSinglePage(root, path string) error {
	return nil
}

func (bm *BlogManager) ConcurrentWalk(root, dir string, jobs chan string) error {
	return nil
}
