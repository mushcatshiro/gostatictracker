package server

import (
	"cmp"
	"fmt"
	"html/template"
	"net/http"
	"slices"

	"github.com/mushcatshiro/gostatictracker/common"
)

func (s *Server) renderBlogView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.PathValue("title")
		page, err := s.blogSiteManager.GetSpecificPageByPath(p)
		if err != nil {
			s.handleError(w, r, 404, err.Error())
			return
		}
		bContent, err := s.muConverter.Convert(page.Content, page)
		if err != nil {
			s.handleError(w, r, 404, fmt.Sprintf("fail to convert %s with error %v", page.Path, err))
			return
		}
		if r.Header.Get("HX-Request") == "true" {
			tmpl.ExecuteTemplate(w, "blog", BlogPageTmplMeta{
				InnerText:  template.HTML(bContent),
				HasMermaid: page.Meta.HasMermaid,
				HasMathJax: page.Meta.HasMathJax,
			})
			return
		}
		val := r.Context().Value(authKey)
		if val == nil {
			s.handleError(w, r, 400, "ok is nil")
			return
		}
		isAuth := val.(bool)
		tmpl.ExecuteTemplate(w, "base", BaseTmplMeta{
			SiteName: "Mushcat`Shiro's Fortress of Solitude",
			IsAuth:   isAuth,
			IsBlog:   true,
			BlogPageTmplMeta: &BlogPageTmplMeta{
				InnerText:  template.HTML(bContent),
				HasMermaid: page.Meta.HasMermaid,
				HasMathJax: page.Meta.HasMathJax,
			},
		})
	}
}

func (s *Server) renderBlogIndexView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, loggedIn := s.getOptionalUser(r)
		var ps []tmplBlogPageMeta
		s.blogSiteManager.Mu.RLock()
		for _, p := range s.blogSiteManager.Pages {
			if s.config.Server.Protected && !loggedIn && p.Meta.IsPrivate {
				continue
			}
			tbpm := tmplBlogPageMeta{
				Title:      p.Meta.Title,
				HasMermaid: p.Meta.HasMermaid,
				HasMathJax: p.Meta.HasMathJax,
				IsDraft:    p.Meta.IsDraft,
				IsPrivate:  p.Meta.IsPrivate,
				Tags:       p.Meta.Tags,
				URL:        p.FullURL,
			}
			if p.Meta.LastModifiedDate != "" {
				lmd, err := common.ParseStringTimeToDifferentFormat(
					p.Meta.LastModifiedDate, common.BlogTimeLayout, common.BlogTimeLayoutNoTtz,
				)
				if err != nil {
					lmd = ""
				}
				tbpm.LastModifiedDate = lmd
			}
			if p.Meta.CreateDate != "" {
				cd, err := common.ParseStringTimeToDifferentFormat(
					p.Meta.CreateDate, common.BlogTimeLayout, common.BlogTimeLayoutNoTtz,
				)
				if err != nil {
					cd = ""
				}
				tbpm.CreateDate = cd
			}
			ps = append(ps, tbpm)
		}
		s.blogSiteManager.Mu.RUnlock()
		slices.SortFunc(ps, func(a, b tmplBlogPageMeta) int {
			return cmp.Compare(b.CreateDate, a.CreateDate)
		})
		tmpl.ExecuteTemplate(w, "bloglist", ps)
	}
}

func (s *Server) renderMarkdownPreviewView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(w, "")
	}
}
