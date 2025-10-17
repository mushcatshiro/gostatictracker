package server

func (s *Server) RegisterRoutes() {
	s.router.HandleFunc("/", s.handleIndex())
	s.router.HandleFunc("/eventForm", s.handleEventFormView())
	s.router.HandleFunc("/bookmarklet", s.renderBookmarkletView)
	s.router.HandleFunc("/gantt", s.renderGanttView)
	s.router.HandleFunc("/list", s.renderListView())
	s.router.HandleFunc("/auth/google/callback", s.handleGoogleCallback())
	s.router.HandleFunc("/bookmarkletsetup", s.renderBookmarkletSetup())
}
