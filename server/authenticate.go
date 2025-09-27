package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenRequest struct {
	Key string `json:"key"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type CustomClaims struct {
	Uid int `json:"uid"`
	jwt.RegisteredClaims
}

func (s *Server) createJWT(uid int) (string, error) {
	cc := CustomClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Minute * time.Duration(s.config.Auth.ExpDuration)),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer:   "server",
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, cc)
	return tok.SignedString([]byte(s.config.Auth.JKey))
}

func (s *Server) handleRequestToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "unexpected request method", http.StatusMethodNotAllowed)
			return
		}
		var req TokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.Key != s.config.Auth.Key {
			http.Error(w, "Failed to authenticate", http.StatusForbidden)
			return
		}
		tok, err := s.createJWT(0)
		if err != nil {
			http.Error(w, "Failed to create token", http.StatusInternalServerError)
		}
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(TokenResponse{Token: tok})
	}
}

func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
		}
		tokStr := headerParts[1]
		claims := &CustomClaims{}
		token, err := jwt.ParseWithClaims(tokStr, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.config.Auth.JKey), nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Error(w, "Token has expired", http.StatusUnauthorized)
			} else {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			}
			return
		}
		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
