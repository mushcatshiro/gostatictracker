package server

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mushcatshiro/gostatictracker/dbop"
	googleOauth2 "google.golang.org/api/oauth2/v2"
)


type CustomClaims struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type OauthState struct {
	Nonce     string `json:"nonce"`
	ReturnURL string `json:"return_url"`
}

var ErrNeedsLogin = errors.New("needs login")
var ErrUnAuthori = errors.New("Forbidden: User not authorized")
var ErrInsufPerm = errors.New("Forbidden: Insufficient permissions")

type contextKey string

const authKey contextKey = "IsAuth"
const userIDKey contextKey = "userID"

func (s *Server) createJWT(userInfo *googleOauth2.Userinfo) (string, error) {
	cc := CustomClaims{
		UID:   userInfo.Id,
		Email: userInfo.Email,
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

func getIPAddress(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func (s *Server) checkUserAuth(ctx context.Context, claims *CustomClaims, ipaddr string) error {
	var userRole string
	query := "SELECT role FROM users WHERE google_id = $1"
	err := s.db.QueryRowContext(ctx, query, claims.UID).Scan(&userRole)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			dbop.InsertUser(s.db, claims.UID, claims.Email, ipaddr, "norole")
			return ErrUnAuthori
		} else {
			log.Printf("Database error during authorization check for user %s: %v", claims.UID, err)
			return err
		}
	}
	if userRole != "admin" {
		return ErrInsufPerm
	}
	return nil
}

// cryptographic, can be expensive
func (s *Server) getOptionalUser(r *http.Request) (*CustomClaims, bool) {
	if !s.config.Server.Protected {
		return nil, true
	}
	jwtCookie, err := r.Cookie("app-jwt")
	if err != nil {
		return nil, false
	}
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(jwtCookie.Value, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Auth.JKey), nil
	})

	if err != nil || !token.Valid {
		return nil, false
	}
	return claims, true
}

func getIsAuth(r *http.Request) bool {
	isAuth, ok := r.Context().Value(authKey).(bool)

	if !ok {
		return false
	}

	return isAuth
}

func (s *Server) forceLogin(w http.ResponseWriter, r *http.Request) {
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
}

func (s *Server) verifyAuth(r *http.Request) (context.Context, error) {
	claims, ok := s.getOptionalUser(r)
	ctx := r.Context()

	if s.config.Server.Protected {
		if !ok {
			return nil, ErrNeedsLogin
		}

		ipaddr := getIPAddress(r)
		err := s.checkUserAuth(ctx, claims, ipaddr)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, userIDKey, claims.UID)
	} else {
		ok = true
	}

	ctx = context.WithValue(ctx, authKey, ok)
	return ctx, nil
}

func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := s.verifyAuth(r)

		if err != nil {
			if errors.Is(err, ErrNeedsLogin) {
				s.forceLogin(w, r)
				return
			}

			if errors.Is(err, ErrUnAuthori) {
				http.Error(w, "Forbidden: User not authorized", http.StatusForbidden)
			} else if errors.Is(err, ErrInsufPerm) {
				http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
