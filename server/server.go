package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/mushcatshiro/gostatictracker/assets"
	"github.com/mushcatshiro/gostatictracker/blog"
	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/render"
	"github.com/spf13/afero"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Server struct {
	config            Config
	router            *http.ServeMux
	db                *dbop.DB
	googleOauthConfig *oauth2.Config
	blogManager       *blog.BlogManager
	renderEngine      *render.RenderEngine
}

func New(config Config) (*Server, error) {
	db, err := dbop.NewDB(
		config.DB.DbType,
		config.DB.User,
		config.DB.Password,
		config.DB.Host,
		config.DB.DbName,
		config.DB.SslMode,
	)
	if err != nil {
		return nil, err
	}
	err = db.InitDB(false, false)
	if err != nil {
		return nil, err
	}
	oauthConfig := &oauth2.Config{
		RedirectURL:  config.GoogleOauthConfig.RedirectURL,
		ClientID:     config.GoogleOauthConfig.ClientID,
		ClientSecret: config.GoogleOauthConfig.ClientSecret,
		Scopes:       config.GoogleOauthConfig.Scopes,
		Endpoint:     google.Endpoint,
	}

	re, err := render.NewRenderEngine(assets.TemplateFS)
	if err != nil {
		return nil, err
	}

	fs := afero.NewReadOnlyFs(afero.NewOsFs())
	bsm, err := blog.NewBlogManager(fs, config.BlogConfig.BlogRoot)
	if err != nil {
		return nil, err
	}

	s := &Server{
		config:            config,
		router:            http.NewServeMux(),
		db:                db,
		googleOauthConfig: oauthConfig,
		blogManager:       bsm,
		renderEngine:      re,
	}

	s.RegisterRoutes()
	return s, nil
}

func (s *Server) Start() {
	addr := ":" + s.config.Server.Port
	slog.Info("Starting server on", "info", addr)
	if err := http.ListenAndServe(addr, s.router); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
