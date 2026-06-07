package server

import (
	"net/http"
	"strings"
)

func cleanDuplicatedPath(path, rawQuery string) (string, bool) {
	// TODO currently only handles duplicated, it should handled any repetition
	if strings.Contains(path, "/editor/editor/") {
		cleanedPath := strings.ReplaceAll(path, "/editor/editor/", "/editor/")
		if rawQuery != "" {
			cleanedPath += "?" + rawQuery
		}
		return cleanedPath, true
	}
	return path, false
}

func CleanDuplicatedPathsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		rawQuery := r.URL.RawQuery
		if cleanedPath, changed := cleanDuplicatedPath(path, rawQuery); changed {
			http.Redirect(w, r, cleanedPath, http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}
