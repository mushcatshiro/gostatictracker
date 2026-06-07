package server

import "html/template"

type tmplBlogEntryMeta struct {
	Title            string
	HasMermaid       bool
	HasMathJax       bool
	IsDraft          bool
	IsPrivate        bool
	Tags             []string
	LastModifiedDate string
	CreateDate       string
	URL              string
	EditURL          string
}

type BaseTmplMeta struct {
	SiteName  string
	IsAuth    bool
	ShowError bool
	IsIndex   bool
	IsBlog    bool
	IsEditor  bool
	*ErrPageTmplMeta
	*BlogPageTmplMeta
}

type ErrPageTmplMeta struct {
	ErrorMessage string
}

type EditorTmplMeta struct {
	FrontMatter []byte
	TextBody template.HTML
}

type BlogPageTmplMeta struct {
	InnerText  template.HTML
	HasMermaid bool
	HasMathJax bool
}

type BlogListTmplMeta struct {
	IsAuth bool
	Tbem   []tmplBlogEntryMeta
}
