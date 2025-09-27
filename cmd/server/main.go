package main

import (
	"log"

	"github.com/mushcatshiro/gostatictracker/server"
)

func main() {
	cfg, err := server.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create new server: %v", err)
	}
	srv.Start()
}
