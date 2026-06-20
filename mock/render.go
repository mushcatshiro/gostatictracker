package mock

import (
	"github.com/mushcatshiro/gostatictracker/models"
)

var MockBaseRenderMeta = models.BaseRenderMeta{
	SiteName: "Test Site",
}

var BlogEntry1 = models.BlogEntry{
	Path:    "/path1/file1.md",
	FullURL: "/blog/path1/file1",
	FrontMatter: models.FrontMatter{
		Title:            "file1",
		LastModifiedDate: "2025-01-01T00:00:00+08:00",
	},
}

var BlogEntry2 = models.BlogEntry{
	Path:    "/path2/file2.md",
	FullURL: "/blog/path2/file2",
	FrontMatter: models.FrontMatter{
		Title:            "file2",
		LastModifiedDate: "2025-01-02T00:00:00+08:00",
	},
}

var BlogEntryMetas = []models.BlogEntryMeta{
	{BlogEntry: BlogEntry1, EditorURL: "/editor/path1/file1"},
	{BlogEntry: BlogEntry1, EditorURL: "/editor/path2/file2"},
}

var quickFormMeta models.FormRenderMeta
var fullFormMeta models.FormRenderMeta
var taskFormMeta models.FormRenderMeta
var ideaFormMeta models.FormRenderMeta
var reminderFormMeta models.FormRenderMeta
var editFormMeta models.FormRenderMeta
var FormTemplateMetaMap map[string]any
var TemplateMetaMap map[string]any

func init() {
	quickFormMeta, _ = models.NewFormMeta(MockBaseRenderMeta, "/form", models.SimpleForm, nil)
	fullFormMeta, _ = models.NewFormMeta(MockBaseRenderMeta, "/form", models.FullForm, nil)
	taskFormMeta, _ = models.NewFormMeta(MockBaseRenderMeta, "/form", models.TaskForm, nil)
	ideaFormMeta, _ = models.NewFormMeta(MockBaseRenderMeta, "/form", models.IdeaForm, nil)
	reminderFormMeta, _ = models.NewFormMeta(MockBaseRenderMeta, "/form", models.ReminderForm, nil)
	editFormMeta, _ = models.NewFormMeta(MockBaseRenderMeta, "/form", models.EditForm, nil)
	FormTemplateMetaMap = map[string]any{
		"quick":    quickFormMeta,
		"full":     fullFormMeta,
		"task":     taskFormMeta,
		"idea":     ideaFormMeta,
		"reminder": reminderFormMeta,
		"edit":     editFormMeta,
	}
	// doesnt include form
	TemplateMetaMap = map[string]any{
		"bloglist.html": models.BlogIndexRenderMeta{
			BaseRenderMeta: MockBaseRenderMeta,
			BlogEntryMetas: BlogEntryMetas,
		},
		"error.html": models.ErrorRenderMeta{
			BaseRenderMeta: MockBaseRenderMeta,
			ErrorMessage: "404 Not Found",
		},
	}
}
