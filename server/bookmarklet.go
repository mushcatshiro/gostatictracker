package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/models"
	"github.com/mushcatshiro/gostatictracker/render"
)

func (s *Server) handleInsertBookmarklet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "unexpected request method", http.StatusMethodNotAllowed)
			return
		}
		userID := r.Context().Value("userID").(string)
		log.Printf("uid: %s", userID)

		title := r.URL.Query().Get("title")
		description := r.URL.Query().Get("desc")
		url := r.URL.Query().Get("url")
		b := models.Bookmarklet{
			Title:       title,
			Description: description,
			URL:         url,
		}
		_, err := dbop.InsertEvent(s.db, b.ToEvent())
		if err != nil {
			errMsg := fmt.Sprintf("failed to insert with error: %v", err)
			http.Error(w, errMsg, http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
	}
}

func (s *Server) renderBookmarkletView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "unexpected request method", http.StatusMethodNotAllowed)
		}
		page, err := render.RenderBookmarklet(s.db)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Failed to render bookmarklet page: %v\n", err)
			fmt.Fprintf(w, "error")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, page)
	}
}

func (s *Server) renderBookmarkletSetup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := bkmkCode

		bookmarkletCode := fmt.Sprintf(c, s.config.Server.Domain)
		h := bkmkSetupHtml
		fmt.Fprintf(w, h, template.HTMLEscapeString(bookmarkletCode))
	}
}

const bkmkSetupHtml = `<h1>Your Bookmarklet</h1>
<p>Drag this link to your bookmarks bar:</p>
<a href="%s" style="padding: 10px 20px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 5px;">Bookmark It!</a>
<p>When you click it on any page, it will save the URL to your collection.</p>`

const bkmkCode = `javascript:void((function(){
	function getMetaValue(propName){
		const metas=document.getElementsByTagName('meta');
		for(let i=0;i<metas.length;i++){
			const metaName=metas[i].getAttribute('name')||metas[i].getAttribute('property');
			if(metaName===propName){
				return metas[i].getAttribute('content');
			}
		}
		return'';
	}
	const metaDescription=getMetaValue('og:description')||getMetaValue('description')||'';
	window.open('%s/api/bookmarklet?url='+encodeURIComponent(window.location.href)+'&title='+encodeURIComponent(document.title)+'&desc='+encodeURIComponent(metaDescription),'save-bookmark','width=500,height=300');
})());`
