package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/mock"
	"github.com/mushcatshiro/gostatictracker/render"
)

func main() {
	var outputDir = flag.String("outputDir", ".", "The targeted output directory of the gantt html(s)")
	var username = flag.String("dbuser", "pgsql", "The postgresql username")
	var password = flag.String("dbpass", "pgsql", "The postgresql password")
	var host = flag.String("dbhost", "localhost", "The postgresql host IP")
	var dbname = flag.String("dbname", "pgsql", "The posgresql database name")
	var insertMock = flag.Bool("insertMock", false, "flag to insert mock data to postgresql")
	var renderGantt = flag.Bool("renderGantt", true, "flag to render gantt html(s)")
	var debug = flag.Bool("debug", false, "flag to enter debug mode")
	flag.Parse()

	info, err := os.Stat(*outputDir)
	if errors.Is(err, os.ErrNotExist) {
		log.Fatalf("path does not exist: %s", *outputDir)
		return
	}
	if err != nil {
		log.Fatalf("error checking path: %v", err)
		return
	}
	if !info.IsDir() {
		log.Fatalf("path given is not a valid file path: %s", *outputDir)
		return
	}

	conn, err := dbop.ConnectDB(*username, *password, *host, *dbname)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer conn.Close()

	dbop.InitDB(conn)

	mockData := slices.Concat(mock.DayViewMockData[:], mock.DayViewOverflowMockData[:], mock.WeekViewMockData[:])

	if *insertMock {
		for _, event := range mockData {
			_, err := dbop.InsertEvent(conn, event.ToEvent())
			if err != nil {
				log.Fatalf("Failed to insert mock event: %v", err)
			}
		}
		log.Printf("inserted %d entries to table `events`\n", len(mockData))
	}

	if *renderGantt {
		render.RenderGantt(*outputDir, conn, *debug)
	}
}
