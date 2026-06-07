package server


func VizRoutes() {
	// rule: to only include `process*` routes, `render` routes covered by
	// RegisterRoutes
	s := Server{}
	s.verifyAuth(nil)
	s.processBlogView("")
	s.processBlogIndexView(true)
	s.processEditorView("")
}
