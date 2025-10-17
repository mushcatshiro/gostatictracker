package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mushcatshiro/gostatictracker/render"
)

func (s *Server) renderGanttView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "unexpected request method", http.StatusMethodNotAllowed)
	}
    groupName := "day view example"
	page, err := render.RenderGanttV2(s.db, groupName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Failed to render gantt page: %v\n", err)
		fmt.Fprintf(w, "error")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, page)
}
