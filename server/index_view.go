package server

import (
	"fmt"
	"net/http"
	"slices"
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

		tmpl.ExecuteTemplate(
			w, "base", BaseTmplMeta{
				SiteName: "Mushcat`Shiro's Fortress of Solitude",
				IsAuth:   isAuth,
				IsIndex:  true,
			},
		)
	}
}

func (s *Server) handleError(w http.ResponseWriter, r *http.Request, status int, msg string) {
	fmt.Println(r.URL.Path)
	w.WriteHeader(status)

	accept := r.Header.Get("Accept")
	acceptFormats := strings.Split(accept, ",")
	if !slices.Contains(acceptFormats, "text/html") && r.Header.Get("HX-Request") == "" {
		fmt.Fprint(w, msg)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		tmpl.ExecuteTemplate(
			w, "error-partial", ErrPageTmplMeta{ErrorMessage: msg},
		)
		return
	}
	tmpl.ExecuteTemplate(w, "base",
		BaseTmplMeta{
			ShowError:       true,
			SiteName:        "Mushcat`Shiro's Fortress of Solitude",
			IsAuth:          false,
			IsIndex:         true,
			ErrPageTmplMeta: &ErrPageTmplMeta{ErrorMessage: msg},
		},
	)
}
