package render

import (
	"fmt"
	"html/template"
	"io"

	"github.com/mushcatshiro/gostatictracker/models"
)

type RenderEngine struct {
	tmpl *template.Template
}

func NewRenderEngine(t *template.Template) *RenderEngine {
	return &RenderEngine{
		tmpl: t,
	}
}

func (r *RenderEngine) verifyTemplate(tmplName string) bool {
	return r.tmpl.Lookup(tmplName) != nil
}

func (r *RenderEngine) Render(w io.Writer, rc models.RenderMeta) error {
	if ok := r.verifyTemplate(rc.TemplateName); !ok {
		err := fmt.Errorf("template %s is not found", rc.TemplateName)
		return err
	}
	return r.tmpl.ExecuteTemplate(w, rc.TemplateName, rc.Data)
}
