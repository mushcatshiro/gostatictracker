package cli

import (
	"errors"
	"fmt"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/spf13/cobra"
)

func createUserCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cu",
		Short: "create user in database",
		RunE: func(cmd *cobra.Command, args []string) error {
			gid, _ := cmd.Flags().GetString("gid")
			email, _ := cmd.Flags().GetString("role")
			ipaddr, _ := cmd.Flags().GetString("ipaddr")
			role, _ := cmd.Flags().GetString("role")
			if gid == "" || email == "" || ipaddr == "" || role == "" {
				return errors.New("inputs must not be empty string")
			}
			if err := dbop.InitDB(app.DB, false, false); err != nil {
				return fmt.Errorf("Failed to initiate database: %v", err)
			}
			err := dbop.InsertUser(app.DB, gid, email, ipaddr, role)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String("gid", "", "Google ID of user")
	cmd.Flags().String("email", "", "Google ID of user")
	cmd.Flags().String("ipaddr", "", "Google ID of user")
	cmd.Flags().String("role", "", "Google ID of user")

	return cmd
}
