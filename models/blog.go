package models

import "time"

type iBlogRecord struct {
	Path         string
	FullURL      string
	SideRepoPath string
	Content      []byte
	iFrontMatter
	EditUrl string
}

type iFrontMatter struct {
	Title            string   `toml:"title"`
	HasMermaid       bool     `toml:"mermaid"`
	HasMathJax       bool     `toml:"math"`
	IsDraft          bool     `toml:"draft"`
	IsPrivate        bool     `toml:"private"`
	Tags             []string `toml:"tags"`
	LastModifiedDate string   `toml:"lastmodified"`
	CreateDate       string   `toml:"date"`
	ContentSidx      int
}

type BlogRecord struct {
	ID         int64
	Title      string
	URL        string
	InsertTime *time.Time
	Metadata   string
}

func (r *Record) ToBlogRecord() BlogRecord {
	return BlogRecord{
		ID:         r.ID,
		Title:      r.Title,
		URL:        r.URL,
		InsertTime: r.InsertTime,
		Metadata:   r.Metadata,
	}
}
