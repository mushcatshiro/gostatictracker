package dbop

import (
	"database/sql"
	"fmt"

	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/models"
)

func GetKanbanGroup(db *sql.DB, status int, groupName string) ([]models.Event, error) {
	var events []models.Event
	query := `SELECT * FROM events WHERE status = $1 AND "group" = $2`
	rows, err := db.Query(query, status, groupName)
	if err != nil {
		return events, fmt.Errorf("failed to get events: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(
			&e.ID, &e.Start, &e.End, &e.ActualStart, &e.ActualEnd, &e.InsertTime,
			&e.Group, &e.AllDay, &e.Title, &e.URL, &e.Description, &e.PID, &e.Priority,
			&e.Metadata, &e.Status,
		); err != nil {
			return events, fmt.Errorf("failed to scan event: %v", err)
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return events, fmt.Errorf("error reading events: %v", err)
	}
	return events, nil
}

func GetKanbanGroups(db *sql.DB, groupName string) (map[string][]models.Event, error) {
	var out map[string][]models.Event

	for idx := common.NOTSTARTED; idx <= common.CANCELLED; idx++ {
		titleName := idx.String()
		listOfEvent, err := GetKanbanGroup(db, int(idx), groupName)
		if err != nil {
			return out, err
		}
		out[titleName] = listOfEvent
	}
	return out, nil
}
