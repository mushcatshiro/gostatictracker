package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var conn *sql.DB

var rootCmd = &cobra.Command{
	Use:   "gst",
	Short: "Static Tracker",
	Long: `Static Tracker is an all-in-one productivity suite designed for ` +
		`those who live in the terminal. Effortlessly manage daily todos, plan ` +
		`and visualize complex projects with Gantt chart support, track ` +
		`important events on calendar, and set personal reminders to stay on ` +
		`top of your goals. It's THE tool for streamlining workflow and ` +
		`centralizing logistics in keyboard-driven interface.`,
}

func init() {
	createCmd.Flags().StringVarP(&cDescription, "description", "d", "", "Extra descriptive information")
	createCmd.Flags().StringVarP(&cUrl, "url", "u", "", "Reference url")
	createCmd.Flags().Int8VarP(&cPriority, "priority", "p", 2, "Do now (0), Do later (1), Delegate (2), Eliminate (3)")
	createCmd.Flags().Int8VarP(&cStatus, "status", "s", 0, "Not started (0), In Progress (1), Completed (2), Cancelled (3)")

	renderCalendarCmd.Flags().IntVarP(&rcMonth, "month", "m", 0, "Target month to render")
	renderCalendarCmd.Flags().IntVarP(&rcYear, "year", "y", 0, "Target year to render")
	renderCalendarCmd.Flags().StringVarP(&rcDir, "output directory", "d", ".", "Target directory to store rendered html file(s)")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(insertMockCmd)
	rootCmd.AddCommand(renderCalendarCmd)

	setupConfig()

	var err error
	conn, err = dbop.ConnectDB(
		viper.GetString("username"),
		viper.GetString("password"),
		viper.GetString("dbhost"),
		viper.GetString("dbname"),
	)
	if err != nil {
		log.Fatalf("Failed to establish connection: %v", err)
	}
}

func setupConfig() {
	viper.SetConfigName("gst")
	viper.SetConfigType("yaml")
	userProfile, _ := os.UserHomeDir()
	viper.AddConfigPath(filepath.Join(userProfile, ".config", "gostatictracker"))
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Faield to read config %v", err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
