package cli

import (
	"fmt"
	"log"
	"strings"

	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/models"
	"github.com/spf13/cobra"
)

func createEntryCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "c",
		Short: "Create an entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			title := strings.Join(args, " ")
			description, _ := cmd.Flags().GetString("description")
			iPriority, _ := cmd.Flags().GetString("priority")
			iStatus, _ := cmd.Flags().GetString("status")
			url, _ := cmd.Flags().GetString("url")

			status, err := common.ParseStatus(iStatus)
			if err != nil {
				return err
			}
			priority, err := common.ParsePriority(iPriority)
			if err != nil {
				return err
			}

			e := models.Event{
				Title:       title,
				Description: description,
				URL:         url,
				Priority:    priority,
				Status:      status,
			}
			id, err := dbop.InsertEvent(app.DB, e)
			if err != nil {
				return fmt.Errorf("Error occured during event creation: %v", err)
			}
			log.Printf("Event %d created", id)
			return nil
		},
	}

	cmd.Flags().StringP("description", "d", "", "Extra descriptive information")
	cmd.Flags().StringP("url", "u", "", "Reference url")
	cmd.Flags().StringP("priority", "p", "2", "Do now (0), Do later (1), Delegate (2), Eliminate (3)")
	cmd.Flags().StringP("status", "s", "0", "Not started (0), In Progress (1), Completed (2), Cancelled (3)")

	return cmd
}
