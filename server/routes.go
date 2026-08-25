package server

import "net/http"

func (s *Server) RegisterRoutes() {
	publicStack := []Middleware{
		s.recoveryMiddleware,
		s.loggingMiddleware,
		s.securityHeadersMiddleware,
		s.CORSMiddleware,
		// to add rate limiter perhaps?
	}
	securedStack := append(publicStack, s.authMiddleware)
	apiMiddlewares := []Middleware{
		s.loggingMiddleware,
		s.CORSMiddleware,
		s.authMiddleware,
	}

	s.router.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	s.router.HandleFunc("GET /.well-known/appspecific/com.chrome.devtools.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	//s.router.HandleFunc("GET /search", s.authMiddleware(s.handleSearch()))
	//s.router.HandleFunc("GET /searchRedirect", s.authMiddleware(s.handleSearchRedirect()))
	//s.router.HandleFunc("GET /eventForm", s.authMiddleware(s.handleEventFormView()))
	//s.router.HandleFunc("GET /gantt", s.renderGanttView)
	// s.router.HandleFunc("/api/error", s.authMiddleware(s.handleApiError()))
	//s.router.HandleFunc("GET /kanban", s.authMiddleware(s.renderKanbanView()))
	// s.router.HandleFunc("/roadmap", s.authMiddleware(s.renderRoadmapView()))
	//s.router.HandleFunc("GET /editor", s.authMiddleware(s.renderEditorView()))
	//s.router.HandleFunc("GET /api/assetUpload", s.authMiddleware(s.handleAssetUpload()))
	// s.router.HandleFunc("GET /siderepo/{path...}", s.authMiddleware(s.renderSideRepoView()))

	s.router.HandleFunc("GET /static/", Chain(s.handleStatic(), publicStack...))

	s.router.HandleFunc("GET /", Chain(s.handleIndex(), securedStack...))

	s.router.HandleFunc("GET /auth/google/callback", s.handleGoogleCallback())

	s.router.HandleFunc("GET /blog/list", Chain(s.renderBlogIndexView(), securedStack...))
	s.router.HandleFunc("GET /blog/{title...}", s.authMiddleware(s.renderBlogView()))

	s.router.HandleFunc("GET /bookmarklet", s.renderBookmarkletView())
	s.router.HandleFunc("GET /api/bookmarklet", Chain(s.handleInsertBookmarklet(), apiMiddlewares...))

	s.router.HandleFunc("GET /editor", Chain(s.renderEditorView(), securedStack...))
	s.router.HandleFunc("GET /editor/{title...}", Chain(s.renderEditorView(), securedStack...))

	s.router.HandleFunc("GET /error", s.handleError())
}
