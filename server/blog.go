package server

import (
	"cmp"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"slices"

	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/markup"
)

func (s *Server) processBlogView(path string) (BlogPageTmplMeta, error) {
	var b BlogPageTmplMeta

	page, err := s.blogSiteManager.GetSpecificPageByPath(path, true)
	if err != nil {
		return b, err
	}
	ctx := markup.RenderContext{
		FilePath: page.Path,
	}
	bContent, err := s.muConverter.Convert(page.Content, ctx)
	if err != nil {
		return b, fmt.Errorf("fail to convert %s with error %w", page.Path, err)
	}

	b.InnerText = template.HTML(bContent)
	b.HasMermaid = page.Meta.HasMermaid
	b.HasMathJax = page.Meta.HasMathJax
	return b, nil
}

func (s *Server) renderBlogView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.PathValue("title")

		b, err := s.processBlogView(p)
		if err != nil {
			s.handleError(w, r, 404, err.Error())
			return
		}

		if r.Header.Get("HX-Request") == "true" {
			s.tmpl.ExecuteTemplate(w, "blog", b)
			return
		}

		isAuth := getIsAuth(r)
		err = s.tmpl.ExecuteTemplate(w, "base", BaseTmplMeta{
			SiteName:         "Mushcat`Shiro's Fortress of Solitude",
			IsAuth:           isAuth,
			IsBlog:           true,
			BlogPageTmplMeta: &b,
		})
		if err != nil {
			slog.Info("fail to execute blog base tmpl: %v", "info", err)
		}
	}
}

func (s *Server) processBlogIndexView(isAuth bool) []tmplBlogEntryMeta {
	tbeml := make([]tmplBlogEntryMeta, 0, len(s.blogSiteManager.Pages))
	s.blogSiteManager.Mu.RLock()
	for _, p := range s.blogSiteManager.Pages {
		if s.config.Server.Protected && !isAuth && p.Meta.IsPrivate {
			continue
		}
		tbpm := tmplBlogEntryMeta{
			Title:      p.Meta.Title,
			HasMermaid: p.Meta.HasMermaid,
			HasMathJax: p.Meta.HasMathJax,
			IsDraft:    p.Meta.IsDraft,
			IsPrivate:  p.Meta.IsPrivate,
			Tags:       p.Meta.Tags,
			URL:        p.FullURL,
			EditURL:    p.EditUrl,
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
		tbeml = append(tbeml, tbpm)
	}
	s.blogSiteManager.Mu.RUnlock()
	slices.SortFunc(tbeml, func(a, b tmplBlogEntryMeta) int {
		return cmp.Compare(b.CreateDate, a.CreateDate)
	})
	return tbeml
}

func (s *Server) renderBlogIndexView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAuth := getIsAuth(r)
		data := struct {
			IsAuth  bool
			Entries []tmplBlogEntryMeta
		}{
			IsAuth:  isAuth,
			Entries: s.processBlogIndexView(isAuth),
		}
		err := s.tmpl.ExecuteTemplate(w, "bloglist", data)
		if err != nil {
			slog.Info("fail to execute bloglist tmpl: %v", "info", err)
		}
	}
}
