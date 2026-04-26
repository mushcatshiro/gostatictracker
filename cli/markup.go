package cli

import (
	"github.com/mushcatshiro/gostatictracker/gvfs"
	"github.com/mushcatshiro/gostatictracker/markup"
	"github.com/spf13/cobra"
)

func runMarkupCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rmm",
		Short: "run markup",
		RunE: func(cmd *cobra.Command, args []string) error {
			pagePath, _ := cmd.Flags().GetString("page")

			p := gvfs.Page{Path: pagePath}
			return markup.GenerateFullHtml(p)
		},
	}
	cmd.Flags().StringP("page", "p", ".testfiles/index.md", "page path")

	return cmd
}
