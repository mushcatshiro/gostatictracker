package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

func (s *Server) handleStatic(staticFS fs.FS, mux *http.ServeMux) http.HandlerFunc {
	// Root file server at "static/public"
	// /static/js/vendor.js maps to static/public/js/vendor.js
	publicFS, err := fs.Sub(staticFS, "static/public")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(publicFS))
	return s.secureFileServer(publicFS, fileServer)
	// mux.Handle("/static/", http.StripPrefix("/static/", handler))
}

func (s *Server) secureFileServer(root fs.FS, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Clean(r.URL.Path)

		// BLOCK ACCESS TO TEMPLATES
		// If someone tries to access /static/index.html directly, reject it.
		if strings.HasSuffix(strings.ToLower(path), ".html") {
			s.handleError(w, r, 403, fmt.Sprintf("Forbidden Access to %s", path))
			return
		}

		f, err := root.Open(strings.TrimPrefix(path, "/"))
		if err != nil {
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
	}
}
