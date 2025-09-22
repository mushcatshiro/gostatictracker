package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleRequestToken(t *testing.T) {
	testCases := []struct {
		tname          string
		expectedResp   *TokenResponse
		expectedStatus int
	}{
		{
			tname:          "Incorrect request method",
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

	testCases := []struct {
		tname          string
		authHeader     string
		expectedStatus int
		expectedResp   string
	}{
		{
			tname: "Success - Valid Token",
			authHeader: "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectedResp: `{"eid": 0}`,
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
