package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/models"
	"github.com/mushcatshiro/gostatictracker/render"
)

type bookmarkletPayload struct {
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
		_, err := dbop.InsertEvent(s.db, b.ToEvent())
		if err != nil {
			errMsg := fmt.Sprintf("failed to insert with error: %v", err)
			http.Error(w, errMsg, http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
	}
}

func (s *Server) renderBookmarkletView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "unexpected request method", http.StatusMethodNotAllowed)
	}
	page, err := render.RenderBookmarklet(s.db)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Failed to render bookmarklet page: %v\n", err)
		fmt.Fprintf(w, "error")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, page)
}
