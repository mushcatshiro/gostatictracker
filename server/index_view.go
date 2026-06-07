package server

import (
	"fmt"
	"net/http"
	"strings"
)

func (s *Server) handleIndex() http.HandlerFunc {
	/*
		- visitor receives base
		- login receives an updated base
		- allow login redirect to update base
	*/

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			msg := fmt.Sprintf("Path not found: %s", r.URL.Path)
			fmt.Println(msg)
			s.handleError(w, r, 404, msg)
			return
		}
		val := r.Context().Value(authKey)
		if val == nil {
			s.handleError(w, r, 400, "ok is nil")
			return
		}
		isAuth := val.(bool)
		if s.config.Server.Protected && !isAuth {
			s.forceLogin(w, r)
		}

		s.tmpl.ExecuteTemplate(
			w, "base", BaseTmplMeta{
				SiteName: "Mushcat`Shiro's Fortress of Solitude",
				IsAuth:   isAuth,
				IsIndex:  true,
			},
		)
	}
}

func (s *Server) handleError(w http.ResponseWriter, r *http.Request, status int, msg string) {
	// TODO need redo; first condition will prevents much more redirects that it
	// should be from actually showing error banner
	isHTMX := r.Header.Get("HX-Request") == "true"
	accept := r.Header.Get("Accept")
	wantsHTML := strings.Contains(accept, "text/html")

	if !isHTMX && wantsHTML && r.URL.Path != "/" {
		// You can optionally pass the error message via a flash cookie or query param
		http.Redirect(w, r, "/", http.StatusSeeOther)
		fmt.Printf("see other for %s with status %d and msg %s", r.URL.Path, status, msg)
		return
	}

	if !wantsHTML && !isHTMX {
		w.WriteHeader(status)
		fmt.Fprintf(w, msg)
		fmt.Printf("dont want html and not htmx for %s with status %d and msg %s", r.URL.Path, status, msg)
		return
	}

	if isHTMX {
		w.Header().Set("HX-Push-Url", "/")
		w.WriteHeader(status)
		s.tmpl.ExecuteTemplate(
			w, "error-partial", ErrPageTmplMeta{ErrorMessage: msg},
		)
		fmt.Printf("is htmx for %s with status %d and msg %s", r.URL.Path, status, msg)
		return
	}

	w.WriteHeader(status)
	s.tmpl.ExecuteTemplate(w, "base",
		BaseTmplMeta{
			ShowError:       true,
			SiteName:        "Mushcat`Shiro's Fortress of Solitude",
			IsAuth:          false,
			IsIndex:         true,
			ErrPageTmplMeta: &ErrPageTmplMeta{ErrorMessage: msg},
		},
	)
	fmt.Printf("base for %s with status %d and msg %s", r.URL.Path, status, msg)
}
