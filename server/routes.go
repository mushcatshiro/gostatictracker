package server

func (s *Server) RegisterRoutes() {
	s.router.HandleFunc("/", s.handleIndex())
	s.router.HandleFunc("/bookmarklet", s.renderBookmarkletView())
	s.router.HandleFunc("/gantt", s.renderGanttView)
	s.router.HandleFunc("/list", s.renderListView())
	s.router.HandleFunc("/auth/google/callback", s.handleGoogleCallback())

	if s.config.Server.Protected {
		s.router.HandleFunc("/eventForm", s.authMiddleware(s.handleEventFormView()))
		s.router.HandleFunc("/api/bookmarklet", s.authMiddleware(s.handleInsertBookmarklet()))
		s.router.HandleFunc("/bookmarkletsetup", s.authMiddleware(s.renderBookmarkletSetup()))
		s.router.HandleFunc("/kanban", s.authMiddleware(s.renderKanbanView()))
	} else {
		s.router.HandleFunc("/eventForm", s.handleEventFormView())
		s.router.HandleFunc("/api/bookmarklet", s.handleInsertBookmarklet())
		s.router.HandleFunc("/bookmarkletsetup", s.renderBookmarkletSetup())
		s.router.HandleFunc("/kanban", s.renderKanbanView())
	}

}
