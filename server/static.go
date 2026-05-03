package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

func (s *Server) RegisterFileServer(mux *http.ServeMux) {
	// Root file server at "static/public"
	// /static/js/vendor.js maps to static/public/js/vendor.js
	publicFS, err := fs.Sub(assets, "static/public")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(publicFS))

	handler := s.secureFileServer(publicFS, fileServer)

	mux.Handle("/static/", http.StripPrefix("/static/", handler))
	/*
			mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		    r.URL.Path = "/static/favicon.ico"
		    handler.ServeHTTP(w, r)
		})
	*/
}

func (s *Server) secureFileServer(root fs.FS, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Clean(r.URL.Path)

		// BLOCK ACCESS TO TEMPLATES
		// If someone tries to access /static/index.html directly, reject it.
		if strings.HasSuffix(path, ".html") {
			s.handleError(w, r, 403, fmt.Sprintf("Forbidden Access to %s", path))
			return
		}

		// ... existing IsDir() and Open() checks ...
		f, err := root.Open(strings.TrimPrefix(path, "/"))
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			s.handleError(w, r, 404, fmt.Sprintf("File %s not found", path))
			return
		}
		defer f.Close()

		stat, _ := f.Stat()
		if stat.IsDir() {
			s.handleError(w, r, 404, fmt.Sprintf("File %s not valid", path))
			return
		}

		next.ServeHTTP(w, r)
	})
}
