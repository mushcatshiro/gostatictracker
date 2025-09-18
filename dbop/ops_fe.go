package dbop

import (
	"database/sql"
	"fmt"

	"github.com/mushcatshiro/gostatictracker/models"
)

func GetUngrouppedEvents(db *sql.DB) ([]models.Event, error) {
	var ungrouppedEvents []models.Event
	query := `SELECT
	  "id",
	  title
	FROM events
	WHERE "group" = 'default'
	`
	rows, err := db.Query(query)
	if err != nil {
		return ungrouppedEvents, fmt.Errorf("failed to get events: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(&e.ID, &e.Title); err != nil {
			return ungrouppedEvents, fmt.Errorf("failed to scan event: %v", err)
		}
		ungrouppedEvents = append(ungrouppedEvents, e)
	}
	if err := rows.Err(); err != nil {
		return ungrouppedEvents, fmt.Errorf("error reading events: %v", err)
	}
	return ungrouppedEvents, nil
}

func BatchUpdateGroup(db *sql.DB) error {
	return nil
}
