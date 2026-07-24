package server

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/mushcatshiro/gostatictracker/blog"
)

const editorTextBody = `---
Title = ""
HasMermaid = false
HasMathJax = false
IsDraft = true
IsPrivate = true
Tags = ["",]
LastModifiedDate = ""
CreateDate = ""
---`

/*
func (s *Server) processAssetUpload(mpf *multipart.Form) (string, error) {
	var resultList strings.Builder

	if files, ok := mpf.File["images"]; ok {
		for _, fileHeader := range files {
			// to change filePath
			safeFname := filepath.Base(fileHeader.Filename)
			filePath := "/static/uploads/" + safeFname
			resultList.WriteString(generateAssetItem(filePath))
		}
	}

	linksRaw := r.FormValue("links")
	links := strings.Split(linksRaw, "\n")
	for _, link := range links {
		trimmedLink := strings.TrimSpace(link)
		if trimmedLink != "" {
			resultList.WriteString(generateAssetItem(trimmedLink))
		}
	}
	return resultList.String(), nil
}
*/

func (s *Server) handleAssetUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const maxMemory = 10 << 20
		if err := r.ParseMultipartForm(maxMemory); err != nil {
			http.Error(w, "Upload processing failed", http.StatusBadRequest)
			return
		}

		var resultList strings.Builder

		if files, ok := r.MultipartForm.File["images"]; ok {
			for _, fileHeader := range files {
				// to change filePath
				safeFname := filepath.Base(fileHeader.Filename)
				filePath := "/static/uploads/" + safeFname
				resultList.WriteString(generateAssetItem(filePath))
			}
		}

		linksRaw := r.FormValue("links")
		links := strings.Split(linksRaw, "\n")
		for _, link := range links {
			trimmedLink := strings.TrimSpace(link)
			if trimmedLink != "" {
				resultList.WriteString(generateAssetItem(trimmedLink))
			}
		}

		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, resultList.String())
	}
}

func generateAssetItem(url string) string {
	return fmt.Sprintf(`<li class="asset-item"
	onclick="insertMarkdownRef('%s')"
	style="cursor:pointer; padding: 8px; border-bottom: 1px solid #eee; font-size: 14px;">
		%s
</li>`, url, url)
}

func (s *Server) processEditorView(path string) (EditorTmplMeta, error) {
	etm := EditorTmplMeta{}
	if path != "" {
		// TODO
		// to handle non standard utf-8 char
		// to include all assets
		p, err := s.blogSiteManager.GetSpecificPageByPath(path, false)
		if errors.Is(err, blog.ErrPageDoesNotExists) {
			return etm, errors.New("page does not exists")
		}
		if err != nil {
			return etm, err
		}
		etm.TextBody = template.HTML(p.Content)
		return etm, nil
	}
	etm.TextBody = template.HTML([]byte(editorTextBody))
	return etm, nil
}

func (s *Server) renderEditorView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isHTMX := r.Header.Get("HX-Request") == "true"
		if !isHTMX {
			s.handleError(w, r, 404, "Forbidden")
		}
		path := r.PathValue("title")
		etm, err := s.processEditorView(path)
		if err != nil {
			s.handleError(w, r, 404, err.Error())
			return
		}
		w.Header().Set("Content-Type", "text/html")
		s.tmpl.ExecuteTemplate(w, "editor", etm)
	}
}

func (s *Server) processEditorSubmit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			s.handleError(w, r, 405, "Method not allowed")
			return
		}
		// get form
		// check if exists?
		// save as tmp?
	}
}
