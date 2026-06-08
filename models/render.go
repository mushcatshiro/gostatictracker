package models

import "html/template"

type RenderMeta struct {
	TemplateName string
	Data         any
}

type BaseRenderMeta struct {
	SiteName      string
	ShowError     bool
	IsAuth        bool
	IsBlog        bool
	IsEditor      bool
	IsIndex       bool
	IsBookmarkler bool
}

type ErrorRender struct {
	BaseRenderMeta
	ErrorMessage string
}

type BookmarkletIndexRenderMeta struct {
	BaseRenderMeta
	Bookmarklets []Bookmarklet
	ServerDomain string
	IsAuth       bool
}

type BlogEntryMeta struct {
	BlogEntry
	URL       string
	EditorURL string
}

type BlogIndexRenderMeta struct {
	BaseRenderMeta
	BlogEntries []BlogEntryMeta
}

type BlogRenderMeta struct {
	BaseRenderMeta
	InnerText  template.HTML
	HasMermaid bool
	HasMathJax bool
}

type EditorRenderMeta struct {
	BaseRenderMeta
	FrontMatter []byte
	TextBody    template.HTML
}

type FormRenderMeta struct{}

type GanttRenderMeta struct{}

type ReminderRenderMeta struct{}

type CalendarRenderMeta struct{}

type KanbanRenderMeta struct{}

type ChecklistRenderMeta struct{}

type TodoRenderMeta struct{}

type TableFormRenderMeta struct{}

var TemplateMetaMap = map[string]any{
	"blogIndex": BlogIndexRenderMeta{},
}
