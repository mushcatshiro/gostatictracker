package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mushcatshiro/gostatictracker/gvfs"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func runWalkWithPrint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rwwp",
		Short: "run walk with print",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("directory")
			if dir == "" {
				return fmt.Errorf("empty directory given")
			}

			fs := afero.NewReadOnlyFs(afero.NewOsFs())
			sm := gvfs.NewSiteManager(fs)
			err := sm.BuildIndex(context.Background(), dir)
			if err != nil {
				return fmt.Errorf("%v", err)
			}
			pgs, err := json.MarshalIndent(sm.Pages, "", "  ")
			fmt.Println(string(pgs))
			return nil
		},
	}
	cmd.Flags().StringP("directory", "d", "", "target director to print `Walk` result")

	return cmd
}
