package cli

import (
	"fmt"
	"strings"

	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/render"
	"github.com/spf13/cobra"
)

func renderKanbanCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rk",
		Short: "Render a Kanban",
		RunE: func(cmd *cobra.Command, args []string) error {
			// current iteration kanban only cares about status
			targetDirectory, _ := cmd.Flags().GetString("targetDir")
			fname, _ := cmd.Flags().GetString("fname")

			if !strings.HasSuffix(fname, ".html") {
				fname += ".html"
			}

			if err := common.ValidateInputDirectory(targetDirectory); err != nil {
				return fmt.Errorf("failed to validate input directory: %v", err)
			}

			htmlString, err := render.RenderKanban(app.DB)
			if err != nil {
				return fmt.Errorf("failed to render kanban: %v", err)
			}

			if err := common.PersistToFileSystem(targetDirectory, fname, htmlString); err != nil {
				return fmt.Errorf("Failed to save file to target directory: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().StringP("targetDir", "d", ".", "Target directory to write generated html file")
	cmd.Flags().StringP("fname", "f", "kanban.html", "File name of generated html file")

	return cmd
}

func renderFormCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rf",
		Short: "Render a form",
		RunE: func(cmd *cobra.Command, args []string) error {
			targetDirectory, _ := cmd.Flags().GetString("targetDir")
			fname, _ := cmd.Flags().GetString("fname")
			id, _ := cmd.Flags().GetInt64("id")

			if !strings.HasSuffix(fname, ".html") {
				fname += ".html"
			}

			if err := common.ValidateInputDirectory(targetDirectory); err != nil {
				return fmt.Errorf("failed to validate input directory: %v", err)
			}

			m, err := dbop.ReadEventById(app.DB, id)
			if err != nil {
				return fmt.Errorf("id %d is not found with error: %v", id, err)
			}
			htmlString, err := render.RenderFormHtml(m, false, "/eventForm")
			if err != nil {
				return fmt.Errorf("failed to render form: %v", err)
			}

			if err := common.PersistToFileSystem(targetDirectory, fname, htmlString); err != nil {
				return fmt.Errorf("Failed to save file to target directory: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().StringP("targetDir", "d", ".", "Target directory to write generated html file")
	cmd.Flags().StringP("fname", "f", "form.html", "File name of generated html file")
	cmd.Flags().Int64P("id", "i", 1, "id of event to be rendered")

	return cmd
}
