package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/models"
)

type bookmarkletPayload struct {
	Token       string `json:"token"`
	Title       string `json:"title"`
	Description string `json:"desc"`
	URL         string `json:"url"`
}

func (bp *bookmarkletPayload) ToBookmarklet() models.Bookmarklet {
	return models.Bookmarklet{
		Title:       bp.Title,
		Description: bp.Description,
		URL:         bp.URL,
	}
}

func (s *Server) handleInsertBookmarklet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "unexpected request method", http.StatusMethodNotAllowed)
		}
		var bp bookmarkletPayload
		if err := json.NewDecoder(r.Body).Decode(&bp); err != nil {
			http.Error(w, "incorrect payload format", http.StatusBadRequest)
		}

		b := bp.ToBookmarklet()
		eid, err := dbop.InsertEvent(s.db, b.ToEvent())
		if err != nil {
			errMsg := fmt.Sprintf("failed to insert with error: %", err)
			http.Error(w, errMsg, http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"eid": %d}`, eid)
		}
	}
}

/*
func renderBookmarklet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "unexpected request method", http.StatusMethodNotAllowed)
	}
	page, err := render.RenderBookmarklet(conn)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, page)
}
*/
