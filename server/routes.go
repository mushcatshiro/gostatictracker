package server

func (s *Server) RegisterRoutes() {
	s.router.HandleFunc("/", s.handleIndex())
	s.router.HandleFunc("/api/auth", s.CORSMiddleware(s.handleRequestToken()))
	s.router.HandleFunc("/api/bookmarklet", s.CORSMiddleware(s.authMiddleware(s.handleInsertBookmarklet())))
	s.router.HandleFunc("/eventForm/{id}", s.handleEventFormView())
	s.router.HandleFunc("/eventForm", s.handleEventFormView())
	s.router.HandleFunc("/bookmarklet", s.renderBookmarkletView)
    s.router.HandleFunc("/gantt", s.renderGanttView)
    s.router.HandleFunc("/list", s.renderListView)
}
