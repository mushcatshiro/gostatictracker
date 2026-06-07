package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/render"
)

func (s *Server) renderGanttView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "unexpected request method", http.StatusMethodNotAllowed)
	}
	// expecting http://localhost:8081/gantt?group=day%20view%20example
	groupName := r.URL.Query().Get("group")
	events, err := dbop.GetGanttGroupEvents(s.db, groupName, true)
	if err != nil {
		s.handleError(w, r, 404, err.Error())
		return
	}
	g, err := dbop.GetGanttRenderMetadata(s.db, groupName)
	if err != nil {
		s.handleError(w, r, 404, err.Error())
		return
	}
	page, err := render.RenderGanttV2(events, g, groupName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Failed to render gantt page: %v\n", err)
		fmt.Fprintf(w, "error")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, page)
}
