package server

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mushcatshiro/gostatictracker/render"
)

func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, "index page")
	}
}

func (s *Server) handleIndexProtected() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
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
		page, err := render.RenderIndexView(s.db, s.config.Server.Domain)
		if err != nil {
			fmt.Fprintf(w, "")
			return
		}
		fmt.Fprint(w, page)
	}
}
