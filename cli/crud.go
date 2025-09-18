package main

import (
	"log"
	"strings"
	"time"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/models"
	"github.com/spf13/cobra"
)

var (
	cDescription string
	cUrl         string
	cPriority    int8
	cStatus      int8
)

var createCmd = &cobra.Command{
	Use:   "c",
	Short: "Create an entry",
	Run:   cliCreate,
}

func cliCreate(cmd *cobra.Command, args []string) {
	title := strings.Join(args, " ")
	it := time.Now()
	e := models.Event{
		Title:       title,
		Description: cDescription,
		InsertTime:  &it,
		URL:         cUrl,
		Priority:    cPriority,
		Status:      cStatus,
	}
	id, err := dbop.InsertEvent(conn, e)
	if err != nil {
		log.Fatalf("Error occured during event creation: %v", err)
	}
	log.Printf("Event %d created", id)
}
