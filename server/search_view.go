package server

import (
	"fmt"
	"net/http"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/render"
)

func (s *Server) handleSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAuth := getIsAuth(r)
		if !isAuth {
			s.handleError(w, r, http.StatusUnauthorized, "Unauthorized request")
			return
		}
		groups, err := dbop.GetUniqueGroups(s.db)
		if err != nil {
			s.handleError(w, r, 404, err.Error())
			return
		}
		page := render.RenderIndexView(groups, "/searchRedirect")
		fmt.Print(w, page)
	}
}

func (s *Server) handleSearchRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderMode := r.URL.Query().Get("render") // e.g., "gantt"
		group := r.URL.Query().Get("group")
		if renderMode == "" || group == "" {
			http.Error(w, render.RenderSimpleView("Not Found", fmt.Sprintf("/%s/%s", renderMode, group)), http.StatusNotFound)
			return
		}
		newURL := fmt.Sprintf("/%s?group=%s", renderMode, group)
		http.Redirect(w, r, newURL, http.StatusFound)
	}
}
