package cli

import (
	"database/sql"

	"github.com/spf13/cobra"
)

type App struct {
	RootCmd *cobra.Command
	DB      *sql.DB
}
