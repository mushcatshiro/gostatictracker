package models

import (
	"errors"
	"html/template"
)

type FormType int

const (
	SimpleForm FormType = iota
	FullForm
	TaskForm
	IdeaForm
	ReminderForm
	EditForm
)

type RenderMeta struct {
	TemplateName string
	Data         any
}

type BaseRenderMeta struct {
	SiteName  string
	ShowError bool
	IsAuth    bool
}

type ErrorRenderMeta struct {
	BaseRenderMeta
	ErrorMessage string
}

type BookmarkletIndexRenderMeta struct {
	BaseRenderMeta
	Bookmarklets []Bookmarklet
	SetupScript  string
}

type BlogEntryMeta struct {
	BaseRenderMeta
	BlogEntry
	EditURL   string
	InnerText template.HTML
}

type BlogIndexRenderMeta struct {
	BaseRenderMeta
	BlogEntryMetas []BlogEntryMeta
}

type BlogRenderMeta struct {
	BaseRenderMeta
	InnerText  template.HTML
	HasMermaid bool
	HasMathJax bool
}

type EditorRenderMeta struct {
	BaseRenderMeta
	FrontMatter template.HTML
	TextBody    template.HTML
}

type FormRenderMeta struct {
	BaseRenderMeta
	PostEndpoint string
	FormTitle    string
	IsSimple     bool
	IsFull       bool
	IsTask       bool
	IsIdea       bool
	IsReminder   bool
	IsEdit       bool
	Event
}

func NewFormMeta(
	baseRenderMeta BaseRenderMeta, endpoint string, formType FormType, e *Event,
) (FormRenderMeta, error) {
	frm := FormRenderMeta{
		BaseRenderMeta: baseRenderMeta,
		PostEndpoint:   endpoint,
	}
	switch formType {
	case SimpleForm:
		frm.FormTitle = "Quick Form"
		frm.IsSimple = true
	case FullForm:
		frm.FormTitle = "Full Form"
		frm.IsFull = true
	case TaskForm:
		frm.FormTitle = "Task Form"
		frm.IsTask = true
	case IdeaForm:
		frm.FormTitle = "Idea Form"
		frm.IsIdea = true
	case ReminderForm:
		frm.FormTitle = "Reminder Form"
		frm.IsReminder = true
	case EditForm:
		frm.FormTitle = "Edit Form"
		frm.IsEdit = true
		if e == nil {
			return frm, errors.New("event must not be null")
		}
		frm.Event = *e
	default:
		return frm, errors.New("undefined form")
	}
	return frm, nil
}

type GanttRenderMeta struct{}

type ReminderRenderMeta struct{}

type CalendarRenderMeta struct{}

type KanbanRenderMeta struct{}

type ChecklistRenderMeta struct{}

type TableFormRenderMeta struct{}

type GroupEntry struct {
	GroupName   string
	DefaultMode string
	URL         string
}

type GroupIndexRenderMeta struct {
	BaseRenderMeta
	GroupEntries []GroupEntry
}
