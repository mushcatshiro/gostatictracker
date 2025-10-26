package server

import (
	"fmt"
	"net/http"

	"github.com/mushcatshiro/gostatictracker/render"
)

func (s *Server) handleIndex() http.HandlerFunc {
	if s.config.Server.Protected {

	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, render.RenderSimpleView("index page", ""))
	}
}
