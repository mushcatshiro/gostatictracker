package server

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed static/*
var assets embed.FS

var tmpl *template.Template

func InitTemplates() {
	// Parse all .html files inside the templates subdirectory
	var err error
	tmpl, err = template.ParseFS(assets, "static/templates/*.html")
	if err != nil {
		panic(err)
	}
}

func RenderHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{"Title": "My Startup"}
	// Executes index.html from the embedded FS
	tmpl.ExecuteTemplate(w, "index.html", data)
}
