package server

func (s *Server) registerRoutes() {
	s.router.HandleFunc("/api/auth", s.handleRequestToken())
	s.router.HandleFunc("/api/bookmarklet", s.authMiddleware(s.handleInsertBookmarklet()))
}
