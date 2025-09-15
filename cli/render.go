package main

import (
	"errors"
	"log"
	"os"

	"github.com/mushcatshiro/gostatictracker/render"
	"github.com/spf13/cobra"
)

var (
	rcMonth int
	rcYear  int
	rcDir   string
)

var renderCalendarCmd = &cobra.Command{
	Use:   "rc",
	Short: "Render calendar",
	Run:   cliRenderCalendar,
}

func cliRenderCalendar(cmd *cobra.Command, args []string) {
	info, err := os.Stat(rcDir)
	if errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Directory specified does not exists: %s", rcDir)
	}
	if err != nil {
		log.Fatalf("Error during directory check %v", err)
	}
	if !info.IsDir() {
		log.Fatalf("Expects a valid path instead got %s", rcDir)
	}
	render.RenderCalendar(rcMonth, rcYear, rcDir, conn)
}
