package server

import (
	"fmt"
	"net/http"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/render"
)

func (s *Server) renderListView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "unexpected request method", http.StatusMethodNotAllowed)
			return
		}
		groupName := r.URL.Query().Get("group")
		le, err := dbop.GetListGroupEntries(s.db, groupName)
		if err != nil {
			s.handleError(w, r, 404, err.Error())
			return
		}
		page := render.RenderList(le, groupName)
		// if err != nil {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	log.Printf("Failed to render gantt page: %v\n", err)
		// 	fmt.Fprintf(w, "error")
		// 	return
		// }
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, page)
	}
}

func (s *Server) processListIndexView(group string) {
	// default to scratchpad
}

func (s *Server) renderListIndexView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
