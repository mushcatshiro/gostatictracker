package render

import (
	"embed"
	"fmt"
	"html/template"
	"io"

	"github.com/mushcatshiro/gostatictracker/models"
)

type RenderEngine struct {
	tmpls map[string]*template.Template
}

func NewRenderEngine(fs embed.FS) (*RenderEngine, error) {
	re := &RenderEngine{
		tmpls: make(map[string]*template.Template),
	}
	entries, err := fs.ReadDir("templates")
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		name := entry.Name()
		if name == "base.html" || entry.IsDir() {
			continue
		}
		tmpl, err := template.ParseFS(fs, "templates/base.html", "templates/"+name)
		if err != nil {
			return nil, err
		}
		re.tmpls[name] = tmpl
	}
	return re, nil
}

func (r *RenderEngine) verifyTemplate(tmplName string) bool {
	return r.tmpls[tmplName] != nil
}

func (r *RenderEngine) Render(w io.Writer, rc models.RenderMeta) error {
	if ok := r.verifyTemplate(rc.TemplateName); !ok {
		err := fmt.Errorf("template %s is not found", rc.TemplateName)
		return err
	}
	return r.tmpls[rc.TemplateName].ExecuteTemplate(w, "base.html", rc.Data)
}
