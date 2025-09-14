package main

import "github.com/spf13/cobra"

var (
	rcMonth int
	rcYear int
)

var renderCalendarCmd = &cobra.Command{
	Use:   "rc",
	Short: "Render calendar",
	Run:   cliRenderCalendar,
}

func cliRenderCalendar(cmd *cobra.Command, args []string) {
}
