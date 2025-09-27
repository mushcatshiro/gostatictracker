package server

func (s *Server) RegisterRoutes() {
	s.router.HandleFunc("/", s.handleIndex())
	s.router.HandleFunc("/api/auth", s.handleRequestToken())
	s.router.HandleFunc("/api/bookmarklet", s.authMiddleware(s.handleInsertBookmarklet()))
	s.router.HandleFunc("/bookmarklet", s.renderBookmarkletView)
}
