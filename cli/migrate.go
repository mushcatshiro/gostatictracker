package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/models"
	"github.com/spf13/cobra"
)

type jsonRecord struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Desc      string `json:"desc"`
	Timestamp string `json:"timestamp"`
}

func migrateCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate entries from JSON file",
		RunE: func(cmd *cobra.Command, args []string) error {
			fpath, _ := cmd.Flags().GetString("filepath")
			if err := common.ValidateInputFilePath(fpath); err != nil {
				return err
			}

			jsonFile, err := os.Open(fpath)
			if err != nil {
				return fmt.Errorf("failed to open file %v", err)
			}
			defer jsonFile.Close()
			byteVal, err := io.ReadAll(jsonFile)
			if err != nil {
				return err
			}

			var records []jsonRecord
			err = json.Unmarshal(byteVal, &records)
			if err != nil {
				return fmt.Errorf("failed to unmarshall json %v", err)
			}

			var bookmarklets []models.Bookmarklet
			for _, fr := range records {
				var insertTime *time.Time
				if fr.Timestamp != "" {
					t, err := time.Parse("2006-01-02 15:04:05.999999", fr.Timestamp)
					if err != nil {
						return fmt.Errorf("failed to process %s's timestamp", fr.Title)
					}
					insertTime = &t
				} else {
					return fmt.Errorf("entry %s have no timestamp", fr.Title)
				}
				bk := models.Bookmarklet{
					InsertTime:  insertTime,
					Title:       fr.Title,
					Description: fr.Desc,
					URL:         fr.URL,
				}
				bookmarklets = append(bookmarklets, bk)
			}

			for _, b := range bookmarklets {
				_, err = dbop.InsertEvent(app.DB, b.ToEvent())
				if err != nil {
					fmt.Printf("failed to insert %s", b.Title)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringP("filepath", "f", "", "json file path to migrate data")
	return cmd
}
