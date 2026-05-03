package server

import "html/template"

type tmplBlogPageMeta struct {
	Title            string
	HasMermaid       bool
	HasMathJax       bool
	IsDraft          bool
	IsPrivate        bool
	Tags             []string
	LastModifiedDate string
	CreateDate       string
	URL              string
}

type BaseTmplMeta struct {
	SiteName  string
	IsAuth    bool
	ShowError bool
	IsIndex   bool
	IsBlog    bool
	*ErrPageTmplMeta
	*BlogPageTmplMeta
}

type ErrPageTmplMeta struct {
	ErrorMessage string
}

type EditorTmplMeta struct {
	TextBody string
}

type BlogPageTmplMeta struct {
	InnerText  template.HTML
	HasMermaid bool
	HasMathJax bool
}
