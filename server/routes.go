package server

import "net/http"

func (s *Server) RegisterRoutes() {
	s.RegisterFileServer(s.router)
	s.router.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	s.router.HandleFunc("GET /.well-known/appspecific/com.chrome.devtools.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	s.router.HandleFunc("GET /search", s.authMiddleware(s.handleSearch()))
	s.router.HandleFunc("GET /searchRedirect", s.authMiddleware(s.handleSearchRedirect()))
	s.router.HandleFunc("GET /eventForm", s.authMiddleware(s.handleEventFormView()))
	s.router.HandleFunc("GET /api/bookmarklet", s.authMiddleware(s.handleInsertBookmarklet()))
	// s.router.HandleFunc("/api/error", s.authMiddleware(s.handleApiError()))
	s.router.HandleFunc("GET /bookmarkletsetup", s.authMiddleware(s.renderBookmarkletSetup()))
	s.router.HandleFunc("GET /kanban", s.authMiddleware(s.renderKanbanView()))
	// s.router.HandleFunc("/roadmap", s.authMiddleware(s.renderRoadmapView()))
	s.router.HandleFunc("GET /blog/{title...}", s.authMiddleware(s.renderBlogView()))
	s.router.HandleFunc("GET /list/blog", s.authMiddleware(s.renderBlogIndexView()))
	s.router.HandleFunc("GET /editor/{title...}", s.authMiddleware(s.renderEditorView()))
	s.router.HandleFunc("GET /editor", s.authMiddleware(s.renderEditorView()))
	s.router.HandleFunc("GET /api/assetUpload", s.authMiddleware(s.handleAssetUpload()))
	s.router.HandleFunc("GET /siderepo/{path...}", s.authMiddleware(s.renderSideRepoView()))

	// TODO refactor by using map to register routes
	s.router.HandleFunc("GET /", s.authMiddleware(s.handleIndex()))
	s.router.HandleFunc("GET /bookmarklet", s.renderBookmarkletView())
	s.router.HandleFunc("GET /gantt", s.renderGanttView)
	s.router.HandleFunc("GET /list", s.renderListView())
	s.router.HandleFunc("GET /auth/google/callback", s.handleGoogleCallback())
	// s.router.HandleFunc("error", s.handleError())

}
