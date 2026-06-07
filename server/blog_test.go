package server

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mushcatshiro/gostatictracker/gvfs"
	"github.com/mushcatshiro/gostatictracker/markup"
	"github.com/spf13/afero"
)

func prepareTestServer() *Server {
	/*
			assemble bundle /blog
		{
		  "valid-path": {
		    "Path": "\\blog\\valid-path.md",
		    "FullURL": "blog/valid-path",
		    "SideRepoPath": "",
		    "Content": null,
		    "Meta": {
		      "Title": "Valid",
		      "HasMermaid": false,
		      "HasMathJax": false,
		      "IsDraft": false,
		      "IsPrivate": false,
		      "Tags": null,
		      "LastModifiedDate": "",
		      "CreateDate": "",
		      "ContentSidx": 3
		    },
		    "EditUrl": "editor/valid-path"
		  }
		}
	*/
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "/blog/valid-path.md", []byte("+++\ntitle = \"Valid\"\n+++\nDummy content"), 0644)

	bsm := gvfs.NewSiteManager(fs)
	err := bsm.BuildIndex(context.Background(), "/blog")
	if err != nil {
		panic(err)
	}

	tmpl := template.Must(template.New("base").Parse(`Base: {{ .BlogPageTmplMeta.InnerText }}`))
	template.Must(tmpl.New("blog").Parse(`BlogFrag: {{ .InnerText }}`))
	template.Must(tmpl.New("bloglist").Parse(`List: {{ len .Entries }} entries`))

	return &Server{
		blogSiteManager: &bsm,
		muConverter:     &markup.MockConverter{ShouldFail: false},
		tmpl:            tmpl,
		config: Config{
			Server: ServerConfig{Protected: true},
		},
	}
}

func TestProcessBlogView(t *testing.T) {
	s := prepareTestServer()

	t.Run("successful page retrieval and conversion", func(t *testing.T) {
		meta, err := s.processBlogView("valid-path")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if meta.InnerText == "" {
			t.Errorf("expected InnerText to be populated")
		}
	})

	t.Run("page not found", func(t *testing.T) {
		_, err := s.processBlogView("non-existent-path")
		if err == nil {
			t.Fatal("expected an error for missing page, got nil")
		}
	})

	t.Run("conversion failure", func(t *testing.T) {
		// Force the mock converter to fail
		s.muConverter.(*markup.MockConverter).ShouldFail = true
		defer func() { s.muConverter.(*markup.MockConverter).ShouldFail = false }() // reset

		_, err := s.processBlogView("valid-path")
		if err == nil || !strings.Contains(err.Error(), "fail to convert") {
			t.Fatalf("expected conversion error, got %v", err)
		}
	})
}

func TestRenderBlogView(t *testing.T) {
	// TODO add test for editor button with/without auth
	s := prepareTestServer()
	handler := s.renderBlogView()

	t.Run("standard HTTP request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/blog", nil)
		req.SetPathValue("title", "valid-path") // Go 1.22 routing mock

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}
		// Check that the "base" template was executed
		if !strings.Contains(rr.Body.String(), "Base:") {
			t.Errorf("expected base template execution, got body: %s", rr.Body.String())
		}
	})

	t.Run("HTMX request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/blog", nil)
		req.SetPathValue("title", "valid-path")
		req.Header.Set("HX-Request", "true")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}
		// Check that the "blog" fragment template was executed
		if !strings.Contains(rr.Body.String(), "BlogFrag:") {
			t.Errorf("expected blog fragment execution, got body: %s", rr.Body.String())
		}
	})

	t.Run("404 error handling", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/blog", nil)
		req.SetPathValue("title", "missing")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		// This assumes s.handleError writes a 404 header.
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rr.Code)
		}
	})

	/*
		current test server uses mock templates, not supporting actual tmpl
		t.Run("authenticated request with editor button", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/blog", nil)
			req.SetPathValue("title", "valid-path")

			ctx := context.WithValue(req.Context(), authKey, true)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", rr.Code)
			}

			expectedStr := "Edit"
			if !strings.Contains(rr.Body.String(), expectedStr) {
				t.Errorf("expected body to contain '%s', got body: %s", expectedStr, rr.Body.String())
			}
		})
	*/
}

func TestProcessBlogIndexView(t *testing.T) {
	s := prepareTestServer()

	// Manually inject pages into the SiteManager for this test, bypassing FS overhead
	// and ensuring we have specific private/public/date variations.
	s.blogSiteManager.Pages = map[string]gvfs.Page{
		"page1": {
			Meta: gvfs.PageMeta{Title: "Public Old", IsPrivate: false, CreateDate: "2023-01-01"},
		},
		"page2": {
			Meta: gvfs.PageMeta{Title: "Private New", IsPrivate: true, CreateDate: "2024-01-01"},
		},
		"page3": {
			Meta: gvfs.PageMeta{Title: "Public Newest", IsPrivate: false, CreateDate: "2024-05-01"},
		},
	}

	t.Run("unauthenticated user sees only public posts", func(t *testing.T) {
		isAuth := false
		entries := s.processBlogIndexView(isAuth)

		if len(entries) != 2 {
			t.Fatalf("expected 2 public entries, got %d", len(entries))
		}
		for _, entry := range entries {
			if entry.IsPrivate {
				t.Errorf("unauthenticated user should not see private entry: %s", entry.Title)
			}
		}
	})

	t.Run("authenticated user sees all posts", func(t *testing.T) {
		isAuth := true
		entries := s.processBlogIndexView(isAuth)

		if len(entries) != 3 {
			t.Fatalf("expected 3 entries, got %d", len(entries))
		}
		hasPrivate := false
		for _, entry := range entries {
			if entry.IsPrivate {
				hasPrivate = true
				break
			}
		}
		if !hasPrivate {
			t.Errorf("expected authenticated user to see private posts, but none were found")
		}
	})

	t.Run("entries are correctly sorted by creation date descending", func(t *testing.T) {
		entries := s.processBlogIndexView(true)

		if len(entries) == 0 {
			t.Fatal("expected entries to sort, got 0")
		}

		for i := 0; i < len(entries)-1; i++ {
			// In Go, standard string comparison works perfectly for YYYY-MM-DD
			if entries[i].CreateDate < entries[i+1].CreateDate {
				t.Errorf("entries are not sorted descending: %s came before %s", entries[i].CreateDate, entries[i+1].CreateDate)
			}
		}
	})
}
