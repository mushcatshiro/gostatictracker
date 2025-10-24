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
			aMode, _ := cmd.Flags().GetBool("append")
			if !aMode {
				if err := dbop.InitDB(app.DB, false, true); err != nil {
					return fmt.Errorf("failed to initiate database: %v", err)
				}
			}
			for _, event := range mockData {
				_, err := dbop.InsertEvent(app.DB, event)
				if err != nil {
					return fmt.Errorf("failed to insert mock event: %v+", err)
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolP("append", "a", false, "append instead of reset")
	return cmd
}
