package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	googleOauth2 "google.golang.org/api/oauth2/v2"
)

type CustomClaims struct {
	UID string `json:"uid"`
	jwt.RegisteredClaims
}

type OauthState struct {
	Nonce     string `json:"nonce"`
	ReturnURL string `json:"return_url"`
}

func (s *Server) createJWT(userInfo *googleOauth2.Userinfo) (string, error) {
	cc := CustomClaims{
		UID: userInfo.Id,
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

func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtCookie, err := r.Cookie("app-jwt")
		if err != nil {
			b := make([]byte, 16)
			rand.Read(b)
			nonce := base64.URLEncoding.EncodeToString(b)
			state := OauthState{
				Nonce:     nonce,
				ReturnURL: r.URL.String(),
			}
			stateJson, _ := json.Marshal(state)
			expiration := time.Now().Add(time.Duration(s.config.Auth.ExpDuration) * time.Minute)
			http.SetCookie(w, &http.Cookie{
				Name:     "oauthstate",
				Value:    base64.URLEncoding.EncodeToString(stateJson),
				Expires:  expiration,
				HttpOnly: true,
				Path:     "/",
				SameSite: http.SameSiteLaxMode,
			})
			loginURL := s.googleOauthConfig.AuthCodeURL(nonce)
			http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
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
		ctx := context.WithValue(r.Context(), "userID", claims.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
