package cli

import (
	"fmt"
	"slices"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/mock"
	"github.com/spf13/cobra"
)

func insertMockCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cm",
		Short: "Create mock entries for testing purposes",
		RunE: func(cmd *cobra.Command, args []string) error {
			mockData := slices.Concat(
				mock.DayViewMockData[:],
				mock.DayViewOverflowMockData[:],
				mock.WeekViewMockData[:],
			)
			if err := dbop.InitDB(app.DB); err != nil {
				return fmt.Errorf("Failed to initiate database: %v", err)
			}
			for _, event := range mockData {
				_, err := dbop.InsertEvent(app.DB, event.ToEvent())
				if err != nil {
					return fmt.Errorf("Failed to insert mock event: %v", err)
				}
			}
			return nil
		},
	}

	return cmd
}
