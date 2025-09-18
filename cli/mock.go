package main

import (
	"log"
	"slices"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/mock"
	"github.com/spf13/cobra"
)

var insertMockCmd = &cobra.Command{
	Use:   "cm",
	Short: "Create mock entries for testing purposes",
	Run:   insertMock,
}

func insertMock(cmd *cobra.Command, args []string) {
	mockData := slices.Concat(
		mock.DayViewMockData[:],
		mock.DayViewOverflowMockData[:],
		mock.WeekViewMockData[:],
	)
	if err := dbop.InitDB(conn); err != nil {
		log.Fatalf("Failed to initiate database: %v", err)
	}
	for _, event := range mockData {
		_, err := dbop.InsertEvent(conn, event.ToEvent())
		if err != nil {
			log.Fatalf("Failed to insert mock event: %v", err)
		}
	}
}
