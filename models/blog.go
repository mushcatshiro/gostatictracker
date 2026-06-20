package models


type BlogEntry struct {
	Path         string
	FullURL      string
	SideRepoPath string
	Content      []byte
	FrontMatter
	EditUrl      string
	Assets       []AssetsMeta
}

type FrontMatter struct {
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

type AssetsMeta struct {
	Fname string
}
