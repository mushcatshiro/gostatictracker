package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/models"
	"github.com/mushcatshiro/gostatictracker/render"
)

func (s *Server) handleInsertBookmarklet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "unexpected request method", http.StatusMethodNotAllowed)
			return
		}

		title := r.URL.Query().Get("title")
		description := r.URL.Query().Get("desc")
		url := r.URL.Query().Get("url")
		b := models.Bookmarklet{
			Title:       title,
			Description: description,
			URL:         url,
		}
		_, err := dbop.InsertEvent(s.db, b.ToEvent())
		if err != nil {
			errMsg := fmt.Sprintf("failed to insert with error: %v", err)
			http.Error(w, errMsg, http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, "Insert is success...")
		}
	}
}

func (s *Server) renderBookmarkletView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprint(w, page)
	}
}

func (s *Server) renderBookmarkletSetup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bkmkSetupHtml := render.RenderBookmarkletSetupHtml(s.config.Server.Domain)
		fmt.Fprint(w, bkmkSetupHtml)
	}
}
