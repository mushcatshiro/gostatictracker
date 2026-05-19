package server

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
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

func (s *Server) renderEditorView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// distinguish between code and markdonw
		pathVar := r.PathValue("title")
		etm := EditorTmplMeta{}
		if pathVar != "" {
			p, err := s.blogSiteManager.GetSpecificPageByPath(pathVar)
			if err != nil {
				s.handleError(w, r, http.StatusBadRequest, err.Error())
				return
			}
			// TODO
			// need to include frontmatter, need to handle non standard utf-8 char
			etm.TextBody = string(p.Content)
		} else {
			etm.TextBody = editorTextBody
		}

		w.Header().Set("Content-Type", "text/html")
		tmpl.ExecuteTemplate(w, "editor", etm)
	}
}

func (s *Server) editorSubmit() http.HandlerFunc {
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
