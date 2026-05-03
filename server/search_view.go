package server

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mushcatshiro/gostatictracker/render"
)

func (s *Server) handleSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtCookie, err := r.Cookie("app-jwt")
		if err != nil {
			fmt.Printf("failed to get token: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		tokStr := jwtCookie.Value
		claims := &CustomClaims{}
		token, err := jwt.ParseWithClaims(tokStr, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.config.Auth.JKey), nil
		})
		if err != nil || !token.Valid {
			fmt.Printf("err: %v || token: %t", err, token.Valid)
			http.Error(w, "Unauthorized request", http.StatusUnauthorized)
			return
		}
		page, err := render.RenderIndexView(s.db, "/searchRedirect")
		if err != nil {
			fmt.Fprintf(w, "")
			return
		}
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
