package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCliApp() *App {
	app := &App{}
	app.RootCmd = &cobra.Command{
		Use:   "gst",
		Short: "Static Tracker",
		Long: `Static Tracker is an all-in-one productivity suite designed for ` +
			`those who live in the terminal. Effortlessly manage daily todos, plan ` +
			`and visualize complex projects with Gantt chart support, track ` +
			`important events on calendar, and set personal reminders to stay on ` +
			`top of your goals. It's THE tool for streamlining workflow and ` +
			`centralizing logistics in keyboard-driven interface.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			connStr, err := dbop.GenerateConnStr(
				viper.GetString("dbtype"),
				viper.GetString("username"),
				viper.GetString("password"),
				viper.GetString("dbhost"),
				viper.GetString("dbname"),
				viper.GetString("sslmode"),
			)
			if err != nil {
				log.Fatalf("Failed to generate connection string: %v", err)
			}
			conn, err := dbop.ConnectDB(connStr)
			if err != nil {
				log.Fatalf("Failed to establish connection: %v", err)
			}
			app.DB = conn
			return nil
		},
	}

	cobra.OnInitialize(setupConfig)

	app.RootCmd.AddCommand(createEntryCmd(app))
	app.RootCmd.AddCommand(insertMockCmd(app))
	app.RootCmd.AddCommand(renderKanbanCmd(app))
	app.RootCmd.AddCommand(renderFormCmd(app))
	return app
}

func setupConfig() {
	viper.SetConfigName("gst")
	viper.SetConfigType("yaml")
	userProfile, _ := os.UserHomeDir()
	viper.AddConfigPath(filepath.Join(userProfile, ".config", "gostatictracker"))
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read config %v", err)
	}
}

func (a *App) Execute() {
	if err := a.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
