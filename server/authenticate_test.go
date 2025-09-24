package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestHandleRequestToken(t *testing.T) {
	testCases := []struct {
		tname          string
		expectedResp   *TokenResponse
		expectedStatus int
	}{
		{
			tname:          "Failure-Incorrect-Request-Method",
			expectedStatus: 405,
		},
	}
	s := &Server{}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/auth", nil)
			rr := httptest.NewRecorder()
			handler := s.handleRequestToken()
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tc.expectedStatus, rr.Code, "staus code mismatched")
		})
	}
}

func TestAuthMiddleWare(t *testing.T) {
	s := &Server{
		config: Config{
			Key:  "1234",
			JKey: "2234",
		},
	}

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"eid: 0"`))
	})

	validToken, err := createJWT(0, "2234")
	if err != nil {
		t.Fatalf("Failed to create valid JWT: %v", err)
	}
	expiredToken, err := createExpiredJWT(s)
	if err != nil {
		t.Fatalf("Failed to create expired JWT: %v", err)
	}

	testCases := []struct {
		tname          string
		authHeader     string
		expectedStatus int
		expectedResp   string
	}{
		{
			tname:          "Success-Valid-Token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectedResp:   `{"eid": 0}`,
		},
		{
			tname:          "Failure-Expired-Token",
			authHeader:     "Bearer " + expiredToken,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/proctedted", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			rr := httptest.NewRecorder()
			handler := s.authMiddleware(dummyHandler)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tc.expectedStatus, rr.Code, "staus code mismatched")
		})
	}
}

func createExpiredJWT(s *Server) (string, error) {
	claims := CustomClaims{
		uid: 999,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour * 1)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JKey))
}
