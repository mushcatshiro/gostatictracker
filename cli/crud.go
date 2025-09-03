package main

import (
	"log"
	"strings"
	"time"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/spf13/cobra"
)

var (
	description string
	url         string
	priority    int8
	status      int8
)

var createCmd = &cobra.Command{
	Use:   "c",
	Short: "Create an entry",
	Run:   cliCreate,
}

func cliCreate(cmd *cobra.Command, args []string) {
	title := strings.Join(args, " ")
	e := dbop.Event{
		Title:       title,
		Description: description,
		InsertTime:  time.Now().Format(dbop.TimeLayout),
		URL:         url,
		Priority:    priority,
		Status:      status,
	}
	id, err := dbop.InsertEvent(conn, e)
	if err != nil {
		log.Fatalf("Error occured during event creation: %v", err)
	}
	log.Printf("Event %d created", id)
}
