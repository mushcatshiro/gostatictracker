package mock

import (
	"html/template"

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

const blogInnerText = `<h1>A Perspective Change</h1>
<p>
  <i>
    <time datetime="2025-08-23">
      23 Aug, 2025
    </time>
  </i>
</p>`

var BlogEntryMetas = []models.BlogEntryMeta{
	{
		BaseRenderMeta: MockBaseRenderMeta,
		BlogEntry:      BlogEntry1,
		EditURL:        "/editor/path1/file1",
		InnerText:      blogInnerText,
	},
	{
		BaseRenderMeta: MockBaseRenderMeta,
		BlogEntry:      BlogEntry1,
		EditURL:        "/editor/path2/file2",
		InnerText:      blogInnerText,
	},
}

var Bookmarklet1 = models.Bookmarklet{
	Title:       "bkmk1",
	Description: "bkmk1",
	Group:       "bookmarklet",
	URL:         "https://bkmk1/bkmk1",
}

var Bookmarklet2 = models.Bookmarklet{
	Title:       "bkmk2",
	Description: "bkmk2",
	Group:       "bookmarklet",
	URL:         "https://bkmk2/bkmk2",
}

var GroupEntry1 = models.GroupEntry{
	GroupName:   "bookmarklet",
	DefaultMode: "bookmarkletlist.html",
	URL:         "https://localhost/bookmarkletlist",
}

var GroupEntry2 = models.GroupEntry{
	GroupName:   "blog",
	DefaultMode: "bloglist.html",
	URL:         "https://localhost/bloglist",
}

var FrontMatter = template.HTML(
	`+++
title = "A Study of bfloat16 for deep learning training"
+++
`)

var TextBody = template.HTML(
	`This paper presents the first comprehensive empirical study demonstrating
the efficacy of the Brain Floating Point (BFLOAT16) half-precision format for
Deep Learning training across image classification, speech recognition,
language modeling, generative networks and industrial recommendation systems.`)

var quickFormMeta models.FormRenderMeta
var fullFormMeta models.FormRenderMeta
var taskFormMeta models.FormRenderMeta
var ideaFormMeta models.FormRenderMeta
var reminderFormMeta models.FormRenderMeta
var editFormMeta models.FormRenderMeta
var FormTemplateMetaMap map[string]any
var TemplateMetaMap map[string]func(isAuth bool) any

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
	TemplateMetaMap = map[string]func(isAuth bool) any{
		"blog.html": func(isAuth bool) any {
			meta := BlogEntryMetas[0]
			meta.IsAuth = isAuth
			return meta
		},
		"bloglist.html": func(isAuth bool) any {
			meta := models.BlogIndexRenderMeta{
				BaseRenderMeta: MockBaseRenderMeta,
				BlogEntryMetas: BlogEntryMetas,
			}
			meta.IsAuth = isAuth
			return meta
		},
		"bookmarkletlist.html": func(isAuth bool) any {
			meta := models.BookmarkletIndexRenderMeta{
				BaseRenderMeta: MockBaseRenderMeta,
				Bookmarklets:   []models.Bookmarklet{Bookmarklet1, Bookmarklet2},
			}
			meta.IsAuth = isAuth
			return meta
		},
		"editor.html": func(isAuth bool) any {
			meta := models.EditorRenderMeta{
				BaseRenderMeta: MockBaseRenderMeta,
				FrontMatter:    FrontMatter,
				TextBody:       TextBody,
			}
			meta.IsAuth = isAuth
			return meta
		},
		"error.html": func(isAuth bool) any {
			meta := models.ErrorRenderMeta{
				BaseRenderMeta: MockBaseRenderMeta,
				ErrorMessage:   "404 Not Found",
			}
			meta.IsAuth = isAuth
			return meta
		},
		"groupindex.html": func(isAuth bool) any {
			meta := models.GroupIndexRenderMeta{
				BaseRenderMeta: MockBaseRenderMeta,
				GroupEntries: []models.GroupEntry{
					GroupEntry1, GroupEntry2,
				},
			}
			meta.IsAuth = isAuth
			return meta
		},
	}
}
