package server

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func cleanDuplicatedPath(path, rawQuery string) (string, bool) {
	// TODO currently only handles duplicated, it should handled any repetition
	if strings.Contains(path, "/editor/editor/") {
		cleanedPath := strings.ReplaceAll(path, "/editor/editor/", "/editor/")
		if rawQuery != "" {
			cleanedPath += "?" + rawQuery
		}
		return cleanedPath, true
	}
	return path, false
}

func CleanDuplicatedPathsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		rawQuery := r.URL.RawQuery
		if cleanedPath, changed := cleanDuplicatedPath(path, rawQuery); changed {
			http.Redirect(w, r, cleanedPath, http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) setCORSHeaders(w http.ResponseWriter) {
	// Restricting allowed origin to your configured domain is critical for cookie-based auth[cite: 38, 40].
	w.Header().Set("Access-Control-Allow-Origin", s.config.Server.Domain)
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
}

func (s *Server) CORSPreflightHandler(w http.ResponseWriter, r *http.Request) {
	s.setCORSHeaders(w)
	w.WriteHeader(http.StatusNoContent) // 204 is the standard success status for OPTIONS
}

func (s *Server) CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.setCORSHeaders(w)
		next.ServeHTTP(w, r)
	}
}

func (s *Server) recoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// recover() catches the panic and stops it from bubbling up
			if err := recover(); err != nil {
				// Log the error and the exact line of code that caused it
				slog.Info("PANIC RECOVERED: %v\n%s", err, debug.Stack())

				// Safely return a 500 error to the client instead of dropping the connection
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (s *Server) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		next.ServeHTTP(recorder, r)
		duration := time.Since(start)
		slog.Info("[%s] %s %s - %d %v", r.RemoteAddr, r.Method, r.URL.Path, recorder.status, duration)
	}
}

func (s *Server) securityHeadersMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		// Allowing own domain + scripts from a specific CDN
		// w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' https://cdn.jsdelivr.net;")

		next.ServeHTTP(w, r)
	}
}
