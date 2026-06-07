package main

import (
	"log/slog"

	"github.com/mushcatshiro/gostatictracker/server"
)

func main() {
	cfg, err := server.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config: ", "error", err.Error())
	}
	srv, err := server.New(cfg)
	if err != nil {
		slog.Error("Failed to create new server: ", "error", err.Error())
	}
	srv.Start()
}
