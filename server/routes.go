package server

import "net/http"

func (s *Server) RegisterRoutes() {
	s.RegisterFileServer(s.router)
	s.router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	s.router.HandleFunc("/.well-known/appspecific/com.chrome.devtools.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	s.router.HandleFunc("/search", s.authMiddleware(s.handleSearch()))
	s.router.HandleFunc("/searchRedirect", s.authMiddleware(s.handleSearchRedirect()))
	s.router.HandleFunc("/eventForm", s.authMiddleware(s.handleEventFormView()))
	s.router.HandleFunc("/api/bookmarklet", s.authMiddleware(s.handleInsertBookmarklet()))
	// s.router.HandleFunc("/api/error", s.authMiddleware(s.handleApiError()))
	s.router.HandleFunc("/bookmarkletsetup", s.authMiddleware(s.renderBookmarkletSetup()))
	s.router.HandleFunc("/kanban", s.authMiddleware(s.renderKanbanView()))
	// s.router.HandleFunc("/roadmap", s.authMiddleware(s.renderRoadmapView()))
	s.router.HandleFunc("/blog/{title...}", s.authMiddleware(s.renderBlogView()))
	s.router.HandleFunc("/list/blog", s.authMiddleware(s.renderBlogIndexView()))
	s.router.HandleFunc("/editor/{title...}", s.authMiddleware(s.renderEditorView()))
	s.router.HandleFunc("/editor", s.authMiddleware(s.renderEditorView()))
	s.router.HandleFunc("/api/assetUpload", s.authMiddleware(s.handleAssetUpload()))
	s.router.HandleFunc("/preview", s.authMiddleware(s.renderBlogPreviewView()))
	s.router.HandleFunc("/siderepo/{path...}", s.authMiddleware(s.renderSideRepoView()))

	// TODO refactor by using map to register routes
	s.router.HandleFunc("/", s.authMiddleware(s.handleIndex()))
	s.router.HandleFunc("/bookmarklet", s.renderBookmarkletView())
	s.router.HandleFunc("/gantt", s.renderGanttView)
	s.router.HandleFunc("/list", s.renderListView())
	s.router.HandleFunc("/auth/google/callback", s.handleGoogleCallback())
	// s.router.HandleFunc("error", s.handleError())

}
