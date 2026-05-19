package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	googleOauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

func (s *Server) handleGoogleCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oauthState, err := r.Cookie("oauthstate")
		if err != nil {
			//|| r.FormValue("state") != oauthState.Value
			log.Printf("invalid oauth google state; %v", err)
			http.Error(w, "Invalid State", http.StatusBadRequest)
			return
		}
		stateJson, err := base64.URLEncoding.DecodeString(oauthState.Value)
		if err != nil {
			http.Error(w, "Invalid state format", http.StatusBadRequest)
			return
		}
		var state OauthState
		if err := json.Unmarshal(stateJson, &state); err != nil {
			log.Printf("%v", err)
			http.Error(w, "Invalid state data", http.StatusBadRequest)
			return
		}

		if r.FormValue("state") != state.Nonce {
			http.Error(w, "Invalid state nonce", http.StatusBadRequest)
			return
		}

		token, err := s.googleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
		if err != nil {
			log.Printf("code exchange failed: %s\n", err.Error())
			http.Error(w, "Code exchange failed", http.StatusInternalServerError)
			return
		}
		userInfo, err := s.getUserInfo(token)
		if err != nil {
			log.Printf("failed to get user info: %s\n", err.Error())
			http.Error(w, "Failed to get user info", http.StatusInternalServerError)
			return
		}
		jwtStr, err := s.createJWT(userInfo)
		if err != nil {
			log.Printf("failed to generate JWT: %s\n", err.Error())
			http.Error(w, "Failed to generate toekn", http.StatusInternalServerError)
			return
		}
		expiration := time.Now().Add(time.Duration(s.config.Auth.ExpDuration) * time.Minute)
		http.SetCookie(w, &http.Cookie{
			Name:     "app-jwt",
			Value:    jwtStr,
			Expires:  expiration,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		})
		http.Redirect(w, r, state.ReturnURL, http.StatusTemporaryRedirect)
	}
}

func (s *Server) getUserInfo(token *oauth2.Token) (*googleOauth2.Userinfo, error) {
	ctx := context.Background()
	oauth2Service, err := googleOauth2.NewService(ctx, option.WithTokenSource(s.googleOauthConfig.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}
	return oauth2Service.Userinfo.Get().Do()
}
