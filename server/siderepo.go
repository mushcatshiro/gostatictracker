package server

import (
	"fmt"
	"net/http"
)

func (s *Server) renderSideRepoView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(w, "")
	}
}
