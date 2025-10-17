package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Server struct {
	config            Config
	router            *http.ServeMux
	db                *sql.DB
	googleOauthConfig *oauth2.Config
}

func New(config Config) (*Server, error) {
	connStr, err := dbop.GenerateConnStr(
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
	conn, err := dbop.ConnectDB(connStr)
	if err != nil {
		return nil, err
	}
	err = dbop.InitDB(conn, false, false)
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
	s := &Server{
		config:            config,
		router:            http.NewServeMux(), // instead of global
		db:                conn,
		googleOauthConfig: oauthConfig,
	}
	s.RegisterRoutes()
	return s, nil
}

func (s *Server) Start() {
	addr := ":" + s.config.Server.Port
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, s.router))
}
